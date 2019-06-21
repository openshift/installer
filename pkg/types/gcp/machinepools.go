package gcp

// MachinePool stores the configuration for a machine pool installed on GCP.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the GCP instance type.
	// eg. n1-standard-4
	InstanceType string `json:"type"`
}
