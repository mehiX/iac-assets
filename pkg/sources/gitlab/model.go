package gitlab

import (
	"fmt"
	"strings"
)

type Result struct {
	Tenant     string
	Zone       string
	CommitID   string
	Machines   []FlatStructMachine
	Aggregates AggregatedResult
	Error      string
}

type Results []Result

// Records ingores Disks for now waiting for a decision how to represent it(multiple lines, all in the same cell)
func (r Results) Records() [][]string {
	header := strings.Split("tenant,zone,v_apps_name,v_app_name,ip_last_octet,tier,cpu_count,cpu_per_socket,memory_size_gb,commitTS", ",")

	recs := make([][]string, 0)
	recs = append(recs, header)
	for _, res := range r {
		for _, m := range res.Machines {
			rec := []string{
				res.Tenant,
				res.Zone,
				m.VAppsName,
				m.VAppName,
				m.IPLastOctet,
				fmt.Sprintf("%d", m.Tier),
				fmt.Sprintf("%d", m.CpuCount),
				fmt.Sprintf("%d", m.CpuPerSocket),
				fmt.Sprintf("%d", m.MemorySizeGB),
				res.CommitID,
			}
			recs = append(recs, rec)
		}
	}

	return recs
}

type AggregatedResult struct {
	CpuCount     int
	MemorySizeGB int
}

type FlatStructMachine struct {
	IPLastOctet  string `json:"ip_last_octet"`
	Tier         int    `json:"tier"`
	CpuCount     int    `json:"cpu_count"`
	CpuPerSocket int    `json:"cpu_per_socker"`
	MemorySizeGB int    `json:"memory_size_gb"`
	Disks        []Disk `json:"disks"`
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

// ToFlatStructMachines returns machines in a flat structure, grouped by zone name
func (z Zones) ToFlatStructMachines() map[string][]FlatStructMachine {
	machines := make(map[string][]FlatStructMachine)

	for zoneName, zone := range z {
		zoneMachines := make([]FlatStructMachine, 0)
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

					zoneMachines = append(zoneMachines, newMachine)
				}
			}
		}
		machines[zoneName] = zoneMachines
	}

	return machines

}
