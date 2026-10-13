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

// VSphereMachineTemplateResource describes the data needed to create a VSphereMachine from a template.
type VSphereMachineTemplateResource struct {
	// Spec is the specification of the desired behavior of the machine.
	Spec VSphereMachineSpec `json:"spec"`
}

// VirtualMachineState describes the state of a VM.
type VirtualMachineState string

const (
	// VirtualMachineStateNotFound is the string representing a VM that cannot be located.
	VirtualMachineStateNotFound = VirtualMachineState("notfound")

	// VirtualMachineStateCreated is the string representing a VM that's been created.
	VirtualMachineStateCreated = VirtualMachineState("created")

	// VirtualMachineStatePoweredOn is the string representing a VM that has successfully powered on.
	VirtualMachineStatePoweredOn = VirtualMachineState("poweredon")

	// VirtualMachineStatePending is the string representing a VM with an in-flight task.
	VirtualMachineStatePending = VirtualMachineState("pending")

	// VirtualMachineStateReady is the string representing a powered-on VM with reported IP addresses.
	VirtualMachineStateReady = VirtualMachineState("ready")

	// VirtualMachineStateDeleting is the string representing a machine that still exists, but has a deleteTimestamp
	// Note that once a VirtualMachine is finally deleted, its state will be VirtualMachineStateNotFound.
	VirtualMachineStateDeleting = VirtualMachineState("deleting")

	// VirtualMachineStateError is reported if an error occurs determining the status.
	VirtualMachineStateError = VirtualMachineState("error")
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
