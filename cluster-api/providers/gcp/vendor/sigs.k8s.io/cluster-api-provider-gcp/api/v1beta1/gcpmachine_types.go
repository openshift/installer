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
	// 4. "pd-balanced" - Balanced Persistent Disk
	// 5. "hyperdisk-balanced" - Hyperdisk Balanced
	// Default is "pd-standard".
	// +optional
	DeviceType *DiskType `json:"deviceType,omitempty"`
	// Size is the size of the disk in GBs.
	// Defaults to 30GB. For "local-ssd" size is always 375GB.
	// +optional
	Size *int64 `json:"size,omitempty"`
	// EncryptionKey defines the KMS key to be used to encrypt the disk.
	// +optional
	EncryptionKey *CustomerEncryptionKey `json:"encryptionKey,omitempty"`
}

// IPForwarding represents the IP forwarding configuration for the GCP machine.
type IPForwarding string

const (
	// IPForwardingEnabled enables the IP forwarding configuration for the GCP machine.
	IPForwardingEnabled IPForwarding = "Enabled"
	// IPForwardingDisabled disables the IP forwarding configuration for the GCP machine.
	IPForwardingDisabled IPForwarding = "Disabled"
)

// SecureBootPolicy represents the secure boot configuration for the GCP machine.
type SecureBootPolicy string

const (
	// SecureBootPolicyEnabled enables the secure boot configuration for the GCP machine.
	SecureBootPolicyEnabled SecureBootPolicy = "Enabled"
	// SecureBootPolicyDisabled disables the secure boot configuration for the GCP machine.
	SecureBootPolicyDisabled SecureBootPolicy = "Disabled"
)

// VirtualizedTrustedPlatformModulePolicy represents the virtualized trusted platform module configuration for the GCP machine.
type VirtualizedTrustedPlatformModulePolicy string

const (
	// VirtualizedTrustedPlatformModulePolicyEnabled enables the virtualized trusted platform module configuration for the GCP machine.
	VirtualizedTrustedPlatformModulePolicyEnabled VirtualizedTrustedPlatformModulePolicy = "Enabled"
	// VirtualizedTrustedPlatformModulePolicyDisabled disables the virtualized trusted platform module configuration for the GCP machine.
	VirtualizedTrustedPlatformModulePolicyDisabled VirtualizedTrustedPlatformModulePolicy = "Disabled"
)

// IntegrityMonitoringPolicy represents the integrity monitoring configuration for the GCP machine.
type IntegrityMonitoringPolicy string

const (
	// IntegrityMonitoringPolicyEnabled enables integrity monitoring for the GCP machine.
	IntegrityMonitoringPolicyEnabled IntegrityMonitoringPolicy = "Enabled"
	// IntegrityMonitoringPolicyDisabled disables integrity monitoring for the GCP machine.
	IntegrityMonitoringPolicyDisabled IntegrityMonitoringPolicy = "Disabled"
)

// GCPShieldedInstanceConfig describes the shielded VM configuration of the instance on GCP.
// Shielded VM configuration allow users to enable and disable Secure Boot, vTPM, and Integrity Monitoring.
type GCPShieldedInstanceConfig struct {
	// SecureBoot Defines whether the instance should have secure boot enabled.
	// Secure Boot verify the digital signature of all boot components, and halting the boot process if signature verification fails.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is Disabled.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	//+optional
	SecureBoot SecureBootPolicy `json:"secureBoot,omitempty"`

	// VirtualizedTrustedPlatformModule enable virtualized trusted platform module measurements to create a known good boot integrity policy baseline.
	// The integrity policy baseline is used for comparison with measurements from subsequent VM boots to determine if anything has changed.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is Enabled.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	VirtualizedTrustedPlatformModule VirtualizedTrustedPlatformModulePolicy `json:"virtualizedTrustedPlatformModule,omitempty"`

	// IntegrityMonitoring determines whether the instance should have integrity monitoring that verify the runtime boot integrity.
	// Compares the most recent boot measurements to the integrity policy baseline and return
	// a pair of pass/fail results depending on whether they match or not.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is Enabled.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	IntegrityMonitoring IntegrityMonitoringPolicy `json:"integrityMonitoring,omitempty"`
}

// ConfidentialComputePolicy represents the confidential compute configuration for the GCP machine.
type ConfidentialComputePolicy string

const (
	// ConfidentialComputePolicyEnabled enables confidential compute for the GCP machine.
	ConfidentialComputePolicyEnabled ConfidentialComputePolicy = "Enabled"
	// ConfidentialComputePolicyDisabled disables confidential compute for the GCP machine.
	ConfidentialComputePolicyDisabled ConfidentialComputePolicy = "Disabled"
)

// ConfidentialVMTechnology represents the confidential computing technology used by a GCP confidential machine.
type ConfidentialVMTechnology string

const (
	// ConfidentialVMTechSEV sets AMD SEV as the VM instance's confidential computing technology of choice.
	ConfidentialVMTechSEV ConfidentialVMTechnology = "sev"
	// ConfidentialVMTechSEVSNP sets AMD SEV-SNP as the VM instance's confidential computing technology of choice.
	ConfidentialVMTechSEVSNP ConfidentialVMTechnology = "sev-snp"
)

// Confidential VM supports Compute Engine machine types in the following series:
// reference: https://cloud.google.com/compute/confidential-vm/docs/os-and-machine-type#machine-type
var confidentialComputeSupportedMachineSeries = map[ConfidentialVMTechnology][]string{
	ConfidentialVMTechSEV:    {"n2d", "c2d", "c3d"},
	ConfidentialVMTechSEVSNP: {"n2d"},
}

// HostMaintenancePolicy represents the desired behavior ase of a host maintenance event.
type HostMaintenancePolicy string

const (
	// HostMaintenancePolicyMigrate causes Compute Engine to live migrate an instance when there is a maintenance event.
	HostMaintenancePolicyMigrate HostMaintenancePolicy = "Migrate"
	// HostMaintenancePolicyTerminate - stops an instance instead of migrating it.
	HostMaintenancePolicyTerminate HostMaintenancePolicy = "Terminate"
)

// KeyType is a type for disk encryption.
type KeyType string

const (
	// CustomerManagedKey (CMEK) references an encryption key stored in Google Cloud KMS.
	CustomerManagedKey KeyType = "Managed"
	// CustomerSuppliedKey (CSEK) specifies an encryption key to use.
	CustomerSuppliedKey KeyType = "Supplied"
)

// ManagedKey is a reference to a key managed by the Cloud Key Management Service.
type ManagedKey struct {
	// KMSKeyName is the name of the encryption key that is stored in Google Cloud KMS. For example:
	// "kmsKeyName": "projects/kms_project_id/locations/region/keyRings/key_region/cryptoKeys/key
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`projects\/[-_[A-Za-z0-9]+\/locations\/[-_[A-Za-z0-9]+\/keyRings\/[-_[A-Za-z0-9]+\/cryptoKeys\/[-_[A-Za-z0-9]+`
	// +kubebuilder:validation:MaxLength=160
	KMSKeyName string `json:"kmsKeyName,omitempty"`
}

// SuppliedKey contains a key for disk encryption. Either RawKey or RSAEncryptedKey must be provided.
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:MaxProperties=1
type SuppliedKey struct {
	// RawKey specifies a 256-bit customer-supplied encryption key, encoded in RFC 4648
	// base64 to either encrypt or decrypt this resource. You can provide either the rawKey or the rsaEncryptedKey.
	// For example: "rawKey": "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0="
	// +optional
	RawKey []byte `json:"rawKey,omitempty"`
	// RSAEncryptedKey specifies an RFC 4648 base64 encoded, RSA-wrapped 2048-bit customer-supplied encryption
	// key to either encrypt or decrypt this resource. You can provide either the rawKey or the
	// rsaEncryptedKey.
	// For example: "rsaEncryptedKey": "ieCx/NcW06PcT7Ep1X6LUTc/hLvUDYyzSZPPVCVPTVEohpeHASqC8uw5TzyO9U+Fka9JFHi
	// z0mBibXUInrC/jEk014kCK/NPjYgEMOyssZ4ZINPKxlUh2zn1bV+MCaTICrdmuSBTWlUUiFoDi
	// D6PYznLwh8ZNdaheCeZ8ewEXgFQ8V+sDroLaN3Xs3MDTXQEMMoNUXMCZEIpg9Vtp9x2oe=="
	// The key must meet the following requirements before you can provide it to Compute Engine:
	// 1. The key is wrapped using a RSA public key certificate provided by Google.
	// 2. After being wrapped, the key must be encoded in RFC 4648 base64 encoding.
	// Gets the RSA public key certificate provided by Google at: https://cloud-certs.storage.googleapis.com/google-cloud-csek-ingress.pem
	// +optional
	RSAEncryptedKey []byte `json:"rsaEncryptedKey,omitempty"`
}

// CustomerEncryptionKey supports both Customer-Managed or Customer-Supplied encryption keys .
type CustomerEncryptionKey struct {
	// KeyType is the type of encryption key. Must be either Managed, aka Customer-Managed Encryption Key (CMEK) or
	// Supplied, aka Customer-Supplied EncryptionKey (CSEK).
	// +kubebuilder:validation:Enum=Managed;Supplied
	KeyType KeyType `json:"keyType"`
	// KMSKeyServiceAccount is the service account being used for the encryption request for the given KMS key.
	// If absent, the Compute Engine default service account is used. For example:
	// "kmsKeyServiceAccount": "name@project_id.iam.gserviceaccount.com.
	// The maximum length is based on the Service Account ID (max 30), Project (max 30), and a valid gcloud email
	// suffix ("iam.gserviceaccount.com").
	// +kubebuilder:validation:MaxLength=85
	// +kubebuilder:validation:Pattern=`[-_[A-Za-z0-9]+@[-_[A-Za-z0-9]+.iam.gserviceaccount.com`
	// +optional
	KMSKeyServiceAccount *string `json:"kmsKeyServiceAccount,omitempty"`
	// ManagedKey references keys managed by the Cloud Key Management Service. This should be set when KeyType is Managed.
	// +optional
	ManagedKey *ManagedKey `json:"managedKey,omitempty"`
	// SuppliedKey provides the key used to create or manage a disk. This should be set when KeyType is Managed.
	// +optional
	SuppliedKey *SuppliedKey `json:"suppliedKey,omitempty"`
}

// ProvisioningModel is a type for Spot VM enablement.
type ProvisioningModel string

const (
	// ProvisioningModelStandard specifies the VM type to NOT be Spot.
	ProvisioningModelStandard ProvisioningModel = "Standard"
	// ProvisioningModelSpot specifies the VM type to be Spot.
	ProvisioningModelSpot ProvisioningModel = "Spot"
)

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

	// ResourceManagerTags is an optional set of tags to apply to GCP resources managed
	// by the GCP provider. GCP supports a maximum of 50 tags per resource.
	// +maxItems=50
	// +optional
	ResourceManagerTags ResourceManagerTags `json:"resourceManagerTags,omitempty"`

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

	// ProvisioningModel defines if instance is spot.
	// If set to "Standard" while preemptible is true, then the VM will be of type "Preemptible".
	// If "Spot", VM type is "Spot". When unspecified, defaults to "Standard".
	// +kubebuilder:validation:Enum=Standard;Spot
	// +optional
	ProvisioningModel *ProvisioningModel `json:"provisioningModel,omitempty"`

	// IPForwarding Allows this instance to send and receive packets with non-matching destination or source IPs.
	// This is required if you plan to use this instance to forward routes. Defaults to enabled.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +kubebuilder:default=Enabled
	// +optional
	IPForwarding *IPForwarding `json:"ipForwarding,omitempty"`

	// ShieldedInstanceConfig is the Shielded VM configuration for this machine
	// +optional
	ShieldedInstanceConfig *GCPShieldedInstanceConfig `json:"shieldedInstanceConfig,omitempty"`

	// OnHostMaintenance determines the behavior when a maintenance event occurs that might cause the instance to reboot.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is "Migrate".
	// +kubebuilder:validation:Enum=Migrate;Terminate;
	// +optional
	OnHostMaintenance *HostMaintenancePolicy `json:"onHostMaintenance,omitempty"`

	// ConfidentialCompute Defines whether the instance should have confidential compute enabled.
	// If enabled OnHostMaintenance is required to be set to "Terminate".
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is false.
	// If ConfidentialInstanceType is configured, even if ConfidentialCompute is Disabled, a confidential compute instance will be configured.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	ConfidentialCompute *ConfidentialComputePolicy `json:"confidentialCompute,omitempty"`

	// confidentialInstanceType determines the required type of confidential computing technology.
	// confidentialInstanceType will precede confidentialCompute. That is, if confidentialCompute is "Disabled" but a valid confidentialInstanceType is specified, a confidential instance will be configured.
	// If confidentialInstanceType isn't set and confidentialCompute is "Enabled" the platform will set the default, which is subject to change over time. Currently the default is "sev" for "c2d", "c3d", and "n2d" machineTypes. For the other machine cases, a valid confidentialInstanceType must be specified.
	// +kubebuilder:validation:Enum=sev;sev-snp;
	// +optional
	ConfidentialInstanceType *ConfidentialVMTechnology `json:"confidentialInstanceType,omitempty"`

	// RootDiskEncryptionKey defines the KMS key to be used to encrypt the root disk.
	// +optional
	RootDiskEncryptionKey *CustomerEncryptionKey `json:"rootDiskEncryptionKey,omitempty"`
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
// +kubebuilder:storageversion
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
