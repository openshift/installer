/*
Copyright 2021.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IBMCloudMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains IBMCloud-specific status information.
type IBMCloudMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// InstanceID is the instance ID of the machine created in IBM Cloud
	// +optional
	InstanceID *string `json:"instanceId,omitempty"`

	// InstanceState is the state of the IBM Cloud instance for this machine
	// +optional
	InstanceState *string `json:"instanceState,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []IBMCloudMachineProviderCondition `json:"conditions,omitempty"`
}

// IBMCloudMachineProviderConditionType is a valid value for IBMCloudMachineProviderCondition.Type
type IBMCloudMachineProviderConditionType string

// Valid conditions for an IBM Cloud machine instance
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated IBMCloudMachineProviderConditionType = "MachineCreated"
)

// IBMCloudMachineProviderConditionReason is reason for the condition's last transition.
type IBMCloudMachineProviderConditionReason string

const (
	// MachineCreationSucceeded indicates machine creation success.
	MachineCreationSucceeded IBMCloudMachineProviderConditionReason = "MachineCreationSucceeded"
	// MachineCreationFailed indicates machine creation failure.
	MachineCreationFailed IBMCloudMachineProviderConditionReason = "MachineCreationFailed"
)

// IBMCloudMachineProviderCondition is a condition in a IBMCloudMachineProviderStatus.
type IBMCloudMachineProviderCondition struct {
	// Type is the type of the condition.
	Type IBMCloudMachineProviderConditionType `json:"type"`
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
	Reason IBMCloudMachineProviderConditionReason `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&IBMCloudMachineProviderStatus{})
}
