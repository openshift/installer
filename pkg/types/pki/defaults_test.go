package pki

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	"github.com/openshift/installer/pkg/types"
)

func TestDefaultPKIProfile(t *testing.T) {
	profile := DefaultPKIProfile()

	assert.Equal(t, configv1alpha1.KeyAlgorithmECDSA, profile.Defaults.Key.Algorithm)
	assert.Equal(t, configv1alpha1.ECDSACurveP256, profile.Defaults.Key.ECDSA.Curve)
	assert.Equal(t, configv1alpha1.KeyAlgorithmECDSA, profile.SignerCertificates.Key.Algorithm)
	assert.Equal(t, configv1alpha1.ECDSACurveP384, profile.SignerCertificates.Key.ECDSA.Curve)
}

func TestEffectiveSignerPKIConfig(t *testing.T) {
	cases := []struct {
		name        string
		ic          *types.InstallConfig
		expectNil   bool
		expectAlgo  configv1alpha1.KeyAlgorithm
		expectSize  int32
		expectCurve configv1alpha1.ECDSACurve
	}{
		{
			name: "feature gate off, pki nil",
			ic: &types.InstallConfig{
				FeatureSet: configv1.Default,
			},
			expectNil: true,
		},
		{
			name: "feature gate on, pki nil - returns ECDSA P-384 from DefaultPKIProfile",
			ic: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
			},
			expectNil:   false,
			expectAlgo:  configv1alpha1.KeyAlgorithmECDSA,
			expectCurve: configv1alpha1.ECDSACurveP384,
		},
		{
			name: "feature gate on, pki specified - returns user config",
			ic: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				PKI: &types.PKIConfig{
					SignerCertificates: configv1alpha1.CertificateConfig{
						Key: configv1alpha1.KeyConfig{
							Algorithm: configv1alpha1.KeyAlgorithmECDSA,
							ECDSA:     configv1alpha1.ECDSAKeyConfig{Curve: configv1alpha1.ECDSACurveP384},
						},
					},
				},
			},
			expectNil: false,
		},
		{
			name: "feature gate off, pki specified - returns user config",
			ic: &types.InstallConfig{
				FeatureSet: configv1.Default,
				PKI: &types.PKIConfig{
					SignerCertificates: configv1alpha1.CertificateConfig{
						Key: configv1alpha1.KeyConfig{
							Algorithm: configv1alpha1.KeyAlgorithmRSA,
							RSA:       configv1alpha1.RSAKeyConfig{KeySize: 4096},
						},
					},
				},
			},
			expectNil: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := EffectiveSignerPKIConfig(tc.ic)

			if tc.expectNil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)

			if tc.ic.PKI != nil {
				// Should return user's config unchanged
				assert.Equal(t, tc.ic.PKI, result)
			} else {
				// Should return synthesized config from DefaultPKIProfile
				assert.Equal(t, tc.expectAlgo, result.SignerCertificates.Key.Algorithm)
				if tc.expectAlgo == configv1alpha1.KeyAlgorithmECDSA {
					assert.Equal(t, tc.expectCurve, result.SignerCertificates.Key.ECDSA.Curve)
				} else {
					assert.Equal(t, tc.expectSize, result.SignerCertificates.Key.RSA.KeySize)
				}
			}
		})
	}
}
