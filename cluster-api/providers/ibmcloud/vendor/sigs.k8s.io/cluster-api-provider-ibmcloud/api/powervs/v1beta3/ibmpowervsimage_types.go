/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// IBMPowerVSImageFinalizer allows IBMPowerVSImageReconciler to clean up resources associated with IBMPowerVSImage before
	// removing it from the apiserver.
	IBMPowerVSImageFinalizer = "ibmpowervsimage.infrastructure.cluster.x-k8s.io"
)

// PowerVSStorageType defines the storage tier for the IBM PowerVS instance.
// +kubebuilder:validation:Enum=tier0;tier1;tier3
type PowerVSStorageType string

const (
	// PowerVSStorageTypeTier0 represents tier 0 storage.
	PowerVSStorageTypeTier0 PowerVSStorageType = "tier0"

	// PowerVSStorageTypeTier1 represents tier 1 storage.
	PowerVSStorageTypeTier1 PowerVSStorageType = "tier1"

	// PowerVSStorageTypeTier3 represents tier 3 storage.
	PowerVSStorageTypeTier3 PowerVSStorageType = "tier3"
)

// PowerVSImageDeletePolicy defines the policy for image retention.
// +kubebuilder:validation:Enum=delete;retain
type PowerVSImageDeletePolicy string

const (
	// PowerVSImageDeletePolicyDelete indicates the image will be deleted when the resource is deleted.
	PowerVSImageDeletePolicyDelete PowerVSImageDeletePolicy = "delete"

	// PowerVSImageDeletePolicyRetain indicates the image will be preserved when the resource is deleted.
	PowerVSImageDeletePolicyRetain PowerVSImageDeletePolicy = "retain"
)

func init() {
	objectTypes = append(objectTypes, &IBMPowerVSImage{}, &IBMPowerVSImageList{})
}

// IBMPowerVSImageSpec defines the desired state of IBMPowerVSImage.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSImageSpec struct {
	// clusterName is the name of the Cluster this object belongs to.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	ClusterName string `json:"clusterName,omitempty"`

	// workspace identifies the PowerVS workspace into which the image will be imported.
	// If omitted, the workspace is inherited from the associated IBMPowerVSCluster.
	// +optional
	Workspace ResourceIdentifier `json:"workspace,omitempty,omitzero"`

	// bucket is the Cloud Object Storage bucket name; bucket-name[/optional/folder]
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Bucket string `json:"bucket,omitempty"`

	// object is the Cloud Object Storage image filename.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Object string `json:"object,omitempty"`

	// region is the Cloud Object Storage region.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	Region string `json:"region,omitempty"`

	// storageType is the type of storage, storage pool with the most available space will be selected.
	// +optional
	StorageType PowerVSStorageType `json:"storageType,omitempty"`

	// deletePolicy defines the policy used to identify images to be preserved beyond the lifecycle of associated cluster.
	// +optional
	DeletePolicy PowerVSImageDeletePolicy `json:"deletePolicy,omitempty"`
}

// IBMPowerVSImageStatus defines the observed state of IBMPowerVSImage.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSImageStatus struct {
	// conditions represents the observations of a IBMPowerVSImage's current state.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// imageID is the id of the imported image.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ImageID string `json:"imageID,omitempty"`

	// imageState is the status of the imported image.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	ImageState PowerVSImageState `json:"imageState,omitempty"`

	// jobID is the job ID of an import operation.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	JobID string `json:"jobID,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *IBMPowerVSImageDeprecatedStatus `json:"deprecated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:path=ibmpowervsimages,scope=Namespaced,categories=cluster-api
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.imageState",description="PowerVS image state"

// IBMPowerVSImage is the Schema for the ibmpowervsimages API.
type IBMPowerVSImage struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of IBMPowerVSImage
	// +required
	Spec IBMPowerVSImageSpec `json:"spec,omitempty,omitzero"`

	// status defines the observed state of IBMPowerVSImage
	// +optional
	Status IBMPowerVSImageStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// IBMPowerVSImageList contains a list of IBMPowerVSImage.
type IBMPowerVSImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []IBMPowerVSImage `json:"items"`
}

// IBMPowerVSImageDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type IBMPowerVSImageDeprecatedStatus struct {
	// v1beta2 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	V1Beta2 *IBMPowerVSImageV1Beta2DeprecatedStatus `json:"v1beta2,omitempty"`
}

// IBMPowerVSImageV1Beta2DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type IBMPowerVSImageV1Beta2DeprecatedStatus struct {
	// conditions defines current service state of the VSphereMachine.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// GetConditions returns the observations of the operational state of the IBMPowerVSImage resource.
func (r *IBMPowerVSImage) GetConditions() []metav1.Condition {
	return r.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (r *IBMPowerVSImage) SetConditions(conditions []metav1.Condition) {
	r.Status.Conditions = conditions
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (r *IBMPowerVSImage) GetV1Beta1Conditions() clusterv1.Conditions {
	if r.Status.Deprecated == nil || r.Status.Deprecated.V1Beta2 == nil {
		return nil
	}
	return r.Status.Deprecated.V1Beta2.Conditions
}

// SetV1Beta1Conditions sets conditions for an API object.
func (r *IBMPowerVSImage) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if r.Status.Deprecated == nil {
		r.Status.Deprecated = &IBMPowerVSImageDeprecatedStatus{}
	}
	if r.Status.Deprecated.V1Beta2 == nil {
		r.Status.Deprecated.V1Beta2 = &IBMPowerVSImageV1Beta2DeprecatedStatus{}
	}
	r.Status.Deprecated.V1Beta2.Conditions = conditions
}
