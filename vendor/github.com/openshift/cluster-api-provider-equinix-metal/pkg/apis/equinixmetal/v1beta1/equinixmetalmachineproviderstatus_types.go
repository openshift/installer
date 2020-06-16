package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EquinixMetalMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains EquinixMetal-specific status information.
// +k8s:openapi-gen=true
type EquinixMetalMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// InstanceID is the ID of the instance in EquinixMetal
	// +optional
	InstanceID *string `json:"instanceId,omitempty"`

	// InstanceState is the provisioning state of the EquinixMetal Instance.
	// +optional
	InstanceState *string `json:"instanceState,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []EquinixMetalMachineProviderCondition `json:"conditions,omitempty"`
}

// EquinixMetalMachineProviderConditionType is a valid value for EquinixMetalMachineProviderCondition.Type.
type EquinixMetalMachineProviderConditionType string

// Valid conditions for an EquinixMetal machine instance.
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated EquinixMetalMachineProviderConditionType = "MachineCreated"
)

// EquinixMetalMachineProviderCondition is a condition in a EquinixMetalMachineProviderStatus.
type EquinixMetalMachineProviderCondition struct {
	// Type is the type of the condition.
	Type EquinixMetalMachineProviderConditionType `json:"type"`
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
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&EquinixMetalMachineProviderStatus{})
}
