package sources

import (
	"fmt"
	"iac-gitlab-assets/pkg/iac"
)

func (c *GitlabCollector) Collect() iac.Result {
	res := c.Query()
	if res.Error != nil {
		return iac.Result{Error: res.Error}
	}

	machines, err := zonesToIacMachines(res.Zones)

	fromName := fmt.Sprintf("%s (%s)",
		c.Name,
		string(res.CommitID[:min(8, len(res.CommitID))]))

	return iac.Result{From: fromName, Machines: machines, Error: err}
}

func zonesToIacMachines(zones Zones) ([]iac.IACMachine, error) {
	machines := make([]iac.IACMachine, 0)

	for zoneName, zone := range zones {
		for vappsName, vapps := range zone {
			for vappName, vapp := range vapps {
				for machineName, machine := range vapp {
					newMachine := iac.IACMachine{
						IPLastOctet:  machine.IPLastOctet,
						Tier:         machine.Tier,
						CpuCount:     machine.CpuCount,
						CpuPerSocket: machine.CpuPerSocket,
						MemorySizeGB: machine.MemorySizeGB,
						Disks:        toIacDisks(machine.Disks),
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

	return machines, nil

}

func toIacDisks(disks []Disk) []iac.Disk {
	iacDisks := make([]iac.Disk, len(disks))
	for i := range disks {
		iacDisks[i] = iac.Disk{
			Bus:    disks[i].Bus,
			Unit:   disks[i].Unit,
			SizeGB: disks[i].SizeGB,
		}
	}

	return iacDisks
}
