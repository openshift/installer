package openstack

// MachinePool stores the configuration for a machine pool installed
// on OpenStack.
type MachinePool struct {
	// FlavorName defines the OpenStack Nova flavor.
	// eg. m1.large
	FlavorName string `json:"type"`

	// RootVolume defines the storage for Nova instance.
	RootVolume RootVolume `json:"rootVolume"`
}

// Set sets the values from `required` to `a`.
func (o *MachinePool) Set(required *MachinePool) {
	if required == nil || o == nil {
		return
	}

	if required.FlavorName != "" {
		o.FlavorName = required.FlavorName
	}

	if required.RootVolume.IOPS != 0 {
		o.RootVolume.IOPS = required.RootVolume.IOPS
	}
	if required.RootVolume.Size != 0 {
		o.RootVolume.Size = required.RootVolume.Size
	}
	if required.RootVolume.Type != "" {
		o.RootVolume.Type = required.RootVolume.Type
	}
}

// RootVolume defines the storage for a Nova instance.
type RootVolume struct {
	// IOPS defines the iops for the instance.
	IOPS int `json:"iops"`
	// Size defines the size of the instance.
	Size int `json:"size"`
	// Type defines the type of the instance.
	Type string `json:"type"`
}
