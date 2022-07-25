package nutanix

// MachinePool stores the configuration for a machine pool installed
// on Nutanix.
type MachinePool struct {
	// NumCPUs is the total number of virtual processor cores to assign a vm.
	//
	// +optional
	NumCPUs int64 `json:"cpus,omitempty"`

	// NumCoresPerSocket is the number of cores per socket in a vm. The number
	// of vCPUs on the vm will be NumCPUs times NumCoresPerSocket.
	// For example: 4 CPUs and 4 Cores per socket will result in 16 VPUs.
	// The AHV scheduler treats socket and core allocation exactly the same
	// so there is no benefit to configuring cores over CPUs.
	//
	// +optional
	NumCoresPerSocket int64 `json:"coresPerSocket,omitempty"`

	// Memory is the size of a VM's memory in MiB.
	//
	// +optional
	MemoryMiB int64 `json:"memoryMiB,omitempty"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk,omitempty"`
}

// OSDisk defines the disk for a virtual machine.
type OSDisk struct {
	// DiskSizeGiB defines the size of disk in GiB.
	//
	// +optional
	DiskSizeGiB int64 `json:"diskSizeGiB,omitempty"`
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

	if required.OSDisk.DiskSizeGiB != 0 {
		p.OSDisk.DiskSizeGiB = required.OSDisk.DiskSizeGiB
	}
}
