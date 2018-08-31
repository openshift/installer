package operators

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"

	kubeaddon "github.com/coreos/tectonic-config/config/kube-addon"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// kubeAddonOperator generates the network-operator-*.yml files
type kubeAddonOperator struct {
	installConfigAsset asset.Asset
	installConfig      *types.InstallConfig
	directory          string
}

var _ asset.Asset = (*kubeAddonOperator)(nil)

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

	addonManifest, err := kao.manifest()
	if err != nil {
		return nil, err
	}
	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: filepath.Join(kao.directory, "kube-addon-operator-config.yml"),
				Data: addonConfig,
			},
			{
				Name: filepath.Join(kao.directory, "kube-addon-operator-manifests.yml"),
				Data: []byte(addonManifest),
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
	registrySecret, err := generateRandomID(16)
	if err != nil {
		return nil, err
	}
	addonConfig.RegistryHTTPSecret = registrySecret
	return yaml.Marshal(addonConfig)
}

func (kao *kubeAddonOperator) manifest() (string, error) {
	return "", nil
}

func (kao *kubeAddonOperator) getAPIServerURL() string {
	return fmt.Sprintf("https://%s-api.%s:6443", kao.installConfig.ClusterName, kao.installConfig.BaseDomain)
}

// generateRandomID reproduce tf random_id behaviour
// TODO: re-evaluate solution
func generateRandomID(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)

	n, err := rand.Reader.Read(bytes)
	if n != byteLength {
		return "", errors.New("generated insufficient random bytes")
	}
	if err != nil {
		return "", err
	}

	b64Str := base64.RawURLEncoding.EncodeToString(bytes)

	return b64Str, nil
}
