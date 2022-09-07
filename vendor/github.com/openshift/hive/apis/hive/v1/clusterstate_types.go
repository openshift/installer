package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/openshift/api/config/v1"
)

// ClusterStateSpec defines the desired state of ClusterState
type ClusterStateSpec struct {
}

// ClusterStateStatus defines the observed state of ClusterState
type ClusterStateStatus struct {
	// LastUpdated is the last time that operator state was updated
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`

	// ClusterOperators contains the state for every cluster operator in the
	// target cluster
	ClusterOperators []ClusterOperatorState `json:"clusterOperators,omitempty"`
}

// ClusterOperatorState summarizes the status of a single cluster operator
type ClusterOperatorState struct {
	// Name is the name of the cluster operator
	Name string `json:"name"`

	// Conditions is the set of conditions in the status of the cluster operator
	// on the target cluster
	Conditions []configv1.ClusterOperatorStatusCondition `json:"conditions,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterState is the Schema for the clusterstates API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
type ClusterState struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterStateSpec   `json:"spec,omitempty"`
	Status ClusterStateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterStateList contains a list of ClusterState
type ClusterStateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterState `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterState{}, &ClusterStateList{})
}
