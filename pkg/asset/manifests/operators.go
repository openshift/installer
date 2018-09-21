// Package manifests deals with creating manifests for all manifests to be installed for the cluster
package manifests

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/content"
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
func (o *manifests) Name() string {
	return "Common Manifests"
}

// Dependencies returns all of the dependencies directly needed by an
// manifests asset.
func (o *manifests) Dependencies() []asset.Asset {
	return []asset.Asset{
		o.installConfig,
		o.assetStock.KubeCoreOperator(),
		o.assetStock.NetworkOperator(),
		o.assetStock.KubeAddonOperator(),
		o.assetStock.Mao(),
		o.rootCA,
		o.etcdCA,
		o.ingressCertKey,
		o.kubeCA,
		o.aggregatorCA,
		o.serviceServingCA,
		o.clusterAPIServerCertKey,
		o.etcdClientCertKey,
		o.apiServerCertKey,
		o.openshiftAPIServerCertKey,
		o.apiServerProxyCertKey,
		o.adminCertKey,
		o.kubeletCertKey,
		o.mcsCertKey,
		o.serviceAccountKeyPair,
		o.kubeconfig,
	}
}

// Generate generates the respective operator config.yml files
func (o *manifests) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	//cvo := dependencies[o.assetStock.ClusterVersionOperator()].Contents[0]
	kco := dependencies[o.assetStock.KubeCoreOperator()].Contents[0]
	no := dependencies[o.assetStock.NetworkOperator()].Contents[0]
	//ingress := dependencies[o.assetStock.IngressOperator()].Contents[0]
	addon := dependencies[o.assetStock.KubeAddonOperator()].Contents[0]
	mao := dependencies[o.assetStock.Mao()].Contents[0]
	installConfig := dependencies[o.installConfig].Contents[0]

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

	templateAssetContents := o.generateTemplateAssets(dependencies)

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
	state.Contents = append(state.Contents, templateAssetContents...)
	return state, nil
}

func (o *manifests) generateTemplateAssets(dependencies map[asset.Asset]*asset.State) []asset.Content {
	ic, err := installconfig.GetInstallConfig(o.installConfig, dependencies)
	if err != nil {
		return nil
	}
	assetContents := make([]asset.Content, 0)
	templateData := &templateData{
		AggregatorCaCert:                string(dependencies[o.aggregatorCA].Contents[certIndex].Data),
		AggregatorCaKey:                 string(dependencies[o.aggregatorCA].Contents[keyIndex].Data),
		ApiserverCert:                   string(dependencies[o.apiServerCertKey].Contents[certIndex].Data),
		ApiserverKey:                    string(dependencies[o.apiServerCertKey].Contents[keyIndex].Data),
		ApiserverProxyCert:              string(dependencies[o.apiServerProxyCertKey].Contents[certIndex].Data),
		ApiserverProxyKey:               string(dependencies[o.apiServerProxyCertKey].Contents[keyIndex].Data),
		Base64encodeCloudProviderConfig: "", // FIXME
		ClusterapiCaCert:                string(dependencies[o.clusterAPIServerCertKey].Contents[certIndex].Data),
		ClusterapiCaKey:                 string(dependencies[o.clusterAPIServerCertKey].Contents[keyIndex].Data),
		EtcdCaCert:                      string(dependencies[o.etcdCA].Contents[certIndex].Data),
		EtcdClientCert:                  string(dependencies[o.etcdClientCertKey].Contents[certIndex].Data),
		EtcdClientKey:                   string(dependencies[o.etcdClientCertKey].Contents[keyIndex].Data),
		KubeCaCert:                      string(dependencies[o.kubeCA].Contents[certIndex].Data),
		KubeCaKey:                       string(dependencies[o.kubeCA].Contents[keyIndex].Data),
		MachineConfigOperatorImage:      "docker.io/openshift/origin-machine-config-operator:v4.0.0",
		McsTLSCert:                      string(dependencies[o.adminCertKey].Contents[certIndex].Data),
		McsTLSKey:                       string(dependencies[o.adminCertKey].Contents[keyIndex].Data),
		OidcCaCert:                      string(dependencies[o.kubeCA].Contents[certIndex].Data),
		OpenshiftApiserverCert:          string(dependencies[o.openshiftAPIServerCertKey].Contents[certIndex].Data),
		OpenshiftApiserverKey:           string(dependencies[o.openshiftAPIServerCertKey].Contents[keyIndex].Data),
		OpenshiftLoopbackKubeconfig:     string(dependencies[o.kubeconfig].Contents[0].Data),
		PullSecret:                      string(ic.PullSecret),
		RootCaCert:                      string(dependencies[o.rootCA].Contents[certIndex].Data),
		ServiceaccountKey:               string(dependencies[o.serviceAccountKeyPair].Contents[keyIndex].Data),
		ServiceaccountPub:               string(dependencies[o.serviceAccountKeyPair].Contents[certIndex].Data),
		ServiceServingCaCert:            string(dependencies[o.serviceServingCA].Contents[certIndex].Data),
		ServiceServingCaKey:             string(dependencies[o.serviceServingCA].Contents[keyIndex].Data),
		TectonicNetworkOperatorImage:    "quay.io/coreos/tectonic-network-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		WorkerIgnConfig:                 "", // FIXME: this means that depending on ignition assets (risk of cyclical dependencies)
	}

	// belongs to machine api operator
	data := applyTemplateData(content.ClusterApiserverCerts, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "cluster-apiserver-certs.yaml"), Data: []byte(data)})

	// machine api operator
	data = applyTemplateData(content.IgnConfig, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "ign-config.yaml"), Data: []byte(data)})

	// kco
	data = applyTemplateData(content.KubeApiserverSecret, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "kube-apiserver-secret.yaml"), Data: []byte(data)})

	// kco
	data = applyTemplateData(content.KubeCloudConfig, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "kube-cloud-config.yaml"), Data: []byte(data)})

	// kco
	data = applyTemplateData(content.KubeControllerManagerSecret, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "kube-controller-manager-secret.yaml"), Data: []byte(data)})

	// mco
	data = applyTemplateData(content.MachineConfigOperator03Deployment, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-operator-03-deployment.yaml"), Data: []byte(data)})

	// mco
	data = applyTemplateData(content.MachineConfigServerTLSSecret, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-server-tls-secret.yaml"), Data: []byte(data)})

	// kube core
	data = applyTemplateData(content.OpenshiftApiserverSecret, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "openshift-apiserver-secret.yaml"), Data: []byte(data)})

	// common
	data = applyTemplateData(content.Pull, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "pull.json"), Data: []byte(data)})

	// network operator
	data = applyTemplateData(content.TectonicNetworkOperator, templateData)
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "tectonic-network-operator.yaml"), Data: []byte(data)})

	// common
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "01-tectonic-namespace.yaml"), Data: []byte(content.TectonicNamespace)})
	// ingress
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "02-ingress-namespace.yaml"), Data: []byte(content.IngressNamespace)})
	// kao
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "03-openshift-web-console-namespace.yaml"), Data: []byte(content.OpenshiftWebConsoleNamespace)})
	// mco
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "openshift-machine-config-operator.yaml"), Data: []byte(content.OpenshiftMachineConfigOperator)})
	// machine api operator
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "openshift-cluster-api-namespace.yaml"), Data: []byte(content.OpenshiftClusterAPINamespace)})
	// common
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "app-version-kind.yaml"), Data: []byte(content.AppVersionKind)})
	// cmacine api operator
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "app-version-mao.yaml"), Data: []byte(content.AppVersionMao)})
	// network
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "app-version-tectonic-network.yaml"), Data: []byte(content.AppVersionTectonicNetwork)})
	// machine api operator
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-api-operator.yaml"), Data: []byte(content.MachineAPIOperator)})

	// mco
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-operator-00-config-crd.yaml"), Data: []byte(content.MachineConfigOperator00ConfigCrd)})
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-operator-01-images-configmap.yaml"), Data: []byte(content.MachineConfigOperator01ImagesConfigmap)})
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "machine-config-operator-02-rbac.yaml"), Data: []byte(content.MachineConfigOperator02Rbac)})
	// common/cvo
	assetContents = append(assetContents, asset.Content{Name: filepath.Join(manifestDir, "operatorstatus-crd.yaml"), Data: []byte(content.OperatorstatusCrd)})
	return assetContents
}

func applyTemplateData(template *template.Template, templateData interface{}) string {
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.String()
}
