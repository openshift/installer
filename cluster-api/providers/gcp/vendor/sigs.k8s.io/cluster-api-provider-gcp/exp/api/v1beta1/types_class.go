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

import infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"

// GCPManagedControlPlaneClassSpec defines the GCPManagedControlPlane properties that may be shared across several gcp managed control planes.
type GCPManagedControlPlaneClassSpec struct {
	// MachineTemplate contains information about how machines
	// should be shaped when creating or updating a control plane.
	// For the GCPManagedControlPlaneTemplate, this field is used
	// only to fulfill the CAPI contract.
	// +optional
	MachineTemplate *GCPManagedControlPlaneTemplateMachineTemplate `json:"machineTemplate,omitempty"`

	// ClusterNetwork define the cluster network.
	// +optional
	ClusterNetwork *ClusterNetwork `json:"clusterNetwork,omitempty"`

	// ClusterSecurity defines the cluster security.
	// +optional
	ClusterSecurity *ClusterSecurity `json:"clusterSecurity,omitempty"`

	// Project is the name of the project to deploy the cluster to.
	Project string `json:"project"`

	// Location represents the location (region or zone) in which the GKE cluster
	// will be created.
	Location string `json:"location"`

	// EnableAutopilot indicates whether to enable autopilot for this GKE cluster.
	// +optional
	EnableAutopilot bool `json:"enableAutopilot"`

	// EnableIdentityService indicates whether to enable Identity Service component for this GKE cluster.
	// +optional
	EnableIdentityService bool `json:"enableIdentityService"`

	// ReleaseChannel represents the release channel of the GKE cluster.
	// +optional
	ReleaseChannel *ReleaseChannel `json:"releaseChannel,omitempty"`

	// BinaryAuthorization represents the mode of operation of Binary Authorization for the GKE cluster.
	// This feature is disabled if this field is not specified.
	// +optional
	BinaryAuthorization *BinaryAuthorization `json:"binaryAuthorization,omitempty"`

	// MasterAuthorizedNetworksConfig represents configuration options for master authorized networks feature of the GKE cluster.
	// This feature is disabled if this field is not specified.
	// +optional
	MasterAuthorizedNetworksConfig *MasterAuthorizedNetworksConfig `json:"master_authorized_networks_config,omitempty"`

	// LoggingService represents configuration of logging service feature of the GKE cluster.
	// Possible values: none, logging.googleapis.com/kubernetes (default).
	// Value is ignored when enableAutopilot = true.
	// +optional
	LoggingService *LoggingService `json:"loggingService,omitempty"`

	// MonitoringService represents configuration of monitoring service feature of the GKE cluster.
	// Possible values: none, monitoring.googleapis.com/kubernetes (default).
	// Value is ignored when enableAutopilot = true.
	// +optional
	MonitoringService *MonitoringService `json:"monitoringService,omitempty"`
}

// GCPManagedMachinePoolClassSpec defines the GCPManagedMachinePool properties that may be shared across several GCP managed machinepools.
type GCPManagedMachinePoolClassSpec struct {
	// NodePoolName specifies the name of the GKE node pool corresponding to this MachinePool. If you don't specify a name
	// then a default name will be created based on the namespace and name of the managed machine pool.
	// +optional
	NodePoolName string `json:"nodePoolName,omitempty"`
	// MachineType is the name of a Google Compute Engine [machine
	// type](https://cloud.google.com/compute/docs/machine-types).
	// If unspecified, the default machine type is `e2-medium`.
	// +optional
	MachineType *string `json:"machineType,omitempty"`
	// DiskSizeGb is the size of the disk attached to each node, specified in GB.
	// The smallest allowed disk size is 10GB. If unspecified, the default disk size is 100GB.
	// +optional
	DiskSizeGb *int32 `json:"diskSizeGb,omitempty"`
	// LocalSsdCount is the number of local SSD disks to be attached to the node.
	// +optional
	LocalSsdCount *int32 `json:"localSsdCount,omitempty"`
	// Scaling specifies scaling for the node pool
	// +optional
	Scaling *NodePoolAutoScaling `json:"scaling,omitempty"`
	// NodeLocations is the list of zones in which the NodePool's
	// nodes should be located.
	// +optional
	NodeLocations []string `json:"nodeLocations,omitempty"`
	// ImageType is image type to use for this nodepool.
	// +optional
	ImageType *string `json:"imageType,omitempty"`
	// InstanceType is name of Compute Engine machine type.
	// +optional
	InstanceType *string `json:"instanceType,omitempty"`
	// DiskType is type of the disk attached to each node.
	// +optional
	DiskType *DiskType `json:"diskType,omitempty"`
	// DiskSizeGB is size of the disk attached to each node,
	// specified in GB.
	// +kubebuilder:validation:Minimum:=10
	// +optional
	DiskSizeGB *int64 `json:"diskSizeGB,omitempty"`
	// MaxPodsPerNode is constraint enforced on the max num of
	// pods per node.
	// +kubebuilder:validation:Minimum:=8
	// +kubebuilder:validation:Maximum:=256
	// +optional
	MaxPodsPerNode *int64 `json:"maxPodsPerNode,omitempty"`
	// NodeNetwork specifies the node network configuration
	// options.
	// +optional
	NodeNetwork NodeNetworkConfig `json:"nodeNetwork,omitempty"`
	// NodeSecurity specifies the node security options.
	// +optional
	NodeSecurity NodeSecurityConfig `json:"nodeSecurity,omitempty"`
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
	// Management specifies the node pool management options.
	// +optional
	Management *NodePoolManagement `json:"management,omitempty"`
	// LinuxNodeConfig specifies the settings for Linux agent nodes.
	// +optional
	LinuxNodeConfig *LinuxNodeConfig `json:"linuxNodeConfig,omitempty"`
}
