package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// EtcdMetricsSignerCertKey is a key/cert pair that signs the etcd-metrics client and peer certs.
type EtcdMetricsSignerCertKey struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*EtcdMetricsSignerCertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *EtcdMetricsSignerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *EtcdMetricsSignerCertKey) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "etcd-metrics-signer", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "etcd-metrics-signer")
}

// Name returns the human-friendly name of the asset.
func (c *EtcdMetricsSignerCertKey) Name() string {
	return "Certificate (etcd-metrics-signer)"
}

// EtcdMetricsCABundle is the asset the generates the etcd-metrics-ca-bundle,
// which contains all the individual client CAs.
type EtcdMetricsCABundle struct {
	CertBundle
}

var _ asset.Asset = (*EtcdMetricsCABundle)(nil)

// Dependencies returns the dependency of the cert bundle.
func (a *EtcdMetricsCABundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdMetricsSignerCertKey{},
	}
}

// Generate generates the cert bundle based on its dependencies.
func (a *EtcdMetricsCABundle) Generate(deps asset.Parents) error {
	var certs []CertInterface
	for _, asset := range a.Dependencies() {
		deps.Get(asset)
		certs = append(certs, asset.(CertInterface))
	}
	return a.CertBundle.Generate("etcd-metrics-ca-bundle", certs...)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdMetricsCABundle) Name() string {
	return "Certificate (etcd-metrics-ca-bundle)"
}

// EtcdMetricsSignerClientCertKey is the asset that generates the etcd-metrics client key/cert pair.
type EtcdMetricsSignerClientCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*EtcdMetricsSignerClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *EtcdMetricsSignerClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdMetricsSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *EtcdMetricsSignerClientCertKey) Generate(dependencies asset.Parents) error {
	ca := &EtcdMetricsSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "etcd-metrics", OrganizationalUnit: []string{"etcd-metrics"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.SignedCertKey.Generate(cfg, ca, "etcd-metrics-signer-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdMetricsSignerClientCertKey) Name() string {
	return "Certificate (etcd-metrics-signer-client)"
}

// EtcdMetricsSignerServerCertKey is the asset that generates the etcd-metrics server key/cert pair.
type EtcdMetricsSignerServerCertKey struct {
	SignedCertKey
}

var _ asset.Asset = (*EtcdMetricsSignerServerCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *EtcdMetricsSignerServerCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdMetricsSignerCertKey{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *EtcdMetricsSignerServerCertKey) Generate(dependencies asset.Parents) error {
	ca := &EtcdMetricsSignerCertKey{}
	dependencies.Get(ca)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "etcd-metrics", OrganizationalUnit: []string{"etcd-metrics"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames: []string{
			"etcd",
			"etcd.kube-system",
			"etcd.kube-system.svc.cluster.local",
			"etcd.kube-system.svc",
			"localhost",
		},
		Validity: ValidityTenYears,
	}

	return a.SignedCertKey.Generate(cfg, ca, "etcd-metrics-signer-server", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdMetricsSignerServerCertKey) Name() string {
	return "Certificate (etcd-metrics-signer-server)"
}
