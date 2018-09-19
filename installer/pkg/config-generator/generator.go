package configgenerator

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/coreos/tectonic-config/config/kube-addon"
	"github.com/coreos/tectonic-config/config/kube-core"
	"github.com/coreos/tectonic-config/config/tectonic-network"
	"github.com/coreos/tectonic-config/config/tectonic-utility"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/config"
)

const (
	authConfigOIDCClientID        = "tectonic-kubectl"
	authConfigOIDCGroupsClaim     = "groups"
	authConfigOIDCUsernameClaim   = "email"
	networkConfigAdvertiseAddress = "0.0.0.0"
	identityConfigConsoleClientID = "tectonic-console"
	identityConfigKubectlClientID = "tectonic-kubectl"
	ingressConfigIngressKind      = "haproxy-router"
	certificatesStrategy          = "userProvidedCA"
	identityAPIService            = "tectonic-identity-api.tectonic-system.svc.cluster.local"
	maoTargetNamespace            = "openshift-cluster-api"
	libvirtPKIPath                = "/etc/pki/libvirt"
)

var errBadLibvirtScheme = errors.New("bad libvirt URI scheme")

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

// maoOperatorConfig contains configuration for mao managed stack
// TODO(enxebre): move up to "github.com/coreos/tectonic-config
type maoOperatorConfig struct {
	metav1.TypeMeta `json:",inline"`
	TargetNamespace string         `json:"targetNamespace"`
	APIServiceCA    string         `json:"apiServiceCA"`
	Provider        string         `json:"provider"`
	AWS             *awsConfig     `json:"aws"`
	Libvirt         *libvirtConfig `json:"libvirt"`
}

type libvirtConfig struct {
	ClusterName string `json:"clusterName"`
	URI         string `json:"uri"`
	NetworkName string `json:"networkName"`
	IPRange     string `json:"iprange"`
	Replicas    int    `json:"replicas"`
}

type awsConfig struct {
	ClusterName      string `json:"clusterName"`
	ClusterID        string `json:"clusterID"`
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone"`
	Image            string `json:"image"`
	Replicas         int    `json:"replicas"`
}

func (c *ConfigGenerator) maoConfig(clusterDir string) (*maoOperatorConfig, error) {
	cfg := maoOperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "machineAPIOperatorConfig",
		},

		TargetNamespace: maoTargetNamespace,
	}

	ca, err := ioutil.ReadFile(filepath.Join(clusterDir, aggregatorCACertPath))
	if err != nil {
		return nil, fmt.Errorf("could not read aggregator CA: %v", err)
	}

	cfg.APIServiceCA = string(ca)
	cfg.Provider = tectonicCloudProvider(c.Platform)

	switch c.Platform {
	case config.PlatformAWS:
		var ami string

		if c.AWS.EC2AMIOverride != "" {
			ami = c.AWS.EC2AMIOverride
		} else {
			ami, err = rhcos.AMI(context.TODO(), rhcos.DefaultChannel, c.Region)
			if err != nil {
				return nil, fmt.Errorf("failed to lookup RHCOS AMI: %v", err)
			}
		}

		cfg.AWS = &awsConfig{
			ClusterName:      c.Name,
			ClusterID:        c.ClusterID,
			Region:           c.Region,
			AvailabilityZone: "",
			Image:            ami,
			Replicas:         c.NodeCount(c.Worker.NodePools),
		}

	case config.PlatformLibvirt:
		uri, err := libvirtURI(c.Libvirt.URI, c.Libvirt.IPRange)
		if err != nil {
			return nil, fmt.Errorf("failed to create libvirt URI: %v", err)
		}

		cfg.Libvirt = &libvirtConfig{
			URI:         uri.String(),
			ClusterName: c.Name,
			NetworkName: c.Libvirt.Network.Name,
			IPRange:     c.Libvirt.IPRange,
			Replicas:    c.NodeCount(c.Worker.NodePools),
		}

	default:
		return nil, fmt.Errorf("unknown provider for machine-api-operator: %v", cfg.Provider)
	}

	return &cfg, nil
}

// KubeSystem returns, if successful, a yaml string for the kube-system.
func (c *ConfigGenerator) KubeSystem(clusterDir string) (string, error) {
	coreConfig, err := c.coreConfig()
	if err != nil {
		return "", err
	}
	installConfig, err := c.installConfig()
	if err != nil {
		return "", err
	}
	maoConfig, err := c.maoConfig(clusterDir)
	if err != nil {
		return "", err
	}

	return configMap("kube-system", genericData{
		"kco-config":     coreConfig,
		"network-config": c.networkConfig(),
		"install-config": installConfig,
		"mao-config":     maoConfig,
	})
}

// TectonicSystem returns, if successful, a yaml string for the tectonic-system.
func (c *ConfigGenerator) TectonicSystem() (string, error) {
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

// InstallConfig returns a YAML-rendered Kubernetes object with the user-supplied cluster configuration.
func (c *ConfigGenerator) InstallConfig() (string, error) {
	ic, err := c.installConfig()
	if err != nil {
		return "", err
	}
	return marshalYAML(ic)
}

func (c *ConfigGenerator) installConfig() (*types.InstallConfig, error) {
	_, podCIDR, err := net.ParseCIDR(c.Networking.PodCIDR)
	if err != nil {
		return nil, err
	}
	_, serviceCIDR, err := net.ParseCIDR(c.Networking.ServiceCIDR)
	if err != nil {
		return nil, err
	}

	var (
		platform       types.Platform
		masterPlatform types.MachinePoolPlatform
		workerPlatform types.MachinePoolPlatform
	)
	switch c.Platform {
	case config.PlatformAWS:
		platform.AWS = &types.AWSPlatform{
			Region:       c.Region,
			VPCID:        c.VPCID,
			VPCCIDRBlock: c.VPCCIDRBlock,
		}
		masterPlatform.AWS = &types.AWSMachinePoolPlatform{
			InstanceType: c.AWS.Master.EC2Type,
			IAMRoleName:  c.AWS.Master.IAMRoleName,
			EC2RootVolume: types.EC2RootVolume{
				IOPS: c.AWS.Master.MasterRootVolume.IOPS,
				Size: c.AWS.Master.MasterRootVolume.Size,
				Type: c.AWS.Master.MasterRootVolume.Type,
			},
		}
		workerPlatform.AWS = &types.AWSMachinePoolPlatform{
			InstanceType: c.AWS.Worker.EC2Type,
			IAMRoleName:  c.AWS.Worker.IAMRoleName,
			EC2RootVolume: types.EC2RootVolume{
				IOPS: c.AWS.Worker.WorkerRootVolume.IOPS,
				Size: c.AWS.Worker.WorkerRootVolume.Size,
				Type: c.AWS.Worker.WorkerRootVolume.Type,
			},
		}
	case config.PlatformLibvirt:
		platform.Libvirt = &types.LibvirtPlatform{
			URI: c.URI,
			Network: types.LibvirtNetwork{
				Name:    c.Network.Name,
				IfName:  c.Network.IfName,
				IPRange: c.Network.IPRange,
			},
		}
		masterPlatform.Libvirt = &types.LibvirtMachinePoolPlatform{
			ImagePool:   "default",
			ImageVolume: "coreos_base",
		}
		workerPlatform.Libvirt = &types.LibvirtMachinePoolPlatform{
			ImagePool:   "default",
			ImageVolume: "coreos_base",
		}
	default:
		return nil, fmt.Errorf("installconfig: invalid platform %s", c.Platform)
	}
	masterCount := int64(c.NodeCount(c.Master.NodePools))
	workerCount := int64(c.NodeCount(c.Worker.NodePools))

	return &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.Name,
		},
		ClusterID: c.ClusterID,
		Admin: types.Admin{
			Email:    c.Admin.Email,
			Password: c.Admin.Password,
			SSHKey:   c.Admin.SSHKey,
		},
		BaseDomain: c.BaseDomain,
		PullSecret: c.PullSecret,
		Networking: types.Networking{
			Type:        types.NetworkType(string(c.Networking.Type)),
			ServiceCIDR: ipnet.IPNet{IPNet: *serviceCIDR},
			PodCIDR:     ipnet.IPNet{IPNet: *podCIDR},
		},
		Platform: platform,
		Machines: []types.MachinePool{{
			Name:     "master",
			Replicas: &masterCount,
			Platform: masterPlatform,
		}, {
			Name:     "worker",
			Replicas: &workerCount,
			Platform: workerPlatform,
		}},
	}, nil
}

// CoreConfig returns, if successful, a yaml string for the on-disk kco-config.
func (c *ConfigGenerator) CoreConfig() (string, error) {
	coreConfig, err := c.coreConfig()
	if err != nil {
		return "", err
	}
	return marshalYAML(coreConfig)
}

func (c *ConfigGenerator) addonConfig() (*kubeaddon.OperatorConfig, error) {
	addonConfig := kubeaddon.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubeaddon.APIVersion,
			Kind:       kubeaddon.Kind,
		},
	}
	addonConfig.CloudProvider = tectonicCloudProvider(c.Platform)
	addonConfig.ClusterConfig.APIServerURL = c.getAPIServerURL()
	registrySecret, err := generateRandomID(16)
	if err != nil {
		return nil, err
	}
	addonConfig.RegistryHTTPSecret = registrySecret
	return &addonConfig, nil
}

func (c *ConfigGenerator) coreConfig() (*kubecore.OperatorConfig, error) {
	coreConfig := kubecore.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: kubecore.APIVersion,
			Kind:       kubecore.Kind,
		},
	}
	coreConfig.ClusterConfig.APIServerURL = c.getAPIServerURL()
	coreConfig.AuthConfig.OIDCClientID = authConfigOIDCClientID
	coreConfig.AuthConfig.OIDCIssuerURL = c.getOicdIssuerURL()
	coreConfig.AuthConfig.OIDCGroupsClaim = authConfigOIDCGroupsClaim
	coreConfig.AuthConfig.OIDCUsernameClaim = authConfigOIDCUsernameClaim

	cidrhost, err := cidrhost(c.Cluster.Networking.ServiceCIDR, 10)
	if err != nil {
		return nil, err
	}
	coreConfig.DNSConfig.ClusterIP = cidrhost

	coreConfig.CloudProviderConfig.CloudConfigPath = ""
	coreConfig.CloudProviderConfig.CloudProviderProfile = k8sCloudProvider(c.Cluster.Platform)

	coreConfig.RoutingConfig.Subdomain = c.getBaseAddress()

	coreConfig.NetworkConfig.ClusterCIDR = c.Cluster.Networking.PodCIDR
	coreConfig.NetworkConfig.ServiceCIDR = c.Cluster.Networking.ServiceCIDR
	coreConfig.NetworkConfig.AdvertiseAddress = networkConfigAdvertiseAddress
	coreConfig.NetworkConfig.EtcdServers = c.getEtcdServersURLs()

	return &coreConfig, nil
}

func (c *ConfigGenerator) networkConfig() *tectonicnetwork.OperatorConfig {
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

func (c *ConfigGenerator) utilityConfig() (*tectonicutility.OperatorConfig, error) {
	utilityConfig := tectonicutility.OperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: tectonicutility.APIVersion,
			Kind:       tectonicutility.Kind,
		},
	}

	utilityConfig.TectonicConfigMapConfig.CertificatesStrategy = certificatesStrategy
	utilityConfig.TectonicConfigMapConfig.ClusterID = c.Cluster.Internal.ClusterID
	utilityConfig.TectonicConfigMapConfig.ClusterName = c.Cluster.Name
	utilityConfig.TectonicConfigMapConfig.InstallerPlatform = tectonicCloudProvider(c.Platform)
	utilityConfig.TectonicConfigMapConfig.KubeAPIServerURL = c.getAPIServerURL()
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

func (c *ConfigGenerator) getEtcdServersURLs() string {
	etcdServers := make([]string, c.Cluster.NodeCount(c.Cluster.Master.NodePools))
	for i := range etcdServers {
		etcdServers[i] = fmt.Sprintf("https://%s-etcd-%v.%s:2379", c.Cluster.Name, i, c.Cluster.BaseDomain)
	}
	return strings.Join(etcdServers, ",")
}

func (c *ConfigGenerator) getAPIServerURL() string {
	return fmt.Sprintf("https://%s-api.%s:6443", c.Cluster.Name, c.Cluster.BaseDomain)
}

func (c *ConfigGenerator) getBaseAddress() string {
	return fmt.Sprintf("%s.%s", c.Cluster.Name, c.Cluster.BaseDomain)
}

func (c *ConfigGenerator) getOicdIssuerURL() string {
	return fmt.Sprintf("https://%s.%s/identity", c.Cluster.Name, c.Cluster.BaseDomain)
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

// GenerateClusterID reproduce tf cluster_id behaviour
// https://github.com/coreos/tectonic-installer/blob/master/modules/tectonic/assets.tf#L81
// TODO: re-evaluate solution
func GenerateClusterID(byteLength int) (string, error) {
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
		return "", fmt.Errorf("invalid CIDR expression (%s): %s", iprange, err)
	}

	ip, err := cidr.Host(network, hostNum)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}

// Converts a platform to the cloudProvider that k8s understands
func k8sCloudProvider(platform config.Platform) string {
	switch platform {
	case config.PlatformAWS:
		return "aws"
	case config.PlatformLibvirt:
		return ""
	}
	panic("invalid platform")
}

// Converts a platform to the cloudProvider that Tectonic understands
func tectonicCloudProvider(platform config.Platform) string {
	switch platform {
	case config.PlatformAWS:
		return "aws"
	case config.PlatformLibvirt:
		return "libvirt"
	}
	panic("invalid platform")
}

// Returns a libvirt URI given for the machine-api-operator given the
// URI in the config and the cluster network in CIDR format.
func libvirtURI(configURI, networkCIDR string) (*url.URL, error) {
	_, ipNet, err := net.ParseCIDR(networkCIDR)
	if err != nil {
		return nil, err
	}

	gateway, err := cidr.Host(ipNet, 1)
	if err != nil {
		return nil, err
	}

	libvirtURI, err := url.Parse(configURI)
	if err != nil {
		return nil, err
	}

	// If there's a transport in the configured URI, replace it with
	// TLS.  Otherwise, if there's none, explicityly add TLS.
	scheme := strings.Split(libvirtURI.Scheme, "+")
	switch len(scheme) {
	case 1, 2: // Replace or add the transport.
		libvirtURI.Scheme = scheme[0] + "+tls"

	default:
		return nil, errBadLibvirtScheme
	}

	query := libvirtURI.Query()
	query.Set("pkipath", libvirtPKIPath)

	libvirtURI.Host = gateway.String()
	libvirtURI.RawQuery = query.Encode()

	return libvirtURI, nil
}
