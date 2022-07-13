package inverter

import (
	"fmt"
	"time"

	"golang.org/x/exp/constraints"
)

type numeric interface {
	constraints.Integer | constraints.Float
}

type (
	Power     float64
	Voltage   float64
	Current   float64
	Energy    float64
	Frequency float64
	Temp      float64
)

func newPower[T numeric](v T) Power         { return Power(float64(v)) }
func newVoltage[T numeric](v T) Voltage     { return Voltage(float64(v) / 10.0) }
func newCurrent[T numeric](v T) Current     { return Current(float64(v) / 10.0) }
func newFrequency[T numeric](v T) Frequency { return Frequency(float64(v) / 100.0) }
func newTemp(v int16) Temp                  { return Temp(float64(v) / 10.0) }
func newEnergy[T numeric](v T) Energy {
	f := float64(v)
	if f == -1 {
		return 0
	}
	return Energy(f / 10.0)
}

func (v Power) String() string     { return fmt.Sprintf("%f W", v) }
func (v Voltage) String() string   { return fmt.Sprintf("%f V", v) }
func (v Current) String() string   { return fmt.Sprintf("%f A", v) }
func (v Energy) String() string    { return fmt.Sprintf("%f kWh", v) }
func (v Frequency) String() string { return fmt.Sprintf("%f Hz", v) }
func (v Temp) String() string      { return fmt.Sprintf("%f C", v) }

// DeviceInfo holds the static information about an inverter.
type DeviceInfo struct {
	ModbusVersion   int    `json:"modbus_version"`
	RatedPower      int    `json:"rated_power"`
	ACOutputType    int    `json:"ac_output_type"`
	SerialNumber    string `json:"serial_number"`
	ModelName       string `json:"model_name"`
	DSP1SWVersion   int    `json:"dsp1_sw_version"`
	DSP2SWVersion   int    `json:"dsp2_sw_version"`
	DSPSVNVersion   int    `json:"dsp_svn_version"`
	ArmSWVersion    int    `json:"arm_sw_version"`
	ArmSVNVersion   int    `json:"arm_svn_version"`
	SoftwareVersion string `json:"software_version"`
	ArmVersion      string `json:"arm_version"`
	SinglePhase     bool   `json:"single_phase"`
}

type ETRuntimeData struct {
	Timestamp              time.Time `json:"timestamp"`
	PV1Voltage             Voltage   `json:"pv1_voltage"`
	PV1Current             Current   `json:"pv1_current"`
	PV1Power               Power     `json:"pv1_power"`
	PV2Voltage             Voltage   `json:"pv2_voltage"`
	PV2Current             Current   `json:"pv2_current"`
	PV2Power               Power     `json:"pv2_power"`
	PVPower                Power     `json:"pv_power"`
	PV2Mode                byte      `json:"pv2_mode"`
	PV1Mode                byte      `json:"pv1_mode"`
	OnGridL1Voltage        Voltage   `json:"on_grid_l1_voltage"`
	OnGridL1Current        Current   `json:"on_grid_l1_current"`
	OnGridL1Frequency      Frequency `json:"on_grid_l1_frequency"`
	OnGridL1Power          Power     `json:"on_grid_l1_power"`
	OnGridL2Voltage        Voltage   `json:"on_grid_l2_voltage"`
	OnGridL2Current        Current   `json:"on_grid_l2_current"`
	OnGridL2Frequency      Frequency `json:"on_grid_l2_frequency"`
	OnGridL2Power          Power     `json:"on_grid_l2_power"`
	OnGridL3Voltage        Voltage   `json:"on_grid_l3_voltage"`
	OnGridL3Current        Current   `json:"on_grid_l3_current"`
	OnGridL3Frequency      Frequency `json:"on_grid_l3_frequency"`
	OnGridL3Power          Power     `json:"on_grid_l3_power"`
	GridMode               int       `json:"grid_mode"`
	TotalInverterPower     Power     `json:"total_inverter_power"`
	ActivePower            Power     `json:"active_power"`
	ReactivePower          int       `json:"reactive_power"`
	ApparentPower          int       `json:"apparent_power"`
	BackupL1Voltage        Voltage   `json:"backup_l1_voltage"`
	BackupL1Current        Current   `json:"backup_l1_current"`
	BackupL1Frequency      Frequency `json:"backup_l1_frequency"`
	LoadModeL1             int       `json:"load_mode_l1"`
	BackupL1Power          Power     `json:"backup_l1_power"`
	BackupL2Voltage        Voltage   `json:"backup_l2_voltage"`
	BackupL2Current        Current   `json:"backup_l2_current"`
	BackupL2Frequency      Frequency `json:"backup_l2_frequency"`
	LoadModeL2             int       `json:"load_mode_l2"`
	BackupL2Power          Power     `json:"backup_l2_power"`
	BackupL3Voltage        Voltage   `json:"backup_l3_voltage"`
	BackupL3Current        Current   `json:"backup_l3_current"`
	BackupL3Frequency      Frequency `json:"backup_l3_frequency"`
	LoadModeL3             int       `json:"load_mode_l3"`
	BackupL3Power          Power     `json:"backup_l3_power"`
	LoadL1                 Power     `json:"load_l1"`
	LoadL2                 Power     `json:"load_l2"`
	LoadL3                 Power     `json:"load_l3"`
	BackupLoad             Power     `json:"backup_load"`
	Load                   Power     `json:"load"`
	UPSLoad                int       `json:"ups_load"`
	TemperatureAir         Temp      `json:"temperature_air"`
	TemperatureModule      Temp      `json:"temperature_module"`
	Temperature            Temp      `json:"temperature"`
	FunctionBit            int       `json:"-"`
	BusVoltage             Voltage   `json:"bus_voltage"`
	NBusVoltage            Voltage   `json:"nbus_voltage"`
	BatteryVoltage         Voltage   `json:"battery_voltage"`
	BatteryCurrent         Current   `json:"battery_current"`
	BatteryMode            int       `json:"battery_mode"`
	WarningCode            int       `json:"warning_code"`
	SafetyCountryCode      int       `json:"safety_country_code"`
	WorkMode               int       `json:"work_mode"`
	OperationCode          int       `json:"operation_code"`
	ErrorCodes             int       `json:"-"`
	EnergyGenerationTotal  Energy    `json:"energy_generation_total"`
	EnergyGenerationToday  Energy    `json:"energy_generation_today"`
	EnergyExportTotal      Energy    `json:"energy_export_total"`
	EnergyExportTotalHours int       `json:"energy_export_total_hours"`
	EnergyExportToday      Energy    `json:"energy_export_today"`
	EnergyImportTotal      Energy    `json:"energy_import_total"`
	EnergyImportToday      Energy    `json:"energy_import_today"`
	EnergyLoadTotal        Energy    `json:"energy_load_total"`
	EnergyLoadDay          Energy    `json:"energy_load_day"`
	BatteryChargeTotal     int       `json:"battery_charge_total"`
	BatteryChargeToday     int       `json:"battery_charge_today"`
	BatteryDischargeTotal  int       `json:"battery_discharge_total"`
	BatteryDischargeToday  int       `json:"battery_discharge_today"`
	DiagStatusCode         int       `json:"-"`
	HouseConsumption       Power     `json:"house_consumption"`
}
