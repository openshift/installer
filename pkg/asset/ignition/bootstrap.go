package ignition

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/coreos/ignition/config/util"
	ignition "github.com/coreos/ignition/config/v2_2/types"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/content"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
)

const (
	// tlsCertDirectory is the directory on the bootstrap node to place the TLS
	// assets.
	tlsCertDirectory = "/opt/tectonic/tls"
)

// bootstrapTemplateData is the data to use to replace values in bootstrap
// template files.
type bootstrapTemplateData struct {
	ClusterDNSIP               string
	CloudProvider              string
	CloudProviderConfig        string
	DebugConfig                string
	KubeCoreRenderImage        string
	MachineConfigOperatorImage string
	EtcdCertSignerImage        string
	EtcdctlImage               string
	BootkubeImage              string
	HyperkubeImage             string
	EtcdCluster                string
}

// bootstrap is an asset that generates the ignition config for bootstrap nodes.
type bootstrap struct {
	directory                 string
	installConfig             asset.Asset
	rootCA                    asset.Asset
	etcdCA                    asset.Asset
	ingressCertKey            asset.Asset
	kubeCA                    asset.Asset
	aggregatorCA              asset.Asset
	serviceServingCA          asset.Asset
	clusterAPIServerCertKey   asset.Asset
	etcdClientCertKey         asset.Asset
	apiServerCertKey          asset.Asset
	openshiftAPIServerCertKey asset.Asset
	apiServerProxyCertKey     asset.Asset
	adminCertKey              asset.Asset
	kubeletCertKey            asset.Asset
	tncCertKey                asset.Asset
	serviceAccountKeyPair     asset.Asset
	kubeconfig                asset.Asset
	kubeconfigKubelet         asset.Asset
}

var _ asset.Asset = (*bootstrap)(nil)

// newBootstrap creates a new bootstrap asset.
func newBootstrap(
	directory string,
	installConfigStock installconfig.Stock,
	tlsStock tls.Stock,
	kubeconfigStock kubeconfig.Stock,
) *bootstrap {
	return &bootstrap{
		directory:                 directory,
		installConfig:             installConfigStock.InstallConfig(),
		rootCA:                    tlsStock.RootCA(),
		etcdCA:                    tlsStock.EtcdCA(),
		ingressCertKey:            tlsStock.IngressCertKey(),
		kubeCA:                    tlsStock.KubeCA(),
		aggregatorCA:              tlsStock.AggregatorCA(),
		serviceServingCA:          tlsStock.ServiceServingCA(),
		clusterAPIServerCertKey:   tlsStock.ClusterAPIServerCertKey(),
		etcdClientCertKey:         tlsStock.EtcdClientCertKey(),
		apiServerCertKey:          tlsStock.APIServerCertKey(),
		openshiftAPIServerCertKey: tlsStock.OpenshiftAPIServerCertKey(),
		apiServerProxyCertKey:     tlsStock.APIServerProxyCertKey(),
		adminCertKey:              tlsStock.AdminCertKey(),
		kubeletCertKey:            tlsStock.KubeletCertKey(),
		tncCertKey:                tlsStock.TNCCertKey(),
		serviceAccountKeyPair:     tlsStock.ServiceAccountKeyPair(),
		kubeconfig:                kubeconfigStock.KubeconfigAdmin(),
		kubeconfigKubelet:         kubeconfigStock.KubeconfigKubelet(),
	}
}

// Dependencies returns the assets on which the bootstrap asset depends.
func (a *bootstrap) Dependencies() []asset.Asset {
	return []asset.Asset{
		a.installConfig,
		a.rootCA,
		a.etcdCA,
		a.ingressCertKey,
		a.kubeCA,
		a.aggregatorCA,
		a.serviceServingCA,
		a.clusterAPIServerCertKey,
		a.etcdClientCertKey,
		a.apiServerCertKey,
		a.openshiftAPIServerCertKey,
		a.apiServerProxyCertKey,
		a.adminCertKey,
		a.kubeletCertKey,
		a.tncCertKey,
		a.serviceAccountKeyPair,
		a.kubeconfig,
		a.kubeconfigKubelet,
	}
}

// Generate generates the ignition config for the bootstrap asset.
func (a *bootstrap) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	installConfig, err := installconfig.GetInstallConfig(a.installConfig, dependencies)
	if err != nil {
		return nil, err
	}

	templateData, err := a.getTemplateData(installConfig)
	if err != nil {
		return nil, err
	}

	config := ignition.Config{}

	a.addBootstrapConfigFiles(&config, dependencies)
	a.addBootstrapCertFiles(&config, dependencies)
	a.addBootkubeFiles(&config, dependencies, templateData)
	a.addTectonicFiles(&config, dependencies, templateData)
	a.addTLSCertFiles(&config, dependencies)

	config.Systemd.Units = append(
		config.Systemd.Units,
		ignition.Unit{Name: "bootkube.service", Contents: content.BootkubeSystemdContents},
		ignition.Unit{Name: "tectonic.service", Contents: content.TectonicSystemdContents, Enabled: util.BoolToPtr(true)},
		ignition.Unit{Name: "kubelet.service", Contents: applyTemplateData(content.KubeletSystemdTemplate, templateData), Enabled: util.BoolToPtr(true)},
	)

	config.Passwd.Users = append(
		config.Passwd.Users,
		ignition.PasswdUser{Name: "core", SSHAuthorizedKeys: []ignition.SSHAuthorizedKey{ignition.SSHAuthorizedKey(installConfig.Admin.SSHKey)}},
	)

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{{
			Name: filepath.Join(a.directory, "bootstrap.ign"),
			Data: data,
		}},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (a *bootstrap) Name() string {
	return "Bootstrap Ignition Config"
}

// getTemplateData returns the data to use to execute bootstrap templates.
func (a *bootstrap) getTemplateData(installConfig *types.InstallConfig) (*bootstrapTemplateData, error) {
	clusterDNSIP, err := installconfig.ClusterDNSIP(installConfig)
	if err != nil {
		return nil, err
	}
	etcdEndpoints := make([]string, masterCount(installConfig))
	for i := range etcdEndpoints {
		etcdEndpoints[i] = fmt.Sprintf("https://%s-etcd-%d.%s:2379", installConfig.Name, i, installConfig.BaseDomain)
	}
	return &bootstrapTemplateData{
		ClusterDNSIP:               clusterDNSIP,
		CloudProvider:              getCloudProvider(installConfig),
		CloudProviderConfig:        getCloudProviderConfig(installConfig),
		DebugConfig:                "",
		KubeCoreRenderImage:        "quay.io/coreos/kube-core-renderer-dev:436b1b4395ae54d866edc88864c9b01797cebac1",
		MachineConfigOperatorImage: "docker.io/openshift/origin-machine-config-operator:v4.0.0",
		EtcdCertSignerImage:        "quay.io/coreos/kube-etcd-signer-server:678cc8e6841e2121ebfdb6e2db568fce290b67d6",
		EtcdctlImage:               "quay.io/coreos/etcd:v3.2.14",
		BootkubeImage:              "quay.io/coreos/bootkube:v0.10.0",
		HyperkubeImage:             "openshift/origin-node:latest",
		EtcdCluster:                strings.Join(etcdEndpoints, ","),
	}, nil
}

func (a *bootstrap) addBootstrapConfigFiles(config *ignition.Config, dependencies map[asset.Asset]*asset.State) {
	// TODO (staebler) - missing the following from assets step
	//     /opt/tectonic/manifests/cluster-config.yaml
	//     /opt/tectonic/tectonic/cluster-config.yaml
	//     /opt/tectonic/tnco-config.yaml
	//     /opt/tectonic/kco-config.yaml
	//     /etc/kubernetes/kubeconfig
	//     /var/lib/kubelet/kubeconfig
}

func (a *bootstrap) addBootstrapCertFiles(config *ignition.Config, dependencies map[asset.Asset]*asset.State) {
	config.Storage.Files = append(
		config.Storage.Files,
		fileFromAsset("/etc/ssl/etcd/ca.crt", 0444, dependencies[a.etcdCA], keyCertAssetCrtIndex),
		fileFromAsset("/etc/ssl/etcd/root-ca.crt", 0444, dependencies[a.rootCA], keyCertAssetCrtIndex),

		// ssl certs
		fileFromAsset("/etc/ssl/certs/root_ca.pem", 0444, dependencies[a.rootCA], keyCertAssetKeyIndex),
		fileFromAsset("/etc/ssl/certs/ingress_ca.pem", 0444, dependencies[a.ingressCertKey], keyCertAssetKeyIndex),
		fileFromAsset("/etc/ssl/certs/etcd_ca.pem", 0444, dependencies[a.etcdCA], keyCertAssetKeyIndex),
	)
}

func (a *bootstrap) addBootkubeFiles(config *ignition.Config, dependencies map[asset.Asset]*asset.State, templateData *bootstrapTemplateData) {
	// TODO (staebler) - missing manifests from bootkube module
	config.Storage.Files = append(
		config.Storage.Files,
		fileFromAsset("/opt/tectonic/auth/kubeconfig", 0400, dependencies[a.kubeconfig], 0),
		fileFromAsset("/opt/tectonic/auth/kubeconfig-kubelet", 0400, dependencies[a.kubeconfigKubelet], 0),
		fileFromString("/opt/tectonic/bootkube.sh", 0555, applyTemplateData(content.BootkubeShFileTemplate, templateData)),
	)
}

func (a *bootstrap) addTectonicFiles(config *ignition.Config, dependencies map[asset.Asset]*asset.State, templateData *bootstrapTemplateData) {
	// TODO (staebler) - missing manifests from tectonic module
	config.Storage.Files = append(
		config.Storage.Files,
		fileFromString("/opt/tectonic/tectonic.sh", 0555, content.TectonicShFileContents),
	)
}

func (a *bootstrap) addTLSCertFiles(config *ignition.Config, dependencies map[asset.Asset]*asset.State) {
	for _, pair := range []struct {
		key   string
		crt   string
		state *asset.State
	}{
		{"", "root-ca.crt", dependencies[a.rootCA]},
		{"kube-ca.key", "kube-ca.crt", dependencies[a.kubeCA]},
		{"aggregator-ca.key", "aggregator-ca.crt", dependencies[a.aggregatorCA]},
		{"service-serving-ca.key", "service-serving-ca.crt", dependencies[a.serviceServingCA]},
		{"etcd-client-ca.key", "etcd-client-ca.crt", dependencies[a.etcdCA]},
		{"cluster-apiserver-ca.key", "cluster-apiserver-ca.crt", dependencies[a.clusterAPIServerCertKey]},

		// etcd cert
		{"etcd-client.key", "etcd-client.crt", dependencies[a.etcdClientCertKey]},

		// kube certs
		{"apiserver.key", "apiserver.crt", dependencies[a.apiServerCertKey]},
		{"openshift-apiserver.key", "openshift-apiserver.crt", dependencies[a.openshiftAPIServerCertKey]},
		{"apiserver-proxy.key", "apiserver-proxy.crt", dependencies[a.apiServerProxyCertKey]},
		{"admin.key", "admin.crt", dependencies[a.adminCertKey]},
		{"kubelet.key", "kubelet.crt", dependencies[a.kubeletCertKey]},

		// tnc cert
		{"tnc.key", "tnc.crt", dependencies[a.tncCertKey]},

		// service account cert
		{"service-account.key", "service-account.crt", dependencies[a.serviceAccountKeyPair]},
	} {
		if pair.key != "" {
			config.Storage.Files = append(config.Storage.Files, fileFromAsset(path.Join(tlsCertDirectory, pair.key), 0600, pair.state, keyCertAssetKeyIndex))
		}
		if pair.crt != "" {
			config.Storage.Files = append(config.Storage.Files, fileFromAsset(path.Join(tlsCertDirectory, pair.crt), 0644, pair.state, keyCertAssetCrtIndex))
		}
	}
}

func getCloudProvider(installConfig *types.InstallConfig) string {
	if installConfig.AWS != nil {
		return "aws"
	}
	return ""
}

func getCloudProviderConfig(installConfig *types.InstallConfig) string {
	return ""
}

func applyTemplateData(template *template.Template, templateData interface{}) string {
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.String()
}
