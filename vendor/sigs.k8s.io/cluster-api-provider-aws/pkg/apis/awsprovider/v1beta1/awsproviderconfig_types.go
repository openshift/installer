/*
Copyright 2018 The Kubernetes Authors.

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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSMachineProviderConfig is the Schema for the awsmachineproviderconfigs API
// +k8s:openapi-gen=true
type AWSMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// AMI is the reference to the AMI from which to create the machine instance.
	AMI AWSResourceReference `json:"ami"`

	// InstanceType is the type of instance to create. Example: m4.xlarge
	InstanceType string `json:"instanceType"`

	// Tags is the set of tags to add to apply to an instance, in addition to the ones
	// added by default by the actuator. These tags are additive. The actuator will ensure
	// these tags are present, but will not remove any other tags that may exist on the
	// instance.
	Tags []TagSpecification `json:"tags,omitempty"`

	// IAMInstanceProfile is a reference to an IAM role to assign to the instance
	IAMInstanceProfile *AWSResourceReference `json:"iamInstanceProfile,omitempty"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is a reference to the secret with AWS credentials. Otherwise, defaults to permissions
	// provided by attached IAM role where the actuator is running.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`

	// KeyName is the name of the KeyPair to use for SSH
	KeyName *string `json:"keyName,omitempty"`

	// DeviceIndex is the index of the device on the instance for the network interface attachment.
	// Defaults to 0.
	DeviceIndex int64 `json:"deviceIndex"`

	// PublicIP specifies whether the instance should get a public IP. If not present,
	// it should use the default of its subnet.
	PublicIP *bool `json:"publicIp,omitempty"`

	// SecurityGroups is an array of references to security groups that should be applied to the
	// instance.
	SecurityGroups []AWSResourceReference `json:"securityGroups,omitempty"`

	// Subnet is a reference to the subnet to use for this instance
	Subnet AWSResourceReference `json:"subnet"`

	// Placement specifies where to create the instance in AWS
	Placement Placement `json:"placement"`

	// LoadBalancers is the set of load balancers to which the new instance
	// should be added once it is created.
	LoadBalancers []LoadBalancerReference `json:"loadBalancers,omitempty"`

	// BlockDevices is the set of block device mapping associated to this instance,
	// block device without a name will be used as a root device and only one device without a name is allowed
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/block-device-mapping-concepts.html
	BlockDevices []BlockDeviceMappingSpec `json:"blockDevices,omitempty"`

	// SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.
	SpotMarketOptions *SpotMarketOptions `json:"spotMarketOptions,omitempty"`

	// Tenancy indicates if instance should run on shared or single-tenant hardware. There are
	// supported 3 options: default, dedicated and host.
	Tenancy InstanceTenancy `json:"tenancy,omitempty"`
}

// BlockDeviceMappingSpec describes a block device mapping
type BlockDeviceMappingSpec struct {

	// The device name exposed to the machine (for example, /dev/sdh or xvdh).
	DeviceName *string `json:"deviceName,omitempty"`

	// Parameters used to automatically set up EBS volumes when the machine is
	// launched.
	EBS *EBSBlockDeviceSpec `json:"ebs,omitempty"`

	// Suppresses the specified device included in the block device mapping of the
	// AMI.
	NoDevice *string `json:"noDevice,omitempty"`

	// The virtual device name (ephemeralN). Machine store volumes are numbered
	// starting from 0. An machine type with 2 available machine store volumes
	// can specify mappings for ephemeral0 and ephemeral1.The number of available
	// machine store volumes depends on the machine type. After you connect to
	// the machine, you must mount the volume.
	//
	// Constraints: For M3 machines, you must specify machine store volumes in
	// the block device mapping for the machine. When you launch an M3 machine,
	// we ignore any machine store volumes specified in the block device mapping
	// for the AMI.
	VirtualName *string `json:"virtualName,omitempty"`
}

// EBSBlockDeviceSpec describes a block device for an EBS volume.
// https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/EbsBlockDevice
type EBSBlockDeviceSpec struct {

	// Indicates whether the EBS volume is deleted on machine termination.
	DeleteOnTermination *bool `json:"deleteOnTermination,omitempty"`

	// Indicates whether the EBS volume is encrypted. Encrypted Amazon EBS volumes
	// may only be attached to machines that support Amazon EBS encryption.
	Encrypted *bool `json:"encrypted,omitempty"`

	// Indicates the KMS key that should be used to encrypt the Amazon EBS volume.
	KMSKey AWSResourceReference `json:"kmsKey,omitempty"`

	// The number of I/O operations per second (IOPS) that the volume supports.
	// For io1, this represents the number of IOPS that are provisioned for the
	// volume. For gp2, this represents the baseline performance of the volume and
	// the rate at which the volume accumulates I/O credits for bursting. For more
	// information about General Purpose SSD baseline performance, I/O credits,
	// and bursting, see Amazon EBS Volume Types (http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html)
	// in the Amazon Elastic Compute Cloud User Guide.
	//
	// Minimal and maximal IOPS for io1 and gp2 are constrained. Please, check
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
	// for precise boundaries for individual volumes.
	//
	// Condition: This parameter is required for requests to create io1 volumes;
	// it is not used in requests to create gp2, st1, sc1, or standard volumes.
	Iops *int64 `json:"iops,omitempty"`

	// The size of the volume, in GiB.
	//
	// Constraints: 1-16384 for General Purpose SSD (gp2), 4-16384 for Provisioned
	// IOPS SSD (io1), 500-16384 for Throughput Optimized HDD (st1), 500-16384 for
	// Cold HDD (sc1), and 1-1024 for Magnetic (standard) volumes. If you specify
	// a snapshot, the volume size must be equal to or larger than the snapshot
	// size.
	//
	// Default: If you're creating the volume from a snapshot and don't specify
	// a volume size, the default is the snapshot size.
	VolumeSize *int64 `json:"volumeSize,omitempty"`

	// The volume type: gp2, io1, st1, sc1, or standard.
	// Default: standard
	VolumeType *string `json:"volumeType,omitempty"`
}

// SpotMarketOptions defines the options available to a user when configuring
// Machines to run on Spot instances.
// Most users should provide an empty struct.
type SpotMarketOptions struct {
	// The maximum price the user is willing to pay for their instances
	// Default: On-Demand price
	MaxPrice *string `json:"maxPrice,omitempty"`
}

// AWSResourceReference is a reference to a specific AWS resource by ID, ARN, or filters.
// Only one of ID, ARN or Filters may be specified. Specifying more than one will result in
// a validation error.
type AWSResourceReference struct {
	// ID of resource
	// +optional
	ID *string `json:"id,omitempty"`

	// ARN of resource
	// +optional
	ARN *string `json:"arn,omitempty"`

	// Filters is a set of filters used to identify a resource
	Filters []Filter `json:"filters,omitempty"`
}

// Placement indicates where to create the instance in AWS
type Placement struct {
	// Region is the region to use to create the instance
	Region string `json:"region,omitempty"`

	// AvailabilityZone is the availability zone of the instance
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

// Filter is a filter used to identify an AWS resource
type Filter struct {
	// Name of the filter. Filter names are case-sensitive.
	Name string `json:"name"`

	// Values includes one or more filter values. Filter values are case-sensitive.
	Values []string `json:"values,omitempty"`
}

// TagSpecification is the name/value pair for a tag
type TagSpecification struct {
	// Name of the tag
	Name string `json:"name"`

	// Value of the tag
	Value string `json:"value"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSMachineProviderConfigList contains a list of AWSMachineProviderConfig
type AWSMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSMachineProviderConfig `json:"items"`
}

// LoadBalancerReference is a reference to a load balancer on AWS.
type LoadBalancerReference struct {
	Name string              `json:"name"`
	Type AWSLoadBalancerType `json:"type"`
}

// AWSLoadBalancerType is the type of LoadBalancer to use when registering
// an instance with load balancers specified in LoadBalancerNames
type AWSLoadBalancerType string

// InstanceTenancy indicates if instance should run on shared or single-tenant hardware.
type InstanceTenancy string

const (
	// DefaultTenancy instance runs on shared hardware
	DefaultTenancy InstanceTenancy = "default"
	// DedicatedTenancy instance runs on single-tenant hardware
	DedicatedTenancy InstanceTenancy = "dedicated"
	// HostTenancy instance runs on a Dedicated Host, which is an isolated server with configurations that you can control.
	HostTenancy InstanceTenancy = "host"
)

// Possible values for AWSLoadBalancerType. Add to this list as other types
// of load balancer are supported by the actuator.
const (
	ClassicLoadBalancerType AWSLoadBalancerType = "classic" // AWS classic ELB
	NetworkLoadBalancerType AWSLoadBalancerType = "network" // AWS Network Load Balancer (NLB)
)

func init() {
	SchemeBuilder.Register(&AWSMachineProviderConfig{}, &AWSMachineProviderConfigList{}, &AWSMachineProviderStatus{})
}
