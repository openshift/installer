package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OperatorHubSpec defines the desired state of OperatorHub
type OperatorHubSpec struct {
	// hubSources is the list of default OperatorSources and their configuration
	HubSources []HubSource `json:"hubSources,omitempty"`
}

// OperatorHubStatus defines the observed state of OperatorHub
type OperatorHubStatus struct {
	// hubSourcesStatus encapsulates the result applying the configuration
	HubSourcesStatus []HubSourceStatus `json:"hubSourcesStatus,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatorHub is the Schema for the operatorhubs API
// +kubebuilder:subresource:status
// +genclient:nonNamespaced
type OperatorHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   OperatorHubSpec   `json:"spec"`
	Status OperatorHubStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatorHubList contains a list of OperatorHub
type OperatorHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []OperatorHub `json:"items"`
}

// HubSource is used to specify the OperatorSource and its configuration
type HubSource struct {
	// name is the name of one of the default OperatorSources
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:Required
	Name string `json:"name"`
	// disabled is used to disable a default OperatorSource on cluster
	// +kubebuilder:Required
	Disabled bool `json:"disabled"`
}

// HubSourceStatus is used to reflect the current state of applying the
// configuration to a default source
type HubSourceStatus struct {
	// name is the name of one of the default OperatorSources
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:Required
	Name string `json:"name,omitempty"`
	// configuration is the state of the default OperatorSources configuration
	Configuration map[string]string `json:"configuration,omitempty"`
	// status indicates success or failure in applying the configuration
	Status string `json:"status,omitempty"`
	// message provides more information regarding failures
	Message string `json:"message,omitempty"`
}
