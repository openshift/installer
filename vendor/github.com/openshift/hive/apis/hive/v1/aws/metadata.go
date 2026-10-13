package aws

// Metadata contains AWS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	// HostedZoneRole is the role to assume when performing operations
	// on a hosted zone owned by another account.
	HostedZoneRole *string `json:"hostedZoneRole,omitempty"`
}
