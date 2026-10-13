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
	"k8s.io/apimachinery/pkg/runtime/schema"

	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// Constants block.
const (
	// LaunchTemplateLatestVersion defines the launching of the latest version of the template.
	LaunchTemplateLatestVersion = "$Latest"
)

// GCPMachinePoolSpec defines the desired state of GCPMachinePool.
type GCPMachinePoolSpec struct {
	// ProviderIDList are the identification IDs of machine instances provided by the provider.
	// This field must match the provider IDs as seen on the node objects corresponding to a machine pool's machine instances.
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`

	// InstanceType is the type of instance to create. Example: n1.standard-2
	InstanceType string `json:"instanceType"`

	// Subnet is a reference to the subnetwork to use for this instance. If not specified,
	// the first subnetwork retrieved from the Cluster Region and Network is picked.
	// +optional
	Subnet *string `json:"subnet,omitempty"`

	// ImageFamily is the full reference to a valid image family to be used for this machine.
	// +optional
	ImageFamily *string `json:"imageFamily,omitempty"`

	// Image is the full reference to a valid image to be used for this machine.
	// Takes precedence over ImageFamily.
	// +optional
	Image *string `json:"image,omitempty"`

	// AdditionalLabels is an optional set of tags to add to an instance, in addition to the ones added by default by the
	// GCP provider. If both the GCPCluster and the GCPMachinePool specify the same tag name with different values, the
	// GCPMachinePool's value takes precedence.
	// +optional
	AdditionalLabels capg.Labels `json:"additionalLabels,omitempty"`

	// AdditionalMetadata is an optional set of metadata to add to an instance, in addition to the ones added by default by the
	// GCP provider.
	// +listType=map
	// +listMapKey=key
	// +optional
	AdditionalMetadata []capg.MetadataItem `json:"additionalMetadata,omitempty"`

	// PublicIP specifies whether the instance should get a public IP.
	// Set this to true if you don't have a NAT instances or Cloud Nat setup.
	// +optional
	PublicIP *bool `json:"publicIP,omitempty"`

	// AdditionalNetworkTags is a list of network tags that should be applied to the
	// instance. These tags are set in addition to any network tags defined
	// at the cluster level or in the actuator.
	// +optional
	AdditionalNetworkTags []string `json:"additionalNetworkTags,omitempty"`

	// ResourceManagerTags is an optional set of tags to apply to GCP resources managed
	// by the GCP provider. GCP supports a maximum of 50 tags per resource.
	// +maxItems=50
	// +optional
	ResourceManagerTags capg.ResourceManagerTags `json:"resourceManagerTags,omitempty"`

	// RootDeviceSize is the size of the root volume in GB.
	// Defaults to 30.
	// +optional
	RootDeviceSize int64 `json:"rootDeviceSize,omitempty"`

	// RootDeviceType is the type of the root volume.
	// Supported types of root volumes:
	// 1. "pd-standard" - Standard (HDD) persistent disk
	// 2. "pd-ssd" - SSD persistent disk
	// 3. "pd-balanced" - Balanced Persistent Disk
	// 4. "hyperdisk-balanced" - Hyperdisk Balanced
	// Default is "pd-standard".
	// +optional
	RootDeviceType *capg.DiskType `json:"rootDeviceType,omitempty"`

	// AdditionalDisks are optional non-boot attached disks.
	// +optional
	AdditionalDisks []capg.AttachedDiskSpec `json:"additionalDisks,omitempty"`

	// ServiceAccount specifies the service account email and which scopes to assign to the machine.
	// Defaults to: email: "default", scope: []{compute.CloudPlatformScope}
	// +optional
	ServiceAccount *capg.ServiceAccount `json:"serviceAccounts,omitempty"`

	// Preemptible defines if instance is preemptible
	// +optional
	Preemptible bool `json:"preemptible,omitempty"`

	// ProvisioningModel defines if instance is spot.
	// If set to "Standard" while preemptible is true, then the VM will be of type "Preemptible".
	// If "Spot", VM type is "Spot". When unspecified, defaults to "Standard".
	// +kubebuilder:validation:Enum=Standard;Spot
	// +optional
	ProvisioningModel *capg.ProvisioningModel `json:"provisioningModel,omitempty"`

	// IPForwarding Allows this instance to send and receive packets with non-matching destination or source IPs.
	// This is required if you plan to use this instance to forward routes. Defaults to enabled.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +kubebuilder:default=Enabled
	// +optional
	IPForwarding *capg.IPForwarding `json:"ipForwarding,omitempty"`

	// ShieldedInstanceConfig is the Shielded VM configuration for this machine
	// +optional
	ShieldedInstanceConfig *capg.GCPShieldedInstanceConfig `json:"shieldedInstanceConfig,omitempty"`

	// OnHostMaintenance determines the behavior when a maintenance event occurs that might cause the instance to reboot.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is "Migrate".
	// +kubebuilder:validation:Enum=Migrate;Terminate;
	// +optional
	OnHostMaintenance *capg.HostMaintenancePolicy `json:"onHostMaintenance,omitempty"`

	// ConfidentialCompute Defines whether the instance should have confidential compute enabled or not, and the confidential computing technology of choice.
	// If Disabled, the machine will not be configured to be a confidential computing instance.
	// If Enabled, confidential computing will be configured and AMD Secure Encrypted Virtualization will be configured by default. That is subject to change over time. If using AMD Secure Encrypted Virtualization is vital, use AMDEncryptedVirtualization explicitly instead.
	// If AMDEncryptedVirtualization, it will configure AMD Secure Encrypted Virtualization (AMD SEV) as the confidential computing technology.
	// If AMDEncryptedVirtualizationNestedPaging, it will configure AMD Secure Encrypted Virtualization Secure Nested Paging (AMD SEV-SNP) as the confidential computing technology.
	// If IntelTrustedDomainExtensions, it will configure Intel TDX as the confidential computing technology.
	// If enabled (any value other than Disabled) OnHostMaintenance is required to be set to "Terminate".
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is false.
	// +kubebuilder:validation:Enum=Enabled;Disabled;AMDEncryptedVirtualization;AMDEncryptedVirtualizationNestedPaging;IntelTrustedDomainExtensions
	// +optional
	ConfidentialCompute *capg.ConfidentialComputePolicy `json:"confidentialCompute,omitempty"`

	// RootDiskEncryptionKey defines the KMS key to be used to encrypt the root disk.
	// +optional
	RootDiskEncryptionKey *capg.CustomerEncryptionKey `json:"rootDiskEncryptionKey,omitempty"`

	// GuestAccelerators is a list of the type and count of accelerator cards
	// attached to the instance.
	// +optional
	GuestAccelerators []capg.Accelerator `json:"guestAccelerators,omitempty"`
}

// GCPMachinePoolStatus defines the observed state of GCPMachinePool.
type GCPMachinePoolStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Replicas is the most recently observed number of replicas
	// +optional
	Replicas int32 `json:"replicas"`

	// Conditions defines current service state of the GCPMachinePool.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=gcpmachinepools,scope=Namespaced,categories=cluster-api,shortName=gcpmp
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="MachinePool ready status"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="Number of replicas"

// GCPMachinePool is the Schema for the gcpmachinepools API.
type GCPMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPMachinePoolSpec   `json:"spec,omitempty"`
	Status GCPMachinePoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// GCPMachinePoolList contains a list of GCPMachinePool.
type GCPMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPMachinePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCPMachinePool{}, &GCPMachinePoolList{})
}

// GetObjectKind will return the ObjectKind of an GCPMachinePool.
func (r *GCPMachinePool) GetObjectKind() schema.ObjectKind {
	return &r.TypeMeta
}

// GetObjectKind will return the ObjectKind of an GCPMachinePoolList.
func (r *GCPMachinePoolList) GetObjectKind() schema.ObjectKind {
	return &r.TypeMeta
}

// GCPMachinePool implements the conditions.Setter interface.
var _ conditions.Setter = &GCPMachinePool{}

// SetConditions sets conditions for a MachinePool.
func (r *GCPMachinePool) SetConditions(conditions []metav1.Condition) {
	r.Status.Conditions = conditions
}

// GetConditions gets conditions for a MachinePool.
func (r *GCPMachinePool) GetConditions() []metav1.Condition {
	return r.Status.Conditions
}
