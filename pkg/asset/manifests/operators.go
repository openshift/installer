// Package manifests deals with creating manifests for all manifests to be installed for the cluster
package manifests

import (
	"bytes"
	"encoding/base64"
	"path/filepath"
	"text/template"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/content/bootkube"
)

const (
	keyIndex    = 0
	certIndex   = 1
	manifestDir = "manifests"
)

// manifests generates the dependent operator config.yaml files
type manifests struct {
	assetStock                Stock
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
	kubeletCertKey            asset.Asset
	mcsCertKey                asset.Asset
	serviceAccountKeyPair     asset.Asset
	kubeconfig                asset.Asset
	workerIgnition            asset.Asset
}

var _ asset.Asset = (*manifests)(nil)

type genericData map[string]string

// Name returns a human friendly name for the operator
func (m *manifests) Name() string {
	return "Common Manifests"
}

// Dependencies returns all of the dependencies directly needed by an
// manifests asset.
func (m *manifests) Dependencies() []asset.Asset {
	return []asset.Asset{
		m.installConfig,
		m.assetStock.KubeCoreOperator(),
		m.assetStock.NetworkOperator(),
		m.assetStock.KubeAddonOperator(),
		m.assetStock.Mao(),
		m.assetStock.Tectonic(),
		m.rootCA,
		m.etcdCA,
		m.ingressCertKey,
		m.kubeCA,
		m.aggregatorCA,
		m.serviceServingCA,
		m.clusterAPIServerCertKey,
		m.etcdClientCertKey,
		m.apiServerCertKey,
		m.openshiftAPIServerCertKey,
		m.apiServerProxyCertKey,
		m.mcsCertKey,
		m.kubeletCertKey,
		m.serviceAccountKeyPair,
		m.kubeconfig,
		m.workerIgnition,
	}
}

// Generate generates the respective operator config.yml files
func (m *manifests) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	//cvo := dependencies[m.assetStock.ClusterVersionOperator()].Contents[0]
	kco := dependencies[m.assetStock.KubeCoreOperator()].Contents[0]
	no := dependencies[m.assetStock.NetworkOperator()].Contents[0]
	addon := dependencies[m.assetStock.KubeAddonOperator()].Contents[0]
	mao := dependencies[m.assetStock.Mao()].Contents[0]
	installConfig := dependencies[m.installConfig].Contents[0]

	// kco+no+mao go to kube-system config map
	kubeSys, err := configMap("kube-system", "cluster-config-v1", genericData{
		"kco-config":     string(kco.Data),
		"network-config": string(no.Data),
		"install-config": string(installConfig.Data),
		"mao-config":     string(mao.Data),
	})
	if err != nil {
		return nil, err
	}

	// addon goes to openshift system
	tectonicSys, err := configMap("tectonic-system", "cluster-config-v1", genericData{
		"addon-config": string(addon.Data),
	})
	if err != nil {
		return nil, err
	}

	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: filepath.Join(manifestDir, "cluster-config.yaml"),
				Data: []byte(kubeSys),
			},
			{
				Name: filepath.Join("tectonic", "cluster-config.yaml"),
				Data: []byte(tectonicSys),
			},
		},
	}
	state.Contents = append(state.Contents, m.generateBootKubeManifests(dependencies)...)
	return state, nil
}

func (m *manifests) generateBootKubeManifests(dependencies map[asset.Asset]*asset.State) []asset.Content {
	ic, err := installconfig.GetInstallConfig(m.installConfig, dependencies)
	if err != nil {
		return nil
	}
	templateData := &bootkubeTemplateData{
		AggregatorCaCert:                base64.StdEncoding.EncodeToString(dependencies[m.aggregatorCA].Contents[certIndex].Data),
		AggregatorCaKey:                 base64.StdEncoding.EncodeToString(dependencies[m.aggregatorCA].Contents[keyIndex].Data),
		ApiserverCert:                   base64.StdEncoding.EncodeToString(dependencies[m.apiServerCertKey].Contents[certIndex].Data),
		ApiserverKey:                    base64.StdEncoding.EncodeToString(dependencies[m.apiServerCertKey].Contents[keyIndex].Data),
		ApiserverProxyCert:              base64.StdEncoding.EncodeToString(dependencies[m.apiServerProxyCertKey].Contents[certIndex].Data),
		ApiserverProxyKey:               base64.StdEncoding.EncodeToString(dependencies[m.apiServerProxyCertKey].Contents[keyIndex].Data),
		Base64encodeCloudProviderConfig: "", // FIXME
		ClusterapiCaCert:                base64.StdEncoding.EncodeToString(dependencies[m.clusterAPIServerCertKey].Contents[certIndex].Data),
		ClusterapiCaKey:                 base64.StdEncoding.EncodeToString(dependencies[m.clusterAPIServerCertKey].Contents[keyIndex].Data),
		EtcdCaCert:                      base64.StdEncoding.EncodeToString(dependencies[m.etcdCA].Contents[certIndex].Data),
		EtcdClientCert:                  base64.StdEncoding.EncodeToString(dependencies[m.etcdClientCertKey].Contents[certIndex].Data),
		EtcdClientKey:                   base64.StdEncoding.EncodeToString(dependencies[m.etcdClientCertKey].Contents[keyIndex].Data),
		KubeCaCert:                      base64.StdEncoding.EncodeToString(dependencies[m.kubeCA].Contents[certIndex].Data),
		KubeCaKey:                       base64.StdEncoding.EncodeToString(dependencies[m.kubeCA].Contents[keyIndex].Data),
		McsTLSCert:                      base64.StdEncoding.EncodeToString(dependencies[m.mcsCertKey].Contents[certIndex].Data),
		McsTLSKey:                       base64.StdEncoding.EncodeToString(dependencies[m.mcsCertKey].Contents[keyIndex].Data),
		OidcCaCert:                      base64.StdEncoding.EncodeToString(dependencies[m.kubeCA].Contents[certIndex].Data),
		OpenshiftApiserverCert:          base64.StdEncoding.EncodeToString(dependencies[m.openshiftAPIServerCertKey].Contents[certIndex].Data),
		OpenshiftApiserverKey:           base64.StdEncoding.EncodeToString(dependencies[m.openshiftAPIServerCertKey].Contents[keyIndex].Data),
		OpenshiftLoopbackKubeconfig:     base64.StdEncoding.EncodeToString(dependencies[m.kubeconfig].Contents[0].Data),
		PullSecret:                      base64.StdEncoding.EncodeToString([]byte(ic.PullSecret)),
		RootCaCert:                      base64.StdEncoding.EncodeToString(dependencies[m.rootCA].Contents[certIndex].Data),
		ServiceaccountKey:               base64.StdEncoding.EncodeToString(dependencies[m.serviceAccountKeyPair].Contents[keyIndex].Data),
		ServiceaccountPub:               base64.StdEncoding.EncodeToString(dependencies[m.serviceAccountKeyPair].Contents[certIndex].Data),
		ServiceServingCaCert:            base64.StdEncoding.EncodeToString(dependencies[m.serviceServingCA].Contents[certIndex].Data),
		ServiceServingCaKey:             base64.StdEncoding.EncodeToString(dependencies[m.serviceServingCA].Contents[keyIndex].Data),
		TectonicNetworkOperatorImage:    "quay.io/coreos/tectonic-network-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		WorkerIgnConfig:                 base64.StdEncoding.EncodeToString(dependencies[m.workerIgnition].Contents[0].Data),
	}

	assetData := map[string][]byte{
		"cluster-apiserver-certs.yaml":          applyTemplateData(bootkube.ClusterApiserverCerts, templateData),
		"ign-config.yaml":                       applyTemplateData(bootkube.IgnConfig, templateData),
		"kube-apiserver-secret.yaml":            applyTemplateData(bootkube.KubeApiserverSecret, templateData),
		"kube-cloud-config.yaml":                applyTemplateData(bootkube.KubeCloudConfig, templateData),
		"kube-controller-manager-secret.yaml":   applyTemplateData(bootkube.KubeControllerManagerSecret, templateData),
		"machine-config-server-tls-secret.yaml": applyTemplateData(bootkube.MachineConfigServerTLSSecret, templateData),
		"openshift-apiserver-secret.yaml":       applyTemplateData(bootkube.OpenshiftApiserverSecret, templateData),
		"pull.json":                             applyTemplateData(bootkube.Pull, templateData),
		"tectonic-network-operator.yaml":        applyTemplateData(bootkube.TectonicNetworkOperator, templateData),

		"01-tectonic-namespace.yaml":                       []byte(bootkube.TectonicNamespace),
		"02-ingress-namespace.yaml":                        []byte(bootkube.IngressNamespace),
		"03-openshift-web-console-namespace.yaml":          []byte(bootkube.OpenshiftWebConsoleNamespace),
		"04-openshift-machine-config-operator.yaml":        []byte(bootkube.OpenshiftMachineConfigOperator),
		"05-openshift-cluster-api-namespace.yaml":          []byte(bootkube.OpenshiftClusterAPINamespace),
		"app-version-kind.yaml":                            []byte(bootkube.AppVersionKind),
		"app-version-mao.yaml":                             []byte(bootkube.AppVersionMao),
		"app-version-tectonic-network.yaml":                []byte(bootkube.AppVersionTectonicNetwork),
		"machine-config-operator-01-images-configmap.yaml": []byte(bootkube.MachineConfigOperator01ImagesConfigmap),
		"operatorstatus-crd.yaml":                          []byte(bootkube.OperatorstatusCrd),
	}

	var assetContents []asset.Content
	for name, data := range assetData {
		assetContents = append(assetContents, asset.Content{
			Name: filepath.Join(manifestDir, name),
			Data: data,
		})
	}

	return assetContents
}

func applyTemplateData(template *template.Template, templateData interface{}) []byte {
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
