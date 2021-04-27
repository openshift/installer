package powervs

// MachinePool stores the configuration for a machine pool installed on Ibm Power VS.
type MachinePool struct {
	// ServiceInstance is Service Instance to install into.
	//
	ServiceInstance string `json:"serviceinstance"`

    // Name is the name of the instance
    //
    Name string `json:"name"`

    // VolumeIDs is the list of volumes attached to the instance.
    //
    VolumeIDs []string `json:"volumeIDs"`

	// Memory defines the memory in GB for the instance.
	//
	Memory int `json:"memory"`

    // CPU defines the processing units for the instance.
    // @TODO: 
	CPU float32 `json:"cpu"`

    // ProcShare defines the processor sharing model for the instance.
	//
    // +optional
    ProcShare string `json:"procShare"`

    // ImageID defines the ImageID for the instance.
	//
	// +optional (does this mean user-optional, or completely?)
    ImageID string `json:"imageID"`

    // Networks defines the network IDs of the instance.
	//
	// +optional
    Networks []string `json:"networks"`

    // KeyPair defines the keypair name for instance.
	//
	// +optional
	KeyPair string `json:"keypair"`

	// SystemType defines the system type for instance.
	//
	// +optional
	SystemType string `json:systemType"`
}

// Sets values taken from passed MachinePool.
/*
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

	if required.OSDisk.DiskSizeGB > 0 {
		a.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}

	if required.OSDisk.DiskType != "" {
		a.OSDisk.DiskType = required.OSDisk.DiskType
	}

	if required.EncryptionKey != nil {
		if a.EncryptionKey == nil {
			a.EncryptionKey = &EncryptionKeyReference{}
		}
		a.EncryptionKey.Set(required.EncryptionKey)
	}
}*/
