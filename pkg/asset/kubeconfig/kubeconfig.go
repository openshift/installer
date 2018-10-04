package kubeconfig

import (
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	clientcmd "k8s.io/client-go/tools/clientcmd/api/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

const (
	// KubeconfigUserNameAdmin is the user name of the admin kubeconfig.
	KubeconfigUserNameAdmin = "admin"
	// KubeconfigUserNameKubelet is the user name of the kubelet kubeconfig.
	KubeconfigUserNameKubelet = "kubelet"
)

// Kubeconfig implements the asset.Asset interface that generates
// the admin kubeconfig and kubelet kubeconfig.
type Kubeconfig struct {
	userName      string // admin or kubelet.
	rootCA        asset.Asset
	certKey       asset.Asset
	installConfig asset.Asset
}

var _ asset.Asset = (*Kubeconfig)(nil)

// Dependencies returns the dependency of the kubeconfig.
func (k *Kubeconfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		k.rootCA,
		k.certKey,
		k.installConfig,
	}
}

// Generate generates the kubeconfig.
func (k *Kubeconfig) Generate(parents map[asset.Asset]*asset.State) (*asset.State, error) {
	installConfig, err := installconfig.GetInstallConfig(k.installConfig, parents)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get InstallConfig from parents")
	}

	kubeconfig := clientcmd.Config{
		Clusters: []clientcmd.NamedCluster{
			{
				Name: installConfig.Name,
				Cluster: clientcmd.Cluster{
					Server: fmt.Sprintf("https://%s-api.%s:6443", installConfig.Name, installConfig.BaseDomain),
					CertificateAuthorityData: parents[k.rootCA].Contents[tls.CertIndex].Data,
				},
			},
		},
		AuthInfos: []clientcmd.NamedAuthInfo{
			{
				Name: k.userName,
				AuthInfo: clientcmd.AuthInfo{
					ClientCertificateData: parents[k.certKey].Contents[tls.CertIndex].Data,
					ClientKeyData:         parents[k.certKey].Contents[tls.KeyIndex].Data,
				},
			},
		},
		Contexts: []clientcmd.NamedContext{
			{
				Name: k.userName,
				Context: clientcmd.Context{
					Cluster:  installConfig.Name,
					AuthInfo: k.userName,
				},
			},
		},
		CurrentContext: k.userName,
	}

	data, err := yaml.Marshal(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Marshal kubeconfig")
	}

	var kubeconfigName string
	switch k.userName {
	case KubeconfigUserNameAdmin:
		kubeconfigName = "kubeconfig"
	default:
		kubeconfigName = fmt.Sprintf("kubeconfig-%s", k.userName)
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: filepath.Join("auth", kubeconfigName),
				Data: data,
			},
		},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (k *Kubeconfig) Name() string {
	return "Kubeconfig"
}
