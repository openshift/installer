package aws

// MachinePool stores the configuration for a machine pool installed
// on AWS.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	//
	// +optional
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the ec2 instance type.
	// eg. m4-large
	//
	// +optional
	InstanceType string `json:"type"`

	// AMIID is the AMI that should be used to boot the ec2 instance.
	// If set, the AMI should belong to the same region as the cluster.
	//
	// +optional
	AMIID string `json:"amiID,omitempty"`

	// EC2RootVolume defines the root volume for EC2 instances in the machine pool.
	//
	// +optional
	EC2RootVolume `json:"rootVolume"`

	// EC2MetadataOptions defines metadata service interaction options for EC2 instances in the machine pool.
	//
	// +optional
	EC2Metadata EC2Metadata `json:"metadataService"`

	// IAMRole is the name of the IAM Role to use for the instance profile of the machine.
	// Leave unset to have the installer create the IAM Role on your behalf.
	// +optional
	IAMRole string `json:"iamRole,omitempty"`

	// AdditionalSecurityGroupIDs contains IDs of additional security groups for machines, where each ID
	// is presented in the format sg-xxxx.
	//
	// +optional
	AdditionalSecurityGroupIDs []string `json:"additionalSecurityGroupIDs,omitempty"`
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

	if required.AMIID != "" {
		a.AMIID = required.AMIID
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

	if required.EC2Metadata.Authentication != "" {
		a.EC2Metadata.Authentication = required.EC2Metadata.Authentication
	}

	if required.IAMRole != "" {
		a.IAMRole = required.IAMRole
	}

	if len(required.AdditionalSecurityGroupIDs) > 0 {
		a.AdditionalSecurityGroupIDs = required.AdditionalSecurityGroupIDs
	}
}

// EC2RootVolume defines the storage for an ec2 instance.
type EC2RootVolume struct {
	// IOPS defines the amount of provisioned IOPS. (KiB/s). IOPS may only be set for
	// io1, io2, & gp3 volume types.
	//
	// +kubebuilder:validation:Minimum=0
	// +optional
	IOPS int `json:"iops"`

	// Size defines the size of the volume in gibibytes (GiB).
	//
	// +kubebuilder:validation:Minimum=0
	Size int `json:"size"`

	// Type defines the type of the volume.
	Type string `json:"type"`

	// The KMS key that will be used to encrypt the EBS volume.
	// If no key is provided the default KMS key for the account will be used.
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_GetEbsDefaultKmsKeyId.html
	// +optional
	KMSKeyARN string `json:"kmsKeyARN,omitempty"`
}

// EC2Metadata defines the metadata service interaction options for an ec2 instance.
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-instance-metadata-service.html
type EC2Metadata struct {
	// Authentication determines whether or not the host requires the use of authentication when interacting with the metadata service.
	// When using authentication, this enforces v2 interaction method (IMDSv2) with the metadata service.
	// When omitted, this means the user has no opinion and the value is left to the platform to choose a good
	// default, which is subject to change over time. The current default is optional.
	// At this point this field represents `HttpTokens` parameter from `InstanceMetadataOptionsRequest` structure in AWS EC2 API
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceMetadataOptionsRequest.html
	// +kubebuilder:validation:Enum=Required;Optional
	// +optional
	Authentication string `json:"authentication,omitempty"`
}
