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

const (
	// ClusterFinalizer allows ReconcileVSphereCluster to clean up vSphere
	// resources associated with VSphereCluster before removing it from the
	// API server.
	ClusterFinalizer = "vspherecluster.infrastructure.cluster.x-k8s.io"
)

// VSphereClusterSpec defines the desired state of VSphereCluster
type VSphereClusterSpec struct {
	// Server is the address of the vSphere endpoint.
	Server string `json:"server,omitempty"`

	// Thumbprint is the colon-separated SHA-1 checksum of the given vCenter server's host certificate
	// +optional
	Thumbprint string `json:"thumbprint,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint"`

	// IdentityRef is a reference to either a Secret or VSphereClusterIdentity that contains
	// the identity to use when reconciling the cluster.
	// +optional
	IdentityRef *VSphereIdentityReference `json:"identityRef,omitempty"`
}

// VSphereClusterStatus defines the observed state of VSphereClusterSpec
type VSphereClusterStatus struct {
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Conditions defines current service state of the VSphereCluster.
	// +optional
	Conditions Conditions `json:"conditions,omitempty"`

	// FailureDomains is a list of failure domain objects synced from the infrastructure provider.
	FailureDomains FailureDomains `json:"failureDomains,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion
// +kubebuilder:resource:path=vsphereclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for VSphereMachine"
// +kubebuilder:printcolumn:name="Server",type="string",JSONPath=".spec.server",description="Server is the address of the vSphere endpoint"
// +kubebuilder:printcolumn:name="ControlPlaneEndpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of Machine"

// VSphereCluster is the Schema for the vsphereclusters API
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereClusterSpec   `json:"spec,omitempty"`
	Status VSphereClusterStatus `json:"status,omitempty"`
}

// GetConditions returns the conditions for the VSphereCluster.
func (c *VSphereCluster) GetConditions() Conditions {
	return c.Status.Conditions
}

// GetConditions sets conditions on the VSphereCluster.
func (c *VSphereCluster) SetConditions(conditions Conditions) {
	c.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VSphereClusterList contains a list of VSphereCluster
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereCluster `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereCluster{}, &VSphereClusterList{})
}
