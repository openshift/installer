/*
Copyright 2025 The Kubernetes Authors.

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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// GKEConfigTemplateSpec defines the desired state of templated GKEConfig GCP GKE Bootstrap Configuration resources.
type GKEConfigTemplateSpec struct {
	Template GKEConfigTemplateResource `json:"template"`
}

// GKEConfigTemplateResource defines the Template structure.
type GKEConfigTemplateResource struct {
	Spec GKEConfigSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=gkeconfigtemplates,scope=Namespaced,categories=cluster-api,shortName=gkect

// GKEConfigTemplate is the GCP GKE Bootstrap Configuration Template API.
type GKEConfigTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GKEConfigTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// GKEConfigTemplateList contains a list of GCP GKE Bootstrap Configuration Templates.
type GKEConfigTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GKEConfigTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GKEConfigTemplate{}, &GKEConfigTemplateList{})
}
