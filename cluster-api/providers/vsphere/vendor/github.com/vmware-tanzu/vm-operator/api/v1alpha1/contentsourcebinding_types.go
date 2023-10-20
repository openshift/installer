// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ContentSourceReference contains info to locate a Kind ContentSource object.
type ContentSourceReference struct {
	// API version of the referent.
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the type of resource being referenced.
	Kind string `json:"kind,omitempty"`
	// Name is the name of resource being referenced.
	Name string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced

// ContentSourceBinding is an object that represents a ContentSource to Namespace mapping.
type ContentSourceBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ContentSourceRef is a reference to a ContentSource object.
	ContentSourceRef ContentSourceReference `json:"contentSourceRef,omitempty"`
}

// +kubebuilder:object:root=true

// ContentSourceBindingList contains a list of ContentSourceBinding.
type ContentSourceBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContentSourceBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContentSourceBinding{}, &ContentSourceBindingList{})
}
