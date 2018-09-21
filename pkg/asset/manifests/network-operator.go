package manifests

import (
	"github.com/ghodss/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"

	tectonicnetwork "github.com/coreos/tectonic-config/config/tectonic-network"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultMTU = "1450"
)

// networkOperator generates the network-operator-*.yml files
type networkOperator struct {
	installConfigAsset asset.Asset
	installConfig      *types.InstallConfig
}

var _ asset.Asset = (*networkOperator)(nil)

// Name returns a human friendly name for the operator
func (no *networkOperator) Name() string {
	return "Network Operator"
}

// Dependencies returns all of the dependencies directly needed by an
// networkOperator asset.
func (no *networkOperator) Dependencies() []asset.Asset {
	return []asset.Asset{
		no.installConfigAsset,
	}
}

// Generate generates the network-operator-config.yml and network-operator-manifest.yml files
func (no *networkOperator) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	ic, err := installconfig.GetInstallConfig(no.installConfigAsset, dependencies)
	if err != nil {
		return nil, err
	}
	no.installConfig = ic

	// installconfig is ready, we can create the core config from it now
	netConfig, err := no.netConfig()
	if err != nil {
		return nil, err
	}

	netManifest, err := no.manifest()
	if err != nil {
		return nil, err
	}
	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: "network-operator-config.yml",
				Data: netConfig,
			},
			{
				Name: "network-operator-manifests.yml",
				Data: netManifest,
			},
		},
	}
	return state, nil
}

func (no *networkOperator) netConfig() ([]byte, error) {
	networkConfig := tectonicnetwork.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: tectonicnetwork.APIVersion,
			Kind:       tectonicnetwork.Kind,
		},
	}

	networkConfig.PodCIDR = no.installConfig.Networking.PodCIDR.String()
	networkConfig.CalicoConfig.MTU = defaultMTU
	networkConfig.NetworkProfile = tectonicnetwork.NetworkType(no.installConfig.Networking.Type)

	return yaml.Marshal(networkConfig)
}

func (no *networkOperator) manifest() ([]byte, error) {
	return []byte(""), nil
}
