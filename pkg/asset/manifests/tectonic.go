package manifests

import (
	"bytes"
	"encoding/base64"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	content "github.com/openshift/installer/pkg/asset/manifests/content/tectonic"
)

// tectonic generates the dependent resource manifests for tectonic (as against bootkube)
type tectonic struct {
	installConfig  asset.Asset
	ingressCertKey asset.Asset
	kubeCA         asset.Asset
}

var _ asset.Asset = (*tectonic)(nil)

// Name returns a human friendly name for the operator
func (t *tectonic) Name() string {
	return "Tectonic Manifests"
}

// Dependencies returns all of the dependencies directly needed by the
// tectonic asset
func (t *tectonic) Dependencies() []asset.Asset {
	return []asset.Asset{
		t.installConfig,
		t.ingressCertKey,
		t.kubeCA,
	}
}

// Generate generates the respective operator config.yml files
func (t *tectonic) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	ic, err := installconfig.GetInstallConfig(t.installConfig, dependencies)
	if err != nil {
		return nil, err
	}
	ingressContents := dependencies[t.ingressCertKey].Contents
	templateData := &tectonicTemplateData{
		IngressCaCert:                          base64.StdEncoding.EncodeToString(dependencies[t.kubeCA].Contents[certIndex].Data),
		IngressKind:                            "haproxy-router",
		IngressStatusPassword:                  ic.Admin.Password, // FIXME: generate a new random one instead?
		IngressTLSBundle:                       base64.StdEncoding.EncodeToString(bytes.Join([][]byte{ingressContents[certIndex].Data, ingressContents[keyIndex].Data}, []byte{})),
		IngressTLSCert:                         base64.StdEncoding.EncodeToString(ingressContents[certIndex].Data),
		IngressTLSKey:                          base64.StdEncoding.EncodeToString(ingressContents[keyIndex].Data),
		KubeAddonOperatorImage:                 "quay.io/coreos/kube-addon-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		KubeCoreOperatorImage:                  "quay.io/coreos/kube-core-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		PullSecret:                             base64.StdEncoding.EncodeToString([]byte(ic.PullSecret)),
		TectonicIngressControllerOperatorImage: "quay.io/coreos/tectonic-ingress-controller-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		TectonicVersion:                        "1.8.4-tectonic.2",
	}

	assetData := map[string][]byte{
		"99_binding-discovery.yaml":                  []byte(content.BindingDiscovery),
		"99_kube-addon-00-appversion.yaml":           []byte(content.AppVersionKubeAddon),
		"99_kube-addon-01-operator.yaml":             applyTemplateData(content.KubeAddonOperator, templateData),
		"99_kube-core-00-appversion.yaml":            []byte(content.AppVersionKubeCore),
		"99_kube-core-00-operator.yaml":              applyTemplateData(content.KubeCoreOperator, templateData),
		"99_role-admin.yaml":                         []byte(content.RoleAdmin),
		"99_role-user.yaml":                          []byte(content.RoleUser),
		"99_tectonic-ingress-00-appversion.yaml":     []byte(content.AppVersionTectonicIngress),
		"99_tectonic-ingress-01-cluster-config.yaml": applyTemplateData(content.ClusterConfigTectonicIngress, templateData),
		"99_tectonic-ingress-02-tls.yaml":            applyTemplateData(content.TLSTectonicIngress, templateData),
		"99_tectonic-ingress-03-pull.json":           applyTemplateData(content.PullTectonicIngress, templateData),
		"99_tectonic-ingress-04-svc-account.yaml":    []byte(content.SvcAccountTectonicIngress),
		"99_tectonic-ingress-05-operator.yaml":       applyTemplateData(content.TectonicIngressControllerOperator, templateData),
		"99_tectonic-system-00-binding-admin.yaml":   []byte(content.BindingAdmin),
		"99_tectonic-system-01-ca-cert.yaml":         applyTemplateData(content.CaCertTectonicSystem, templateData),
		"99_tectonic-system-02-privileged-scc.yaml":  []byte(content.PriviledgedSccTectonicSystem),
		"99_tectonic-system-03-pull.json":            applyTemplateData(content.PullTectonicSystem, templateData),
	}

	var assetContents []asset.Content
	for name, data := range assetData {
		assetContents = append(assetContents, asset.Content{
			Name: filepath.Join("tectonic", name),
			Data: data,
		})
	}

	return &asset.State{Contents: assetContents}, nil
}
