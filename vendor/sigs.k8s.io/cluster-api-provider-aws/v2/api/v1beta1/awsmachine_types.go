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
	// MachineFinalizer allows ReconcileAWSMachine to clean up AWS resources associated with AWSMachine before
	// removing it from the apiserver.
	MachineFinalizer = "awsmachine.infrastructure.cluster.x-k8s.io"

	// DefaultIgnitionVersion represents default Ignition version generated for machine userdata.
	DefaultIgnitionVersion = "2.3"
)

// SecretBackend defines variants for backend secret storage.
type SecretBackend string

var (
	// SecretBackendSSMParameterStore defines AWS Systems Manager Parameter Store as the secret backend.
	SecretBackendSSMParameterStore = SecretBackend("ssm-parameter-store")

	// SecretBackendSecretsManager defines AWS Secrets Manager as the secret backend.
	SecretBackendSecretsManager = SecretBackend("secrets-manager")
)

// AWSMachineSpec defines the desired state of an Amazon EC2 instance.
type AWSMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	ProviderID *string `json:"providerID,omitempty"`

	// InstanceID is the EC2 instance ID for this machine.
	InstanceID *string `json:"instanceID,omitempty"`

	// AMI is the reference to the AMI from which to create the machine instance.
	AMI AMIReference `json:"ami,omitempty"`

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
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=2
	InstanceType string `json:"instanceType"`

	// AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
	// AWS provider. If both the AWSCluster and the AWSMachine specify the same tag name with different values, the
	// AWSMachine's value takes precedence.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// IAMInstanceProfile is a name of an IAM instance profile to assign to the instance
	// +optional
	IAMInstanceProfile string `json:"iamInstanceProfile,omitempty"`

	// PublicIP specifies whether the instance should get a public IP.
	// Precedence for this setting is as follows:
	// 1. This field if set
	// 2. Cluster/flavor setting
	// 3. Subnet default
	// +optional
	PublicIP *bool `json:"publicIP,omitempty"`

	// AdditionalSecurityGroups is an array of references to security groups that should be applied to the
	// instance. These security groups would be set in addition to any security groups defined
	// at the cluster level or in the actuator. It is possible to specify either IDs of Filters. Using Filters
	// will cause additional requests to AWS API and if tags change the attached security groups might change too.
	// +optional
	AdditionalSecurityGroups []AWSResourceReference `json:"additionalSecurityGroups,omitempty"`

	// FailureDomain is the failure domain unique identifier this Machine should be attached to, as defined in Cluster API.
	// For this infrastructure provider, the ID is equivalent to an AWS Availability Zone.
	// If multiple subnets are matched for the availability zone, the first one returned is picked.
	FailureDomain *string `json:"failureDomain,omitempty"`

	// Subnet is a reference to the subnet to use for this instance. If not specified,
	// the cluster subnet will be used.
	// +optional
	Subnet *AWSResourceReference `json:"subnet,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the instance. Valid values are empty string (do not use SSH keys), a valid SSH key name, or omitted (use the default SSH key name)
	// +optional
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// RootVolume encapsulates the configuration options for the root volume
	// +optional
	RootVolume *Volume `json:"rootVolume,omitempty"`

	// Configuration options for the non root storage volumes.
	// +optional
	NonRootVolumes []Volume `json:"nonRootVolumes,omitempty"`

	// NetworkInterfaces is a list of ENIs to associate with the instance.
	// A maximum of 2 may be specified.
	// +optional
	// +kubebuilder:validation:MaxItems=2
	NetworkInterfaces []string `json:"networkInterfaces,omitempty"`

	// UncompressedUserData specify whether the user data is gzip-compressed before it is sent to ec2 instance.
	// cloud-init has built-in support for gzip-compressed user data
	// user data stored in aws secret manager is always gzip-compressed.
	//
	// +optional
	UncompressedUserData *bool `json:"uncompressedUserData,omitempty"`

	// CloudInit defines options related to the bootstrapping systems where
	// CloudInit is used.
	// +optional
	CloudInit CloudInit `json:"cloudInit,omitempty"`

	// Ignition defined options related to the bootstrapping systems where Ignition is used.
	// +optional
	Ignition *Ignition `json:"ignition,omitempty"`

	// SpotMarketOptions allows users to configure instances to be run using AWS Spot instances.
	// +optional
	SpotMarketOptions *SpotMarketOptions `json:"spotMarketOptions,omitempty"`

	// Tenancy indicates if instance should run on shared or single-tenant hardware.
	// +optional
	// +kubebuilder:validation:Enum:=default;dedicated;host
	Tenancy string `json:"tenancy,omitempty"`
}

// CloudInit defines options related to the bootstrapping systems where
// CloudInit is used.
type CloudInit struct {
	// InsecureSkipSecretsManager, when set to true will not use AWS Secrets Manager
	// or AWS Systems Manager Parameter Store to ensure privacy of userdata.
	// By default, a cloud-init boothook shell script is prepended to download
	// the userdata from Secrets Manager and additionally delete the secret.
	InsecureSkipSecretsManager bool `json:"insecureSkipSecretsManager,omitempty"`

	// SecretCount is the number of secrets used to form the complete secret
	// +optional
	SecretCount int32 `json:"secretCount,omitempty"`

	// SecretPrefix is the prefix for the secret name. This is stored
	// temporarily, and deleted when the machine registers as a node against
	// the workload cluster.
	// +optional
	SecretPrefix string `json:"secretPrefix,omitempty"`

	// SecureSecretsBackend, when set to parameter-store will utilize the AWS Systems Manager
	// Parameter Storage to distribute secrets. By default or with the value of secrets-manager,
	// will use AWS Secrets Manager instead.
	// +optional
	// +kubebuilder:validation:Enum=secrets-manager;ssm-parameter-store
	SecureSecretsBackend SecretBackend `json:"secureSecretsBackend,omitempty"`
}

// Ignition defines options related to the bootstrapping systems where Ignition is used.
type Ignition struct {
	// Version defines which version of Ignition will be used to generate bootstrap data.
	//
	// +optional
	// +kubebuilder:default="2.3"
	// +kubebuilder:validation:Enum="2.3"
	Version string `json:"version,omitempty"`
}

// AWSMachineStatus defines the observed state of AWSMachine.
type AWSMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Interruptible reports that this machine is using spot instances and can therefore be interrupted by CAPI when it receives a notice that the spot instance is to be terminated by AWS.
	// This will be set to true when SpotMarketOptions is not nil (i.e. this machine is using a spot instance).
	// +optional
	Interruptible bool `json:"interruptible,omitempty"`

	// Addresses contains the AWS instance associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// InstanceState is the state of the AWS instance for this machine.
	// +optional
	InstanceState *InstanceState `json:"instanceState,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
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
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *string `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
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
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// Conditions defines current service state of the AWSMachine.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:resource:path=awsmachines,scope=Namespaced,categories=cluster-api,shortName=awsm
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AWSMachine belongs"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.instanceState",description="EC2 instance state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="InstanceID",type="string",JSONPath=".spec.providerID",description="EC2 instance ID"
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns with this AWSMachine"

// AWSMachine is the schema for Amazon EC2 machines.
type AWSMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSMachineSpec   `json:"spec,omitempty"`
	Status AWSMachineStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the AWSMachine resource.
func (r *AWSMachine) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the AWSMachine to the predescribed clusterv1.Conditions.
func (r *AWSMachine) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion

// AWSMachineList contains a list of Amazon EC2 machines.
type AWSMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSMachine{}, &AWSMachineList{})
}
