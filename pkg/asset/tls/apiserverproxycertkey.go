package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// APIServerProxyCertKey is the asset that generates the API server proxy key/cert pair.
type APIServerProxyCertKey struct {
	CertKey
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
		Subject:      pkix.Name{CommonName: "kube-apiserver-proxy", Organization: []string{"kube-master"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.CertKey.Generate(cfg, aggregatorCA, "apiserver-proxy", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *APIServerProxyCertKey) Name() string {
	return "Certificate (kube-apiserver-proxy)"
}
