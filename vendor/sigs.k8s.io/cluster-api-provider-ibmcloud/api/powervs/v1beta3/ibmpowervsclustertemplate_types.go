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

package v1beta3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func init() {
	objectTypes = append(objectTypes, &IBMPowerVSClusterTemplate{}, &IBMPowerVSClusterTemplateList{})
}

// IBMPowerVSClusterTemplateSpec defines the desired state of IBMPowerVSClusterTemplate.
type IBMPowerVSClusterTemplateSpec struct {
	// template is the IBMPowerVSClusterTemplateResource.
	// +required
	Template IBMPowerVSClusterTemplateResource `json:"template,omitempty,omitzero"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=ibmpowervsclustertemplates,scope=Namespaced,categories=cluster-api,shortName=ibmpowervsct
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of IBMPowerVSClusterTemplate"

// IBMPowerVSClusterTemplate is the Schema for the ibmpowervsclustertemplates API.
type IBMPowerVSClusterTemplate struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of IBMPowerVSClusterTemplate
	// +required
	Spec IBMPowerVSClusterTemplateSpec `json:"spec,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// IBMPowerVSClusterTemplateList contains a list of IBMPowerVSClusterTemplate.
type IBMPowerVSClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []IBMPowerVSClusterTemplate `json:"items"`
}

// IBMPowerVSClusterTemplateResource describes the data needed to create an IBMPowerVSCluster from a template.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSClusterTemplateResource struct {
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty,omitzero"`
	// spec is the IBMPowerVSClusterSpec.
	// +required
	Spec IBMPowerVSClusterSpec `json:"spec,omitempty,omitzero"`
}
