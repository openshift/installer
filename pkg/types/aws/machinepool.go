package aws

// MachinePool stores the configuration for a machine pool installed
// on AWS.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the ec2 instance type.
	// eg. m4-large
	InstanceType string `json:"type"`

	// EC2RootVolume defines the root volume for EC2 instances in the machine pool.
	EC2RootVolume `json:"rootVolume"`
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

	if required.EC2RootVolume.IOPS != 0 {
		a.EC2RootVolume.IOPS = required.EC2RootVolume.IOPS
	}
	if required.EC2RootVolume.Size != 0 {
		a.EC2RootVolume.Size = required.EC2RootVolume.Size
	}
	if required.EC2RootVolume.Type != "" {
		a.EC2RootVolume.Type = required.EC2RootVolume.Type
	}
	if required.EC2RootVolume.KMSKeyARN != "" {
		a.EC2RootVolume.KMSKeyARN = required.EC2RootVolume.KMSKeyARN
	}
}

// EC2RootVolume defines the storage for an ec2 instance.
type EC2RootVolume struct {
	// IOPS defines the amount of provisioned IOPS. This is only valid
	// for type io1.
	IOPS int `json:"iops"`
	// Size defines the size of the volume in gibibytes (GiB).
	Size int `json:"size"`
	// Type defines the type of the volume.
	Type string `json:"type"`

	// The KMS key that will be used to encrypt the EBS volume.
	// If no key is provided the default KMS key for the account will be used.
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_GetEbsDefaultKmsKeyId.html
	// +optional
	KMSKeyARN string `json:"kmsKeyARN,omitempty"`
}
