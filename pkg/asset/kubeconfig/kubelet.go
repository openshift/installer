package kubeconfig

import (
	"context"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	kubeconfigKubeletPath = filepath.Join("auth", "kubeconfig-kubelet")
)

// Kubelet is the asset for the kubelet kubeconfig.
type Kubelet struct {
	kubeconfig
}

var _ asset.WritableAsset = (*Kubelet)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *Kubelet) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.KubeAPIServerCompleteCABundle{},
		&tls.KubeletClientCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *Kubelet) Generate(_ context.Context, parents asset.Parents) error {
	ca := &tls.KubeAPIServerCompleteCABundle{}
	clientcertkey := &tls.KubeletClientCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(ca, clientcertkey, installConfig)

	return k.kubeconfig.generate(
		ca,
		clientcertkey,
		getIntAPIServerURL(installConfig.Config),
		installConfig.Config.GetName(),
		"kubelet",
		kubeconfigKubeletPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *Kubelet) Name() string {
	return "Kubeconfig Kubelet"
}

// Load is a no-op because kubelet kubeconfig is not written to disk.
func (k *Kubelet) Load(asset.FileFetcher) (bool, error) {
	return false, nil
}
