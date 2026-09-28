package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/openshift/installer/pkg/asset"
)

// RootCA is the key/cert that signs the Machine Config Server serving cert and
// the internal-release-image serving cert as well as the journal-gatewayd cert,
// which is bootstrap-only.
//
// Despite the "root" name, RootCA is a flat single-tier signer today: it signs
// those leaf certs directly and roots no CA hierarchy. The name is a vestige of
// the original OpenShift 4 design, where it was the apex of an intermediate-CA
// hierarchy that was flattened into independent per-component signers in 2019.
//
// The private key is (as best we know) discarded after installation completes.
type RootCA struct {
	SelfSignedCertKey
}

var _ asset.WritableAsset = (*RootCA)(nil)

// Dependencies returns SignerKeyParams. Configurable PKI requires
// reading the PKI config from InstallConfig, but adding InstallConfig
// as a dependency here would break codepaths that generate signer certs
// without an install-config on disk (e.g. agent create certificates,
// node-joiner). SignerKeyParams reads the config directly from disk
// without triggering InstallConfig validation.
func (c *RootCA) Dependencies() []asset.Asset {
	return []asset.Asset{&SignerKeyParams{}}
}

// Generate generates the MCS/Ignition CA.
func (c *RootCA) Generate(ctx context.Context, parents asset.Parents) error {
	signerKeyParams := &SignerKeyParams{}
	parents.Get(signerKeyParams)

	if !signerKeyParams.ConfigurablePKIEnabled {
		cfg := &CertCfg{
			Subject:   pkix.Name{CommonName: "root-ca", OrganizationalUnit: []string{"openshift"}},
			KeyUsages: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			Validity:  ValidityTenYears(),
			IsCA:      true,
		}
		return c.SelfSignedCertKey.Generate(ctx, cfg, "root-ca", nil)
	}

	keyGen, err := signerKeyParams.ResolveSignerKeyGen("machine-config.machine-config-server-signer")
	if err != nil {
		return err
	}
	cfg := &CertCfg{
		Subject:  pkix.Name{CommonName: "root-ca", OrganizationalUnit: []string{"openshift"}},
		Validity: ValidityTenYears(),
		IsCA:     true,
	}
	return c.SelfSignedCertKey.Generate(ctx, cfg, "root-ca", keyGen)
}

// Name returns the human-friendly name of the asset.
func (c *RootCA) Name() string {
	return "Machine Config Server Root CA"
}
