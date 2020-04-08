package openstack

// MachinePool stores the configuration for a machine pool installed
// on OpenStack.
type MachinePool struct {
	// FlavorName defines the OpenStack Nova flavor.
	// eg. m1.large
	FlavorName string `json:"type"`

	// RootVolume defines the root volume for instances in the machine pool.
	// The instances use ephemeral disks if not set.
	// +optional
	RootVolume *RootVolume `json:"rootVolume,omitempty"`

	// AdditionalNetworkIDs contains IDs of additional networks for machines,
	// where each ID is presented in UUID v4 format.
	// Allowed address pairs won't be created for the additional networks.
	// +optional
	AdditionalNetworkIDs []string `json:"additionalNetworkIDs,omitempty"`

	// AdditionalSecurityGroupIDs contains IDs of additional security groups for machines,
	// where each ID is presented in UUID v4 format.
	// +optional
	AdditionalSecurityGroupIDs []string `json:"additionalSecurityGroupIDs,omitempty"`
}

// Set sets the values from `required` to `a`.
func (o *MachinePool) Set(required *MachinePool) {
	if required == nil || o == nil {
		return
	}

	if required.FlavorName != "" {
		o.FlavorName = required.FlavorName
	}

	if required.RootVolume != nil {
		if o.RootVolume == nil {
			o.RootVolume = new(RootVolume)
		}
		o.RootVolume.Size = required.RootVolume.Size
		o.RootVolume.Type = required.RootVolume.Type
	}

	if required.AdditionalNetworkIDs != nil {
		o.AdditionalNetworkIDs = append(required.AdditionalNetworkIDs[:0:0], required.AdditionalNetworkIDs...)
	}

	if required.AdditionalSecurityGroupIDs != nil {
		o.AdditionalSecurityGroupIDs = append(required.AdditionalSecurityGroupIDs[:0:0], required.AdditionalSecurityGroupIDs...)
	}
}

// RootVolume defines the storage for an instance.
type RootVolume struct {
	// Size defines the size of the volume in gibibytes (GiB).
	// Required
	Size int `json:"size"`
	// Type defines the type of the volume.
	// Required
	Type string `json:"type"`
}
