package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// EtcdClientCertKey is the asset that generates the etcd client key/cert pair.
type EtcdClientCertKey struct {
	CertKey
}

var _ asset.Asset = (*EtcdClientCertKey)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *EtcdClientCertKey) Dependencies() []asset.Asset {
	return []asset.Asset{
		&EtcdCA{},
	}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *EtcdClientCertKey) Generate(dependencies asset.Parents) error {
	etcdCA := &EtcdCA{}
	dependencies.Get(etcdCA)

	cfg := &CertCfg{
		Subject:      pkix.Name{CommonName: "etcd", OrganizationalUnit: []string{"etcd"}},
		KeyUsages:    x509.KeyUsageKeyEncipherment,
		ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Validity:     ValidityTenYears,
	}

	return a.CertKey.Generate(cfg, etcdCA, "etcd-client", DoNotAppendParent)
}

// Name returns the human-friendly name of the asset.
func (a *EtcdClientCertKey) Name() string {
	return "Certificate (etcd)"
}
