package kubeconfig

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Stock is the stock of the kubeconfig assets that can be generated.
type Stock interface {
	Kubeconfig() asset.Asset
	KubeconfigKubelet() asset.Asset
}

// StockImpl implements the Stock interface for kubeconfig assets.
type StockImpl struct {
	kubeconfig        asset.Asset
	kubeconfigKubelet asset.Asset
}

var _ Stock = (*StockImpl)(nil)

// EstablishStock establishes the stock of assets in the specified directory.
func (s *StockImpl) EstablishStock(rootDir string, installConfigStock installconfig.Stock, tlsStock tls.Stock) {
	s.kubeconfig = &Kubeconfig{
		rootDir:       rootDir,
		rootCA:        tlsStock.RootCA(),
		certKey:       tlsStock.AdminCertKey(),
		installConfig: installConfigStock.InstallConfig(),
		userName:      KubeconfigUserNameAdmin,
	}

	s.kubeconfigKubelet = &Kubeconfig{
		rootDir:       rootDir,
		rootCA:        tlsStock.RootCA(),
		certKey:       tlsStock.KubeletCertKey(),
		installConfig: installConfigStock.InstallConfig(),
		userName:      KubeconfigUserNamekubelet,
	}
}

// Kubeconfig is the asset that generates the admin kubeconfig.
func (s *StockImpl) Kubeconfig() asset.Asset { return s.kubeconfig }

// KubeconfigKubelet is the asset that generates the kubelet kubeconfig.
func (s *StockImpl) KubeconfigKubelet() asset.Asset { return s.kubeconfigKubelet }
