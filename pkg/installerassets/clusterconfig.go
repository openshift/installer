package installerassets

import (
	"context"
	"net"
	"strconv"

	"github.com/ghodss/yaml"
	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func clusterConfigRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "manifests/cluster-config.yaml",
		RebuildHelper: clusterConfigRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"admin/email",
		"admin/password",
		"base-domain",
		"cluster-id",
		"cluster-name",
		"network/cluster-cidr",
		"network/host-subnet-length",
		"network/service-cidr",
		"platform",
		"pull-secret",
		"ssh.pub",
	)
	if err != nil {
		return nil, err
	}

	_, serviceCIDR, err := net.ParseCIDR(string(parents["network/service-cidr"].Data))
	if err != nil {
		return nil, errors.Wrap(err, "parse service CIDR")
	}

	hostSubnetLength, err := strconv.ParseUint(string(parents["network/host-subnet-length"].Data), 10, 32)
	if err != nil {
		return nil, errors.Wrap(err, "parse host subnet length")
	}

	config := &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: string(parents["cluster-name"].Data),
		},
		ClusterID: string(parents["cluster-id"].Data),
		Admin: types.Admin{
			Email:    string(parents["admin/email"].Data),
			Password: string(parents["admin/password"].Data),
			SSHKey:   string(parents["ssh.pub"].Data),
		},
		BaseDomain: string(parents["base-domain"].Data),
		Networking: types.Networking{
			Type: "OpenshiftSDN",
			ServiceCIDR: ipnet.IPNet{
				IPNet: *serviceCIDR,
			},
			ClusterNetworks: []netopv1.ClusterNetwork{
				{
					CIDR:             string(parents["network/cluster-cidr"].Data),
					HostSubnetLength: uint32(hostSubnetLength),
				},
			},
		},
	}

	// support the machine-config operator:
	// Dec 01 00:38:21 wking-bootstrap bootkube.sh[5569]: panic: invalid platform
	// Dec 01 00:38:21 wking-bootstrap bootkube.sh[5569]: goroutine 1 [running]:
	// Dec 01 00:38:21 wking-bootstrap bootkube.sh[5569]: github.com/openshift/machine-config-operator/pkg/operator.platformFromInstallConfig(0x0, 0x0, 0x0, 0x0, 0xc4203d34b0, 0x5, 0x0, 0x0, 0x0, 0x0, ...)
	platform := string(parents["platform"].Data)
	switch platform {
	case "aws":
		config.Platform.AWS = &aws.Platform{}
	case "libvirt":
		config.Platform.Libvirt = &libvirt.Platform{}
	case "openstack":
		config.Platform.OpenStack = &openstack.Platform{}
	default:
		return nil, errors.Errorf("unrecognized platform %q", platform)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-config-v1",
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			"install-config": string(data),
		},
	}

	asset.Data, err = yaml.Marshal(configMap)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	Rebuilders["manifests/cluster-config.yaml"] = clusterConfigRebuilder
}
