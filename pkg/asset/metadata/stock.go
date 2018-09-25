package metadata

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// Stock is the stock of the cluster assets that can be generated.
type Stock interface {
	// Metadata is the asset that generates the metadata.json file
	Metadata() asset.Asset
}

// StockImpl is the implementation of the cluster asset stock.
type StockImpl struct {
	metadata asset.Asset
}

var _ Stock = (*StockImpl)(nil)

// EstablishStock establishes the stock of assets in the specified directory.
func (s *StockImpl) EstablishStock(installConfigStock installconfig.Stock, clusterStock cluster.Stock) {
	s.metadata = &Metadata{
		installConfig: installConfigStock.InstallConfig(),
		cluster:       clusterStock.Cluster(),
	}
}

// Metadata returns the terraform tfvar asset.
func (s *StockImpl) Metadata() asset.Asset { return s.metadata }
