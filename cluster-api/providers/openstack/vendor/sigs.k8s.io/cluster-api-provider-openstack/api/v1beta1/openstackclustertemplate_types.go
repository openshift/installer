/*
Copyright 2023 The Kubernetes Authors.

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

// OpenStackClusterTemplateResource describes the data needed to create a OpenStackCluster from a template.
type OpenStackClusterTemplateResource struct {
	Spec OpenStackClusterSpec `json:"spec"`
}

// OpenStackClusterTemplateSpec defines the desired state of OpenStackClusterTemplate.
type OpenStackClusterTemplateSpec struct {
	Template OpenStackClusterTemplateResource `json:"template"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=openstackclustertemplates,scope=Namespaced,categories=cluster-api,shortName=osct

// OpenStackClusterTemplate is the Schema for the openstackclustertemplates API.
type OpenStackClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OpenStackClusterTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// OpenStackClusterTemplateList contains a list of OpenStackClusterTemplate.
type OpenStackClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackClusterTemplate `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &OpenStackClusterTemplate{}, &OpenStackClusterTemplateList{})
}
