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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GCPManagedControlPlaneTemplateSpec defines the desired state of GCPManagedControlPlaneTemplate.
type GCPManagedControlPlaneTemplateSpec struct {
	Template GCPManagedControlPlaneTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=gcpmanagedcontrolplanetemplates,scope=Namespaced,categories=cluster-api,shortName=amcpt
// +kubebuilder:storageversion

// GCPManagedControlPlaneTemplate is the Schema for the GCPManagedControlPlaneTemplates API.
type GCPManagedControlPlaneTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GCPManagedControlPlaneTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// GCPManagedControlPlaneTemplateList contains a list of GCPManagedControlPlaneTemplates.
type GCPManagedControlPlaneTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPManagedControlPlaneTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCPManagedControlPlaneTemplate{}, &GCPManagedControlPlaneTemplateList{})
}

// GCPManagedControlPlaneTemplateResource describes the data needed to create an GCPManagedCluster from a template.
type GCPManagedControlPlaneTemplateResource struct {
	Spec GCPManagedControlPlaneTemplateResourceSpec `json:"spec"`
}
