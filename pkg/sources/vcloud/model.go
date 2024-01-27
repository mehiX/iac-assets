package vcloud

type Result struct {
	Tenant     string
	Name       string
	Endpoint   string
	Machines   []VM
	Aggregates AggregatedResult
	Error      string
}

type Results []Result

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
