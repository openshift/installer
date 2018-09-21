package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
)

// RootCA contains the private key and the cert that's
// self-signed as the root CA.
type RootCA struct{}

var _ asset.Asset = (*CertKey)(nil)

// Dependencies returns the dependency of the root-ca, which is empty.
func (c *RootCA) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the root-ca key and cert pair.
func (c *RootCA) Generate(parents map[asset.Asset]*asset.State) (*asset.State, error) {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "root-ca", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	key, crt, err := GenerateRootCertKey(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RootCA %v", err)
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: assetFilePath(RootCAKeyName),
				Data: []byte(PrivateKeyToPem(key)),
			},
			{
				Name: assetFilePath(RootCACertName),
				Data: []byte(CertToPem(crt)),
			},
		},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (c *RootCA) Name() string {
	return "Root CA"
}
