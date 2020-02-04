package kubeconfig

import (
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	kubeconfigEtcdPath = filepath.Join("auth", "kubeconfig-etcd")
)

// Etcd is the asset for the etcd kubeconfig.
type Etcd struct {
	kubeconfig
}

var _ asset.WritableAsset = (*Etcd)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *Etcd) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.KubeAPIServerCompleteCABundle{},
		&tls.EtcdKubeAPIServerClientCert{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *Etcd) Generate(parents asset.Parents) error {
	ca := &tls.KubeAPIServerCompleteCABundle{}
	clientcertkey := &tls.EtcdKubeAPIServerClientCert{}
	installConfig := &installconfig.InstallConfig{}
	parents.Get(ca, clientcertkey, installConfig)

	return k.kubeconfig.generate(
		ca,
		clientcertkey,
		getIntAPIServerURL(installConfig.Config),
		installConfig.Config.GetName(),
		"etcd-sa",
		kubeconfigEtcdPath,
	)
}

// Name returns the human-friendly name of the asset.
func (k *Etcd) Name() string {
	return "Kubeconfig Etcd"
}

// Load is a no-op because kubelet kubeconfig is not written to disk.
func (k *Etcd) Load(asset.FileFetcher) (bool, error) {
	return false, nil
}
