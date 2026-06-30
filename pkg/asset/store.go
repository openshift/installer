package asset

import (
	"context"
)

// Store is a store for the states of assets.
type Store interface {
	// Fetch retrieves the state of the given asset, generating it and its
	// dependencies if necessary. When purging consumed assets, none of the
	// assets in assetsToPreserve will be purged.
	Fetch(ctx context.Context, assetToFetch Asset, assetsToPreserve ...WritableAsset) error

	// Destroy removes the asset from all its internal state and also from
	// disk if possible.
	Destroy(Asset) error

	// DestroyState removes everything from the internal state and the internal
	// state file
	DestroyState() error

	// Load retrieves the state of the given asset but does not generate it if it
	// does not exist and instead will return nil if not found.
	Load(Asset) (Asset, error)

	// PersistToFile writes the asset's files to disk for the user to edit
	PersistToFile(wa WritableAsset) error
}
