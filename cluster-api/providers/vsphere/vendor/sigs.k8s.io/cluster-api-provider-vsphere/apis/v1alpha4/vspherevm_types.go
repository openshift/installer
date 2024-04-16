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

//nolint:godot
package v1alpha4

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	// VMFinalizer allows the reconciler to clean up resources associated
	// with a VSphereVM before removing it from the API Server.
	VMFinalizer = "vspherevm.infrastructure.cluster.x-k8s.io"
)

// VSphereVMSpec defines the desired state of VSphereVM.
type VSphereVMSpec struct {
	VirtualMachineCloneSpec `json:",inline"`

	// BootstrapRef is a reference to a bootstrap provider-specific resource
	// that holds configuration details.
	// This field is optional in case no bootstrap data is required to create
	// a VM.
	// +optional
	BootstrapRef *corev1.ObjectReference `json:"bootstrapRef,omitempty"`

	// BiosUUID is the VM's BIOS UUID that is assigned at runtime after
	// the VM has been created.
	// This field is required at runtime for other controllers that read
	// this CRD as unstructured data.
	// +optional
	BiosUUID string `json:"biosUUID,omitempty"`
}

// VSphereVMStatus defines the observed state of VSphereVM
type VSphereVMStatus struct {
	// Ready is true when the provider resource is ready.
	// This field is required at runtime for other controllers that read
	// this CRD as unstructured data.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Addresses is a list of the VM's IP addresses.
	// This field is required at runtime for other controllers that read
	// this CRD as unstructured data.
	// +optional
	Addresses []string `json:"addresses,omitempty"`

	// CloneMode is the type of clone operation used to clone this VM. Since
	// LinkedMode is the default but fails gracefully if the source of the
	// clone has no snapshots, this field may be used to determine the actual
	// type of clone operation used to create this VM.
	// +optional
	CloneMode CloneMode `json:"cloneMode,omitempty"`

	// Snapshot is the name of the snapshot from which the VM was cloned if
	// LinkedMode is enabled.
	// +optional
	Snapshot string `json:"snapshot,omitempty"`

	// RetryAfter tracks the time we can retry queueing a task
	// +optional
	RetryAfter metav1.Time `json:"retryAfter,omitempty"`

	// TaskRef is a managed object reference to a Task related to the machine.
	// This value is set automatically at runtime and should not be set or
	// modified by users.
	// +optional
	TaskRef string `json:"taskRef,omitempty"`

	// Network returns the network status for each of the machine's configured
	// network interfaces.
	// +optional
	Network []NetworkStatus `json:"network,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the vspherevm and will contain a succinct value suitable
	// for vm interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the vm.
	//
	// Any transient errors that occur during the reconciliation of vspherevms
	// can be added as events to the vspherevm object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the vspherevm and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the vm.
	//
	// Any transient errors that occur during the reconciliation of vspherevms
	// can be added as events to the vspherevm object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// Conditions defines current service state of the VSphereVM.
	// +optional
	Conditions Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion
// +kubebuilder:resource:path=vspherevms,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status

// VSphereVM is the Schema for the vspherevms API
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereVM struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereVMSpec   `json:"spec,omitempty"`
	Status VSphereVMStatus `json:"status,omitempty"`
}

func (r *VSphereVM) GetConditions() Conditions {
	return r.Status.Conditions
}

func (r *VSphereVM) SetConditions(conditions Conditions) {
	r.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VSphereVMList contains a list of VSphereVM
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereVMList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereVM `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereVM{}, &VSphereVMList{})
}
