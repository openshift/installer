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
	// ClusterFinalizer allows ReconcileGCPCluster to clean up GCP resources associated with GCPCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "gcpcluster.infrastructure.cluster.x-k8s.io"
)

// GCPClusterSpec defines the desired state of GCPCluster.
type GCPClusterSpec struct {
	// Project is the name of the project to deploy the cluster to.
	Project string `json:"project"`

	// The GCP Region the cluster lives in.
	Region string `json:"region"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// NetworkSpec encapsulates all things related to GCP network.
	// +optional
	Network NetworkSpec `json:"network"`

	// FailureDomains is an optional field which is used to assign selected availability zones to a cluster
	// FailureDomains if empty, defaults to all the zones in the selected region and if specified would override
	// the default zones.
	// +optional
	FailureDomains []string `json:"failureDomains,omitempty"`

	// AdditionalLabels is an optional set of tags to add to GCP resources managed by the GCP provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalLabels Labels `json:"additionalLabels,omitempty"`

	// CredentialsRef is a reference to a Secret that contains the credentials to use for provisioning this cluster. If not
	// supplied then the credentials of the controller will be used.
	// +optional
	CredentialsRef *ObjectReference `json:"credentialsRef,omitempty"`
}

// GCPClusterStatus defines the observed state of GCPCluster.
type GCPClusterStatus struct {
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
	Network        Network                  `json:"network,omitempty"`

	// Bastion Instance `json:"bastion,omitempty"`
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=gcpclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this GCPCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for GCE instances"
// +kubebuilder:printcolumn:name="Network",type="string",JSONPath=".spec.network.name",description="GCP network the cluster is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".status.apiEndpoints[0]",description="API Endpoint",priority=1

// GCPCluster is the Schema for the gcpclusters API.
type GCPCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPClusterSpec   `json:"spec,omitempty"`
	Status GCPClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GCPClusterList contains a list of GCPCluster.
type GCPClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCPCluster{}, &GCPClusterList{})
}
