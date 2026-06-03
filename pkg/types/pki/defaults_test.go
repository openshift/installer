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

	assert.Equal(t, configv1alpha1.KeyAlgorithmRSA, profile.Defaults.Key.Algorithm)
	assert.Equal(t, int32(4096), profile.Defaults.Key.RSA.KeySize)
	assert.Equal(t, configv1alpha1.KeyAlgorithmRSA, profile.SignerCertificates.Key.Algorithm)
	assert.Equal(t, int32(4096), profile.SignerCertificates.Key.RSA.KeySize)
}

func TestEffectiveSignerPKIConfig(t *testing.T) {
	cases := []struct {
		name        string
		ic          *types.InstallConfig
		expectNil   bool
		expectAlgo  types.KeyAlgorithm
		expectSize  int32
		expectCurve types.ECDSACurve
	}{
		{
			name: "feature gate off, pki nil",
			ic: &types.InstallConfig{
				FeatureSet: configv1.Default,
			},
			expectNil: true,
		},
		{
			name: "feature gate on, pki nil - returns RSA-4096 from DefaultPKIProfile",
			ic: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
			},
			expectNil:  false,
			expectAlgo: types.KeyAlgorithmRSA,
			expectSize: 4096,
		},
		{
			name: "feature gate on, pki specified - returns user config",
			ic: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				PKI: &types.PKIConfig{
					SignerCertificates: types.CertificateConfig{
						Key: types.KeyConfig{
							Algorithm: types.KeyAlgorithmECDSA,
							ECDSA:     &types.ECDSAKeyConfig{Curve: types.ECDSACurveP384},
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
					SignerCertificates: types.CertificateConfig{
						Key: types.KeyConfig{
							Algorithm: types.KeyAlgorithmRSA,
							RSA:       &types.RSAKeyConfig{KeySize: 4096},
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
				assert.Equal(t, tc.ic.PKI, result)
			} else {
				assert.Equal(t, tc.expectAlgo, result.SignerCertificates.Key.Algorithm)
				if tc.expectAlgo == types.KeyAlgorithmECDSA {
					assert.Equal(t, tc.expectCurve, result.SignerCertificates.Key.ECDSA.Curve)
				} else {
					assert.Equal(t, tc.expectSize, result.SignerCertificates.Key.RSA.KeySize)
				}
			}
		})
	}
}
