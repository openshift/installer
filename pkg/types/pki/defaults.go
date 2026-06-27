package pki

import (
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
)

// DefaultPKIProfile returns the default PKI profile for OpenShift clusters.
// https://github.com/openshift/library-go/blob/12d8376369b7c5b76f688d01089882ca28e351c3/pkg/pki/profile.go#L11-L26
// TODO: This is a local copy until we have an effective PKI profile default
// defined in openshift/api that can be consumed here.
func DefaultPKIProfile() configv1alpha1.PKIProfile {
	return configv1alpha1.PKIProfile{
		Defaults: configv1alpha1.DefaultCertificateConfig{
			Key: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmECDSA,
				ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP256},
			},
		},
		SignerCertificates: configv1alpha1.CertificateConfig{
			Key: configv1alpha1.KeyConfig{
				Algorithm: configv1alpha1.KeyAlgorithmECDSA,
				ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP384},
			},
		},
	}
}

// EffectiveSignerPKIConfig returns the effective PKI config for signer certificate generation.
//   - If user specified pki in install-config, returns that config unchanged.
//   - If ConfigurablePKI feature gate is enabled and pki is nil, returns a synthetic
//     PKIConfig with ECDSA P-384 derived from DefaultPKIProfile().SignerCertificates.
//   - If feature gate is disabled, returns nil (RSA-2048 legacy path).
func EffectiveSignerPKIConfig(ic *types.InstallConfig) *types.PKIConfig {
	if ic.PKI != nil {
		return ic.PKI
	}

	if ic.Enabled(features.FeatureGateConfigurablePKI) {
		// Default signer config matches DefaultPKIProfile().SignerCertificates
		return &types.PKIConfig{
			SignerCertificates: types.CertificateConfig{
				Key: types.KeyConfig{
					Algorithm: types.KeyAlgorithmECDSA,
					ECDSA:     types.ECDSAKeyConfig{Curve: types.ECDSACurveP384},
				},
			},
		}
	}

	return nil
}
