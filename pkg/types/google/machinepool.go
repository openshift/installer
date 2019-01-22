package google

// MachinePool stores the configuration for a machine pool installed
// on GCP.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	Zones []string `json:"zones,omitempty"`

	// ImageName defines the Image that should be used.
	ImageName string `json:"imageName,omitempty"`

	// InstanceType defines the GCP instance type.
	// eg. m4-large
	InstanceType string `json:"type"`

	// RootVolume defines the storage for instance.
	RootVolume `json:"rootVolume"`
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

	if required.RootVolume.Size != 0 {
		a.RootVolume.Size = required.RootVolume.Size
	}
	if required.RootVolume.Type != "" {
		a.RootVolume.Type = required.RootVolume.Type
	}
}

// RootVolume defines the storage for an instance.
type RootVolume struct {
	// Size defines the size of the storage.
	Size int `json:"size"`
	// Type defines the type of the storage.
	Type string `json:"type"`
}
