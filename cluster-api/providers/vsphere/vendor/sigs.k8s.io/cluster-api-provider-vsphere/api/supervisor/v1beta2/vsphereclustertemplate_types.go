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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// VSphereClusterTemplateSpec defines the desired state of VSphereClusterTemplate.
type VSphereClusterTemplateSpec struct {
	// template defines the desired state of VSphereClusterTemplate.
	// +required
	Template VSphereClusterTemplateResource `json:"template,omitempty,omitzero"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclustertemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="ClusterClass",type="string",JSONPath=`.metadata.ownerReferences[?(@.kind=="ClusterClass")].name`,description="Name of the ClusterClass owning this template"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of VSphereClusterTemplate"

// VSphereClusterTemplate is the Schema for the vsphereclustertemplates API.
type VSphereClusterTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereClusterTemplate.
	// +optional
	Spec VSphereClusterTemplateSpec `json:"spec,omitempty,omitzero"`
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

// VSphereClusterTemplateResource describes the data for creating a VSphereCluster from a template.
// +kubebuilder:validation:MinProperties=1
type VSphereClusterTemplateResource struct {
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec is the desired state of VSphereClusterTemplateResource.
	// +optional
	Spec VSphereClusterSpec `json:"spec,omitempty,omitzero"`
}
