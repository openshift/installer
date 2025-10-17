package dataplane

const (
	// MsiIdentityURLHeader is provided by ARM in responses for resource creation
	// to specify the URL at which clients can get credentials for a managed identity
	// associated with the ARM resource being created.
	MsiIdentityURLHeader = "x-ms-identity-url"
	// MsiPrincipalIDHeader is provided by ARM in responses for resource creation
	// to specify the service principal ID for a managed identity associated with
	// the ARM resource being created.
	MsiPrincipalIDHeader = "x-ms-identity-principal-id"
	// MsiTenantHeader is provided by ARM in responses for resource creation to specify
	// the tenant id for a managed identity associated with the ARM resource being created.
	MsiTenantHeader = "x-ms-home-tenant-id"
)

const (
	// ManagedIdentityCredentialsStoragePrefix is a suggested prefix to use when
	// storing a ManagedIdentityCredentials object in Azure KeyVault.
	ManagedIdentityCredentialsStoragePrefix = "msi-"
	// UserAssignedIdentityCredentialsStoragePrefix is a suggested prefix to use when
	// storing a UserAssignedIdentityCredentials object in Azure KeyVault.
	UserAssignedIdentityCredentialsStoragePrefix = "uamsi-"
)
