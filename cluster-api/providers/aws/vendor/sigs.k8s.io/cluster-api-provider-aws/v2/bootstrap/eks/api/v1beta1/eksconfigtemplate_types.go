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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EKSConfigTemplateSpec defines the desired state of templated EKSConfig Amazon EKS Bootstrap Configuration resources.
type EKSConfigTemplateSpec struct {
	Template EKSConfigTemplateResource `json:"template"`
}

// EKSConfigTemplateResource defines the Template structure.
type EKSConfigTemplateResource struct {
	Spec EKSConfigSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:resource:path=eksconfigtemplates,scope=Namespaced,categories=cluster-api,shortName=eksct

// EKSConfigTemplate is the Amazon EKS Bootstrap Configuration Template API.
type EKSConfigTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EKSConfigTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion

// EKSConfigTemplateList contains a list of Amazon EKS Bootstrap Configuration Templates.
type EKSConfigTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSConfigTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EKSConfigTemplate{}, &EKSConfigTemplateList{})
}
