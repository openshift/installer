/*
Copyright 2021 Nutanix Inc.

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
)

// NutanixMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains nutanix-specific status information.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NutanixMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// The Nutanix VM's UUID
	// +optional
	VmUUID *string `json:"vmUUID,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []NutanixMachineProviderCondition `json:"conditions,omitempty"`
}

// NutanixMachineProviderConditionType is a valid value for NutanixMachineProviderCondition.Type
type NutanixMachineProviderConditionType string

// Valid conditions for an Nutanix machine instance.
const (
	// MachineCreation indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreation NutanixMachineProviderConditionType = "MachineCreation"
	MachineUpdate   NutanixMachineProviderConditionType = "MachineUpdate"
	MachineDeletion NutanixMachineProviderConditionType = "MachineDeletion"
)

// NutanixMachineProviderConditionReason is reason for the condition's last transition.
type NutanixMachineProviderConditionReason string

const (
	// MachineCreationSucceeded indicates machine creation success.
	MachineCreationSucceeded NutanixMachineProviderConditionReason = "MachineCreationSucceeded"
	// MachineCreationFailed indicates machine creation failure.
	MachineCreationFailed NutanixMachineProviderConditionReason = "MachineCreationFailed"
	// MachineCreationSucceeded indicates machine update success.
	MachineUpdateSucceeded NutanixMachineProviderConditionReason = "MachineUpdateSucceeded"
	// MachineCreationFailed indicates machine update failure.
	MachineUpdateFailed NutanixMachineProviderConditionReason = "MachineUpdateFailed"
	// MachineCreationSucceeded indicates machine deletion success.
	MachineDeletionSucceeded NutanixMachineProviderConditionReason = "MachineDeletionSucceeded"
	// MachineCreationFailed indicates machine deletion failure.
	MachineDeletionFailed NutanixMachineProviderConditionReason = "MachineDeletionFailed"
)

// NutanixMachineProviderCondition is a condition in a NutanixMachineProviderStatus.
type NutanixMachineProviderCondition struct {
	// Type is the type of the condition.
	Type NutanixMachineProviderConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason NutanixMachineProviderConditionReason `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}
