package nodejoiner

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/store"
)

// NewMonitorAddNodesCommand creates a new command for monitor add nodes.
func NewMonitorAddNodesCommand(directory string) error {
	fetcher := store.NewAssetsFetcher(directory)
	return fetcher.FetchAndPersist([]asset.WritableAsset{})
}
