package asset

// Store is a store for the states of assets.
type Store interface {
	// Fetch retrieves the state of the given assets, generating them and their
	// dependencies if necessary.
	Fetch(...WritableAsset) error

	// Destroy removes the asset from all its internal state and also from
	// disk if possible.
	Destroy(Asset) error

	// DestroyState removes everything from the internal state and the internal
	// state file
	DestroyState() error
}
