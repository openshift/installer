package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

func TestGatedFeatures(t *testing.T) {
	tests := []struct {
		name                 string
		installConfig        *types.InstallConfig
		expectedFeatureGates []configv1.FeatureGateName
	}{
		{
			name: "no gated features",
			installConfig: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				Platform: types.Platform{
					AWS: &aws.Platform{
						Region: "us-east-1",
					},
				},
				Compute: []types.MachinePool{
					{
						Name: types.MachinePoolComputeRoleName,
						Platform: types.MachinePoolPlatform{
							AWS: &aws.MachinePool{},
						},
					},
				},
			},
			expectedFeatureGates: []configv1.FeatureGateName{},
		},
		{
			name: "dedicated hosts configured",
			installConfig: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				Platform: types.Platform{
					AWS: &aws.Platform{
						Region: "us-east-1",
					},
				},
				Compute: []types.MachinePool{
					{
						Name: types.MachinePoolComputeRoleName,
						Platform: types.MachinePoolPlatform{
							AWS: &aws.MachinePool{
								HostPlacement: &aws.HostPlacement{
									Affinity: func() *aws.HostAffinity {
										a := aws.HostAffinityDedicatedHost
										return &a
									}(),
									DedicatedHost: []aws.DedicatedHost{
										{ID: "h-1234567890abcdef0"},
									},
								},
							},
						},
					},
				},
			},
			expectedFeatureGates: []configv1.FeatureGateName{features.FeatureGateAWSDedicatedHosts},
		},
		{
			name: "dedicated hosts configured on second compute pool",
			installConfig: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				Platform: types.Platform{
					AWS: &aws.Platform{
						Region: "us-east-1",
					},
				},
				Compute: []types.MachinePool{
					{
						Name: types.MachinePoolComputeRoleName,
						Platform: types.MachinePoolPlatform{
							AWS: &aws.MachinePool{},
						},
					},
					{
						Name: "worker-special",
						Platform: types.MachinePoolPlatform{
							AWS: &aws.MachinePool{
								HostPlacement: &aws.HostPlacement{
									Affinity: func() *aws.HostAffinity {
										a := aws.HostAffinityDedicatedHost
										return &a
									}(),
									DedicatedHost: []aws.DedicatedHost{
										{ID: "h-1234567890abcdef0"},
									},
								},
							},
						},
					},
				},
			},
			expectedFeatureGates: []configv1.FeatureGateName{features.FeatureGateAWSDedicatedHosts},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gatedFeatures := GatedFeatures(tc.installConfig)

			// Extract feature gate names from the result
			var actualFeatureGates []configv1.FeatureGateName
			for _, gf := range gatedFeatures {
				if gf.Condition {
					actualFeatureGates = append(actualFeatureGates, gf.FeatureGateName)
				}
			}

			assert.ElementsMatch(t, tc.expectedFeatureGates, actualFeatureGates,
				"Expected feature gates %v but got %v", tc.expectedFeatureGates, actualFeatureGates)
		})
	}
}

func TestGatedFeatures_DedicatedHostsFieldPath(t *testing.T) {
	installConfig := &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		Platform: types.Platform{
			AWS: &aws.Platform{
				Region: "us-east-1",
			},
		},
		Compute: []types.MachinePool{
			{
				Name: types.MachinePoolComputeRoleName,
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						HostPlacement: &aws.HostPlacement{
							Affinity: func() *aws.HostAffinity {
								a := aws.HostAffinityDedicatedHost
								return &a
							}(),
							DedicatedHost: []aws.DedicatedHost{
								{ID: "h-1234567890abcdef0"},
							},
						},
					},
				},
			},
		},
	}

	gatedFeatures := GatedFeatures(installConfig)

	// Find the dedicated hosts feature
	var dedicatedHostFeature *field.Path
	for _, gf := range gatedFeatures {
		if gf.FeatureGateName == features.FeatureGateAWSDedicatedHosts {
			dedicatedHostFeature = gf.Field
			break
		}
	}

	assert.NotNil(t, dedicatedHostFeature, "Expected to find dedicated hosts feature gate")
	expectedPath := field.NewPath("compute").Index(0).Child("platform", "aws", "hostPlacement")
	assert.Equal(t, expectedPath.String(), dedicatedHostFeature.String(),
		"Field path should point to the hostPlacement field")
}
