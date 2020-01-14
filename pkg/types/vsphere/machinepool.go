package vsphere

// MachinePool stores the configuration for a machine pool installed
// on vSphere.
type MachinePool struct {
	// OSDisk defines the storage for instance.
	OSDisk `json:"osDisk"`
}

// OSDisk defines the disk for a virtual machine.
type OSDisk struct {
	// DiskSizeGB defines the size of disk in GB.
	DiskSizeGB int32 `json:"diskSizeGB"`
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}

	if required.OSDisk.DiskSizeGB != 0 {
		p.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}
}
