package kubeconfig

import (
	"context"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	kubeconfigLoopbackPath = filepath.Join("auth", "kubeconfig-loopback")
)

// LoopbackClient is the asset for the admin kubeconfig.
type LoopbackClient struct {
	kubeconfig
}

var _ asset.WritableAsset = (*LoopbackClient)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *LoopbackClient) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.AdminKubeConfigClientCertKey{},
		&tls.KubeAPIServerLocalhostCABundle{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *LoopbackClient) Generate(_ context.Context, parents asset.Parents) error {
	ca := &tls.KubeAPIServerLocalhostCABundle{}
	clientCertKey := &tls.AdminKubeConfigClientCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(ca, clientCertKey, installConfig)

	return k.kubeconfig.generate(
		ca,
		clientCertKey,
		getLoopbackAPIServerURL(installConfig.Config),
		installConfig.Config.GetName(),
		"loopback",
		kubeconfigLoopbackPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *LoopbackClient) Name() string {
	return "Kubeconfig Admin Client (Loopback)"
}

// Load returns the kubeconfig from disk.
func (k *LoopbackClient) Load(f asset.FileFetcher) (found bool, err error) {
	return k.load(f, kubeconfigLoopbackPath)
}
