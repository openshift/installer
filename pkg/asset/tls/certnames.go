package tls

import (
	"fmt"

	libcrypto "github.com/openshift/library-go/pkg/crypto"
	libpki "github.com/openshift/library-go/pkg/pki"
)

// resolveSignerKeyGen resolves the KeyPairGenerator for a signer certificate
// from the SignerPKIConfig's profile.
func resolveSignerKeyGen(pkiCfg *SignerPKIConfig, certName string) (libcrypto.KeyPairGenerator, error) {
	provider := libpki.NewStaticPKIProfileProvider(&pkiCfg.Profile)
	resolved, err := libpki.ResolveCertificateConfig(provider, libpki.CertificateTypeSigner, certName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve PKI config for signer %q: %w", certName, err)
	}
	return resolved.Key, nil
}
