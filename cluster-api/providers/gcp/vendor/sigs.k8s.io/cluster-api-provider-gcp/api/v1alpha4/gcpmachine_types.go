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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	// MachineFinalizer allows ReconcileGCPMachine to clean up GCP resources associated with GCPMachine before
	// removing it from the apiserver.
	MachineFinalizer = "gcpmachine.infrastructure.cluster.x-k8s.io"
)

// DiskType is a type to use to define with disk type will be used.
type DiskType string

const (
	// PdStandardDiskType defines the name for the standard disk.
	PdStandardDiskType DiskType = "pd-standard"
	// PdSsdDiskType defines the name for the ssd disk.
	PdSsdDiskType DiskType = "pd-ssd"
	// LocalSsdDiskType defines the name for the local ssd disk.
	LocalSsdDiskType DiskType = "local-ssd"
)

// AttachedDiskSpec degined GCP machine disk.
type AttachedDiskSpec struct {
	// DeviceType is a device type of the attached disk.
	// Supported types of non-root attached volumes:
	// 1. "pd-standard" - Standard (HDD) persistent disk
	// 2. "pd-ssd" - SSD persistent disk
	// 3. "local-ssd" - Local SSD disk (https://cloud.google.com/compute/docs/disks/local-ssd).
	// Default is "pd-standard".
	// +optional
	DeviceType *DiskType `json:"deviceType,omitempty"`
	// Size is the size of the disk in GBs.
	// Defaults to 30GB. For "local-ssd" size is always 375GB.
	// +optional
	Size *int64 `json:"size,omitempty"`
}

// GCPMachineSpec defines the desired state of GCPMachine.
type GCPMachineSpec struct {
	// InstanceType is the type of instance to create. Example: n1.standard-2
	InstanceType string `json:"instanceType"`

	// Subnet is a reference to the subnetwork to use for this instance. If not specified,
	// the first subnetwork retrieved from the Cluster Region and Network is picked.
	// +optional
	Subnet *string `json:"subnet,omitempty"`

	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// ImageFamily is the full reference to a valid image family to be used for this machine.
	// +optional
	ImageFamily *string `json:"imageFamily,omitempty"`

	// Image is the full reference to a valid image to be used for this machine.
	// Takes precedence over ImageFamily.
	// +optional
	Image *string `json:"image,omitempty"`

	// AdditionalLabels is an optional set of tags to add to an instance, in addition to the ones added by default by the
	// GCP provider. If both the GCPCluster and the GCPMachine specify the same tag name with different values, the
	// GCPMachine's value takes precedence.
	// +optional
	AdditionalLabels Labels `json:"additionalLabels,omitempty"`

	// AdditionalMetadata is an optional set of metadata to add to an instance, in addition to the ones added by default by the
	// GCP provider.
	// +listType=map
	// +listMapKey=key
	// +optional
	AdditionalMetadata []MetadataItem `json:"additionalMetadata,omitempty"`

	// IAMInstanceProfile is a name of an IAM instance profile to assign to the instance
	// +optional
	// IAMInstanceProfile string `json:"iamInstanceProfile,omitempty"`

	// PublicIP specifies whether the instance should get a public IP.
	// Set this to true if you don't have a NAT instances or Cloud Nat setup.
	// +optional
	PublicIP *bool `json:"publicIP,omitempty"`

	// AdditionalNetworkTags is a list of network tags that should be applied to the
	// instance. These tags are set in addition to any network tags defined
	// at the cluster level or in the actuator.
	// +optional
	AdditionalNetworkTags []string `json:"additionalNetworkTags,omitempty"`

	// RootDeviceSize is the size of the root volume in GB.
	// Defaults to 30.
	// +optional
	RootDeviceSize int64 `json:"rootDeviceSize,omitempty"`

	// RootDeviceType is the type of the root volume.
	// Supported types of root volumes:
	// 1. "pd-standard" - Standard (HDD) persistent disk
	// 2. "pd-ssd" - SSD persistent disk
	// Default is "pd-standard".
	// +optional
	RootDeviceType *DiskType `json:"rootDeviceType,omitempty"`

	// AdditionalDisks are optional non-boot attached disks.
	// +optional
	AdditionalDisks []AttachedDiskSpec `json:"additionalDisks,omitempty"`

	// ServiceAccount specifies the service account email and which scopes to assign to the machine.
	// Defaults to: email: "default", scope: []{compute.CloudPlatformScope}
	// +optional
	ServiceAccount *ServiceAccount `json:"serviceAccounts,omitempty"`

	// Preemptible defines if instance is preemptible
	// +optional
	Preemptible bool `json:"preemptible,omitempty"`
}

// MetadataItem defines a single piece of metadata associated with an instance.
type MetadataItem struct {
	// Key is the identifier for the metadata entry.
	Key string `json:"key"`
	// Value is the value of the metadata entry.
	Value *string `json:"value,omitempty"`
}

// GCPMachineStatus defines the observed state of GCPMachine.
type GCPMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the GCP instance associated addresses.
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// InstanceStatus is the status of the GCP instance for this machine.
	// +optional
	InstanceStatus *InstanceStatus `json:"instanceState,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=gcpmachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this GCPMachine belongs"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.instanceState",description="GCE instance state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="InstanceID",type="string",JSONPath=".spec.providerID",description="GCE instance ID"
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns with this GCPMachine"

// GCPMachine is the Schema for the gcpmachines API.
type GCPMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPMachineSpec   `json:"spec,omitempty"`
	Status GCPMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GCPMachineList contains a list of GCPMachine.
type GCPMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCPMachine{}, &GCPMachineList{})
}
