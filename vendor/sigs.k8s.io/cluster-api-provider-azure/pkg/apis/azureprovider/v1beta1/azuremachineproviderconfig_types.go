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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Annotation constants
const (
	// ClusterIDLabel is the label that a machineset must have to identify the
	// cluster to which it belongs.
	ClusterIDLabel = "machine.openshift.io/cluster-api-cluster"
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

	Location     string            `json:"location,omitempty"`
	VMSize       string            `json:"vmSize,omitempty"`
	Image        Image             `json:"image"`
	OSDisk       OSDisk            `json:"osDisk"`
	SSHPublicKey string            `json:"sshPublicKey,omitempty"`
	PublicIP     bool              `json:"publicIP"`
	Tags         map[string]string `json:"tags,omitempty"`

	// Network Security Group that needs to be attached to the machine's interface.
	// No security group will be attached if empty.
	SecurityGroup string `json:"securityGroup,omitempty"`

	// Application Security Groups that need to be attached to the machine's interface.
	// No application security groups will be attached if zero-length.
	ApplicationSecurityGroups []string `json:"applicationSecurityGroups,omitempty"`

	// Subnet to use for this instance
	Subnet string `json:"subnet"`

	// PublicLoadBalancer to use for this instance
	PublicLoadBalancer string `json:"publicLoadBalancer,omitempty"`

	// InternalLoadBalancerName to use for this instance
	InternalLoadBalancer string `json:"internalLoadBalancer,omitempty"`

	// NatRule to set inbound NAT rule of the load balancer
	NatRule *int `json:"natRule,omitempty"`

	// ManagedIdentity to set managed identity name
	ManagedIdentity string `json:"managedIdentity,omitempty"`

	// Vnet to set virtual network name
	Vnet string `json:"vnet,omitempty"`

	// Availability Zone for the virtual machine.
	// If nil, the virtual machine should be deployed to no zone
	Zone *string `json:"zone,omitempty"`

	NetworkResourceGroup string `json:"networkResourceGroup,omitempty"`
	ResourceGroup        string `json:"resourceGroup,omitempty"`

	// SpotVMOptions allows the ability to specify the Machine should use a Spot VM
	SpotVMOptions *SpotVMOptions `json:"spotVMOptions,omitempty"`
}

// SpotVMOptions defines the options relevant to running the Machine on Spot VMs
type SpotVMOptions struct {
	// MaxPrice defines the maximum price the user is willing to pay for Spot VM instances
	MaxPrice *string `json:"maxPrice,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&AzureMachineProviderSpec{})
}
