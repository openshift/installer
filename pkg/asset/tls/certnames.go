package tls

import (
	"fmt"

	libcrypto "github.com/openshift/library-go/pkg/crypto"
	libpki "github.com/openshift/library-go/pkg/pki"
)

// resolveSignerKeyGen resolves the KeyPairGenerator for a signer certificate
// from the SignerPKIConfig's profile.
func resolveSignerKeyGen(pkiCfg *SignerPKIConfig, certName string) (libcrypto.KeyPairGenerator, error) {
	return resolveKeyGen(pkiCfg, libpki.CertificateTypeSigner, certName)
}

// resolveKeyGen resolves the KeyPairGenerator for a certificate of the given type.
func resolveKeyGen(pkiCfg *SignerPKIConfig, certType libpki.CertificateType, certName string) (libcrypto.KeyPairGenerator, error) {
	provider := libpki.NewStaticPKIProfileProvider(&pkiCfg.Profile)
	resolved, err := libpki.ResolveCertificateConfig(provider, certType, certName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve PKI config for %s certificate %q: %w", certType, certName, err)
	}
	return resolved.Key, nil
}
