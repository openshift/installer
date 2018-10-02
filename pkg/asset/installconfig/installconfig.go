package installconfig

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/ghodss/yaml"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

var (
	defaultServiceCIDR = parseCIDR("10.3.0.0/16")
	defaultPodCIDR     = parseCIDR("10.2.0.0/16")
)

// installConfig generates the install-config.yml file.
type installConfig struct {
	assetStock Stock
}

var _ asset.Asset = (*installConfig)(nil)

// Dependencies returns all of the dependencies directly needed by an
// installConfig asset.
func (a *installConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		a.assetStock.ClusterID(),
		a.assetStock.EmailAddress(),
		a.assetStock.Password(),
		a.assetStock.SSHKey(),
		a.assetStock.BaseDomain(),
		a.assetStock.ClusterName(),
		a.assetStock.PullSecret(),
		a.assetStock.Platform(),
	}
}

// Generate generates the install-config.yml file.
func (a *installConfig) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	clusterID := string(dependencies[a.assetStock.ClusterID()].Contents[0].Data)
	emailAddress := string(dependencies[a.assetStock.EmailAddress()].Contents[0].Data)
	password := string(dependencies[a.assetStock.Password()].Contents[0].Data)
	sshKey := string(dependencies[a.assetStock.SSHKey()].Contents[0].Data)
	baseDomain := string(dependencies[a.assetStock.BaseDomain()].Contents[0].Data)
	clusterName := string(dependencies[a.assetStock.ClusterName()].Contents[0].Data)
	pullSecret := string(dependencies[a.assetStock.PullSecret()].Contents[0].Data)

	installConfig := types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName,
		},
		ClusterID: clusterID,
		Admin: types.Admin{
			Email:    emailAddress,
			Password: password,
			SSHKey:   sshKey,
		},
		BaseDomain: baseDomain,
		Networking: types.Networking{
			// TODO(yifan): Flannel is the temporal default network type for now,
			// Need to update it to the new types.
			Type: "flannel",

			ServiceCIDR: ipnet.IPNet{
				IPNet: defaultServiceCIDR,
			},
			PodCIDR: ipnet.IPNet{
				IPNet: defaultPodCIDR,
			},
		},
		PullSecret: pullSecret,
	}

	platformState := dependencies[a.assetStock.Platform()]
	platform := string(platformState.Contents[0].Data)
	switch platform {
	case AWSPlatformType:
		if err := json.Unmarshal(platformState.Contents[1].Data, &installConfig.AWS); err != nil {
			return nil, err
		}

		installConfig.Machines = []types.MachinePool{
			{
				Name:     "master",
				Replicas: func(x int64) *int64 { return &x }(3),
			},
			{
				Name:     "worker",
				Replicas: func(x int64) *int64 { return &x }(3),
			},
		}
	case OpenStackPlatformType:
		if err := json.Unmarshal(platformState.Contents[1].Data, &installConfig.OpenStack); err != nil {
			return nil, err
		}
		installConfig.Machines = []types.MachinePool{
			{
				Name:     "master",
				Replicas: func(x int64) *int64 { return &x }(3),
			},
			{
				Name:     "worker",
				Replicas: func(x int64) *int64 { return &x }(3),
			},
		}
	case LibvirtPlatformType:
		if err := json.Unmarshal(platformState.Contents[1].Data, &installConfig.Libvirt); err != nil {
			return nil, err
		}
		installConfig.Libvirt.Network.Name = clusterName
		installConfig.Machines = []types.MachinePool{
			{
				Name:     "master",
				Replicas: func(x int64) *int64 { return &x }(1),
			},
			{
				Name:     "worker",
				Replicas: func(x int64) *int64 { return &x }(1),
			},
		}
	default:
		return nil, fmt.Errorf("unknown platform type %q", platform)
	}

	data, err := yaml.Marshal(installConfig)
	if err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: "install-config.yml",
				Data: data,
			},
		},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (a installConfig) Name() string {
	return "Install Config"
}

// GetInstallConfig returns the *types.InstallConfig from the parent asset map.
func GetInstallConfig(installConfig asset.Asset, parents map[asset.Asset]*asset.State) (*types.InstallConfig, error) {
	var cfg types.InstallConfig

	st, ok := parents[installConfig]
	if !ok {
		return nil, fmt.Errorf("failed to find %T in parents", installConfig)
	}

	if err := yaml.Unmarshal(st.Contents[0].Data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the installconfig: %v", err)
	}

	return &cfg, nil
}

// ClusterDNSIP returns the string representation of the DNS server's IP
// address.
func ClusterDNSIP(installConfig *types.InstallConfig) (string, error) {
	ip, err := cidr.Host(&installConfig.ServiceCIDR.IPNet, 10)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}

func parseCIDR(s string) net.IPNet {
	_, cidr, _ := net.ParseCIDR(s)
	return *cidr
}
