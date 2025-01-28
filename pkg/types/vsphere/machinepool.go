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

	// DataDisks are additional disks to add to the VM that are not part of the VM's OVA template.
	//
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=29
	DataDisks []DataDisk `json:"dataDisks"`

	// Zones defines available zones
	// Zones is available in TechPreview.
	//
	// +omitempty
	Zones []string `json:"zones,omitempty"`
}

// OSDisk defines the disk for a virtual machine.
type OSDisk struct {
	// DiskSizeGB defines the size of disk in GB.
	//
	// +optional
	DiskSizeGB int32 `json:"diskSizeGB"`
}

// DataDisk defines a data disk to add to the VM that is not part of the VM OVA template.
type DataDisk struct {
	// name is used to identify the disk definition. name is required needs to be unique so that it can be used to
	// clearly identify purpose of the disk.
	// +kubebuilder:example=images_1
	// +kubebuilder:validation:MaxLength=80
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9]([-_a-zA-Z0-9]*[a-zA-Z0-9])?$"
	// +required
	Name string `json:"name"`
	// sizeGiB is the size of the disk in GiB.
	// The maximum supported size is 16384 GiB.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=16384
	// +required
	SizeGiB int32 `json:"sizeGiB"`
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

	if len(required.Zones) > 0 {
		p.Zones = required.Zones
	}

	if len(required.DataDisks) > 0 {
		p.DataDisks = required.DataDisks
	}
}
