// Package gcp contains GCP-specific structures for installer
// configuration and management.
// +k8s:deepcopy-gen=package
package gcp

// Name is name for the gcp platform.
const Name string = "gcp"

// AuthorizationMode is the mode or type of authentication indicated in the google credentials struct.
type AuthorizationMode string

const (
	// AuthorizedUserMode indicates that an authorized user without a service account has been used
	// for authentication with the gcloud.
	AuthorizedUserMode AuthorizationMode = "authorized_user"

	// ServiceAccountMode indicates that a service account has been used for authentication with
	// the gcloud.
	ServiceAccountMode AuthorizationMode = "service_account"

	// ExternalAccountMode indicates that an external user such as AWS, Azure, etc. has been used for
	// authentication with gcloud.
	ExternalAccountMode AuthorizationMode = "external_account"
)
