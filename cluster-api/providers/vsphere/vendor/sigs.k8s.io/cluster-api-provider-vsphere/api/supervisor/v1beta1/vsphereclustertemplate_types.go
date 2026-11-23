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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereClusterTemplateSpec defines the desired state of VSphereClusterTemplate.
type VSphereClusterTemplateSpec struct {
	// template defines the desired state of VSphereClusterTemplate.
	Template VSphereClusterTemplateResource `json:"template"`
}

//+kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclustertemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:deprecatedversion

// VSphereClusterTemplate is the Schema for the vsphereclustertemplates API.
type VSphereClusterTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereClusterTemplate.
	Spec VSphereClusterTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// VSphereClusterTemplateList contains a list of VSphereClusterTemplate.
type VSphereClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereClusterTemplate `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereClusterTemplate{}, &VSphereClusterTemplateList{})
}

// VSphereClusterTemplateResource defines the template structure.
type VSphereClusterTemplateResource struct {
	// spec is the desired state of VSphereClusterTemplateResource.
	Spec VSphereClusterSpec `json:"spec"`
}
