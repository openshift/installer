package powervs

// MachinePool stores the configuration for a machine pool installed on IBM Power VS.
type MachinePool struct {
	// VolumeIDs is the list of volumes attached to the instance.
	//
	// +optional
	VolumeIDs []string `json:"volumeIDs"`

	// Memory defines the memory in GB for the instance.
	//
	// +optional
	Memory int `json:"memory"`

	// Processors defines the processing units for the instance.
	//
	// +optional
	Processors float64 `json:"processors"`

	// ProcType defines the processor sharing model for the instance.
	//
	// +optional
	ProcType string `json:"procType"`

	// SysType defines the system type for instance.
	//
	// +optional
	SysType string `json:"sysType"`
}

// Set stores values from required into a
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}
	if len(required.VolumeIDs) != 0 {
		a.VolumeIDs = required.VolumeIDs
	}
	if required.Memory != 0 {
		a.Memory = required.Memory
	}
	if required.Processors != 0 {
		a.Processors = required.Processors
	}
	if required.ProcType != "" {
		a.ProcType = required.ProcType
	}
	if required.SysType != "" {
		a.SysType = required.SysType
	}
}
