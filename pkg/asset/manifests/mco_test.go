package manifests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/gcp"
)

func TestGenerateMCO(t *testing.T) {
	cases := []struct {
		name          string
		installConfig *types.InstallConfig
		expectedMCO   *operatorv1.MachineConfiguration
	}{
		{
			name: "minimal install config doesn't panic",
			installConfig: func() *types.InstallConfig {
				ic := icBuild.build()
				ic.ControlPlane = nil
				return ic
			}(),
			expectedMCO: nil,
		},
		{
			name:          "vanilla aws produces no mco cfg",
			installConfig: icBuild.build(icBuild.forAWS()),
			expectedMCO:   nil,
		},
		{
			name:          "aws with a custom compute image disables mco management",
			installConfig: icBuild.build(icBuild.withAWSComputeAMI()),
			expectedMCO:   mcoBuild.build(mcoBuild.withComputeBootImageMgmtDisabled()),
		},
		{
			name:          "gcp with a custom compute image disables mco management",
			installConfig: icBuild.build(icBuild.withGCPComputeAMI()),
			expectedMCO:   mcoBuild.build(mcoBuild.withComputeBootImageMgmtDisabled()),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fileData, err := os.ReadFile("../../../data/data/manifests/manifests/cluster-mco-02-config.yaml.template")
			if err != nil {
				t.Errorf("failed reading mco template: %v", err)
			}

			mcoData := generateMCOManifest(tc.installConfig, []*asset.File{{Data: fileData}})
			var actualMCO *operatorv1.MachineConfiguration
			if mcoData != nil {
				if err = yaml.Unmarshal(mcoData[0].Data, &actualMCO); err != nil {
					t.Errorf("failed to serialize mco operator configuration: %v", err)
				}
			}
			assert.Equal(t, tc.expectedMCO, actualMCO)
		})
	}
}

type mcoOption func(*operatorv1.MachineConfiguration)

type mcoBuildNamespace struct{}

var mcoBuild mcoBuildNamespace

func (b mcoBuildNamespace) build(opts ...mcoOption) *operatorv1.MachineConfiguration {
	mco := &operatorv1.MachineConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfiguration",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Spec: operatorv1.MachineConfigurationSpec{
			StaticPodOperatorSpec: operatorv1.StaticPodOperatorSpec{
				OperatorSpec: operatorv1.OperatorSpec{
					LogLevel:         operatorv1.Normal,
					OperatorLogLevel: operatorv1.Normal,
				},
			},
		},
	}
	for _, opt := range opts {
		opt(mco)
	}
	return mco
}

func (b mcoBuildNamespace) withComputeBootImageMgmtDisabled() mcoOption {
	return func(mco *operatorv1.MachineConfiguration) {
		mco.Spec.ManagedBootImages = operatorv1.ManagedBootImages{
			MachineManagers: []operatorv1.MachineManager{
				{
					Resource: operatorv1.MachineSets,
					APIGroup: operatorv1.MachineAPI,
					Selection: operatorv1.MachineManagerSelector{
						Mode: operatorv1.None,
					},
				},
			},
		}
	}
}

func (b icBuildNamespace) withAWSComputeAMI() icOption {
	return func(ic *types.InstallConfig) {
		b.forAWS()(ic)
		ic.Compute = []types.MachinePool{
			{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						AMIID: "ami-xxxxxxxxxxxxx",
					},
				},
			},
		}
	}
}

func (b icBuildNamespace) withGCPComputeAMI() icOption {
	return func(ic *types.InstallConfig) {
		b.forGCP()(ic)
		ic.Compute = []types.MachinePool{
			{
				Platform: types.MachinePoolPlatform{
					GCP: &gcp.MachinePool{
						OSImage: &gcp.OSImage{
							Name:    "myMostFavoriteOSImage",
							Project: "myMostFavoriteProject",
						},
					},
				},
			},
		}
	}
}
