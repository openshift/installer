package kubeconfig

import (
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	kubeconfigAdminPath = filepath.Join("auth", "kubeconfig")
)

// Admin is the asset for the admin kubeconfig.
type Admin struct {
	kubeconfig
}

var _ asset.WritableAsset = (*Admin)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *Admin) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.RootCA{},
		&tls.AdminCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *Admin) Generate(parents asset.Parents) error {
	rootCA := &tls.RootCA{}
	adminCertKey := &tls.AdminCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(rootCA, adminCertKey, installConfig)

	return k.kubeconfig.generate(
		rootCA,
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
