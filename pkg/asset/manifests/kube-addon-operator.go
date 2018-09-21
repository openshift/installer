package manifests

import (
	"fmt"

	"github.com/ghodss/yaml"

	kubeaddon "github.com/coreos/tectonic-config/config/kube-addon"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

// kubeAddonOperator generates the network-operator-*.yml files
type kubeAddonOperator struct {
	installConfigAsset asset.Asset
	installConfig      *types.InstallConfig
}

var _ asset.Asset = (*kubeAddonOperator)(nil)

// Name returns a human friendly name for the operator
func (kao *kubeAddonOperator) Name() string {
	return "Kube Addon Operator"
}

// Dependencies returns all of the dependencies directly needed by an
// kubeAddonOperator asset.
func (kao *kubeAddonOperator) Dependencies() []asset.Asset {
	return []asset.Asset{
		kao.installConfigAsset,
	}
}

// Generate generates the network-operator-config.yml and network-operator-manifest.yml files
func (kao *kubeAddonOperator) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	ic, err := installconfig.GetInstallConfig(kao.installConfigAsset, dependencies)
	if err != nil {
		return nil, err
	}
	kao.installConfig = ic

	// installconfig is ready, we can create the addon config from it now
	addonConfig, err := kao.addonConfig()
	if err != nil {
		return nil, err
	}

	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: "kube-addon-operator-config.yml",
				Data: addonConfig,
			},
		},
	}
	return state, nil
}

func (kao *kubeAddonOperator) addonConfig() ([]byte, error) {
	addonConfig := kubeaddon.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubeaddon.APIVersion,
			Kind:       kubeaddon.Kind,
		},
	}
	addonConfig.CloudProvider = tectonicCloudProvider(kao.installConfig.Platform)
	addonConfig.ClusterConfig.APIServerURL = kao.getAPIServerURL()
	addonConfig.RegistryHTTPSecret = rand.String(16)
	return yaml.Marshal(addonConfig)
}

func (kao *kubeAddonOperator) getAPIServerURL() string {
	return fmt.Sprintf("https://%s-api.%s:6443", kao.installConfig.Name, kao.installConfig.BaseDomain)
}
