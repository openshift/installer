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

// OvirtMachineProviderSpec is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an Ovirt VM. It is used by the Ovirt machine actuator to create a single machine instance.
type OvirtMachineProviderSpec struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// CredentialsSecret is a reference to the secret with oVirt credentials.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`

	Id string `json:"id"`
	Name string `json:"name"`
	// The VM template this instance will be created from
	TemplateId string `json:"template_id"`

	// the oVirt cluster this VM instance belongs too
	ClusterId string `json:"cluster_id"`


}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OvirtClusterProviderSpec of an oVirt cluster
// +k8s:openapi-gen=true
type OvirtClusterProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OvirtClusterProviderStatus
// +k8s:openapi-gen=true
type OvirtClusterProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// CACertificate is a PEM encoded CA Certificate for the control plane nodes.
	CACertificate []byte

	// CAPrivateKey is a PEM encoded PKCS1 CA PrivateKey for the control plane nodes.
	CAPrivateKey []byte

	// Network contains all information about the created OpenStack Network.
	// It includes Subnets and Router.
}

func init() {
	SchemeBuilder.Register(&OvirtMachineProviderSpec{})
	SchemeBuilder.Register(&OvirtClusterProviderSpec{})
	SchemeBuilder.Register(&OvirtClusterProviderStatus{})
}
