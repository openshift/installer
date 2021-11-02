/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PowerVSMachineProviderConfig is the Schema for the powervsmachineproviderconfigs API
type PowerVSMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ServiceInstanceID is the PowerVS service ID
	ServiceInstanceID string `json:"serviceInstanceID"`

	// Image is the reference to the Image from which to create the machine instance.
	Image PowerVSResourceReference `json:"image"`

	// UserDataSecret is the k8s secret contains the user data script
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is the k8s secret contains the API Key for IBM Cloud authentication
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`

	// SysType is the System type used to host the vsi
	SysType string `json:"sysType"`

	// ProcessorType is the processor type, e.g: dedicated, shared, capped
	ProcType string `json:"procType"`

	// Processors is Number of processors allocated
	Processors string `json:"processors"`

	// Memory is Amount of memory allocated (in GB)
	Memory string `json:"memory"`

	// Network is the reference to the Network to use for this instance.
	Network PowerVSResourceReference `json:"network"`

	// KeyPairName is the name of the SSH key pair provided to the server for authenticating users
	KeyPairName string `json:"keyPairName,omitempty"`
}

// PowerVSResourceReference is a reference to a specific PowerVS resource by ID or Name
// Only one of ID or Name may be specified. Specifying more than one will result in
// a validation error.
type PowerVSResourceReference struct {
	// ID of resource
	// +optional
	ID *string `json:"id,omitempty"`

	// Name of resource
	// +optional
	Name *string `json:"name,omitempty"`
}

//+kubebuilder:object:root=true

// PowerVSMachineProviderConfigList contains a list of PowerVSMachineProviderConfig
type PowerVSMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PowerVSMachineProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PowerVSMachineProviderConfig{}, &PowerVSMachineProviderConfigList{}, &PowerVSMachineProviderStatus{})
}
