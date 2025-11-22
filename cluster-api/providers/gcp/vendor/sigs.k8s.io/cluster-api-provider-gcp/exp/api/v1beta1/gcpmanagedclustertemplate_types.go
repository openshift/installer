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

// GCPManagedClusterTemplateSpec defines the desired state of GCPManagedClusterTemplate.
type GCPManagedClusterTemplateSpec struct {
	Template GCPManagedClusterTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=gcpmanagedclustertemplates,scope=Namespaced,categories=cluster-api,shortName=amct
// +kubebuilder:storageversion

// GCPManagedClusterTemplate is the Schema for the GCPManagedClusterTemplates API.
type GCPManagedClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GCPManagedClusterTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// GCPManagedClusterTemplateList contains a list of GCPManagedClusterTemplates.
type GCPManagedClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPManagedClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCPManagedClusterTemplate{}, &GCPManagedClusterTemplateList{})
}

// GCPManagedClusterTemplateResource describes the data needed to create an GCPManagedCluster from a template.
type GCPManagedClusterTemplateResource struct {
	Spec GCPManagedClusterTemplateResourceSpec `json:"spec"`
}
