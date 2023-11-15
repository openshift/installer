package store

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
)

// AssetsFetcher it's used to retrieve and resolve a specified set of assets.
type AssetsFetcher interface {
	// Fetchs and persists all the writable assets from the configured assets store.
	FetchAndPersist(context.Context, []asset.WritableAsset) error
}

type fetcher struct {
	storeDir string
}

// NewAssetsFetcher creates a new AssetsFetcher instance for the specified assets store folder.
func NewAssetsFetcher(storeDir string) AssetsFetcher {
	return &fetcher{
		storeDir: storeDir,
	}
}

// Fetchs all the writable assets from the configured assets store.
func (f *fetcher) FetchAndPersist(ctx context.Context, assets []asset.WritableAsset) error {
	assetStore, err := NewStore(f.storeDir)
	if err != nil {
		return fmt.Errorf("failed to create asset store: %w", err)
	}

	for _, a := range assets {
		err := assetStore.Fetch(ctx, a, assets...)
		if err != nil {
			err = errors.Wrapf(err, "failed to fetch %s", a.Name())
		}

		err2 := assetStore.PersistToFile(a)
		if err2 != nil {
			err2 = errors.Wrapf(err2, "failed to write asset (%s) to disk", a.Name())
			if err != nil {
				logrus.Error(err2)
				return err
			}
			return err2
		}

		if err != nil {
			return err
		}
	}

	return nil
}
