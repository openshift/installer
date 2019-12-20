package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VSphereMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains VSphere-specific status information.
// +k8s:openapi-gen=true
type VSphereMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// TODO: populate what we need here:
	// InstanceID is the ID of the instance in VSphere
	// +optional
	//InstanceID *string `json:"instanceId,omitempty"`

	// InstanceState is the provisioning state of the VSphere Instance.
	// +optional
	//InstanceState *string `json:"instanceState,omitempty"`
	//
	// TaskRef?
	// Ready?
	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	// Conditions []VSphereMachineProviderCondition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&VSphereMachineProviderStatus{})
}
