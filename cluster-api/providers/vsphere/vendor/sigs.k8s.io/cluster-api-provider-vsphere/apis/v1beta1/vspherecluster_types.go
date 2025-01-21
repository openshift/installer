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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ClusterFinalizer allows ReconcileVSphereCluster to clean up vSphere
	// resources associated with VSphereCluster before removing it from the
	// API server.
	ClusterFinalizer = "vspherecluster.infrastructure.cluster.x-k8s.io"
)

// VCenterVersion conveys the API version of the vCenter instance.
type VCenterVersion string

// NewVCenterVersion returns a VCenterVersion for the passed string.
func NewVCenterVersion(version string) VCenterVersion {
	return VCenterVersion(version)
}

// VSphereClusterSpec defines the desired state of VSphereCluster.
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

	// ClusterModules hosts information regarding the anti-affinity vSphere constructs
	// for each of the objects responsible for creation of VM objects belonging to the cluster.
	// +optional
	ClusterModules []ClusterModule `json:"clusterModules,omitempty"`

	// DisableClusterModule is used to explicitly turn off the ClusterModule feature.
	// This should work along side NodeAntiAffinity feature flag.
	// If the NodeAntiAffinity feature flag is turned off, this will be disregarded.
	// +optional
	DisableClusterModule bool `json:"disableClusterModule,omitempty"`

	// FailureDomainSelector is the label selector to use for failure domain selection
	// for the control plane nodes of the cluster.
	// If not set (`nil`), selecting failure domains will be disabled.
	// An empty value (`{}`) selects all existing failure domains.
	// A valid selector will select all failure domains which match the selector.
	// +optional
	FailureDomainSelector *metav1.LabelSelector `json:"failureDomainSelector,omitempty"`
}

// ClusterModule holds the anti affinity construct `ClusterModule` identifier
// in use by the VMs owned by the object referred by the TargetObjectName field.
type ClusterModule struct {
	// ControlPlane indicates whether the referred object is responsible for control plane nodes.
	// Currently, only the KubeadmControlPlane objects have this flag set to true.
	// Only a single object in the slice can have this value set to true.
	ControlPlane bool `json:"controlPlane"`

	// TargetObjectName points to the object that uses the Cluster Module information to enforce
	// anti-affinity amongst its descendant VM objects.
	TargetObjectName string `json:"targetObjectName"`

	// ModuleUUID is the unique identifier of the `ClusterModule` used by the object.
	ModuleUUID string `json:"moduleUUID"`
}

// VSphereClusterStatus defines the observed state of VSphereClusterSpec.
type VSphereClusterStatus struct {
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Conditions defines current service state of the VSphereCluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// FailureDomains is a list of failure domain objects synced from the infrastructure provider.
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`

	// VCenterVersion defines the version of the vCenter server defined in the spec.
	VCenterVersion VCenterVersion `json:"vCenterVersion,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for VSphereMachine"
// +kubebuilder:printcolumn:name="Server",type="string",JSONPath=".spec.server",description="Server is the address of the vSphere endpoint."
// +kubebuilder:printcolumn:name="ControlPlaneEndpoint",type="string",JSONPath=".spec.controlPlaneEndpoint[0]",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of Machine"

// VSphereCluster is the Schema for the vsphereclusters API.
type VSphereCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereClusterSpec   `json:"spec,omitempty"`
	Status VSphereClusterStatus `json:"status,omitempty"`
}

// GetConditions returns the conditions for the VSphereCluster.
func (c *VSphereCluster) GetConditions() clusterv1.Conditions {
	return c.Status.Conditions
}

// SetConditions sets conditions on the VSphereCluster.
func (c *VSphereCluster) SetConditions(conditions clusterv1.Conditions) {
	c.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VSphereClusterList contains a list of VSphereCluster.
type VSphereClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereCluster `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereCluster{}, &VSphereClusterList{})
}
