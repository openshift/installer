package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"

	"github.com/coreos/matchbox/matchbox/storage/storagepb"
	bootkube "github.com/kubernetes-incubator/bootkube/pkg/asset"
	bootkubeTLS "github.com/kubernetes-incubator/bootkube/pkg/tlsutil"
	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/binassets"
	"github.com/coreos/tectonic-installer/installer/server/asset"
	"github.com/coreos/tectonic-installer/installer/server/defaults"
)

var (
	// TectonicMetalCluster defaults
	metalDefaults = struct {
		installIgnTmpl    []byte
		controllerIgnTmpl []byte
		workerIgnTmpl     []byte
	}{
		installIgnTmpl:    binassets.MustAsset("provisioning/install-reboot.yaml.tmpl"),
		controllerIgnTmpl: binassets.MustAsset("provisioning/tectonic-metal-controller.yaml.tmpl"),
		workerIgnTmpl:     binassets.MustAsset("provisioning/tectonic-metal-worker.yaml.tmpl"),
	}
)

// TectonicMetalCluster provisions a Tectonic self-hosted Kuberntes cluster on
// physical machines (bare metal).
type TectonicMetalCluster struct {
	// Matchbox HTTP name/IP and port
	MatchboxHTTP string `json:"matchboxHTTP"`
	// Matchbox gRPC API name/IP and port
	MatchboxRPC string `json:"matchboxRPC"`
	// Matchbox certificate authority for verifying the server's identity
	MatchboxCA string `json:"matchboxCA"`
	// Matchbox client certificate and key for authentication
	MatchboxClientCert string `json:"matchboxClientCert"`
	MatchboxClientKey  string `json:"matchboxClientKey"`

	// CoreOS PXE and install channel/version
	Channel string `json:"channel"`
	Version string `json:"version"`

	// External etcd client endpoint, e.g. etcd.example.com:2379
	ExternalETCDClient string `json:"externalETCDClient"`

	// Kubernetes Control Plane nodes
	ControllerDomain string `json:"controllerDomain"`
	Controllers      []Node `json:"controllers"`
	// Kuberntes Worker nodes
	Workers []Node `json:"workers"`
	// Admin SSH Public Keys
	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys"`

	// Custom Certificate Authority (optional)
	CACertificate string `json:"caCertificate"`
	CAPrivateKey  string `json:"caPrivateKey"`

	PodCIDR     string `json:"podCIDR"`
	ServiceCIDR string `json:"serviceCIDR"`

	// Computed IPs for self-hosted Kubernetes
	APIServiceIP net.IP
	DNSServiceIP net.IP

	// Tectonic
	TectonicDomain string          `json:"tectonicDomain"`
	Tectonic       *TectonicConfig `json:"tectonic"`

	// Generated
	kubeconfig *kubeConfig

	// Assets for bootkube
	assets []asset.Asset
}

// Initialize validates cluster data and sets defaults.
func (c *TectonicMetalCluster) Initialize() error {
	if len(c.Controllers) < 1 {
		return errTooFewControllers
	}
	if c.size() < 2 {
		return errClusterTooSmall
	}
	if c.MatchboxHTTP == "" {
		return errMissingMatchboxEndpoint
	}

	if c.ControllerDomain == "" {
		return errMissingControllerDomain
	}
	if c.TectonicDomain == "" {
		return errMissingTectonicDomain
	}

	for _, node := range c.nodes() {
		if node.MAC == nil {
			return errMissingMACAddress
		}
	}

	for _, node := range c.nodes() {
		if node.Name == "" {
			return errMissingNodeName
		}
	}

	if c.Channel == "" {
		return errMissingChannel
	}

	if c.ExternalETCDClient != "" {
		if _, err := url.Parse(c.ExternalETCDClient); err != nil {
			return errInvalidExternalETCD
		}
	}

	if c.Version == "" {
		return errMissingVersion
	}

	if c.Tectonic == nil {
		return errMissingTectonicConfig
	}

	// Kubernetes CIDR customization
	if c.PodCIDR == "" {
		c.PodCIDR = defaults.PodCIDR
	}
	if c.ServiceCIDR == "" {
		c.ServiceCIDR = defaults.ServiceCIDR
	}

	var err error
	c.APIServiceIP, err = defaults.APIServiceIP(c.ServiceCIDR)
	if err != nil {
		return err
	}
	c.DNSServiceIP, err = defaults.DNSServiceIP(c.ServiceCIDR)
	if err != nil {
		return err
	}

	// Tectonic asset pipeline requires both domains
	c.Tectonic.ControllerDomain = c.ControllerDomain
	c.Tectonic.TectonicDomain = c.TectonicDomain
	return c.Tectonic.AssertValid()
}

// GenerateAssets generates cluster assets.
func (c *TectonicMetalCluster) GenerateAssets() ([]asset.Asset, error) {
	config, err := c.getBootkubeConfig()
	if err != nil {
		return nil, err
	}

	// Self-hosted (kube-system) manifests
	assets, err := NewBootkubeAssets(config)
	if err != nil {
		return nil, err
	}

	// Tectonic (tectonic-system) manifests and add-ons
	m := metrics{
		certificatesStrategy:   getCertificatesStrategy(c.CACertificate),
		installerPlatform:      c.Kind(),
		tectonicUpdaterEnabled: c.Tectonic.Updater.Enabled,
	}
	assets, err = NewTectonicAssets(assets, c.Tectonic, m)
	if err != nil {
		return nil, err
	}

	// Read generated credentials to inject into host metadata
	kubecfg, err := parseKubeConfig(assets)
	if err != nil {
		return nil, err
	}

	// Store generated kubeconfig
	c.kubeconfig = kubecfg

	// Save assets to apply to template in Publish()
	c.assets = assets

	return assets, err
}

// StatusChecker returns a StatusChecker for Tectonic metal clusters.
func (c *TectonicMetalCluster) StatusChecker() (StatusChecker, error) {
	return TectonicMetalChecker{
		Controllers:    c.Controllers,
		Workers:        c.Workers,
		TectonicDomain: c.Tectonic.TectonicDomain,
	}, nil
}

// Kind returns the kind name.
func (c *TectonicMetalCluster) Kind() string {
	return "tectonic-metal"
}

// Publish writes profiles, groups, and Ignition to a matchbox service.
func (c *TectonicMetalCluster) Publish(ctx context.Context) error {
	client, err := NewMatchboxClient(&MatchboxConfig{
		Endpoint:   c.MatchboxRPC,
		CA:         []byte(c.MatchboxCA),
		ClientCert: []byte(c.MatchboxClientCert),
		ClientKey:  []byte(c.MatchboxClientKey),
	})
	if err != nil {
		return err
	}
	defer client.Close()

	groups, err := c.groups()
	if err != nil {
		return err
	}
	if err = client.Push(ctx, groups, c.profiles(), c.ignitionTemplates()); err != nil {
		return err
	}

	return nil
}

// groups returns a machine group for each cluster node.
func (c *TectonicMetalCluster) groups() (groups []*storagepb.Group, err error) {
	// Match to CoreOS install-reboot Profile
	for _, node := range c.nodes() {
		install := &storagepb.Group{
			Id:      fmt.Sprintf("tectonic-install-%s", node.MAC.DashString()),
			Name:    "CoreOS Install",
			Profile: "install-reboot",
			Selector: map[string]string{
				"mac": node.MAC.String(),
			},
		}
		install.Metadata, err = c.installMetadata()
		if err != nil {
			return nil, err
		}
		groups = append(groups, install)
	}

	// Match controllers to tectonic-controller Profile
	for i, controller := range c.Controllers {
		group := &storagepb.Group{
			Id:      fmt.Sprintf("tectonic-node-%s", controller.MAC.DashString()),
			Profile: "tectonic-controller",
			Selector: map[string]string{
				"mac": controller.MAC.String(),
				"os":  "installed",
			},
		}
		group.Metadata, err = c.controllerMetadata(i)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	// Match workers to tectonic-worker Profile
	for i, worker := range c.Workers {
		group := &storagepb.Group{
			Id:      fmt.Sprintf("tectonic-node-%s", worker.MAC.DashString()),
			Profile: "tectonic-worker",
			Selector: map[string]string{
				"mac": worker.MAC.String(),
				"os":  "installed",
			},
		}
		group.Metadata, err = c.workerMetadata(i)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// profiles returns the install, controller, and worker profiles.
func (c *TectonicMetalCluster) profiles() []*storagepb.Profile {
	install := c.installProfile()
	controller := controllerProfile()
	worker := workerProfile()
	return []*storagepb.Profile{install, controller, worker}
}

// ignitionTemplates returns the cluster Ignition templates.
func (c TectonicMetalCluster) ignitionTemplates() []asset.Asset {
	return []asset.Asset{
		asset.New("tectonic-controller.yaml.tmpl", metalDefaults.controllerIgnTmpl),
		asset.New("tectonic-worker.yaml.tmpl", metalDefaults.workerIgnTmpl),
		asset.New("install-reboot.yaml.tmpl", metalDefaults.installIgnTmpl),
	}
}

func (c *TectonicMetalCluster) installProfile() *storagepb.Profile {
	return &storagepb.Profile{
		Id:         "install-reboot",
		Name:       "CoreOS Install and Reboot",
		IgnitionId: "install-reboot.yaml.tmpl",
		Boot: &storagepb.NetBoot{
			Kernel: fmt.Sprintf("/assets/coreos/%s/coreos_production_pxe.vmlinuz", c.Version),
			Initrd: []string{
				fmt.Sprintf("/assets/coreos/%s/coreos_production_pxe_image.cpio.gz", c.Version),
			},
			Cmdline: map[string]string{
				"coreos.first_boot": "1",
				"coreos.config.url": ignitionPath(c.MatchboxHTTP),
			},
		},
	}
}

func controllerProfile() *storagepb.Profile {
	return &storagepb.Profile{
		Id:         "tectonic-controller",
		Name:       "Tectonic Controller",
		IgnitionId: "tectonic-controller.yaml.tmpl",
		Boot:       &storagepb.NetBoot{},
	}
}

func workerProfile() *storagepb.Profile {
	return &storagepb.Profile{
		Id:         "tectonic-worker",
		Name:       "Tectonic Worker",
		IgnitionId: "tectonic-worker.yaml.tmpl",
		Boot:       &storagepb.NetBoot{},
	}
}

// size returns the total number of nodes in the cluster.
func (c *TectonicMetalCluster) size() int {
	return len(c.Controllers) + len(c.Workers)
}

// nodes returns the slice of the nodes in the cluster.
func (c *TectonicMetalCluster) nodes() []Node {
	return append(c.Controllers[:], c.Workers[:]...)
}

// installMetadata returns the group metadata for installing CoreOS.
func (c *TectonicMetalCluster) installMetadata() ([]byte, error) {
	data := map[string]interface{}{
		"coreos_channel":      c.Channel,
		"coreos_version":      c.Version,
		"ignition_endpoint":   fmt.Sprintf("http://%s/ignition", c.MatchboxHTTP),
		"baseurl":             coreosAssetsPath(c.MatchboxHTTP),
		"ssh_authorized_keys": c.SSHAuthorizedKeys,
	}
	return json.Marshal(data)
}

// controllerMetadata returns the group metadata for controller nodes. Requires
// validated cluster data.
func (c *TectonicMetalCluster) controllerMetadata(i int) ([]byte, error) {
	etcdInitialCluster := etcdInitialCluster(c.Controllers)
	etcdEndpoints := etcdEndpoints(c.Controllers)
	if c.ExternalETCDClient != "" {
		etcdEndpoints = c.ExternalETCDClient
	}

	// Save assets in ignition template
	if c.assets == nil {
		return nil, fmt.Errorf("Failed to determine cluster assets")
	}
	buf, err := asset.ZipAssets(c.assets)
	if err != nil {
		return nil, fmt.Errorf("Failed to zip assets: %v", err)
	}

	node := c.Controllers[i]
	data := map[string]interface{}{
		"domain_name":               node.Name,
		"etcd_endpoints":            etcdEndpoints,
		"etcd_initial_cluster":      etcdInitialCluster,
		"etcd_name":                 etcdNodeName(i),
		"external_etcd":             c.ExternalETCDClient != "",
		"k8s_controller_endpoint":   fmt.Sprintf("https://%s:443", c.ControllerDomain),
		"k8s_dns_service_ip":        c.DNSServiceIP,
		"k8s_pod_network":           c.PodCIDR,
		"k8s_service_ip_range":      c.ServiceCIDR,
		"k8s_certificate_authority": c.kubeconfig.certificateAuthority,
		"k8s_client_certificate":    c.kubeconfig.clientCertificate,
		"k8s_client_key":            c.kubeconfig.clientKey,
		"ssh_authorized_keys":       c.SSHAuthorizedKeys,
		"zippedAssets":              base64.StdEncoding.EncodeToString(buf),
	}
	return json.Marshal(data)
}

// workerMetadata returns the group metadata for worker nodes. Requires
// validated cluster data.
func (c *TectonicMetalCluster) workerMetadata(i int) ([]byte, error) {
	etcdEndpoints := etcdEndpoints(c.Controllers)
	if c.ExternalETCDClient != "" {
		etcdEndpoints = c.ExternalETCDClient
	}

	node := c.Workers[i]
	data := map[string]interface{}{
		"domain_name":               node.Name,
		"etcd_endpoints":            etcdEndpoints,
		"external_etcd":             c.ExternalETCDClient != "",
		"k8s_controller_endpoint":   fmt.Sprintf("https://%s:443", c.ControllerDomain),
		"k8s_dns_service_ip":        c.DNSServiceIP,
		"k8s_pod_network":           c.PodCIDR,
		"k8s_service_ip_range":      c.ServiceCIDR,
		"k8s_certificate_authority": c.kubeconfig.certificateAuthority,
		"k8s_client_certificate":    c.kubeconfig.clientCertificate,
		"k8s_client_key":            c.kubeconfig.clientKey,
		"ssh_authorized_keys":       c.SSHAuthorizedKeys,
	}
	return json.Marshal(data)
}

// getBootkubeConfig returns a BootkubeConfig for asset generation.
func (c *TectonicMetalCluster) getBootkubeConfig() (BootkubeConfig, error) {
	// etcd client net.URLs (e.g. http://etcd.example.com:2379)
	etcds, err := etcdURLs(c.Controllers)
	if err != nil {
		return BootkubeConfig{}, err
	}
	if c.ExternalETCDClient != "" {
		etcd, err := url.Parse(c.ExternalETCDClient)
		if err != nil {
			return BootkubeConfig{}, err
		}
		etcds = []*url.URL{etcd}
	}

	// kube-apiserver net.URLs (e.g. https://cluster.example.com:443)
	controllerURL, err := url.Parse(fmt.Sprintf("https://%s:443", c.ControllerDomain))
	if err != nil {
		return BootkubeConfig{}, err
	}

	_, podNet, err := net.ParseCIDR(c.PodCIDR)
	if err != nil {
		return BootkubeConfig{}, err
	}

	_, serviceNet, err := net.ParseCIDR(c.ServiceCIDR)
	if err != nil {
		return BootkubeConfig{}, err
	}

	// (optional) user-provided certificate authority
	caCert, err := getCACertificate(c.CACertificate)
	if err != nil {
		return BootkubeConfig{}, err
	}
	caPrivateKey, err := getCAPrivateKey(c.CAPrivateKey)
	if err != nil {
		return BootkubeConfig{}, err
	}

	// Configure bootkube asset rendering
	return BootkubeConfig{
		Config: bootkube.Config{
			EtcdServers: etcds,
			CACert:      caCert,
			CAPrivKey:   caPrivateKey,
			APIServers:  []*url.URL{controllerURL},
			AltNames: &bootkubeTLS.AltNames{
				DNSNames: []string{c.ControllerDomain},
			},
			PodCIDR:      podNet,
			ServiceCIDR:  serviceNet,
			APIServiceIP: c.APIServiceIP,
			DNSServiceIP: c.DNSServiceIP,
		},
		OIDCIssuer: &OIDCIssuer{
			IssuerURL:     fmt.Sprintf("https://%s/identity", c.Tectonic.TectonicDomain),
			ClientID:      oidcKubectlClientID,
			UsernameClaim: oidcUsernameClaim,
			// CACert is written as a secret and mounted by the kube-apiserver
			CAPath: "/etc/kubernetes/secrets/ca.crt",
		},
	}, nil
}
