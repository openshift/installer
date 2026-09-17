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
	// MachineFinalizer allows ReconcileVSphereMachine to clean up VSphere
	// resources associated with VSphereMachine before removing it from the
	// API Server.
	MachineFinalizer = "vspheremachine.infrastructure.cluster.x-k8s.io"
)

// VSphereMachine's Ready condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereMachineReadyCondition is true if the VSphereMachine's deletionTimestamp is not set, VSphereMachine's
	// VirtualMachineProvisioned is true.
	VSphereMachineReadyCondition = clusterv1.ReadyCondition

	// VSphereMachineReadyReason surfaces when the VSphereMachine readiness criteria is met.
	VSphereMachineReadyReason = clusterv1.ReadyReason

	// VSphereMachineNotReadyReason surfaces when the VSphereMachine readiness criteria is not met.
	VSphereMachineNotReadyReason = clusterv1.NotReadyReason

	// VSphereMachineReadyUnknownReason surfaces when at least one VSphereMachine readiness criteria is unknown
	// and no VSphereMachine readiness criteria is not met.
	VSphereMachineReadyUnknownReason = clusterv1.ReadyUnknownReason
)

// VSphereMachine's VirtualMachineProvisioned condition and corresponding reasons that will be used in v1Beta2 API version.
//
// NOTE:
//   - In supervisor mode, before creating the VM the VirtualMachine goes trough a series of preflight checks; if one is failing, the
//     reason for this failure and the message are surfaced in the VSphereMachine's VirtualMachineProvisioned condition.
//   - In govmomi mode, in some cases, reason and message from the VSphereVM are surfaced in the VSphereMachine's
//     VirtualMachineProvisioned condition.
const (
	// VSphereMachineVirtualMachineProvisionedCondition documents the status of the VirtualMachine that is controlled
	// by the VSphereMachine.
	VSphereMachineVirtualMachineProvisionedCondition = "VirtualMachineProvisioned"

	// VSphereMachineVirtualMachineWaitingForClusterInfrastructureReadyReason documents the VirtualMachine that is controlled
	// by the VSphereMachine waiting for the cluster infrastructure to be ready.
	// Note: This reason is used only in govmomi mode.
	VSphereMachineVirtualMachineWaitingForClusterInfrastructureReadyReason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// VSphereMachineVirtualMachineWaitingForControlPlaneInitializedReason documents the VirtualMachine that is controlled
	// by the VSphereMachine waiting for the control plane to be initialized.
	VSphereMachineVirtualMachineWaitingForControlPlaneInitializedReason = clusterv1.WaitingForControlPlaneInitializedReason

	// VSphereMachineVirtualMachineWaitingForBootstrapDataReason documents the VirtualMachine that is controlled
	// by the VSphereMachine waiting for the bootstrap data to be ready.
	VSphereMachineVirtualMachineWaitingForBootstrapDataReason = clusterv1.WaitingForBootstrapDataReason

	// VSphereMachineVirtualMachineProvisioningReason surfaces when the VirtualMachine that is controlled
	// by the VSphereMachine is provisioning.
	// Note: This reason is used only in supervisor mode.
	VSphereMachineVirtualMachineProvisioningReason = "Provisioning"

	// VSphereMachineVirtualMachinePoweringOnReason surfaces when the VirtualMachine that is controlled
	// by the VSphereMachine is executing the power on sequence.
	// Note: This reason is used only in supervisor mode.
	VSphereMachineVirtualMachinePoweringOnReason = "PoweringOn"

	// VSphereMachineVirtualMachineWaitingForVirtualMachineGroupReason surfaces that the VirtualMachine
	// is waiting for its corresponding VirtualMachineGroup to be created and to include this VM as a member.
	VSphereMachineVirtualMachineWaitingForVirtualMachineGroupReason = "WaitingForVirtualMachineGroup"

	// VSphereMachineVirtualMachineWaitingForNetworkAddressReason surfaces when the VirtualMachine that is controlled
	// by the VSphereMachine waiting for the machine network settings to be reported after machine being powered on.
	VSphereMachineVirtualMachineWaitingForNetworkAddressReason = "WaitingForNetworkAddress"

	// VSphereMachineVirtualMachineWaitingForBIOSUUIDReason surfaces when the VirtualMachine that is controlled
	// by the VSphereMachine waiting for the machine to have a BIOS UUID.
	// Note: This reason is used only in supervisor mode.
	VSphereMachineVirtualMachineWaitingForBIOSUUIDReason = "WaitingForBIOSUUID"

	// VSphereMachineVirtualMachineProvisionedReason surfaces when the VirtualMachine that is controlled
	// by the VSphereMachine is provisioned.
	VSphereMachineVirtualMachineProvisionedReason = clusterv1.ProvisionedReason

	// VSphereMachineVirtualMachineNotProvisionedReason surfaces when the VirtualMachine that is controlled
	// by the VSphereMachine is not provisioned.
	VSphereMachineVirtualMachineNotProvisionedReason = clusterv1.NotProvisionedReason

	// VSphereMachineVirtualMachineDeletingReason surfaces when the VirtualMachine that is controlled
	// by the VSphereMachine is being deleted.
	VSphereMachineVirtualMachineDeletingReason = clusterv1.DeletingReason
)

// VSphereMachineSpec defines the desired state of VSphereMachine.
type VSphereMachineSpec struct {
	VirtualMachineCloneSpec `json:",inline"`

	// providerID is the virtual machine's BIOS UUID formatted as
	// vsphere://12345678-1234-1234-1234-123456789abc
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	ProviderID string `json:"providerID,omitempty"`

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

	// naming allows configuring the naming strategy used when calculating the name of the VSphereVM.
	// +optional
	Naming VSphereVMNamingSpec `json:"naming,omitempty,omitzero"`
}

// VSphereVMNamingSpec defines the naming strategy for the VSphereVMs.
// +kubebuilder:validation:MinProperties=1
type VSphereVMNamingSpec struct {
	// template defines the template to use for generating the name of the VSphereVM object.
	// If not defined, it will fall back to `{{ .machine.name }}`.
	// The templating has the following data available:
	// * `.machine.name`: The name of the Machine object.
	// The templating also has the following funcs available:
	// * `trimSuffix`: same as strings.TrimSuffix
	// * `trunc`: truncates a string, e.g. `trunc 2 "hello"` or `trunc -2 "hello"`
	// Notes:
	// * While the template offers some flexibility, we would like the name to link to the Machine name
	//   to ensure better user experience when troubleshooting
	// * Generated names must be valid Kubernetes names as they are used to create a VSphereVM object
	//   and usually also as the name of the Node object.
	// * Names are automatically truncated at 63 characters. Please note that this can lead to name conflicts,
	//   so we highly recommend to use a template which leads to a name shorter than 63 characters.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Template string `json:"template,omitempty"`
}

// VSphereMachineStatus defines the observed state of VSphereMachine.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineStatus struct {
	// conditions represents the observations of a VSphereMachine's current state.
	// Known condition types are Ready, VirtualMachineProvisioned and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// initialization provides observations of the VSphereMachine initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial Machine provisioning.
	// +optional
	Initialization VSphereMachineInitializationStatus `json:"initialization,omitempty,omitzero"`

	// addresses contains the VSphere instance associated addresses.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// network returns the network status for each of the machine's configured
	// network interfaces.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	Network []NetworkStatus `json:"network,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *VSphereMachineDeprecatedStatus `json:"deprecated,omitempty"`
}

// VSphereMachineInitializationStatus provides observations of the VSphereMachine initialization process.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineInitializationStatus struct {
	// provisioned is true when the infrastructure provider reports that the Machine's infrastructure is fully provisioned.
	// NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Machine provisioning.
	// +optional
	Provisioned *bool `json:"provisioned,omitempty"`
}

// VSphereMachineDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereMachineDeprecatedStatus struct {
	// v1beta1 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	V1Beta1 *VSphereMachineV1Beta1DeprecatedStatus `json:"v1beta1,omitempty"`
}

// VSphereMachineV1Beta1DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereMachineV1Beta1DeprecatedStatus struct {
	// conditions defines current service state of the VSphereMachine.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// failureReason will be set in the event that there is a terminal problem
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
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

	// failureMessage will be set in the event that there is a terminal problem
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
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=10240
	FailureMessage *string `json:"failureMessage,omitempty"` //nolint:kubeapilinter // field will be removed when v1beta1 is removed
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspheremachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels['cluster\\.x-k8s\\.io/cluster-name']",description="Cluster"
// +kubebuilder:printcolumn:name="Provider ID",type="string",JSONPath=".spec.providerID",description="Provider ID",priority=10
// +kubebuilder:printcolumn:name="Template",type="string",JSONPath=".spec.template",description="Template is the name, inventory path, managed object reference or the managed object ID of the template used to clone the VM",priority=10
// +kubebuilder:printcolumn:name="Paused",type="string",JSONPath=`.status.conditions[?(@.type=="Paused")].status`,description="Reconciliation paused",priority=10
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.initialization.provisioned",description="VSphereMachine is provisioned"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of VSphereMachine"

// VSphereMachine is the Schema for the vspheremachines API.
type VSphereMachine struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereMachine.
	// +required
	Spec VSphereMachineSpec `json:"spec,omitempty,omitzero"`

	// status is the observed state of VSphereMachine.
	// +optional
	Status VSphereMachineStatus `json:"status,omitempty,omitzero"`
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (c *VSphereMachine) GetV1Beta1Conditions() clusterv1.Conditions {
	if c.Status.Deprecated == nil || c.Status.Deprecated.V1Beta1 == nil {
		return nil
	}
	return c.Status.Deprecated.V1Beta1.Conditions
}

// SetV1Beta1Conditions sets the conditions on this object.
func (c *VSphereMachine) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if c.Status.Deprecated == nil {
		c.Status.Deprecated = &VSphereMachineDeprecatedStatus{}
	}
	if c.Status.Deprecated.V1Beta1 == nil {
		c.Status.Deprecated.V1Beta1 = &VSphereMachineV1Beta1DeprecatedStatus{}
	}
	c.Status.Deprecated.V1Beta1.Conditions = conditions
}

// GetConditions returns the set of conditions for this object.
func (c *VSphereMachine) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (c *VSphereMachine) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VSphereMachineList contains a list of VSphereMachine.
type VSphereMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereMachine `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereMachine{}, &VSphereMachineList{})
}
