/*
Copyright 2022 The Kubernetes Authors.

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
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ManagedMachinePoolFinalizer allows Reconcile to clean up GCP resources associated with the GCPManagedMachinePool before
	// removing it from the apiserver.
	ManagedMachinePoolFinalizer = "gcpmanagedmachinepool.infrastructure.cluster.x-k8s.io"
)

// GCPManagedMachinePoolSpec defines the desired state of GCPManagedMachinePool.
type GCPManagedMachinePoolSpec struct {
	// NodePoolName specifies the name of the GKE node pool corresponding to this MachinePool. If you don't specify a name
	// then a default name will be created based on the namespace and name of the managed machine pool.
	// +optional
	NodePoolName string `json:"nodePoolName,omitempty"`
	// Scaling specifies scaling for the node pool
	// +optional
	Scaling *NodePoolAutoScaling `json:"scaling,omitempty"`
	// KubernetesLabels specifies the labels to apply to the nodes of the node pool.
	// +optional
	KubernetesLabels infrav1.Labels `json:"kubernetesLabels,omitempty"`
	// KubernetesTaints specifies the taints to apply to the nodes of the node pool.
	// +optional
	KubernetesTaints Taints `json:"kubernetesTaints,omitempty"`
	// AdditionalLabels is an optional set of tags to add to GCP resources managed by the GCP provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalLabels infrav1.Labels `json:"additionalLabels,omitempty"`
	// ProviderIDList are the provider IDs of instances in the
	// managed instance group corresponding to the nodegroup represented by this
	// machine pool
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`
}

// GCPManagedMachinePoolStatus defines the observed state of GCPManagedMachinePool.
type GCPManagedMachinePoolStatus struct {
	Ready bool `json:"ready"`
	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`
	// Conditions specifies the cpnditions for the managed machine pool
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Mode",type="string",JSONPath=".spec.mode"
// +kubebuilder:resource:path=gcpmanagedmachinepools,scope=Namespaced,categories=cluster-api,shortName=gcpmmp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// GCPManagedMachinePool is the Schema for the gcpmanagedmachinepools API.
type GCPManagedMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPManagedMachinePoolSpec   `json:"spec,omitempty"`
	Status GCPManagedMachinePoolStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GCPManagedMachinePoolList contains a list of GCPManagedMachinePool.
type GCPManagedMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPManagedMachinePool `json:"items"`
}

// NodePoolAutoScaling specifies scaling options.
type NodePoolAutoScaling struct {
	MinCount *int32 `json:"minCount,omitempty"`
	MaxCount *int32 `json:"maxCount,omitempty"`
}

// GetConditions returns the machine pool conditions.
func (r *GCPManagedMachinePool) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the status conditions for the GCPManagedMachinePool.
func (r *GCPManagedMachinePool) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&GCPManagedMachinePool{}, &GCPManagedMachinePoolList{})
}
