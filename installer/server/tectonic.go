package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"text/template"

	"golang.org/x/crypto/bcrypt"

	bootkube "github.com/kubernetes-incubator/bootkube/pkg/asset"
	"github.com/pborman/uuid"

	"github.com/coreos/tectonic-installer/installer/binassets"
	"github.com/coreos/tectonic-installer/installer/server/asset"
	"github.com/coreos/tectonic-installer/installer/server/tlsutil"
)

const (
	tectonicVersion              = "1.5.6-tectonic.1"
	tectonicNamespace            = "tectonic-system"
	tectonicCASecretName         = "tectonic-ca-cert-secret"
	tectonicTLSSecret            = "tectonic-tls-secret"
	tectonicConfigName           = "tectonic-config"
	tectonicIdentityServerSecret = "tectonic-identity-grpc-server-secret"
	tectonicIdentityClientSecret = "tectonic-identity-grpc-client-secret"

	// output paths for generated assets
	assetPathTectonicCASecret              = "tectonic/tectonic-ca-cert-secret.yaml"
	assetPathTectonicTLSSecret             = "tectonic/tectonic-tls-secret.yaml"
	assetPathTectonicConfig                = "tectonic/tectonic-config.yaml"
	assetPathInstallerProgress             = "tectonic.progress"
	assetPathTectonicLicense               = "tectonic/tectonic-license-secret.json"
	assetPathTectonicPullSecret            = "tectonic/coreos-pull-secret.json"
	assetPathRBACAdminBinding              = "tectonic/binding-admin.yaml"
	assetPathIngressRules                  = "tectonic/ingress-rules.yaml"
	assetPathIdentityConfig                = "tectonic/identity-config.yaml"
	assetPathTectonicChannelOperatorConfig = "updater/tectonic-channel-operator-config.json"
	assetPathTectonicIdentityServerSecret  = "tectonic/tectonic-identity-grpc-server-secret.yaml"
	assetPathTectonicIdentityClientSecret  = "tectonic/tectonic-identity-grpc-client-secret.yaml"

	// the two secret types that kubernetes supports for docker configs
	dockerformatv1 = "dockercfg"
	dockerformatv2 = "dockerconfigjson"

	secretByteLength = 16 // 128 bit secrets
	bcryptCost       = 12

	// When dex and console are deployed in cluster it uses well known client IDs.
	oidcConsoleClientID = "tectonic-console"
	oidcKubectlClientID = "tectonic-kubectl"

	// DNS host name of the Tectonic Identity API service.
	tectonicIdentityAPIService = "tectonic-identity-api.tectonic-system.svc.cluster.local"
)

// mustStaticAsset takes a file from the "assets" directory and creates an asset
// of the desired name. It panics if it can't find the asset.
func mustStaticAsset(name, binasset string) asset.Asset {
	return asset.New(name, binassets.MustAsset(binasset))
}

// mustTemplateAsset parses a named binasset as a Template and panics if the
// asset is not found or parses with a non-nil error.
func mustTemplateAsset(binasset string) *template.Template {
	return template.Must(template.New("").Parse(string(binassets.MustAsset(binasset))))
}

var (
	errMissingTectonicConfig = errors.New("installer: Missing tectonic config section")
	errMissingTectonicDomain = errors.New("installer: Missing tectonic domain name")
	errInvalidDockercfg      = errors.New("installer: Invalid dockercfg")

	// Tectonic Kubernetes static assets
	tectonicStaticAssets = []asset.Asset{
		// Bootkube start
		mustStaticAsset("bootkube-start", "bootkube-start"),

		// Tectonic components.
		mustStaticAsset("tectonic/tectonic-namespace.yaml", "tectonic-namespace.yaml"),
		mustStaticAsset("tectonic/heapster-deployment.yaml", "heapster-deployment.yaml"),
		mustStaticAsset("tectonic/heapster-svc.yaml", "heapster-svc.yaml"),
		mustStaticAsset("tectonic/identity-deployment.yaml", "identity-deployment.yaml"),
		mustStaticAsset("tectonic/identity-svc.yaml", "identity-svc.yaml"),
		mustStaticAsset("tectonic/console-svc.yaml", "console-svc.yaml"),
		mustStaticAsset("tectonic/console-deployment.yaml", "console-deployment.yaml"),
		mustStaticAsset("tectonic/stats-emitter-deployment.yaml", "stats-emitter-deployment.yaml"),
		mustStaticAsset("tectonic/identity-api-svc.yaml", "identity-api-svc.yaml"),

		// RBAC roles and role bindings.
		mustStaticAsset("tectonic/role-admin.yaml", "rbac/role-admin.yaml"),
		mustStaticAsset("tectonic/role-readonly.yaml", "rbac/role-readonly.yaml"),
		mustStaticAsset("tectonic/role-user.yaml", "rbac/role-user.yaml"),
		mustStaticAsset("tectonic/binding-discovery.yaml", "rbac/binding-discovery.yaml"),

		// Monitoring
		mustStaticAsset("tectonic/node-exporter-ds.yaml", "monitoring/node-exporter-ds.yaml"),
		mustStaticAsset("tectonic/node-exporter-svc.yaml", "monitoring/node-exporter-svc.yaml"),
		mustStaticAsset("tectonic/prometheus-k8s-cluster-role-binding.yaml", "monitoring/prometheus-k8s-cluster-role-binding.yaml"),
		mustStaticAsset("tectonic/prometheus-k8s-cluster-role.yaml", "monitoring/prometheus-k8s-cluster-role.yaml"),
		mustStaticAsset("tectonic/prometheus-k8s-config.yaml", "monitoring/prometheus-k8s-config.yaml"),
		mustStaticAsset("tectonic/prometheus-k8s.json", "monitoring/prometheus-k8s.json"),
		mustStaticAsset("tectonic/prometheus-k8s-rules.yaml", "monitoring/prometheus-k8s-rules.yaml"),
		mustStaticAsset("tectonic/prometheus-k8s-service-account.yaml", "monitoring/prometheus-k8s-service-account.yaml"),
		mustStaticAsset("tectonic/prometheus-operator-cluster-role-binding.yaml", "monitoring/prometheus-operator-cluster-role-binding.yaml"),
		mustStaticAsset("tectonic/prometheus-operator-cluster-role.yaml", "monitoring/prometheus-operator-cluster-role.yaml"),
		mustStaticAsset("tectonic/prometheus-operator-service-account.yaml", "monitoring/prometheus-operator-service-account.yaml"),
		mustStaticAsset("tectonic/prometheus-operator.yaml", "monitoring/prometheus-operator.yaml"),
		mustStaticAsset("tectonic/prometheus-svc.yaml", "monitoring/prometheus-svc.yaml"),

		// Ingress (common)
		mustStaticAsset("tectonic/default-backend-deployment.yaml", "ingress/default-backend-deployment.yaml"),
		mustStaticAsset("tectonic/default-backend-service.yaml", "ingress/default-backend-service.yaml"),
		mustStaticAsset("tectonic/default-backend-configmap.yaml", "ingress/default-backend-configmap.yaml"),
	}

	tectonicUpdaterStaticAssets = []asset.Asset{
		// Role bindings
		mustStaticAsset("updater/admin-binding.yaml", "candidate/admin-binding.yaml"),

		// Third-party resources definitions
		mustStaticAsset("updater/tectonic-channel-operator-config-kind.yaml", "candidate/tectonic-channel-operator-config-kind.yaml"),
		mustStaticAsset("updater/app-version-kind.yaml", "candidate/app-version-kind.yaml"),
		mustStaticAsset("updater/migration-status-kind.yaml", "candidate/migration-status-kind.yaml"),

		// Daemonsets and Deployments
		mustStaticAsset("updater/node-agent.yaml", "candidate/node-agent.yaml"),
		mustStaticAsset("updater/kube-version-operator-deployment.yaml", "candidate/kube-version-operator-deployment.yaml"),
		mustStaticAsset("updater/tectonic-channel-operator-deployment.yaml", "candidate/tectonic-channel-operator-deployment.yaml"),

		// Third-party resources
		// TODO: Revert back to YAML as soon as github.com/kubernetes/kubernetes/issues/37455 is fixed.
		mustStaticAsset("updater/app-version-tectonic-cluster.json", "candidate/app-version-tectonic-cluster.json"),
		mustStaticAsset("updater/app-version-kubernetes.json", "candidate/app-version-kubernetes.json"),
	}

	// Ingress
	ingressDeployment      = mustStaticAsset("tectonic/nginx-ingress-deployment.yaml", "ingress/nodeport/nginx-ingress-deployment.yaml")
	ingressDaemonset       = mustStaticAsset("tectonic/nginx-ingress-daemonset.yaml", "ingress/hostport/nginx-ingress-daemonset.yaml")
	ingressNodePortService = mustStaticAsset("tectonic/nginx-ingress-service.yaml", "ingress/nodeport/nginx-ingress-service.yaml")
	ingressHostPortService = mustStaticAsset("tectonic/nginx-ingress-service.yaml", "ingress/hostport/nginx-ingress-service.yaml")

	// Tectonic templated assets
	ingressRulesTmpl       = mustTemplateAsset("ingress/ingress-rules.yaml.tmpl")
	tectonicLicenseTmpl    = mustTemplateAsset("tectonic-license-secret.json.tmpl")
	tectonicPullSecretTmpl = mustTemplateAsset("coreos-pull-secret.json.tmpl")
	tectonicAdminBindTmpl  = mustTemplateAsset("rbac/binding-admin.yaml.tmpl")

	tectonicIdentityConfigTmpl = mustTemplateAsset("identity-config.yaml.tmpl")

	// TODO: Revert back to YAML as soon as github.com/kubernetes/kubernetes/issues/37455 is fixed.
	tectonicChannelOperatorConfigTmpl = mustTemplateAsset("candidate/tectonic-channel-operator-config.json.tmpl")

	// Default TectonicUpdaterConfig values
	defaultUpdaterConfig = TectonicUpdaterConfig{
		Server:  "https://tectonic.update.core-os.net",
		Channel: "tectonic-1.5",
		AppID:   "6bc7b986-4654-4a0f-94b3-84ce6feb1db4",
	}
)

// TectonicConfig holds variables needed when generating Tectonic templates
// or assets.
type TectonicConfig struct {
	ControllerDomain string `json:"-"`
	TectonicDomain   string `json:"-"`
	License          string `json:"license"`
	Dockercfg        string `json:"dockercfg"`

	// Identity
	IdentityAdminUser     string `json:"identityAdminUser"`
	IdentityAdminPassword []byte `json:"identityAdminPassword"`

	// Ingress
	IngressKind string `json:"ingressKind"`

	// Updater
	Updater TectonicUpdaterConfig `json:"updater"`
}

// TectonicUpdaterConfig represents the configuration for the
// Tectonic Channel Operator.
type TectonicUpdaterConfig struct {
	Enabled bool `json:"enabled"`

	// Omaha configuration
	Server  string `json:"server"`
	Channel string `json:"channel"`
	AppID   string `json:"appID"`
}

// AssertValid validates the Tectonic data for common errors.
func (t *TectonicConfig) AssertValid() error {
	if t.ControllerDomain == "" {
		return errMissingControllerDomain
	}
	if t.TectonicDomain == "" {
		return errMissingTectonicDomain
	}

	// Set defaults for the Updater.
	if t.Updater.Server == "" {
		t.Updater.Server = defaultUpdaterConfig.Server
	}
	if t.Updater.Channel == "" {
		t.Updater.Channel = defaultUpdaterConfig.Channel
	}
	if t.Updater.AppID == "" {
		t.Updater.AppID = defaultUpdaterConfig.AppID
	}

	return nil
}

// NewTectonicAssets generates Kubernetes manifests for Tectonic clusters.
func NewTectonicAssets(assets []asset.Asset, config *TectonicConfig, m metrics) ([]asset.Asset, error) {
	assets = append(assets, tectonicStaticAssets...)

	// Tectonic CA Secret
	tectonicCASecret, err := newTectonicCASecretAsset(assets)
	if err != nil {
		return []asset.Asset{}, err
	}
	assets = append(assets, tectonicCASecret)

	// Tectonic TLS Secret
	tectonicTLSSecret, err := newTectonicTLSSecretAsset(assets, config)
	if err != nil {
		return []asset.Asset{}, err
	}
	assets = append(assets, tectonicTLSSecret)

	// Tectonic Identity Secrets
	tectonicIdentityServerSecret, tectonicIdentityClientSecret, err := newTectonicIdentitySecretAsset(assets)
	if err != nil {
		return []asset.Asset{}, err
	}
	assets = append(assets, tectonicIdentityServerSecret, tectonicIdentityClientSecret)

	// Tectonic ConfigMap
	tectonicConfig, err := newTectonicConfigAsset(config.TectonicDomain, config.ControllerDomain, m)
	if err != nil {
		return []asset.Asset{}, err
	}
	assets = append(assets, tectonicConfig)

	// Tectonic License
	license, err := newTectonicLicenseSecretAsset(config.License)
	if err != nil {
		return []asset.Asset{}, err
	}
	assets = append(assets, license)

	// Tectonic Pull Secret
	pullSecret, err := newTectonicPullSecretAsset(config.Dockercfg)
	if err != nil {
		return []asset.Asset{}, err
	}
	assets = append(assets, pullSecret)

	// Create a role binding from the admin user to the admin role. That way this
	// user will be an admin.
	adminBinding, err := newTectonicRoleBinding(config.IdentityAdminUser)
	if err != nil {
		return []asset.Asset{}, err
	}
	assets = append(assets, adminBinding)

	// Create the tectonic ingress manifests
	ingressAssets, err := newTectonicIngressAssets(config.IngressKind, config.TectonicDomain)
	if err != nil {
		return nil, err
	}
	assets = append(assets, ingressAssets...)

	// Create the tectonic identity config
	tectonicIdentity, err := newTectonicIdentityAssets(config.TectonicDomain, config.IdentityAdminUser, config.IdentityAdminPassword)
	if err != nil {
		return nil, err
	}
	assets = append(assets, tectonicIdentity)

	// Create the Tectonic updater
	if config.Updater.Enabled {
		tectonicUpdater, err := newTectonicUpdaterAssets(&config.Updater)
		if err != nil {
			return nil, err
		}
		assets = append(assets, tectonicUpdater...)
	}

	return assets, nil
}

// newTectonicCASecretAsset creates the tectonic-ca-cert-secret Secret.
func newTectonicCASecretAsset(assets []asset.Asset) (asset.Asset, error) {
	caCertAsset, err := asset.Find(assets, bootkube.AssetPathCACert)
	if err != nil {
		return nil, err
	}

	generated := map[string][]byte{
		"ca-cert": caCertAsset.Data(),
	}

	secretYAML, err := secretFromData(tectonicCASecretName, tectonicNamespace, generated)
	if err != nil {
		return nil, err
	}

	return asset.New(assetPathTectonicCASecret, secretYAML), nil
}

// newTectonicTLSSecretAsset creates the tectonic-tls-secret.
func newTectonicTLSSecretAsset(assets []asset.Asset, config *TectonicConfig) (asset.Asset, error) {
	caCert, caKey, err := parseCAFiles(assets)
	if err != nil {
		return nil, err
	}

	key, err := tlsutil.NewPrivateKey()
	if err != nil {
		return nil, err
	}

	cfg := tlsutil.CertConfig{
		CommonName: config.TectonicDomain,
	}

	cert, err := tlsutil.NewServerCertificate(cfg, key, caCert, caKey, nil)
	if err != nil {
		return nil, err
	}

	generated := map[string][]byte{
		"tls.crt": tlsutil.EncodeCertificatePEM(cert),
		"tls.key": tlsutil.EncodePrivateKeyPEM(key),
	}

	secretYAML, err := secretFromData(tectonicTLSSecret, tectonicNamespace, generated)
	if err != nil {
		return nil, err
	}

	return asset.New(assetPathTectonicTLSSecret, secretYAML), nil
}

// Ingess configuration options
var ingressKinds = []string{"HostPort", "NodePort"}

// newTectonicIngressAssets create the Tectonic Ingress assets.
func newTectonicIngressAssets(kind string, tectonicDomain string) ([]asset.Asset, error) {
	var assets []asset.Asset
	switch kind {
	case "NodePort":
		assets = []asset.Asset{
			ingressDeployment,
			ingressNodePortService,
		}
	case "HostPort":
		assets = []asset.Asset{
			ingressDaemonset,
			ingressHostPortService,
		}
	default:
		return nil, fmt.Errorf("IngressKind %s not in %v", kind, ingressKinds)
	}

	data := struct {
		TectonicDomain string
	}{
		TectonicDomain: tectonicDomain,
	}
	b, err := renderTemplate(ingressRulesTmpl, data)
	if err != nil {
		return nil, fmt.Errorf("render Ingress rules: %v", err)
	}
	return append(assets, asset.New(assetPathIngressRules, b)), nil
}

// Configuration information required by console. This is currently generated by
// dex's initialization, but in the future may refer to a remote provider.
type consoleOIDCConfig struct {
	// Issuer is dex's fully qualified URL. For example "https://foo.com:32000/identity"
	Issuer string

	// OAuth2 client credentials for console itself. These will be used when logging
	// in users to the console.
	ConsoleClientID string
	ConsoleSecret   string

	// OAuth2 client credentials for kubeclt. This is the client trusted by the API server
	// and the credentials put in the generated kubeconfig files.
	KubectlClientID string
	KubectlSecret   string
}

func newSecret(length int) string {
	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// newTectonicIdentitySecretAsset creates the tectonicIdentityServerSecret and tectonicIdentityClientSecret.
func newTectonicIdentitySecretAsset(assets []asset.Asset) (asset.Asset, asset.Asset, error) {
	caCert, caKey, err := parseCAFiles(assets)
	if err != nil {
		return nil, nil, err
	}

	cfg := tlsutil.CertConfig{
		CommonName: tectonicIdentityAPIService,
	}

	// Tectonic Identity Server Secret.
	sKey, err := tlsutil.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	sCert, err := tlsutil.NewServerCertificate(cfg, sKey, caCert, caKey, nil)
	if err != nil {
		return nil, nil, err
	}

	sGenerated := map[string][]byte{
		"tls-cert": tlsutil.EncodeCertificatePEM(sCert),
		"tls-key":  tlsutil.EncodePrivateKeyPEM(sKey),
		"ca-cert":  tlsutil.EncodeCertificatePEM(caCert),
	}

	serverSecretYAML, err := secretFromData(tectonicIdentityServerSecret, tectonicNamespace, sGenerated)
	if err != nil {
		return nil, nil, err
	}

	// Tectonic Identity Client Secret.
	cKey, err := tlsutil.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	cCert, err := tlsutil.NewClientCertificate(cfg, cKey, caCert, caKey, nil)
	if err != nil {
		return nil, nil, err
	}

	cGenerated := map[string][]byte{
		"tls-cert": tlsutil.EncodeCertificatePEM(cCert),
		"tls-key":  tlsutil.EncodePrivateKeyPEM(cKey),
		"ca-cert":  tlsutil.EncodeCertificatePEM(caCert),
	}

	clientSecretYAML, err := secretFromData(tectonicIdentityClientSecret, tectonicNamespace, cGenerated)
	if err != nil {
		return nil, nil, err
	}

	return asset.New(assetPathTectonicIdentityServerSecret, serverSecretYAML), asset.New(assetPathTectonicIdentityClientSecret, clientSecretYAML), nil

}

// When run in cluster, dex determines the client IDs and secrets used by console.
func newTectonicIdentityAssets(tectonicDomain, user string, password []byte) (asset.Asset, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt failed: %v", err)
	}

	data := struct {
		Issuer string

		ConsoleCallback string
		ConsoleSecret   string
		KubectlSecret   string

		AdminEmail        string
		AdminPasswordHash string
		AdminUserID       string
	}{
		Issuer:            fmt.Sprintf("https://%s/identity", tectonicDomain),
		ConsoleCallback:   fmt.Sprintf("https://%s/auth/callback", tectonicDomain),
		ConsoleSecret:     newSecret(secretByteLength),
		KubectlSecret:     newSecret(secretByteLength),
		AdminUserID:       newSecret(secretByteLength),
		AdminEmail:        user,
		AdminPasswordHash: string(hash),
	}

	tectonicConfig, err := renderTemplate(tectonicIdentityConfigTmpl, data)
	if err != nil {
		return nil, fmt.Errorf("render identity config: %v", err)
	}
	return asset.New(assetPathIdentityConfig, tectonicConfig), nil
}

// newTectonicConfigAsset creates the tectonic-config ConfigMap.
func newTectonicConfigAsset(tectonicDomain, controllerDomain string, m metrics) (asset.Asset, error) {
	data := make(map[string]string)
	// The following are data used for metrics.
	data["clusterID"] = uuid.NewRandom().String()
	data["installerPlatform"] = m.installerPlatform
	data["certificatesStrategy"] = string(m.certificatesStrategy)
	data["tectonicUpdaterEnabled"] = strconv.FormatBool(m.tectonicUpdaterEnabled)
	// The following are data used in the console.
	data["consoleBaseAddress"] = fmt.Sprintf("https://%s", tectonicDomain)
	data["apiServerEndpoint"] = fmt.Sprintf("https://%s:443", controllerDomain)
	data["tectonicVersion"] = tectonicVersion
	data["dexAPIHost"] = fmt.Sprintf("%s:5557", tectonicIdentityAPIService)

	configYAML, err := configMapFromData(tectonicConfigName, tectonicNamespace, data)
	if err != nil {
		return nil, err
	}

	return asset.New(assetPathTectonicConfig, configYAML), nil
}

// newTectonicLicenseSecretAsset templates the tectonic-license Secret.
func newTectonicLicenseSecretAsset(license string) (asset.Asset, error) {
	data := struct {
		License string
	}{base64.StdEncoding.EncodeToString([]byte(license))}
	secretJSON, err := renderTemplate(tectonicLicenseTmpl, data)
	if err != nil {
		return nil, err
	}
	return asset.New(assetPathTectonicLicense, secretJSON), nil
}

// newTectonicPullSecretAsset templates the coreos-pull-secret Secret.
func newTectonicPullSecretAsset(dockercfg string) (asset.Asset, error) {
	dockercfgFormat, err := getDockercfgFormat([]byte(dockercfg))
	if err != nil {
		return nil, err
	}

	data := struct {
		Dockercfg       string
		DockercfgFormat string
	}{
		Dockercfg:       base64.StdEncoding.EncodeToString([]byte(dockercfg)),
		DockercfgFormat: dockercfgFormat,
	}
	secretJSON, err := renderTemplate(tectonicPullSecretTmpl, data)
	if err != nil {
		return nil, err
	}
	return asset.New(assetPathTectonicPullSecret, secretJSON), nil
}

// newTectonicRoleBinding templates the admin ClusterRoleBinding.
func newTectonicRoleBinding(user string) (asset.Asset, error) {
	data := struct {
		Email string
	}{user}
	roleBinding, err := renderTemplate(tectonicAdminBindTmpl, data)
	if err != nil {
		return nil, err
	}
	return asset.New(assetPathRBACAdminBinding, roleBinding), nil
}

func newTectonicUpdaterAssets(config *TectonicUpdaterConfig) (assets []asset.Asset, err error) {
	// Add Static assets
	assets = append(assets, tectonicUpdaterStaticAssets...)

	// Render the Tectonic Channel Operator configuration and add it as an asset.
	operatorConfig, err := renderTemplate(tectonicChannelOperatorConfigTmpl, config)
	if err != nil {
		return nil, err
	}
	assets = append(assets, asset.New(assetPathTectonicChannelOperatorConfig, operatorConfig))

	return assets, nil
}

// renderTemplate executes the give template with data.
func renderTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

// parseCAFiles parses the CA certificate and key assets.
func parseCAFiles(assets []asset.Asset) (*x509.Certificate, *rsa.PrivateKey, error) {
	caCertAsset, err := asset.Find(assets, bootkube.AssetPathCACert)
	if err != nil {
		return nil, nil, err
	}
	caKeyAsset, err := asset.Find(assets, bootkube.AssetPathCAKey)
	if err != nil {
		return nil, nil, err
	}

	caCert, err := tlsutil.ParsePEMEncodedCert(caCertAsset.Data())
	if err != nil {
		return nil, nil, err
	}

	caKey, err := tlsutil.ParsePEMEncodedPrivateKey(caKeyAsset.Data())
	if err != nil {
		return nil, nil, err
	}

	return caCert, caKey, nil
}

func getDockercfgFormat(dockercfg []byte) (string, error) {
	dockercfgObj := make(map[string]interface{})
	err := json.Unmarshal([]byte(dockercfg), &dockercfgObj)
	if err != nil {
		return "", err
	}

	// If an "auths" key exists at the top level, it's v2, othewerwise we
	// expect just a top level quay.io key
	if _, isV2Dockercfg := dockercfgObj["auths"]; isV2Dockercfg {
		return dockerformatv2, nil
	} else if _, isV1Dockercfg := dockercfgObj["quay.io"]; isV1Dockercfg {
		return dockerformatv1, nil
	} else {
		return "", errInvalidDockercfg
	}
}
