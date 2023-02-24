package bdev

type BdevNvmeControllerInfo struct {
	Name   string               `json:"name"`
	ctrlrs []NvmeControllerInfo `json:"ctrlrs"`
}

type NvmeControllerInfo struct {
	State  string             `json:"state"`
	Cntlid uint16             `json:"cntlid"`
	Trid   NvmeControllerTrid `json:"trid"`
	Host   NvmeControllerHost `json:"host"`
}

type NvmeControllerTrid struct {
	Trtype  string `json:"trtype"`
	Adrfam  string `json:"adrfam"`
	Traddr  string `json:"traddr"`
	Trsvcid string `json:"trsvcid"`
	Subnqn  string `json:"subnqn"`
}

type NvmeControllerHost struct {
	Nqn   string `json:"nqn"`
	Addr  string `json:"addr"`
	Svcid string `json:"svcid"`
}

type BdevNvmeAttachControllerRequest struct {
	Name   string `json:"name"`
	Trtype string `json:"trtype"`
	Traddr string `json:"traddr"`

	Subnqn    string `json:"subnqn,omitempty"`
	Trsvcid   string `json:"trsvcid,omitempty"`
	Adrfam    string `json:"adrfam,omitempty"`
	Hostaddr  string `json:"hostaddr,omitempty"`
	Hostsvcid string `json:"hostsvcid,omitempty"`
}

type BdevNvmeDetachControllerRequest struct {
	Name string `json:"name"`

	Trtype    string `json:"trtype,omitempty"`
	Traddr    string `json:"traddr,omitempty"`
	Subnqn    string `json:"subnqn,omitempty"`
	Trsvcid   string `json:"trsvcid,omitempty"`
	Adrfam    string `json:"adrfam,omitempty"`
	Hostaddr  string `json:"hostaddr,omitempty"`
	Hostsvcid string `json:"hostsvcid,omitempty"`
}

type BdevNvmeGetControllersRequest struct {
	Name string `json:"name,omitempty"`
}

//type BdevNvmeControllerHealthInfo struct {
//	ModelNumber                             string  `json:"model_number"`
//	SerialNumber                            string  `json:"serial_number"`
//	FirmwareRevision                        string  `json:"firmware_revision"`
//	Traddr                                  string  `json:"traddr"`
//	TemperatureCelsius                      uint64  `json:"temperature_celsius"`
//	AvailableSparePercentage                uint64  `json:"available_spare_percentage"`
//	AvailableSpareThresholdPercentage       uint64  `json:"available_spare_threshold_percentage"`
//	PercentageUsed                          uint64  `json:"percentage_used"`
//	DataUnitsRead                           uint128 `json:"data_units_read"`
//	DataUnitsWritten                        uint128 `json:"data_units_written"`
//	HostReadCommands                        uint128 `json:"host_read_commands"`
//	HostWriteCommands                       uint128 `json:"host_write_commands"`
//	ControllerBusyTime                      uint128 `json:"controller_busy_time"`
//	PowerCycles                             uint128 `json:"power_cycles"`
//	PowerOnHours                            uint128 `json:"power_on_hours"`
//	UnsafeShutdowns                         uint128 `json:"unsafe_shutdowns"`
//	MediaErrors                             uint128 `json:"media_errors"`
//	NumErrLogEntries                        uint128 `json:"num_err_log_entries"`
//	WarningTemperatureTimeMinutes           uint64  `json:"warning_temperature_time_minutes"`
//	CriticalCompositeTemperatureTimeMinutes uint64  `json:"critical_composite_temperature_time_minutes"`
//}
//
//type BdevNvmeGetControllerHealthInfoRequest struct {
//	Name string `json:"name"`
//}
