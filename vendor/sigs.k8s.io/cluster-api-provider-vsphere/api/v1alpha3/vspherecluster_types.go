/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/cluster-api-provider-vsphere/api/v1alpha3/cloudprovider"
)

const (
	// ClusterFinalizer allows ReconcileVSphereCluster to clean up vSphere
	// resources associated with VSphereCluster before removing it from the
	// API server.
	ClusterFinalizer = "vspherecluster.infrastructure.cluster.x-k8s.io"
)

// VSphereClusterSpec defines the desired state of VSphereCluster
type VSphereClusterSpec struct {
	// Server is the address of the vSphere endpoint.
	Server string `json:"server,omitempty"`

	// Insecure is a flag that controls whether or not to validate the
	// vSphere server's certificate.
	// +optional
	Insecure *bool `json:"insecure,omitempty"`

	// CloudProviderConfiguration holds the cluster-wide configuration for the
	// vSphere cloud provider.
	CloudProviderConfiguration cloudprovider.Config `json:"cloudProviderConfiguration,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint"`
}

// VSphereClusterStatus defines the observed state of VSphereClusterSpec
type VSphereClusterStatus struct {
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status

// VSphereCluster is the Schema for the vsphereclusters API
type VSphereCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereClusterSpec   `json:"spec,omitempty"`
	Status VSphereClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereClusterList contains a list of VSphereCluster
type VSphereClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VSphereCluster{}, &VSphereClusterList{})
}
