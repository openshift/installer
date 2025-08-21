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

// VSphereVM's Ready condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMReadyV1Beta2Condition is true if the VSphereVM's deletionTimestamp is not set, VSphereVM's
	// VirtualMachineProvisioned, VCenterAvailable and IPAddressClaimsFulfilled are true.
	VSphereVMReadyV1Beta2Condition = clusterv1.ReadyV1Beta2Condition

	// VSphereVMReadyV1Beta2Reason surfaces when the VSphereVM readiness criteria is met.
	VSphereVMReadyV1Beta2Reason = clusterv1.ReadyV1Beta2Reason

	// VSphereVMNotReadyV1Beta2Reason surfaces when the VSphereVM readiness criteria is not met.
	VSphereVMNotReadyV1Beta2Reason = clusterv1.NotReadyV1Beta2Reason

	// VSphereVMReadyUnknownV1Beta2Reason surfaces when at least one VSphereVM readiness criteria is unknown
	// and no VSphereVM readiness criteria is not met.
	VSphereVMReadyUnknownV1Beta2Reason = clusterv1.ReadyUnknownV1Beta2Reason
)

// VSphereVM's VirtualMachineProvisioned condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMVirtualMachineProvisionedV1Beta2Condition documents the status of the VirtualMachine that is controlled
	// by the VSphereVM.
	VSphereVMVirtualMachineProvisionedV1Beta2Condition = "VirtualMachineProvisioned"

	// VSphereVMVirtualMachineWaitingForCloneV1Beta2Reason documents the VirtualMachine that is controlled
	// by the VSphereVM waiting for the clone operation to complete.
	VSphereVMVirtualMachineWaitingForCloneV1Beta2Reason = "WaitingForClone"

	// VSphereVMVirtualMachineWaitingForStaticIPAllocationV1Beta2Reason documents the VirtualMachine that is controlled
	// by the VSphereVM waiting for the allocation of a static IP address.
	VSphereVMVirtualMachineWaitingForStaticIPAllocationV1Beta2Reason = "WaitingForStaticIPAllocation"

	// VSphereVMVirtualMachineWaitingForIPAddressV1Beta2Reason documents the VirtualMachine that is controlled
	// by the VSphereVM waiting for an IP address to be provisioned from the IPAM provider.
	VSphereVMVirtualMachineWaitingForIPAddressV1Beta2Reason = "WaitingForIPAddress"

	// VSphereVMVirtualMachineWaitingForIPAllocationV1Beta2Reason documents the VirtualMachine that is controlled
	// by the VSphereVM waiting for the allocation of an IP address.
	// This is used when the dhcp4 or dhcp6 for a VirtualMachine is set and the VirtualMachine is waiting for the
	// relevant IP address to show up on the VM.
	VSphereVMVirtualMachineWaitingForIPAllocationV1Beta2Reason = "WaitingForIPAllocation"

	// VSphereVMVirtualMachinePoweringOnV1Beta2Reason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM is executing the power on sequence.
	VSphereVMVirtualMachinePoweringOnV1Beta2Reason = "PoweringOn"

	// VSphereVMVirtualMachineProvisionedV1Beta2Reason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM is provisioned.
	VSphereVMVirtualMachineProvisionedV1Beta2Reason = clusterv1.ProvisionedV1Beta2Reason

	// VSphereVMVirtualMachineTaskFailedV1Beta2Reason surfaces when a task for the VirtualMachine that is controlled
	// by the VSphereVM failed; the reconcile look will automatically retry the operation,
	// but a user intervention might be required to fix the problem.
	VSphereVMVirtualMachineTaskFailedV1Beta2Reason = "TaskFailed"

	// VSphereVMVirtualMachineNotFoundByBIOSUUIDV1Beta2Reason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM can't be found by BIOS UUID.
	// Those kind of errors could be transient sometimes and failed VSphereVM are automatically
	// reconciled by the controller.
	VSphereVMVirtualMachineNotFoundByBIOSUUIDV1Beta2Reason = "NotFoundByBIOSUUID"

	// VSphereVMVirtualMachineNotProvisionedV1Beta2Reason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM is not provisioned.
	VSphereVMVirtualMachineNotProvisionedV1Beta2Reason = clusterv1.NotProvisionedV1Beta2Reason

	// VSphereVMVirtualMachineDeletingV1Beta2Reason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM is being deleted.
	VSphereVMVirtualMachineDeletingV1Beta2Reason = clusterv1.DeletingV1Beta2Reason
)

// VSphereVM's VCenterAvailable condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMVCenterAvailableV1Beta2Condition documents the availability of the VCenter hosting the VSphereVM.
	VSphereVMVCenterAvailableV1Beta2Condition = "VCenterAvailable"

	// VSphereVMVCenterAvailableV1Beta2Reason documents the VCenter hosting the VSphereVM
	// being available.
	VSphereVMVCenterAvailableV1Beta2Reason = clusterv1.AvailableV1Beta2Reason

	// VSphereVMVCenterUnreachableV1Beta2Reason documents the VCenter hosting the VSphereVM
	// cannot be reached.
	VSphereVMVCenterUnreachableV1Beta2Reason = "VCenterUnreachable"
)

// VSphereVM's IPAddressClaimsFulfilled condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMIPAddressClaimsFulfilledV1Beta2Condition documents the status of claiming an IP address
	// from an IPAM provider.
	VSphereVMIPAddressClaimsFulfilledV1Beta2Condition = "IPAddressClaimsFulfilled"

	// VSphereVMIPAddressClaimsBeingCreatedV1Beta2Reason documents that claims for the
	// IP addresses required by the VSphereVM are being created.
	VSphereVMIPAddressClaimsBeingCreatedV1Beta2Reason = "IPAddressClaimsBeingCreated"

	// VSphereVMIPAddressClaimsWaitingForIPAddressV1Beta2Reason documents that claims for the
	// IP addresses required by the VSphereVM are waiting for IP addresses.
	VSphereVMIPAddressClaimsWaitingForIPAddressV1Beta2Reason = "WaitingForIPAddress"

	// VSphereVMIPAddressClaimsFulfilledV1Beta2Reason documents that claims for the
	// IP addresses required by the VSphereVM are fulfilled.
	VSphereVMIPAddressClaimsFulfilledV1Beta2Reason = "Fulfilled"

	// VSphereVMIPAddressClaimsNotFulfilledV1Beta2Reason documents that claims for the
	// IP addresses required by the VSphereVM are not fulfilled.
	VSphereVMIPAddressClaimsNotFulfilledV1Beta2Reason = "NotFulfilled"
)

// VSphereVM's GuestSoftPowerOffSucceeded condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMGuestSoftPowerOffSucceededV1Beta2Condition documents the status of performing guest initiated
	// graceful shutdown.
	VSphereVMGuestSoftPowerOffSucceededV1Beta2Condition string = "GuestSoftPowerOffSucceeded"

	// VSphereVMGuestSoftPowerOffInProgressV1Beta2Reason documents that the guest receives
	// a graceful shutdown request.
	VSphereVMGuestSoftPowerOffInProgressV1Beta2Reason = "InProgress"

	// GuestSoftPowerOffFailedV1Beta2Reason documents that the graceful
	// shutdown request failed.
	VSphereVMGuestSoftPowerOffFailedV1Beta2Reason = "Failed"

	// GuestSoftPowerOffSucceededV1Beta2Reason documents that the graceful
	// shutdown request succeeded.
	VSphereVMGuestSoftPowerOffSucceededV1Beta2Reason = "Succeeded"
)

// VSphereVM's PCIDevicesDetached condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMPCIDevicesDetachedV1Beta2Condition documents the status of the attached PCI devices on the VSphereVM.
	// It is a negative condition to notify the user that the device(s) is no longer attached to
	// the underlying VM and would require manual intervention to fix the situation.
	VSphereVMPCIDevicesDetachedV1Beta2Condition string = "PCIDevicesDetached"

	// VSphereVMPCIDevicesDetachedNotFoundV1Beta2Reason documents the VSphereVM not having the PCI device attached during VM startup.
	// This would indicate that the PCI devices were removed out of band by an external entity.
	VSphereVMPCIDevicesDetachedNotFoundV1Beta2Reason = "NotFound"
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

	// v1beta2 groups all the fields that will be added or modified in VSphereVM's status with the V1Beta2 version.
	// +optional
	V1Beta2 *VSphereVMV1Beta2Status `json:"v1beta2,omitempty"`
}

// VSphereVMV1Beta2Status groups all the fields that will be added or modified in VSphereVMStatus with the V1Beta2 version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereVMV1Beta2Status struct {
	// conditions represents the observations of a VSphereVM's current state.
	// Known condition types are Ready, VirtualMachineProvisioned, VCenterAvailable and IPAddressClaimsFulfilled,
	// GuestSoftPowerOffSucceeded, PCIDevicesDetached and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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

// GetV1Beta2Conditions returns the set of conditions for this object.
func (r *VSphereVM) GetV1Beta2Conditions() []metav1.Condition {
	if r.Status.V1Beta2 == nil {
		return nil
	}
	return r.Status.V1Beta2.Conditions
}

// SetV1Beta2Conditions sets conditions for an API object.
func (r *VSphereVM) SetV1Beta2Conditions(conditions []metav1.Condition) {
	if r.Status.V1Beta2 == nil {
		r.Status.V1Beta2 = &VSphereVMV1Beta2Status{}
	}
	r.Status.V1Beta2.Conditions = conditions
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
