package operators

import (
	"net"
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/apparentlymart/go-cidr/cidr"
	tnc "github.com/coreos/tectonic-config/config/tectonic-node-controller"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// tectonicNodeControllerOperator generates the tnco-operator.yaml files
type tectonicNodeControllerOperator struct {
	installConfigAsset asset.Asset
	installConfig      *types.InstallConfig
	directory          string
}

var _ asset.Asset = (*tectonicNodeControllerOperator)(nil)

// Dependencies returns all of the dependencies directly needed by an
// tectonicNodeControllerOperator asset.
func (tnco *tectonicNodeControllerOperator) Dependencies() []asset.Asset {
	return []asset.Asset{
		tnco.installConfigAsset,
	}
}

// Generate generates the tnco-operator-config.yml files
func (tnco *tectonicNodeControllerOperator) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	ic, err := installconfig.GetInstallConfig(tnco.installConfigAsset, dependencies)
	if err != nil {
		return nil, err
	}
	tnco.installConfig = ic

	// installconfig is ready, we can create the tnco config from it now
	tncoConfig, err := tnco.tncoConfig()
	if err != nil {
		return nil, err
	}

	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: filepath.Join(tnco.directory, "tnco-config.yml"),
				Data: tncoConfig,
			},
		},
	}
	return state, nil
}

func (tnco *tectonicNodeControllerOperator) tncoConfig() ([]byte, error) {
	tncoConfig := tnc.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: tnc.TNCOConfigAPIVersion,
			Kind:       tnc.TNCOConfigKind,
		},
	}

	tncoConfig.ControllerConfig = tnc.ControllerConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: tnc.TNCConfigAPIVersion,
			Kind:       tnc.TNCConfigKind,
		},
	}

	svcCidr := tnco.installConfig.Networking.ServiceCIDR
	ip, err := cidr.Host(&net.IPNet{IP: svcCidr.IP, Mask: svcCidr.Mask}, 10)
	if err != nil {
		return nil, err
	}
	tncoConfig.ControllerConfig.ClusterDNSIP = ip.String()

	tncoConfig.ControllerConfig.Platform = tectonicCloudProvider(tnco.installConfig.Platform)
	tncoConfig.ControllerConfig.CloudProviderConfig = "" // TODO(yifan): Get CloudProviderConfig.
	tncoConfig.ControllerConfig.ClusterName = tnco.installConfig.ClusterName
	tncoConfig.ControllerConfig.BaseDomain = tnco.installConfig.BaseDomain
	tncoConfig.ControllerConfig.EtcdInitialCount = 1           // TODO (rchopra): confirmed?
	tncoConfig.ControllerConfig.AdditionalConfigs = []string{} // TODO(yifan): Get additional configs.
	tncoConfig.ControllerConfig.NodePoolUpdateLimit = nil      // TODO(yifan): Get the node pool update limit.

	return yaml.Marshal(tncoConfig)
}
