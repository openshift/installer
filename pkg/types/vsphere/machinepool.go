package vsphere

// MachinePool stores the configuration for a machine pool installed
// on vSphere.
type MachinePool struct {
	// NumCPUs is the total number of virtual processor cores to assign a vm.
	//
	// +optional
	NumCPUs int32 `json:"cpus"`

	// NumCoresPerSocket is the number of cores per socket in a vm. The number
	// of vCPUs on the vm will be NumCPUs/NumCoresPerSocket.
	//
	// +optional
	NumCoresPerSocket int32 `json:"coresPerSocket"`

	// Memory is the size of a VM's memory in MB.
	//
	// +optional
	MemoryMiB int64 `json:"memoryMB"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk"`
}

// OSDisk defines the disk for a virtual machine.
type OSDisk struct {
	// DiskSizeGB defines the size of disk in GB.
	//
	// +optional
	DiskSizeGB int32 `json:"diskSizeGB"`
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}

	if required.NumCPUs != 0 {
		p.NumCPUs = required.NumCPUs
	}

	if required.NumCoresPerSocket != 0 {
		p.NumCoresPerSocket = required.NumCoresPerSocket
	}

	if required.MemoryMiB != 0 {
		p.MemoryMiB = required.MemoryMiB
	}

	if required.OSDisk.DiskSizeGB != 0 {
		p.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}
}
