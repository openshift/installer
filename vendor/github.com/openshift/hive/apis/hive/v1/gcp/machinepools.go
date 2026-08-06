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

	// OnHostMaintenance determines the behavior when a maintenance event occurs that might cause the instance to reboot.
	// This is required to be set to "Terminate" if you want to provision machine with attached GPUs.
	// Otherwise, allowed values are "Migrate" and "Terminate".
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is "Migrate".
	// +kubebuilder:validation:Enum=Migrate;Terminate;
	// +optional
	OnHostMaintenance string `json:"onHostMaintenance,omitempty"`

	// ServiceAccount is the email of a gcp service account to be attached to worker nodes
	// in order to provide the permissions required by the cloud provider. For the default
	// worker MachinePool, it is the user's responsibility to match this to the value
	// provided in the install-config.
	//
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`

	// userTags has additional keys and values that we will add as tags to the providerSpec of
	// MachineSets that we creates on GCP. Tag key and tag value should be the shortnames of the
	// tag key and tag value resource. Consumer is responsible for using this only for spokes
	// where custom tags are supported.
	UserTags []UserTag `json:"userTags,omitempty"`

	// Tags defines a set of network tags which will be added to instances in the machineset.
	// Not to be confused with UserTags.
	//
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// OSDisk defines the disk for machines on GCP.
type OSDisk struct {
	// DiskType defines the type of disk.
	// The valid values at this time are: pd-standard, pd-ssd, local-ssd, pd-balanced, hyperdisk-balanced.
	// Defaulted internally to pd-ssd.
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

// UserTag is a tag to apply to GCP resources created for the cluster.
type UserTag struct {
	// parentID is the ID of the hierarchical resource where the tags are defined,
	// e.g. at the Organization or the Project level. To find the Organization ID or Project ID refer to the following pages:
	// https://cloud.google.com/resource-manager/docs/creating-managing-organization#retrieving_your_organization_id,
	// https://cloud.google.com/resource-manager/docs/creating-managing-projects#identifying_projects.
	// An OrganizationID must consist of decimal numbers, and cannot have leading zeroes.
	// A ProjectID must be 6 to 30 characters in length, can only contain lowercase letters,
	// numbers, and hyphens, and must start with a letter, and cannot end with a hyphen.
	ParentID string `json:"parentID"`

	// key is the key part of the tag. A tag key can have a maximum of 63 characters and
	// cannot be empty. Tag key must begin and end with an alphanumeric character, and
	// must contain only uppercase, lowercase alphanumeric characters, and the following
	// special characters `._-`.
	Key string `json:"key"`

	// value is the value part of the tag. A tag value can have a maximum of 63 characters
	// and cannot be empty. Tag value must begin and end with an alphanumeric character, and
	// must contain only uppercase, lowercase alphanumeric characters, and the following
	// special characters `_-.@%=+:,*#&(){}[]` and spaces.
	Value string `json:"value"`
}
