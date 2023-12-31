ALTER TABLE et_runtime_data ADD COLUMN meter_test_status INT NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_comm_status INT NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN active_power_l1 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN active_power_l2 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN active_power_l3 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN active_power_total DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN reactive_power_total DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_power_factor1 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_power_factor2 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_power_factor3 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_power_factor DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_frequency DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_energy_export_total DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_energy_import_total DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_active_power1 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_active_power2 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_active_power3 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_active_power_total DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_reactive_power1 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_reactive_power2 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_reactive_power3 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_reactive_power_total DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_apparent_power1 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_apparent_power2 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_apparent_power3 DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_apparent_power_total DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE et_runtime_data ADD COLUMN meter_software_version INT NOT NULL DEFAULT 0;
