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
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	// VMFinalizer allows the reconciler to clean up resources associated
	// with a VSphereVM before removing it from the API Server.
	VMFinalizer = "vspherevm.infrastructure.cluster.x-k8s.io"

	// IPAddressClaimFinalizer allows the reconciler to prevent deletion of an
	// IPAddressClaim that is in use.
	IPAddressClaimFinalizer = "vspherevm.infrastructure.cluster.x-k8s.io/ip-claim-protection"

	// GuestSoftPowerOffDefaultTimeout is the default timeout to wait for
	// shutdown finishes in the guest VM before powering off the VM forcibly
	// Only effective when the powerOffMode is set to trySoft.
	GuestSoftPowerOffDefaultTimeout = 5 * time.Minute
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

	// PowerOffMode describes the desired behavior when powering off a VM.
	//
	// There are three, supported power off modes: hard, soft, and
	// trySoft. The first mode, hard, is the equivalent of a physical
	// system's power cord being ripped from the wall. The soft mode
	// requires the VM's guest to have VM Tools installed and attempts to
	// gracefully shut down the VM. Its variant, trySoft, first attempts
	// a graceful shutdown, and if that fails or the VM is not in a powered off
	// state after reaching the GuestSoftPowerOffTimeout, the VM is halted.
	//
	// If omitted, the mode defaults to hard.
	//
	// +optional
	// +kubebuilder:default=hard
	PowerOffMode VirtualMachinePowerOpMode `json:"powerOffMode,omitempty"`

	// GuestSoftPowerOffTimeout sets the wait timeout for shutdown in the VM guest.
	// The VM will be powered off forcibly after the timeout if the VM is still
	// up and running when the PowerOffMode is set to trySoft.
	//
	// This parameter only applies when the PowerOffMode is set to trySoft.
	//
	// If omitted, the timeout defaults to 5 minutes.
	//
	// +optional
	GuestSoftPowerOffTimeout *metav1.Duration `json:"guestSoftPowerOffTimeout,omitempty"`
}

// VSphereVMStatus defines the observed state of VSphereVM.
type VSphereVMStatus struct {
	// Host describes the hostname or IP address of the infrastructure host
	// that the VSphereVM is residing on.
	// +optional
	Host string `json:"host,omitempty"`

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
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// ModuleUUID is the unique identifier for the vCenter cluster module construct
	// which is used to configure anti-affinity. Objects with the same ModuleUUID
	// will be anti-affined, meaning that the vCenter DRS will best effort schedule
	// the VMs on separate hosts.
	// +optional
	ModuleUUID *string `json:"moduleUUID,omitempty"`

	// VMRef is the VM's Managed Object Reference on vSphere. It can be used by consumers
	// to programatically get this VM representation on vSphere in case of the need to retrieve informations.
	// This field is set once the machine is created and should not be changed
	// +optional
	VMRef string `json:"vmRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspherevms,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VSphereVM is the Schema for the vspherevms API.
type VSphereVM struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereVMSpec   `json:"spec,omitempty"`
	Status VSphereVMStatus `json:"status,omitempty"`
}

// GetConditions returns the conditions for a VSphereVM.
func (r *VSphereVM) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the conditions on a VSphereVM.
func (r *VSphereVM) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VSphereVMList contains a list of VSphereVM.
type VSphereVMList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereVM `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereVM{}, &VSphereVMList{})
}
