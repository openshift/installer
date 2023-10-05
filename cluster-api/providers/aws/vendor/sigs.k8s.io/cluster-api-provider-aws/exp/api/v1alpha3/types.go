/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
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

// AWSLaunchTemplate defines the desired state of AWSLaunchTemplate
type AWSLaunchTemplate struct {
	// The name of the launch template.
	Name string `json:"name,omitempty"`

	// The name or the Amazon Resource Name (ARN) of the instance profile associated
	// with the IAM role for the instance. The instance profile contains the IAM
	// role.
	IamInstanceProfile string `json:"iamInstanceProfile,omitempty"`

	// AMI is the reference to the AMI from which to create the machine instance.
	// +optional
	AMI infrav1alpha3.AWSResourceReference `json:"ami,omitempty"`

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
	RootVolume *infrav1alpha3.Volume `json:"rootVolume,omitempty"`

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
	AdditionalSecurityGroups []infrav1alpha3.AWSResourceReference `json:"additionalSecurityGroups,omitempty"`
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
)

// InstancesDistribution to configure distribution of On-Demand Instances and Spot Instances.
type InstancesDistribution struct {
	// +kubebuilder:validation:Enum=prioritized
	// +kubebuilder:default=prioritized
	OnDemandAllocationStrategy OnDemandAllocationStrategy `json:"onDemandAllocationStrategy,omitempty"`

	// +kubebuilder:validation:Enum=lowest-price;capacity-optimized
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
	ID                string             `json:"id,omitempty"`
	Tags              infrav1alpha3.Tags `json:"tags,omitempty"`
	Name              string             `json:"name,omitempty"`
	DesiredCapacity   *int32             `json:"desiredCapacity,omitempty"`
	MaxSize           int32              `json:"maxSize,omitempty"`
	MinSize           int32              `json:"minSize,omitempty"`
	PlacementGroup    string             `json:"placementGroup,omitempty"`
	Subnets           []string           `json:"subnets,omitempty"`
	DefaultCoolDown   metav1.Duration    `json:"defaultCoolDown,omitempty"`
	CapacityRebalance bool               `json:"capacityRebalance,omitempty"`

	MixedInstancesPolicy *MixedInstancesPolicy `json:"mixedInstancesPolicy,omitempty"`
	Status               ASGStatus
	Instances            []infrav1alpha3.Instance `json:"instances,omitempty"`
}

// ASGStatus is a status string returned by the autoscaling API
type ASGStatus string

var (
	// ASGStatusDeleteInProgress is the string representing an ASG that is currently deleting.
	ASGStatusDeleteInProgress = ASGStatus("Delete in progress")
)
