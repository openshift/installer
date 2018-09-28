package bootstrap

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Stock is the stock of the bootstrap ignition asset.
type Stock interface {
	// BootstrapIgnition is the asset that generates the bootstrap.ign ignition
	// config file for the bootstrap node.
	BootstrapIgnition() asset.Asset
}

// StockImpl is the implementation of the master and worker ignition
// asset stock.
type StockImpl struct {
	boostrap asset.Asset
}

// EstablishStock establishes the stock of assets.
func (s *StockImpl) EstablishStock(
	installConfigStock installconfig.Stock,
	tlsStock tls.Stock,
	kubeconfigStock kubeconfig.Stock,
	manifestStock manifests.Stock,
) {
	s.boostrap = newBootstrap(installConfigStock, tlsStock, kubeconfigStock, manifestStock)
}

// BootstrapIgnition returns the bootstrap asset.
func (s *StockImpl) BootstrapIgnition() asset.Asset { return s.boostrap }
