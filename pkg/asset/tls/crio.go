package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// CrioCSRSignerCertKey is a key/cert pair that signs the CRI-O client certs.
type CrioCSRSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*CrioCSRSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *CrioCSRSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *CrioCSRSignerCertKey) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "crio-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityOneDay,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "crio-signer")
}

// Name returns the human-friendly name of the asset.
func (c *CrioCSRSignerCertKey) Name() string {
	return "Certificate (crio-signer)"
}

// CrioClientCABundle is the asset the generates the crio-client-ca-bundle,
// which contains all the individual client CAs.
type CrioClientCABundle struct {
	CertBundle
}

var _ asset.Asset = (*CrioClientCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *CrioClientCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&CrioCSRSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *CrioClientCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("crio-client-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *CrioClientCABundle) Name() string {
	return "Certificate (crio-client-ca-bundle)"
}

// CrioServingCABundle is the asset the generates the crio-serving-ca-bundle,
// which contains all the individual serving CAs.
type CrioServingCABundle struct {
	CertBundle
}

var _ asset.Asset = (*CrioServingCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *CrioServingCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&CrioCSRSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *CrioServingCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("crio-serving-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *CrioServingCABundle) Name() string {
	return "Certificate (crio-serving-ca-bundle)"
}

// CrioBootstrapCertSigner is a key/cert pair that signs the CRI-O bootstrap client certs.
type CrioBootstrapCertSigner struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*CrioBootstrapCertSigner)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *CrioBootstrapCertSigner) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *CrioBootstrapCertSigner) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "crio-bootstrap-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "crio-bootstrap-signer")
}

// Name returns the human-friendly name of the asset.
func (c *CrioBootstrapCertSigner) Name() string {
	return "Certificate (crio-bootstrap-signer)"
}

// CrioBootstrapCABundle is the asset the generates the crio-bootstrap-ca-bundle,
// which contains all the individual client CAs.
type CrioBootstrapCABundle struct {
	CertBundle
}

var _ asset.Asset = (*CrioBootstrapCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *CrioBootstrapCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&CrioBootstrapCertSigner{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *CrioBootstrapCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("crio-bootstrap-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *CrioBootstrapCABundle) Name() string {
	return "Certificate (crio-bootstrap-ca-bundle)"
}

// CrioClientCertKey is the asset that generates the key/cert pair for the CRI-O client.
// This credential can be revoked by deleting the configmap containing its signer.
type CrioClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*CrioClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *CrioClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&CrioBootstrapCertSigner{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *CrioClientCertKey) Generate(dependencies asset.Parents) error {
	ca := &CrioBootstrapCertSigner{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject: pkix.Name{CommonName: "crio", Organization: []string{"openshift"}},
		DNSNames: []string{
			"metrics",
			"metrics.crio-metrics",
			"metrics.crio-metrics.svc",
			"metrics.crio-metrics.svc.cluster.local",
		},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.SignedCertKey.Generate(cfg, ca, "crio-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *CrioClientCertKey) Name() string {
	return "Certificate (crio-client)"
}
