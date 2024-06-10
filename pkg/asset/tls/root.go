package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// RootCA contains the private key and the cert that acts as a certificate
// authority, which is in turn really only used to generate a certificate
// for the Machine Config Server.  More in
// https://docs.openshift.com/container-platform/4.13/security/certificate_types_descriptions/machine-config-operator-certificates.html
// and
// https://github.com/openshift/api/tree/master/tls/docs/MachineConfig%20Operator%20Certificates
// This logic dates back to the very creation of OpenShift 4 and the initial code for this project.
// The private key is (as best we know) completely discarded after an installation is complete.
type RootCA struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*RootCA)(nil)

// Dependencies returns nothing.
func (c *RootCA) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the MCS/Ignition CA.
func (c *RootCA) Generate(ctx context.Context, parents asset.Parents) error {
	cfg := &CertCfg{
		Subject:   pkix.Name{CommonName: "root-ca", OrganizationalUnit: []string{"openshift"}},
		KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		Validity:  ValidityTenYears,
		IsCA:      true,
	}

	return c.SelfSignedCertKey.Generate(ctx, cfg, "root-ca")
}

// Name returns the human-friendly name of the asset.
func (c *RootCA) Name() string {
	return "Machine Config Server Root CA"
}
