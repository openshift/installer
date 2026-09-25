/*
Copyright 2025 The Kubernetes Authors.

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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// ClusterFinalizer allows ReconcileVSphereCluster to clean up vSphere
	// resources associated with VSphereCluster before removing it from the
	// API server.
	ClusterFinalizer = "vspherecluster.infrastructure.cluster.x-k8s.io"
)

// VSphereCluster's Ready condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterReadyCondition is true if the VSphereCluster's deletionTimestamp is not set, VSphereCluster's
	// FailureDomainsReady, VCenterAvailable and ClusterModulesReady conditions are true.
	VSphereClusterReadyCondition = clusterv1.ReadyCondition

	// VSphereClusterReadyReason surfaces when the VSphereCluster readiness criteria is met.
	VSphereClusterReadyReason = clusterv1.ReadyReason

	// VSphereClusterNotReadyReason surfaces when the VSphereCluster readiness criteria is not met.
	VSphereClusterNotReadyReason = clusterv1.NotReadyReason

	// VSphereClusterReadyUnknownReason surfaces when at least one VSphereCluster readiness criteria is unknown
	// and no VSphereCluster readiness criteria is not met.
	VSphereClusterReadyUnknownReason = clusterv1.ReadyUnknownReason
)

// VSphereCluster's FailureDomainsReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterFailureDomainsReadyCondition documents the status of failure domains for a VSphereCluster.
	VSphereClusterFailureDomainsReadyCondition = "FailureDomainsReady"

	// VSphereClusterFailureDomainsReadyReason surfaces when the failure domains for a VSphereCluster are ready.
	VSphereClusterFailureDomainsReadyReason = clusterv1.ReadyReason

	// VSphereClusterFailureDomainsWaitingForFailureDomainStatusReason surfaces when not all VSphereFailureDomains for a VSphereCluster are ready.
	VSphereClusterFailureDomainsWaitingForFailureDomainStatusReason = "WaitingForFailureDomainStatus"

	// VSphereClusterFailureDomainsNotReadyReason surfaces when the failure domains for a VSphereCluster are not ready.
	VSphereClusterFailureDomainsNotReadyReason = clusterv1.NotReadyReason

	// VSphereClusterFailureDomainsDeletingReason surfaces when the failure domains for a VSphereCluster are being deleted.
	VSphereClusterFailureDomainsDeletingReason = clusterv1.DeletingReason
)

// VSphereCluster's VCenterAvailable condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterVCenterAvailableCondition documents the status of vCenter for a VSphereCluster.
	VSphereClusterVCenterAvailableCondition = "VCenterAvailable"

	// VSphereClusterVCenterAvailableReason surfaces when the vCenter for a VSphereCluster is available.
	VSphereClusterVCenterAvailableReason = clusterv1.AvailableReason

	// VSphereClusterVCenterUnreachableReason surfaces when the vCenter for a VSphereCluster is unreachable.
	VSphereClusterVCenterUnreachableReason = "VCenterUnreachable"

	// VSphereClusterVCenterAvailableDeletingReason surfaces when the VSphereCluster is being deleted.
	VSphereClusterVCenterAvailableDeletingReason = clusterv1.DeletingReason
)

// VSphereCluster's ClusterModulesReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterClusterModulesReadyCondition documents the status of vCenter for a VSphereCluster.
	VSphereClusterClusterModulesReadyCondition = "ClusterModulesReady"

	// VSphereClusterClusterModulesReadyReason surfaces when the cluster modules for a VSphereCluster are ready.
	VSphereClusterClusterModulesReadyReason = clusterv1.ReadyReason

	// VSphereClusterModulesInvalidVCenterVersionReason surfaces when the cluster modules for a VSphereCluster can't be reconciled
	// due to an invalid vCenter version.
	VSphereClusterModulesInvalidVCenterVersionReason = "InvalidVCenterVersion"

	// VSphereClusterClusterModulesNotReadyReason surfaces when the cluster modules for a VSphereCluster are not ready.
	VSphereClusterClusterModulesNotReadyReason = clusterv1.NotReadyReason

	// VSphereClusterClusterModulesDeletingReason surfaces when the cluster modules for a VSphereCluster are being deleted.
	VSphereClusterClusterModulesDeletingReason = clusterv1.DeletingReason
)

// VCenterVersion conveys the API version of the vCenter instance.
type VCenterVersion string

// NewVCenterVersion returns a VCenterVersion for the passed string.
func NewVCenterVersion(version string) VCenterVersion {
	return VCenterVersion(version)
}

// VSphereClusterSpec defines the desired state of VSphereCluster.
type VSphereClusterSpec struct {
	// server is the address of the vSphere endpoint.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Server string `json:"server,omitempty"`

	// thumbprint is the colon-separated SHA-1 checksum of the given vCenter server's host certificate
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Thumbprint string `json:"thumbprint,omitempty"`

	// controlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint,omitempty,omitzero"`

	// identityRef is a reference to either a Secret or VSphereClusterIdentity that contains
	// the identity to use when reconciling the cluster.
	// +optional
	IdentityRef VSphereIdentityReference `json:"identityRef,omitempty,omitzero"`

	// clusterModules hosts information regarding the anti-affinity vSphere constructs
	// for each of the objects responsible for creation of VM objects belonging to the cluster.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=512
	ClusterModules []ClusterModule `json:"clusterModules,omitempty"`

	// disableClusterModule is used to explicitly turn off the ClusterModule feature.
	// This should work along side NodeAntiAffinity feature flag.
	// If the NodeAntiAffinity feature flag is turned off, this will be disregarded.
	// +optional
	DisableClusterModule *bool `json:"disableClusterModule,omitempty"`

	// failureDomainSelector is the label selector to use for failure domain selection
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
	// controlPlane indicates whether the referred object is responsible for control plane nodes.
	// Currently, only the KubeadmControlPlane objects have this flag set to true.
	// Only a single object in the slice can have this value set to true.
	// +required
	ControlPlane *bool `json:"controlPlane,omitempty"`

	// targetObjectName points to the object that uses the Cluster Module information to enforce
	// anti-affinity amongst its descendant VM objects.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	TargetObjectName string `json:"targetObjectName,omitempty"`

	// moduleUUID is the unique identifier of the `ClusterModule` used by the object.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	ModuleUUID string `json:"moduleUUID,omitempty"`
}

// VSphereClusterStatus defines the observed state of VSphereClusterSpec.
// +kubebuilder:validation:MinProperties=1
type VSphereClusterStatus struct {
	// conditions represents the observations of a VSphereCluster's current state.
	// Known condition types are Ready, FailureDomainsReady, VCenterAvailable, ClusterModulesReady and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// initialization provides observations of the VSphereCluster initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial Cluster provisioning.
	// +optional
	Initialization VSphereClusterInitializationStatus `json:"initialization,omitempty,omitzero"`

	// failureDomains is a list of failure domain objects synced from the infrastructure provider.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	FailureDomains []clusterv1.FailureDomain `json:"failureDomains,omitempty"`

	// vCenterVersion defines the version of the vCenter server defined in the spec.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	VCenterVersion VCenterVersion `json:"vCenterVersion,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *VSphereClusterDeprecatedStatus `json:"deprecated,omitempty"`
}

// VSphereClusterInitializationStatus provides observations of the VSphereCluster initialization process.
// +kubebuilder:validation:MinProperties=1
type VSphereClusterInitializationStatus struct {
	// provisioned is true when the infrastructure provider reports that the Cluster's infrastructure is fully provisioned.
	// NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Cluster provisioning.
	// +optional
	Provisioned *bool `json:"provisioned,omitempty"`
}

// VSphereClusterDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterDeprecatedStatus struct {
	// v1beta1 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	// +optional
	V1Beta1 *VSphereClusterV1Beta1DeprecatedStatus `json:"v1beta1,omitempty"`
}

// VSphereClusterV1Beta1DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterV1Beta1DeprecatedStatus struct {
	// conditions defines current service state of the VSphereCluster.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
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
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels['cluster\\.x-k8s\\.io/cluster-name']",description="Cluster"
// +kubebuilder:printcolumn:name="Server",type="string",JSONPath=".spec.server",description="Server is the address of the vSphere endpoint."
// +kubebuilder:printcolumn:name="ControlPlane Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint"
// +kubebuilder:printcolumn:name="Paused",type="string",JSONPath=`.status.conditions[?(@.type=="Paused")].status`,description="Reconciliation paused",priority=10
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.initialization.provisioned",description="VSphereCluster is provisioned"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of VSphereCluster"

// VSphereCluster is the Schema for the vsphereclusters API.
type VSphereCluster struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereCluster.
	// +required
	Spec VSphereClusterSpec `json:"spec,omitempty,omitzero"`

	// status is the observed state of VSphereCluster.
	// +optional
	Status VSphereClusterStatus `json:"status,omitempty,omitzero"`
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (c *VSphereCluster) GetV1Beta1Conditions() clusterv1.Conditions {
	if c.Status.Deprecated == nil || c.Status.Deprecated.V1Beta1 == nil {
		return nil
	}
	return c.Status.Deprecated.V1Beta1.Conditions
}

// SetV1Beta1Conditions sets the conditions on this object.
func (c *VSphereCluster) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if c.Status.Deprecated == nil {
		c.Status.Deprecated = &VSphereClusterDeprecatedStatus{}
	}
	if c.Status.Deprecated.V1Beta1 == nil {
		c.Status.Deprecated.V1Beta1 = &VSphereClusterV1Beta1DeprecatedStatus{}
	}
	c.Status.Deprecated.V1Beta1.Conditions = conditions
}

// GetConditions returns the set of conditions for this object.
func (c *VSphereCluster) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (c *VSphereCluster) SetConditions(conditions []metav1.Condition) {
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
