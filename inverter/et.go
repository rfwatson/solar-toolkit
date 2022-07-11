package inverter

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"

	"git.netflux.io/rob/goodwe-go/command"
)

type ET struct {
	SerialNumber string
	ModelName    string
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
type etRuntimeData struct {
	_                      [6]byte
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
	PVGenerationTotal      int32
	PVGenerationToday      int32
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

func (data *etRuntimeData) toRuntimeData(singlePhase bool) *ETRuntimeData {
	filterSinglePhase := func(i int) int {
		if singlePhase {
			return 0
		}
		return i
	}

	return &ETRuntimeData{
		PV1Voltage:             int(data.PV1Voltage),
		PV1Current:             int(data.PV1Current),
		PV1Power:               int(data.PV1Power),
		PV2Voltage:             int(data.PV2Voltage),
		PV2Current:             int(data.PV2Current),
		PV2Power:               int(data.PV2Power),
		PVPower:                int(data.PV1Power) + int(data.PV2Power),
		PV2Mode:                data.PV2Mode,
		PV1Mode:                data.PV1Mode,
		OnGridL1Voltage:        int(data.OnGridL1Voltage),
		OnGridL1Current:        int(data.OnGridL1Current),
		OnGridL1Frequency:      int(data.OnGridL1Frequency),
		OnGridL1Power:          int(data.OnGridL1Power),
		OnGridL2Voltage:        filterSinglePhase(int(data.OnGridL2Voltage)),
		OnGridL2Current:        filterSinglePhase(int(data.OnGridL2Current)),
		OnGridL2Frequency:      filterSinglePhase(int(data.OnGridL2Frequency)),
		OnGridL2Power:          filterSinglePhase(int(data.OnGridL2Power)),
		OnGridL3Voltage:        filterSinglePhase(int(data.OnGridL3Voltage)),
		OnGridL3Current:        filterSinglePhase(int(data.OnGridL3Current)),
		OnGridL3Frequency:      filterSinglePhase(int(data.OnGridL3Frequency)),
		OnGridL3Power:          filterSinglePhase(int(data.OnGridL3Power)),
		GridMode:               int(data.GridMode),
		TotalInverterPower:     int(data.TotalInverterPower),
		ActivePower:            int(data.ActivePower),
		ReactivePower:          int(data.ReactivePower),
		ApparentPower:          int(data.ApparentPower),
		BackupL1Voltage:        int(data.BackupL1Voltage),
		BackupL1Current:        int(data.BackupL1Current),
		BackupL1Frequency:      int(data.BackupL1Frequency),
		LoadModeL1:             int(data.LoadModeL1),
		BackupL1Power:          int(data.BackupL1Power),
		BackupL2Voltage:        filterSinglePhase(int(data.BackupL2Voltage)),
		BackupL2Current:        filterSinglePhase(int(data.BackupL2Current)),
		BackupL2Frequency:      filterSinglePhase(int(data.BackupL2Frequency)),
		LoadModeL2:             filterSinglePhase(int(data.LoadModeL2)),
		BackupL2Power:          filterSinglePhase(int(data.BackupL2Power)),
		BackupL3Voltage:        filterSinglePhase(int(data.BackupL3Voltage)),
		BackupL3Current:        filterSinglePhase(int(data.BackupL3Current)),
		BackupL3Frequency:      filterSinglePhase(int(data.BackupL3Frequency)),
		LoadModeL3:             filterSinglePhase(int(data.LoadModeL3)),
		BackupL3Power:          filterSinglePhase(int(data.BackupL3Power)),
		LoadL1:                 int(data.LoadL1),
		LoadL2:                 filterSinglePhase(int(data.LoadL2)),
		LoadL3:                 filterSinglePhase(int(data.LoadL3)),
		BackupLoad:             int(data.BackupLoad),
		Load:                   int(data.Load),
		UPSLoad:                int(data.UPSLoad),
		TemperatureAir:         int(data.TemperatureAir),
		TemperatureModule:      int(data.TemperatureModule),
		Temperature:            int(data.Temperature),
		FunctionBit:            int(data.FunctionBit),
		BusVoltage:             int(data.BusVoltage),
		NBusVoltage:            int(data.NBusVoltage),
		BatteryVoltage:         int(data.BatteryVoltage),
		BatteryCurrent:         int(data.BatteryCurrent),
		BatteryMode:            int(data.BatteryMode),
		WarningCode:            int(data.WarningCode),
		SafetyCountryCode:      int(data.SafetyCountryCode),
		WorkMode:               int(data.WorkMode),
		OperationCode:          int(data.OperationCode),
		ErrorCodes:             int(data.ErrorCodes),
		PVGenerationTotal:      int(data.PVGenerationTotal),
		PVGenerationToday:      int(data.PVGenerationToday),
		EnergyExportTotal:      int(data.EnergyExportTotal),
		EnergyExportTotalHours: int(data.EnergyExportTotalHours),
		EnergyExportToday:      int(data.EnergyExportToday),
		EnergyImportTotal:      int(data.EnergyImportTotal),
		EnergyImportToday:      int(data.EnergyImportToday),
		EnergyLoadTotal:        int(data.EnergyLoadTotal),
		EnergyLoadDay:          int(data.EnergyLoadDay),
		BatteryChargeTotal:     int(data.BatteryChargeTotal),
		BatteryChargeToday:     int(data.BatteryChargeToday),
		BatteryDischargeTotal:  int(data.BatteryDischargeTotal),
		BatteryDischargeToday:  int(data.BatteryDischargeToday),
		DiagStatusCode:         int(data.DiagStatusCode),
		HouseConsumption:       int32(float64(data.PV1Power) + float64(data.PV2Power) + math.Round(float64(data.BatteryVoltage)*float64(data.BatteryCurrent)) - float64(data.ActivePower)),
	}
}

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
