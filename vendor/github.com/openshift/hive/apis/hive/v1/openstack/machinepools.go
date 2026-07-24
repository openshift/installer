package openstack

// MachinePool stores the configuration for a machine pool installed
// on OpenStack.
type MachinePool struct {
	// Flavor defines the OpenStack Nova flavor.
	// eg. m1.large
	// The json key here differs from the installer which uses both "computeFlavor" and type "type" depending on which
	// type you're looking at, and the resulting field on the MachineSet is "flavor". We are opting to stay consistent
	// with the end result.
	Flavor string `json:"flavor"`

	// RootVolume defines the root volume for instances in the machine pool.
	// The instances use ephemeral disks if not set.
	// +optional
	RootVolume *RootVolume `json:"rootVolume,omitempty"`

	// AdditionalSecurityGroupIDs contains IDs of additional security groups for machines, where each ID
	// is presented in the UUID format.
	//
	// +optional
	AdditionalSecurityGroupIDs []string `json:"additionalSecurityGroupIDs,omitempty"`
}

// Set sets the values from `required` to `a`.
func (o *MachinePool) Set(required *MachinePool) {
	if required == nil || o == nil {
		return
	}

	if required.Flavor != "" {
		o.Flavor = required.Flavor
	}

	if required.RootVolume != nil {
		if o.RootVolume == nil {
			o.RootVolume = new(RootVolume)
		}
		o.RootVolume.Size = required.RootVolume.Size
		o.RootVolume.Type = required.RootVolume.Type
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
