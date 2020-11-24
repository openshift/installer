package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// AggregatorCA is the asset that generates the aggregator-ca key/cert pair.
// [DEPRECATED]
type AggregatorCA struct {
	SelfSignedCertKey
}

var _ asset.Asset = (*AggregatorCA)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *AggregatorCA) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *AggregatorCA) Generate(dependencies asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "aggregator", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityOneMonth,
		IsCA:      true,
	}

	return a.SelfSignedCertKey.Generate(cfg, "aggregator-ca")
}

// Name returns the human-friendly name of the asset.
func (a *AggregatorCA) Name() string {
	return "Certificate (aggregator)"
}

// APIServerProxyCertKey is the asset that generates the API server proxy key/cert pair.
// [DEPRECATED]
type APIServerProxyCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*APIServerProxyCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *APIServerProxyCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AggregatorCA{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *APIServerProxyCertKey) Generate(dependencies asset.Parents) error {
	aggregatorCA := &AggregatorCA{}
	dependencies.Get(aggregatorCA)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:kube-apiserver-proxy", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityOneMonth,
	}

	return a.SignedCertKey.Generate(cfg, aggregatorCA, "apiserver-proxy", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *APIServerProxyCertKey) Name() string {
	return "Certificate (system:kube-apiserver-proxy)"
}

// AggregatorSignerCertKey is a key/cert pair that signs the aggregator client certs.
type AggregatorSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*AggregatorSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *AggregatorSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *AggregatorSignerCertKey) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "aggregator-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityOneMonth,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "aggregator-signer")
}

// Name returns the human-friendly name of the asset.
func (c *AggregatorSignerCertKey) Name() string {
	return "Certificate (aggregator-signer)"
}

// AggregatorCABundle is the asset the generates the aggregator-ca-bundle,
// which contains all the individual client CAs.
type AggregatorCABundle struct {
	CertBundle
}

var _ asset.Asset = (*AggregatorCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *AggregatorCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AggregatorSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *AggregatorCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("aggregator-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *AggregatorCABundle) Name() string {
	return "Certificate (aggregator-ca-bundle)"
}

// AggregatorClientCertKey is the asset that generates the API server proxy key/cert pair.
type AggregatorClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*AggregatorClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair
func (a *AggregatorClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AggregatorSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *AggregatorClientCertKey) Generate(dependencies asset.Parents) error {
	ca := &AggregatorSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "system:kube-apiserver-proxy", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityOneMonth,
	}

	return a.SignedCertKey.Generate(cfg, ca, "aggregator-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *AggregatorClientCertKey) Name() string {
	return "Certificate (system:kube-apiserver-proxy)"
}
