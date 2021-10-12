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
)

// Valid conditions for an alibabacloud machine instance.
const (
	// MachineCreation indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreation AlibabaCloudMachineProviderConditionType = "MachineCreation"
)

// AlibabaCloudMachineProviderConditionReason is reason for the condition's last transition.
type AlibabaCloudMachineProviderConditionReason string

const (
	// MachineCreationSucceeded indicates machine creation success.
	MachineCreationSucceeded AlibabaCloudMachineProviderConditionReason = "MachineCreationSucceeded"
	// MachineCreationFailed indicates machine creation failure.
	MachineCreationFailed AlibabaCloudMachineProviderConditionReason = "MachineCreationFailed"
)

// AlibabaCloudMachineProviderStatus is the Schema for the alibabacloudmachineproviderconfig API
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AlibabaCloudMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// InstanceID is the instance ID of the machine created in alibabacloud
	// +optional
	InstanceID *string `json:"instanceId,omitempty"`

	// InstanceState is the state of the alibabacloud instance for this machine
	// +optional
	InstanceState *string `json:"instanceState,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []AlibabaCloudMachineProviderCondition `json:"conditions,omitempty"`
}

// AlibabaCloudMachineProviderConditionType is a valid value for AlibabaCloudMachineProviderCondition.Type
type AlibabaCloudMachineProviderConditionType string

// AlibabaCloudMachineProviderCondition is a condition in a AlibabaCloudMachineProviderStatus
type AlibabaCloudMachineProviderCondition struct {
	// Type is the type of the condition.
	Type AlibabaCloudMachineProviderConditionType `json:"type"`

	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`

	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason AlibabaCloudMachineProviderConditionReason `json:"reason"`
	// Message is a human-readable message indicating details about last transition.

	// +optional
	Message string `json:"message"`
}

func init() {
	SchemeBuilder.Register(&AlibabaCloudMachineProviderStatus{})
}
