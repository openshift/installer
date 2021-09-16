package powervs

// MachinePool stores the configuration for a machine pool installed on IBM Power VS.
type MachinePool struct {
	// ServiceInstance is Service Instance to install into.
	//
	ServiceInstance string `json:"serviceinstance"`

	// Name is the name of the instance
	//
	Name string `json:"name"`

	// KeyPairName is the name of an SSH key pair stored in the Power VS
	// Service Instance
	KeyPairName string `json:"keypairname"`

	// VolumeIDs is the list of volumes attached to the instance.
	//
	// +optional
	VolumeIDs []string `json:"volumeIDs"`

	// Memory defines the memory in GB for the instance.
	//
	Memory int `json:"memory"`

	// Processors defines the processing units for the instance.
	// @TODO:
	Processors float64 `json:"processors"`

	// ProcType defines the processor sharing model for the instance.
	//
	// +optional
	ProcType string `json:"procType"`

	// ImageName defines the boot image name for the instance.
	//
	ImageName string `json:"imageName"`

	// NetworkIDs defines the network IDs of the instance.
	//
	NetworkIDs []string `json:"networkIDs"`

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
	if required.ImageName != "" {
		a.ImageName = required.ImageName
	}
	if required.ServiceInstance != "" {
		a.ServiceInstance = required.ServiceInstance
	}
	if len(required.NetworkIDs) > 0 {
		a.NetworkIDs = required.NetworkIDs
	}
	if required.ServiceInstance != "" {
		a.ServiceInstance = required.ServiceInstance
	}
	if required.KeyPairName != "" {
		a.KeyPairName = required.KeyPairName
	}
}
