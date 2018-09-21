package kubeconfig

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Stock is the stock of the kubeconfig assets that can be generated.
type Stock interface {
	// KubeconfigAdmin is the asset that generates the admin kubeconfig.
	KubeconfigAdmin() asset.Asset
	// KubeconfigKubelet is the asset that generates the kubelet kubeconfig.
	KubeconfigKubelet() asset.Asset
}

// StockImpl implements the Stock interface for kubeconfig assets.
type StockImpl struct {
	kubeconfig        asset.Asset
	kubeconfigKubelet asset.Asset
}

var _ Stock = (*StockImpl)(nil)

// EstablishStock establishes the stock of assets.
func (s *StockImpl) EstablishStock(installConfigStock installconfig.Stock, tlsStock tls.Stock) {
	s.kubeconfig = &Kubeconfig{
		rootCA:        tlsStock.RootCA(),
		certKey:       tlsStock.AdminCertKey(),
		installConfig: installConfigStock.InstallConfig(),
		userName:      KubeconfigUserNameAdmin,
	}

	s.kubeconfigKubelet = &Kubeconfig{
		rootCA:        tlsStock.RootCA(),
		certKey:       tlsStock.KubeletCertKey(),
		installConfig: installConfigStock.InstallConfig(),
		userName:      KubeconfigUserNameKubelet,
	}
}

// KubeconfigAdmin is the asset that generates the admin kubeconfig.
func (s *StockImpl) KubeconfigAdmin() asset.Asset { return s.kubeconfig }

// KubeconfigKubelet is the asset that generates the kubelet kubeconfig.
func (s *StockImpl) KubeconfigKubelet() asset.Asset { return s.kubeconfigKubelet }
