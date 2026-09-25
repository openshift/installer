package pki

import (
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	features "github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	libpki "github.com/openshift/library-go/pkg/pki"
)

// EffectiveProfile returns the resolved PKI profile and whether ConfigurablePKI is enabled.
//   - Feature gate off: returns an explicit RSA-2048 profile (legacy behavior) and false.
//   - Feature gate on, no user PKI: returns library-go's DefaultPKIProfile() and true.
//   - Feature gate on, user PKI set: returns DefaultPKIProfile() with user's signerCertificates overlaid, and true.
func EffectiveProfile(ic *types.InstallConfig) (configv1alpha1.PKIProfile, bool) {
	if ic == nil || !ic.Enabled(features.FeatureGateConfigurablePKI) {
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
		}, false
	}

	profile := libpki.DefaultPKIProfile()

	if ic.PKI != nil {
		profile.SignerCertificates = toAPICertificateConfig(ic.PKI.SignerCertificates)
	}

	return profile, true
}

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
