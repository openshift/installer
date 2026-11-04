package manifests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

func TestAPIServerGenerate(t *testing.T) {
	cases := []struct {
		name         string
		profileType  string
		expectedType configv1.TLSProfileType
		expectErr    bool
		assertSpec   func(t *testing.T, profile *configv1.TLSSecurityProfile)
	}{
		{
			name:         "unset defaults to Modern",
			profileType:  "",
			expectedType: configv1.TLSProfileModernType,
			assertSpec: func(t *testing.T, p *configv1.TLSSecurityProfile) {
				t.Helper()
				assert.NotNil(t, p.Modern)
			},
		},
		{
			name:         "Modern",
			profileType:  string(configv1.TLSProfileModernType),
			expectedType: configv1.TLSProfileModernType,
			assertSpec: func(t *testing.T, p *configv1.TLSSecurityProfile) {
				t.Helper()
				assert.NotNil(t, p.Modern)
			},
		},
		{
			name:         "Intermediate",
			profileType:  string(configv1.TLSProfileIntermediateType),
			expectedType: configv1.TLSProfileIntermediateType,
			assertSpec: func(t *testing.T, p *configv1.TLSSecurityProfile) {
				t.Helper()
				assert.NotNil(t, p.Intermediate)
			},
		},
		{
			name:         "Old",
			profileType:  string(configv1.TLSProfileOldType),
			expectedType: configv1.TLSProfileOldType,
			assertSpec: func(t *testing.T, p *configv1.TLSSecurityProfile) {
				t.Helper()
				assert.NotNil(t, p.Old)
			},
		},
		{
			name:        "unsupported profile type errors",
			profileType: "bogus",
			expectErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ic := icBuild.build(func(ic *types.InstallConfig) {
				ic.TLSProfileType = tc.profileType
			})

			parents := asset.Parents{}
			parents.Add(installconfig.MakeAsset(ic))

			a := &APIServer{}
			err := a.Generate(t.Context(), parents)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			files := a.Files()
			require.Len(t, files, 1)
			assert.Equal(t, "manifests/cluster-apiserver-02-config.yml", files[0].Filename)

			var got configv1.APIServer
			require.NoError(t, yaml.Unmarshal(files[0].Data, &got))

			assert.Equal(t, configv1.TLSAdherencePolicyStrictAllComponents, got.Spec.TLSAdherence)
			require.NotNil(t, got.Spec.TLSSecurityProfile)
			assert.Equal(t, tc.expectedType, got.Spec.TLSSecurityProfile.Type)
			tc.assertSpec(t, got.Spec.TLSSecurityProfile)
		})
	}
}
