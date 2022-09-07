/*
Copyright 2021.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Interface struct {
	// nic name used in the yaml, which relates 1:1 to the mac address.
	// Name in REST API: logicalNICName
	Name string `json:"name"`
	// mac address present on the host.
	// +kubebuilder:validation:Pattern=`^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$`
	MacAddress string `json:"macAddress"`
}

type RawNetConfig []byte

// NetConfig contains the namestatectl yaml [1] as string instead of golang struct
// so we don't need to be in sync with the schema.
//
// [1] https://github.com/nmstate/nmstate/blob/base/libnmstate/schemas/operational-state.yaml
// +kubebuilder:validation:Type=object
type NetConfig struct {
	Raw RawNetConfig `json:"-"`
}

type NMStateConfigSpec struct {
	// Interfaces is an array of interface objects containing the name and MAC
	// address for interfaces that are referenced in the raw nmstate config YAML.
	// Interfaces listed here will be automatically renamed in the nmstate config
	// YAML to match the real device name that is observed to have the
	// corresponding MAC address. At least one interface must be listed so that it
	// can be used to identify the correct host, which is done by matching any MAC
	// address in this list to any MAC address observed on the host.
	// +kubebuilder:validation:MinItems=1
	Interfaces []*Interface `json:"interfaces,omitempty"`
	// yaml that can be processed by nmstate, using custom marshaling/unmarshaling that will allow to populate nmstate config as plain yaml.
	// +kubebuilder:validation:XPreserveUnknownFields
	NetConfig NetConfig `json:"config,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type NMStateConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NMStateConfigSpec `json:"spec,omitempty"`
	// No status
}

// +kubebuilder:object:root=true

// NMStateConfigList contains a list of NMStateConfigs
type NMStateConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NMStateConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NMStateConfig{}, &NMStateConfigList{})
}

// This override the NetConfig type [1] so we can do a custom marshalling of
// nmstate yaml without the need to have golang code representing the nmstate schema

// [1] https://github.com/kubernetes/kube-openapi/tree/master/pkg/generators
func (_ NetConfig) OpenAPISchemaType() []string { return []string{"object"} }
