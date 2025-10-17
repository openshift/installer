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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// AzureASOManagedClusterKind is the kind for AzureASOManagedCluster.
	AzureASOManagedClusterKind = "AzureASOManagedCluster"

	// AzureASOManagedControlPlaneFinalizer is the finalizer added to AzureASOManagedControlPlanes.
	AzureASOManagedControlPlaneFinalizer = "azureasomanagedcontrolplane.infrastructure.cluster.x-k8s.io"
)

// AzureASOManagedClusterSpec defines the desired state of AzureASOManagedCluster.
type AzureASOManagedClusterSpec struct {
	AzureASOManagedClusterTemplateResourceSpec `json:",inline"`

	// ControlPlaneEndpoint is the location of the API server within the control plane. CAPZ manages this field
	// and it should not be set by the user. It fulfills Cluster API's cluster infrastructure provider contract.
	// Because this field is programmatically set by CAPZ after resource creation, we define it as +optional
	// in the API schema to permit resource admission.
	//+optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// AzureASOManagedClusterStatus defines the observed state of AzureASOManagedCluster.
type AzureASOManagedClusterStatus struct {
	// Ready represents whether or not the cluster has been provisioned and is ready. It fulfills Cluster
	// API's cluster infrastructure provider contract.
	//+optional
	Ready bool `json:"ready"`

	//+optional
	Resources []ResourceStatus `json:"resources,omitempty"`
}

// ResourceStatus represents the status of a resource.
type ResourceStatus struct {
	Resource StatusResource `json:"resource"`
	Ready    bool           `json:"ready"`
}

// StatusResource is a handle to a resource.
type StatusResource struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
	Name    string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// AzureASOManagedCluster is the Schema for the azureasomanagedclusters API.
type AzureASOManagedCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureASOManagedClusterSpec   `json:"spec,omitempty"`
	Status AzureASOManagedClusterStatus `json:"status,omitempty"`
}

// GetResourceStatuses returns the status of resources.
func (a *AzureASOManagedCluster) GetResourceStatuses() []ResourceStatus {
	return a.Status.Resources
}

// SetResourceStatuses sets the status of resources.
func (a *AzureASOManagedCluster) SetResourceStatuses(r []ResourceStatus) {
	a.Status.Resources = r
}

//+kubebuilder:object:root=true

// AzureASOManagedClusterList contains a list of AzureASOManagedCluster.
type AzureASOManagedClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureASOManagedCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureASOManagedCluster{}, &AzureASOManagedClusterList{})
}
