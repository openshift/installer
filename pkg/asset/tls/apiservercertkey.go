package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// APIServerCertKey is the asset that generates the API server key/cert pair.
type APIServerCertKey struct {
	CertKey
}

var _ asset.Asset = (*APIServerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *APIServerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeCA{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *APIServerCertKey) Generate(dependencies asset.Parents) error {
	kubeCA := &KubeCA{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(kubeCA, installConfig)

	apiServerAddress, err := cidrhost(installConfig.Config.Networking.ServiceCIDR.IPNet, 1)
	if err != nil {
		return errors.Wrap(err, "failed to get API Server address from InstallConfig")
	}

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "kube-apiserver", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
		DNSNames: []string{
			apiAddress(installConfig.Config),
			"kubernetes", "kubernetes.default",
			"kubernetes.default.svc",
			"kubernetes.default.svc.cluster.local",
			"localhost",
		},
		IPAddresses: []net.IP{net.ParseIP(apiServerAddress), net.ParseIP("127.0.0.1")},
	}

	return a.CertKey.Generate(cfg, kubeCA, "apiserver", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *APIServerCertKey) Name() string {
	return "Certificate (kube-apiaserver)"
}
