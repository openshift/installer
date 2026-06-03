package manifests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	configv1alpha1 "github.com/openshift/api/config/v1alpha1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

func TestPKIConfigurationGenerate(t *testing.T) {
	cases := []struct {
		name                string
		installConfig       *types.InstallConfig
		expectEmpty         bool
		expectMode          configv1alpha1.PKICertificateManagementMode
		expectSignerAlgo    configv1alpha1.KeyAlgorithm
		expectSignerRSA     int32
		expectSignerCurve   configv1alpha1.ECDSACurve
		expectDefaultsAlgo  configv1alpha1.KeyAlgorithm
		expectDefaultsRSA   int32
		expectDefaultsCurve configv1alpha1.ECDSACurve
	}{
		{
			name: "feature gate disabled - no manifest generated",
			installConfig: &types.InstallConfig{
				FeatureSet: configv1.Default,
			},
			expectEmpty: true,
		},
		{
			name: "feature gate enabled, pki nil - mode Default",
			installConfig: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
			},
			expectEmpty: false,
			expectMode:  configv1alpha1.PKICertificateManagementModeDefault,
		},
		{
			name: "feature gate enabled, pki RSA-4096",
			installConfig: &types.InstallConfig{
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
			expectEmpty:        false,
			expectMode:         configv1alpha1.PKICertificateManagementModeCustom,
			expectSignerAlgo:   configv1alpha1.KeyAlgorithmRSA,
			expectSignerRSA:    4096,
			expectDefaultsAlgo: configv1alpha1.KeyAlgorithmRSA,
			expectDefaultsRSA:  4096,
		},
		{
			name: "feature gate enabled, pki ECDSA P-384",
			installConfig: &types.InstallConfig{
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
			expectEmpty:        false,
			expectMode:         configv1alpha1.PKICertificateManagementModeCustom,
			expectSignerAlgo:   configv1alpha1.KeyAlgorithmECDSA,
			expectSignerCurve:  configv1alpha1.ECDSACurveP384,
			expectDefaultsAlgo: configv1alpha1.KeyAlgorithmRSA,
			expectDefaultsRSA:  4096,
		},
		{
			name: "feature gate enabled, pki RSA-2048 explicit",
			installConfig: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				PKI: &types.PKIConfig{
					SignerCertificates: types.CertificateConfig{
						Key: types.KeyConfig{
							Algorithm: types.KeyAlgorithmRSA,
							RSA:       &types.RSAKeyConfig{KeySize: 2048},
						},
					},
				},
			},
			expectEmpty:        false,
			expectMode:         configv1alpha1.PKICertificateManagementModeCustom,
			expectSignerAlgo:   configv1alpha1.KeyAlgorithmRSA,
			expectSignerRSA:    2048,
			expectDefaultsAlgo: configv1alpha1.KeyAlgorithmRSA,
			expectDefaultsRSA:  4096,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(installconfig.MakeAsset(tc.installConfig))

			pkiAsset := &PKIConfiguration{}
			err := pkiAsset.Generate(context.Background(), parents)
			if !assert.NoError(t, err) {
				return
			}

			if tc.expectEmpty {
				assert.Empty(t, pkiAsset.Files())
				return
			}

			if !assert.Len(t, pkiAsset.Files(), 1) {
				return
			}
			assert.Equal(t, "manifests/cluster-pki-02-config.yaml", pkiAsset.Files()[0].Filename)

			// Unmarshal and verify the CR structure
			var pkiCR configv1alpha1.PKI
			err = yaml.Unmarshal(pkiAsset.Files()[0].Data, &pkiCR)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, "config.openshift.io/v1alpha1", pkiCR.APIVersion)
			assert.Equal(t, "PKI", pkiCR.Kind)
			assert.Equal(t, "cluster", pkiCR.Name)
			assert.Equal(t, tc.expectMode, pkiCR.Spec.CertificateManagement.Mode)

			if tc.expectMode != configv1alpha1.PKICertificateManagementModeCustom {
				return
			}

			profile := pkiCR.Spec.CertificateManagement.Custom.PKIProfile

			// Verify defaults
			assert.Equal(t, tc.expectDefaultsAlgo, profile.Defaults.Key.Algorithm)
			if tc.expectDefaultsAlgo == configv1alpha1.KeyAlgorithmRSA {
				assert.Equal(t, tc.expectDefaultsRSA, profile.Defaults.Key.RSA.KeySize)
			}
			if tc.expectDefaultsAlgo == configv1alpha1.KeyAlgorithmECDSA {
				assert.Equal(t, tc.expectDefaultsCurve, profile.Defaults.Key.ECDSA.Curve)
			}

			// Verify signerCertificates
			assert.Equal(t, tc.expectSignerAlgo, profile.SignerCertificates.Key.Algorithm)
			if tc.expectSignerAlgo == configv1alpha1.KeyAlgorithmRSA {
				assert.Equal(t, tc.expectSignerRSA, profile.SignerCertificates.Key.RSA.KeySize)
			}
			if tc.expectSignerAlgo == configv1alpha1.KeyAlgorithmECDSA {
				assert.Equal(t, tc.expectSignerCurve, profile.SignerCertificates.Key.ECDSA.Curve)
			}
		})
	}
}
