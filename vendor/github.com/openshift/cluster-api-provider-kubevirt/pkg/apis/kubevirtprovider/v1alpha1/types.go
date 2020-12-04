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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

// KubevirtMachineProviderSpec is the Schema for the KubevirtMachineProviderSpec API
// +k8s:openapi-gen=true
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KubevirtMachineProviderSpec struct {
	metav1.TypeMeta            `json:",inline"`
	SourcePvcName              string `json:"sourcePvcName,omitempty"`
	CredentialsSecretName      string `json:"credentialsSecretName,omitempty"`
	RequestedMemory            string `json:"requestedMemory,omitempty"`
	RequestedCPU               uint32 `json:"requestedCPU,omitempty"`
	RequestedStorage           string `json:"requestedStorage,omitempty"`
	StorageClassName           string `json:"storageClassName,omitempty"`
	IgnitionSecretName         string `json:"ignitionSecretName,omitempty"`
	NetworkName                string `json:"networkName,omitempty"`
	PersistentVolumeAccessMode string `json:"persistentVolumeAccessMode,omitempty"`
}

// KubevirtMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains Kubevirt-specific status information.
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KubevirtMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
	kubevirtapiv1.VirtualMachineStatus
}

func init() {
	SchemeBuilder.Register(&KubevirtMachineProviderSpec{}, &KubevirtMachineProviderStatus{})
}
