package installerassets

import (
	"context"

	"github.com/ghodss/yaml"
	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"github.com/openshift/installer/pkg/assets"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1a1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func clusterAPIClusterRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/99_openshift-cluster-api_cluster.yaml",
		RebuildHelper: clusterAPIClusterRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"cluster-name",
		"manifests/cluster-network-02-config.yaml",
		"network/service-cidr",
	)
	if err != nil {
		return nil, err
	}

	var netConfig *netopv1.NetworkConfig
	err = yaml.Unmarshal(parents["manifests/cluster-network-02-config.yaml"].Data, &netConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal network config")
	}

	pods := []string{}
	for _, clusterNetwork := range netConfig.Spec.ClusterNetworks {
		pods = append(pods, clusterNetwork.CIDR)
	}

	cluster := clusterv1a1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cluster.k8s.io/v1alpha1",
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      string(parents["cluster-name"].Data),
			Namespace: "openshift-cluster-api",
		},
		Spec: clusterv1a1.ClusterSpec{
			ClusterNetwork: clusterv1a1.ClusterNetworkingConfig{
				Services: clusterv1a1.NetworkRanges{
					CIDRBlocks: []string{string(parents["network/service-cidr"].Data)},
				},
				Pods: clusterv1a1.NetworkRanges{
					CIDRBlocks: pods,
				},
			},
		},
	}

	asset.Data, err = yaml.Marshal(cluster)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func init() {
	Rebuilders["manifests/99_openshift-cluster-api_cluster.yaml"] = clusterAPIClusterRebuilder
}
