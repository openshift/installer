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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenStackMachineTemplateSpec defines the desired state of OpenStackMachineTemplate.
type OpenStackMachineTemplateSpec struct {
	Template OpenStackMachineTemplateResource `json:"template"`
}

// OpenStackMachineTemplateStatus defines the observed state of OpenStackMachineTemplate.
type OpenStackMachineTemplateStatus struct {
	// Capacity defines the resource capacity for this machine.
	// This value is used for autoscaling from zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`
	// +optional
	NodeInfo NodeInfo `json:"nodeInfo,omitempty,omitzero"`
}

// NodeInfo contains information about the node's architecture and operating system.
// +kubebuilder:validation:MinProperties=1
type NodeInfo struct {
	// operatingSystem is a string representing the operating system of the node.
	// This may be a string like 'linux' or 'windows'.
	// +optional
	OperatingSystem string `json:"operatingSystem,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=openstackmachinetemplates,scope=Namespaced,categories=cluster-api,shortName=osmt
// +kubebuilder:subresource:status

// OpenStackMachineTemplate is the Schema for the openstackmachinetemplates API.
type OpenStackMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackMachineTemplateSpec   `json:"spec,omitempty"`
	Status OpenStackMachineTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackMachineTemplateList contains a list of OpenStackMachineTemplate.
type OpenStackMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackMachineTemplate `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &OpenStackMachineTemplate{}, &OpenStackMachineTemplateList{})
}

// GetIdentifyRef returns the object's namespace and IdentityRef if it has an IdentityRef, or nulls if it does not.
func (r *OpenStackMachineTemplate) GetIdentityRef() (*string, *OpenStackIdentityReference) {
	if r.Spec.Template.Spec.IdentityRef != nil {
		return &r.Namespace, r.Spec.Template.Spec.IdentityRef
	}
	return nil, nil
}
