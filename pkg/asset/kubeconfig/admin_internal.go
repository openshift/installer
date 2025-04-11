package kubeconfig

import (
	"context"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	kubeconfigAdminInternalPath = filepath.Join("auth", "kubeconfig")
)

// AdminInternalClient is the asset for the admin kubeconfig.
type AdminInternalClient struct {
	kubeconfig
}

var _ asset.WritableAsset = (*AdminInternalClient)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *AdminInternalClient) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.AdminKubeConfigClientCertKey{},
		&tls.KubeAPIServerCompleteCABundle{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *AdminInternalClient) Generate(_ context.Context, parents asset.Parents) error {
	ca := &tls.KubeAPIServerCompleteCABundle{}
	clientCertKey := &tls.AdminKubeConfigClientCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(ca, clientCertKey, installConfig)

	return k.kubeconfig.generate(
		ca,
		clientCertKey,
		getIntAPIServerURL(installConfig.Config),
		installConfig.Config.GetName(),
		"admin",
		kubeconfigAdminInternalPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *AdminInternalClient) Name() string {
	return "Kubeconfig Admin Internal Client"
}

// Load returns the kubeconfig from disk.
func (k *AdminInternalClient) Load(f asset.FileFetcher) (found bool, err error) {
	return false, nil
}
