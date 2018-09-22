package manifests

import (
	"fmt"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/apparentlymart/go-cidr/cidr"
	kubecore "github.com/coreos/tectonic-config/config/kube-core"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	authConfigOIDCClientID        = "tectonic-kubectl"
	authConfigOIDCGroupsClaim     = "groups"
	authConfigOIDCUsernameClaim   = "email"
	networkConfigAdvertiseAddress = "0.0.0.0"
)

// kubeCoreOperator generates the kube-core-operator.yaml files
type kubeCoreOperator struct {
	installConfigAsset asset.Asset
	installConfig      *types.InstallConfig
}

var _ asset.Asset = (*kubeCoreOperator)(nil)

// Name returns a human friendly name for the operator
func (kco *kubeCoreOperator) Name() string {
	return "Kube Core Operator"
}

// Dependencies returns all of the dependencies directly needed by an
// kubeCoreOperator asset.
func (kco *kubeCoreOperator) Dependencies() []asset.Asset {
	return []asset.Asset{
		kco.installConfigAsset,
	}
}

// Generate generates the kube-core-operator-config.yml files
func (kco *kubeCoreOperator) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	ic, err := installconfig.GetInstallConfig(kco.installConfigAsset, dependencies)
	if err != nil {
		return nil, err
	}
	kco.installConfig = ic

	// installconfig is ready, we can create the core config from it now
	coreConfig, err := kco.coreConfig()
	if err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(coreConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal core config: %v", err)
	}
	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: "kco-config.yaml",
				Data: data,
			},
		},
	}
	return state, nil
}

func (kco *kubeCoreOperator) coreConfig() (*kubecore.OperatorConfig, error) {
	coreConfig := kubecore.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubecore.APIVersion,
			Kind:       kubecore.Kind,
		},
	}
	coreConfig.ClusterConfig.APIServerURL = kco.getAPIServerURL()
	coreConfig.AuthConfig.OIDCClientID = authConfigOIDCClientID
	coreConfig.AuthConfig.OIDCIssuerURL = kco.getOicdIssuerURL()
	coreConfig.AuthConfig.OIDCGroupsClaim = authConfigOIDCGroupsClaim
	coreConfig.AuthConfig.OIDCUsernameClaim = authConfigOIDCUsernameClaim

	ip, err := cidr.Host(&kco.installConfig.Networking.ServiceCIDR.IPNet, 10)
	if err != nil {
		return nil, err
	}
	coreConfig.DNSConfig.ClusterIP = ip.String()

	coreConfig.CloudProviderConfig.CloudConfigPath = ""
	coreConfig.CloudProviderConfig.CloudProviderProfile = k8sCloudProvider(kco.installConfig.Platform)

	coreConfig.RoutingConfig.Subdomain = kco.getBaseAddress()

	coreConfig.NetworkConfig.ClusterCIDR = kco.installConfig.Networking.PodCIDR.String()
	coreConfig.NetworkConfig.ServiceCIDR = kco.installConfig.Networking.ServiceCIDR.String()
	coreConfig.NetworkConfig.AdvertiseAddress = networkConfigAdvertiseAddress
	coreConfig.NetworkConfig.EtcdServers = strings.Join(kco.getEtcdServersURLs(), ",")

	return &coreConfig, nil
}

func (kco *kubeCoreOperator) getAPIServerURL() string {
	return fmt.Sprintf("https://%s-api.%s:6443", kco.installConfig.Name, kco.installConfig.BaseDomain)
}

func (kco *kubeCoreOperator) getEtcdServersURLs() []string {
	var urls []string
	for i := 0; i < kco.installConfig.MasterCount(); i++ {
		urls = append(urls, fmt.Sprintf("https://%s-etcd-%d.%s:2379", kco.installConfig.Name, i, kco.installConfig.BaseDomain))
	}
	return urls
}

func (kco *kubeCoreOperator) getOicdIssuerURL() string {
	return fmt.Sprintf("https://%s.%s/identity", kco.installConfig.Name, kco.installConfig.BaseDomain)
}

func (kco *kubeCoreOperator) getBaseAddress() string {
	return fmt.Sprintf("%s.%s", kco.installConfig.Name, kco.installConfig.BaseDomain)
}

// Converts a platform to the cloudProvider that k8s understands
func k8sCloudProvider(platform types.Platform) string {
	if platform.AWS != nil {
		return "aws"
	}
	if platform.Libvirt != nil {
		//return "libvirt"
	}
	return ""
}
