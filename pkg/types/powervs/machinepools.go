package powervs

// ProcType defines valid types for a ppc64le processor in Power VS
// +kubebuilder:validation:Enum="";capped;dedicated;shared
type ProcType string

// Capped type for capped processor consumption
const Capped ProcType = "capped"

// Dedicated type for dedicated processor(s)
const Dedicated ProcType = "dedicated"

// Shared type shared type for shared processor(s)
const Shared ProcType = "shared"

// MachinePool stores the configuration for a machine pool installed on IBM Power VS.
type MachinePool struct {
	// VolumeIDs is the list of volumes attached to the instance.
	//
	// +optional
	VolumeIDs []string `json:"volumeIDs,omitempty"`

	// Memory defines the memory in GB for the instance.
	//
	// +optional
	Memory string `json:"memory,omitempty"`

	// Processors defines the processing units for the instance.
	//
	// +optional
	Processors string `json:"processors,omitempty"`

	// ProcType defines the processor sharing model for the instance.
	// Must be one of {capped, dedicated, shared}.
	//
	// +optional
	ProcType ProcType `json:"procType,omitempty"`

	// SysType defines the system type for instance.
	//
	// +optional
	SysType string `json:"sysType,omitempty"`
}

// Set stores values from required into a
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}
	if len(required.VolumeIDs) != 0 {
		a.VolumeIDs = required.VolumeIDs
	}
	if required.Memory != "" {
		a.Memory = required.Memory
	}
	if required.Processors != "" {
		a.Processors = required.Processors
	}
	if required.ProcType != "" {
		a.ProcType = required.ProcType
	}
	if required.SysType != "" {
		a.SysType = required.SysType
	}
}
