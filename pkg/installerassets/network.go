package installerassets

import (
	"context"
	"strconv"

	"github.com/ghodss/yaml"
	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func networkConfigRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/cluster-network-02-config.yaml",
		RebuildHelper: networkConfigRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"manifests/cluster-config.yaml",
		"network/host-subnet-length",
	)
	if err != nil {
		return nil, err
	}

	hostSubnetLength, err := strconv.ParseUint(string(parents["network/host-subnet-length"].Data), 10, 32)
	if err != nil {
		return nil, errors.Wrap(err, "parse host subnet length")
	}

	var clusterConfig *corev1.ConfigMap
	err = yaml.Unmarshal(parents["manifests/cluster-config.yaml"].Data, &clusterConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal cluster-config")
	}

	var installConfig *types.InstallConfig
	err = yaml.Unmarshal([]byte(clusterConfig.Data["install-config"]), &installConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal install-config")
	}

	netConfig := installConfig.Networking

	// determine pod address space.
	// This can go away when we get rid of PodCIDR
	// entirely in favor of ClusterNetworks
	var clusterNets []netopv1.ClusterNetwork
	if len(netConfig.ClusterNetworks) > 0 {
		clusterNets = netConfig.ClusterNetworks
	} else if netConfig.PodCIDR == nil || netConfig.PodCIDR.IPNet.IP.IsUnspecified() {
		return nil, errors.Errorf("either PodCIDR or ClusterNetworks must be specified")
	} else {
		clusterNets = []netopv1.ClusterNetwork{
			{
				CIDR:             netConfig.PodCIDR.String(),
				HostSubnetLength: uint32(hostSubnetLength),
			},
		}
	}

	defaultNet := netopv1.DefaultNetworkDefinition{
		Type: netConfig.Type,
	}

	// Add any network-specific configuration defaults here.
	switch netConfig.Type {
	case netopv1.NetworkTypeOpenshiftSDN:
		defaultNet.OpenshiftSDNConfig = &netopv1.OpenshiftSDNConfig{
			// Default to network policy, operator provides all other defaults.
			Mode: netopv1.SDNModePolicy,
		}
	}

	config := &netopv1.NetworkConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: netopv1.SchemeGroupVersion.String(),
			Kind:       "NetworkConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			// not namespaced
		},

		Spec: netopv1.NetworkConfigSpec{
			ServiceNetwork:  netConfig.ServiceCIDR.String(),
			ClusterNetworks: clusterNets,
			DefaultNetwork:  defaultNet,
		},
	}

	asset.Data, err = yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	Rebuilders["manifests/cluster-network-02-config.yaml"] = networkConfigRebuilder
	Rebuilders["network/node-cidr"] = PlatformOverrideRebuilder(
		"network/node-cidr",
		ConstantDefault([]byte("10.0.0.0/16")),
	)
	Defaults["network/cluster-cidr"] = ConstantDefault([]byte("10.128.0.0/14"))
	Defaults["network/host-subnet-length"] = ConstantDefault([]byte("9"))
	Defaults["network/service-cidr"] = ConstantDefault([]byte("172.30.0.0/16"))
}
