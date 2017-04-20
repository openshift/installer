package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"

	bootkube "github.com/kubernetes-incubator/bootkube/pkg/asset"

	"github.com/coreos/tectonic-installer/installer/server/asset"
	"github.com/coreos/tectonic-installer/installer/server/tlsutil"
)

const (
	// See mapping of service accounts to kubernetes names:
	// http://kubernetes.io/docs/admin/authentication/#service-account-tokens
	kubeServiceAccount = "system:serviceaccount:kube-system:default"

	// NOTE(ericchiang): if the username claim is anything other than "email" it
	// will be prefixed by tectonic identity's issuer URL.
	//
	// See: https://github.com/kubernetes/kubernetes/issues/31380
	oidcUsernameClaim = "email"

	// Output paths for generated assets
	assetPathBootkubeAPIServer = "manifests/kube-apiserver.yaml"
)

var (
	// Bootkube Self-hosted assets
	bootkubeAPIserverTmpl = mustTemplateAsset("self-hosted/kube-apiserver.yaml.tmpl")
)

// BootkubeConfig represents the configuration needed to generate Bootkube
// assets.
type BootkubeConfig struct {
	bootkube.Config
	OIDCIssuer *OIDCIssuer
}

// OIDCIssuer is the OIDC configuration for the Bootkube assets.
type OIDCIssuer struct {
	IssuerURL     string
	ClientID      string
	UsernameClaim string
	CAPath        string
}

// kubeConfig represents Kubernetes kubeconfig credentials.
type kubeConfig struct {
	certificateAuthority string
	clientCertificate    string
	clientKey            string
}

// parseKubeConfig parses the Kubernetes credentials from assets.
func parseKubeConfig(assets []asset.Asset) (*kubeConfig, error) {
	caCert, err := asset.Find(assets, bootkube.AssetPathCACert)
	if err != nil {
		return nil, err
	}

	clientCert, err := asset.Find(assets, bootkube.AssetPathKubeletCert)
	if err != nil {
		return nil, err
	}

	clientKey, err := asset.Find(assets, bootkube.AssetPathKubeletKey)
	if err != nil {
		return nil, err
	}

	return &kubeConfig{
		certificateAuthority: base64.StdEncoding.EncodeToString(caCert.Data()),
		clientCertificate:    base64.StdEncoding.EncodeToString(clientCert.Data()),
		clientKey:            base64.StdEncoding.EncodeToString(clientKey.Data()),
	}, nil
}

// NewBootkubeAssets wraps bootkube default asset generation and replace
// specific assets to fit our needs (e.g. OIDC).
func NewBootkubeAssets(cfg BootkubeConfig) ([]asset.Asset, error) {
	assets := []asset.Asset{}

	// Generate default assets using Bootkube.
	defaultAssets, err := bootkube.NewDefaultAssets(cfg.Config)
	if err != nil {
		return []asset.Asset{}, err
	}
	for _, defaultAsset := range defaultAssets {
		assets = append(assets, asset.New(defaultAsset.Name, defaultAsset.Data))
	}

	// Replace kube-apiserver asset and add AuthZ asset.
	apiServerAssetByte, err := renderTemplate(bootkubeAPIserverTmpl, cfg)
	if err != nil {
		return []asset.Asset{}, err
	}
	assets, err = asset.Replace(assets, asset.New(assetPathBootkubeAPIServer, apiServerAssetByte))
	if err != nil {
		return []asset.Asset{}, err
	}

	return assets, nil
}

// Returns a custom certificate authority certificate or nil.
func getCACertificate(cert string) (*x509.Certificate, error) {
	if cert == "" {
		return nil, nil
	}
	return tlsutil.ParsePEMEncodedCert([]byte(cert))
}

// Returns a custom certificate authority private key or nil.
func getCAPrivateKey(key string) (*rsa.PrivateKey, error) {
	if key == "" {
		return nil, nil
	}
	return tlsutil.ParsePEMEncodedPrivateKey([]byte(key))
}
