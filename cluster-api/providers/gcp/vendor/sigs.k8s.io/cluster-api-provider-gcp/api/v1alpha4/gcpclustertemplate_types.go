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

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GCPClusterTemplateSpec defines the desired state of GCPClusterTemplate.
type GCPClusterTemplateSpec struct {
	Template GCPClusterTemplateResource `json:"template"`
}

// GCPClusterTemplateResource contains spec for GCPClusterSpec.
type GCPClusterTemplateResource struct {
	Spec GCPClusterSpec `json:"spec"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:path=gcpclustertemplates,scope=Namespaced,categories=cluster-api,shortName=gcpct

// GCPClusterTemplate is the Schema for the gcpclustertemplates API.
type GCPClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GCPClusterTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// GCPClusterTemplateList contains a list of GCPClusterTemplate.
type GCPClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCPClusterTemplate{}, &GCPClusterTemplateList{})
}
