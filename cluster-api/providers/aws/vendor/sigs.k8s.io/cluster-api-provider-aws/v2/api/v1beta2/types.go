/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta2

import (
	"k8s.io/apimachinery/pkg/util/sets"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// AWSResourceReference is a reference to a specific AWS resource by ID or filters.
// Only one of ID or Filters may be specified. Specifying more than one will result in
// a validation error.
type AWSResourceReference struct {
	// ID of resource
	// +optional
	ID *string `json:"id,omitempty"`

	// Filters is a set of key/value pairs used to identify a resource
	// They are applied according to the rules defined by the AWS API:
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html
	// +optional
	Filters []Filter `json:"filters,omitempty"`
}

// AMIReference is a reference to a specific AWS resource by ID, ARN, or filters.
// Only one of ID, ARN or Filters may be specified. Specifying more than one will result in
// a validation error.
type AMIReference struct {
	// ID of resource
	// +optional
	ID *string `json:"id,omitempty"`

	// EKSOptimizedLookupType If specified, will look up an EKS Optimized image in SSM Parameter store
	// +kubebuilder:validation:Enum:=AmazonLinux;AmazonLinuxGPU
	// +optional
	EKSOptimizedLookupType *EKSAMILookupType `json:"eksLookupType,omitempty"`
}

// Filter is a filter used to identify an AWS resource.
type Filter struct {
	// Name of the filter. Filter names are case-sensitive.
	Name string `json:"name"`

	// Values includes one or more filter values. Filter values are case-sensitive.
	Values []string `json:"values"`
}

// AWSMachineProviderConditionType is a valid value for AWSMachineProviderCondition.Type.
type AWSMachineProviderConditionType string

// Valid conditions for an AWS machine instance.
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated AWSMachineProviderConditionType = "MachineCreated"
)

const (
	// ExternalResourceGCAnnotation is the name of an annotation that indicates if
	// external resources should be garbage collected for the cluster.
	ExternalResourceGCAnnotation = "aws.cluster.x-k8s.io/external-resource-gc"

	// ExternalResourceGCTasksAnnotation is the name of an annotation that indicates what
	// external resources tasks should be executed by garbage collector for the cluster.
	ExternalResourceGCTasksAnnotation = "aws.cluster.x-k8s.io/external-resource-tasks-gc"
)

// GCTask defines a task to be executed by the garbage collector.
type GCTask string

var (
	// GCTaskLoadBalancer defines a task to cleaning up resources for AWS load balancers.
	GCTaskLoadBalancer = GCTask("load-balancer")

	// GCTaskTargetGroup defines a task to cleaning up resources for AWS target groups.
	GCTaskTargetGroup = GCTask("target-group")

	// GCTaskSecurityGroup defines a task to cleaning up resources for AWS security groups.
	GCTaskSecurityGroup = GCTask("security-group")
)

// AZSelectionScheme defines the scheme of selecting AZs.
type AZSelectionScheme string

var (
	// AZSelectionSchemeOrdered will select AZs based on alphabetical order.
	AZSelectionSchemeOrdered = AZSelectionScheme("Ordered")

	// AZSelectionSchemeRandom will select AZs randomly.
	AZSelectionSchemeRandom = AZSelectionScheme("Random")
)

// InstanceState describes the state of an AWS instance.
type InstanceState string

var (
	// InstanceStatePending is the string representing an instance in a pending state.
	InstanceStatePending = InstanceState("pending")

	// InstanceStateRunning is the string representing an instance in a running state.
	InstanceStateRunning = InstanceState("running")

	// InstanceStateShuttingDown is the string representing an instance shutting down.
	InstanceStateShuttingDown = InstanceState("shutting-down")

	// InstanceStateTerminated is the string representing an instance that has been terminated.
	InstanceStateTerminated = InstanceState("terminated")

	// InstanceStateStopping is the string representing an instance
	// that is in the process of being stopped and can be restarted.
	InstanceStateStopping = InstanceState("stopping")

	// InstanceStateStopped is the string representing an instance
	// that has been stopped and can be restarted.
	InstanceStateStopped = InstanceState("stopped")

	// InstanceRunningStates defines the set of states in which an EC2 instance is
	// running or going to be running soon.
	InstanceRunningStates = sets.NewString(
		string(InstanceStatePending),
		string(InstanceStateRunning),
	)

	// InstanceOperationalStates defines the set of states in which an EC2 instance is
	// or can return to running, and supports all EC2 operations.
	InstanceOperationalStates = InstanceRunningStates.Union(
		sets.NewString(
			string(InstanceStateStopping),
			string(InstanceStateStopped),
		),
	)

	// InstanceKnownStates represents all known EC2 instance states.
	InstanceKnownStates = InstanceOperationalStates.Union(
		sets.NewString(
			string(InstanceStateShuttingDown),
			string(InstanceStateTerminated),
		),
	)
)

// Instance describes an AWS instance.
type Instance struct {
	ID string `json:"id"`

	// The current state of the instance.
	State InstanceState `json:"instanceState,omitempty"`

	// The instance type.
	Type string `json:"type,omitempty"`

	// The ID of the subnet of the instance.
	SubnetID string `json:"subnetId,omitempty"`

	// The ID of the AMI used to launch the instance.
	ImageID string `json:"imageId,omitempty"`

	// The name of the SSH key pair.
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// SecurityGroupIDs are one or more security group IDs this instance belongs to.
	SecurityGroupIDs []string `json:"securityGroupIds,omitempty"`

	// UserData is the raw data script passed to the instance which is run upon bootstrap.
	// This field must not be base64 encoded and should only be used when running a new instance.
	UserData *string `json:"userData,omitempty"`

	// The name of the IAM instance profile associated with the instance, if applicable.
	IAMProfile string `json:"iamProfile,omitempty"`

	// Addresses contains the AWS instance associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// The private IPv4 address assigned to the instance.
	PrivateIP *string `json:"privateIp,omitempty"`

	// The public IPv4 address assigned to the instance, if applicable.
	PublicIP *string `json:"publicIp,omitempty"`

	// Specifies whether enhanced networking with ENA is enabled.
	ENASupport *bool `json:"enaSupport,omitempty"`

	// Indicates whether the instance is optimized for Amazon EBS I/O.
	EBSOptimized *bool `json:"ebsOptimized,omitempty"`

	// Configuration options for the root storage volume.
	// +optional
	RootVolume *Volume `json:"rootVolume,omitempty"`

	// Configuration options for the non root storage volumes.
	// +optional
	NonRootVolumes []Volume `json:"nonRootVolumes,omitempty"`

	// Specifies ENIs attached to instance
	NetworkInterfaces []string `json:"networkInterfaces,omitempty"`

	// The tags associated with the instance.
	Tags map[string]string `json:"tags,omitempty"`

	// Availability zone of instance
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// SpotMarketOptions option for configuring instances to be run using AWS Spot instances.
	SpotMarketOptions *SpotMarketOptions `json:"spotMarketOptions,omitempty"`

	// PlacementGroupName specifies the name of the placement group in which to launch the instance.
	// +optional
	PlacementGroupName string `json:"placementGroupName,omitempty"`

	// PlacementGroupPartition is the partition number within the placement group in which to launch the instance.
	// This value is only valid if the placement group, referred in `PlacementGroupName`, was created with
	// strategy set to partition.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=7
	// +optional
	PlacementGroupPartition int64 `json:"placementGroupPartition,omitempty"`

	// Tenancy indicates if instance should run on shared or single-tenant hardware.
	// +optional
	Tenancy string `json:"tenancy,omitempty"`

	// IDs of the instance's volumes
	// +optional
	VolumeIDs []string `json:"volumeIDs,omitempty"`

	// InstanceMetadataOptions is the metadata options for the EC2 instance.
	// +optional
	InstanceMetadataOptions *InstanceMetadataOptions `json:"instanceMetadataOptions,omitempty"`

	// PrivateDNSName is the options for the instance hostname.
	// +optional
	PrivateDNSName *PrivateDNSName `json:"privateDnsName,omitempty"`

	// PublicIPOnLaunch is the option to associate a public IP on instance launch
	// +optional
	PublicIPOnLaunch *bool `json:"publicIPOnLaunch,omitempty"`
}

// InstanceMetadataState describes the state of InstanceMetadataOptions.HttpEndpoint and InstanceMetadataOptions.InstanceMetadataTags
type InstanceMetadataState string

const (
	// InstanceMetadataEndpointStateDisabled represents the disabled state
	InstanceMetadataEndpointStateDisabled = InstanceMetadataState("disabled")

	// InstanceMetadataEndpointStateEnabled represents the enabled state
	InstanceMetadataEndpointStateEnabled = InstanceMetadataState("enabled")
)

// HTTPTokensState describes the state of InstanceMetadataOptions.HTTPTokensState
type HTTPTokensState string

const (
	// HTTPTokensStateOptional represents the optional state
	HTTPTokensStateOptional = HTTPTokensState("optional")

	// HTTPTokensStateRequired represents the required state (IMDSv2)
	HTTPTokensStateRequired = HTTPTokensState("required")
)

// InstanceMetadataOptions describes metadata options for the EC2 instance.
type InstanceMetadataOptions struct {
	// Enables or disables the HTTP metadata endpoint on your instances.
	//
	// If you specify a value of disabled, you cannot access your instance metadata.
	//
	// Default: enabled
	//
	// +kubebuilder:validation:Enum:=enabled;disabled
	// +kubebuilder:default=enabled
	HTTPEndpoint InstanceMetadataState `json:"httpEndpoint,omitempty"`

	// The desired HTTP PUT response hop limit for instance metadata requests. The
	// larger the number, the further instance metadata requests can travel.
	//
	// Default: 1
	//
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=64
	// +kubebuilder:default=1
	HTTPPutResponseHopLimit int64 `json:"httpPutResponseHopLimit,omitempty"`

	// The state of token usage for your instance metadata requests.
	//
	// If the state is optional, you can choose to retrieve instance metadata with
	// or without a session token on your request. If you retrieve the IAM role
	// credentials without a token, the version 1.0 role credentials are returned.
	// If you retrieve the IAM role credentials using a valid session token, the
	// version 2.0 role credentials are returned.
	//
	// If the state is required, you must send a session token with any instance
	// metadata retrieval requests. In this state, retrieving the IAM role credentials
	// always returns the version 2.0 credentials; the version 1.0 credentials are
	// not available.
	//
	// Default: optional
	//
	// +kubebuilder:validation:Enum:=optional;required
	// +kubebuilder:default=optional
	HTTPTokens HTTPTokensState `json:"httpTokens,omitempty"`

	// Set to enabled to allow access to instance tags from the instance metadata.
	// Set to disabled to turn off access to instance tags from the instance metadata.
	// For more information, see Work with instance tags using the instance metadata
	// (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#work-with-tags-in-IMDS).
	//
	// Default: disabled
	//
	// +kubebuilder:validation:Enum:=enabled;disabled
	// +kubebuilder:default=disabled
	InstanceMetadataTags InstanceMetadataState `json:"instanceMetadataTags,omitempty"`
}

// SetDefaults sets the default values for the InstanceMetadataOptions.
func (obj *InstanceMetadataOptions) SetDefaults() {
	if obj.HTTPEndpoint == "" {
		obj.HTTPEndpoint = InstanceMetadataEndpointStateEnabled
	}
	if obj.HTTPPutResponseHopLimit == 0 {
		obj.HTTPPutResponseHopLimit = 1
	}
	if obj.HTTPTokens == "" {
		obj.HTTPTokens = HTTPTokensStateOptional // Defaults to IMDSv1
	}
	if obj.InstanceMetadataTags == "" {
		obj.InstanceMetadataTags = InstanceMetadataEndpointStateDisabled
	}
}

// Volume encapsulates the configuration options for the storage device.
type Volume struct {
	// Device name
	// +optional
	DeviceName string `json:"deviceName,omitempty"`

	// Size specifies size (in Gi) of the storage device.
	// Must be greater than the image snapshot size or 8 (whichever is greater).
	// +kubebuilder:validation:Minimum=8
	Size int64 `json:"size"`

	// Type is the type of the volume (e.g. gp2, io1, etc...).
	// +optional
	Type VolumeType `json:"type,omitempty"`

	// IOPS is the number of IOPS requested for the disk. Not applicable to all types.
	// +optional
	IOPS int64 `json:"iops,omitempty"`

	// Throughput to provision in MiB/s supported for the volume type. Not applicable to all types.
	// +optional
	Throughput *int64 `json:"throughput,omitempty"`

	// Encrypted is whether the volume should be encrypted or not.
	// +optional
	Encrypted *bool `json:"encrypted,omitempty"`

	// EncryptionKey is the KMS key to use to encrypt the volume. Can be either a KMS key ID or ARN.
	// If Encrypted is set and this is omitted, the default AWS key will be used.
	// The key must already exist and be accessible by the controller.
	// +optional
	EncryptionKey string `json:"encryptionKey,omitempty"`
}

// VolumeType describes the EBS volume type.
// See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html
type VolumeType string

var (
	// VolumeTypeIO1 is the string representing a provisioned iops ssd io1 volume.
	VolumeTypeIO1 = VolumeType("io1")

	// VolumeTypeIO2 is the string representing a provisioned iops ssd io2 volume.
	VolumeTypeIO2 = VolumeType("io2")

	// VolumeTypeGP2 is the string representing a general purpose ssd gp2 volume.
	VolumeTypeGP2 = VolumeType("gp2")

	// VolumeTypeGP3 is the string representing a general purpose ssd gp3 volume.
	VolumeTypeGP3 = VolumeType("gp3")

	// VolumeTypesGP are volume types provisioned for general purpose io.
	VolumeTypesGP = sets.NewString(
		string(VolumeTypeIO1),
		string(VolumeTypeIO2),
	)

	// VolumeTypesProvisioned are volume types provisioned for high performance io.
	VolumeTypesProvisioned = sets.NewString(
		string(VolumeTypeIO1),
		string(VolumeTypeIO2),
	)
)

// SpotMarketOptions defines the options available to a user when configuring
// Machines to run on Spot instances.
// Most users should provide an empty struct.
type SpotMarketOptions struct {
	// MaxPrice defines the maximum price the user is willing to pay for Spot VM instances
	// +optional
	// +kubebuilder:validation:pattern="^[0-9]+(\.[0-9]+)?$"
	MaxPrice *string `json:"maxPrice,omitempty"`
}

// EKSAMILookupType specifies which AWS AMI to use for a AWSMachine and AWSMachinePool.
type EKSAMILookupType string

const (
	// AmazonLinux is the default AMI type.
	AmazonLinux EKSAMILookupType = "AmazonLinux"
	// AmazonLinuxGPU is the AmazonLinux GPU AMI type.
	AmazonLinuxGPU EKSAMILookupType = "AmazonLinuxGPU"
)

// PrivateDNSName is the options for the instance hostname.
type PrivateDNSName struct {
	// EnableResourceNameDNSAAAARecord indicates whether to respond to DNS queries for instance hostnames with DNS AAAA records.
	// +optional
	EnableResourceNameDNSAAAARecord *bool `json:"enableResourceNameDnsAAAARecord,omitempty"`
	// EnableResourceNameDNSARecord indicates whether to respond to DNS queries for instance hostnames with DNS A records.
	// +optional
	EnableResourceNameDNSARecord *bool `json:"enableResourceNameDnsARecord,omitempty"`
	// The type of hostname to assign to an instance.
	// +optional
	// +kubebuilder:validation:Enum:=ip-name;resource-name
	HostnameType *string `json:"hostnameType,omitempty"`
}
