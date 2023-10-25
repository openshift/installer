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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ManagedControlPlaneFinalizer allows the controller to clean up resources on delete.
	ManagedControlPlaneFinalizer = "awsmanagedcontrolplane.controlplane.cluster.x-k8s.io"

	// AWSManagedControlPlaneKind is the Kind of AWSManagedControlPlane.
	AWSManagedControlPlaneKind = "AWSManagedControlPlane"
)

// AWSManagedControlPlaneSpec defines the desired state of an Amazon EKS Cluster.
type AWSManagedControlPlaneSpec struct { //nolint: maligned
	// EKSClusterName allows you to specify the name of the EKS cluster in
	// AWS. If you don't specify a name then a default name will be created
	// based on the namespace and name of the managed control plane.
	// +optional
	EKSClusterName string `json:"eksClusterName,omitempty"`

	// IdentityRef is a reference to a identity to be used when reconciling the managed control plane.
	// +optional
	IdentityRef *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`

	// NetworkSpec encapsulates all things related to AWS network.
	NetworkSpec infrav1.NetworkSpec `json:"network,omitempty"`

	// SecondaryCidrBlock is the additional CIDR range to use for pod IPs.
	// Must be within the 100.64.0.0/10 or 198.19.0.0/16 range.
	// +optional
	SecondaryCidrBlock *string `json:"secondaryCidrBlock,omitempty"`

	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// Partition is the AWS security partition being used. Defaults to "aws"
	// +optional
	Partition string `json:"partition,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// Version defines the desired Kubernetes version. If no version number
	// is supplied then the latest version of Kubernetes that EKS supports
	// will be used.
	// +kubebuilder:validation:MinLength:=2
	// +kubebuilder:validation:Pattern:=^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.?(\.0|[1-9][0-9]*)?$
	// +optional
	Version *string `json:"version,omitempty"`

	// RoleName specifies the name of IAM role that gives EKS
	// permission to make API calls. If the role is pre-existing
	// we will treat it as unmanaged and not delete it on
	// deletion. If the EKSEnableIAM feature flag is true
	// and no name is supplied then a role is created.
	// +kubebuilder:validation:MinLength:=2
	// +optional
	RoleName *string `json:"roleName,omitempty"`

	// RoleAdditionalPolicies allows you to attach additional polices to
	// the control plane role. You must enable the EKSAllowAddRoles
	// feature flag to incorporate these into the created role.
	// +optional
	RoleAdditionalPolicies *[]string `json:"roleAdditionalPolicies,omitempty"`

	// Logging specifies which EKS Cluster logs should be enabled. Entries for
	// each of the enabled logs will be sent to CloudWatch
	// +optional
	Logging *ControlPlaneLoggingSpec `json:"logging,omitempty"`

	// EncryptionConfig specifies the encryption configuration for the cluster
	// +optional
	EncryptionConfig *EncryptionConfig `json:"encryptionConfig,omitempty"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// IAMAuthenticatorConfig allows the specification of any additional user or role mappings
	// for use when generating the aws-iam-authenticator configuration. If this is nil the
	// default configuration is still generated for the cluster.
	// +optional
	IAMAuthenticatorConfig *IAMAuthenticatorConfig `json:"iamAuthenticatorConfig,omitempty"`

	// Endpoints specifies access to this cluster's control plane endpoints
	// +optional
	EndpointAccess EndpointAccess `json:"endpointAccess,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// ImageLookupFormat is the AMI naming format to look up machine images when
	// a machine does not specify an AMI. When set, this will be used for all
	// cluster machines unless a machine specifies a different ImageLookupOrg.
	// Supports substitutions for {{.BaseOS}} and {{.K8sVersion}} with the base
	// OS and kubernetes version, respectively. The BaseOS will be the value in
	// ImageLookupBaseOS or ubuntu (the default), and the kubernetes version as
	// defined by the packages produced by kubernetes/release without v as a
	// prefix: 1.13.0, 1.12.5-mybuild.1, or 1.17.3. For example, the default
	// image format of capa-ami-{{.BaseOS}}-?{{.K8sVersion}}-* will end up
	// searching for AMIs that match the pattern capa-ami-ubuntu-?1.18.0-* for a
	// Machine that is targeting kubernetes v1.18.0 and the ubuntu base OS. See
	// also: https://golang.org/pkg/text/template/
	// +optional
	ImageLookupFormat string `json:"imageLookupFormat,omitempty"`

	// ImageLookupOrg is the AWS Organization ID to look up machine images when a
	// machine does not specify an AMI. When set, this will be used for all
	// cluster machines unless a machine specifies a different ImageLookupOrg.
	// +optional
	ImageLookupOrg string `json:"imageLookupOrg,omitempty"`

	// ImageLookupBaseOS is the name of the base operating system used to look
	// up machine images when a machine does not specify an AMI. When set, this
	// will be used for all cluster machines unless a machine specifies a
	// different ImageLookupBaseOS.
	ImageLookupBaseOS string `json:"imageLookupBaseOS,omitempty"`

	// Bastion contains options to configure the bastion host.
	// +optional
	Bastion infrav1.Bastion `json:"bastion"`

	// TokenMethod is used to specify the method for obtaining a client token for communicating with EKS
	// iam-authenticator - obtains a client token using iam-authentictor
	// aws-cli - obtains a client token using the AWS CLI
	// Defaults to iam-authenticator
	// +kubebuilder:default=iam-authenticator
	// +kubebuilder:validation:Enum=iam-authenticator;aws-cli
	TokenMethod *EKSTokenMethod `json:"tokenMethod,omitempty"`

	// AssociateOIDCProvider can be enabled to automatically create an identity
	// provider for the controller for use with IAM roles for service accounts
	// +kubebuilder:default=false
	AssociateOIDCProvider bool `json:"associateOIDCProvider,omitempty"`

	// Addons defines the EKS addons to enable with the EKS cluster.
	// +optional
	Addons *[]Addon `json:"addons,omitempty"`

	// IdentityProviderconfig is used to specify the oidc provider config
	// to be attached with this eks cluster
	// +optional
	OIDCIdentityProviderConfig *OIDCIdentityProviderConfig `json:"oidcIdentityProviderConfig,omitempty"`

	// VpcCni is used to set configuration options for the VPC CNI plugin
	// +optional
	VpcCni VpcCni `json:"vpcCni,omitempty"`

	// KubeProxy defines managed attributes of the kube-proxy daemonset
	KubeProxy KubeProxy `json:"kubeProxy,omitempty"`
}

// KubeProxy specifies how the kube-proxy daemonset is managed.
type KubeProxy struct {
	// Disable set to true indicates that kube-proxy should be disabled. With EKS clusters
	// kube-proxy is automatically installed into the cluster. For clusters where you want
	// to use kube-proxy functionality that is provided with an alternate CNI, this option
	// provides a way to specify that the kube-proxy daemonset should be deleted. You cannot
	// set this to true if you are using the Amazon kube-proxy addon.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`
}

// VpcCni specifies configuration related to the VPC CNI.
type VpcCni struct {
	// Disable indicates that the Amazon VPC CNI should be disabled. With EKS clusters the
	// Amazon VPC CNI is automatically installed into the cluster. For clusters where you want
	// to use an alternate CNI this option provides a way to specify that the Amazon VPC CNI
	// should be deleted. You cannot set this to true if you are using the
	// Amazon VPC CNI addon.
	// +kubebuilder:default=false
	Disable bool `json:"disable,omitempty"`
	// Env defines a list of environment variables to apply to the `aws-node` DaemonSet
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// EndpointAccess specifies how control plane endpoints are accessible.
type EndpointAccess struct {
	// Public controls whether control plane endpoints are publicly accessible
	// +optional
	Public *bool `json:"public,omitempty"`
	// PublicCIDRs specifies which blocks can access the public endpoint
	// +optional
	PublicCIDRs []*string `json:"publicCIDRs,omitempty"`
	// Private points VPC-internal control plane access to the private endpoint
	// +optional
	Private *bool `json:"private,omitempty"`
}

// EncryptionConfig specifies the encryption configuration for the EKS clsuter.
type EncryptionConfig struct {
	// Provider specifies the ARN or alias of the CMK (in AWS KMS)
	Provider *string `json:"provider,omitempty"`
	// Resources specifies the resources to be encrypted
	Resources []*string `json:"resources,omitempty"`
}

// OIDCProviderStatus holds the status of the AWS OIDC identity provider.
type OIDCProviderStatus struct {
	// ARN holds the ARN of the provider
	ARN string `json:"arn,omitempty"`
	// TrustPolicy contains the boilerplate IAM trust policy to use for IRSA
	TrustPolicy string `json:"trustPolicy,omitempty"`
}

type IdentityProviderStatus struct {
	// ARN holds the ARN of associated identity provider
	ARN string `json:"arn,omitempty"`

	// Status holds current status of associated identity provider
	Status string `json:"status,omitempty"`
}

// AWSManagedControlPlaneStatus defines the observed state of an Amazon EKS Cluster.
type AWSManagedControlPlaneStatus struct {
	// Networks holds details about the AWS networking resources used by the control plane
	// +optional
	Network infrav1.NetworkStatus `json:"networkStatus,omitempty"`
	// FailureDomains specifies a list fo available availability zones that can be used
	// +optional
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
	// Bastion holds details of the instance that is used as a bastion jump box
	// +optional
	Bastion *infrav1.Instance `json:"bastion,omitempty"`
	// OIDCProvider holds the status of the identity provider for this cluster
	// +optional
	OIDCProvider OIDCProviderStatus `json:"oidcProvider,omitempty"`
	// ExternalManagedControlPlane indicates to cluster-api that the control plane
	// is managed by an external service such as AKS, EKS, GKE, etc.
	// +kubebuilder:default=true
	ExternalManagedControlPlane *bool `json:"externalManagedControlPlane,omitempty"`
	// Initialized denotes whether or not the control plane has the
	// uploaded kubernetes config-map.
	// +optional
	Initialized bool `json:"initialized"`
	// Ready denotes that the AWSManagedControlPlane API Server is ready to
	// receive requests and that the VPC infra is ready.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
	// ErrorMessage indicates that there is a terminal problem reconciling the
	// state, and will be set to a descriptive error message.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
	// Conditions specifies the cpnditions for the managed control plane
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
	// Addons holds the current status of the EKS addons
	// +optional
	Addons []AddonState `json:"addons,omitempty"`
	// IdentityProviderStatus holds the status for
	// associated identity provider
	// +optional
	IdentityProviderStatus IdentityProviderStatus `json:"identityProviderStatus,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedcontrolplanes,shortName=awsmcp,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AWSManagedControl belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Control plane infrastructure is ready for worker nodes"
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.network.vpc.id",description="AWS VPC the control plane is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.publicIp",description="Bastion IP address for breakglass access"

// AWSManagedControlPlane is the schema for the Amazon EKS Managed Control Plane API.
type AWSManagedControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSManagedControlPlaneSpec   `json:"spec,omitempty"`
	Status AWSManagedControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSManagedControlPlaneList contains a list of Amazon EKS Managed Control Planes.
type AWSManagedControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedControlPlane `json:"items"`
}

// GetConditions returns the control planes conditions.
func (r *AWSManagedControlPlane) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the status conditions for the AWSManagedControlPlane.
func (r *AWSManagedControlPlane) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&AWSManagedControlPlane{}, &AWSManagedControlPlaneList{})
}
