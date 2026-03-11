/*
Copyright 2022 The Kubernetes Authors.

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

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// IBMPowerVSClusterTemplateSpec defines the desired state of IBMPowerVSClusterTemplate.
type IBMPowerVSClusterTemplateSpec struct {
	Template IBMPowerVSClusterTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
//+kubebuilder:storageversion
// +kubebuilder:resource:path=ibmpowervsclustertemplates,scope=Namespaced,categories=cluster-api,shortName=ibmpowervsct
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of IBMPowerVSClusterTemplate"

// IBMPowerVSClusterTemplate is the schema for IBM Power VS Kubernetes Cluster Templates.
type IBMPowerVSClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IBMPowerVSClusterTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// IBMPowerVSClusterTemplateList contains a list of IBMPowerVSClusterTemplate.
type IBMPowerVSClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMPowerVSClusterTemplate `json:"items"`
}

// IBMPowerVSClusterTemplateResource describes the data needed to create an IBMPowerVSCluster from a template.
type IBMPowerVSClusterTemplateResource struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta capiv1beta1.ObjectMeta `json:"metadata,omitempty"`
	Spec       IBMPowerVSClusterSpec  `json:"spec"`
}

func init() {
	objectTypes = append(objectTypes, &IBMPowerVSClusterTemplate{}, &IBMPowerVSClusterTemplateList{})
}
