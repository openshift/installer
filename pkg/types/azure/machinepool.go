package azure

// MachinePool stores the configuration for a machine pool installed
// on Azure.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	// eg. ["1", "2", "3"]
	//
	// +optional
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the azure instance type.
	// eg. Standard_DS_V2
	//
	// +optional
	InstanceType string `json:"type"`

	// EncryptionAtHost enables encryption at the VM host.
	//
	// +optional
	EncryptionAtHost bool `json:"encryptionAtHost,omitempty"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk"`
}

// Set sets the values from `required` to `a`.
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}

	if len(required.Zones) > 0 {
		a.Zones = required.Zones
	}

	if required.InstanceType != "" {
		a.InstanceType = required.InstanceType
	}

	if required.EncryptionAtHost {
		a.EncryptionAtHost = required.EncryptionAtHost
	}

	if required.OSDisk.DiskSizeGB != 0 {
		a.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}

	if required.OSDisk.DiskType != "" {
		a.OSDisk.DiskType = required.OSDisk.DiskType
	}

	if required.DiskEncryptionSet != nil {
		a.DiskEncryptionSet = required.DiskEncryptionSet
	}
}
