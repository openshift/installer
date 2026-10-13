/*
Copyright 2024 The ORC Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=openstack
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Available",type="string",JSONPath=".status.conditions[?(@.type=='Available')].status",description="Availability status of resource"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Progressing')].message",description="Message describing current progress status"

// RouterInterface is the Schema for an ORC resource.
type RouterInterface struct {
	metav1.TypeMeta `json:",inline"`

	// metadata contains the object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec specifies the desired state of the resource.
	// +optional
	Spec RouterInterfaceSpec `json:"spec,omitempty"`

	// status defines the observed state of the resource.
	// +optional
	Status RouterInterfaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RouterInterfaceList contains a list of RouterInterface.
type RouterInterfaceList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata contains the list metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// items contains a list of RouterInterface.
	// +kubebuilder:validation:MaxItems:=64
	// +required
	Items []RouterInterface `json:"items"`
}

func (l *RouterInterfaceList) GetItems() []RouterInterface {
	return l.Items
}

// +kubebuilder:validation:Enum:=Subnet
// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=8
type RouterInterfaceType string

const (
	RouterInterfaceTypeSubnet RouterInterfaceType = "Subnet"
)

// +kubebuilder:validation:XValidation:rule="self.type == 'Subnet' ? has(self.subnetRef) : !has(self.subnetRef)",message="subnetRef is required when type is 'Subnet' and not permitted otherwise"
// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="RouterInterfaceResourceSpec is immutable"
type RouterInterfaceSpec struct {
	// type specifies the type of the router interface.
	// +required
	// +unionDiscriminator
	Type RouterInterfaceType `json:"type,omitempty"`

	// routerRef references the router to which this interface belongs.
	// +required
	RouterRef KubernetesNameRef `json:"routerRef,omitempty"`

	// subnetRef references the subnet the router interface is created on.
	// +unionMember
	// +optional
	SubnetRef *KubernetesNameRef `json:"subnetRef,omitempty"`
}

type RouterInterfaceStatus struct {
	// conditions represents the observed status of the object.
	// Known .status.conditions.type are: "Available", "Progressing"
	//
	// Available represents the availability of the OpenStack resource. If it is
	// true then the resource is ready for use.
	//
	// Progressing indicates whether the controller is still attempting to
	// reconcile the current state of the OpenStack resource to the desired
	// state. Progressing will be False either because the desired state has
	// been achieved, or because some terminal error prevents it from ever being
	// achieved and the controller is no longer attempting to reconcile. If
	// Progressing is True, an observer waiting on the resource should continue
	// to wait.
	//
	// +kubebuilder:validation:MaxItems:=32
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// id is the unique identifier of the port created for the router interface
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ID *string `json:"id,omitempty"`
}

var _ ObjectWithConditions = &Router{}

func (i *RouterInterface) GetConditions() []metav1.Condition {
	return i.Status.Conditions
}

func init() {
	SchemeBuilder.Register(&RouterInterface{}, &RouterInterfaceList{})
}
