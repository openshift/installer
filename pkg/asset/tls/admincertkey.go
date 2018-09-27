package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

//AdminCertKey is the asset that generates the admin key/cert pair.
type AdminCertKey struct {
	CertKey
}

var _ asset.WritableAsset = (*AdminCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *AdminCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&KubeCA{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *AdminCertKey) Generate(dependencies asset.Parents) error {
	kubeCA := &KubeCA{}
	dependencies.Get(kubeCA)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:admin", Organization: []string{"system:masters"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.CertKey.Generate(cfg, kubeCA, "admin", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *AdminCertKey) Name() string {
	return "Certificate (system:admin)"
}
