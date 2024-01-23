package gitlab

type FlatStructMachine struct {
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

type Machine struct {
	IPLastOctet  string `yaml:"ip_last_octet"`
	Tier         int    `yaml:"tier"`
	CpuCount     int    `yaml:"cpu_count"`
	CpuPerSocket int    `yaml:"cpu_per_socker"`
	MemorySizeGB int    `yaml:"memory_size_gb"`
	Disks        []Disk `yaml:"disks"`
}

type Disk struct {
	Bus    int `yaml:"bus"`
	Unit   int `yaml:"unit"`
	SizeGB int `yaml:"size_gb"`
}

type Vapp map[string]Machine

type Vapps map[string]Vapp

type Zone map[string]Vapps

type Zones map[string]Zone

func (z Zones) ToFlatStructMachines() []FlatStructMachine {
	machines := make([]FlatStructMachine, 0)

	for zoneName, zone := range z {
		for vappsName, vapps := range zone {
			for vappName, vapp := range vapps {
				for machineName, machine := range vapp {
					newMachine := FlatStructMachine{
						IPLastOctet:  machine.IPLastOctet,
						Tier:         machine.Tier,
						CpuCount:     machine.CpuCount,
						CpuPerSocket: machine.CpuPerSocket,
						MemorySizeGB: machine.MemorySizeGB,
						Disks:        machine.Disks,
						Name:         machineName,
						VAppName:     vappName,
						VAppsName:    vappsName,
						ZoneName:     zoneName,
					}

					machines = append(machines, newMachine)
				}
			}
		}
	}

	return machines

}