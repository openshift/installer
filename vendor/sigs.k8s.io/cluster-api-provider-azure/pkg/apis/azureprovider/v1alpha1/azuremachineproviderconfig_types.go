/*
Copyright 2018 The Kubernetes Authors.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzureMachineProviderSpec is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an Azure virtual machine. It is used by the Azure machine actuator to create a single Machine.
// Required parameters such as location that are not specified by this configuration, will be defaulted
// by the actuator.
// TODO: Update type
// +k8s:openapi-gen=true
type AzureMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.SecretReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is a reference to the secret with Azure credentials.
	CredentialsSecret *corev1.SecretReference `json:"credentialsSecret,omitempty"`

	Location      string `json:"location"`
	VMSize        string `json:"vmSize"`
	Image         Image  `json:"image"`
	OSDisk        OSDisk `json:"osDisk"`
	SSHPublicKey  string `json:"sshPublicKey"`
	SSHPrivateKey string `json:"sshPrivateKey"`
	PublicIP      bool   `json:"publicIP"`

	// Subnet to use for this instance
	Subnet string `json:"subnet"`

	// PublicLoadBalancer to use for this instance
	PublicLoadBalancer string `json:"publicLoadBalancer"`

	// InternalLoadBalancerName to use for this instance
	InternalLoadBalancer string `json:"internalLoadBalancer"`

	// NatRule to set inbound NAT rule of the load balancer
	NatRule *int `json:"natRule"`

	// ManagedIdentity to set managed identity name
	ManagedIdentity string `json:"managedIdentity"`

	// Vnet to set virtual network name
	Vnet string `json:"vnet"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&AzureMachineProviderSpec{})
}
