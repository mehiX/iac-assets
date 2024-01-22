package vcloud

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

type Collector struct {
	Name      string
	Endpoints []string
	Username  string
	Password  string
	Insecure  bool // skip SSL verification
	Tenants   []string
}

type Result struct {
	VirtualMachines []VM
	Time            time.Time
	Error           error
}

func (c *Collector) Query(endpoint, tenant string) Result {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return Result{Error: fmt.Errorf("unable to parse url: %w", err)}
	}

	client := govcd.NewVCDClient(*u, c.Insecure)
	err = client.Authenticate(c.Username, c.Password, tenant)
	if err != nil {
		return Result{Error: fmt.Errorf("unable to authenticate: %w", err)}
	}

	// Get the org (tenant)
	org, err := client.GetOrgByName(tenant)
	if err != nil {
		return Result{Error: err}
	}

	// Get the compute policies
	computePolicies, err := org.GetAllVdcComputePolicies(url.Values{})
	if err != nil {
		return Result{Error: err}
	}

	// Get all the VDCs in the tenants
	vdcNames, err := org.QueryOrgVdcList()
	if err != nil {
		return Result{Error: err}
	}

	vmStream := make(chan VM)

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(vdcNames))

		// Loop through the VDCs and get all the VMs
		for _, vdcName := range vdcNames {
			go func(vdcName *types.QueryResultOrgVdcRecordType) {
				defer wg.Done()
				// Get the VDC
				vdc, err := org.GetVDCByName(vdcName.Name, false)
				if err != nil {
					vmStream <- VM{Error: err}
					return
				}

				// Get all VMs in the VDC
				var filter types.VmQueryFilter
				vms, err := vdc.QueryVmList(filter)
				if err != nil {
					vmStream <- VM{Error: err}
					return
				}

				// Loop through the VMs and add them to the exportLines slice
				for _, vm := range vms {
					vmStream <- createVMFromVCDOutput(tenant, vm, vdcName.Name, computePolicies)
				}
			}(vdcName)
		}

		wg.Wait()
		close(vmStream)
	}()

	vms := make([]VM, 0)
	for vm := range vmStream {
		vms = append(vms, vm)
	}

	return Result{VirtualMachines: vms, Time: time.Now()}
}

// createVMFromVCDOutput creates a VM struct from the output of the vcd api
func createVMFromVCDOutput(tenant string, vmFromVCD *types.QueryResultVMRecordType, vdcName string, computePolicies []*govcd.VdcComputePolicy) VM {
	return VM{
		Tenant:          tenant,
		VPCName:         vdcName,
		VMName:          vmFromVCD.Name,
		CPUs:            vmFromVCD.Cpus,
		Memory:          vmFromVCD.MemoryMB,
		Storage:         fmt.Sprintf("%v", vmFromVCD.TotalStorageAllocatedMb),
		StorageProfile:  vmFromVCD.StorageProfileName,
		OS:              vmFromVCD.GuestOS,
		Status:          vmFromVCD.Status,
		CreationDate:    vmFromVCD.DateCreated,
		PlacementPolicy: GetComputePolicyName(vmFromVCD.VmPlacementPolicyId, computePolicies),
		SizingPolicy:    GetComputePolicyName(vmFromVCD.VmSizingPolicyId, computePolicies),
	}
}

// GetComputePolicyName gets the name of the compute policy based on the policy ID
func GetComputePolicyName(policyID string, computePolicies []*govcd.VdcComputePolicy) string {
	for _, policy := range computePolicies {
		if strings.Contains(policy.VdcComputePolicy.ID, policyID) {
			return policy.VdcComputePolicy.Name
		}
	}
	return ""
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
