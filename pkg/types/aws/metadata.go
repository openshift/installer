package aws

// Metadata contains AWS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Region string `json:"region"`

	// ServiceEndpoints list contains custom endpoints which will override default
	// service endpoint of AWS Services.
	// There must be only one ServiceEndpoint for a service.
	// +optional
	ServiceEndpoints []ServiceEndpoint `json:"serviceEndpoints,omitempty"`

	// Identifier holds a slice of filter maps.  The maps hold the
	// key/value pairs for the tags we will be matching against.  A
	// resource matches the map if all of the key/value pairs are in its
	// tags.  A resource matches Identifier if it matches any of the maps.
	Identifier []map[string]string `json:"identifier"`
}
