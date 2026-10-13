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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

const (
	// KindMachinePool is a MachinePool resource Kind
	KindMachinePool string = "MachinePool"
)

// EBS can be used to automatically set up EBS volumes when an instance is launched.
type EBS struct {
	// Encrypted is whether the volume should be encrypted or not.
	// +optional
	Encrypted bool `json:"encrypted,omitempty"`

	// The size of the volume, in GiB.
	// This can be a number from 1-1,024 for standard, 4-16,384 for io1, 1-16,384
	// for gp2, and 500-16,384 for st1 and sc1. If you specify a snapshot, the volume
	// size must be equal to or larger than the snapshot size.
	// +optional
	VolumeSize int64 `json:"volumeSize,omitempty"`

	// The volume type
	// For more information, see Amazon EBS Volume Types (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html)
	// +kubebuilder:validation:Enum=standard;io1;gp2;st1;sc1;io2
	// +optional
	VolumeType string `json:"volumeType,omitempty"`
}

// BlockDeviceMapping specifies the block devices for the instance.
// You can specify virtual devices and EBS volumes.
type BlockDeviceMapping struct {
	// The device name exposed to the EC2 instance (for example, /dev/sdh or xvdh).
	// +kubebuilder:validation:Required
	DeviceName string `json:"deviceName,omitempty"`

	// You can specify either VirtualName or Ebs, but not both.
	// +optional
	Ebs EBS `json:"ebs,omitempty"`
}

// AWSLaunchTemplate defines the desired state of AWSLaunchTemplate.
type AWSLaunchTemplate struct {
	// The name of the launch template.
	Name string `json:"name,omitempty"`

	// The name or the Amazon Resource Name (ARN) of the instance profile associated
	// with the IAM role for the instance. The instance profile contains the IAM
	// role.
	IamInstanceProfile string `json:"iamInstanceProfile,omitempty"`

	// AMI is the reference to the AMI from which to create the machine instance.
	// +optional
	AMI infrav1.AMIReference `json:"ami,omitempty"`

	// ImageLookupFormat is the AMI naming format to look up the image for this
	// machine It will be ignored if an explicit AMI is set. Supports
	// substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base OS and
	// kubernetes version, respectively. The BaseOS will be the value in
	// ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
	// defined by the packages produced by kubernetes/release without v as a
	// prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
	// image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
	// searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
	// Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
	// also: https://golang.org/pkg/text/template/
	// +optional
	ImageLookupFormat string `json:"imageLookupFormat,omitempty"`

	// ImageLookupOrg is the AWS Organization ID to use for image lookup if AMI is not set.
	ImageLookupOrg string `json:"imageLookupOrg,omitempty"`

	// ImageLookupBaseOS is the name of the base operating system to use for
	// image lookup the AMI is not set.
	ImageLookupBaseOS string `json:"imageLookupBaseOS,omitempty"`

	// InstanceType is the type of instance to create. Example: m4.xlarge
	InstanceType string `json:"instanceType,omitempty"`

	// RootVolume encapsulates the configuration options for the root volume
	// +optional
	RootVolume *infrav1.Volume `json:"rootVolume,omitempty"`

	// Configuration options for the non root storage volumes.
	// +optional
	NonRootVolumes []infrav1.Volume `json:"nonRootVolumes,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string
	// (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// VersionNumber is the version of the launch template that is applied.
	// Typically a new version is created when at least one of the following happens:
	// 1) A new launch template spec is applied.
	// 2) One or more parameters in an existing template is changed.
	// 3) A new AMI is discovered.
	VersionNumber *int64 `json:"versionNumber,omitempty"`

	// AdditionalSecurityGroups is an array of references to security groups that should be applied to the
	// instances. These security groups would be set in addition to any security groups defined
	// at the cluster level or in the actuator.
	// +optional
	AdditionalSecurityGroups []infrav1.AWSResourceReference `json:"additionalSecurityGroups,omitempty"`

	// SpotMarketOptions are options for configuring AWSMachinePool instances to be run using AWS Spot instances.
	SpotMarketOptions *infrav1.SpotMarketOptions `json:"spotMarketOptions,omitempty"`

	// InstanceMetadataOptions defines the behavior for applying metadata to instances.
	// +optional
	InstanceMetadataOptions *infrav1.InstanceMetadataOptions `json:"instanceMetadataOptions,omitempty"`

	// PrivateDNSName is the options for the instance hostname.
	// +optional
	PrivateDNSName *infrav1.PrivateDNSName `json:"privateDnsName,omitempty"`

	// CapacityReservationID specifies the target Capacity Reservation into which the instance should be launched.
	// +optional
	CapacityReservationID *string `json:"capacityReservationId,omitempty"`

	// MarketType specifies the type of market for the EC2 instance. Valid values include:
	// "OnDemand" (default): The instance runs as a standard OnDemand instance.
	// "Spot": The instance runs as a Spot instance. When SpotMarketOptions is provided, the marketType defaults to "Spot".
	// "CapacityBlock": The instance utilizes pre-purchased compute capacity (capacity blocks) with AWS Capacity Reservations.
	//  If this value is selected, CapacityReservationID must be specified to identify the target reservation.
	// If marketType is not specified and spotMarketOptions is provided, the marketType defaults to "Spot".
	// +optional
	MarketType infrav1.MarketType `json:"marketType,omitempty"`

	// CapacityReservationPreference specifies the preference for use of Capacity Reservations by the instance. Valid values include:
	// "Open": The instance may make use of open Capacity Reservations that match its AZ and InstanceType
	// "None": The instance may not make use of any Capacity Reservations. This is to conserve open reservations for desired workloads
	// "CapacityReservationsOnly": The instance will only run if matched or targeted to a Capacity Reservation
	// +kubebuilder:validation:Enum="";None;CapacityReservationsOnly;Open
	// +optional
	CapacityReservationPreference infrav1.CapacityReservationPreference `json:"capacityReservationPreference,omitempty"`
}

// Overrides are used to override the instance type specified by the launch template with multiple
// instance types that can be used to launch On-Demand Instances and Spot Instances.
type Overrides struct {
	InstanceType string `json:"instanceType"`
}

// OnDemandAllocationStrategy indicates how to allocate instance types to fulfill On-Demand capacity.
type OnDemandAllocationStrategy string

var (
	// OnDemandAllocationStrategyPrioritized uses the order of instance type overrides
	// for the LaunchTemplate to define the launch priority of each instance type.
	OnDemandAllocationStrategyPrioritized = OnDemandAllocationStrategy("prioritized")

	// OnDemandAllocationStrategyLowestPrice will make the Auto Scaling group launch
	// instances using the On-Demand pools with the lowest price, and evenly allocates
	// your instances across the On-Demand pools that you specify.
	OnDemandAllocationStrategyLowestPrice = OnDemandAllocationStrategy("lowest-price")
)

// SpotAllocationStrategy indicates how to allocate instances across Spot Instance pools.
type SpotAllocationStrategy string

var (
	// SpotAllocationStrategyLowestPrice will make the Auto Scaling group launch
	// instances using the Spot pools with the lowest price, and evenly allocates
	// your instances across the number of Spot pools that you specify.
	SpotAllocationStrategyLowestPrice = SpotAllocationStrategy("lowest-price")

	// SpotAllocationStrategyCapacityOptimized will make the Auto Scaling group launch
	// instances using Spot pools that are optimally chosen based on the available Spot capacity.
	SpotAllocationStrategyCapacityOptimized = SpotAllocationStrategy("capacity-optimized")

	// SpotAllocationStrategyCapacityOptimizedPrioritized will make the Auto Scaling group launch
	// instances using Spot pools that are optimally chosen based on the available Spot capacity
	// while also taking into account the priority order specified by the user for Instance Types.
	SpotAllocationStrategyCapacityOptimizedPrioritized = SpotAllocationStrategy("capacity-optimized-prioritized")

	// SpotAllocationStrategyPriceCapacityOptimized will make the Auto Scaling group launch
	// instances using Spot pools that consider both price and available Spot capacity to
	// provide a balance between cost savings and allocation reliability.
	SpotAllocationStrategyPriceCapacityOptimized = SpotAllocationStrategy("price-capacity-optimized")
)

// InstancesDistribution to configure distribution of On-Demand Instances and Spot Instances.
type InstancesDistribution struct {
	// +kubebuilder:validation:Enum=prioritized;lowest-price
	// +kubebuilder:default=prioritized
	OnDemandAllocationStrategy OnDemandAllocationStrategy `json:"onDemandAllocationStrategy,omitempty"`

	// +kubebuilder:validation:Enum=lowest-price;capacity-optimized;capacity-optimized-prioritized;price-capacity-optimized
	// +kubebuilder:default=lowest-price
	SpotAllocationStrategy SpotAllocationStrategy `json:"spotAllocationStrategy,omitempty"`

	// +kubebuilder:default=0
	OnDemandBaseCapacity *int64 `json:"onDemandBaseCapacity,omitempty"`

	// +kubebuilder:default=100
	OnDemandPercentageAboveBaseCapacity *int64 `json:"onDemandPercentageAboveBaseCapacity,omitempty"`
}

// MixedInstancesPolicy for an Auto Scaling group.
type MixedInstancesPolicy struct {
	InstancesDistribution *InstancesDistribution `json:"instancesDistribution,omitempty"`
	Overrides             []Overrides            `json:"overrides,omitempty"`
}

// Tags is a mapping for tags.
type Tags map[string]string

// AutoScalingGroup describes an AWS autoscaling group.
type AutoScalingGroup struct {
	// The tags associated with the instance.
	ID                    string          `json:"id,omitempty"`
	Tags                  infrav1.Tags    `json:"tags,omitempty"`
	Name                  string          `json:"name,omitempty"`
	DesiredCapacity       *int32          `json:"desiredCapacity,omitempty"`
	MaxSize               int32           `json:"maxSize,omitempty"`
	MinSize               int32           `json:"minSize,omitempty"`
	PlacementGroup        string          `json:"placementGroup,omitempty"`
	Subnets               []string        `json:"subnets,omitempty"`
	DefaultCoolDown       metav1.Duration `json:"defaultCoolDown,omitempty"`
	DefaultInstanceWarmup metav1.Duration `json:"defaultInstanceWarmup,omitempty"`
	CapacityRebalance     bool            `json:"capacityRebalance,omitempty"`

	MixedInstancesPolicy      *MixedInstancesPolicy `json:"mixedInstancesPolicy,omitempty"`
	Status                    ASGStatus
	Instances                 []infrav1.Instance `json:"instances,omitempty"`
	CurrentlySuspendProcesses []string           `json:"currentlySuspendProcesses,omitempty"`
}

// AWSLifecycleHook describes an AWS lifecycle hook
type AWSLifecycleHook struct {
	// The name of the lifecycle hook.
	Name string `json:"name"`

	// The ARN of the notification target that Amazon EC2 Auto Scaling uses to
	// notify you when an instance is in the transition state for the lifecycle hook.
	// +optional
	NotificationTargetARN *string `json:"notificationTargetARN,omitempty"`

	// The ARN of the IAM role that allows the Auto Scaling group to publish to the
	// specified notification target.
	// +optional
	RoleARN *string `json:"roleARN,omitempty"`

	// The state of the EC2 instance to which to attach the lifecycle hook.
	// +kubebuilder:validation:Enum="autoscaling:EC2_INSTANCE_LAUNCHING";"autoscaling:EC2_INSTANCE_TERMINATING"
	LifecycleTransition LifecycleTransition `json:"lifecycleTransition"`

	// The maximum time, in seconds, that an instance can remain in a Pending:Wait or
	// Terminating:Wait state. The maximum is 172800 seconds (48 hours) or 100 times
	// HeartbeatTimeout, whichever is smaller.
	// +optional
	// +kubebuilder:validation:Format=duration
	HeartbeatTimeout *metav1.Duration `json:"heartbeatTimeout,omitempty"`

	// The default result for the lifecycle hook. The possible values are CONTINUE and ABANDON.
	// +optional
	// +kubebuilder:validation:Enum=CONTINUE;ABANDON
	// +kubebuilder:validation:default:=none
	DefaultResult *LifecycleHookDefaultResult `json:"defaultResult,omitempty"`

	// Contains additional metadata that will be passed to the notification target.
	// +optional
	NotificationMetadata *string `json:"notificationMetadata,omitempty"`
}

// LifecycleTransition is the state of the EC2 instance to which to attach the lifecycle hook.
type LifecycleTransition string

const (
	// LifecycleHookTransitionInstanceLaunching is the launching state of the EC2 instance.
	LifecycleHookTransitionInstanceLaunching LifecycleTransition = "autoscaling:EC2_INSTANCE_LAUNCHING"
	// LifecycleHookTransitionInstanceTerminating is the terminating state of the EC2 instance.
	LifecycleHookTransitionInstanceTerminating LifecycleTransition = "autoscaling:EC2_INSTANCE_TERMINATING"
)

func (l LifecycleTransition) String() string {
	return string(l)
}

// LifecycleHookDefaultResult is the default result for the lifecycle hook.
type LifecycleHookDefaultResult string

const (
	// LifecycleHookDefaultResultContinue is the default result for the lifecycle hook to continue.
	LifecycleHookDefaultResultContinue LifecycleHookDefaultResult = "CONTINUE"
	// LifecycleHookDefaultResultAbandon is the default result for the lifecycle hook to abandon.
	LifecycleHookDefaultResultAbandon LifecycleHookDefaultResult = "ABANDON"
)

func (d LifecycleHookDefaultResult) String() string {
	return string(d)
}

// ASGStatus is a status string returned by the autoscaling API.
type ASGStatus string

// ASGStatusDeleteInProgress is the string representing an ASG that is currently deleting.
var ASGStatusDeleteInProgress = ASGStatus("Delete in progress")

// TaintEffect is the effect for a Kubernetes taint.
type TaintEffect string

var (
	// TaintEffectNoSchedule is a taint that indicates that a pod shouldn't be scheduled on a node
	// unless it can tolerate the taint.
	TaintEffectNoSchedule = TaintEffect("no-schedule")
	// TaintEffectNoExecute is a taint that indicates that a pod shouldn't be schedule on a node
	// unless it can tolerate it. And if its already running on the node it will be evicted.
	TaintEffectNoExecute = TaintEffect("no-execute")
	// TaintEffectPreferNoSchedule is a taint that indicates that there is a "preference" that pods shouldn't
	// be scheduled on a node unless it can tolerate the taint. the scheduler will try to avoid placing the pod
	// but it may still run on the node if there is no other option.
	TaintEffectPreferNoSchedule = TaintEffect("prefer-no-schedule")
)

// Taint defines the specs for a Kubernetes taint.
type Taint struct {
	// Effect specifies the effect for the taint
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=no-schedule;no-execute;prefer-no-schedule
	Effect TaintEffect `json:"effect"`
	// Key is the key of the taint
	// +kubebuilder:validation:Required
	Key string `json:"key"`
	// Value is the value of the taint
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// Equals is used to test if 2 taints are equal.
func (t *Taint) Equals(other *Taint) bool {
	if t == nil || other == nil {
		return t == other
	}

	return t.Effect == other.Effect &&
		t.Key == other.Key &&
		t.Value == other.Value
}

// Taints is an array of Taints.
type Taints []Taint

// Contains checks for existence of a matching taint.
func (t *Taints) Contains(taint *Taint) bool {
	for _, t := range *t {
		if t.Equals(taint) {
			return true
		}
	}

	return false
}

// UpdateConfig is the configuration options for updating a nodegroup. Only one of MaxUnavailable
// and MaxUnavailablePercentage should be specified.
type UpdateConfig struct {
	// MaxUnavailable is the maximum number of nodes unavailable at once during a version update.
	// Nodes will be updated in parallel. The maximum number is 100.
	// +optional
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	MaxUnavailable *int `json:"maxUnavailable,omitempty"`

	// MaxUnavailablePercentage is the maximum percentage of nodes unavailable during a version update. This
	// percentage of nodes will be updated in parallel, up to 100 nodes at once.
	// +optional
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	MaxUnavailablePercentage *int `json:"maxUnavailablePercentage,omitempty"`
}

// AZSubnetType is the type of subnet to use when an availability zone is specified.
type AZSubnetType string

const (
	// AZSubnetTypePublic is a public subnet.
	AZSubnetTypePublic AZSubnetType = "public"
	// AZSubnetTypePrivate is a private subnet.
	AZSubnetTypePrivate AZSubnetType = "private"
	// AZSubnetTypeAll is all subnets in an availability zone.
	AZSubnetTypeAll AZSubnetType = "all"
)

// NewAZSubnetType returns a pointer to an AZSubnetType.
func NewAZSubnetType(t AZSubnetType) *AZSubnetType {
	return &t
}
