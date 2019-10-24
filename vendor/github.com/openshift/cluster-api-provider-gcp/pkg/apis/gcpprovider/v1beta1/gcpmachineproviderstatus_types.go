package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains GCP-specific status information.
// +k8s:openapi-gen=true
type GCPMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// InstanceID is the ID of the instance in GCP
	// +optional
	InstanceID *string `json:"instanceId,omitempty"`

	// InstanceState is the provisioning state of the GCP Instance.
	// +optional
	InstanceState *string `json:"instanceState,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []GCPMachineProviderCondition `json:"conditions,omitempty"`
}

// GCPMachineProviderConditionType is a valid value for GCPMachineProviderCondition.Type
type GCPMachineProviderConditionType string

// Valid conditions for an GCP machine instance
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated GCPMachineProviderConditionType = "MachineCreated"
)

// GCPMachineProviderCondition is a condition in a GCPMachineProviderStatus
type GCPMachineProviderCondition struct {
	// Type is the type of the condition.
	Type GCPMachineProviderConditionType `json:"type"`
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
	SchemeBuilder.Register(&GCPMachineProviderStatus{})
}
