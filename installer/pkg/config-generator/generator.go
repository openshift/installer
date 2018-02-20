package configgenerator

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/coreos/tectonic-installer/installer/pkg/config"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/coreos/tectonic-config/config/kube-addon"
	"github.com/coreos/tectonic-config/config/kube-core"
	"github.com/coreos/tectonic-config/config/tectonic-network"
	"github.com/coreos/tectonic-config/config/tectonic-utility"
	"github.com/ghodss/yaml"
	"golang.org/x/crypto/bcrypt"
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
		"kco-config":     c.coreConfig(),
		"network-config": c.networkConfig(),
	})
}

// TectonicSystem returns, if successful, a yaml string for the tectonic-system.
func (c ConfigGenerator) TectonicSystem() (string, error) {
	utilityConfig, err := c.utilityConfig()
	if err != nil {
		return "", err
	}
	addonConfig, err := c.addonConfig()
	if err != nil {
		return "", err
	}
	return configMap("tectonic-system", genericData{
		"addon-config":   addonConfig,
		"utility-config": utilityConfig,
	})
}

func (c ConfigGenerator) addonConfig() (*kubeaddon.OperatorConfig, error) {
	addonConfig := kubeaddon.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubeaddon.APIVersion,
			// TODO: get Kind from kubeaddon.Kind. Operator doesn't like KubeAddonOperatorConfig
			Kind: "AddonConfig",
		},
	}
	cidrhost, err := cidrhost(c.Cluster.Networking.ServiceCIDR, 10)
	if err != nil {
		return nil, err
	}
	addonConfig.DNSConfig.ClusterIP = cidrhost
	addonConfig.CloudProvider = c.Platform
	return &addonConfig, nil
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

	coreConfig.ClusterConfig.APIServerURL = fmt.Sprintf("https://%s-api.%s:443", c.Cluster.Name, c.Cluster.DNS.BaseDomain)
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
	networkConfig.NetworkProfile = tectonicnetwork.NetworkType(c.Cluster.Networking.Type)

	return &networkConfig
}

func (c ConfigGenerator) utilityConfig() (*tectonicutility.OperatorConfig, error) {
	utilityConfig := tectonicutility.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: tectonicutility.APIVersion,
			// TODO: get Kind from tectonicutility.Kind. Operator doesn't like TectonicUtilityOperatorConfig
			Kind: "UtilityConfig",
		},
	}

	var err error
	bytes, err := bcrypt.GenerateFromPassword([]byte(c.Console.AdminPassword), 12)
	if err != nil {
		return nil, err
	}
	hashedAdminPassword := string(bytes)
	adminUserID, err := generateRandomID(16)
	if err != nil {
		return nil, err
	}
	consoleSecret, err := generateRandomID(16)
	if err != nil {
		return nil, err
	}
	KubectlSecret, err := generateRandomID(16)
	if err != nil {
		return nil, err
	}
	clusterID, err := generateClusterID(16)
	if err != nil {
		return nil, err
	}
	utilityConfig.IdentityConfig.AdminEmail = c.Console.AdminEmail
	utilityConfig.IdentityConfig.AdminPasswordHash = hashedAdminPassword
	utilityConfig.IdentityConfig.AdminUserID = adminUserID
	utilityConfig.IdentityConfig.ConsoleClientID = "tectonic-console"
	utilityConfig.IdentityConfig.ConsoleSecret = consoleSecret
	utilityConfig.IdentityConfig.KubectlClientID = "tectonic-kubectl"
	utilityConfig.IdentityConfig.KubectlSecret = KubectlSecret

	utilityConfig.IngressConfig.ConsoleBaseHost = fmt.Sprintf("%s.%s", c.Cluster.Name, c.Cluster.DNS.BaseDomain)
	utilityConfig.IngressConfig.IngressKind = "NodePort"

	utilityConfig.StatsEmitterConfig.StatsURL = "https://stats-collector.tectonic.com"

	utilityConfig.TectonicConfigMapConfig.BaseAddress = fmt.Sprintf("%s.%s", c.Cluster.Name, c.Cluster.DNS.BaseDomain)
	utilityConfig.TectonicConfigMapConfig.CertificatesStrategy = "userProvidedCA"
	// TODO: Consolidate ClusterID with the one genereated by terraform and and passed to the bootstrap step.
	utilityConfig.TectonicConfigMapConfig.ClusterID = clusterID
	utilityConfig.TectonicConfigMapConfig.ClusterName = c.Cluster.Name
	utilityConfig.TectonicConfigMapConfig.IdentityAPIService = "tectonic-identity-api.tectonic-system.svc.cluster.local"
	utilityConfig.TectonicConfigMapConfig.InstallerPlatform = c.Cluster.Platform
	utilityConfig.TectonicConfigMapConfig.KubeAPIServerURL = fmt.Sprintf("https://%s-api.%s:443", c.Cluster.Name, c.Cluster.DNS.BaseDomain)
	// TODO: Speficy what's a version in ut2 and set it here
	utilityConfig.TectonicConfigMapConfig.TectonicVersion = "ut2"

	return &utilityConfig, nil
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

// generateClusterID reproduce tf cluster_id behaviour
// https://github.com/coreos/tectonic-installer/blob/master/modules/tectonic/assets.tf#L81
// TODO: re-evaluate solution
func generateClusterID(byteLength int) (string, error) {
	randomID, err := generateRandomID(byteLength)
	if err != nil {
		return "", err
	}
	bytes, err := base64.RawURLEncoding.DecodeString(randomID)
	hexStr := hex.EncodeToString(bytes)
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexStr[0:8],
		hexStr[8:12],
		hexStr[12:16],
		hexStr[16:20],
		hexStr[20:32]), nil
}

// cidrhost takes an IP address range in CIDR notation
// and creates an IP address with the given host number.
// If given host number is negative, the count starts from the end of the range
// Reproduces tf behaviour.
// TODO: re-evaluate solution
func cidrhost(iprange string, hostNum int) (string, error) {
	_, network, err := net.ParseCIDR(iprange)
	if err != nil {
		return "", fmt.Errorf("invalid CIDR expression: %s", err)
	}

	ip, err := cidr.Host(network, hostNum)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
