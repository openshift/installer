package configgenerator

import (
	"fmt"
	"strings"

	"github.com/coreos/tectonic-installer/installer/pkg/config"

	"github.com/coreos/tectonic-config/config/kube-addon"
	"github.com/coreos/tectonic-config/config/kube-core"
	"github.com/coreos/tectonic-config/config/tectonic-network"
	"github.com/coreos/tectonic-config/config/tectonic-utility"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigGenerator defines the cluster config generation for a cluster.
type ConfigGenerator struct {
	config.Cluster
}

type configurationObject struct {
	metav1.TypeMeta

	Metadata metadata `json:"metadata,omitempty"`
	Data     data     `json:"data,omitempty"`
}

type data map[string]string

type metadata struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type genericData map[string]interface{}

// New returns a ConfigGenerator for a cluster.
func New(cluster config.Cluster) ConfigGenerator {
	return ConfigGenerator{
		Cluster: cluster,
	}
}

// KubeSystem returns, if successful, a yaml string for the kube-system.
func (c ConfigGenerator) KubeSystem() (string, error) {
	return configMap("kube-system", genericData{
		"core-config":    c.coreConfig(),
		"network-config": c.networkConfig(),
	})
}

// TectonicSystem returns, if successful, a yaml string for the tectonic-system.
func (c ConfigGenerator) TectonicSystem() (string, error) {
	return configMap("tectonic-system", genericData{
		"addon-config":   c.addonConfig(),
		"utility-config": c.utilityConfig(),
	})
}

func (c ConfigGenerator) addonConfig() *kubeaddon.OperatorConfig {
	return &kubeaddon.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubeaddon.APIVersion,
			Kind:       kubeaddon.Kind,
		},
	}
}

func (c ConfigGenerator) coreConfig() *kubecore.OperatorConfig {
	coreConfig := kubecore.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubecore.APIVersion,
			Kind:       kubecore.Kind,
		},
	}
	coreConfig.ClusterConfig.APIServerURL = fmt.Sprintf("%s-api.%s", c.Cluster.Name, c.Cluster.DNS.BaseDomain)
	coreConfig.AuthConfig.OIDCClientID = "tectonic-kubectl"
	coreConfig.AuthConfig.OIDCIssuerURL = fmt.Sprintf("%s.%s/identity", c.Cluster.Name, c.Cluster.DNS.BaseDomain)
	coreConfig.AuthConfig.OIDCGroupsClaim = "groups"
	coreConfig.AuthConfig.OIDCUsernameClaim = "email"

	coreConfig.CloudProviderConfig.CloudConfigPath = ""
	coreConfig.CloudProviderConfig.CloudProviderProfile = c.Cluster.Platform

	coreConfig.NetworkConfig.ClusterCIDR = c.Cluster.Networking.NodeCIDR
	coreConfig.NetworkConfig.ServiceCIDR = c.Cluster.Networking.ServiceCIDR
	coreConfig.NetworkConfig.AdvertiseAddress = "0.0.0.0"
	if len(c.Cluster.Etcd.ExternalServers) > 0 {
		coreConfig.NetworkConfig.EtcdServers = strings.Join(c.Cluster.Etcd.ExternalServers, ",")
	} else {
		var etcdServers []string
		for i := 0; i < c.Etcd.NodeCount; i++ {
			etcdServers = append(etcdServers, fmt.Sprintf("https://%s-etcd-%v.%s:2379", c.Cluster.Name, i, c.Cluster.DNS.BaseDomain))
		}
		coreConfig.NetworkConfig.EtcdServers = strings.Join(etcdServers, ",")
	}

	return &coreConfig
}

func (c ConfigGenerator) networkConfig() *tectonicnetwork.OperatorConfig {
	networkConfig := tectonicnetwork.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: tectonicnetwork.APIVersion,
			Kind:       tectonicnetwork.Kind,
		},
	}

	networkConfig.PodCIDR = c.Cluster.Networking.PodCIDR
	networkConfig.CalicoConfig.MTU = c.Cluster.Networking.MTU

	return &networkConfig
}

func (c ConfigGenerator) utilityConfig() *tectonicutility.OperatorConfig {
	utilityConfig := tectonicutility.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: tectonicutility.APIVersion,
			Kind:       tectonicutility.Kind,
		},
	}

	return &utilityConfig
}

func configMap(namespace string, unmarshaledData genericData) (string, error) {
	data := make(data)

	for key, obj := range unmarshaledData {
		str, err := marshalYAML(obj)
		if err != nil {
			return "", err
		}
		data[key] = str
	}

	configurationObject := configurationObject{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Metadata: metadata{
			Name:      "cluster-config-v1",
			Namespace: namespace,
		},
		Data: data,
	}

	str, err := marshalYAML(configurationObject)
	if err != nil {
		return "", err
	}
	return str, nil
}

func marshalYAML(obj interface{}) (string, error) {
	data, err := yaml.Marshal(&obj)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
