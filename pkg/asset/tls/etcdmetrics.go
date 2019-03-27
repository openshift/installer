package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// EtcdMetricSignerCertKey is a key/cert pair that signs the etcd-metrics client and server certs.
type EtcdMetricSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*EtcdMetricSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *EtcdMetricSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *EtcdMetricSignerCertKey) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "etcd-metric-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "etcd-metric-signer")
}

// Name returns the human-friendly name of the asset.
func (c *EtcdMetricSignerCertKey) Name() string {
	return "Certificate (etcd-metric-signer)"
}

// EtcdMetricCABundle is the asset the generates the etcd-metrics-ca-bundle,
// which contains all the individual client CAs.
type EtcdMetricCABundle struct {
	CertBundle
}

var _ asset.Asset = (*EtcdMetricCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *EtcdMetricCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdMetricSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *EtcdMetricCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("etcd-metric-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdMetricCABundle) Name() string {
	return "Certificate (etcd-metric-ca-bundle)"
}

// EtcdMetricSignerClientCertKey is the asset that generates the etcd-metrics client key/cert pair.
type EtcdMetricSignerClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*EtcdMetricSignerClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *EtcdMetricSignerClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdMetricSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *EtcdMetricSignerClientCertKey) Generate(dependencies asset.Parents) error {
	ca := &EtcdMetricSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "etcd-metric", OrganizationalUnit: []string{"etcd-metric"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.SignedCertKey.Generate(cfg, ca, "etcd-metric-signer-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdMetricSignerClientCertKey) Name() string {
	return "Certificate (etcd-metric-signer-client)"
}
