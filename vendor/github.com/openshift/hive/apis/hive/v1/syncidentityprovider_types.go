package v1

import (
	openshiftapiv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SyncIdentityProviderCommonSpec defines the identity providers to sync
type SyncIdentityProviderCommonSpec struct {
	//IdentityProviders is an ordered list of ways for a user to identify themselves
	// +required
	IdentityProviders []openshiftapiv1.IdentityProvider `json:"identityProviders"`
}

// SelectorSyncIdentityProviderSpec defines the SyncIdentityProviderCommonSpec to sync to
// ClusterDeploymentSelector indicating which clusters the SelectorSyncIdentityProvider applies
// to in any namespace.
type SelectorSyncIdentityProviderSpec struct {
	SyncIdentityProviderCommonSpec `json:",inline"`

	// ClusterDeploymentSelector is a LabelSelector indicating which clusters the SelectorIdentityProvider
	// applies to in any namespace.
	// +optional
	ClusterDeploymentSelector metav1.LabelSelector `json:"clusterDeploymentSelector,omitempty"`
}

// SyncIdentityProviderSpec defines the SyncIdentityProviderCommonSpec identity providers to sync along with
// ClusterDeploymentRefs indicating which clusters the SyncIdentityProvider applies to in the
// SyncIdentityProvider's namespace.
type SyncIdentityProviderSpec struct {
	SyncIdentityProviderCommonSpec `json:",inline"`

	// ClusterDeploymentRefs is the list of LocalObjectReference indicating which clusters the
	// SyncSet applies to in the SyncSet's namespace.
	// +required
	ClusterDeploymentRefs []corev1.LocalObjectReference `json:"clusterDeploymentRefs"`
}

// IdentityProviderStatus defines the observed state of SyncSet
type IdentityProviderStatus struct {
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SelectorSyncIdentityProvider is the Schema for the SelectorSyncSet API
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Cluster
type SelectorSyncIdentityProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SelectorSyncIdentityProviderSpec `json:"spec,omitempty"`
	Status IdentityProviderStatus           `json:"status,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SyncIdentityProvider is the Schema for the SyncIdentityProvider API
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced
type SyncIdentityProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SyncIdentityProviderSpec `json:"spec,omitempty"`
	Status IdentityProviderStatus   `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SelectorSyncIdentityProviderList contains a list of SelectorSyncIdentityProviders
type SelectorSyncIdentityProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SelectorSyncIdentityProvider `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SyncIdentityProviderList contains a list of SyncIdentityProviders
type SyncIdentityProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SyncIdentityProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&SyncIdentityProvider{},
		&SyncIdentityProviderList{},
		&SelectorSyncIdentityProvider{},
		&SelectorSyncIdentityProviderList{},
	)
}
