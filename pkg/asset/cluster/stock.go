package cluster

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
)

// Stock is the stock of the cluster assets that can be generated.
type Stock interface {
	// TFVars is the asset that generates the terraform.tfvar file
	TFVars() asset.Asset
	// Cluster is the asset that creates the cluster.
	Cluster() asset.Asset
}

// StockImpl is the implementation of the cluster asset stock.
type StockImpl struct {
	tfvars  asset.Asset
	cluster asset.Asset
}

var _ Stock = (*StockImpl)(nil)

// EstablishStock establishes the stock of assets in the specified directory.
func (s *StockImpl) EstablishStock(installConfigStock installconfig.Stock, bootstrapStock bootstrap.Stock, machineStock machine.Stock, kubeconfigStock kubeconfig.Stock) {
	s.tfvars = &TerraformVariables{
		installConfig:     installConfigStock.InstallConfig(),
		bootstrapIgnition: bootstrapStock.BootstrapIgnition(),
		masterIgnition:    machineStock.MasterIgnition(),
		workerIgnition:    machineStock.WorkerIgnition(),
	}

	s.cluster = &Cluster{
		tfvars:     s.tfvars,
		kubeconfig: kubeconfigStock.KubeconfigAdmin(),
	}
}

// TFVars returns the terraform tfvar asset.
func (s *StockImpl) TFVars() asset.Asset { return s.tfvars }

// Cluster returns the cluster asset.
func (s *StockImpl) Cluster() asset.Asset { return s.cluster }
