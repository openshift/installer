package kubeconfig

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Kubelet is the asset for the kubelet kubeconfig.
type Kubelet struct {
	kubeconfig
}

var _ asset.WritableAsset = (*Kubelet)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *Kubelet) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.RootCA{},
		&tls.KubeletCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *Kubelet) Generate(parents asset.Parents) error {
	rootCA := &tls.RootCA{}
	kubeletCertKey := &tls.KubeletCertKey{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(rootCA, kubeletCertKey, installConfig)

	return k.kubeconfig.generate(
		rootCA,
		kubeletCertKey,
		installConfig.Config,
		"kubelet",
		"kubeconfig-kubelet",
	)
}

// Name returns the human-friendly name of the asset.
func (k *Kubelet) Name() string {
	return "Kubeconfig Kubelet"
}
