package vcloud

import (
	"fmt"
	"strings"
)

type Result struct {
	Tenant     string
	Name       string
	Endpoint   string
	Machines   []VM
	Aggregates AggregatedResult
	Error      string
}

type Results []Result

func (r Results) Records() [][]string {
	header := strings.Split("tenant,vpc_name,vm_name,cpus,memory,storage,storage_profile,os,status,placement_policy,sizing_policy,creation_date,endpoint", ",")

	recs := make([][]string, 0)
	recs = append(recs, header)
	for _, res := range r {
		for _, m := range res.Machines {
			rec := []string{
				res.Tenant,
				m.VPCName,
				m.VMName,
				fmt.Sprintf("%d", m.CPUs),
				fmt.Sprintf("%d", m.Memory),
				m.Storage,
				m.StorageProfile,
				m.OS,
				m.Status,
				m.PlacementPolicy,
				m.SizingPolicy,
				m.CreationDate,
				res.Endpoint,
			}
			recs = append(recs, rec)
		}
	}

	return recs
}

type AggregatedResult struct {
	CPUs      int
	Memory    int
	StorageMB int
}

type VM struct {
	Tenant          string // The name of the tenant this VM belongs to
	VPCName         string // The name of the VPC this VM belongs to
	VMName          string // The name (and hostname) of the VM
	CPUs            int    // The number of CPUs
	Memory          int    // The amount of memory in MB
	Storage         string // The amount of storage in MB
	StorageProfile  string // The storage profile
	OS              string // The OS
	Status          string // The status of the VM (e.g. POWERED_ON. POWERED_OFF, FAILED_CREATION)
	PlacementPolicy string // The placement policy (e.g. Redhat, Windows)
	SizingPolicy    string // The sizing policy (e.g. Extra large server)
	CreationDate    string // The creation date
	Error           error
}
