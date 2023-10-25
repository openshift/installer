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

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ClusterFinalizer allows ReconcileAWSCluster to clean up AWS resources associated with AWSCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "awscluster.infrastructure.cluster.x-k8s.io"

	// AWSClusterControllerIdentityName is the name of the AWSClusterControllerIdentity singleton.
	AWSClusterControllerIdentityName = "default"
)

// AWSClusterSpec defines the desired state of an EC2-based Kubernetes cluster.
type AWSClusterSpec struct {
	// NetworkSpec encapsulates all things related to AWS network.
	NetworkSpec NetworkSpec `json:"network,omitempty"`

	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// ControlPlaneLoadBalancer is optional configuration for customizing control plane behavior.
	// +optional
	ControlPlaneLoadBalancer *AWSLoadBalancerSpec `json:"controlPlaneLoadBalancer,omitempty"`

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
	Bastion Bastion `json:"bastion"`

	// IdentityRef is a reference to a identity to be used when reconciling this cluster
	// +optional
	IdentityRef *AWSIdentityReference `json:"identityRef,omitempty"`

	// S3Bucket contains options to configure a supporting S3 bucket for this
	// cluster - currently used for nodes requiring Ignition
	// (https://coreos.github.io/ignition/) for bootstrapping (requires
	// BootstrapFormatIgnition feature flag to be enabled).
	// +optional
	S3Bucket *S3Bucket `json:"s3Bucket,omitempty"`
}

// AWSIdentityKind defines allowed AWS identity types.
type AWSIdentityKind string

var (
	// ControllerIdentityKind defines identity reference kind as AWSClusterControllerIdentity.
	ControllerIdentityKind = AWSIdentityKind("AWSClusterControllerIdentity")

	// ClusterRoleIdentityKind defines identity reference kind as AWSClusterRoleIdentity.
	ClusterRoleIdentityKind = AWSIdentityKind("AWSClusterRoleIdentity")

	// ClusterStaticIdentityKind defines identity reference kind as AWSClusterStaticIdentity.
	ClusterStaticIdentityKind = AWSIdentityKind("AWSClusterStaticIdentity")
)

// AWSIdentityReference specifies a identity.
type AWSIdentityReference struct {
	// Name of the identity.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Kind of the identity.
	// +kubebuilder:validation:Enum=AWSClusterControllerIdentity;AWSClusterRoleIdentity;AWSClusterStaticIdentity
	Kind AWSIdentityKind `json:"kind"`
}

// Bastion defines a bastion host.
type Bastion struct {
	// Enabled allows this provider to create a bastion host instance
	// with a public ip to access the VPC private network.
	// +optional
	Enabled bool `json:"enabled"`

	// DisableIngressRules will ensure there are no Ingress rules in the bastion host's security group.
	// Requires AllowedCIDRBlocks to be empty.
	// +optional
	DisableIngressRules bool `json:"disableIngressRules,omitempty"`

	// AllowedCIDRBlocks is a list of CIDR blocks allowed to access the bastion host.
	// They are set as ingress rules for the Bastion host's Security Group (defaults to 0.0.0.0/0).
	// +optional
	AllowedCIDRBlocks []string `json:"allowedCIDRBlocks,omitempty"`

	// InstanceType will use the specified instance type for the bastion. If not specified,
	// Cluster API Provider AWS will use t3.micro for all regions except us-east-1, where t2.micro
	// will be the default.
	InstanceType string `json:"instanceType,omitempty"`

	// AMI will use the specified AMI to boot the bastion. If not specified,
	// the AMI will default to one picked out in public space.
	// +optional
	AMI string `json:"ami,omitempty"`
}

// AWSLoadBalancerSpec defines the desired state of an AWS load balancer.
type AWSLoadBalancerSpec struct {
	// Name sets the name of the classic ELB load balancer. As per AWS, the name must be unique
	// within your set of load balancers for the region, must have a maximum of 32 characters, must
	// contain only alphanumeric characters or hyphens, and cannot begin or end with a hyphen. Once
	// set, the value cannot be changed.
	// +kubebuilder:validation:MaxLength:=32
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9]([A-Za-z0-9]{0,31}|[-A-Za-z0-9]{0,30}[A-Za-z0-9])$`
	// +optional
	Name *string `json:"name,omitempty"`

	// Scheme sets the scheme of the load balancer (defaults to internet-facing)
	// +kubebuilder:default=internet-facing
	// +kubebuilder:validation:Enum=internet-facing;internal
	// +optional
	Scheme *ClassicELBScheme `json:"scheme,omitempty"`

	// CrossZoneLoadBalancing enables the classic ELB cross availability zone balancing.
	//
	// With cross-zone load balancing, each load balancer node for your Classic Load Balancer
	// distributes requests evenly across the registered instances in all enabled Availability Zones.
	// If cross-zone load balancing is disabled, each load balancer node distributes requests evenly across
	// the registered instances in its Availability Zone only.
	//
	// Defaults to false.
	// +optional
	CrossZoneLoadBalancing bool `json:"crossZoneLoadBalancing"`

	// Subnets sets the subnets that should be applied to the control plane load balancer (defaults to discovered subnets for managed VPCs or an empty set for unmanaged VPCs)
	// +optional
	Subnets []string `json:"subnets,omitempty"`

	// HealthCheckProtocol sets the protocol type for classic ELB health check target
	// default value is ClassicELBProtocolSSL
	// +optional
	HealthCheckProtocol *ClassicELBProtocol `json:"healthCheckProtocol,omitempty"`

	// AdditionalSecurityGroups sets the security groups used by the load balancer. Expected to be security group IDs
	// This is optional - if not provided new security groups will be created for the load balancer
	// +optional
	AdditionalSecurityGroups []string `json:"additionalSecurityGroups,omitempty"`
}

// AWSClusterStatus defines the observed state of AWSCluster.
type AWSClusterStatus struct {
	// +kubebuilder:default=false
	Ready          bool                     `json:"ready"`
	Network        NetworkStatus            `json:"networkStatus,omitempty"`
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
	Bastion        *Instance                `json:"bastion,omitempty"`
	Conditions     clusterv1.Conditions     `json:"conditions,omitempty"`
}

type S3Bucket struct {
	// ControlPlaneIAMInstanceProfile is a name of the IAMInstanceProfile, which will be allowed
	// to read control-plane node bootstrap data from S3 Bucket.
	ControlPlaneIAMInstanceProfile string `json:"controlPlaneIAMInstanceProfile"`

	// NodesIAMInstanceProfiles is a list of IAM instance profiles, which will be allowed to read
	// worker nodes bootstrap data from S3 Bucket.
	NodesIAMInstanceProfiles []string `json:"nodesIAMInstanceProfiles"`

	// Name defines name of S3 Bucket to be created.
	// +kubebuilder:validation:MinLength:=3
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`
	Name string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsclusters,scope=Namespaced,categories=cluster-api,shortName=awsc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AWSCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for EC2 instances"
// +kubebuilder:printcolumn:name="VPC",type="string",JSONPath=".spec.network.vpc.id",description="AWS VPC the cluster is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.publicIp",description="Bastion IP address for breakglass access"

// AWSCluster is the schema for Amazon EC2 based Kubernetes Cluster API.
type AWSCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSClusterSpec   `json:"spec,omitempty"`
	Status AWSClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSClusterList contains a list of AWSCluster.
type AWSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSCluster `json:"items"`
}

// GetConditions returns the observations of the operational state of the AWSCluster resource.
func (r *AWSCluster) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the AWSCluster to the predescribed clusterv1.Conditions.
func (r *AWSCluster) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&AWSCluster{}, &AWSClusterList{})
}
