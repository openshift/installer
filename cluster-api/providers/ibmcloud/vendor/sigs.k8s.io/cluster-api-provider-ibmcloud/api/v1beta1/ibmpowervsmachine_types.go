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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// IBMPowerVSMachineFinalizer allows IBMPowerVSMachineReconciler to clean up resources associated with IBMPowerVSMachine before
	// removing it from the apiserver.
	IBMPowerVSMachineFinalizer = "ibmpowervsmachine.infrastructure.cluster.x-k8s.io"
)

// IBMPowerVSMachineSpec defines the desired state of IBMPowerVSMachine.
type IBMPowerVSMachineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ServiceInstanceID is the id of the power cloud instance where the vsi instance will get deployed.
	// +kubebuilder:validation:MinLength=1
	ServiceInstanceID string `json:"serviceInstanceID"`

	// SSHKey is the name of the SSH key pair provided to the vsi for authenticating users.
	SSHKey string `json:"sshKey,omitempty"`

	// Image is the reference to the Image from which to create the machine instance.
	// +optional
	Image *IBMPowerVSResourceReference `json:"image,omitempty"`

	// ImageRef is an optional reference to a provider-specific resource that holds
	// the details for provisioning the Image for a Cluster.
	// +optional
	ImageRef *corev1.LocalObjectReference `json:"imageRef,omitempty"`

	// SysType is the System type used to host the vsi.
	// +optional
	SysType string `json:"sysType,omitempty"`

	// ProcType is the processor type, e.g: dedicated, shared, capped
	// +optional
	ProcType string `json:"procType,omitempty"`

	// Processors is Number of processors allocated.
	// +optional
	// +kubebuilder:validation:Pattern=^\d+(\.)?(\d)?(\d)?$
	Processors string `json:"processors,omitempty"`

	// Memory is Amount of memory allocated (in GB)
	// +optional
	Memory string `json:"memory,omitempty"`

	// Network is the reference to the Network to use for this instance.
	Network IBMPowerVSResourceReference `json:"network"`

	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`
}

// IBMPowerVSResourceReference is a reference to a specific PowerVS resource by ID or Name
// Only one of ID or Name may be specified. Specifying more than one will result in
// a validation error.
type IBMPowerVSResourceReference struct {
	// ID of resource
	// +kubebuilder:validation:MinLength=1
	// +optional
	ID *string `json:"id,omitempty"`

	// Name of resource
	// +kubebuilder:validation:MinLength=1
	// +optional
	Name *string `json:"name,omitempty"`

	// Regular expression to match resource,
	// In case of multiple resources matches the provided regular expression the first matched resource will be selected
	// +kubebuilder:validation:MinLength=1
	// +optional
	RegEx *string `json:"regex,omitempty"`
}

// IBMPowerVSMachineStatus defines the observed state of IBMPowerVSMachine.
type IBMPowerVSMachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	InstanceID string `json:"instanceID,omitempty"`

	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the vsi associated addresses.
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// Health is the health of the vsi.
	// +optional
	Health string `json:"health,omitempty"`

	// InstanceState is the status of the vsi.
	// +optional
	InstanceState PowerVSInstanceState `json:"instanceState,omitempty"`

	// Fault will report if any fault messages for the vsi.
	// +optional
	Fault string `json:"fault,omitempty"`

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

	// Conditions defines current service state of the IBMPowerVSMachine.
	// +optional
	Conditions capiv1beta1.Conditions `json:"conditions,omitempty"`

	// Region specifies the Power VS Service instance region.
	Region *string `json:"region,omitempty"`

	// Zone specifies the Power VS Service instance zone.
	Zone *string `json:"zone,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this IBMPowerVSMachine belongs"
// +kubebuilder:printcolumn:name="Machine",type="string",priority=1,JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object to which this IBMPowerVSMachine belongs"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of IBMPowerVSMachine"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for IBM PowerVS instances"
// +kubebuilder:printcolumn:name="Internal-IP",type="string",priority=1,JSONPath=".status.addresses[?(@.type==\"InternalIP\")].address",description="Instance Internal Addresses"
// +kubebuilder:printcolumn:name="External-IP",type="string",priority=1,JSONPath=".status.addresses[?(@.type==\"ExternalIP\")].address",description="Instance External Addresses"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.instanceState",description="PowerVS instance state"
// +kubebuilder:printcolumn:name="Health",type="string",JSONPath=".status.health",description="PowerVS instance health"

// IBMPowerVSMachine is the Schema for the ibmpowervsmachines API.
type IBMPowerVSMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IBMPowerVSMachineSpec   `json:"spec,omitempty"`
	Status IBMPowerVSMachineStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the IBMPowerVSMachine resource.
func (r *IBMPowerVSMachine) GetConditions() capiv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the IBMPowerVSMachine to the predescribed clusterv1.Conditions.
func (r *IBMPowerVSMachine) SetConditions(conditions capiv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

// IBMPowerVSMachineList contains a list of IBMPowerVSMachine.
type IBMPowerVSMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMPowerVSMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IBMPowerVSMachine{}, &IBMPowerVSMachineList{})
}
