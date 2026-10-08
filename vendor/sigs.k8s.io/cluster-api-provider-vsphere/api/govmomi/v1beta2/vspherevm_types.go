/*
Copyright 2025 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/api/deprecated/errors"
)

const (
	// VMFinalizer allows the reconciler to clean up resources associated
	// with a VSphereVM before removing it from the API Server.
	VMFinalizer = "vspherevm.infrastructure.cluster.x-k8s.io"

	// IPAddressClaimFinalizer allows the reconciler to prevent deletion of an
	// IPAddressClaim that is in use.
	IPAddressClaimFinalizer = "vspherevm.infrastructure.cluster.x-k8s.io/ip-claim-protection"

	// GuestSoftPowerOffDefaultTimeoutSeconds is the default timeout to wait for
	// shutdown finishes in the guest VM before powering off the VM forcibly
	// Only effective when the powerOffMode is set to trySoft.
	GuestSoftPowerOffDefaultTimeoutSeconds = 5 * 60
)

// VSphereVM's Ready condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMReadyCondition is true if the VSphereVM's deletionTimestamp is not set, VSphereVM's
	// VirtualMachineProvisioned, VCenterAvailable and IPAddressClaimsFulfilled are true.
	VSphereVMReadyCondition = clusterv1.ReadyCondition

	// VSphereVMReadyReason surfaces when the VSphereVM readiness criteria is met.
	VSphereVMReadyReason = clusterv1.ReadyReason

	// VSphereVMNotReadyReason surfaces when the VSphereVM readiness criteria is not met.
	VSphereVMNotReadyReason = clusterv1.NotReadyReason

	// VSphereVMReadyUnknownReason surfaces when at least one VSphereVM readiness criteria is unknown
	// and no VSphereVM readiness criteria is not met.
	VSphereVMReadyUnknownReason = clusterv1.ReadyUnknownReason
)

// VSphereVM's VirtualMachineProvisioned condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMVirtualMachineProvisionedCondition documents the status of the VirtualMachine that is controlled
	// by the VSphereVM.
	VSphereVMVirtualMachineProvisionedCondition = "VirtualMachineProvisioned"

	// VSphereVMVirtualMachineWaitingForCloneReason documents the VirtualMachine that is controlled
	// by the VSphereVM waiting for the clone operation to complete.
	VSphereVMVirtualMachineWaitingForCloneReason = "WaitingForClone"

	// VSphereVMVirtualMachineWaitingForStaticIPAllocationReason documents the VirtualMachine that is controlled
	// by the VSphereVM waiting for the allocation of a static IP address.
	VSphereVMVirtualMachineWaitingForStaticIPAllocationReason = "WaitingForStaticIPAllocation"

	// VSphereVMVirtualMachineWaitingForIPAddressReason documents the VirtualMachine that is controlled
	// by the VSphereVM waiting for an IP address to be provisioned from the IPAM provider.
	VSphereVMVirtualMachineWaitingForIPAddressReason = "WaitingForIPAddress"

	// VSphereVMVirtualMachineWaitingForIPAllocationReason documents the VirtualMachine that is controlled
	// by the VSphereVM waiting for the allocation of an IP address.
	// This is used when the dhcp4 or dhcp6 for a VirtualMachine is set and the VirtualMachine is waiting for the
	// relevant IP address to show up on the VM.
	VSphereVMVirtualMachineWaitingForIPAllocationReason = "WaitingForIPAllocation"

	// VSphereVMVirtualMachinePoweringOnReason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM is executing the power on sequence.
	VSphereVMVirtualMachinePoweringOnReason = "PoweringOn"

	// VSphereVMVirtualMachineProvisionedReason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM is provisioned.
	VSphereVMVirtualMachineProvisionedReason = clusterv1.ProvisionedReason

	// VSphereVMVirtualMachineTaskFailedReason surfaces when a task for the VirtualMachine that is controlled
	// by the VSphereVM failed; the reconcile look will automatically retry the operation,
	// but a user intervention might be required to fix the problem.
	VSphereVMVirtualMachineTaskFailedReason = "TaskFailed"

	// VSphereVMVirtualMachineNotFoundByBIOSUUIDReason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM can't be found by BIOS UUID.
	// Those kind of errors could be transient sometimes and failed VSphereVM are automatically
	// reconciled by the controller.
	VSphereVMVirtualMachineNotFoundByBIOSUUIDReason = "NotFoundByBIOSUUID"

	// VSphereVMVirtualMachineNotProvisionedReason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM is not provisioned.
	VSphereVMVirtualMachineNotProvisionedReason = clusterv1.NotProvisionedReason

	// VSphereVMVirtualMachineDeletingReason surfaces when the VirtualMachine that is controlled
	// by the VSphereVM is being deleted.
	VSphereVMVirtualMachineDeletingReason = clusterv1.DeletingReason
)

// VSphereVM's VCenterAvailable condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMVCenterAvailableCondition documents the availability of the VCenter hosting the VSphereVM.
	VSphereVMVCenterAvailableCondition = "VCenterAvailable"

	// VSphereVMVCenterAvailableReason documents the VCenter hosting the VSphereVM
	// being available.
	VSphereVMVCenterAvailableReason = clusterv1.AvailableReason

	// VSphereVMVCenterUnreachableReason documents the VCenter hosting the VSphereVM
	// cannot be reached.
	VSphereVMVCenterUnreachableReason = "VCenterUnreachable"
)

// VSphereVM's IPAddressClaimsFulfilled condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMIPAddressClaimsFulfilledCondition documents the status of claiming an IP address
	// from an IPAM provider.
	VSphereVMIPAddressClaimsFulfilledCondition = "IPAddressClaimsFulfilled"

	// VSphereVMIPAddressClaimsBeingCreatedReason documents that claims for the
	// IP addresses required by the VSphereVM are being created.
	VSphereVMIPAddressClaimsBeingCreatedReason = "IPAddressClaimsBeingCreated"

	// VSphereVMIPAddressClaimsWaitingForIPAddressReason documents that claims for the
	// IP addresses required by the VSphereVM are waiting for IP addresses.
	VSphereVMIPAddressClaimsWaitingForIPAddressReason = "WaitingForIPAddress"

	// VSphereVMIPAddressClaimsFulfilledReason documents that claims for the
	// IP addresses required by the VSphereVM are fulfilled.
	VSphereVMIPAddressClaimsFulfilledReason = "Fulfilled"

	// VSphereVMIPAddressClaimsNotFulfilledReason documents that claims for the
	// IP addresses required by the VSphereVM are not fulfilled.
	VSphereVMIPAddressClaimsNotFulfilledReason = "NotFulfilled"
)

// VSphereVM's GuestSoftPowerOffSucceeded condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMGuestSoftPowerOffSucceededCondition documents the status of performing guest initiated
	// graceful shutdown.
	VSphereVMGuestSoftPowerOffSucceededCondition string = "GuestSoftPowerOffSucceeded"

	// VSphereVMGuestSoftPowerOffInProgressReason documents that the guest receives
	// a graceful shutdown request.
	VSphereVMGuestSoftPowerOffInProgressReason = "InProgress"

	// VSphereVMGuestSoftPowerOffFailedReason documents that the graceful
	// shutdown request failed.
	VSphereVMGuestSoftPowerOffFailedReason = "Failed"

	// VSphereVMGuestSoftPowerOffSucceededReason documents that the graceful
	// shutdown request succeeded.
	VSphereVMGuestSoftPowerOffSucceededReason = "Succeeded"
)

// VSphereVM's PCIDevicesDetached condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereVMPCIDevicesDetachedCondition documents the status of the attached PCI devices on the VSphereVM.
	// It is a negative condition to notify the user that the device(s) is no longer attached to
	// the underlying VM and would require manual intervention to fix the situation.
	VSphereVMPCIDevicesDetachedCondition string = "PCIDevicesDetached"

	// VSphereVMPCIDevicesDetachedNotFoundReason documents the VSphereVM not having the PCI device attached during VM startup.
	// This would indicate that the PCI devices were removed out of band by an external entity.
	VSphereVMPCIDevicesDetachedNotFoundReason = "NotFound"
)

// VSphereVMSpec defines the desired state of VSphereVM.
type VSphereVMSpec struct {
	VirtualMachineCloneSpec `json:",inline"`

	// bootstrapRef is a reference to a bootstrap provider-specific resource
	// that holds configuration details.
	// This field is optional in case no bootstrap data is required to create
	// a VM.
	// +optional
	BootstrapRef VSphereVMBootstrapReference `json:"bootstrapRef,omitempty,omitzero"`

	// biosUUID is the VM's BIOS UUID that is assigned at runtime after
	// the VM has been created.
	// This field is required at runtime for other controllers that read
	// this CRD as unstructured data.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	BiosUUID string `json:"biosUUID,omitempty"`

	// powerOffMode describes the desired behavior when powering off a VM.
	//
	// There are three, supported power off modes: hard, soft, and
	// trySoft. The first mode, hard, is the equivalent of a physical
	// system's power cord being ripped from the wall. The soft mode
	// requires the VM's guest to have VM Tools installed and attempts to
	// gracefully shut down the VM. Its variant, trySoft, first attempts
	// a graceful shutdown, and if that fails or the VM is not in a powered off
	// state after reaching the GuestSoftPowerOffTimeoutSeconds, the VM is halted.
	//
	// If omitted, the mode defaults to hard.
	//
	// +optional
	PowerOffMode VirtualMachinePowerOpMode `json:"powerOffMode,omitempty"`

	// guestSoftPowerOffTimeoutSeconds sets the wait timeout for shutdown in the VM guest.
	// The VM will be powered off forcibly after the timeout if the VM is still
	// up and running when the PowerOffMode is set to trySoft.
	//
	// This parameter only applies when the PowerOffMode is set to trySoft.
	//
	// If omitted, the timeout defaults to 5 minutes.
	//
	// +optional
	// +kubebuilder:validation:Minimum=1
	GuestSoftPowerOffTimeoutSeconds int32 `json:"guestSoftPowerOffTimeoutSeconds,omitempty"`
}

// VSphereVMBootstrapReference is a reference to a Secret with the bootstrap data.
type VSphereVMBootstrapReference struct {
	// name of the Secret being referenced.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name,omitempty"`
}

// VSphereVMStatus defines the observed state of VSphereVM.
// +kubebuilder:validation:MinProperties=1
type VSphereVMStatus struct {
	// conditions represents the observations of a VSphereVM's current state.
	// Known condition types are Ready, VirtualMachineProvisioned, VCenterAvailable and IPAddressClaimsFulfilled,
	// GuestSoftPowerOffSucceeded, PCIDevicesDetached and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// host describes the hostname or IP address of the infrastructure host
	// that the VSphereVM is residing on.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Host string `json:"host,omitempty"`

	// ready is true when the provider resource is ready.
	// This field is required at runtime for other controllers that read
	// this CRD as unstructured data.
	// +optional
	Ready *bool `json:"ready,omitempty"`

	// addresses is a list of the VM's IP addresses.
	// This field is required at runtime for other controllers that read
	// this CRD as unstructured data.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=39
	Addresses []string `json:"addresses,omitempty"`

	// cloneMode is the type of clone operation used to clone this VM. Since
	// linkedClone is the default but fails gracefully if the source of the
	// clone has no snapshots, this field may be used to determine the actual
	// type of clone operation used to create this VM.
	// +optional
	CloneMode CloneMode `json:"cloneMode,omitempty"`

	// snapshot is the name of the snapshot from which the VM was cloned if
	// linkedClone is enabled.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Snapshot string `json:"snapshot,omitempty"`

	// retryAfter tracks the time we can retry queueing a task
	// +optional
	RetryAfter metav1.Time `json:"retryAfter,omitempty"`

	// taskRef is a managed object reference to a Task related to the machine.
	// This value is set automatically at runtime and should not be set or
	// modified by users.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	TaskRef string `json:"taskRef,omitempty"`

	// network returns the network status for each of the machine's configured
	// network interfaces.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	Network []NetworkStatus `json:"network,omitempty"`

	// moduleUUID is the unique identifier for the vCenter cluster module construct
	// which is used to configure anti-affinity. Objects with the same ModuleUUID
	// will be anti-affined, meaning that the vCenter DRS will best effort schedule
	// the VMs on separate hosts.
	// +optional
	// +kubebuilder:validation:MaxLength=2048
	ModuleUUID *string `json:"moduleUUID,omitempty"`

	// vmRef is the VM's Managed Object Reference on vSphere. It can be used by consumers
	// to programatically get this VM representation on vSphere in case of the need to retrieve informations.
	// This field is set once the machine is created and should not be changed
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	VMRef string `json:"vmRef,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *VSphereVMDeprecatedStatus `json:"deprecated,omitempty"`
}

// VSphereVMDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereVMDeprecatedStatus struct {
	// v1beta1 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	// +optional
	V1Beta1 *VSphereVMV1Beta1DeprecatedStatus `json:"v1beta1,omitempty"`
}

// VSphereVMV1Beta1DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereVMV1Beta1DeprecatedStatus struct {
	// conditions defines current service state of the VSphereVM.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// failureReason will be set in the event that there is a terminal problem
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
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

	// failureMessage will be set in the event that there is a terminal problem
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
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=10240
	FailureMessage *string `json:"failureMessage,omitempty"` //nolint:kubeapilinter // field will be removed when v1beta1 is removed
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspherevms,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VSphereVM is the Schema for the vspherevms API.
type VSphereVM struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereVM.
	// +required
	Spec VSphereVMSpec `json:"spec,omitempty,omitzero"`

	// status is the observed state of VSphereVM.
	// +optional
	Status VSphereVMStatus `json:"status,omitempty,omitzero"`
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (c *VSphereVM) GetV1Beta1Conditions() clusterv1.Conditions {
	if c.Status.Deprecated == nil || c.Status.Deprecated.V1Beta1 == nil {
		return nil
	}
	return c.Status.Deprecated.V1Beta1.Conditions
}

// SetV1Beta1Conditions sets the conditions on this object.
func (c *VSphereVM) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if c.Status.Deprecated == nil {
		c.Status.Deprecated = &VSphereVMDeprecatedStatus{}
	}
	if c.Status.Deprecated.V1Beta1 == nil {
		c.Status.Deprecated.V1Beta1 = &VSphereVMV1Beta1DeprecatedStatus{}
	}
	c.Status.Deprecated.V1Beta1.Conditions = conditions
}

// GetConditions returns the set of conditions for this object.
func (c *VSphereVM) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (c *VSphereVM) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
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
