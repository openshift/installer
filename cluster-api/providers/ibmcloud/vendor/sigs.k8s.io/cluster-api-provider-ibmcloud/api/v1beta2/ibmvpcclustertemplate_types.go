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

// IBMVPCClusterTemplateSpec defines the desired state of IBMVPCClusterTemplate.
type IBMVPCClusterTemplateSpec struct {
	Template IBMVPCClusterTemplateResource `json:"template,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=ibmvpcclustertemplates,scope=Namespaced,categories=cluster-api,shortName=ibmvpcct
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of IBMVPCClusterTemplate"

// IBMVPCClusterTemplate is the Schema for the ibmvpcclustertemplates API.
type IBMVPCClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec IBMVPCClusterTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// IBMVPCClusterTemplateList contains a list of IBMVPCClusterTemplate.
type IBMVPCClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMVPCClusterTemplate `json:"items"`
}

// IBMVPCClusterTemplateResource describes the data needed to create an IBMVPCCluster from a template.
type IBMVPCClusterTemplateResource struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta capiv1beta1.ObjectMeta `json:"metadata,omitempty"`
	Spec       IBMVPCClusterSpec      `json:"spec"`
}

func init() {
	objectTypes = append(objectTypes, &IBMVPCClusterTemplate{}, &IBMVPCClusterTemplateList{})
}
