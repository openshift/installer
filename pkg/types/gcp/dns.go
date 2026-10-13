package gcp

// DNSZoneParams is a set of parameters used to find a DNS zone.
type DNSZoneParams struct {
	// Name is the name of the DNS zone. When provided, the name will be
	// used for the search. When empty any zone matching the other
	// parameters will be returned. Note that either `Name` or `BaseDomain`
	// must be provided.
	Name string

	// Project is the project of the DNS zone.
	Project string

	// IsPublic is true if the DNS zone is public.
	IsPublic bool

	// BaseDomain is the base domain of the DNS zone.
	// Note that either `Name` or `BaseDomain` must be provided.
	BaseDomain string

	// InstallerCreated is true when the DNS zone should be created
	// by the OpenShift Installer (and will be owned by the
	// OpenShift Installer).
	InstallerCreated bool
}
