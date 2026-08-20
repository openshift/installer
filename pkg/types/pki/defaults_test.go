package pki

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	"github.com/openshift/installer/pkg/types"
)

func TestEffectiveProfile(t *testing.T) {
	testCases := []struct {
		name                  string
		ic                    *types.InstallConfig
		expectEnabled         bool
		expectDefaultsAlgo    configv1alpha1.KeyAlgorithm
		expectDefaultsCurve   configv1alpha1.ECDSACurve
		expectDefaultsRSASize int32
		expectSignerAlgo      configv1alpha1.KeyAlgorithm
		expectSignerCurve     configv1alpha1.ECDSACurve
		expectSignerRSASize   int32
	}{
		{
			name: "feature gate off - legacy RSA-2048",
			ic: &types.InstallConfig{
				FeatureSet: configv1.Default,
			},
			expectEnabled:         false,
			expectDefaultsAlgo:    configv1alpha1.KeyAlgorithmRSA,
			expectDefaultsRSASize: 2048,
			expectSignerAlgo:      configv1alpha1.KeyAlgorithmRSA,
			expectSignerRSASize:   2048,
		},
		{
			name: "feature gate on, pki nil - library-go defaults",
			ic: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
			},
			expectEnabled:       true,
			expectDefaultsAlgo:  configv1alpha1.KeyAlgorithmECDSA,
			expectDefaultsCurve: configv1alpha1.ECDSACurveP256,
			expectSignerAlgo:    configv1alpha1.KeyAlgorithmECDSA,
			expectSignerCurve:   configv1alpha1.ECDSACurveP384,
		},
		{
			name: "feature gate on, user RSA-4096 signers",
			ic: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				PKI: &types.PKIConfig{
					SignerCertificates: types.CertificateConfig{
						Key: types.KeyConfig{
							Algorithm: types.KeyAlgorithmRSA,
							RSA:       &types.RSAKeyConfig{KeySize: 4096},
						},
					},
				},
			},
			expectEnabled:       true,
			expectDefaultsAlgo:  configv1alpha1.KeyAlgorithmECDSA,
			expectDefaultsCurve: configv1alpha1.ECDSACurveP256,
			expectSignerAlgo:    configv1alpha1.KeyAlgorithmRSA,
			expectSignerRSASize: 4096,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			profile, enabled := EffectiveProfile(tc.ic)
			assert.Equal(t, tc.expectEnabled, enabled)
			assert.Equal(t, tc.expectDefaultsAlgo, profile.Defaults.Key.Algorithm)
			if tc.expectDefaultsAlgo == configv1alpha1.KeyAlgorithmECDSA {
				assert.Equal(t, tc.expectDefaultsCurve, profile.Defaults.Key.ECDSA.Curve)
			} else {
				assert.Equal(t, tc.expectDefaultsRSASize, profile.Defaults.Key.RSA.KeySize)
			}
			assert.Equal(t, tc.expectSignerAlgo, profile.SignerCertificates.Key.Algorithm)
			if tc.expectSignerAlgo == configv1alpha1.KeyAlgorithmECDSA {
				assert.Equal(t, tc.expectSignerCurve, profile.SignerCertificates.Key.ECDSA.Curve)
			} else {
				assert.Equal(t, tc.expectSignerRSASize, profile.SignerCertificates.Key.RSA.KeySize)
			}
		})
	}
}
