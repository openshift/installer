package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// AggregatorCA is the asset that generates the aggregator-ca key/cert pair.
type AggregatorCA struct {
	CertKey
}

var _ asset.Asset = (*AggregatorCA)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *AggregatorCA) Dependencies() []asset.Asset {
	return []asset.Asset{
		&RootCA{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *AggregatorCA) Generate(dependencies asset.Parents) error {
	rootCA := &RootCA{}
	dependencies.Get(rootCA)

	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "aggregator", OrganizationalUnit: []string{"bootkube"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return a.CertKey.Generate(cfg, rootCA, "aggregator-ca", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *AggregatorCA) Name() string {
	return "Certificate (aggregator)"
}
