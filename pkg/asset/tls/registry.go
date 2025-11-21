package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// InternalReleaseRegistrySignerCertKey is a key/cert pair that signs the internal release registry server certs.
type InternalReleaseRegistrySignerCertKey struct {
	SelfSignedCertKey
	LoadedFromDisk bool
}

var _ asset.WritableAsset = (*InternalReleaseRegistrySignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *InternalReleaseRegistrySignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *InternalReleaseRegistrySignerCertKey) Generate(ctx context.Context, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "internal-release-registry-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears(),
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "internal-release-registry-signer")
}

// Load reads the asset files from disk.
func (c *InternalReleaseRegistrySignerCertKey) Load(f asset.FileFetcher) (bool, error) {
	loaded, err := c.loadCertKey(f, "internal-release-registry-signer")
	if err != nil {
		return false, err
	}
	c.LoadedFromDisk = loaded
	return loaded, nil
}

// Name returns the human-friendly name of the asset.
func (c *InternalReleaseRegistrySignerCertKey) Name() string {
	return "Certificate (internal-release-registry-signer)"
}

// InternalReleaseRegistryCertKey is the asset that generates the internal release registry
// serving key/cert pair for both localhost and api-int.
type InternalReleaseRegistryCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*InternalReleaseRegistryCertKey)(nil)

// Dependencies returns the dependency of the cert/key pair.
func (a *InternalReleaseRegistryCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InternalReleaseRegistrySignerCertKey{},
		&installconfig.InstallConfig{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *InternalReleaseRegistryCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &InternalReleaseRegistrySignerCertKey{}
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(ca, installConfig)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "internal-release-registry", Organization: []string{"openshift"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityTenYears(),
		DNSNames: []string{
			"localhost",
			internalAPIAddress(installConfig.Config),
		},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "internal-release-registry", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *InternalReleaseRegistryCertKey) Name() string {
	return "Certificate (internal-release-registry)"
}

// InternalReleaseRegistryLocalhostCertKey is the asset that generates the internal release registry
// serving key/cert pair for localhost only (used in unconfigured-ignition before cluster name is known).
type InternalReleaseRegistryLocalhostCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*InternalReleaseRegistryLocalhostCertKey)(nil)

// Dependencies returns the dependency of the cert/key pair.
func (a *InternalReleaseRegistryLocalhostCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InternalReleaseRegistrySignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *InternalReleaseRegistryLocalhostCertKey) Generate(ctx context.Context, dependencies asset.Parents) error {
	ca := &InternalReleaseRegistrySignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "internal-release-registry-localhost", Organization: []string{"openshift"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Validity:     ValidityTenYears(),
		DNSNames: []string{
			"localhost",
		},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	return a.SignedCertKey.Generate(ctx, cfg, ca, "internal-release-registry-localhost", AppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *InternalReleaseRegistryLocalhostCertKey) Name() string {
	return "Certificate (internal-release-registry-localhost)"
}
