package kubeconfig

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/imagebased/configimage"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// ImageBasedAdminClient is the asset for the image-based admin kubeconfig.
type ImageBasedAdminClient struct {
	AdminClient
}

// Dependencies returns the dependency of the kubeconfig.
func (k *ImageBasedAdminClient) Dependencies() []asset.Asset {
	return []asset.Asset{
		&tls.AdminKubeConfigClientCertKey{},
		&configimage.ImageBasedKubeAPIServerCompleteCABundle{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the kubeconfig.
func (k *ImageBasedAdminClient) Generate(parents asset.Parents) error {
	ca := &configimage.ImageBasedKubeAPIServerCompleteCABundle{}
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
