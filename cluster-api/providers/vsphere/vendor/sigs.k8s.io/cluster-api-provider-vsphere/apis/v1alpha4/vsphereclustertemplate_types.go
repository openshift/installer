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

//nolint:godot
package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereClusterTemplateSpec defines the desired state of VSphereClusterTemplate
type VSphereClusterTemplateSpec struct {
	Template VSphereClusterTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion
// +kubebuilder:resource:path=vsphereclustertemplates,scope=Namespaced,categories=cluster-api

// VSphereClusterTemplate is the Schema for the vsphereclustertemplates API
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VSphereClusterTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereClusterTemplateList contains a list of VSphereClusterTemplate.
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereClusterTemplate `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereClusterTemplate{}, &VSphereClusterTemplateList{})
}

type VSphereClusterTemplateResource struct {
	Spec VSphereClusterSpec `json:"spec"`
}
