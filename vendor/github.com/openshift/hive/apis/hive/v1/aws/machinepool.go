package aws

// MachinePoolPlatform stores the configuration for a machine pool
// installed on AWS.
type MachinePoolPlatform struct {
	// Zones is list of availability zones that can be used.
	Zones []string `json:"zones,omitempty"`

	// Subnets is the list of IDs of subnets to which to attach the machines.
	// There must be exactly one subnet for each availability zone used.
	// These subnets may be public or private.
	// As a special case, for consistency with install-config, you may specify exactly one
	// private and one public subnet for each availability zone. In this case, the public
	// subnets will be filtered out and only the private subnets will be used.
	// If empty/omitted, we will look for subnets in each availability zone tagged with
	// Name=<clusterID>-private-<az> (legacy terraform) or <clusterID>-subnet-private-<az>
	// (CAPA).
	Subnets []string `json:"subnets,omitempty"`

	// InstanceType defines the ec2 instance type.
	// eg. m4-large
	InstanceType string `json:"type"`

	// EC2RootVolume defines the storage for ec2 instance.
	EC2RootVolume `json:"rootVolume"`

	// SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.
	// +optional
	SpotMarketOptions *SpotMarketOptions `json:"spotMarketOptions,omitempty"`

	// EC2MetadataOptions defines metadata service interaction options for EC2 instances in the machine pool.
	// +optional
	EC2Metadata *EC2Metadata `json:"metadataService,omitempty"`

	// AdditionalSecurityGroupIDs contains IDs of additional security groups for machines, where each ID
	// is presented in the format sg-xxxx.
	//
	// +optional
	AdditionalSecurityGroupIDs []string `json:"additionalSecurityGroupIDs,omitempty"`

	// UserTags contains the user defined tags to be supplied for the ec2 instance.
	// Note that these will be merged with ClusterDeployment.Spec.Platform.AWS.UserTags, with
	// this field taking precedence when keys collide.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`
}

// SpotMarketOptions defines the options available to a user when configuring
// Machines to run on Spot instances.
// Most users should provide an empty struct.
type SpotMarketOptions struct {
	// The maximum price the user is willing to pay for their instances
	// Default: On-Demand price
	// +optional
	MaxPrice *string `json:"maxPrice,omitempty"`
}

// EC2RootVolume defines the storage for an ec2 instance.
type EC2RootVolume struct {
	// IOPS defines the iops for the storage.
	// +optional
	IOPS int `json:"iops,omitempty"`
	// Size defines the size of the storage.
	Size int `json:"size"`
	// Type defines the type of the storage.
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
	// +optional
	Authentication string `json:"authentication,omitempty"`
}
