package kubeconfig

import (
	"context"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	kubeconfigAdminPath = filepath.Join("auth", "kubeconfig")
)

// AdminClient is the asset for the admin kubeconfig.
type AdminClient struct {
	kubeconfig
}

var _ asset.WritableAsset = (*AdminClient)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *AdminClient) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.AdminKubeConfigClientCertKey{},
		&tls.KubeAPIServerCompleteCABundle{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *AdminClient) Generate(_ context.Context, parents asset.Parents) error {
	ca := &tls.KubeAPIServerCompleteCABundle{}
	clientCertKey := &tls.AdminKubeConfigClientCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(ca, clientCertKey, installConfig)

	return k.kubeconfig.generate(
		ca,
		clientCertKey,
		getExtAPIServerURL(installConfig.Config),
		installConfig.Config.GetName(),
		"admin",
		kubeconfigAdminPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *AdminClient) Name() string {
	return "Kubeconfig Admin Client"
}

// Load returns the kubeconfig from disk.
func (k *AdminClient) Load(f asset.FileFetcher) (found bool, err error) {
	return k.load(f, kubeconfigAdminPath)
}
