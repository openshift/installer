/*
Copyright 2021 Nutanix Inc.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NutanixMachineProviderConfig is the Schema for the nutanixmachineproviderconfigs API
// +k8s:openapi-gen=true
type NutanixMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	ClusterReferenceUUID string `json:"clusterReferenceUuid,omitempty"`
	ImageUUID            string `json:"imageUuid,omitempty"`
	ImageName            string `json:"imageName,omitempty"`
	SubnetUUID           string `json:"subnetUuid,omitempty"`
	NumVcpusPerSocket    int64  `json:"numVcpusPerSocket,omitempty"`
	NumSockets           int64  `json:"numSockets,omitempty"`
	MemorySizeMib        int64  `json:"memorySizeMib,omitempty"`
	DiskSizeMib          int64  `json:"diskSizeMib,omitempty"`
	PowerState           string `json:"powerState,omitempty"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret contains a local reference to a secret that contains the
	// credentials data to access Nutanix client
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`
}

func init() {
	SchemeBuilder.Register(&NutanixMachineProviderConfig{}, &NutanixMachineProviderStatus{})
}
