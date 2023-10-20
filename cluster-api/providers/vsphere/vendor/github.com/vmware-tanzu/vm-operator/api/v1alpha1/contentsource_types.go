// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ContentProviderReference contains the info to locate a content provider resource.
type ContentProviderReference struct {
	// API version of the referent.
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind is the type of resource being referenced.
	Kind string `json:"kind"`
	// Name is the name of resource being referenced.
	Name string `json:"name"`
	// Namespace of the resource being referenced. If empty, cluster scoped resource is assumed.
	Namespace string `json:"namespace,omitempty"`
}

// ContentSourceSpec defines the desired state of ContentSource.
type ContentSourceSpec struct {
	// ProviderRef is a reference to a content provider object that describes a provider.
	ProviderRef ContentProviderReference `json:"providerRef,omitempty"`
}

// ContentSourceStatus defines the observed state of ContentSource.
type ContentSourceStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ContentSource is the Schema for the contentsources API.
// A ContentSource represents the desired specification and the observed status of a ContentSource instance.
type ContentSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContentSourceSpec   `json:"spec,omitempty"`
	Status ContentSourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ContentSourceList contains a list of ContentSource.
type ContentSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContentSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContentSource{}, &ContentSourceList{})
}
