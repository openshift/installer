// Package manifests deals with creating manifests for all manifests to be installed for the cluster
package manifests

import (
	"bytes"
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
	adminCertKey              asset.Asset
	kubeletCertKey            asset.Asset
	mcsCertKey                asset.Asset
	serviceAccountKeyPair     asset.Asset
	kubeconfig                asset.Asset
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
		m.adminCertKey,
		m.kubeletCertKey,
		m.serviceAccountKeyPair,
		m.kubeconfig,
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

	bootkubeContents := m.generateBootKubeManifests(dependencies)

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
	state.Contents = append(state.Contents, bootkubeContents...)
	return state, nil
}

func (m *manifests) generateBootKubeManifests(dependencies map[asset.Asset]*asset.State) []asset.Content {
	ic, err := installconfig.GetInstallConfig(m.installConfig, dependencies)
	if err != nil {
		return nil
	}
	assetContents := make([]asset.Content, 0)
	templateData := &bootkubeTemplateData{
		AggregatorCaCert:                string(dependencies[m.aggregatorCA].Contents[certIndex].Data),
		AggregatorCaKey:                 string(dependencies[m.aggregatorCA].Contents[keyIndex].Data),
		ApiserverCert:                   string(dependencies[m.apiServerCertKey].Contents[certIndex].Data),
		ApiserverKey:                    string(dependencies[m.apiServerCertKey].Contents[keyIndex].Data),
		ApiserverProxyCert:              string(dependencies[m.apiServerProxyCertKey].Contents[certIndex].Data),
		ApiserverProxyKey:               string(dependencies[m.apiServerProxyCertKey].Contents[keyIndex].Data),
		Base64encodeCloudProviderConfig: "", // FIXME
		ClusterapiCaCert:                string(dependencies[m.clusterAPIServerCertKey].Contents[certIndex].Data),
		ClusterapiCaKey:                 string(dependencies[m.clusterAPIServerCertKey].Contents[keyIndex].Data),
		EtcdCaCert:                      string(dependencies[m.etcdCA].Contents[certIndex].Data),
		EtcdClientCert:                  string(dependencies[m.etcdClientCertKey].Contents[certIndex].Data),
		EtcdClientKey:                   string(dependencies[m.etcdClientCertKey].Contents[keyIndex].Data),
		KubeCaCert:                      string(dependencies[m.kubeCA].Contents[certIndex].Data),
		KubeCaKey:                       string(dependencies[m.kubeCA].Contents[keyIndex].Data),
		MachineConfigOperatorImage:      "docker.io/openshift/origin-machine-config-operator:v4.0.0",
		McsTLSCert:                      string(dependencies[m.adminCertKey].Contents[certIndex].Data),
		McsTLSKey:                       string(dependencies[m.adminCertKey].Contents[keyIndex].Data),
		OidcCaCert:                      string(dependencies[m.kubeCA].Contents[certIndex].Data),
		OpenshiftApiserverCert:          string(dependencies[m.openshiftAPIServerCertKey].Contents[certIndex].Data),
		OpenshiftApiserverKey:           string(dependencies[m.openshiftAPIServerCertKey].Contents[keyIndex].Data),
		OpenshiftLoopbackKubeconfig:     string(dependencies[m.kubeconfig].Contents[0].Data),
		PullSecret:                      string(ic.PullSecret),
		RootCaCert:                      string(dependencies[m.rootCA].Contents[certIndex].Data),
		ServiceaccountKey:               string(dependencies[m.serviceAccountKeyPair].Contents[keyIndex].Data),
		ServiceaccountPub:               string(dependencies[m.serviceAccountKeyPair].Contents[certIndex].Data),
		ServiceServingCaCert:            string(dependencies[m.serviceServingCA].Contents[certIndex].Data),
		ServiceServingCaKey:             string(dependencies[m.serviceServingCA].Contents[keyIndex].Data),
		TectonicNetworkOperatorImage:    "quay.io/coreos/tectonic-network-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		WorkerIgnConfig:                 "", // FIXME: this means depending on ignition assets (risk of cyclical dependencies)
	}

	// belongs to machine api operator
	data := applyTemplateData(bootkube.ClusterApiserverCerts, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "cluster-apiserver-certs.yaml"), Data: []byte(data)})

	// machine api operator
	data = applyTemplateData(bootkube.IgnConfig, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "ign-config.yaml"), Data: []byte(data)})

	// kco
	data = applyTemplateData(bootkube.KubeApiserverSecret, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "kube-apiserver-secret.yaml"), Data: []byte(data)})

	// kco
	data = applyTemplateData(bootkube.KubeCloudConfig, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "kube-cloud-config.yaml"), Data: []byte(data)})

	// kco
	data = applyTemplateData(bootkube.KubeControllerManagerSecret, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "kube-controller-manager-secret.yaml"), Data: []byte(data)})

	// mco
	data = applyTemplateData(bootkube.MachineConfigOperator03Deployment, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-operator-03-deployment.yaml"), Data: []byte(data)})

	// mco
	data = applyTemplateData(bootkube.MachineConfigServerTLSSecret, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-server-tls-secret.yaml"), Data: []byte(data)})

	// kube core
	data = applyTemplateData(bootkube.OpenshiftApiserverSecret, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "openshift-apiserver-secret.yaml"), Data: []byte(data)})

	// common
	data = applyTemplateData(bootkube.Pull, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "pull.json"), Data: []byte(data)})

	// network operator
	data = applyTemplateData(bootkube.TectonicNetworkOperator, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "tectonic-network-operator.yaml"), Data: []byte(data)})

	// common
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "01-tectonic-namespace.yaml"), Data: []byte(bootkube.TectonicNamespace)})
	// ingress
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "02-ingress-namespace.yaml"), Data: []byte(bootkube.IngressNamespace)})
	// kao
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "03-openshift-web-console-namespace.yaml"), Data: []byte(bootkube.OpenshiftWebConsoleNamespace)})
	// mco
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "openshift-machine-config-operator.yaml"), Data: []byte(bootkube.OpenshiftMachineConfigOperator)})
	// machine api operator
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "openshift-cluster-api-namespace.yaml"), Data: []byte(bootkube.OpenshiftClusterAPINamespace)})
	// common
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "app-version-kind.yaml"), Data: []byte(bootkube.AppVersionKind)})
	// cmacine api operator
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "app-version-mao.yaml"), Data: []byte(bootkube.AppVersionMao)})
	// network
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "app-version-tectonic-network.yaml"), Data: []byte(bootkube.AppVersionTectonicNetwork)})
	// machine api operator
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-api-operator.yaml"), Data: []byte(bootkube.MachineAPIOperator)})

	// mco
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-operator-00-config-crd.yaml"), Data: []byte(bootkube.MachineConfigOperator00ConfigCrd)})
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-operator-01-images-configmap.yaml"), Data: []byte(bootkube.MachineConfigOperator01ImagesConfigmap)})
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-operator-02-rbac.yaml"), Data: []byte(bootkube.MachineConfigOperator02Rbac)})
	// common/cvo
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "operatorstatus-crd.yaml"), Data: []byte(bootkube.OperatorstatusCrd)})
	return assetContents
}

func applyTemplateData(template *template.Template, templateData interface{}) string {
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.String()
}
