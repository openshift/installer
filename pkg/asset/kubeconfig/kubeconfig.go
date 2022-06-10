package kubeconfig

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	clientcmd "k8s.io/client-go/tools/clientcmd/api/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
)

type kubeconfig struct {
	Config *clientcmd.Config
	asset.DefaultFileWriter
}

// generate generates the kubeconfig.
func (k *kubeconfig) generate(
	ca tls.CertInterface,
	clientCertKey tls.CertKeyInterface,
	apiURL string,
	cluster string,
	userName string,
	kubeconfigPath string,
) error {
	k.Config = &clientcmd.Config{
		Clusters: []clientcmd.NamedCluster{
			{
				Name: cluster,
				Cluster: clientcmd.Cluster{
					Server:                   apiURL,
					CertificateAuthorityData: ca.Cert(),
				},
			},
		},
		AuthInfos: []clientcmd.NamedAuthInfo{
			{
				Name: userName,
				AuthInfo: clientcmd.AuthInfo{
					ClientCertificateData: clientCertKey.Cert(),
					ClientKeyData:         clientCertKey.Key(),
				},
			},
		},
		Contexts: []clientcmd.NamedContext{
			{
				Name: userName,
				Context: clientcmd.Context{
					Cluster:  cluster,
					AuthInfo: userName,
				},
			},
		},
		CurrentContext: userName,
	}

	data, err := yaml.Marshal(k.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal kubeconfig")
	}

	k.File = &asset.File{
		Filename: kubeconfigPath,
		Data:     data,
	}

	return nil
}

// load returns the kubeconfig from disk.
func (k *kubeconfig) load(f asset.FileFetcher, name string) (found bool, err error) {
	file, err := f.FetchByName(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	config := &clientcmd.Config{}
	if err := yaml.Unmarshal(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", name)
	}

	k.File, k.Config = file, config
	return true, nil
}

func getExtAPIServerURL(ic *types.InstallConfig) string {
	return fmt.Sprintf("https://api.%s:6443", ic.ClusterDomain())
}

func getIntAPIServerURL(ic *types.InstallConfig) string {
	return fmt.Sprintf("https://api-int.%s:6443", ic.ClusterDomain())
}

func getLoopbackAPIServerURL(ic *types.InstallConfig) string {
	return "https://localhost:6443"
}
