package validation

//go:generate mockgen -source=./validvaluesfetcher.go -destination=./mock/validvaluesfetcher_generated.go -package=mock

// ValidValuesFetcher is used to retrieve valid values for fields in Platform.
type ValidValuesFetcher interface {
	// GetCloudNames gets the valid cloud names.
	GetCloudNames() ([]string, error)
	// GetRegionNames gets the valid region names.
	GetRegionNames(cloud string) ([]string, error)
	// GetImageNames gets the valid image names.
	GetImageNames(cloud string) ([]string, error)
	// GetNetworkNames gets the valid network names.
	GetNetworkNames(cloud string) ([]string, error)
	// GetFlavorNames gets the valid flavor names.
	GetFlavorNames(cloud string) ([]string, error)
}
