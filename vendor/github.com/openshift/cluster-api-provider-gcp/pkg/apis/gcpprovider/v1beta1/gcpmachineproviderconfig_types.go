package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GCPMachineProviderSpec is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an GCP virtual machine. It is used by the GCP machine actuator to create a single Machine.
// +k8s:openapi-gen=true
type GCPMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is a reference to the secret with GCP credentials.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`

	CanIPForward       bool                   `json:"canIPForward"`
	DeletionProtection bool                   `json:"deletionProtection"`
	Disks              []*GCPDisk             `json:"disks,omitempty"`
	Labels             map[string]string      `json:"labels,omitempty"`
	Metadata           []*GCPMetadata         `json:"gcpMetadata,omitempty"`
	NetworkInterfaces  []*GCPNetworkInterface `json:"networkInterfaces,omitempty"`
	ServiceAccounts    []GCPServiceAccount    `json:"serviceAccounts"`
	Tags               []string               `json:"tags,omitempty"`
	TargetPools        []string               `json:"targetPools,omitempty"`
	MachineType        string                 `json:"machineType"`
	Region             string                 `json:"region"`
	Zone               string                 `json:"zone"`
	ProjectID          string                 `json:"projectID,omitempty"`

	// Preemptible indicates if created instance is preemptible
	Preemptible bool `json:"preemptible,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&GCPMachineProviderSpec{})
}

// GCPDisk describes disks for GCP.
type GCPDisk struct {
	AutoDelete bool              `json:"autoDelete"`
	Boot       bool              `json:"boot"`
	SizeGb     int64             `json:"sizeGb"`
	Type       string            `json:"type"`
	Image      string            `json:"image"`
	Labels     map[string]string `json:"labels"`
}

// GCPMetadata describes metadata for GCP.
type GCPMetadata struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

// GCPNetworkInterface describes network interfaces for GCP
type GCPNetworkInterface struct {
	PublicIP   bool   `json:"publicIP,omitempty"`
	Network    string `json:"network,omitempty"`
	ProjectID  string `json:"projectID,omitempty"`
	Subnetwork string `json:"subnetwork,omitempty"`
}

// GCPServiceAccount describes service accounts for GCP.
type GCPServiceAccount struct {
	Email  string   `json:"email"`
	Scopes []string `json:"scopes"`
}
