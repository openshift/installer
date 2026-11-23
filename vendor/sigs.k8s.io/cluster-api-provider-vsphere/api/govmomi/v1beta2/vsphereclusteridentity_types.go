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
	// SecretIdentitySetFinalizer is the finalizer for VSphereCluster credentials secrets .
	SecretIdentitySetFinalizer = "vspherecluster/infrastructure.cluster.x-k8s.io"
	// VSphereClusterIdentityFinalizer is the finalizer for VSphereClusterIdentity credentials secrets.
	VSphereClusterIdentityFinalizer = "vsphereclusteridentity/infrastructure.cluster.x-k8s.io"
)

// VSphereClusterIdentity's Available condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterIdentityAvailableCondition documents the availability for a VSphereClusterIdentity.
	VSphereClusterIdentityAvailableCondition = clusterv1.AvailableCondition

	// VSphereClusterIdentityAvailableReason surfaces when the VSphereClusterIdentity is available.
	VSphereClusterIdentityAvailableReason = clusterv1.AvailableReason

	// VSphereClusterIdentitySecretNotAvailableReason surfaces when the VSphereClusterIdentity secret is not available.
	VSphereClusterIdentitySecretNotAvailableReason = "SecretNotAvailable"

	// VSphereClusterIdentitySecretAlreadyInUseReason surfaces when the VSphereClusterIdentity secret is already in use.
	VSphereClusterIdentitySecretAlreadyInUseReason = "SecretAlreadyInUse"

	// VSphereClusterIdentitySettingSecretOwnerReferenceFailedReason surfaces when setting the owner reference on the VSphereClusterIdentity secret failed.
	VSphereClusterIdentitySettingSecretOwnerReferenceFailedReason = "SettingSecretOwnerReferenceFailed"

	// VSphereClusterIdentityDeletingReason surfaces when the VSphereClusterIdentity is being deleted.
	VSphereClusterIdentityDeletingReason = clusterv1.DeletingReason
)

// VSphereClusterIdentitySpec contains a secret reference and a group of allowed namespaces.
type VSphereClusterIdentitySpec struct {
	// secretName references a Secret inside the controller namespace with the credentials to use
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	SecretName string `json:"secretName,omitempty"`

	// allowedNamespaces is used to identify which namespaces are allowed to use this account.
	// Namespaces can be selected with a label selector.
	// If this object is nil, no namespaces will be allowed
	// +optional
	AllowedNamespaces *AllowedNamespaces `json:"allowedNamespaces,omitempty,omitzero"`
}

// VSphereClusterIdentityStatus contains the status of the VSphereClusterIdentity.
// +kubebuilder:validation:MinProperties=1
type VSphereClusterIdentityStatus struct {
	// conditions represents the observations of a VSphereClusterIdentity's current state.
	// Known condition types are Available and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ready is true when the VSphereClusterIdentity is ready.
	// +optional
	Ready *bool `json:"ready,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *VSphereClusterIdentityDeprecatedStatus `json:"deprecated,omitempty"`
}

// VSphereClusterIdentityDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterIdentityDeprecatedStatus struct {
	// v1beta1 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	// +optional
	V1Beta1 *VSphereClusterIdentityV1Beta1DeprecatedStatus `json:"v1beta1,omitempty"`
}

// VSphereClusterIdentityV1Beta1DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterIdentityV1Beta1DeprecatedStatus struct {
	// conditions defines current service state of the VSphereClusterIdentity.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// AllowedNamespaces restricts the namespaces this VSphereClusterIdentity can be used from.
type AllowedNamespaces struct {
	// selector is a standard Kubernetes LabelSelector. A label query over a set of resources.
	// +optional
	Selector metav1.LabelSelector `json:"selector,omitempty,omitzero"`
}

// VSphereIdentityKind is the kind of mechanism used to handle credentials for the VCenter API.
type VSphereIdentityKind string

var (
	// VSphereClusterIdentityKind is used when a VSphereClusterIdentity is referenced in a VSphereCluster.
	VSphereClusterIdentityKind = VSphereIdentityKind("VSphereClusterIdentity")
	// SecretKind is used when a secret is referenced directly in a VSphereCluster.
	SecretKind = VSphereIdentityKind("Secret")
)

// VSphereIdentityReference is the mechanism used to handle credentials for the VCenter API.
type VSphereIdentityReference struct {
	// kind of the identity. Can either be VSphereClusterIdentity or Secret
	// +required
	// +kubebuilder:validation:Enum=VSphereClusterIdentity;Secret
	Kind VSphereIdentityKind `json:"kind,omitempty"`

	// name of the identity.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name,omitempty"`
}

// IsDefined returns true if the ref is defined.
func (m *VSphereIdentityReference) IsDefined() bool {
	return m.Kind != "" || m.Name != ""
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (c *VSphereClusterIdentity) GetV1Beta1Conditions() clusterv1.Conditions {
	if c.Status.Deprecated == nil || c.Status.Deprecated.V1Beta1 == nil {
		return nil
	}
	return c.Status.Deprecated.V1Beta1.Conditions
}

// SetV1Beta1Conditions sets the conditions on this object.
func (c *VSphereClusterIdentity) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if c.Status.Deprecated == nil {
		c.Status.Deprecated = &VSphereClusterIdentityDeprecatedStatus{}
	}
	if c.Status.Deprecated.V1Beta1 == nil {
		c.Status.Deprecated.V1Beta1 = &VSphereClusterIdentityV1Beta1DeprecatedStatus{}
	}
	c.Status.Deprecated.V1Beta1.Conditions = conditions
}

// GetConditions returns the set of conditions for this object.
func (c *VSphereClusterIdentity) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (c *VSphereClusterIdentity) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclusteridentities,scope=Cluster,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VSphereClusterIdentity defines the account to be used for reconciling clusters.
type VSphereClusterIdentity struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereClusterIdentity.
	// +required
	Spec VSphereClusterIdentitySpec `json:"spec,omitempty,omitzero"`

	// status is the observed state of VSphereClusterIdentity.
	// +optional
	Status VSphereClusterIdentityStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// VSphereClusterIdentityList contains a list of VSphereClusterIdentity.
type VSphereClusterIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereClusterIdentity `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereClusterIdentity{}, &VSphereClusterIdentityList{})
}
