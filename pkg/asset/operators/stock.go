package operators

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// Stock is the stock of operator assets that can be generated.
type Stock interface {
	// ClusterVersionOperator returns the cvo asset object
	ClusterVersionOperator() asset.Asset

	// KubeCoreOperator returns the kco asset object
	KubeCoreOperator() asset.Asset

	// NetworkOperator returns the network operator asset object
	NetworkOperator() asset.Asset

	// Tnco returns the TNCO asset object
	Tnco() asset.Asset

	// KubeAddonOperator returns the addon asset object
	KubeAddonOperator() asset.Asset
}

// StockImpl implements the Stock interface for operators
type StockImpl struct {
	operators              asset.Asset
	clusterVersionOperator asset.Asset
	kubeCoreOperator       asset.Asset
	networkOperator        asset.Asset
	tnco                   asset.Asset
	addonOperator          asset.Asset
}

var _ Stock = (*StockImpl)(nil)

// EstablishStock establishes the stock of assets in the specified directory.
func (s *StockImpl) EstablishStock(rootDir string, stock installconfig.Stock) {
	s.operators = &operators{
		assetStock: s,
		directory:  rootDir,
	}
	s.kubeCoreOperator = &kubeCoreOperator{
		installConfigAsset: stock.InstallConfig(),
		directory:          rootDir,
	}
	s.addonOperator = &kubeAddonOperator{
		installConfigAsset: stock.InstallConfig(),
		directory:          rootDir,
	}
	s.tnco = &tectonicNodeControllerOperator{
		installConfigAsset: stock.InstallConfig(),
		directory:          rootDir,
	}
	s.networkOperator = &networkOperator{
		installConfigAsset: stock.InstallConfig(),
		directory:          rootDir,
	}
	// TODO:
	//s.clusterVersionOperator = &clusterVersionOperator{}
}

// Operators returns the operators asset
func (s *StockImpl) Operators() asset.Asset { return s.operators }

// ClusterVersionOperator returns the cvo asset object
func (s *StockImpl) ClusterVersionOperator() asset.Asset { return s.clusterVersionOperator }

// KubeCoreOperator returns the kco asset object
func (s *StockImpl) KubeCoreOperator() asset.Asset { return s.kubeCoreOperator }

// NetworkOperator returns the network operator asset object
func (s *StockImpl) NetworkOperator() asset.Asset { return s.networkOperator }

// KubeAddonOperator returns the addon operator asset object
func (s *StockImpl) KubeAddonOperator() asset.Asset { return s.addonOperator }

// Tnco returns the tnc operator asset object
func (s *StockImpl) Tnco() asset.Asset { return s.tnco }
