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

package v1alpha6

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	// MachineFinalizer allows ReconcileOpenStackMachine to clean up OpenStack resources associated with OpenStackMachine before
	// removing it from the apiserver.
	MachineFinalizer = "openstackmachine.infrastructure.cluster.x-k8s.io"
)

// OpenStackMachineSpec defines the desired state of OpenStackMachine.
type OpenStackMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	ProviderID *string `json:"providerID,omitempty"`

	// InstanceID is the OpenStack instance ID for this machine.
	InstanceID *string `json:"instanceID,omitempty"`

	// The name of the cloud to use from the clouds secret
	// +optional
	CloudName string `json:"cloudName"`

	// The flavor reference for the flavor for your server instance.
	// +kubebuilder:validation:MinLength=1
	Flavor *string `json:"flavor,omitempty"`

	// FlavorID allows flavors to be specified by ID.  This field takes precedence
	// over Flavor.
	// +kubebuilder:validation:MinLength=1
	FlavorID *string `json:"flavorID,omitempty"`

	// The name of the image to use for your server instance.
	// If the RootVolume is specified, this will be ignored and use rootVolume directly.
	Image string `json:"image,omitempty"`

	// The uuid of the image to use for your server instance.
	// if it's empty, Image name will be used
	ImageUUID string `json:"imageUUID,omitempty"`

	// The ssh key to inject in the instance
	SSHKeyName string `json:"sshKeyName,omitempty"`

	// A networks object. Required parameter when there are multiple networks defined for the tenant.
	// When you do not specify both networks and ports parameters, the server attaches to the only network created for the current tenant.
	Networks []NetworkParam `json:"networks,omitempty"`

	// Ports to be attached to the server instance. They are created if a port with the given name does not already exist.
	// When you do not specify both networks and ports parameters, the server attaches to the only network created for the current tenant.
	Ports []PortOpts `json:"ports,omitempty"`

	// UUID, IP address of a port from this subnet will be marked as AccessIPv4 on the created compute instance
	Subnet string `json:"subnet,omitempty"`

	// The floatingIP which will be associated to the machine, only used for master.
	// The floatingIP should have been created and haven't been associated.
	FloatingIP string `json:"floatingIP,omitempty"`

	// The names of the security groups to assign to the instance
	SecurityGroups []SecurityGroupParam `json:"securityGroups,omitempty"`

	// Whether the server instance is created on a trunk port or not.
	Trunk bool `json:"trunk,omitempty"`

	// Machine tags
	// Requires Nova api 2.52 minimum!
	// +listType=set
	Tags []string `json:"tags,omitempty"`

	// Metadata mapping. Allows you to create a map of key value pairs to add to the server instance.
	ServerMetadata map[string]string `json:"serverMetadata,omitempty"`

	// Config Drive support
	ConfigDrive *bool `json:"configDrive,omitempty"`

	// The volume metadata to boot from
	RootVolume *RootVolume `json:"rootVolume,omitempty"`

	// The server group to assign the machine to
	ServerGroupID string `json:"serverGroupID,omitempty"`

	// IdentityRef is a reference to a identity to be used when reconciling this cluster
	// +optional
	IdentityRef *OpenStackIdentityReference `json:"identityRef,omitempty"`
}

// OpenStackMachineStatus defines the observed state of OpenStackMachine.
type OpenStackMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the OpenStack instance associated addresses.
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// InstanceState is the state of the OpenStack instance for this machine.
	// +optional
	InstanceState *InstanceState `json:"instanceState,omitempty"`

	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

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

	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion:warning="The v1alpha6 version of OpenStackMachine has been deprecated and will be removed in a future release of the API. Please upgrade."
// +kubebuilder:resource:path=openstackmachines,scope=Namespaced,categories=cluster-api,shortName=osm
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this OpenStackMachine belongs"
// +kubebuilder:printcolumn:name="InstanceState",type="string",JSONPath=".status.instanceState",description="OpenStack instance state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="ProviderID",type="string",JSONPath=".spec.providerID",description="OpenStack instance ID"
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns with this OpenStackMachine"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of OpenStackMachine"

// OpenStackMachine is the Schema for the openstackmachines API.
//
// Deprecated: This type will be removed in one of the next releases.
type OpenStackMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackMachineSpec   `json:"spec,omitempty"`
	Status OpenStackMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion

// OpenStackMachineList contains a list of OpenStackMachine.
//
// Deprecated: This type will be removed in one of the next releases.
type OpenStackMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackMachine `json:"items"`
}

// GetConditions returns the observations of the operational state of the OpenStackMachine resource.
func (r *OpenStackMachine) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the OpenStackMachine to the predescribed clusterv1.Conditions.
func (r *OpenStackMachine) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// SetFailure sets the OpenStackMachine status failure reason and failure message.
func (r *OpenStackMachine) SetFailure(failureReason errors.MachineStatusError, failureMessage error) {
	r.Status.FailureReason = &failureReason
	r.Status.FailureMessage = ptr.To(failureMessage.Error())
}

func init() {
	SchemeBuilder.Register(&OpenStackMachine{}, &OpenStackMachineList{})
}
