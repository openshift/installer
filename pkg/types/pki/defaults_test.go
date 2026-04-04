package pki

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	"github.com/openshift/installer/pkg/types"
)

func TestEffectiveSignerPKIProfile(t *testing.T) {
	cases := []struct {
		name                string
		ic                  *types.InstallConfig
		expectSignerAlgo    configv1alpha1.KeyAlgorithm
		expectSignerRSA     int32
		expectSignerCurve   configv1alpha1.ECDSACurve
		expectDefaultsAlgo  configv1alpha1.KeyAlgorithm
		expectDefaultsRSA   int32
		expectDefaultsCurve configv1alpha1.ECDSACurve
	}{
		{
			name: "feature gate off - RSA 2048 legacy",
			ic: &types.InstallConfig{
				FeatureSet: configv1.Default,
			},
			expectSignerAlgo:   configv1alpha1.KeyAlgorithmRSA,
			expectSignerRSA:    2048,
			expectDefaultsAlgo: configv1alpha1.KeyAlgorithmRSA,
			expectDefaultsRSA:  2048,
		},
		{
			name: "feature gate on, pki nil - DefaultPKIProfile",
			ic: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
			},
			expectSignerAlgo:    configv1alpha1.KeyAlgorithmECDSA,
			expectSignerCurve:   configv1alpha1.ECDSACurveP384,
			expectDefaultsAlgo:  configv1alpha1.KeyAlgorithmECDSA,
			expectDefaultsCurve: configv1alpha1.ECDSACurveP256,
		},
		{
			name: "feature gate on, pki RSA 4096 - overlay signer only",
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
			expectSignerAlgo:    configv1alpha1.KeyAlgorithmRSA,
			expectSignerRSA:     4096,
			expectDefaultsAlgo:  configv1alpha1.KeyAlgorithmECDSA,
			expectDefaultsCurve: configv1alpha1.ECDSACurveP256,
		},
		{
			name: "feature gate on, pki ECDSA P-384",
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
			expectSignerAlgo:    configv1alpha1.KeyAlgorithmECDSA,
			expectSignerCurve:   configv1alpha1.ECDSACurveP384,
			expectDefaultsAlgo:  configv1alpha1.KeyAlgorithmECDSA,
			expectDefaultsCurve: configv1alpha1.ECDSACurveP256,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			profile := EffectiveSignerPKIProfile(tc.ic)

			// Verify signer config
			assert.Equal(t, tc.expectSignerAlgo, profile.SignerCertificates.Key.Algorithm)
			if tc.expectSignerAlgo == configv1alpha1.KeyAlgorithmRSA {
				assert.Equal(t, tc.expectSignerRSA, profile.SignerCertificates.Key.RSA.KeySize)
			}
			if tc.expectSignerAlgo == configv1alpha1.KeyAlgorithmECDSA {
				assert.Equal(t, tc.expectSignerCurve, profile.SignerCertificates.Key.ECDSA.Curve)
			}

			// Verify defaults
			assert.Equal(t, tc.expectDefaultsAlgo, profile.Defaults.Key.Algorithm)
			if tc.expectDefaultsAlgo == configv1alpha1.KeyAlgorithmRSA {
				assert.Equal(t, tc.expectDefaultsRSA, profile.Defaults.Key.RSA.KeySize)
			}
			if tc.expectDefaultsAlgo == configv1alpha1.KeyAlgorithmECDSA {
				assert.Equal(t, tc.expectDefaultsCurve, profile.Defaults.Key.ECDSA.Curve)
			}
		})
	}
}