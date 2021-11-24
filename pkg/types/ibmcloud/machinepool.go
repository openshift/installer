package ibmcloud

// MachinePool stores the configuration for a machine pool installed on IBM Cloud.
type MachinePool struct {
	// InstanceType is the VSI machine profile.
	InstanceType string `json:"type,omitempty"`

	// Zones is the list of availability zones used for machines in the pool.
	// +optional
	Zones []string `json:"zones,omitempty"`

	// BootVolume is the configuration for the machine's boot volume.
	// +optional
	BootVolume *BootVolume `json:"bootVolume,omitempty"`

	// DedicatedHosts is the configuration for the machine's dedicated host and profile.
	// +optional
	DedicatedHosts []DedicatedHost `json:"dedicatedHosts,omitempty"`
}

// BootVolume stores the configuration for an individual machine's boot volume.
type BootVolume struct {
	// EncryptionKey is the CRN referencing a Key Protect or Hyper Protect
	// Crypto Services key to use for volume encryption. If not specified, a
	// provider managed encryption key will be used.
	// +optional
	EncryptionKey string `json:"encryptionKey,omitempty"`
}

// DedicatedHost stores the configuration for the machine's dedicated host platform.
type DedicatedHost struct {
	// Name is the name of the dedicated host to provision the machine on. If
	// specified, machines will be created on pre-existing dedicated host.
	// +optional
	Name string `json:"name,omitempty"`

	// Profile is the profile ID for the dedicated host. If specified, new
	// dedicated host will be created for machines.
	// +optional
	Profile string `json:"profile,omitempty"`
}

// Set sets the values from `required` to `a`.
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}

	if required.InstanceType != "" {
		a.InstanceType = required.InstanceType
	}

	if len(required.Zones) > 0 {
		a.Zones = required.Zones
	}

	if required.BootVolume != nil {
		if a.BootVolume == nil {
			a.BootVolume = &BootVolume{}
		}
		if required.BootVolume.EncryptionKey != "" {
			a.BootVolume.EncryptionKey = required.BootVolume.EncryptionKey
		}
	}

	if len(required.DedicatedHosts) > 0 {
		a.DedicatedHosts = required.DedicatedHosts
	}
}
