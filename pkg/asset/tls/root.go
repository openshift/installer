package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
)

// RootCA contains the private key and the cert that's
// self-signed as the root CA.
type RootCA struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*RootCA)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *RootCA) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *RootCA) Generate(log *logrus.Entry, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "root-ca", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(cfg, "root-ca")
}

// Name returns the human-friendly name of the asset.
func (c *RootCA) Name() string {
	return "Root CA"
}
