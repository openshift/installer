package azure

// MachinePool stores the configuration for a machine pool installed
// on AWS.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the azure instance type.
	// eg. m4-large
	InstanceType string `json:"type"`
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
}
