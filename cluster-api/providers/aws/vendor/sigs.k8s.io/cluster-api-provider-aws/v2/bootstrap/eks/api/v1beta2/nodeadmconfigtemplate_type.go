/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeadmConfigTemplateSpec defines the desired state of templated NodeadmConfig Amazon EKS Configuration resources.
type NodeadmConfigTemplateSpec struct {
	Template NodeadmConfigTemplateResource `json:"template"`
}

// NodeadmConfigTemplateResource defines the Template structure.
type NodeadmConfigTemplateResource struct {
	// Spec represents the NodeadmConfig each object created from the template will become.
	// We are setting nullable to avoid this issue:
	// https://github.com/kubernetes/kubernetes/issues/117447#issuecomment-2127733969
	// where we cannot remove all fields with an SSA patch if they were previously set.
	// +nullable
	Spec NodeadmConfigSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=nodeadmconfigtemplates,scope=Namespaced,categories=cluster-api,shortName=nodeadmct
// +kubebuilder:storageversion

// NodeadmConfigTemplate is the Amazon EKS Bootstrap Configuration Template API.
type NodeadmConfigTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec NodeadmConfigTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// NodeadmConfigTemplateList contains a list of Amazon EKS Bootstrap Configuration Templates.
type NodeadmConfigTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeadmConfigTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeadmConfigTemplate{}, &NodeadmConfigTemplateList{})
}
