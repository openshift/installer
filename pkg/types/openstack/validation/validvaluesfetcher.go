package validation

//go:generate mockgen -source=./validvaluesfetcher.go -destination=./mock/validvaluesfetcher_generated.go -package=mock

// ValidValuesFetcher is used to retrieve valid values for fields in Platform.
type ValidValuesFetcher interface {
	// GetCloudNames gets the valid cloud names.
	GetCloudNames() ([]string, error)
	// GetNetworkNames gets the valid network names.
	GetNetworkNames(cloud string) ([]string, error)
	// GetFlavorNames gets the valid flavor names.
	GetFlavorNames(cloud string) ([]string, error)
	// GetNetworkExtensionsAliases gets the aliases for all the networking enabled extensions
	GetNetworkExtensionsAliases(cloud string) ([]string, error)
	// GetServiceCatalog gets the catalog service names
	GetServiceCatalog(cloud string) ([]string, error)
	// GetFloatingIPNames gets the floating IPs
	GetFloatingIPNames(cloud string, floatingNetwork string) ([]string, error)
}
