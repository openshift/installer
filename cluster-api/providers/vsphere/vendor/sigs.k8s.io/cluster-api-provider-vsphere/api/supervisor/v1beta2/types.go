/*
Copyright 2026 The Kubernetes Authors.

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

import clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

// VSphereMachineTemplateResource describes the data needed to create a VSphereMachine from a template.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineTemplateResource struct {
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec is the specification of the desired behavior of the machine.
	// +optional
	Spec VSphereMachineSpec `json:"spec,omitempty,omitzero"`
}

// VSphereMachinePhase is the phase of the VSphereMachine.
// +kubebuilder:validation:Enum=NotFound;Created;PoweredOn;Pending;Ready;Deleting;Error
type VSphereMachinePhase string

const (
	// VSphereMachinePhaseNotFound is the string representing a VM that cannot be located.
	VSphereMachinePhaseNotFound = VSphereMachinePhase("NotFound")

	// VSphereMachinePhaseCreated is the string representing a VM that's been created.
	VSphereMachinePhaseCreated = VSphereMachinePhase("Created")

	// VSphereMachinePhasePoweredOn is the string representing a VM that has successfully powered on.
	VSphereMachinePhasePoweredOn = VSphereMachinePhase("PoweredOn")

	// VSphereMachinePhasePending is the string representing a VM with an in-flight task.
	VSphereMachinePhasePending = VSphereMachinePhase("Pending")

	// VSphereMachinePhaseReady is the string representing a powered-on VM with reported IP addresses.
	VSphereMachinePhaseReady = VSphereMachinePhase("Ready")

	// VSphereMachinePhaseDeleting is the string representing a machine that still exists, but has a deletionTimestamp.
	// Note that once a VirtualMachine is finally deleted, its state will be VSphereMachinePhaseNotFound.
	VSphereMachinePhaseDeleting = VSphereMachinePhase("Deleting")

	// VSphereMachinePhaseError is reported if an error occurs determining the status.
	VSphereMachinePhaseError = VSphereMachinePhase("Error")
)

// VirtualMachinePowerOpMode represents the various power operation modes
// when powering off or suspending a VM.
// +kubebuilder:validation:Enum=hard;soft;trySoft
type VirtualMachinePowerOpMode string

const (
	// VirtualMachinePowerOpModeHard indicates to halt a VM when powering it
	// off or when suspending a VM to not involve the guest.
	VirtualMachinePowerOpModeHard VirtualMachinePowerOpMode = "hard"

	// VirtualMachinePowerOpModeSoft indicates to ask VM Tools running
	// inside of a VM's guest to shutdown the guest gracefully when powering
	// off a VM or when suspending a VM to allow the guest to participate.
	//
	// If this mode is set on a VM whose guest does not have VM Tools or if
	// VM Tools is present but the operation fails, the VM may never realize
	// the desired power state. This can prevent a VM from being deleted as well
	// as many other unexpected issues. It is recommended to use trySoft
	// instead.
	VirtualMachinePowerOpModeSoft VirtualMachinePowerOpMode = "soft"

	// VirtualMachinePowerOpModeTrySoft indicates to first attempt a Soft
	// operation and fall back to hard if VM Tools is not present in the guest,
	// if the soft operation fails, or if the VM is not in the desired power
	// state within the configured timeout (default 5m).
	VirtualMachinePowerOpModeTrySoft VirtualMachinePowerOpMode = "trySoft"
)
