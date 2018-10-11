package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
)

// RootCA contains the private key and the cert that's
// self-signed as the root CA.
type RootCA struct {
	CertKey
}

var _ asset.WritableAsset = (*RootCA)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *RootCA) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *RootCA) Generate(parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "root-ca", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	key, crt, err := GenerateRootCertKey(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to generate RootCA")
	}

	c.KeyRaw = PrivateKeyToPem(key)
	c.CertRaw = CertToPem(crt)

	c.generateFiles("root-ca")

	return nil
}

// Name returns the human-friendly name of the asset.
func (c *RootCA) Name() string {
	return "Root CA"
}
