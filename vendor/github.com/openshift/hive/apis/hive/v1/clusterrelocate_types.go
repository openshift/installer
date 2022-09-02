package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterRelocateSpec defines the relocation of clusters from one Hive instance to another.
type ClusterRelocateSpec struct {
	// KubeconfigSecretRef is a reference to the secret containing the kubeconfig for the destination Hive instance.
	// The kubeconfig must be in a data field where the key is "kubeconfig".
	KubeconfigSecretRef KubeconfigSecretReference `json:"kubeconfigSecretRef"`

	// ClusterDeploymentSelector is a LabelSelector indicating which clusters will be relocated.
	ClusterDeploymentSelector metav1.LabelSelector `json:"clusterDeploymentSelector"`
}

// KubeconfigSecretReference is a reference to a secret containing the kubeconfig for a remote cluster.
type KubeconfigSecretReference struct {
	// Name is the name of the secret.
	Name string `json:"name"`
	// Namespace is the namespace where the secret lives.
	Namespace string `json:"namespace"`
}

// ClusterRelocateStatus defines the observed state of ClusterRelocate.
type ClusterRelocateStatus struct{}

// +genclient:nonNamespaced
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterRelocate is the Schema for the ClusterRelocates API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Selector",type="string",JSONPath=".spec.clusterDeploymentSelector"
// +kubebuilder:resource:path=clusterrelocates
type ClusterRelocate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterRelocateSpec   `json:"spec,omitempty"`
	Status ClusterRelocateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterRelocateList contains a list of ClusterRelocate
type ClusterRelocateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterRelocate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterRelocate{}, &ClusterRelocateList{})
}
