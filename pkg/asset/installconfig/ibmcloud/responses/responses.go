package responses

// DNSZoneResponse represents a DNS zone response.
type DNSZoneResponse struct {
	// Name is the domain name of the zone.
	Name string

	// ID is the zone's ID.
	ID string

	// InstanceID is the IBM Cloud Resource ID for the service instance where
	// the DNS zone is managed.
	InstanceID string

	// InstanceCRN is the IBM Cloud Resource CRN for the service instance where
	// the DNS zone is managed.
	InstanceCRN string

	// InstanceName is the display name of the service instance where the DNS zone
	// is managed.
	InstanceName string

	// ResourceGroupID is the resource group ID of the service instance.
	ResourceGroupID string
}

// EncryptionKeyResponse represents an encryption key response.
type EncryptionKeyResponse struct {
	// ID is the key's instance Id.
	ID string

	// Type is the type of key (root, standard).
	Type string

	// CRN is the IBM Cloud CRN representation of the key.
	CRN string

	// State is an integer representing whether the key is enabled.
	State int

	// Deleted is whether the key has been deleted.
	Deleted *bool
}
