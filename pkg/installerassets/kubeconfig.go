package installerassets

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/openshift/installer/pkg/assets"
	clientcmd "k8s.io/client-go/tools/clientcmd/api/v1"
)

func kubeconfigRebuilder(role string, clientKey string, clientCert string) assets.Rebuild {
	return func(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
		asset := &assets.Asset{
			Name:          "auth/kubeconfig-" + role,
			RebuildHelper: kubeconfigRebuilder(role, clientKey, clientCert),
		}

		parents, err := asset.GetParents(
			ctx,
			getByName,
			"base-domain",
			"cluster-name",
			clientKey,
			clientCert,
			"tls/root-ca.crt",
		)
		if err != nil {
			return nil, err
		}

		asset.Data, err = kubeconfig(
			parents["tls/root-ca.crt"].Data,
			parents[clientKey].Data,
			parents[clientCert].Data,
			string(parents["cluster-name"].Data),
			string(parents["base-domain"].Data),
			role,
		)

		return asset, err
	}
}

// kubeconfig renders a YAML kubeconfig from the given parameters.
func kubeconfig(
	rootCACert []byte,
	clientKey []byte,
	clientCert []byte,
	clusterName string,
	baseDomain string,
	userName string,
) ([]byte, error) {
	return yaml.Marshal(&clientcmd.Config{
		Clusters: []clientcmd.NamedCluster{
			{
				Name: clusterName,
				Cluster: clientcmd.Cluster{
					Server: fmt.Sprintf("https://%s-api.%s:6443", clusterName, baseDomain),
					CertificateAuthorityData: rootCACert,
				},
			},
		},
		AuthInfos: []clientcmd.NamedAuthInfo{
			{
				Name: userName,
				AuthInfo: clientcmd.AuthInfo{
					ClientKeyData:         clientKey,
					ClientCertificateData: clientCert,
				},
			},
		},
		Contexts: []clientcmd.NamedContext{
			{
				Name: userName,
				Context: clientcmd.Context{
					Cluster:  clusterName,
					AuthInfo: userName,
				},
			},
		},
		CurrentContext: userName,
	})
}

func init() {
	Rebuilders["auth/kubeconfig-admin"] = kubeconfigRebuilder(
		"admin",
		"tls/admin-client.key",
		"tls/admin-client.crt",
	)

	Rebuilders["auth/kubeconfig-kubelet"] = kubeconfigRebuilder(
		"kubelet",
		"tls/kubelet-client.key",
		"tls/kubelet-client.crt",
	)
}
