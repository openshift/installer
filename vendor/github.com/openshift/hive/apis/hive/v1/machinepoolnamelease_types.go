package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachinePoolNameLeaseSpec is a minimal resource for obtaining unique machine pool names of a limited length.
type MachinePoolNameLeaseSpec struct {
}

// MachinePoolNameLeaseStatus defines the observed state of MachinePoolNameLease.
type MachinePoolNameLeaseStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachinePoolNameLease is the Schema for the MachinePoolNameLeases API. This resource is mostly empty
// as we're primarily relying on the name to determine if a lease is available.
// Note that not all cloud providers require the use of a lease for naming, at present this
// is only required for GCP where we're extremely restricted on name lengths.
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="MachinePool",type="string",JSONPath=".metadata.labels.hive\\.openshift\\.io/machine-pool-name"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.hive\\.openshift\\.io/cluster-deployment-name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Namespaced
type MachinePoolNameLease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachinePoolNameLeaseSpec   `json:"spec,omitempty"`
	Status MachinePoolNameLeaseStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachinePoolNameLeaseList contains a list of MachinePoolNameLeases.
type MachinePoolNameLeaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachinePoolNameLease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachinePoolNameLease{}, &MachinePoolNameLeaseList{})
}
