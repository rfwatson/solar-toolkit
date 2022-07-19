package store

import (
	"fmt"

	"git.netflux.io/rob/solar-toolkit/inverter"
	"github.com/jmoiron/sqlx"
)

type PostgresStore struct {
	db *sqlx.DB
}

func NewSQL(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

const insertSql = `INSERT INTO et_runtime_data (timestamp, pv1_voltage, pv1_current, pv1_power, pv2_voltage, pv2_current, pv2_power, pv_power, pv2_mode, pv1_mode, on_grid_l1_voltage, on_grid_l1_current, on_grid_l1_frequency, on_grid_l1_power, on_grid_l2_voltage, on_grid_l2_current, on_grid_l2_frequency, on_grid_l2_power, on_grid_l3_voltage, on_grid_l3_current, on_grid_l3_frequency, on_grid_l3_power, grid_mode, total_inverter_power, active_power, reactive_power, apparent_power, backup_l1_voltage, backup_l1_current, backup_l1_frequency, load_mode_l1, backup_l1_power, backup_l2_voltage, backup_l2_current, backup_l2_frequency, load_mode_l2, backup_l2_power, backup_l3_voltage, backup_l3_current, backup_l3_frequency, load_mode_l3, backup_l3_power, load_l1, load_l2, load_l3, backup_load, load, ups_load, temperature_air, temperature_module, temperature, bus_voltage, nbus_voltage, battery_voltage, battery_current, battery_mode, warning_code, safety_country_code, work_mode, operation_code, energy_generation_total, energy_generation_today, energy_export_total, energy_export_total_hours, energy_export_today, energy_import_total, energy_import_today, energy_load_total, energy_load_day, battery_charge_total, battery_charge_today, battery_discharge_total, battery_discharge_today, house_consumption, meter_test_status, meter_comm_status, active_power_l1, active_power_l2, active_power_l3, active_power_total, reactive_power_total, meter_power_factor1, meter_power_factor2, meter_power_factor3, meter_power_factor, meter_frequency, meter_energy_export_total, meter_energy_import_total, meter_active_power1, meter_active_power2, meter_active_power3, meter_active_power_total, meter_reactive_power1, meter_reactive_power2, meter_reactive_power3, meter_reactive_power_total, meter_apparent_power1, meter_apparent_power2, meter_apparent_power3, meter_apparent_power_total, meter_software_version, created_at) VALUES (:timestamp, :pv1_voltage, :pv1_current, :pv1_power, :pv2_voltage, :pv2_current, :pv2_power, :pv_power, :pv2_mode, :pv1_mode, :on_grid_l1_voltage, :on_grid_l1_current, :on_grid_l1_frequency, :on_grid_l1_power, :on_grid_l2_voltage, :on_grid_l2_current, :on_grid_l2_frequency, :on_grid_l2_power, :on_grid_l3_voltage, :on_grid_l3_current, :on_grid_l3_frequency, :on_grid_l3_power, :grid_mode, :total_inverter_power, :active_power, :reactive_power, :apparent_power, :backup_l1_voltage, :backup_l1_current, :backup_l1_frequency, :load_mode_l1, :backup_l1_power, :backup_l2_voltage, :backup_l2_current, :backup_l2_frequency, :load_mode_l2, :backup_l2_power, :backup_l3_voltage, :backup_l3_current, :backup_l3_frequency, :load_mode_l3, :backup_l3_power, :load_l1, :load_l2, :load_l3, :backup_load, :load, :ups_load, :temperature_air, :temperature_module, :temperature, :bus_voltage, :nbus_voltage, :battery_voltage, :battery_current, :battery_mode, :warning_code, :safety_country_code, :work_mode, :operation_code, :energy_generation_total, :energy_generation_today, :energy_export_total, :energy_export_total_hours, :energy_export_today, :energy_import_total, :energy_import_today, :energy_load_total, :energy_load_day, :battery_charge_total, :battery_charge_today, :battery_discharge_total, :battery_discharge_today, :house_consumption, :meter_test_status, :meter_comm_status, :active_power_l1, :active_power_l2, :active_power_l3, :active_power_total, :reactive_power_total, :meter_power_factor1, :meter_power_factor2, :meter_power_factor3, :meter_power_factor, :meter_frequency, :meter_energy_export_total, :meter_energy_import_total, :meter_active_power1, :meter_active_power2, :meter_active_power3, :meter_active_power_total, :meter_reactive_power1, :meter_reactive_power2, :meter_reactive_power3, :meter_reactive_power_total, :meter_apparent_power1, :meter_apparent_power2, :meter_apparent_power3, :meter_apparent_power_total, :meter_software_version, NOW());`

func (s *PostgresStore) InsertDataFrame(frame *inverter.ETDataFrame) error {
	if _, err := s.db.NamedExec(insertSql, frame); err != nil {
		return fmt.Errorf("error inserting data: %s", err)
	}

	return nil
}
