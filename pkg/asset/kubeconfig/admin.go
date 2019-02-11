package kubeconfig

import (
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	kubeconfigAdminPath       = filepath.Join("auth", "kubeconfig")
	kubeconfigAdminClientPath = filepath.Join("auth", "kubeconfig-admin")
)

// Admin is the asset for the admin kubeconfig.
// [DEPRECATED]
type Admin struct {
	kubeconfig
}

var _ asset.WritableAsset = (*Admin)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *Admin) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.KubeCA{},
		&tls.AdminCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *Admin) Generate(parents asset.Parents) error {
	kubeCA := &tls.KubeCA{}
	adminCertKey := &tls.AdminCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(kubeCA, adminCertKey, installConfig)

	return k.kubeconfig.generate(
		kubeCA,
		adminCertKey,
		installConfig.Config,
		"admin",
		kubeconfigAdminPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *Admin) Name() string {
	return "Kubeconfig Admin"
}

// Load returns the kubeconfig from disk.
func (k *Admin) Load(f asset.FileFetcher) (found bool, err error) {
	return k.load(f, kubeconfigAdminPath)
}

// AdminClient is the asset for the admin kubeconfig.
type AdminClient struct {
	kubeconfig
}

var _ asset.WritableAsset = (*AdminClient)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *AdminClient) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.AdminKubeConfigClientCertKey{},
		&tls.AdminKubeConfigCABundle{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *AdminClient) Generate(parents asset.Parents) error {
	ca := &tls.AdminKubeConfigCABundle{}
	clientCertKey := &tls.AdminKubeConfigClientCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(ca, clientCertKey, installConfig)

	return k.kubeconfig.generate(
		ca,
		clientCertKey,
		installConfig.Config,
		"admin",
		kubeconfigAdminClientPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *AdminClient) Name() string {
	return "Kubeconfig Admin Client"
}

// Load returns the kubeconfig from disk.
func (k *AdminClient) Load(f asset.FileFetcher) (found bool, err error) {
	return k.load(f, kubeconfigAdminClientPath)
}
