/*
Copyright 2021 The Kubernetes Authors.

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

//nolint:godot
package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereDeploymentZoneSpec defines the desired state of VSphereDeploymentZone
type VSphereDeploymentZoneSpec struct {

	// Server is the address of the vSphere endpoint.
	Server string `json:"server,omitempty"`

	// FailureDomain is the name of the VSphereFailureDomain used for this VSphereDeploymentZone
	FailureDomain string `json:"failureDomain,omitempty"`

	// ControlPlane determines if this failure domain is suitable for use by control plane machines.
	// +optional
	ControlPlane *bool `json:"controlPlane,omitempty"`

	// PlacementConstraint encapsulates the placement constraints
	// used within this deployment zone.
	PlacementConstraint PlacementConstraint `json:"placementConstraint"`
}

// PlacementConstraint is the context information for VM placements within a failure domain
type PlacementConstraint struct {
	// ResourcePool is the name or inventory path of the resource pool in which
	// the virtual machine is created/located.
	// +optional
	ResourcePool string `json:"resourcePool,omitempty"`

	// Folder is the name or inventory path of the folder in which the
	// virtual machine is created/located.
	// +optional
	Folder string `json:"folder,omitempty"`
}

type Network struct {
	// Name is the network name for this machine's VM.
	Name string `json:"name,omitempty"`

	// DHCP4 is a flag that indicates whether or not to use DHCP for IPv4
	// +optional
	DHCP4 *bool `json:"dhcp4,omitempty"`

	// DHCP6 indicates whether or not to use DHCP for IPv6
	// +optional
	DHCP6 *bool `json:"dhcp6,omitempty"`
}

type VSphereDeploymentZoneStatus struct {
	// Ready is true when the VSphereDeploymentZone resource is ready.
	// If set to false, it will be ignored by VSphereClusters
	// +optional
	Ready *bool `json:"ready,omitempty"`

	// Conditions defines current service state of the VSphereMachine.
	// +optional
	Conditions Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion
// +kubebuilder:resource:path=vspheredeploymentzones,scope=Cluster,categories=cluster-api
// +kubebuilder:subresource:status

// VSphereDeploymentZone is the Schema for the vspheredeploymentzones API
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereDeploymentZone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereDeploymentZoneSpec   `json:"spec,omitempty"`
	Status VSphereDeploymentZoneStatus `json:"status,omitempty"`
}

func (z *VSphereDeploymentZone) GetConditions() Conditions {
	return z.Status.Conditions
}

func (z *VSphereDeploymentZone) SetConditions(conditions Conditions) {
	z.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VSphereDeploymentZoneList contains a list of VSphereDeploymentZone
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereDeploymentZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereDeploymentZone `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereDeploymentZone{}, &VSphereDeploymentZoneList{})
}
