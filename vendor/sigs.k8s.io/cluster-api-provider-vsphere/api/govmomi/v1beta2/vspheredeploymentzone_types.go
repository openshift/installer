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
	// DeploymentZoneFinalizer allows ReconcileVSphereDeploymentZone to
	// check for dependents associated with VSphereDeploymentZone
	// before removing it from the API Server.
	DeploymentZoneFinalizer = "vspheredeploymentzone.infrastructure.cluster.x-k8s.io"
)

// VSphereDeploymentZone's Ready condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereDeploymentZoneReadyCondition is true if the VSphereDeploymentZone's deletionTimestamp is not set, VSphereDeploymentZone's
	// VCenterAvailable, PlacementConstraintReady and FailureDomainValidated conditions are true.
	VSphereDeploymentZoneReadyCondition = clusterv1.ReadyCondition

	// VSphereDeploymentZoneReadyReason surfaces when the VSphereDeploymentZone readiness criteria is met.
	VSphereDeploymentZoneReadyReason = clusterv1.ReadyReason

	// VSphereDeploymentZoneNotReadyReason surfaces when the VSphereDeploymentZone readiness criteria is not met.
	VSphereDeploymentZoneNotReadyReason = clusterv1.NotReadyReason

	// VSphereDeploymentZoneReadyUnknownReason surfaces when at least one VSphereDeploymentZone readiness criteria is unknown
	// and no VSphereDeploymentZone readiness criteria is not met.
	VSphereDeploymentZoneReadyUnknownReason = clusterv1.ReadyUnknownReason
)

// VSphereDeploymentZone's VCenterAvailable condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereDeploymentZoneVCenterAvailableCondition documents the status of vCenter for a VSphereDeploymentZone.
	VSphereDeploymentZoneVCenterAvailableCondition = "VCenterAvailable"

	// VSphereDeploymentZoneVCenterAvailableReason surfaces when the vCenter for a VSphereDeploymentZone is available.
	VSphereDeploymentZoneVCenterAvailableReason = clusterv1.AvailableReason

	// VSphereDeploymentZoneVCenterUnreachableReason surfaces when the vCenter for a VSphereDeploymentZone is unreachable.
	VSphereDeploymentZoneVCenterUnreachableReason = "VCenterUnreachable"

	// VSphereDeploymentZoneVCenterAvailableDeletingReason surfaces when the vCenter for a VSphereDeploymentZone is being deleted.
	VSphereDeploymentZoneVCenterAvailableDeletingReason = clusterv1.DeletingReason
)

// VSphereDeploymentZone's PlacementConstraintReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereDeploymentZonePlacementConstraintReadyCondition documents the placement constraint status for a VSphereDeploymentZone.
	VSphereDeploymentZonePlacementConstraintReadyCondition = "PlacementConstraintReady"

	// VSphereDeploymentZonePlacementConstraintReadyReason surfaces when the placement status for a VSphereDeploymentZone is ready.
	VSphereDeploymentZonePlacementConstraintReadyReason = clusterv1.ReadyReason

	// VSphereDeploymentZonePlacementConstraintResourcePoolNotFoundReason surfaces when the resource pool for a VSphereDeploymentZone is not found.
	VSphereDeploymentZonePlacementConstraintResourcePoolNotFoundReason = "ResourcePoolNotFound"

	// VSphereDeploymentZonePlacementConstraintFolderNotFoundReason surfaces when the folder for a VSphereDeploymentZone is not found.
	VSphereDeploymentZonePlacementConstraintFolderNotFoundReason = "FolderNotFound"

	// VSphereDeploymentZonePlacementConstraintDeletingReason surfaces when the VSphereDeploymentZone is being deleted.
	VSphereDeploymentZonePlacementConstraintDeletingReason = clusterv1.DeletingReason
)

// VSphereDeploymentZone's FailureDomainValidated condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereDeploymentZoneFailureDomainValidatedCondition documents failure domain validation status for a VSphereDeploymentZone.
	VSphereDeploymentZoneFailureDomainValidatedCondition = "FailureDomainValidated"

	// VSphereDeploymentZoneFailureDomainValidatedReason surfaces when the failure domain for a VSphereDeploymentZone is validated.
	VSphereDeploymentZoneFailureDomainValidatedReason = "Validated"

	// VSphereDeploymentZoneFailureDomainValidationFailedReason surfaces when the failure domain's validation for a VSphereDeploymentZone failed.
	VSphereDeploymentZoneFailureDomainValidationFailedReason = "ValidationFailed"

	// VSphereDeploymentZoneFailureDomainRegionMisconfiguredReason surfaces when the failure domain's region for a VSphereDeploymentZone is misconfigured.
	VSphereDeploymentZoneFailureDomainRegionMisconfiguredReason = "RegionMisconfigured"

	// VSphereDeploymentZoneFailureDomainZoneMisconfiguredReason surfaces when the failure domain's zone for a VSphereDeploymentZone is misconfigured.
	VSphereDeploymentZoneFailureDomainZoneMisconfiguredReason = "ZoneMisconfigured"

	// VSphereDeploymentZoneFailureDomainHostsMisconfiguredReason surfaces when the failure domain's hosts for a VSphereDeploymentZone are misconfigured.
	VSphereDeploymentZoneFailureDomainHostsMisconfiguredReason = "HostsMisconfigured"

	// VSphereDeploymentZoneFailureDomainDatastoreNotFoundReason surfaces when the failure domain's datastore for a VSphereDeploymentZone is not found.
	VSphereDeploymentZoneFailureDomainDatastoreNotFoundReason = "DatastoreNotFound"

	// VSphereDeploymentZoneFailureDomainNetworkNotFoundReason surfaces when the failure domain's network for a VSphereDeploymentZone is not found.
	VSphereDeploymentZoneFailureDomainNetworkNotFoundReason = "NetworkNotFound"

	// VSphereDeploymentZoneFailureDomainComputeClusterNotFoundReason surfaces when the failure domain's compute cluster for a VSphereDeploymentZone is not found.
	VSphereDeploymentZoneFailureDomainComputeClusterNotFoundReason = "ComputeClusterNotFound"

	// VSphereDeploymentZoneFailureDomainResourcePoolNotFoundReason surfaces when the failure domain's resource pool for a VSphereDeploymentZone is not found.
	VSphereDeploymentZoneFailureDomainResourcePoolNotFoundReason = "ResourcePoolNotFound"

	// VSphereDeploymentZoneFailureDomainDeletingReason surfaces when the VSphereDeploymentZone is being deleted.
	VSphereDeploymentZoneFailureDomainDeletingReason = clusterv1.DeletingReason
)

// VSphereDeploymentZoneSpec defines the desired state of VSphereDeploymentZone.
type VSphereDeploymentZoneSpec struct {
	// server is the address of the vSphere endpoint.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Server string `json:"server,omitempty"`

	// failureDomain is the name of the VSphereFailureDomain used for this VSphereDeploymentZone
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	FailureDomain string `json:"failureDomain,omitempty"`

	// controlPlane determines if this failure domain is suitable for use by control plane machines.
	// +optional
	ControlPlane *bool `json:"controlPlane,omitempty"`

	// placementConstraint encapsulates the placement constraints
	// used within this deployment zone.
	// +optional
	PlacementConstraint PlacementConstraint `json:"placementConstraint,omitempty,omitzero"`
}

// PlacementConstraint is the context information for VM placements within a failure domain.
// +kubebuilder:validation:MinProperties=1
type PlacementConstraint struct {
	// resourcePool is the name or inventory path of the resource pool in which
	// the virtual machine is created/located.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	ResourcePool string `json:"resourcePool,omitempty"`

	// folder is the name or inventory path of the folder in which the
	// virtual machine is created/located.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Folder string `json:"folder,omitempty"`
}

// VSphereDeploymentZoneStatus contains the status for a VSphereDeploymentZone.
// +kubebuilder:validation:MinProperties=1
type VSphereDeploymentZoneStatus struct {
	// conditions represents the observations of a VSphereDeploymentZone's current state.
	// Known condition types are Ready, VCenterAvailable, PlacementConstraintReady, FailureDomainValidated and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ready is true when the VSphereDeploymentZone resource is ready.
	// If set to false, it will be ignored by VSphereClusters
	// +optional
	Ready *bool `json:"ready,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *VSphereDeploymentZoneDeprecatedStatus `json:"deprecated,omitempty"`
}

// VSphereDeploymentZoneDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereDeploymentZoneDeprecatedStatus struct {
	// v1beta1 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	// +optional
	V1Beta1 *VSphereDeploymentZoneV1Beta1DeprecatedStatus `json:"v1beta1,omitempty"`
}

// VSphereDeploymentZoneV1Beta1DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereDeploymentZoneV1Beta1DeprecatedStatus struct {
	// conditions defines current service state of the VSphereDeploymentZone.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=vspheredeploymentzones,scope=Cluster,categories=cluster-api
// +kubebuilder:subresource:status

// VSphereDeploymentZone is the Schema for the vspheredeploymentzones API.
type VSphereDeploymentZone struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereDeploymentZone.
	// +required
	Spec VSphereDeploymentZoneSpec `json:"spec,omitempty,omitzero"`

	// status is the observed state of VSphereDeploymentZone.
	// +optional
	Status VSphereDeploymentZoneStatus `json:"status,omitempty,omitzero"`
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (c *VSphereDeploymentZone) GetV1Beta1Conditions() clusterv1.Conditions {
	if c.Status.Deprecated == nil || c.Status.Deprecated.V1Beta1 == nil {
		return nil
	}
	return c.Status.Deprecated.V1Beta1.Conditions
}

// SetV1Beta1Conditions sets the conditions on this object.
func (c *VSphereDeploymentZone) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if c.Status.Deprecated == nil {
		c.Status.Deprecated = &VSphereDeploymentZoneDeprecatedStatus{}
	}
	if c.Status.Deprecated.V1Beta1 == nil {
		c.Status.Deprecated.V1Beta1 = &VSphereDeploymentZoneV1Beta1DeprecatedStatus{}
	}
	c.Status.Deprecated.V1Beta1.Conditions = conditions
}

// GetConditions returns the set of conditions for this object.
func (c *VSphereDeploymentZone) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (c *VSphereDeploymentZone) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VSphereDeploymentZoneList contains a list of VSphereDeploymentZone.
type VSphereDeploymentZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereDeploymentZone `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereDeploymentZone{}, &VSphereDeploymentZoneList{})
}
