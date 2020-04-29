package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereMachineProviderConditionType is a valid value for VSphereMachineProviderCondition.Type.
type VSphereMachineProviderConditionType string

// Valid conditions for an vSphere machine instance.
const (
	// MachineCreation indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreation VSphereMachineProviderConditionType = "MachineCreation"
)

// VSphereMachineProviderConditionReason is reason for the condition's last transition.
type VSphereMachineProviderConditionReason string

const (
	// MachineCreationSucceeded indicates machine creation success.
	MachineCreationSucceeded VSphereMachineProviderConditionReason = "MachineCreationSucceeded"
	// MachineCreationFailed indicates machine creation failure.
	MachineCreationFailed VSphereMachineProviderConditionReason = "MachineCreationFailed"
)

// VSphereMachineProviderCondition is a condition in a VSphereMachineProviderStatus.
type VSphereMachineProviderCondition struct {
	// Type is the type of the condition.
	Type VSphereMachineProviderConditionType `json:"type"`
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
	Reason VSphereMachineProviderConditionReason `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VSphereMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains VSphere-specific status information.
// +k8s:openapi-gen=true
type VSphereMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// TODO: populate what we need here:
	// InstanceID is the ID of the instance in VSphere
	// +optional
	InstanceID *string `json:"instanceId,omitempty"`

	// InstanceState is the provisioning state of the VSphere Instance.
	// +optional
	InstanceState *string `json:"instanceState,omitempty"`
	//
	// TaskRef?
	// Ready?
	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []VSphereMachineProviderCondition `json:"conditions,omitempty"`

	// TaskRef is a managed object reference to a Task related to the machine.
	// This value is set automatically at runtime and should not be set or
	// modified by users.
	// +optional
	TaskRef string `json:"taskRef,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&VSphereMachineProviderStatus{})
}
