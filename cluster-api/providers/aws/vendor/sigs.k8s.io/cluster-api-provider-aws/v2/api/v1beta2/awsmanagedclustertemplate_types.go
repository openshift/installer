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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSManagedClusterTemplateSpec defines the desired state of AWSManagedClusterTemplate.
type AWSManagedClusterTemplateSpec struct {
	Template AWSManagedClusterTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmanagedclustertemplates,scope=Namespaced,categories=cluster-api,shortName=amct
// +kubebuilder:storageversion

// AWSManagedClusterTemplate is the Schema for the AWSManagedClusterTemplates API.
type AWSManagedClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AWSManagedClusterTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AWSManagedClusterTemplateList contains a list of AWSManagedClusterTemplates.
type AWSManagedClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSManagedClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSManagedClusterTemplate{}, &AWSManagedClusterTemplateList{})
}

// AWSManagedClusterTemplateResource describes the data needed to create an AWSManagedCluster from a template.
type AWSManagedClusterTemplateResource struct {
	Spec AWSManagedClusterSpec `json:"spec"`
}
