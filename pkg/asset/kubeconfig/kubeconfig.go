package kubeconfig

import (
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	clientcmd "k8s.io/client-go/tools/clientcmd/api/v1"
)

const (
	// KubeconfigUserNameAdmin is the user name of the admin kubeconfig.
	KubeconfigUserNameAdmin = "admin"
	// KubeconfigUserNamekubelet is the user name of the kubelet kubeconfig.
	KubeconfigUserNamekubelet = "kubelet"
)

// Kubeconfig implements the asset.Asset interface that generates
// the admin kubeconfig and kubelet kubeconfig.
type Kubeconfig struct {
	rootDir       string
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
	var err error

	caCertData, err := getCertKeyData(k.rootCA, parents, ".crt")
	if err != nil {
		return nil, err
	}
	clientKeyData, err := getCertKeyData(k.certKey, parents, ".key")
	if err != nil {
		return nil, err
	}
	clientCertData, err := getCertKeyData(k.certKey, parents, ".crt")
	if err != nil {
		return nil, err
	}
	installConfig, err := getInstallConfig(k.installConfig, parents)
	if err != nil {
		return nil, err
	}

	kubeconfig := clientcmd.Config{
		Clusters: []clientcmd.NamedCluster{
			{
				Name: installConfig.Name,
				Cluster: clientcmd.Cluster{
					Server: fmt.Sprintf("https://%s-api.%s:6443", installConfig.Name, installConfig.BaseDomain),
					CertificateAuthorityData: caCertData,
				},
			},
		},
		AuthInfos: []clientcmd.NamedAuthInfo{
			{
				Name: k.userName,
				AuthInfo: clientcmd.AuthInfo{
					ClientCertificateData: clientCertData,
					ClientKeyData:         clientKeyData,
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
		return nil, err
	}

	st := &asset.State{
		Contents: []asset.Content{
			{
				Name: assetFilePath(k.rootDir, k.userName),
				Data: data,
			},
		},
	}

	if err := st.PersistToFile(); err != nil {
		return nil, err
	}

	return st, nil
}

func assetFilePath(rootDir, userName string) string {
	if userName == KubeconfigUserNamekubelet {
		return filepath.Join(rootDir, "auth", "kubeconfig-kubelet")
	}
	return filepath.Join(rootDir, "auth", "kubeconfig")
}

// getCertKeyData extracts the cert or key data from the parent map based on the ext (".crt" or ".key").
// It returns a base64 encoded []bye for the data.
func getCertKeyData(a asset.Asset, parents map[asset.Asset]*asset.State, ext string) ([]byte, error) {
	st, ok := parents[a]
	if !ok {
		return nil, fmt.Errorf("failed to find %T in parents", a)
	}

	var data []byte
	for _, c := range st.Contents {
		if filepath.Ext(c.Name) == ext {
			data = c.Data
			break
		}
	}
	if data == nil {
		return nil, fmt.Errorf("failed to find data in %v with extension == %q", st, ext)
	}

	return []byte(base64.StdEncoding.EncodeToString(data)), nil
}

func getInstallConfig(a asset.Asset, parents map[asset.Asset]*asset.State) (*types.InstallConfig, error) {
	var cfg types.InstallConfig

	st, ok := parents[a]
	if !ok {
		return nil, fmt.Errorf("failed to find %T in parents", a)
	}

	if err := yaml.Unmarshal(st.Contents[0].Data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the installconfig: %v", err)
	}

	return &cfg, nil
}
