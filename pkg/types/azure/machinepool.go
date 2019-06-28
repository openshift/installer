package azure

// MachinePool stores the configuration for a machine pool installed
// on Azure.
type MachinePool struct {
	// InstanceType defines the azure instance type.
	// eg. Standard_DS_V2
	InstanceType string `json:"type"`

	// OSDisk defines the storage for instance.
	OSDisk `json:"osDisk"`
}

// OSDisk defines the disk for machines on Azure.
type OSDisk struct {
	// DiskSizeGB defines the size of disk in GB.
	DiskSizeGB int32 `json:"diskSizeGB"`
}

// Set sets the values from `required` to `a`.
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}

	if required.InstanceType != "" {
		a.InstanceType = required.InstanceType
	}

	if required.OSDisk.DiskSizeGB != 0 {
		a.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}
}
