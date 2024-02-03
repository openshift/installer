package nodejoiner

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/store"
)

// NewAddNodesCommand creates a new command for add nodes.
func NewAddNodesCommand(directory string) error {
	fetcher := store.NewAssetsFetcher(directory)
	return fetcher.FetchAndPersist([]asset.WritableAsset{})
}
