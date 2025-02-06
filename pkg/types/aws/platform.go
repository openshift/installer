package aws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/endpoints"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/dns"
)

const (
	// VolumeTypeGp2 is the type of EBS volume for General Purpose SSD gp2.
	VolumeTypeGp2 = "gp2"
	// VolumeTypeGp3 is the type of EBS volume for General Purpose SSD gp3.
	VolumeTypeGp3 = "gp3"
)

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// The field is deprecated. AMIID is the AMI that should be used to boot
	// machines for the cluster. If set, the AMI should belong to the same
	// region as the cluster.
	//
	// +optional
	AMIID string `json:"amiID,omitempty"`

	// Region specifies the AWS region where the cluster will be created.
	Region string `json:"region"`

	// Subnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	//
	// +optional
	Subnets []string `json:"subnets,omitempty"`

	// HostedZone is the ID of an existing hosted zone into which to add DNS
	// records for the cluster's internal API. An existing hosted zone can
	// only be used when also using existing subnets. The hosted zone must be
	// associated with the VPC containing the subnets.
	// Leave the hosted zone unset to have the installer create the hosted zone
	// on your behalf.
	// +optional
	HostedZone string `json:"hostedZone,omitempty"`

	// HostedZoneRole is the ARN of an IAM role to be assumed when performing
	// operations on the provided HostedZone. HostedZoneRole can be used
	// in a shared VPC scenario when the private hosted zone belongs to a
	// different account than the rest of the cluster resources.
	// If HostedZoneRole is set, HostedZone must also be set.
	//
	// +optional
	HostedZoneRole string `json:"hostedZoneRole,omitempty"`

	// UserTags additional keys and values that the installer will add
	// as tags to all resources that it creates. Resources created by the
	// cluster itself may not include these tags.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

	// ServiceEndpoints list contains custom endpoints which will override default
	// service endpoint of AWS Services.
	// There must be only one ServiceEndpoint for a service.
	// +optional
	ServiceEndpoints []ServiceEndpoint `json:"serviceEndpoints,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on AWS for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// The field is deprecated. ExperimentalPropagateUserTags is an experimental
	// flag that directs in-cluster operators to include the specified
	// user tags in the tags of the AWS resources that the operators create.
	// +optional
	ExperimentalPropagateUserTag *bool `json:"experimentalPropagateUserTags,omitempty"`

	// PropagateUserTags is a flag that directs in-cluster operators
	// to include the specified user tags in the tags of the
	// AWS resources that the operators create.
	// +optional
	PropagateUserTag bool `json:"propagateUserTags,omitempty"`

	// LBType is an optional field to specify a load balancer type.
	// When this field is specified, all ingresscontrollers (including the
	// default ingresscontroller) will be created using the specified load-balancer
	// type by default.
	//
	// Following are the accepted values:
	//
	// * "Classic": A Classic Load Balancer that makes routing decisions at
	// either the transport layer (TCP/SSL) or the application layer
	// (HTTP/HTTPS). See the following for additional details:
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/load-balancer-types.html#clb
	//
	// * "NLB": A Network Load Balancer that makes routing decisions at the
	// transport layer (TCP/SSL). See the following for additional details:
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/load-balancer-types.html#nlb
	//
	// If this field is not set explicitly, it defaults to "Classic".  This
	// default is subject to change over time.
	//
	// +optional
	LBType configv1.AWSLBType `json:"lbType,omitempty"`

	// eipAllocations contains Elastic IP (EIP) allocations for AWS resources
	// within the cluster.
	//
	// +optional
	EIPAllocations *EIPAllocations `json:"eipAllocations,omitempty"`

	// PreserveBootstrapIgnition is deprecated. Use bestEffortDeleteIgnition instead.
	// +optional
	PreserveBootstrapIgnition bool `json:"preserveBootstrapIgnition,omitempty"`

	// BestEffortDeleteIgnition is an optional field that can be used to ignore errors from S3 deletion of ignition
	// objects during cluster bootstrap. The default behavior is to fail the installation if ignition objects cannot be
	// deleted. Enable this functionality when there are known reasons disallowing their deletion.
	// +optional
	BestEffortDeleteIgnition bool `json:"bestEffortDeleteIgnition,omitempty"`

	// PublicIpv4Pool is an optional field that can be used to tell the installation process to use
	// Public IPv4 address that you bring to your AWS account with BYOIP.
	// +optional
	PublicIpv4Pool string `json:"publicIpv4Pool,omitempty"`

	// UserProvisionedDNS indicates if the customer is providing their own DNS solution in place of the default
	// provisioned by the Installer.
	// +kubebuilder:default:="Disabled"
	// +default="Disabled"
	// +kubebuilder:validation:Enum="Enabled";"Disabled"
	UserProvisionedDNS dns.UserProvisionedDNS `json:"userProvisionedDNS,omitempty"`
}

// EIPAllocations contains Elastic IP (EIP) allocations for AWS resources
// within the cluster.
type EIPAllocations struct {
	// IngressNetworkLoadBalancer is a list of IDs for Elastic IP (EIP) addresses that
	// are assigned to the default AWS NLB IngressController.
	// The following restrictions apply:
	//
	// eipAllocations can only be used with external scope, not internal.
	// An EIP can be allocated to only a single IngressController.
	// The number of EIP allocations must match the number of subnets that are used for the load balancer.
	// Each EIP allocation must be unique.
	// A maximum of 10 EIP allocations are permitted.
	//
	// See https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html for general
	// information about configuration, characteristics, and limitations of Elastic IP addresses.
	//
	// +openshift:enable:FeatureGate=InstallEIPForDefaultNLBIngressController
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:XValidation:rule=`self.all(x, self.exists_one(y, x == y))`,message="eipAllocations cannot contain duplicates"
	// +kubebuilder:validation:MaxItems=10
	IngressNetworkLoadBalancer []EIPAllocation `json:"ingressNetworkLoadBalancer"`
}

// EIPAllocation is an ID for an Elastic IP (EIP) address that can be allocated to an ELB in the AWS environment.
// Values must begin with `eipalloc-` followed by exactly 17 hexadecimal (`[0-9a-fA-F]`) characters.
// + Explanation of the regex `^eipalloc-[0-9a-fA-F]{17}$` for validating value of the EIPAllocation:
// + ^eipalloc- ensures the string starts with "eipalloc-".
// + [0-9a-fA-F]{17} matches exactly 17 hexadecimal characters (0-9, a-f, A-F).
// + $ ensures the string ends after the 17 hexadecimal characters.
// + Example of Valid and Invalid values:
// + eipalloc-1234567890abcdef1 is valid.
// + eipalloc-1234567890abcde is not valid (too short).
// + eipalloc-1234567890abcdefg is not valid (contains a non-hex character 'g').
// + Max length is calculated as follows:
// + eipalloc- = 9 chars and 17 hexadecimal chars after `-`
// + So, total is 17 + 9 = 26 chars required for value of an EIPAllocation.
//
// +kubebuilder:validation:MinLength=26
// +kubebuilder:validation:MaxLength=26
// +kubebuilder:validation:XValidation:rule=`self.startsWith('eipalloc-')`,message="eipAllocations should start with 'eipalloc-'"
// +kubebuilder:validation:XValidation:rule=`self.split("-", 2)[1].matches('[0-9a-fA-F]{17}$')`,message="eipAllocations must be 'eipalloc-' followed by exactly 17 hexadecimal characters (0-9, a-f, A-F)"
type EIPAllocation string

// ServiceEndpoint store the configuration for services to
// override existing defaults of AWS Services.
type ServiceEndpoint struct {
	// Name is the name of the AWS service.
	// This must be provided and cannot be empty.
	Name string `json:"name"`

	// URL is fully qualified URI with scheme https, that overrides the default generated
	// endpoint for a client.
	// This must be provided and cannot be empty.
	//
	// +kubebuilder:validation:Pattern=`^https://`
	URL string `json:"url"`
}

// IsSecretRegion returns true if the region is part of either the ISO or ISOB partitions.
func IsSecretRegion(region string) bool {
	partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region)
	if !ok {
		return false
	}
	switch partition.ID() {
	case endpoints.AwsIsoPartitionID, endpoints.AwsIsoBPartitionID:
		return true
	}
	return false
}

// IsPublicOnlySubnetsEnabled returns whether the public-only subnets feature has been enabled via env var.
func IsPublicOnlySubnetsEnabled() bool {
	// Even though this looks too simple for a function, it's better than having to update the logic everywhere it's
	// used in case we decide to check for specific values set in the env var.
	return os.Getenv("OPENSHIFT_INSTALL_AWS_PUBLIC_ONLY") != ""
}
