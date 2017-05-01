package asset

import (
	"bytes"
	"encoding/base64"
	"path/filepath"
	"text/template"

	"github.com/ghodss/yaml"

	"github.com/kubernetes-incubator/bootkube/pkg/asset/internal"
)

const (
	secretNamespace     = "kube-system"
	secretAPIServerName = "kube-apiserver"
	secretCMName        = "kube-controller-manager"
)

func newStaticAssets() Assets {
	var noData interface{}
	assets := Assets{
		mustCreateAssetFromTemplate(AssetPathScheduler, internal.SchedulerTemplate, noData),
		mustCreateAssetFromTemplate(AssetPathSchedulerDisruption, internal.SchedulerDisruptionTemplate, noData),
		mustCreateAssetFromTemplate(AssetPathControllerManagerDisruption, internal.ControllerManagerDisruptionTemplate, noData),
		mustCreateAssetFromTemplate(AssetPathKubeDNSDeployment, internal.DNSDeploymentTemplate, noData),
		mustCreateAssetFromTemplate(AssetPathCheckpointer, internal.CheckpointerTemplate, noData),
		mustCreateAssetFromTemplate(AssetPathKubeFlannel, internal.KubeFlannelTemplate, noData),
		mustCreateAssetFromTemplate(AssetPathKubeSystemSARoleBinding, internal.KubeSystemSARoleBindingTemplate, noData),
	}
	return assets
}

func newDynamicAssets(conf Config) Assets {
	assets := Assets{
		mustCreateAssetFromTemplate(AssetPathControllerManager, internal.ControllerManagerTemplate, conf),
		mustCreateAssetFromTemplate(AssetPathAPIServer, internal.APIServerTemplate, conf),
		mustCreateAssetFromTemplate(AssetPathProxy, internal.ProxyTemplate, conf),
		mustCreateAssetFromTemplate(AssetPathKubeFlannelCfg, internal.KubeFlannelCfgTemplate, conf),
		mustCreateAssetFromTemplate(AssetPathKubeDNSSvc, internal.DNSSvcTemplate, conf),
	}
	if conf.SelfHostKubelet {
		assets = append(assets, mustCreateAssetFromTemplate(AssetPathKubelet, internal.KubeletTemplate, conf))
	}
	if conf.SelfHostedEtcd {
		assets = append(assets,
			mustCreateAssetFromTemplate(AssetPathEtcdOperator, internal.EtcdOperatorTemplate, conf),
			mustCreateAssetFromTemplate(AssetPathEtcdSvc, internal.EtcdSvcTemplate, conf),
			mustCreateAssetFromTemplate(AssetPathKenc, internal.KencTemplate, conf),
		)
	}
	return assets
}

func newKubeConfigAsset(assets Assets, conf Config) (Asset, error) {
	caCert, err := assets.Get(AssetPathCACert)
	if err != nil {
		return Asset{}, err
	}

	kubeletCert, err := assets.Get(AssetPathKubeletCert)
	if err != nil {
		return Asset{}, err
	}

	kubeletKey, err := assets.Get(AssetPathKubeletKey)
	if err != nil {
		return Asset{}, err
	}

	type templateCfg struct {
		Server      string
		CACert      string
		KubeletCert string
		KubeletKey  string
	}

	return assetFromTemplate(AssetPathKubeConfig, internal.KubeConfigTemplate, templateCfg{
		Server:      conf.APIServers[0].String(),
		CACert:      base64.StdEncoding.EncodeToString(caCert.Data),
		KubeletCert: base64.StdEncoding.EncodeToString(kubeletCert.Data),
		KubeletKey:  base64.StdEncoding.EncodeToString(kubeletKey.Data),
	})
}

func newAPIServerSecretAsset(assets Assets) (Asset, error) {
	secretAssets := []string{
		AssetPathAPIServerKey,
		AssetPathAPIServerCert,
		AssetPathServiceAccountPubKey,
		AssetPathCACert,
	}

	secretYAML, err := secretFromAssets(secretAPIServerName, secretNamespace, secretAssets, assets)
	if err != nil {
		return Asset{}, err
	}

	return Asset{Name: AssetPathAPIServerSecret, Data: secretYAML}, nil
}

func newControllerManagerSecretAsset(assets Assets) (Asset, error) {
	secretAssets := []string{
		AssetPathServiceAccountPrivKey,
		AssetPathCACert, //TODO(aaron): do we want this also distributed as secret? or expect available on host?
	}

	secretYAML, err := secretFromAssets(secretCMName, secretNamespace, secretAssets, assets)
	if err != nil {
		return Asset{}, err
	}

	return Asset{Name: AssetPathControllerManagerSecret, Data: secretYAML}, nil
}

// TODO(aaron): use actual secret object (need to wrap in apiversion/type)
type secret struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Type       string            `json:"type"`
	Data       map[string]string `json:"data"`
}

func secretFromAssets(name, namespace string, assetNames []string, assets Assets) ([]byte, error) {
	data := make(map[string]string)
	for _, an := range assetNames {
		a, err := assets.Get(an)
		if err != nil {
			return []byte{}, err
		}
		data[filepath.Base(a.Name)] = base64.StdEncoding.EncodeToString(a.Data)
	}
	return yaml.Marshal(secret{
		ApiVersion: "v1",
		Kind:       "Secret",
		Type:       "Opaque",
		Metadata: map[string]string{
			"name":      name,
			"namespace": namespace,
		},
		Data: data,
	})
}

func mustCreateAssetFromTemplate(name string, template []byte, data interface{}) Asset {
	a, err := assetFromTemplate(name, template, data)
	if err != nil {
		panic(err)
	}
	return a
}

func assetFromTemplate(name string, tb []byte, data interface{}) (Asset, error) {
	tmpl, err := template.New(name).Parse(string(tb))
	if err != nil {
		return Asset{}, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return Asset{}, err
	}
	return Asset{Name: name, Data: buf.Bytes()}, nil
}
