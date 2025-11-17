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
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

const (
	// ManagedMachinePoolFinalizer allows Reconcile to clean up GCP resources associated with the GCPManagedMachinePool before
	// removing it from the apiserver.
	ManagedMachinePoolFinalizer = "gcpmanagedmachinepool.infrastructure.cluster.x-k8s.io"

	// GCPManagedMachinePoolMachineKind indicates the kind of an GCPManagedMachinePoolMachine.
	GCPManagedMachinePoolMachineKind = "GCPManagedMachinePool"
)

// DiskType is type of the disk attached to node.
// +kubebuilder:validation:Enum=pd-standard;pd-ssd;pd-balanced
type DiskType string

const (
	// Standard disk type.
	Standard DiskType = "pd-standard"
	// SSD disk type.
	SSD DiskType = "pd-ssd"
	// Balanced disk type.
	Balanced DiskType = "pd-balanced"
	// HyperdiskBalanced disk type
	HyperdiskBalanced DiskType = "hyperdisk-balanced"
)

// GCPManagedMachinePoolSpec defines the desired state of GCPManagedMachinePool.
type GCPManagedMachinePoolSpec struct {
	GCPManagedMachinePoolClassSpec `json:",inline"`

	// ProviderIDList are the provider IDs of instances in the
	// managed instance group corresponding to the nodegroup represented by this
	// machine pool
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`
}

// NodeNetworkConfig encapsulates node network configurations.
type NodeNetworkConfig struct {
	// Tags is list of instance tags applied to all nodes. Tags
	// are used to identify valid sources or targets for network
	// firewalls.
	// +optional
	Tags []string `json:"tags,omitempty"`
	// CreatePodRange specifies whether to create a new range for
	// pod IPs in this node pool.
	// +optional
	CreatePodRange *bool `json:"createPodRange,omitempty"`
	// PodRangeName is ID of the secondary range for pod IPs.
	// +optional
	PodRangeName *string `json:"podRangeName,omitempty"`
	// PodRangeCidrBlock is the IP address range for pod IPs in
	// this node pool.
	// +optional
	PodRangeCidrBlock *string `json:"podRangeCidrBlock,omitempty"`
}

// NodeSecurityConfig encapsulates node security configurations.
type NodeSecurityConfig struct {
	// ServiceAccount specifies the identity details for node
	// pool.
	// +optional
	ServiceAccount ServiceAccountConfig `json:"serviceAccount,omitempty"`
	// SandboxType is type of the sandbox to use for the node.
	// +optional
	SandboxType *string `json:"sandboxType,omitempty"`
	// EnableSecureBoot defines whether the instance has Secure
	// Boot enabled.
	// +optional
	EnableSecureBoot *bool `json:"enableSecureBoot,omitempty"`
	// EnableIntegrityMonitoring defines whether the instance has
	// integrity monitoring enabled.
	// +optional
	EnableIntegrityMonitoring *bool `json:"enableIntegrityMonitoring,omitempty"`
}

// ServiceAccountConfig encapsulates service account options.
type ServiceAccountConfig struct {
	// Email is the Google Cloud Platform Service Account to be
	// used by the node VMs.
	// +optional
	Email *string `json:"email,omitempty"`
	// Scopes is a set of Google API scopes to be made available
	// on all of the node VMs under the "default" service account.
	// +optional
	Scopes []string `json:"scopes,omitempty"`
}

// GCPManagedMachinePoolStatus defines the observed state of GCPManagedMachinePool.
type GCPManagedMachinePoolStatus struct {
	// Ready denotes that the GCPManagedMachinePool has joined the cluster
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`
	// Conditions specifies the cpnditions for the managed machine pool
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
	// InfrastructureMachineKind is the kind of the infrastructure resources behind MachinePool Machines.
	// +optional
	InfrastructureMachineKind string `json:"infrastructureMachineKind,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Replicas",type="string",JSONPath=".status.replicas"
// +kubebuilder:resource:path=gcpmanagedmachinepools,scope=Namespaced,categories=cluster-api,shortName=gcpmmp
// +kubebuilder:storageversion

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
	// MinCount specifies the minimum number of nodes in the node pool
	// +optional
	MinCount *int32 `json:"minCount,omitempty"`
	// MaxCount specifies the maximum number of nodes in the node pool
	// +optional
	MaxCount *int32 `json:"maxCount,omitempty"`
	// Is autoscaling enabled for this node pool. If unspecified, the default value is true.
	// +optional
	EnableAutoscaling *bool `json:"enableAutoscaling,omitempty"`
	// Location policy used when scaling up a nodepool.
	// +kubebuilder:validation:Enum=balanced;any
	// +optional
	LocationPolicy *ManagedNodePoolLocationPolicy `json:"locationPolicy,omitempty"`
}

// NodePoolManagement specifies auto-upgrade and auto-repair options.
type NodePoolManagement struct {
	// AutoUpgrade specifies whether node auto-upgrade is enabled for the node
	// pool. If enabled, node auto-upgrade helps keep the nodes in your node pool
	// up to date with the latest release version of Kubernetes.
	AutoUpgrade bool `json:"autoUpgrade,omitempty"`
	// AutoRepair specifies whether the node auto-repair is enabled for the node
	// pool. If enabled, the nodes in this node pool will be monitored and, if
	// they fail health checks too many times, an automatic repair action will be
	// triggered.
	AutoRepair bool `json:"autoRepair,omitempty"`
}

// ManagedNodePoolLocationPolicy specifies the location policy of the node pool when autoscaling is enabled.
type ManagedNodePoolLocationPolicy string

// LinuxNodeConfig specifies the settings for Linux agent nodes.
type LinuxNodeConfig struct {
	// Sysctls specifies the sysctl settings for this node pool.
	// +optional
	Sysctls []SysctlConfig `json:"sysctls,omitempty"`
	// CgroupMode specifies the cgroup mode for this node pool.
	// +optional
	CgroupMode *ManagedNodePoolCgroupMode `json:"cgroupMode,omitempty"`
}

// SysctlConfig specifies the sysctl settings for Linux nodes.
type SysctlConfig struct {
	// Parameter specifies sysctl parameter name.
	// +optional
	Parameter string `json:"parameter,omitempty"`
	// Value specifies sysctl parameter value.
	// +optional
	Value string `json:"value,omitempty"`
}

// ManagedNodePoolCgroupMode specifies the cgroup mode of the node pool when autoscaling is enabled.
type ManagedNodePoolCgroupMode int32

const (
	// ManagedNodePoolLocationPolicyBalanced aims to balance the sizes of different zones.
	ManagedNodePoolLocationPolicyBalanced ManagedNodePoolLocationPolicy = "balanced"
	// ManagedNodePoolLocationPolicyAny picks zones that have the highest capacity available.
	ManagedNodePoolLocationPolicyAny ManagedNodePoolLocationPolicy = "any"
)

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
