package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// AdminKubeConfigSignerCertKey is a key/cert pair that signs the admin kubeconfig client certs.
type AdminKubeConfigSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*AdminKubeConfigSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *AdminKubeConfigSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *AdminKubeConfigSignerCertKey) Generate(ctx context.Context, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "admin-kubeconfig-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears(),
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "admin-kubeconfig-signer")
}

// Load reads the asset files from disk.
func (c *AdminKubeConfigSignerCertKey) Load(f asset.FileFetcher) (bool, error) {
	return c.loadCertKey(f, "admin-kubeconfig-signer")
}

// Name returns the human-friendly name of the asset.
func (c *AdminKubeConfigSignerCertKey) Name() string {
	return "Certificate (admin-kubeconfig-signer)"
}

// AdminKubeConfigCABundle is the asset the generates the admin-kubeconfig-ca-bundle,
// which contains all the individual client CAs.
type AdminKubeConfigCABundle struct {
	CertBundle
}

var _ asset.Asset = (*AdminKubeConfigCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *AdminKubeConfigCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AdminKubeConfigSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *AdminKubeConfigCABundle) Generate(ctx context.Context, deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate(ctx, "admin-kubeconfig-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *AdminKubeConfigCABundle) Name() string {
	return "Certificate (admin-kubeconfig-ca-bundle)"
}

// AdminKubeConfigClientCertKey is the asset that generates the key/cert pair for admin client to apiserver.
type AdminKubeConfigClientCertKey struct {
	SignedCertKey
}

var _ asset.WritableAsset = (*AdminKubeConfigClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *AdminKubeConfigClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AdminKubeConfigSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *AdminKubeConfigClientCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &AdminKubeConfigSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:admin", Organization: []string{"system:masters"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears(),
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "admin-kubeconfig-client", DoNotAppendParent)
}

// Load reads the asset files from disk.
func (a *AdminKubeConfigClientCertKey) Load(f asset.FileFetcher) (bool, error) {
	return a.loadCertKey(f, "admin-kubeconfig-client")
}

// Name returns the human-friendly name of the asset.
func (a *AdminKubeConfigClientCertKey) Name() string {
	return "Certificate (admin-kubeconfig-client)"
}
