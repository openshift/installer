// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClassReference contains info to locate a Kind VirtualMachineClass object.
type ClassReference struct {
	// API version of the referent.
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the type of resource being referenced.
	Kind string `json:"kind,omitempty"`
	// Name is the name of resource being referenced.
	Name string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmclassbinding
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// VirtualMachineClassBinding is a binding object responsible for
// defining a VirtualMachineClass and a Namespace associated with it.
type VirtualMachineClassBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ClassReference is a reference to a VirtualMachineClass object
	ClassRef ClassReference `json:"classRef,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineClassBindingList contains a list of VirtualMachineClassBinding.
type VirtualMachineClassBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineClassBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachineClassBinding{}, &VirtualMachineClassBindingList{})
}
