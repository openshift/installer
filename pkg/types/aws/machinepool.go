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
	// Cannot be specified together with iamProfile.
	// +optional
	IAMRole string `json:"iamRole,omitempty"`

	// IAMProfile is the name of the IAM instance profile to use for the machine.
	// Leave unset to have the installer create the IAM Profile on your behalf.
	// Cannot be specified together with iamRole.
	// +optional
	IAMProfile string `json:"iamProfile,omitempty"`

	// AdditionalSecurityGroupIDs contains IDs of additional security groups for machines, where each ID
	// is presented in the format sg-xxxx.
	//
	// +kubebuilder:validation:MaxItems=10
	// +optional
	AdditionalSecurityGroupIDs []string `json:"additionalSecurityGroupIDs,omitempty"`

	// CPUOptions defines CPU-related settings for the instance, including the confidential computing policy.
	// When omitted, this means no opinion and the AWS platform is left to choose a reasonable default.
	// More info:
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CpuOptionsRequest.html,
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/cpu-options-supported-instances-values.html
	// +optional
	CPUOptions *CPUOptions `json:"cpuOptions,omitempty,omitzero"`

	// hostPlacement configures placement on AWS Dedicated Hosts. This allows admins to assign instances to specific host
	// for a variety of needs including for regulatory compliance, to leverage existing per-socket or per-core software licenses (BYOL),
	// and to gain visibility and control over instance placement on a physical server.
	// When omitted, the instance is not constrained to a dedicated host.
	// +openshift:enable:FeatureGate=AWSDedicatedHosts
	// +optional
	HostPlacement *HostPlacement `json:"hostPlacement,omitempty"`
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

	if required.IAMProfile != "" {
		a.IAMProfile = required.IAMProfile
	}

	if len(required.AdditionalSecurityGroupIDs) > 0 {
		a.AdditionalSecurityGroupIDs = required.AdditionalSecurityGroupIDs
	}

	if required.CPUOptions != nil {
		a.CPUOptions = required.CPUOptions
	}

	if required.HostPlacement != nil {
		a.HostPlacement = required.HostPlacement
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

// ConfidentialComputePolicy represents the confidential compute configuration for the instance.
// +kubebuilder:validation:Enum=Disabled;AMDEncryptedVirtualizationNestedPaging
type ConfidentialComputePolicy string

const (
	// ConfidentialComputePolicyDisabled disables confidential computing for the instance.
	ConfidentialComputePolicyDisabled ConfidentialComputePolicy = "Disabled"
	// ConfidentialComputePolicySEVSNP enables AMD SEV-SNP as the confidential computing technology for the instance.
	ConfidentialComputePolicySEVSNP ConfidentialComputePolicy = "AMDEncryptedVirtualizationNestedPaging"
)

// CPUOptions defines CPU-related settings for the instance, including the confidential computing policy.
// If provided, it must not be empty — at least one field must be set.
// +kubebuilder:validation:MinProperties=1
type CPUOptions struct {
	// ConfidentialCompute specifies whether confidential computing should be enabled for the instance,
	// and, if so, which confidential computing technology to use.
	// Valid values are: Disabled, AMDEncryptedVirtualizationNestedPaging and omitted.
	// When set to Disabled, confidential computing will be disabled for the instance.
	// When set to AMDEncryptedVirtualizationNestedPaging, AMD SEV-SNP will be used as the confidential computing technology for the instance.
	// In this case, ensure the following conditions are met:
	// 1) The selected instance type supports AMD SEV-SNP.
	// 2) The selected AWS region supports AMD SEV-SNP.
	// 3) The selected AMI supports AMD SEV-SNP.
	// More details can be checked at https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/sev-snp.html
	// When omitted, this means no opinion and the AWS platform is left to choose a reasonable default,
	// which is subject to change without notice. The current default is Disabled.
	// +optional
	ConfidentialCompute *ConfidentialComputePolicy `json:"confidentialCompute,omitempty"`
}

// HostPlacement is the type that will be used to configure the placement of AWS instances.
// This can be configured for default placement (AnyAvailable) and dedicated hosts (DedicatedHost).
// +kubebuilder:validation:XValidation:rule="has(self.affinity) && self.affinity == 'DedicatedHost' ?  has(self.dedicatedHost) : !has(self.dedicatedHost)",message="dedicatedHost is required when affinity is DedicatedHost, and forbidden otherwise"
type HostPlacement struct {
	// affinity specifies the affinity setting for the instance.
	// Allowed values are AnyAvailable and DedicatedHost.
	// When Affinity is set to DedicatedHost, an instance started onto a specific host always restarts on the same host if stopped. In this scenario, the `dedicatedHost` field must be set.
	// When Affinity is set to AnyAvailable, and you stop and restart the instance, it can be restarted on any available host.
	// +required
	// +unionDiscriminator
	Affinity *HostAffinity `json:"affinity,omitempty"`

	// dedicatedHost specifies the exact host that an instance should be restarted on if stopped.
	// dedicatedHost is required when 'affinity' is set to DedicatedHost, and forbidden otherwise.
	// +optional
	// +unionMember
	DedicatedHost []DedicatedHost `json:"dedicatedHost,omitempty"`
}

// HostAffinity selects how an instance should be placed on AWS Dedicated Hosts.
// +kubebuilder:validation:Enum:=DedicatedHost;AnyAvailable
type HostAffinity string

const (
	// HostAffinityAnyAvailable lets the platform select any available dedicated host.
	HostAffinityAnyAvailable HostAffinity = "AnyAvailable"

	// HostAffinityDedicatedHost requires specifying a particular host via dedicatedHost.host.hostID.
	HostAffinityDedicatedHost HostAffinity = "DedicatedHost"
)

// DedicatedHost represents the configuration for the usage of dedicated host.
type DedicatedHost struct {
	// id identifies the AWS Dedicated Host on which the instance must run.
	// The value must start with "h-" followed by 17 lowercase hexadecimal characters (0-9 and a-f).
	// Must be exactly 19 characters in length.
	// +kubebuilder:validation:XValidation:rule="self.matches('^h-[0-9a-f]{17}$')",message="hostID must start with 'h-' followed by 17 lowercase hexadecimal characters (0-9 and a-f)"
	// +kubebuilder:validation:MinLength=19
	// +kubebuilder:validation:MaxLength=19
	// +required
	ID string `json:"id,omitempty"`

	// zone is the availability zone that the dedicated host belongs to
	// +optional
	Zone string `json:"zone,omitempty"`
}
