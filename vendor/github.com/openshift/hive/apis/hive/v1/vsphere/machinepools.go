package vsphere

// MachinePool stores the configuration for a machine pool installed
// on vSphere.
type MachinePool struct {
	// NumCPUs is the total number of virtual processor cores to assign a vm.
	NumCPUs int32 `json:"cpus"`

	// NumCoresPerSocket is the number of cores per socket in a vm. The number
	// of vCPUs on the vm will be NumCPUs/NumCoresPerSocket.
	NumCoresPerSocket int32 `json:"coresPerSocket"`

	// Memory is the size of a VM's memory in MB.
	MemoryMiB int64 `json:"memoryMB"`

	// OSDisk defines the storage for instance.
	OSDisk `json:"osDisk"`
}

// OSDisk defines the disk for a virtual machine.
type OSDisk struct {
	// DiskSizeGB defines the size of disk in GB.
	DiskSizeGB int32 `json:"diskSizeGB"`
}
