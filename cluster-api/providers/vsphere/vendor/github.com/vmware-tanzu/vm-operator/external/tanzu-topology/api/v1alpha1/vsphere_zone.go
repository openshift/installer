// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereZoneSpec defines the desired state of VSphereZone.
type VSphereZoneSpec struct {
	// Description is the description of the vSphere Zone.
	Description string `json:"description,omitempty"`
}

// VSphereZoneStatus defines the observed state of VSphereZone.
type VSphereZoneStatus struct {
}

// VSphereZone is the schema for the VSphereZone resource for the
// vSphere Zone.
//
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspherezones,scope=Cluster
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type VSphereZone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereZoneSpec   `json:"spec,omitempty"`
	Status VSphereZoneStatus `json:"status,omitempty"`
}

// VSphereZoneList contains a list of VSphereZone resources.
//
// +kubebuilder:object:root=true
type VSphereZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereZone `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereZone{}, &VSphereZoneList{})
}
