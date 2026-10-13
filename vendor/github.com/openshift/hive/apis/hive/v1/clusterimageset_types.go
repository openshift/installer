package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterImageSetSpec defines the desired state of ClusterImageSet
type ClusterImageSetSpec struct {
	// ReleaseImage is the image that contains the payload to use when installing
	// a cluster.
	ReleaseImage string `json:"releaseImage"`
}

// ClusterImageSetStatus defines the observed state of ClusterImageSet
type ClusterImageSetStatus struct{}

// +genclient:nonNamespaced
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterImageSet is the Schema for the clusterimagesets API
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Release",type="string",JSONPath=".spec.releaseImage"
// +kubebuilder:resource:path=clusterimagesets,shortName=imgset,scope=Cluster
type ClusterImageSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterImageSetSpec   `json:"spec,omitempty"`
	Status ClusterImageSetStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterImageSetList contains a list of ClusterImageSet
type ClusterImageSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterImageSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterImageSet{}, &ClusterImageSetList{})
}
