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
	// SecretIdentitySetFinalizer is the finalizer for VSphereCluster credentials secrets .
	SecretIdentitySetFinalizer = "vspherecluster/infrastructure.cluster.x-k8s.io"
	// VSphereClusterIdentityFinalizer is the finalizer for VSphereClusterIdentity credentials secrets.
	VSphereClusterIdentityFinalizer = "vsphereclusteridentity/infrastructure.cluster.x-k8s.io"
)

// VSphereClusterIdentity's Available condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterIdentityAvailableV1Beta2Condition documents the availability for a VSphereClusterIdentity.
	VSphereClusterIdentityAvailableV1Beta2Condition = clusterv1beta1.AvailableV1Beta2Condition

	// VSphereClusterIdentityAvailableV1Beta2Reason surfaces when the VSphereClusterIdentity is available.
	VSphereClusterIdentityAvailableV1Beta2Reason = clusterv1beta1.AvailableV1Beta2Reason

	// VSphereClusterIdentitySecretNotAvailableV1Beta2Reason surfaces when the VSphereClusterIdentity secret is not available.
	VSphereClusterIdentitySecretNotAvailableV1Beta2Reason = "SecretNotAvailable"

	// VSphereClusterIdentitySecretAlreadyInUseV1Beta2Reason surfaces when the VSphereClusterIdentity secret is already in use.
	VSphereClusterIdentitySecretAlreadyInUseV1Beta2Reason = "SecretAlreadyInUse"

	// VSphereClusterIdentitySettingSecretOwnerReferenceFailedV1Beta2Reason surfaces when setting the owner reference on the VSphereClusterIdentity secret failed.
	VSphereClusterIdentitySettingSecretOwnerReferenceFailedV1Beta2Reason = "SettingSecretOwnerReferenceFailed"

	// VSphereClusterIdentityDeletingV1Beta2Reason surfaces when the VSphereClusterIdentity is being deleted.
	VSphereClusterIdentityDeletingV1Beta2Reason = clusterv1beta1.DeletingV1Beta2Reason
)

// VSphereClusterIdentitySpec contains a secret reference and a group of allowed namespaces.
type VSphereClusterIdentitySpec struct {
	// SecretName references a Secret inside the controller namespace with the credentials to use
	// +kubebuilder:validation:MinLength=1
	SecretName string `json:"secretName,omitempty"`

	// AllowedNamespaces is used to identify which namespaces are allowed to use this account.
	// Namespaces can be selected with a label selector.
	// If this object is nil, no namespaces will be allowed
	// +optional
	AllowedNamespaces *AllowedNamespaces `json:"allowedNamespaces,omitempty"`
}

// VSphereClusterIdentityStatus contains the status of the VSphereClusterIdentity.
type VSphereClusterIdentityStatus struct {
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Conditions defines current service state of the VSphereCluster.
	// +optional
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`

	// v1beta2 groups all the fields that will be added or modified in VSphereClusterIdentity's status with the V1Beta2 version.
	// +optional
	V1Beta2 *VSphereClusterIdentityV1Beta2Status `json:"v1beta2,omitempty"`
}

// VSphereClusterIdentityV1Beta2Status groups all the fields that will be added or modified in VSphereClusterIdentityStatus with the V1Beta2 version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterIdentityV1Beta2Status struct {
	// conditions represents the observations of a VSphereClusterIdentity's current state.
	// Known condition types are Available and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// AllowedNamespaces restricts the namespaces this VSphereClusterIdentity can be used from.
type AllowedNamespaces struct {
	// Selector is a standard Kubernetes LabelSelector. A label query over a set of resources.
	// +optional
	Selector metav1.LabelSelector `json:"selector"`
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
	// Kind of the identity. Can either be VSphereClusterIdentity or Secret
	// +kubebuilder:validation:Enum=VSphereClusterIdentity;Secret
	Kind VSphereIdentityKind `json:"kind"`

	// Name of the identity.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// GetConditions returns the conditions for the VSphereClusterIdentity.
func (c *VSphereClusterIdentity) GetConditions() clusterv1beta1.Conditions {
	return c.Status.Conditions
}

// SetConditions sets the conditions on the VSphereClusterIdentity.
func (c *VSphereClusterIdentity) SetConditions(conditions clusterv1beta1.Conditions) {
	c.Status.Conditions = conditions
}

// GetV1Beta2Conditions returns the set of conditions for this object.
func (c *VSphereClusterIdentity) GetV1Beta2Conditions() []metav1.Condition {
	if c.Status.V1Beta2 == nil {
		return nil
	}
	return c.Status.V1Beta2.Conditions
}

// SetV1Beta2Conditions sets conditions for an API object.
func (c *VSphereClusterIdentity) SetV1Beta2Conditions(conditions []metav1.Condition) {
	if c.Status.V1Beta2 == nil {
		c.Status.V1Beta2 = &VSphereClusterIdentityV1Beta2Status{}
	}
	c.Status.V1Beta2.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclusteridentities,scope=Cluster,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VSphereClusterIdentity defines the account to be used for reconciling clusters.
type VSphereClusterIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereClusterIdentitySpec   `json:"spec,omitempty"`
	Status VSphereClusterIdentityStatus `json:"status,omitempty"`
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
