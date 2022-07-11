package inverter

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
	SinglePhase     bool
}

type ETRuntimeData struct {
	PV1Voltage             int   `json:"pv1_voltage"`
	PV1Current             int   `json:"pv1_current"`
	PV1Power               int   `json:"pv1_power"`
	PV2Voltage             int   `json:"pv2_voltage"`
	PV2Current             int   `json:"pv2_current"`
	PV2Power               int   `json:"pv2_power"`
	PVPower                int   `json:"pv_power"`
	PV2Mode                byte  `json:"pv2_mode"`
	PV1Mode                byte  `json:"pv1_mode"`
	OnGridL1Voltage        int   `json:"on_grid_l1_voltage"`
	OnGridL1Current        int   `json:"on_grid_l1_current"`
	OnGridL1Frequency      int   `json:"on_grid_l1_frequency"`
	OnGridL1Power          int   `json:"on_grid_l1_power"`
	OnGridL2Voltage        int   `json:"on_grid_l2_voltage"`
	OnGridL2Current        int   `json:"on_grid_l2_current"`
	OnGridL2Frequency      int   `json:"on_grid_l2_frequency"`
	OnGridL2Power          int   `json:"on_grid_l2_power"`
	OnGridL3Voltage        int   `json:"on_grid_l3_voltage"`
	OnGridL3Current        int   `json:"on_grid_l3_current"`
	OnGridL3Frequency      int   `json:"on_grid_l3_frequency"`
	OnGridL3Power          int   `json:"on_grid_l3_power"`
	GridMode               int   `json:"grid_mode"`
	TotalInverterPower     int   `json:"total_inverter_power"`
	ActivePower            int   `json:"active_power"`
	ReactivePower          int   `json:"reactive_power"`
	ApparentPower          int   `json:"apparent_power"`
	BackupL1Voltage        int   `json:"backup_l1_voltage"`
	BackupL1Current        int   `json:"backup_l1_current"`
	BackupL1Frequency      int   `json:"backup_l1_frequency"`
	LoadModeL1             int   `json:"load_mode_l1"`
	BackupL1Power          int   `json:"backup_l1_power"`
	BackupL2Voltage        int   `json:"backup_l2_voltage"`
	BackupL2Current        int   `json:"backup_l2_current"`
	BackupL2Frequency      int   `json:"backup_l2_frequency"`
	LoadModeL2             int   `json:"load_mode_l2"`
	BackupL2Power          int   `json:"backup_l2_power"`
	BackupL3Voltage        int   `json:"backup_l3_voltage"`
	BackupL3Current        int   `json:"backup_l3_current"`
	BackupL3Frequency      int   `json:"backup_l3_frequency"`
	LoadModeL3             int   `json:"load_mode_l3"`
	BackupL3Power          int   `json:"backup_l3_power"`
	LoadL1                 int   `json:"load_l1"`
	LoadL2                 int   `json:"load_l2"`
	LoadL3                 int   `json:"load_l3"`
	BackupLoad             int   `json:"backup_load"`
	Load                   int   `json:"load"`
	UPSLoad                int   `json:"ups_load"`
	TemperatureAir         int   `json:"temperature_air"`
	TemperatureModule      int   `json:"temperature_module"`
	Temperature            int   `json:"temperature"`
	FunctionBit            int   `json:"-"`
	BusVoltage             int   `json:"bus_voltage"`
	NBusVoltage            int   `json:"nbus_voltage"`
	BatteryVoltage         int   `json:"battery_voltage"`
	BatteryCurrent         int   `json:"battery_current"`
	BatteryMode            int   `json:"battery_mode"`
	WarningCode            int   `json:"warning_code"`
	SafetyCountryCode      int   `json:"safety_country_code"`
	WorkMode               int   `json:"work_mode"`
	OperationCode          int   `json:"operation_code"`
	ErrorCodes             int   `json:"-"`
	PVGenerationTotal      int   `json:"pv_generation_total"`
	PVGenerationToday      int   `json:"pv_generation_today"`
	EnergyExportTotal      int   `json:"energy_export_total"`
	EnergyExportTotalHours int   `json:"energy_export_total_hours"`
	EnergyExportToday      int   `json:"energy_export_today"`
	EnergyImportTotal      int   `json:"energy_import_total"`
	EnergyImportToday      int   `json:"energy_import_today"`
	EnergyLoadTotal        int   `json:"energy_load_total"`
	EnergyLoadDay          int   `json:"energy_load_day"`
	BatteryChargeTotal     int   `json:"battery_charge_total"`
	BatteryChargeToday     int   `json:"battery_charge_today"`
	BatteryDischargeTotal  int   `json:"battery_discharge_total"`
	BatteryDischargeToday  int   `json:"battery_discharge_today"`
	DiagStatusCode         int   `json:"-"`
	HouseConsumption       int32 `json:"house_consumption"`
}
