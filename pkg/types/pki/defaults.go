package pki

import (
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
)

// DefaultPKIProfile returns the default PKI profile for OpenShift clusters.
// Currently uses RSA-4096 until all day-2 operators (CKAO, CKMO, etc.) support
// ECDSA certificate rotation. Once operator support lands, switch to ECDSA P-384
// signers and ECDSA P-256 defaults to match the upstream library-go profile:
// https://github.com/openshift/library-go/blob/12d8376369b7c5b76f688d01089882ca28e351c3/pkg/pki/profile.go#L11-L26
func DefaultPKIProfile() configv1alpha1.PKIProfile {
	return configv1alpha1.PKIProfile{
		Defaults: configv1alpha1.DefaultCertificateConfig{
			Key: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 4096},
			},
		},
		SignerCertificates: configv1alpha1.CertificateConfig{
			Key: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmRSA,
				RSA:       configv1alpha1.RSAKeyConfig{KeySize: 4096},
			},
		},
	}
}

// EffectiveSignerPKIConfig returns the effective PKI config for signer certificate generation.
//   - If ConfigurablePKI feature gate is disabled, returns nil (RSA-2048 legacy path).
//   - If user specified pki in install-config, returns that config unchanged.
//   - If pki is nil, returns a PKIConfig derived from DefaultPKIProfile().SignerCertificates.
func EffectiveSignerPKIConfig(ic *types.InstallConfig) *types.PKIConfig {
	if ic == nil {
		return nil
	}

	if !ic.Enabled(features.FeatureGateConfigurablePKI) {
		return nil
	}

	if ic.PKI != nil {
		return ic.PKI
	}

	profile := DefaultPKIProfile()
	keyConfig := types.KeyConfig{
		Algorithm: types.KeyAlgorithm(profile.SignerCertificates.Key.Algorithm),
	}
	switch keyConfig.Algorithm {
	case types.KeyAlgorithmRSA:
		keyConfig.RSA = &types.RSAKeyConfig{KeySize: profile.SignerCertificates.Key.RSA.KeySize}
	case types.KeyAlgorithmECDSA:
		keyConfig.ECDSA = &types.ECDSAKeyConfig{Curve: types.ECDSACurve(profile.SignerCertificates.Key.ECDSA.Curve)}
	}
	return &types.PKIConfig{
		SignerCertificates: types.CertificateConfig{
			Key: keyConfig,
		},
	}
}
