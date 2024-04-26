/*
Copyright 2023 The Kubernetes Authors.

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
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/errors"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

const (
	// MachineFinalizer allows ReconcileOpenStackMachine to clean up OpenStack resources associated with OpenStackMachine before
	// removing it from the apiserver.
	MachineFinalizer        = "openstackmachine.infrastructure.cluster.x-k8s.io"
	IPClaimMachineFinalizer = "openstackmachine.infrastructure.cluster.x-k8s.io/ip-claim"
)

// OpenStackMachineSpec defines the desired state of OpenStackMachine.
type OpenStackMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	ProviderID *string `json:"providerID,omitempty"`

	// The flavor reference for the flavor for your server instance.
	Flavor string `json:"flavor"`

	// The image to use for your server instance.
	// If the rootVolume is specified, this will be used when creating the root volume.
	// +required
	Image ImageParam `json:"image"`

	// The ssh key to inject in the instance
	SSHKeyName string `json:"sshKeyName,omitempty"`

	// Ports to be attached to the server instance. They are created if a port with the given name does not already exist.
	// If not specified a default port will be added for the default cluster network.
	Ports []PortOpts `json:"ports,omitempty"`

	// The names of the security groups to assign to the instance
	SecurityGroups []SecurityGroupParam `json:"securityGroups,omitempty"`

	// Whether the server instance is created on a trunk port or not.
	Trunk bool `json:"trunk,omitempty"`

	// Tags which will be added to the machine and all dependent resources
	// which support them. These are in addition to Tags defined on the
	// cluster.
	// Requires Nova api 2.52 minimum!
	// +listType=set
	Tags []string `json:"tags,omitempty"`

	// Metadata mapping. Allows you to create a map of key value pairs to add to the server instance.
	// +listType=map
	// +listMapKey=key
	ServerMetadata []ServerMetadata `json:"serverMetadata,omitempty"`

	// Config Drive support
	ConfigDrive *bool `json:"configDrive,omitempty"`

	// The volume metadata to boot from
	RootVolume *RootVolume `json:"rootVolume,omitempty"`

	// AdditionalBlockDevices is a list of specifications for additional block devices to attach to the server instance
	// +listType=map
	// +listMapKey=name
	// +optional
	AdditionalBlockDevices []AdditionalBlockDevice `json:"additionalBlockDevices,omitempty"`

	// The server group to assign the machine to.
	// +optional
	ServerGroup *ServerGroupParam `json:"serverGroup,omitempty"`

	// IdentityRef is a reference to a secret holding OpenStack credentials
	// to be used when reconciling this machine. If not specified, the
	// credentials specified in the cluster will be used.
	// +optional
	IdentityRef *OpenStackIdentityReference `json:"identityRef,omitempty"`

	// floatingIPPoolRef is a reference to a IPPool that will be assigned
	// to an IPAddressClaim. Once the IPAddressClaim is fulfilled, the FloatingIP
	// will be assigned to the OpenStackMachine.
	// +optional
	FloatingIPPoolRef *corev1.TypedLocalObjectReference `json:"floatingIPPoolRef,omitempty"`
}

type ServerMetadata struct {
	// Key is the server metadata key
	// +kubebuilder:validation:MaxLength:=255
	// +kubebuilder:validation:Required
	Key string `json:"key"`

	// Value is the server metadata value
	// +kubebuilder:validation:MaxLength:=255
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// OpenStackMachineStatus defines the observed state of OpenStackMachine.
type OpenStackMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// InstanceID is the OpenStack instance ID for this machine.
	// +optional
	InstanceID optional.String `json:"instanceID,omitempty"`

	// Addresses contains the OpenStack instance associated addresses.
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// InstanceState is the state of the OpenStack instance for this machine.
	// +optional
	InstanceState *InstanceState `json:"instanceState,omitempty"`

	// Resolved contains parts of the machine spec with all external
	// references fully resolved.
	// +optional
	Resolved *ResolvedMachineSpec `json:"resolved,omitempty"`

	// Resources contains references to OpenStack resources created for the machine.
	// +optional
	Resources *MachineResources `json:"resources,omitempty"`

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
// +genclient:Namespaced
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=openstackmachines,scope=Namespaced,categories=cluster-api,shortName=osm
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this OpenStackMachine belongs"
// +kubebuilder:printcolumn:name="InstanceState",type="string",JSONPath=".status.instanceState",description="OpenStack instance state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="ProviderID",type="string",JSONPath=".spec.providerID",description="OpenStack instance ID"
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns with this OpenStackMachine"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of OpenStackMachine"

// OpenStackMachine is the Schema for the openstackmachines API.
type OpenStackMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackMachineSpec   `json:"spec,omitempty"`
	Status OpenStackMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackMachineList contains a list of OpenStackMachine.
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
	objectTypes = append(objectTypes, &OpenStackMachine{}, &OpenStackMachineList{})
}
