/*
Copyright 2019 The Kubernetes Authors.

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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/selection"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BareMetalMachineProviderSpec holds data that the actuator needs to provision
// and manage a Machine.
// +k8s:openapi-gen=true
type BareMetalMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Image is the image to be provisioned.
	Image Image `json:"image"`

	// UserData references the Secret that holds user data needed by the bare metal
	// operator. The Namespace is optional; it will default to the Machine's
	// namespace if not specified.
	UserData *corev1.SecretReference `json:"userData,omitempty"`

	// HostSelector specifies matching criteria for labels on BareMetalHosts.
	// This is used to limit the set of BareMetalHost objects considered for
	// claiming for a Machine.
	HostSelector HostSelector `json:"hostSelector,omitempty"`
}

// HostSelector specifies matching criteria for labels on BareMetalHosts.
// This is used to limit the set of BareMetalHost objects considered for
// claiming for a Machine.
type HostSelector struct {
	// Key/value pairs of labels that must exist on a chosen BareMetalHost
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// Label match expressions that must be true on a chosen BareMetalHost
	MatchExpressions []HostSelectorRequirement `json:"matchExpressions,omitempty"`
}

type HostSelectorRequirement struct {
	Key      string             `json:"key"`
	Operator selection.Operator `json:"operator"`
	Values   []string           `json:"values"`
}

// Image holds the details of an image to use during provisioning.
type Image struct {
	// URL is a location of an image to deploy.
	URL string `json:"url"`

	// Checksum is a md5sum value or a URL to retrieve one.
	Checksum string `json:"checksum"`
}

// IsValid returns an error if the object is not valid, otherwise nil. The
// string representation of the error is suitable for human consumption.
func (s *BareMetalMachineProviderSpec) IsValid() error {
	missing := []string{}
	if s.Image.URL == "" {
		missing = append(missing, "Image.URL")
	}
	if s.Image.Checksum == "" {
		missing = append(missing, "Image.Checksum")
	}
	if len(missing) > 0 {
		return fmt.Errorf("Missing fields from ProviderSpec: %v", missing)
	}
	return nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BareMetalMachineProviderSpecList contains a list of BareMetalMachineProviderSpec
type BareMetalMachineProviderSpecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BareMetalMachineProviderSpec `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BareMetalMachineProviderSpec{}, &BareMetalMachineProviderSpecList{})
}
