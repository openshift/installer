package pki

import (
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	features "github.com/openshift/api/features"
	libpki "github.com/openshift/library-go/pkg/pki"

	"github.com/openshift/installer/pkg/types"
)

// EffectiveSignerPKIProfile returns the resolved PKI profile for certificate generation.
// The returned profile is never nil.
//   - If ConfigurablePKI is disabled, returns a legacy profile with RSA 2048 for all certs.
//   - If ConfigurablePKI is enabled, returns DefaultPKIProfile from library-go,
//     overlaying user's signerCertificates if specified in install-config.
func EffectiveSignerPKIProfile(ic *types.InstallConfig) configv1alpha1.PKIProfile {
	if !ic.Enabled(features.FeatureGateConfigurablePKI) {
		// Legacy: explicit RSA 2048 for everything
		rsa2048 := configv1alpha1.KeyConfig{
			Algorithm: configv1alpha1.KeyAlgorithmRSA,
			RSA:       configv1alpha1.RSAKeyConfig{KeySize: 2048},
		}
		return configv1alpha1.PKIProfile{
			Defaults: configv1alpha1.DefaultCertificateConfig{
				Key: rsa2048,
			},
			SignerCertificates: configv1alpha1.CertificateConfig{
				Key: rsa2048,
			},
		}
	}

	// Feature gate on: start from library-go's DefaultPKIProfile
	profile := libpki.DefaultPKIProfile()

	// Overlay user's signerCertificates if specified
	if ic.PKI != nil {
		profile.SignerCertificates = toAPICertificateConfig(ic.PKI.SignerCertificates)
	}

	return profile
}

// EffectiveSignerPKIConfig returns the effective PKI config for signer certificate generation.
// Deprecated: Use EffectiveSignerPKIProfile via the SignerPKIConfig asset instead.
// This function is retained temporarily for callers that have not yet been migrated.
func EffectiveSignerPKIConfig(ic *types.InstallConfig) *types.PKIConfig {
	if ic.PKI != nil {
		return ic.PKI
	}

	if ic.Enabled(features.FeatureGateConfigurablePKI) {
		return &types.PKIConfig{
			SignerCertificates: types.CertificateConfig{
				Key: types.KeyConfig{
					Algorithm: types.KeyAlgorithmECDSA,
					ECDSA:     &types.ECDSAKeyConfig{Curve: types.ECDSACurveP384},
				},
			},
		}
	}

	return nil
}

// toAPICertificateConfig converts installer-local CertificateConfig to the upstream API type.
func toAPICertificateConfig(local types.CertificateConfig) configv1alpha1.CertificateConfig {
	apiKey := configv1alpha1.KeyConfig{
		Algorithm: configv1alpha1.KeyAlgorithm(local.Key.Algorithm),
	}
	if local.Key.RSA != nil {
		apiKey.RSA = configv1alpha1.RSAKeyConfig{KeySize: local.Key.RSA.KeySize}
	}
	if local.Key.ECDSA != nil {
		apiKey.ECDSA = configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurve(local.Key.ECDSA.Curve)}
	}
	return configv1alpha1.CertificateConfig{Key: apiKey}
}