package aws

// Metadata contains AWS metadata (e.g. for uninstalling the cluster).
type Metadata struct {
	Region string `json:"region"`

	// Identifier holds a slice of filter maps.  The maps hold the
	// key/value pairs for the tags we will be matching against.  A
	// resource matches the map if all of the key/value pairs are in its
	// tags.  A resource matches Identifier if it matches any of the maps.
	Identifier []map[string]string `json:"identifier"`
}
