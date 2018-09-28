package manifests

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Stock is the stock of operator assets that can be generated.
type Stock interface {
	// Manifests returns the cluster manifests
	Manifests() asset.Asset

	// ClusterVersionOperator returns the cvo asset object
	ClusterVersionOperator() asset.Asset

	// KubeCoreOperator returns the kco asset object
	KubeCoreOperator() asset.Asset

	// NetworkOperator returns the network operator asset object
	NetworkOperator() asset.Asset

	// KubeAddonOperator returns the addon asset object
	KubeAddonOperator() asset.Asset

	// Mao returns the machine api operator asset object
	Mao() asset.Asset

	// Tectonic returns the tectonic manfests asset object
	Tectonic() asset.Asset
}

// StockImpl implements the Stock interface for manifests
type StockImpl struct {
	manifests              asset.Asset
	clusterVersionOperator asset.Asset
	kubeCoreOperator       asset.Asset
	networkOperator        asset.Asset
	addonOperator          asset.Asset
	mao                    asset.Asset
	tectonic               asset.Asset
}

var _ Stock = (*StockImpl)(nil)

// EstablishStock establishes the stock of assets.
func (s *StockImpl) EstablishStock(stock installconfig.Stock, tlsStock tls.Stock, kubeConfigStock kubeconfig.Stock, machineStock machine.Stock) {
	s.manifests = &manifests{
		assetStock:                s,
		installConfig:             stock.InstallConfig(),
		rootCA:                    tlsStock.RootCA(),
		etcdCA:                    tlsStock.EtcdCA(),
		ingressCertKey:            tlsStock.IngressCertKey(),
		kubeCA:                    tlsStock.KubeCA(),
		aggregatorCA:              tlsStock.AggregatorCA(),
		serviceServingCA:          tlsStock.ServiceServingCA(),
		clusterAPIServerCertKey:   tlsStock.ClusterAPIServerCertKey(),
		etcdClientCertKey:         tlsStock.EtcdClientCertKey(),
		apiServerCertKey:          tlsStock.APIServerCertKey(),
		openshiftAPIServerCertKey: tlsStock.OpenshiftAPIServerCertKey(),
		apiServerProxyCertKey:     tlsStock.APIServerProxyCertKey(),
		kubeletCertKey:            tlsStock.KubeletCertKey(),
		mcsCertKey:                tlsStock.MCSCertKey(),
		serviceAccountKeyPair:     tlsStock.ServiceAccountKeyPair(),
		kubeconfig:                kubeConfigStock.KubeconfigAdmin(),
		workerIgnition:            machineStock.WorkerIgnition(),
	}
	s.kubeCoreOperator = &kubeCoreOperator{
		installConfigAsset: stock.InstallConfig(),
	}
	s.addonOperator = &kubeAddonOperator{
		installConfigAsset: stock.InstallConfig(),
	}
	s.networkOperator = &networkOperator{
		installConfigAsset: stock.InstallConfig(),
	}
	s.mao = &machineAPIOperator{
		installConfigAsset: stock.InstallConfig(),
		aggregatorCA:       tlsStock.AggregatorCA(),
	}
	s.tectonic = &tectonic{
		installConfig:  stock.InstallConfig(),
		ingressCertKey: tlsStock.IngressCertKey(),
		kubeCA:         tlsStock.KubeCA(),
	}
	// TODO:
	//s.clusterVersionOperator = &clusterVersionOperator{}
}

// Manifests returns the manifests asset
func (s *StockImpl) Manifests() asset.Asset { return s.manifests }

// ClusterVersionOperator returns the cvo asset object
func (s *StockImpl) ClusterVersionOperator() asset.Asset { return s.clusterVersionOperator }

// KubeCoreOperator returns the kco asset object
func (s *StockImpl) KubeCoreOperator() asset.Asset { return s.kubeCoreOperator }

// NetworkOperator returns the network operator asset object
func (s *StockImpl) NetworkOperator() asset.Asset { return s.networkOperator }

// KubeAddonOperator returns the addon operator asset object
func (s *StockImpl) KubeAddonOperator() asset.Asset { return s.addonOperator }

// Mao returns the machine API operator asset object
func (s *StockImpl) Mao() asset.Asset { return s.mao }

// Tectonic returns the tectonic manifests asset object
func (s *StockImpl) Tectonic() asset.Asset { return s.tectonic }
