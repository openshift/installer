package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EquinixMetalMachineProviderConfig is the Shema for the equinixmetalmachineproviderconfigs API.
// +k8s:openapi-gen=true
type EquinixMetalMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is a reference to the secret with EquinixMetal credentials.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`

	MachineType   string   `json:"machineType"`
	Facility      string   `json:"facility"`
	ProjectID     string   `json:"projectID,omitempty"`
	BillingCycle  string   `json:"billingCycle"`
	OS            string   `json:"os"`
	CustomData    string   `json:"customData,omitempty"`
	IPXEScriptURL string   `json:"ipxeScriptURL,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&EquinixMetalMachineProviderConfig{})
}
