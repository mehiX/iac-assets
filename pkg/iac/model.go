package iac

type IACMachine struct {
	IPLastOctet  string `yaml:"ip_last_octet"`
	Tier         int    `yaml:"tier"`
	CpuCount     int    `yaml:"cpu_count"`
	CpuPerSocket int    `yaml:"cpu_per_socker"`
	MemorySizeGB int    `yaml:"memory_size_gb"`
	Disks        []Disk `yaml:"disks"`
	Name         string `json:"name"`
	VAppName     string `json:"vapp"`
	VAppsName    string `json:"vapps"`
	ZoneName     string `json:"zone"`
}

type Disk struct {
	Bus    int `yaml:"bus"`
	Unit   int `yaml:"unit"`
	SizeGB int `yaml:"size_gb"`
}
