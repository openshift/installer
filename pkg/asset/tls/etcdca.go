package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
)

// EtcdCA is the asset that generates the etcd-ca key/cert pair.
type EtcdCA struct {
	CertKey
}

var _ asset.Asset = (*EtcdCA)(nil)

// Dependencies returns the dependency of the the cert/key pair, which includes
// the parent CA, and install config if it depends on the install config for
// DNS names, etc.
func (a *EtcdCA) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the cert/key pair based on its dependencies.
func (a *EtcdCA) Generate(dependencies asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "etcd", OrganizationalUnit: []string{"etcd"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	key, crt, err := GenerateRootCertKey(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to generate ETCD client CA")
	}

	a.KeyRaw = PrivateKeyToPem(key)
	a.CertRaw = CertToPem(crt)

	a.generateFiles("etcd-client-ca")

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *EtcdCA) Name() string {
	return "Certificate (etcd)"
}
