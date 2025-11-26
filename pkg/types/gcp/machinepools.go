package gcp

import "k8s.io/apimachinery/pkg/util/sets"

// FeatureSwitch indicates whether the feature is enabled or disabled.
type FeatureSwitch string

// OnHostMaintenanceType indicates the setting for the OnHostMaintenance feature, but this is only
// applicable when ConfidentialCompute is Enabled.
type OnHostMaintenanceType string

// ConfidentialComputePolicy indicates the setting for the ConfidentialCompute feature.
type ConfidentialComputePolicy string

const (
	// PDSSD is the constant string representation for persistent disk ssd disk types.
	PDSSD = "pd-ssd"
	// PDStandard is the constant string representation for persistent disk standard disk types.
	PDStandard = "pd-standard"
	// PDBalanced is the constant string representation for persistent disk balanced disk types.
	PDBalanced = "pd-balanced"
	// HyperDiskBalanced is the constant string representation for hyperdisk balanced disk types.
	HyperDiskBalanced = "hyperdisk-balanced"
)

var (
	// ControlPlaneSupportedDisks contains the supported disk types for control plane nodes.
	ControlPlaneSupportedDisks = sets.New(HyperDiskBalanced, PDBalanced, PDSSD)

	// ComputeSupportedDisks contains the supported disk types for control plane nodes.
	ComputeSupportedDisks = sets.New(HyperDiskBalanced, PDBalanced, PDSSD, PDStandard)

	// DefaultCustomInstanceType is the default instance type on the GCP server side. The default custom
	// instance type can be changed on the client side with gcloud.
	DefaultCustomInstanceType = "n1"

	// InstanceTypeToDiskTypeMap contains a map where the key is the Instance Type, and the
	// values are a list of disk types that are supported by the installer and correlate to the Instance Type.
	InstanceTypeToDiskTypeMap = map[string][]string{
		"a2":  {PDStandard, PDSSD, PDBalanced},
		"a3":  {PDSSD, PDBalanced},
		"c2":  {PDStandard, PDSSD, PDBalanced},
		"c2d": {PDStandard, PDSSD, PDBalanced},
		"c3":  {PDSSD, PDBalanced, HyperDiskBalanced},
		"c3d": {PDSSD, PDBalanced, HyperDiskBalanced},
		"c4":  {HyperDiskBalanced},
		"c4a": {HyperDiskBalanced},
		"e2":  {PDStandard, PDSSD, PDBalanced},
		"g2":  {PDStandard, PDSSD, PDBalanced},
		"m1":  {PDSSD, PDBalanced, HyperDiskBalanced},
		"n1":  {PDStandard, PDSSD, PDBalanced},
		"n2":  {PDStandard, PDSSD, PDBalanced},
		"n2d": {PDStandard, PDSSD, PDBalanced},
		"n4":  {HyperDiskBalanced},
		"t2a": {PDStandard, PDSSD, PDBalanced},
		"t2d": {PDStandard, PDSSD, PDBalanced},
	}
)

const (
	// EnabledFeature indicates that the feature is configured as enabled.
	EnabledFeature FeatureSwitch = "Enabled"

	// DisabledFeature indicates that the feature is configured as disabled.
	DisabledFeature FeatureSwitch = "Disabled"

	// OnHostMaintenanceMigrate is the default, and it indicates that the OnHostMaintenance feature is set to Migrate.
	OnHostMaintenanceMigrate OnHostMaintenanceType = "Migrate"

	// OnHostMaintenanceTerminate indicates that the OnHostMaintenance feature is set to Terminate.
	OnHostMaintenanceTerminate OnHostMaintenanceType = "Terminate"

	// ConfidentialComputePolicySEV indicates that the ConfidentialCompute feature is set to AMDEncryptedVirtualization.
	ConfidentialComputePolicySEV ConfidentialComputePolicy = "AMDEncryptedVirtualization"

	// ConfidentialComputePolicySEVSNP indicates that the ConfidentialCompute feature is set to AMDEncryptedVirtualizationNestedPaging.
	ConfidentialComputePolicySEVSNP ConfidentialComputePolicy = "AMDEncryptedVirtualizationNestedPaging"

	// ConfidentialComputePolicyTDX indicates that the ConfidentialCompute feature is set to IntelTrustedDomainExtensions.
	ConfidentialComputePolicyTDX ConfidentialComputePolicy = "IntelTrustedDomainExtensions"
)

var (
	// ConfidentialComputePolicyToSupportedInstanceType is a map containing machine types and the list of confidential computing technologies each of them support.
	ConfidentialComputePolicyToSupportedInstanceType = map[ConfidentialComputePolicy][]string{
		ConfidentialComputePolicySEV:    {"c2d", "n2d", "c3d"},
		ConfidentialComputePolicySEVSNP: {"n2d"},
		ConfidentialComputePolicyTDX:    {"c3"},
	}
)

// MachinePool stores the configuration for a machine pool installed on GCP.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	//
	// +optional
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the GCP instance type.
	// eg. n1-standard-4
	//
	// +optional
	InstanceType string `json:"type"`

	// OSDisk defines the storage for instance.
	//
	// +optional
	OSDisk `json:"osDisk"`

	// OSImage defines a custom image for instance.
	//
	// +optional
	OSImage *OSImage `json:"osImage,omitempty"`

	// Tags defines a set of network tags which will be added to instances in the machineset
	//
	// +optional
	Tags []string `json:"tags,omitempty"`

	// SecureBoot Defines whether the instance should have secure boot enabled.
	// secure boot Verify the digital signature of all boot components, and halt the boot process if signature verification fails.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is false.
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	SecureBoot string `json:"secureBoot,omitempty"`

	// OnHostMaintenance determines the behavior when a maintenance event occurs that might cause the instance to reboot.
	// Allowed values are "Migrate" and "Terminate".
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is "Migrate".
	// +kubebuilder:default="Migrate"
	// +default="Migrate"
	// +kubebuilder:validation:Enum=Migrate;Terminate;
	// +optional
	OnHostMaintenance string `json:"onHostMaintenance,omitempty"`

	// confidentialCompute is an optional field defining whether the instance should have
	// Confidential Computing enabled or not, and the Confidential Computing technology of choice.
	//     With Disabled, Confidential Computing is disabled.
	//     With Enabled, Confidential Computing is enabled with no preference on the
	// Confidential Computing technology. The platform chooses a default i.e. AMD SEV,
	// which is subject to change over time.
	//     With AMDEncryptedVirtualization, Confidential Computing is enabled with
	// AMD Secure Encrypted Virtualization (AMD SEV).
	//     With AMDEncryptedVirtualizationNestedPaging, Confidential Computing is
	// enabled with AMD Secure Encrypted Virtualization Secure Nested Paging
	// (AMD SEV-SNP).
	//     With IntelTrustedDomainExtensions, Confidential Computing is enabled with
	// Intel Trusted Domain Extensions (Intel TDX).
	//     If any value other than Disabled is set, a machine type and region that supports
	// Confidential Computing must be specified. Machine series and regions supporting
	// Confidential Computing technologies can be checked at
	// https://cloud.google.com/confidential-computing/confidential-vm/docs/supported-configurations#machine-type-cpu-zone
	//     If any value other than Disabled is set, onHostMaintenance is required to be set
	// to "Terminate".
	// +kubebuilder:default="Disabled"
	// +default="Disabled"
	// +kubebuilder:validation:Enum="";Enabled;Disabled;AMDEncryptedVirtualization;AMDEncryptedVirtualizationNestedPaging;IntelTrustedDomainExtensions
	// +optional
	ConfidentialCompute string `json:"confidentialCompute,omitempty"`

	// ServiceAccount is the email of a gcp service account to be used during installations.
	// The provided service account can be attached to both control-plane nodes
	// and worker nodes in order to provide the permissions required by the cloud provider.
	//
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

// OSDisk defines the disk for machines on GCP.
type OSDisk struct {
	// DiskType defines the type of disk.
	// For control plane nodes, the valid values are pd-balanced, pd-ssd, and hyperdisk-balanced.
	// +optional
	// +kubebuilder:validation:Enum=pd-balanced;pd-ssd;pd-standard;hyperdisk-balanced
	DiskType string `json:"diskType"`

	// DiskSizeGB defines the size of disk in GB.
	//
	// +kubebuilder:validation:Minimum=16
	// +kubebuilder:validation:Maximum=65536
	DiskSizeGB int64 `json:"DiskSizeGB"`

	// EncryptionKey defines the KMS key to be used to encrypt the disk.
	//
	// +optional
	EncryptionKey *EncryptionKeyReference `json:"encryptionKey,omitempty"`
}

// OSImage defines the image to use for the OS.
type OSImage struct {
	// Name defines the name of the image.
	//
	// +required
	Name string `json:"name"`

	// Project defines the name of the project containing the image.
	//
	// +required
	Project string `json:"project"`
}

// Set sets the values from `required` to `a`.
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}

	if len(required.Zones) > 0 {
		a.Zones = required.Zones
	}

	if required.InstanceType != "" {
		a.InstanceType = required.InstanceType
	}

	if required.Tags != nil {
		a.Tags = required.Tags
	}

	if required.OSDisk.DiskSizeGB > 0 {
		a.OSDisk.DiskSizeGB = required.OSDisk.DiskSizeGB
	}

	if required.OSDisk.DiskType != "" {
		a.OSDisk.DiskType = required.OSDisk.DiskType
	}

	if required.OSImage != nil {
		a.OSImage = required.OSImage
	}

	if required.EncryptionKey != nil {
		if a.EncryptionKey == nil {
			a.EncryptionKey = &EncryptionKeyReference{}
		}
		a.EncryptionKey.Set(required.EncryptionKey)
	}
	if required.SecureBoot != "" {
		a.SecureBoot = required.SecureBoot
	}

	if required.OnHostMaintenance != "" {
		a.OnHostMaintenance = required.OnHostMaintenance
	}

	if required.ConfidentialCompute != "" {
		a.ConfidentialCompute = required.ConfidentialCompute
	}

	if required.ServiceAccount != "" {
		a.ServiceAccount = required.ServiceAccount
	}
}

// EncryptionKeyReference describes the encryptionKey to use for a disk's encryption.
type EncryptionKeyReference struct {
	// KMSKey is a reference to a KMS Key to use for the encryption.
	//
	// +optional
	KMSKey *KMSKeyReference `json:"kmsKey,omitempty"`

	// KMSKeyServiceAccount is the service account being used for the
	// encryption request for the given KMS key. If absent, the Compute
	// Engine default service account is used.
	// See https://cloud.google.com/compute/docs/access/service-accounts#compute_engine_service_account
	// for details on the default service account.
	//
	// +optional
	KMSKeyServiceAccount string `json:"kmsKeyServiceAccount,omitempty"`
}

// Set sets the values from `required` to `e`.
func (e *EncryptionKeyReference) Set(required *EncryptionKeyReference) {
	if required == nil || e == nil {
		return
	}

	if required.KMSKeyServiceAccount != "" {
		e.KMSKeyServiceAccount = required.KMSKeyServiceAccount
	}

	if required.KMSKey != nil {
		if e.KMSKey == nil {
			e.KMSKey = &KMSKeyReference{}
		}
		e.KMSKey.Set(required.KMSKey)
	}
}

// KMSKeyReference gathers required fields for looking up a GCP KMS Key
type KMSKeyReference struct {
	// Name is the name of the customer managed encryption key to be used for the disk encryption.
	Name string `json:"name"`

	// KeyRing is the name of the KMS Key Ring which the KMS Key belongs to.
	KeyRing string `json:"keyRing"`

	// ProjectID is the ID of the Project in which the KMS Key Ring exists.
	// Defaults to the VM ProjectID if not set.
	//
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// Location is the GCP location in which the Key Ring exists.
	Location string `json:"location"`
}

// Set sets the values from `required` to `k`.
func (k *KMSKeyReference) Set(required *KMSKeyReference) {
	if required == nil || k == nil {
		return
	}

	if required.Name != "" {
		k.Name = required.Name
	}

	if required.KeyRing != "" {
		k.KeyRing = required.KeyRing
	}

	if required.ProjectID != "" {
		k.ProjectID = required.ProjectID
	}

	if required.Location != "" {
		k.Location = required.Location
	}
}
