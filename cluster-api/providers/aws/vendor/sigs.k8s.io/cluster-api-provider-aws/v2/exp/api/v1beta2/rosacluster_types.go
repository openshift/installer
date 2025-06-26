/*
Copyright 2023 The Kubernetes Authors.

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

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ROSAClusterSpec defines the desired state of ROSACluster.
type ROSAClusterSpec struct {
	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// ROSAClusterStatus defines the observed state of ROSACluster.
type ROSAClusterStatus struct {
	// Ready is when the ROSAControlPlane has a API server URL.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// FailureDomains specifies a list fo available availability zones that can be used
	// +optional
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`

	// Conditions defines current service state of the ROSACluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosaclusters,scope=Namespaced,categories=cluster-api,shortName=rosac
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AWSManagedControl belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Control plane infrastructure is ready for worker nodes"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1

// ROSACluster is the Schema for the ROSAClusters API.
type ROSACluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ROSAClusterSpec   `json:"spec,omitempty"`
	Status ROSAClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ROSAClusterList contains a list of ROSACluster.
type ROSAClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSACluster `json:"items"`
}

// GetConditions returns the observations of the operational state of the
// ROSACluster resource.
func (r *ROSACluster) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the ROSACluster to the
// predescribed clusterv1.Conditions.
func (r *ROSACluster) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&ROSACluster{}, &ROSAClusterList{})
}
