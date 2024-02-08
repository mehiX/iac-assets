package vcloud

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"sync"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

type Source struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Tenant   string `json:"tenant"`   // Org
	Endpoint string `json:"endpoint"` // Href
	Insecure bool   `json:"insecure"`
}

type Response struct {
	VirtualMachines []VM
	Tenant          string
	Endpoint        string
	Error           error
}

func (s Source) Query() Response {

	lg := slog.With("tenant", s.Tenant, "endpoint", s.Endpoint)

	lg.Info("query vcloud")

	u, err := url.ParseRequestURI(s.Endpoint)
	if err != nil {
		return Response{Error: fmt.Errorf("unable to parse url: %w", err)}
	}

	client := govcd.NewVCDClient(*u, s.Insecure, govcd.WithHttpTimeout(10))
	err = client.Authenticate(s.User, s.Password, s.Tenant)
	if err != nil {
		return Response{Error: err}
	}
	lg.Info("client authenticated")

	defer client.Disconnect()

	// Get the org (tenant)
	org, err := client.GetOrgByName(s.Tenant)
	if err != nil {
		return Response{Error: err}
	}
	lg.Info("got tenant")

	// Get the compute policies
	computePolicies, err := org.GetAllVdcComputePolicies(url.Values{})
	if err != nil {
		return Response{Error: err}
	}

	// Get all the VDCs in the tenants
	vdcNames, err := org.QueryOrgVdcList()
	if err != nil {
		return Response{Error: err}
	}
	lg.Info("got vdc list", "size", len(vdcNames))

	vmStream := make(chan VM, 10)

	var wg sync.WaitGroup
	wg.Add(len(vdcNames))

	// Loop through the VDCs and get all the VMs
	for _, vdcName := range vdcNames {
		go func(vdcName *types.QueryResultOrgVdcRecordType) {
			defer wg.Done()
			// Get the VDC
			vdc, err := org.GetVDCByName(vdcName.Name, false)
			if err != nil {
				lg.Error("getting VDC by name", "error", err.Error())
				vmStream <- VM{Tenant: s.Tenant, Error: err}
				return
			}

			// Get all VMs in the VDC
			var filter types.VmQueryFilter
			vms, err := vdc.QueryVmList(filter)
			if err != nil {
				lg.Error("getting VM list", "error", err.Error())
				vmStream <- VM{Tenant: s.Tenant, Error: err}
				return
			}

			// Loop through the VMs and add them to the exportLines slice
			for _, vm := range vms {
				vmStream <- createVMFromVCDOutput(s.Tenant, vm, vdcName.Name, computePolicies)
			}
		}(vdcName)
	}

	go func() {
		wg.Wait()
		close(vmStream)
	}()

	vms := make([]VM, 0)
	for vm := range vmStream {
		vms = append(vms, vm)
	}

	lg.Info("reponding with VMS", "size", len(vms))
	return Response{VirtualMachines: vms, Tenant: s.Tenant, Endpoint: s.Endpoint}
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
		PlacementPolicy: getComputePolicyName(vmFromVCD.VmPlacementPolicyId, computePolicies),
		SizingPolicy:    getComputePolicyName(vmFromVCD.VmSizingPolicyId, computePolicies),
	}
}

// getComputePolicyName gets the name of the compute policy based on the policy ID
func getComputePolicyName(policyID string, computePolicies []*govcd.VdcComputePolicy) string {
	for _, policy := range computePolicies {
		if strings.Contains(policy.VdcComputePolicy.ID, policyID) {
			return policy.VdcComputePolicy.Name
		}
	}
	return ""
}
