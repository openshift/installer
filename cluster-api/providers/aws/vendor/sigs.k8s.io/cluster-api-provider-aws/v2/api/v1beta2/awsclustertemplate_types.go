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

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// AWSClusterTemplateSpec defines the desired state of AWSClusterTemplate.
type AWSClusterTemplateSpec struct {
	Template AWSClusterTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsclustertemplates,scope=Namespaced,categories=cluster-api,shortName=awsct
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of AWSClusterTemplate"
// +k8s:defaulter-gen=true

// AWSClusterTemplate is the schema for Amazon EC2 based Kubernetes Cluster Templates.
type AWSClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AWSClusterTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// AWSClusterTemplateList contains a list of AWSClusterTemplate.
type AWSClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSClusterTemplate{}, &AWSClusterTemplateList{})
}

type AWSClusterTemplateResource struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty"`
	Spec       AWSClusterSpec       `json:"spec"`
}
