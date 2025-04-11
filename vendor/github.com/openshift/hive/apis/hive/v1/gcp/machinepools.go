package gcp

// MachinePool stores the configuration for a machine pool installed on GCP.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the GCP instance type.
	// eg. n1-standard-4
	InstanceType string `json:"type"`

	// OSDisk defines the storage for instances.
	//
	// +optional
	OSDisk OSDisk `json:"osDisk"`

	// NetworkProjectID specifies which project the network and subnets exist in when
	// they are not in the main ProjectID.
	// +optional
	NetworkProjectID string `json:"networkProjectID,omitempty"`

	// SecureBoot Defines whether the instance should have secure boot enabled.
	// Verifies the digital signature of all boot components, and halts the boot process if signature verification fails.
	// If omitted, the platform chooses a default, which is subject to change over time. Currently that default is "Disabled".
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	SecureBoot string `json:"secureBoot,omitempty"`
}

// OSDisk defines the disk for machines on GCP.
type OSDisk struct {
	// DiskType defines the type of disk.
	// The valid values are pd-standard and pd-ssd.
	// Defaulted internally to pd-ssd.
	// +kubebuilder:validation:Enum=pd-ssd;pd-standard
	// +optional
	DiskType string `json:"diskType,omitempty"`

	// DiskSizeGB defines the size of disk in GB.
	// Defaulted internally to 128.
	//
	// +kubebuilder:validation:Minimum=16
	// +kubebuilder:validation:Maximum=65536
	// +optional
	DiskSizeGB int64 `json:"diskSizeGB,omitempty"`

	// EncryptionKey defines the KMS key to be used to encrypt the disk.
	//
	// +optional
	EncryptionKey *EncryptionKeyReference `json:"encryptionKey,omitempty"`
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
