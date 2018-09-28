package machine

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Stock is the stock of master and worker ignition assets.
type Stock interface {
	// MasterIgnition is the asset that generates the master.ign ignition config
	// files for the master nodes.
	MasterIgnition() asset.Asset
	// WorkerIgnition is the asset that generates the worker.ign ignition config
	// file for the worker nodes.
	WorkerIgnition() asset.Asset
}

// StockImpl is the implementation of the bootstrap ignition asset stock.
type StockImpl struct {
	master asset.Asset
	worker asset.Asset
}

// EstablishStock establishes the stock of assets.
func (s *StockImpl) EstablishStock(
	installConfigStock installconfig.Stock,
	tlsStock tls.Stock,
) {
	s.master = newMaster(installConfigStock, tlsStock)
	s.worker = newWorker(installConfigStock, tlsStock)
}

// MasterIgnition returns the master asset.
func (s *StockImpl) MasterIgnition() asset.Asset { return s.master }

// WorkerIgnition returns the worker asset.
func (s *StockImpl) WorkerIgnition() asset.Asset { return s.worker }
