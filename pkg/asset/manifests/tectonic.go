package manifests

import (
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/content/tectonic/ingress"
	"github.com/openshift/installer/pkg/asset/manifests/content/tectonic/rbac"
	"github.com/openshift/installer/pkg/asset/manifests/content/tectonic/secrets"
	"github.com/openshift/installer/pkg/asset/manifests/content/tectonic/security"
	"github.com/openshift/installer/pkg/asset/manifests/content/tectonic/updater"
	"github.com/openshift/installer/pkg/asset/manifests/content/tectonic/updater/appversions"
	"github.com/openshift/installer/pkg/asset/manifests/content/tectonic/updater/operators"
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
	manifestDir := "tectonic"
	assetContents := make([]asset.Content, 0)

	ingressContents := dependencies[t.ingressCertKey].Contents
	templateData := &tectonicTemplateData{
		IngressCaCert:                          string(dependencies[t.kubeCA].Contents[certIndex].Data),
		IngressKind:                            "haproxy-router",
		IngressStatusPassword:                  ic.Admin.Password, // FIXME: generate a new random one instead?
		IngressTLSBundle:                       string(ingressContents[certIndex].Data),
		IngressTLSCert:                         string(ingressContents[certIndex].Data),
		IngressTLSKey:                          string(ingressContents[keyIndex].Data),
		KubeAddonOperatorImage:                 "quay.io/coreos/kube-addon-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		KubeCoreOperatorImage:                  "quay.io/coreos/kube-core-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		PullSecret:                             ic.PullSecret,
		TectonicIngressControllerOperatorImage: "quay.io/coreos/tectonic-ingress-controller-operator-dev:3b6952f5a1ba89bb32dd0630faddeaf2779c9a85",
		TectonicVersion:                        "1.8.4-tectonic.2",
	}

	assetContentMap := map[string]string{
		// template files
		"secrets/ingress-tls.yaml":                                    applyTemplateData(secrets.IngressTLS, templateData),
		"secrets/ca-cert.yaml":                                        applyTemplateData(secrets.CaCert, templateData),
		"secrets/pull.json":                                           applyTemplateData(secrets.Pull, templateData),
		"updater/operators/tectonic-ingress-controller-operator.yaml": applyTemplateData(operators.TectonicIngressControllerOperator, templateData),
		"updater/operators/kube-addon-operator.yaml":                  applyTemplateData(operators.KubeAddonOperator, templateData),
		"updater/operators/kube-core-operator.yaml":                   applyTemplateData(operators.KubeCoreOperator, templateData),
		"updater/app_versions/app-version-tectonic-cluster.yaml":      applyTemplateData(appversions.AppVersionTectonicCluster, templateData),
		"ingress/pull.json":                                           applyTemplateData(ingress.Pull, templateData),
		"ingress/cluster-config.yaml":                                 applyTemplateData(ingress.ClusterConfig, templateData),

		// constant files
		"security/priviledged-scc-tectonic.yaml":                 security.PriviledgedSccTectonic,
		"rbac/role-admin.yaml":                                   rbac.RoleAdmin,
		"rbac/binding-admin.yaml":                                rbac.BindingAdmin,
		"rbac/binding-discovery.yaml":                            rbac.BindingDiscovery,
		"rbac/role-user.yaml":                                    rbac.RoleUser,
		"updater/migration-status-kind.yaml":                     updater.MigrationStatusKind,
		"updater/app_versions/app-version-kube-addon.yaml":       appversions.AppVersionKubeAddon,
		"updater/app_versions/app-version-tectonic-ingress.yaml": appversions.AppVersionTectonicIngress,
		"updater/app_versions/app-version-kube-core.yaml":        appversions.AppVersionKubeCore,
		"updater/app-version-kind.yaml":                          updater.AppVersionKind,
		"ingress/svc-account.yaml":                               ingress.SvcAccount,
	}

	for k, v := range assetContentMap {
		assetContent := asset.Content{
			Name: filepath.Join(manifestDir, k),
			Data: []byte(v),
		}
		assetContents = append(assetContents, assetContent)
	}
	state := &asset.State{
		Contents: assetContents,
	}
	return state, nil
}
