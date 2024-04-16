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
	// SecretIdentitySetFinalizer is the finalizer for VSphereCluster credentials secrets .
	SecretIdentitySetFinalizer = "vspherecluster/infrastructure.cluster.x-k8s.io"
	// VSphereClusterIdentityFinalizer is the finalizer for VSphereClusterIdentity credentials secrets.
	VSphereClusterIdentityFinalizer = "vsphereclusteridentity/infrastructure.cluster.x-k8s.io"
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
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
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
func (c *VSphereClusterIdentity) GetConditions() clusterv1.Conditions {
	return c.Status.Conditions
}

// SetConditions sets the conditions on the VSphereClusterIdentity.
func (c *VSphereClusterIdentity) SetConditions(conditions clusterv1.Conditions) {
	c.Status.Conditions = conditions
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
