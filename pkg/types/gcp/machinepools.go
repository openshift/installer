package gcp

import "github.com/openshift/api/machine/v1beta1"

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

	// Tags defines a set of network tags which will be added to instances in the machineset
	//
	// +optional
	Tags []string `json:"tags,omitempty"`

	// OnHostMaintenance determines the behavior when a maintenance event occurs that might cause the instance to reboot.
	// +optional
	OnHostMaintenance v1beta1.GCPHostMaintenanceType `json:"onHostMaintenance,omitempty"`

	// EnableConfidentialCompute Defines whether the instance should have confidential compute enabled.
	// If enabled OnHostMaintenance is required to be set to "Terminate".
	// +optional
	EnableConfidentialCompute bool `json:"enableConfidentialCompute,omitempty"`
}

// OSDisk defines the disk for machines on GCP.
type OSDisk struct {
	// DiskType defines the type of disk.
	// For control plane nodes, the valid value is pd-ssd.
	// +optional
	// +kubebuilder:validation:Enum=pd-ssd;pd-standard
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

	if required.EncryptionKey != nil {
		if a.EncryptionKey == nil {
			a.EncryptionKey = &EncryptionKeyReference{}
		}
		a.EncryptionKey.Set(required.EncryptionKey)
	}

	if required.OnHostMaintenance != "" {
		a.OnHostMaintenance = required.OnHostMaintenance
	}

	if required.EnableConfidentialCompute {
		a.EnableConfidentialCompute = required.EnableConfidentialCompute
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
