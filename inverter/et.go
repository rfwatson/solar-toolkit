package inverter

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"git.netflux.io/rob/goodwe-go/command"
)

// The timezone used to parse timestamps.
const locationName = "Europe/Madrid"

type ET struct {
	SerialNumber string
	ModelName    string
}

func (inv ET) isSinglePhase() bool {
	return strings.Contains(inv.SerialNumber, "EHU")
}

// Unexported struct used for parsing binary data only.
type etDeviceInfo struct {
	ModbusVersion   uint16
	RatedPower      uint16
	ACOutputType    uint16
	SerialNumber    [16]byte
	ModelName       [10]byte
	DSP1SWVersion   uint16
	DSP2SWVersion   uint16
	DSPSVNVersion   uint16
	ArmSWVersion    uint16
	ArmSVNVersion   uint16
	SoftwareVersion [12]byte
	ArmVersion      [12]byte
}

func (info *etDeviceInfo) toDeviceInfo() *DeviceInfo {
	serialNumber := string(info.SerialNumber[:])
	return &DeviceInfo{
		ModbusVersion:   int(info.ModbusVersion),
		RatedPower:      int(info.RatedPower),
		ACOutputType:    int(info.ACOutputType),
		SerialNumber:    serialNumber,
		ModelName:       strings.TrimSpace(string(info.ModelName[:])),
		DSP1SWVersion:   int(info.DSP1SWVersion),
		DSP2SWVersion:   int(info.DSP2SWVersion),
		DSPSVNVersion:   int(info.DSPSVNVersion),
		ArmSWVersion:    int(info.ArmSWVersion),
		ArmSVNVersion:   int(info.ArmSVNVersion),
		SoftwareVersion: string(info.SoftwareVersion[:]),
		ArmVersion:      string(info.ArmVersion[:]),
		SinglePhase:     strings.Contains(serialNumber, "EHU"),
	}
}

// Unexported struct used for parsing binary data only.
//
// Raw types are based partly on the the PyPI library, and partly on the
// third-party online documentation:
//
// https://github.com/marcelblijleven/goodwe/blob/327c7803e8415baeb4b6252431db91e1fc6f2fb3
// https://github.com/tkubec/GoodWe/wiki/ET-Series-Registers
//
// It's especially unclear whether fields should be parsed signed or unsigned.
// Handling differs in the above two sources. In most cases, overflowing a
// uint16 max value is unlikely but it may have an impact on handling negative
// values. To allow for the latter case, signed types are mostly preferred
// below.
type etRuntimeData struct {
	Timestamp              [6]byte
	PV1Voltage             int16
	PV1Current             int16
	PV1Power               int32
	PV2Voltage             int16
	PV2Current             int16
	PV2Power               int32
	_                      [18]byte
	PV2Mode                byte
	PV1Mode                byte
	OnGridL1Voltage        int16
	OnGridL1Current        int16
	OnGridL1Frequency      int16
	OnGridL1Power          int32
	OnGridL2Voltage        int16
	OnGridL2Current        int16
	OnGridL2Frequency      int16
	OnGridL2Power          int32
	OnGridL3Voltage        int16
	OnGridL3Current        int16
	OnGridL3Frequency      int16
	OnGridL3Power          int32
	GridMode               int16
	TotalInverterPower     int32
	ActivePower            int32
	ReactivePower          int32
	ApparentPower          int32
	BackupL1Voltage        int16
	BackupL1Current        int16
	BackupL1Frequency      int16
	LoadModeL1             int16
	BackupL1Power          int32
	BackupL2Voltage        int16
	BackupL2Current        int16
	BackupL2Frequency      int16
	LoadModeL2             int16
	BackupL2Power          int32
	BackupL3Voltage        int16
	BackupL3Current        int16
	BackupL3Frequency      int16
	LoadModeL3             int16
	BackupL3Power          int32
	LoadL1                 int32
	LoadL2                 int32
	LoadL3                 int32
	BackupLoad             int32
	Load                   int32
	UPSLoad                int16
	TemperatureAir         int16
	TemperatureModule      int16
	Temperature            int16
	FunctionBit            int16
	BusVoltage             int16
	NBusVoltage            int16
	BatteryVoltage         int16
	BatteryCurrent         int16
	_                      [2]byte
	BatteryMode            int32
	WarningCode            int16
	SafetyCountryCode      int16
	WorkMode               int32
	OperationCode          int16
	ErrorCodes             int16
	EnergyGenerationTotal  int32
	EnergyGenerationToday  int32
	EnergyExportTotal      int32
	EnergyExportTotalHours int32
	EnergyExportToday      int16
	EnergyImportTotal      int32
	EnergyImportToday      int16
	EnergyLoadTotal        int32
	EnergyLoadDay          int16
	BatteryChargeTotal     int32
	BatteryChargeToday     int16
	BatteryDischargeTotal  int32
	BatteryDischargeToday  int16
	_                      [16]byte
	DiagStatusCode         int32
}

func filterSinglePhase[T numeric](v T, singlePhase bool) T {
	if singlePhase {
		return 0
	}
	return v
}

// toRuntimeData panics if the `locationName` constant cannot be resolved to a
// time.Location.
func (data *etRuntimeData) toRuntimeData(singlePhase bool) *ETRuntimeData {
	yr := data.Timestamp[0]
	mon := data.Timestamp[1]
	day := data.Timestamp[2]
	hr := data.Timestamp[3]
	min := data.Timestamp[4]
	sec := data.Timestamp[5]
	loc, err := time.LoadLocation(locationName)
	if err != nil {
		panic(fmt.Sprintf("unknown location: %s", locationName))
	}

	return &ETRuntimeData{
		Timestamp:              time.Date(2000+int(yr), time.Month(mon), int(day), int(hr), int(min), int(sec), 0, loc),
		PV1Voltage:             newVoltage(data.PV1Voltage),
		PV1Current:             newCurrent(data.PV1Current),
		PV1Power:               newPower(data.PV1Power),
		PV2Voltage:             newVoltage(data.PV2Voltage),
		PV2Current:             newCurrent(data.PV2Current),
		PV2Power:               newPower(data.PV2Power),
		PVPower:                newPower(data.PV1Power + data.PV2Power),
		PV2Mode:                data.PV2Mode,
		PV1Mode:                data.PV1Mode,
		OnGridL1Voltage:        newVoltage(data.OnGridL1Voltage),
		OnGridL1Current:        newCurrent(data.OnGridL1Current),
		OnGridL1Frequency:      newFrequency(data.OnGridL1Frequency),
		OnGridL1Power:          newPower(data.OnGridL1Power),
		OnGridL2Voltage:        newVoltage(filterSinglePhase(data.OnGridL2Voltage, singlePhase)),
		OnGridL2Current:        newCurrent(filterSinglePhase(data.OnGridL2Current, singlePhase)),
		OnGridL2Frequency:      newFrequency(filterSinglePhase(data.OnGridL2Frequency, singlePhase)),
		OnGridL2Power:          newPower(filterSinglePhase(data.OnGridL2Power, singlePhase)),
		OnGridL3Voltage:        newVoltage(filterSinglePhase(data.OnGridL3Voltage, singlePhase)),
		OnGridL3Current:        newCurrent(filterSinglePhase(data.OnGridL3Current, singlePhase)),
		OnGridL3Frequency:      newFrequency(filterSinglePhase(data.OnGridL3Frequency, singlePhase)),
		OnGridL3Power:          newPower(filterSinglePhase(data.OnGridL3Power, singlePhase)),
		GridMode:               int(data.GridMode),
		TotalInverterPower:     newPower(data.TotalInverterPower),
		ActivePower:            newPower(data.ActivePower),
		ReactivePower:          int(data.ReactivePower),
		ApparentPower:          int(data.ApparentPower),
		BackupL1Voltage:        newVoltage(data.BackupL1Voltage),
		BackupL1Current:        newCurrent(data.BackupL1Current),
		BackupL1Frequency:      newFrequency(data.BackupL1Frequency),
		LoadModeL1:             int(data.LoadModeL1),
		BackupL1Power:          newPower(data.BackupL1Power),
		BackupL2Voltage:        newVoltage(filterSinglePhase(data.BackupL2Voltage, singlePhase)),
		BackupL2Current:        newCurrent(filterSinglePhase(data.BackupL2Current, singlePhase)),
		BackupL2Frequency:      newFrequency(filterSinglePhase(data.BackupL2Frequency, singlePhase)),
		LoadModeL2:             int(filterSinglePhase(data.LoadModeL2, singlePhase)),
		BackupL2Power:          newPower(filterSinglePhase(data.BackupL2Power, singlePhase)),
		BackupL3Voltage:        newVoltage(filterSinglePhase(data.BackupL3Voltage, singlePhase)),
		BackupL3Current:        newCurrent(filterSinglePhase(data.BackupL3Current, singlePhase)),
		BackupL3Frequency:      newFrequency(filterSinglePhase(data.BackupL3Frequency, singlePhase)),
		LoadModeL3:             int(filterSinglePhase(data.LoadModeL3, singlePhase)),
		BackupL3Power:          newPower(filterSinglePhase(data.BackupL3Power, singlePhase)),
		LoadL1:                 newPower(data.LoadL1),
		LoadL2:                 newPower(filterSinglePhase(data.LoadL2, singlePhase)),
		LoadL3:                 newPower(filterSinglePhase(data.LoadL3, singlePhase)),
		BackupLoad:             newPower(data.BackupLoad),
		Load:                   newPower(data.Load),
		UPSLoad:                int(data.UPSLoad),
		TemperatureAir:         newTemp(data.TemperatureAir),
		TemperatureModule:      newTemp(data.TemperatureModule),
		Temperature:            newTemp(data.Temperature),
		FunctionBit:            int(data.FunctionBit),
		BusVoltage:             newVoltage(data.BusVoltage),
		NBusVoltage:            newVoltage(data.NBusVoltage),
		BatteryVoltage:         newVoltage(data.BatteryVoltage),
		BatteryCurrent:         newCurrent(data.BatteryCurrent),
		BatteryMode:            int(data.BatteryMode),
		WarningCode:            int(data.WarningCode),
		SafetyCountryCode:      int(data.SafetyCountryCode),
		WorkMode:               int(data.WorkMode),
		OperationCode:          int(data.OperationCode),
		ErrorCodes:             int(data.ErrorCodes),
		EnergyGenerationTotal:  newEnergy(data.EnergyGenerationTotal),
		EnergyGenerationToday:  newEnergy(data.EnergyGenerationToday),
		EnergyExportTotal:      newEnergy(data.EnergyExportTotal),
		EnergyExportTotalHours: int(data.EnergyExportTotalHours),
		EnergyExportToday:      newEnergy(data.EnergyExportToday),
		EnergyImportTotal:      newEnergy(data.EnergyImportTotal),
		EnergyImportToday:      newEnergy(data.EnergyImportToday),
		EnergyLoadTotal:        newEnergy(data.EnergyLoadTotal),
		EnergyLoadDay:          newEnergy(data.EnergyLoadDay),
		BatteryChargeTotal:     int(data.BatteryChargeTotal),
		BatteryChargeToday:     int(data.BatteryChargeToday),
		BatteryDischargeTotal:  int(data.BatteryDischargeTotal),
		BatteryDischargeToday:  int(data.BatteryDischargeToday),
		DiagStatusCode:         int(data.DiagStatusCode),
		HouseConsumption:       Power(int32(float64(data.PV1Power) + float64(data.PV2Power) + math.Round(float64(data.BatteryVoltage)*float64(data.BatteryCurrent)) - float64(data.ActivePower))),
	}
}

func (inv ET) DecodeRuntimeData(p []byte) (*ETRuntimeData, error) {
	var runtimeData etRuntimeData
	if err := binary.Read(bytes.NewReader(p), binary.BigEndian, &runtimeData); err != nil {
		return nil, fmt.Errorf("error parsing response: %s", err)
	}

	return runtimeData.toRuntimeData(inv.isSinglePhase()), nil
}

// DEPRECATED
func (inv ET) DeviceInfo(ctx context.Context, conn io.ReadWriter) (*DeviceInfo, error) {
	resp, err := command.Send(command.NewModbus(command.ModbusCommandTypeRead, 0x88b8, 0x0021), conn)
	if err != nil {
		return nil, fmt.Errorf("error sending command: %s", err)
	}

	var deviceInfo etDeviceInfo
	if err := binary.Read(bytes.NewReader(resp), binary.BigEndian, &deviceInfo); err != nil {
		return nil, fmt.Errorf("error parsing response: %s", err)
	}

	return deviceInfo.toDeviceInfo(), nil
}

// DEPRECATED
func (inv ET) RuntimeData(ctx context.Context, conn io.ReadWriter) (*ETRuntimeData, error) {
	deviceInfo, err := inv.DeviceInfo(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("error fetching device info: %s", err)
	}

	resp, err := command.Send(command.NewModbus(command.ModbusCommandTypeRead, 0x891c, 0x007d), conn)
	if err != nil {
		return nil, fmt.Errorf("error sending command: %s", err)
	}

	var runtimeData etRuntimeData
	if err := binary.Read(bytes.NewReader(resp), binary.BigEndian, &runtimeData); err != nil {
		return nil, fmt.Errorf("error parsing response: %s", err)
	}

	return runtimeData.toRuntimeData(deviceInfo.SinglePhase), nil
}
