package cluster

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// Stock is the stock of the cluster assets that can be generated.
type Stock interface {
	// TFVars is the asset that generates the terraform.tfvar file
	TFVars() asset.Asset
}

// StockImpl is the implementation of the cluster asset stock.
type StockImpl struct {
	tfvars asset.Asset
}

var _ Stock = (*StockImpl)(nil)

// EstablishStock establishes the stock of assets in the specified directory.
func (s *StockImpl) EstablishStock(rootDir string, installConfigStock installconfig.Stock, ignitionStock ignition.Stock) {
	s.tfvars = &TerraformVariables{
		rootDir:           rootDir,
		installConfig:     installConfigStock.InstallConfig(),
		bootstrapIgnition: ignitionStock.BootstrapIgnition(),
		masterIgnition:    ignitionStock.MasterIgnition(),
		workerIgnition:    ignitionStock.WorkerIgnition(),
	}
}

// TFVars returns the terraform tfvar asset.
func (s *StockImpl) TFVars() asset.Asset { return s.tfvars }
