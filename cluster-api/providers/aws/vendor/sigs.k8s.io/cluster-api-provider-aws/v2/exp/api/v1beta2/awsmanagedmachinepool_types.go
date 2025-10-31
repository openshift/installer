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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ManagedMachineAMIType specifies which AWS AMI to use for a managed MachinePool.
// Source of truth can be found using the link below:
// https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateNodegroup.html#AmazonEKS-CreateNodegroup-request-amiType
type ManagedMachineAMIType string

const (
	// Al2x86_64 is the default AMI type.
	Al2x86_64 ManagedMachineAMIType = "AL2_x86_64"
	// Al2x86_64GPU is the x86-64 GPU AMI type.
	Al2x86_64GPU ManagedMachineAMIType = "AL2_x86_64_GPU"
	// Al2Arm64 is the Arm AMI type.
	Al2Arm64 ManagedMachineAMIType = "AL2_ARM_64"
	// Custom is the custom AMI type.
	Custom ManagedMachineAMIType = "CUSTOM"
	// BottleRocketArm64 is the Arm AMI type.
	BottleRocketArm64 ManagedMachineAMIType = "BOTTLEROCKET_ARM_64"
	// BottleRocketx86_64 is the BottleRocket x86-64 AMI type.
	BottleRocketx86_64 ManagedMachineAMIType = "BOTTLEROCKET_x86_64"
	// BottleRocketArm64Fips is the BottleRocket Arm Fips AMI type.
	BottleRocketArm64Fips ManagedMachineAMIType = "BOTTLEROCKET_ARM_64_FIPS"
	// BottleRocketx86_64Fips is the BottleRocket x86-64 Fips AMI type.
	BottleRocketx86_64Fips ManagedMachineAMIType = "BOTTLEROCKET_x86_64_FIPS"
	// BottleRocketArm64Nvidia is the BottleRocket Arm Nvidia AMI type.
	BottleRocketArm64Nvidia ManagedMachineAMIType = "BOTTLEROCKET_ARM_64_NVIDIA"
	// BottleRocketx86_64Nvidia is the BottleRocket x86-64 Nvidia AMI type.
	BottleRocketx86_64Nvidia ManagedMachineAMIType = "BOTTLEROCKET_x86_64_NVIDIA"
	// WindowsCore2019x86_64 is the Windows Core 2019 x86-64 AMI type.
	WindowsCore2019x86_64 ManagedMachineAMIType = "WINDOWS_CORE_2019_x86_64"
	// WindowsFull2019x86_64 is the Windows Full 2019 x86-64 AMI type.
	WindowsFull2019x86_64 ManagedMachineAMIType = "WINDOWS_FULL_2019_x86_64"
	// WindowsCore2022x86_64 is the Windows Core 2022 x86-64 AMI type.
	WindowsCore2022x86_64 ManagedMachineAMIType = "WINDOWS_CORE_2022_x86_64"
	// WindowsFull2022x86_64 is the Windows Full 2022 x86-64 AMI type.
	WindowsFull2022x86_64 ManagedMachineAMIType = "WINDOWS_FULL_2022_x86_64"
	// Al2023x86_64 is the AL2023 x86-64 AMI type.
	Al2023x86_64 ManagedMachineAMIType = "AL2023_x86_64_STANDARD"
	// Al2023Arm64 is the AL2023 Arm AMI type.
	Al2023Arm64 ManagedMachineAMIType = "AL2023_ARM_64_STANDARD"
	// Al2023x86_64Neuron is the AL2023 x86-64 Neuron AMI type.
	Al2023x86_64Neuron ManagedMachineAMIType = "AL2023_x86_64_NEURON"
	// Al2023x86_64Nvidia is the AL2023 x86-64 Nvidia AMI type.
	Al2023x86_64Nvidia ManagedMachineAMIType = "AL2023_x86_64_NVIDIA"
	// Al2023Arm64Nvidia is the AL2023 Arm Nvidia AMI type.
	Al2023Arm64Nvidia ManagedMachineAMIType = "AL2023_ARM_64_NVIDIA"
)

// ManagedMachinePoolCapacityType specifies the capacity type to be used for the managed MachinePool.
type ManagedMachinePoolCapacityType string

const (
	// ManagedMachinePoolCapacityTypeOnDemand is the default capacity type, to launch on-demand instances.
	ManagedMachinePoolCapacityTypeOnDemand ManagedMachinePoolCapacityType = "onDemand"
	// ManagedMachinePoolCapacityTypeSpot is the spot instance capacity type to launch spot instances.
	ManagedMachinePoolCapacityTypeSpot ManagedMachinePoolCapacityType = "spot"
)

var (
	// DefaultEKSNodegroupRole is the name of the default IAM role to use for EKS nodegroups
	// if no other role is supplied in the spec and if iam role creation is not enabled. The default
	// can be created using clusterawsadm or created manually.
	DefaultEKSNodegroupRole = fmt.Sprintf("eks-nodegroup%s", iamv1.DefaultNameSuffix)
)

// AWSManagedMachinePoolSpec defines the desired state of AWSManagedMachinePool.
type AWSManagedMachinePoolSpec struct {
	// EKSNodegroupName specifies the name of the nodegroup in AWS
	// corresponding to this MachinePool. If you don't specify a name
	// then a default name will be created based on the namespace and
	// name of the managed machine pool.
	// +optional
	EKSNodegroupName string `json:"eksNodegroupName,omitempty"`

	// AvailabilityZones is an array of availability zones instances can run in
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// AvailabilityZoneSubnetType specifies which type of subnets to use when an availability zone is specified.
	// +kubebuilder:validation:Enum:=public;private;all
	// +optional
	AvailabilityZoneSubnetType *AZSubnetType `json:"availabilityZoneSubnetType,omitempty"`

	// SubnetIDs specifies which subnets are used for the
	// auto scaling group of this nodegroup
	// +optional
	SubnetIDs []string `json:"subnetIDs,omitempty"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// RoleAdditionalPolicies allows you to attach additional polices to
	// the node group role. You must enable the EKSAllowAddRoles
	// feature flag to incorporate these into the created role.
	// +optional
	RoleAdditionalPolicies []string `json:"roleAdditionalPolicies,omitempty"`

	// RoleName specifies the name of IAM role for the node group.
	// If the role is pre-existing we will treat it as unmanaged
	// and not delete it on deletion. If the EKSEnableIAM feature
	// flag is true and no name is supplied then a role is created.
	// +optional
	RoleName string `json:"roleName,omitempty"`

	// RolePath sets the path to the role. For more information about paths, see IAM Identifiers
	// (https://docs.aws.amazon.com/IAM/latest/UserGuide/Using_Identifiers.html)
	// in the IAM User Guide.
	//
	// This parameter is optional. If it is not included, it defaults to a slash
	// (/).
	RolePath string `json:"rolePath,omitempty"`

	// RolePermissionsBoundary sets the ARN of the managed policy that is used
	// to set the permissions boundary for the role.
	//
	// A permissions boundary policy defines the maximum permissions that identity-based
	// policies can grant to an entity, but does not grant permissions. Permissions
	// boundaries do not define the maximum permissions that a resource-based policy
	// can grant to an entity. To learn more, see Permissions boundaries for IAM
	// entities (https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_boundaries.html)
	// in the IAM User Guide.
	//
	// For more information about policy types, see Policy types (https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html#access_policy-types)
	// in the IAM User Guide.
	RolePermissionsBoundary string `json:"rolePermissionsBoundary,omitempty"`

	// AMIVersion defines the desired AMI release version. If no version number
	// is supplied then the latest version for the Kubernetes version
	// will be used
	// +kubebuilder:validation:MinLength:=2
	// +optional
	AMIVersion *string `json:"amiVersion,omitempty"`

	// AMIType defines the AMI type
	// +kubebuilder:validation:Enum:=AL2_x86_64;AL2_x86_64_GPU;AL2_ARM_64;CUSTOM;BOTTLEROCKET_ARM_64;BOTTLEROCKET_x86_64;BOTTLEROCKET_ARM_64_FIPS;BOTTLEROCKET_x86_64_FIPS;BOTTLEROCKET_ARM_64_NVIDIA;BOTTLEROCKET_x86_64_NVIDIA;WINDOWS_CORE_2019_x86_64;WINDOWS_FULL_2019_x86_64;WINDOWS_CORE_2022_x86_64;WINDOWS_FULL_2022_x86_64;AL2023_x86_64_STANDARD;AL2023_ARM_64_STANDARD;AL2023_x86_64_NEURON;AL2023_x86_64_NVIDIA;AL2023_ARM_64_NVIDIA
	// +kubebuilder:default:=AL2_x86_64
	// +optional
	AMIType *ManagedMachineAMIType `json:"amiType,omitempty"`

	// Labels specifies labels for the Kubernetes node objects
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Taints specifies the taints to apply to the nodes of the machine pool
	// +optional
	Taints Taints `json:"taints,omitempty"`

	// DiskSize specifies the root disk size
	// +optional
	DiskSize *int32 `json:"diskSize,omitempty"`

	// InstanceType specifies the AWS instance type
	// +optional
	InstanceType *string `json:"instanceType,omitempty"`

	// Scaling specifies scaling for the ASG behind this pool
	// +optional
	Scaling *ManagedMachinePoolScaling `json:"scaling,omitempty"`

	// RemoteAccess specifies how machines can be accessed remotely
	// +optional
	RemoteAccess *ManagedRemoteAccess `json:"remoteAccess,omitempty"`

	// ProviderIDList are the provider IDs of instances in the
	// autoscaling group corresponding to the nodegroup represented by this
	// machine pool
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`

	// CapacityType specifies the capacity type for the ASG behind this pool
	// +kubebuilder:validation:Enum:=onDemand;spot
	// +kubebuilder:default:=onDemand
	// +optional
	CapacityType *ManagedMachinePoolCapacityType `json:"capacityType,omitempty"`

	// UpdateConfig holds the optional config to control the behaviour of the update
	// to the nodegroup.
	// +optional
	UpdateConfig *UpdateConfig `json:"updateConfig,omitempty"`

	// AWSLaunchTemplate specifies the launch template to use to create the managed node group.
	// If AWSLaunchTemplate is specified, certain node group configuraions outside of launch template
	// are prohibited (https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html).
	// +optional
	AWSLaunchTemplate *AWSLaunchTemplate `json:"awsLaunchTemplate,omitempty"`

	// AWSLifecycleHooks specifies lifecycle hooks for the managed node group.
	// +optional
	AWSLifecycleHooks []AWSLifecycleHook `json:"lifecycleHooks,omitempty"`
}

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	MinSize *int32 `json:"minSize,omitempty"`
	MaxSize *int32 `json:"maxSize,omitempty"`
}

// ManagedRemoteAccess specifies remote access settings for EC2 instances.
type ManagedRemoteAccess struct {
	// SSHKeyName specifies which EC2 SSH key can be used to access machines.
	// If left empty, the key from the control plane is used.
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// SourceSecurityGroups specifies which security groups are allowed access
	SourceSecurityGroups []string `json:"sourceSecurityGroups,omitempty"`

	// Public specifies whether to open port 22 to the public internet
	Public bool `json:"public,omitempty"`
}

// AWSManagedMachinePoolStatus defines the observed state of AWSManagedMachinePool.
type AWSManagedMachinePoolStatus struct {
	// Ready denotes that the AWSManagedMachinePool nodegroup has joined
	// the cluster
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`

	// The ID of the launch template
	// +optional
	LaunchTemplateID *string `json:"launchTemplateID,omitempty"`

	// The version of the launch template
	// +optional
	LaunchTemplateVersion *string `json:"launchTemplateVersion,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the MachinePool and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of MachinePools
	// can be added as events to the MachinePool object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *string `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the MachinePool and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the MachinePool's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of MachinePools
	// can be added as events to the MachinePool object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// Conditions defines current service state of the managed machine pool
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedmachinepools,scope=Namespaced,categories=cluster-api,shortName=awsmmp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="MachinePool ready status"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="Number of replicas"

// AWSManagedMachinePool is the Schema for the awsmanagedmachinepools API.
type AWSManagedMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSManagedMachinePoolSpec   `json:"spec,omitempty"`
	Status AWSManagedMachinePoolStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the AWSManagedMachinePool resource.
func (r *AWSManagedMachinePool) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the AWSManagedMachinePool to the predescribed clusterv1.Conditions.
func (r *AWSManagedMachinePool) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// AWSManagedMachinePoolList contains a list of AWSManagedMachinePools.
type AWSManagedMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedMachinePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSManagedMachinePool{}, &AWSManagedMachinePoolList{})
}
