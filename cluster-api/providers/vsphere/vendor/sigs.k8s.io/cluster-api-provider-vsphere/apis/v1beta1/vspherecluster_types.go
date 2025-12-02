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
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

const (
	// ClusterFinalizer allows ReconcileVSphereCluster to clean up vSphere
	// resources associated with VSphereCluster before removing it from the
	// API server.
	ClusterFinalizer = "vspherecluster.infrastructure.cluster.x-k8s.io"
)

// VSphereCluster's Ready condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterReadyV1Beta2Condition is true if the VSphereCluster's deletionTimestamp is not set, VSphereCluster's
	// FailureDomainsReady, VCenterAvailable and ClusterModulesReady conditions are true.
	VSphereClusterReadyV1Beta2Condition = clusterv1beta1.ReadyV1Beta2Condition

	// VSphereClusterReadyV1Beta2Reason surfaces when the VSphereCluster readiness criteria is met.
	VSphereClusterReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterNotReadyV1Beta2Reason surfaces when the VSphereCluster readiness criteria is not met.
	VSphereClusterNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason

	// VSphereClusterReadyUnknownV1Beta2Reason surfaces when at least one VSphereCluster readiness criteria is unknown
	// and no VSphereCluster readiness criteria is not met.
	VSphereClusterReadyUnknownV1Beta2Reason = clusterv1beta1.ReadyUnknownV1Beta2Reason
)

// VSphereCluster's FailureDomainsReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterFailureDomainsReadyV1Beta2Condition documents the status of failure domains for a VSphereCluster.
	VSphereClusterFailureDomainsReadyV1Beta2Condition = "FailureDomainsReady"

	// VSphereClusterFailureDomainsReadyV1Beta2Reason surfaces when the failure domains for a VSphereCluster are ready.
	VSphereClusterFailureDomainsReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterFailureDomainsWaitingForFailureDomainStatusV1Beta2Reason surfaces when not all VSphereFailureDomains for a VSphereCluster are ready.
	VSphereClusterFailureDomainsWaitingForFailureDomainStatusV1Beta2Reason = "WaitingForFailureDomainStatus"

	// VSphereClusterFailureDomainsNotReadyV1Beta2Reason surfaces when the failure domains for a VSphereCluster are not ready.
	VSphereClusterFailureDomainsNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason

	// VSphereClusterFailureDomainsDeletingV1Beta2Reason surfaces when the failure domains for a VSphereCluster are being deleted.
	VSphereClusterFailureDomainsDeletingV1Beta2Reason = clusterv1beta1.DeletingV1Beta2Reason
)

// VSphereCluster's VCenterAvailable condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterVCenterAvailableV1Beta2Condition documents the status of vCenter for a VSphereCluster.
	VSphereClusterVCenterAvailableV1Beta2Condition = "VCenterAvailable"

	// VSphereClusterVCenterAvailableV1Beta2Reason surfaces when the vCenter for a VSphereCluster is available.
	VSphereClusterVCenterAvailableV1Beta2Reason = clusterv1beta1.AvailableV1Beta2Reason

	// VSphereClusterVCenterUnreachableV1Beta2Reason surfaces when the vCenter for a VSphereCluster is unreachable.
	VSphereClusterVCenterUnreachableV1Beta2Reason = "VCenterUnreachable"

	// VSphereClusterVCenterAvailableDeletingV1Beta2Reason surfaces when the VSphereCluster is being deleted.
	VSphereClusterVCenterAvailableDeletingV1Beta2Reason = clusterv1beta1.DeletingV1Beta2Reason
)

// VSphereCluster's ClusterModulesReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterClusterModulesReadyV1Beta2Condition documents the status of vCenter for a VSphereCluster.
	VSphereClusterClusterModulesReadyV1Beta2Condition = "ClusterModulesReady"

	// VSphereClusterClusterModulesReadyV1Beta2Reason surfaces when the cluster modules for a VSphereCluster are ready.
	VSphereClusterClusterModulesReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterModulesInvalidVCenterVersionV1Beta2Reason surfaces when the cluster modules for a VSphereCluster can't be reconciled
	// due to an invalid vCenter version.
	VSphereClusterModulesInvalidVCenterVersionV1Beta2Reason = "InvalidVCenterVersion"

	// VSphereClusterClusterModulesNotReadyV1Beta2Reason surfaces when the cluster modules for a VSphereCluster are not ready.
	VSphereClusterClusterModulesNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason

	// VSphereClusterClusterModulesDeletingV1Beta2Reason surfaces when the cluster modules for a VSphereCluster are being deleted.
	VSphereClusterClusterModulesDeletingV1Beta2Reason = clusterv1beta1.DeletingV1Beta2Reason
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
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`

	// FailureDomains is a list of failure domain objects synced from the infrastructure provider.
	FailureDomains clusterv1beta1.FailureDomains `json:"failureDomains,omitempty"`

	// VCenterVersion defines the version of the vCenter server defined in the spec.
	VCenterVersion VCenterVersion `json:"vCenterVersion,omitempty"`

	// v1beta2 groups all the fields that will be added or modified in VSphereCluster's status with the V1Beta2 version.
	// +optional
	V1Beta2 *VSphereClusterV1Beta2Status `json:"v1beta2,omitempty"`
}

// VSphereClusterV1Beta2Status groups all the fields that will be added or modified in VSphereClusterStatus with the V1Beta2 version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterV1Beta2Status struct {
	// conditions represents the observations of a VSphereCluster's current state.
	// Known condition types are Ready, FailureDomainsReady, VCenterAvailable, ClusterModulesReady and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for VSphereMachine"
// +kubebuilder:printcolumn:name="Server",type="string",JSONPath=".spec.server",description="Server is the address of the vSphere endpoint."
// +kubebuilder:printcolumn:name="ControlPlaneEndpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of Machine"

// VSphereCluster is the Schema for the vsphereclusters API.
type VSphereCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereClusterSpec   `json:"spec,omitempty"`
	Status VSphereClusterStatus `json:"status,omitempty"`
}

// GetConditions returns the conditions for the VSphereCluster.
func (c *VSphereCluster) GetConditions() clusterv1beta1.Conditions {
	return c.Status.Conditions
}

// SetConditions sets conditions on the VSphereCluster.
func (c *VSphereCluster) SetConditions(conditions clusterv1beta1.Conditions) {
	c.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VSphereClusterList contains a list of VSphereCluster.
type VSphereClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereCluster `json:"items"`
}

// GetV1Beta2Conditions returns the set of conditions for this object.
func (c *VSphereCluster) GetV1Beta2Conditions() []metav1.Condition {
	if c.Status.V1Beta2 == nil {
		return nil
	}
	return c.Status.V1Beta2.Conditions
}

// SetV1Beta2Conditions sets conditions for an API object.
func (c *VSphereCluster) SetV1Beta2Conditions(conditions []metav1.Condition) {
	if c.Status.V1Beta2 == nil {
		c.Status.V1Beta2 = &VSphereClusterV1Beta2Status{}
	}
	c.Status.V1Beta2.Conditions = conditions
}

func init() {
	objectTypes = append(objectTypes, &VSphereCluster{}, &VSphereClusterList{})
}
