/*
Copyright 2021 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

// BootstrapUser contains a list of elements that is specific
// to the configuration and enablement of an IAM user.
type BootstrapUser struct {

	// Enable controls whether or not a bootstrap AWS IAM user will be created.
	// This can be used to scope down the initial credentials used to bootstrap the
	// cluster.
	// Defaults to false.
	Enable bool `json:"enable"`

	// UserName controls the username of the bootstrap user. Defaults to
	// "bootstrapper.cluster-api-provider-aws.sigs.k8s.io"
	UserName string `json:"userName,omitempty"`

	// GroupName controls the group the user will belong to. Defaults to
	// "bootstrapper.cluster-api-provider-aws.sigs.k8s.io"
	GroupName string `json:"groupName,omitempty"`

	// ExtraPolicyAttachments is a list of additional policies to be attached to the IAM user.
	ExtraPolicyAttachments []string `json:"extraPolicyAttachments,omitempty"`

	// ExtraGroups is a list of groups to add this user to.
	ExtraGroups []string `json:"extraGroups,omitempty"`

	// ExtraStatements are additional AWS IAM policy document statements to be included inline for the user.
	ExtraStatements []iamv1.StatementEntry `json:"extraStatements,omitempty"`

	// Tags is a map of tags to be applied to the AWS IAM user.
	Tags infrav1.Tags `json:"tags,omitempty"`
}

// ControlPlane controls the configuration of the AWS IAM role for
// the control plane of provisioned Kubernetes clusters.
type ControlPlane struct {
	AWSIAMRoleSpec `json:",inline"`

	// DisableClusterAPIControllerPolicyAttachment, if set to true, will not attach the AWS IAM policy for Cluster
	// API Provider AWS to the control plane role. Defaults to false.
	DisableClusterAPIControllerPolicyAttachment bool `json:"disableClusterAPIControllerPolicyAttachment,omitempty"`

	// DisableCloudProviderPolicy if set to true, will not generate and attach the AWS IAM policy for the AWS Cloud Provider.
	DisableCloudProviderPolicy bool `json:"disableCloudProviderPolicy"`

	// EnableCSIPolicy if set to true, will generate and attach the AWS IAM policy for the EBS CSI Driver.
	EnableCSIPolicy bool `json:"enableCSIPolicy"`
}

// AWSIAMRoleSpec defines common configuration for AWS IAM roles created by
// Kubernetes Cluster API Provider AWS.
type AWSIAMRoleSpec struct {
	// Disable if set to true will not create the AWS IAM role. Defaults to false.
	Disable bool `json:"disable"` // default: false

	// ExtraPolicyAttachments is a list of additional policies to be attached to the IAM role.
	ExtraPolicyAttachments []string `json:"extraPolicyAttachments,omitempty"`

	// ExtraStatements are additional IAM statements to be included inline for the role.
	ExtraStatements []iamv1.StatementEntry `json:"extraStatements,omitempty"`

	// Path sets the path to the role.
	// +optional
	Path string `json:"path,omitempty"`

	// PermissionsBoundary sets the ARN of the managed policy that is used to set the permissions boundary for the role.
	// +optional
	PermissionsBoundary string `json:"permissionsBoundary,omitempty"`

	// TrustStatements is an IAM PolicyDocument defining what identities are allowed to assume this role.
	// See "sigs.k8s.io/cluster-api-provider-aws/v2/cmd/clusterawsadm/api/iam/v1beta1" for more documentation.
	TrustStatements []iamv1.StatementEntry `json:"trustStatements,omitempty"`

	// Tags is a map of tags to be applied to the AWS IAM role.
	Tags infrav1.Tags `json:"tags,omitempty"`
}

// EKSConfig represents the EKS related configuration config.
type EKSConfig struct {
	// Disable controls whether EKS-related permissions are granted
	Disable bool `json:"disable"`
	// AllowIAMRoleCreation controls whether the EKS controllers have permissions for creating IAM
	// roles per cluster
	AllowIAMRoleCreation bool `json:"iamRoleCreation,omitempty"`
	// EnableUserEKSConsolePolicy controls the creation of the policy to view EKS nodes and workloads.
	EnableUserEKSConsolePolicy bool `json:"enableUserEKSConsolePolicy,omitempty"`
	// DefaultControlPlaneRole controls the configuration of the AWS IAM role for
	// the EKS control plane. This is the default role that will be used if
	// no role is included in the spec and automatic creation of the role
	// isn't enabled
	DefaultControlPlaneRole AWSIAMRoleSpec `json:"defaultControlPlaneRole,omitempty"`
	// ManagedMachinePool controls the configuration of the AWS IAM role for
	// used by EKS managed machine pools.
	ManagedMachinePool *AWSIAMRoleSpec `json:"managedMachinePool,omitempty"`
	// Fargate controls the configuration of the AWS IAM role for
	// used by EKS managed machine pools.
	Fargate *AWSIAMRoleSpec `json:"fargate,omitempty"`
	// KMSAliasPrefix is prefix to use to restrict permission to KMS keys to only those that have an alias
	// name that is prefixed by this.
	// Defaults to cluster-api-provider-aws-*
	KMSAliasPrefix string `json:"kmsAliasPrefix,omitempty"`
}

// EventBridgeConfig represents configuration for enabling experimental feature to consume
// EventBridge EC2 events.
type EventBridgeConfig struct {
	// Enable controls whether permissions are granted to consume EC2 events
	Enable bool `json:"enable,omitempty"`
}

// ClusterAPIControllers controls the configuration of the AWS IAM role for
// the Kubernetes Cluster API Provider AWS controller.
type ClusterAPIControllers struct {
	AWSIAMRoleSpec `json:",inline"`
	// AllowedEC2InstanceProfiles controls which EC2 roles are allowed to be
	// consumed by Cluster API when creating an ec2 instance. Defaults to
	// *.<suffix>, where suffix is defaulted to .cluster-api-provider-aws.sigs.k8s.io
	AllowedEC2InstanceProfiles []string `json:"allowedEC2InstanceProfiles,omitempty"`
}

// Nodes controls the configuration of the AWS IAM role for worker nodes
// in a cluster created by Kubernetes Cluster API Provider AWS.
type Nodes struct {
	AWSIAMRoleSpec `json:",inline"`

	// DisableCloudProviderPolicy if set to true, will not generate and attach the policy for the AWS Cloud Provider.
	// Defaults to false.
	DisableCloudProviderPolicy bool `json:"disableCloudProviderPolicy"`

	// EC2ContainerRegistryReadOnly controls whether the node has read-only access to the
	// EC2 container registry
	EC2ContainerRegistryReadOnly bool `json:"ec2ContainerRegistryReadOnly"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSIAMConfiguration controls the creation of AWS Identity and Access Management (IAM) resources for use
// by Kubernetes clusters and Kubernetes Cluster API Provider AWS.
type AWSIAMConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	Spec AWSIAMConfigurationSpec `json:"spec,omitempty"`
}

// S3Buckets controls the configuration of the AWS IAM role for S3 buckets
// which can be created for storing bootstrap data for nodes requiring it.
type S3Buckets struct {
	// Enable controls whether permissions are granted to manage S3 buckets.
	Enable bool `json:"enable"`

	// NamePrefix will be prepended to every AWS IAM role bucket name. Defaults to "cluster-api-provider-aws-".
	// AWSCluster S3 Bucket name must be prefixed with the same prefix.
	NamePrefix string `json:"namePrefix"`
}

// AWSIAMConfigurationSpec defines the specification of the AWSIAMConfiguration.
type AWSIAMConfigurationSpec struct {
	// NamePrefix will be prepended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to "".
	NamePrefix string `json:"namePrefix,omitempty"`

	// NameSuffix will be appended to every AWS IAM role, user and policy created by clusterawsadm. Defaults to
	// ".cluster-api-provider-aws.sigs.k8s.io".
	NameSuffix *string `json:"nameSuffix,omitempty"`

	// ControlPlane controls the configuration of the AWS IAM role for a Kubernetes cluster's control plane nodes.
	ControlPlane ControlPlane `json:"controlPlane,omitempty"`

	// ClusterAPIControllers controls the configuration of an IAM role and policy specifically for Kubernetes Cluster API Provider AWS.
	ClusterAPIControllers ClusterAPIControllers `json:"clusterAPIControllers,omitempty"`

	// Nodes controls the configuration of the AWS IAM role for all nodes in a Kubernetes cluster.
	Nodes Nodes `json:"nodes,omitempty"`

	// BootstrapUser contains a list of elements that is specific
	// to the configuration and enablement of an IAM user.
	BootstrapUser BootstrapUser `json:"bootstrapUser,omitempty"`

	// StackName defines the name of the AWS CloudFormation stack.
	StackName string `json:"stackName,omitempty"`

	// StackTags defines the tags of the AWS CloudFormation stack.
	// +optional
	StackTags map[string]string `json:"stackTags,omitempty"`

	// Region controls which region the control-plane is created in if not specified on the command line or
	// via environment variables.
	Region string `json:"region,omitempty"`

	// EKS controls the configuration related to EKS. Settings in here affect the control plane
	// and nodes roles
	EKS *EKSConfig `json:"eks,omitempty"`

	// EventBridge controls configuration for consuming EventBridge events
	EventBridge *EventBridgeConfig `json:"eventBridge,omitempty"`

	// Partition is the AWS security partition being used. Defaults to "aws"
	Partition string `json:"partition,omitempty"`

	// SecureSecretsBackend, when set to parameter-store will create AWS Systems Manager
	// Parameter Storage policies. By default or with the value of secrets-manager,
	// will generate AWS Secrets Manager policies instead.
	// +kubebuilder:validation:Enum=secrets-manager;ssm-parameter-store
	SecureSecretsBackends []infrav1.SecretBackend `json:"secureSecretBackends,omitempty"`

	// S3Buckets, when enabled, will add controller nodes permissions to
	// create S3 Buckets for workload clusters.
	// TODO: This field could be a pointer, but it seems it breaks setting default values?
	// +optional
	S3Buckets S3Buckets `json:"s3Buckets,omitempty"`

	// AllowAssumeRole enables the sts:AssumeRole permission within the CAPA policies
	AllowAssumeRole bool `json:"allowAssumeRole,omitempty"`
}

// GetObjectKind returns the AAWSIAMConfiguration's TypeMeta.
func (obj *AWSIAMConfiguration) GetObjectKind() schema.ObjectKind {
	return &obj.TypeMeta
}

// NewAWSIAMConfiguration will generate a new default AWSIAMConfiguration.
func NewAWSIAMConfiguration() *AWSIAMConfiguration {
	conf := &AWSIAMConfiguration{}
	SetObjectDefaults_AWSIAMConfiguration(conf)
	return conf
}
