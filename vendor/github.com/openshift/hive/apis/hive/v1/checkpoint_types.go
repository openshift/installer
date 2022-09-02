package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CheckpointSpec defines the metadata around the Hive objects state in the namespace at the time of the last backup.
type CheckpointSpec struct {
	// LastBackupChecksum is the checksum of all Hive objects in the namespace at the time of the last backup.
	LastBackupChecksum string `json:"lastBackupChecksum"`

	// LastBackupTime is the last time we performed a backup of the namespace
	LastBackupTime metav1.Time `json:"lastBackupTime"`

	// LastBackupRef is a reference to last backup object created
	LastBackupRef BackupReference `json:"lastBackupRef"`
}

// BackupReference is a reference to a backup resource
type BackupReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// CheckpointStatus defines the observed state of Checkpoint
type CheckpointStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Checkpoint is the Schema for the backup of Hive objects.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
type Checkpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CheckpointSpec   `json:"spec,omitempty"`
	Status CheckpointStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CheckpointList contains a list of Checkpoint
type CheckpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Checkpoint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Checkpoint{}, &CheckpointList{})
}
