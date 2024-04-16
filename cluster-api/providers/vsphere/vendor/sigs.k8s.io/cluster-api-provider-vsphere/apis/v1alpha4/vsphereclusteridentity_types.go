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

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// SecretIdentitySetFinalizer is the finalizer for VSphereCluster credentials secrets .
	SecretIdentitySetFinalizer = "vspherecluster/infrastructure.cluster.x-k8s.io"
)

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

type VSphereClusterIdentityStatus struct {
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Conditions defines current service state of the VSphereCluster.
	// +optional
	Conditions Conditions `json:"conditions,omitempty"`
}

type AllowedNamespaces struct {
	// Selector is a standard Kubernetes LabelSelector. A label query over a set of resources.
	// +optional
	Selector metav1.LabelSelector `json:"selector"`
}

type VSphereIdentityKind string

var (
	VSphereClusterIdentityKind = VSphereIdentityKind("VSphereClusterIdentity")
	SecretKind                 = VSphereIdentityKind("Secret")
)

type VSphereIdentityReference struct {
	// Kind of the identity. Can either be VSphereClusterIdentity or Secret
	// +kubebuilder:validation:Enum=VSphereClusterIdentity;Secret
	Kind VSphereIdentityKind `json:"kind"`

	// Name of the identity.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

func (c *VSphereClusterIdentity) GetConditions() Conditions {
	return c.Status.Conditions
}

func (c *VSphereClusterIdentity) SetConditions(conditions Conditions) {
	c.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion
// +kubebuilder:resource:path=vsphereclusteridentities,scope=Cluster,categories=cluster-api
// +kubebuilder:subresource:status

// VSphereClusterIdentity defines the account to be used for reconciling clusters
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereClusterIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereClusterIdentitySpec   `json:"spec,omitempty"`
	Status VSphereClusterIdentityStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereClusterIdentityList contains a list of VSphereClusterIdentity
//
// Deprecated: This type will be removed in one of the next releases.
type VSphereClusterIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereClusterIdentity `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereClusterIdentity{}, &VSphereClusterIdentityList{})
}
