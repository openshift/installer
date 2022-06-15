package manifests

import (
	"testing"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
)

func TestAgentManifests_Generate(t *testing.T) {

	fakeSecret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{Name: "fake-secret"},
	}
	fakeInfraEnv := &aiv1beta1.InfraEnv{
		ObjectMeta: v1.ObjectMeta{Name: "fake-infraEnv"},
	}
	fakeStaticNetworkConfig := []*models.HostStaticNetworkConfig{{NetworkYaml: "some-yaml"}}
	fakeNMStatConfig := []*aiv1beta1.NMStateConfig{{ObjectMeta: v1.ObjectMeta{Name: "fake-nmState"}}}
	fakeAgentClusterInstall := &hiveext.AgentClusterInstall{ObjectMeta: v1.ObjectMeta{Name: "fake-agentClusterInstall"}}
	fakeClusterDeployment := &hivev1.ClusterDeployment{ObjectMeta: v1.ObjectMeta{Name: "fake-clusterDeployment"}}
	fakeClusterImageSet := &hivev1.ClusterImageSet{ObjectMeta: v1.ObjectMeta{Name: "fake-clusterImageSet"}}
	fakeAgent := []*aiv1beta1.Agent{{ObjectMeta: v1.ObjectMeta{Name: "fake-agent"}}}

	tests := []struct {
		Name                        string
		Assets                      []asset.WritableAsset
		ExpectedPullSecret          *corev1.Secret
		ExpectedInfraEnv            *aiv1beta1.InfraEnv
		ExpectedStaticNetworkConfig []*models.HostStaticNetworkConfig
		ExpectedNMStateConfig       []*aiv1beta1.NMStateConfig
		ExpectedAgentClusterInstall *hiveext.AgentClusterInstall
		ExpectedClusterDeployment   *hivev1.ClusterDeployment
		ExpectedClusterImageSet     *hivev1.ClusterImageSet
		ExpectedAgent               []*aiv1beta1.Agent
		ExpectedError               string
	}{
		{
			Name: "default",
			Assets: []asset.WritableAsset{
				&AgentPullSecret{Config: fakeSecret},
				&InfraEnv{Config: fakeInfraEnv},
				&NMStateConfig{
					StaticNetworkConfig: fakeStaticNetworkConfig,
					Config:              fakeNMStatConfig,
				},
				&AgentClusterInstall{Config: fakeAgentClusterInstall},
				&ClusterDeployment{Config: fakeClusterDeployment},
				&ClusterImageSet{Config: fakeClusterImageSet},
				&Agent{Config: fakeAgent},
			},
			ExpectedPullSecret:          fakeSecret,
			ExpectedInfraEnv:            fakeInfraEnv,
			ExpectedStaticNetworkConfig: fakeStaticNetworkConfig,
			ExpectedNMStateConfig:       fakeNMStatConfig,
			ExpectedAgentClusterInstall: fakeAgentClusterInstall,
			ExpectedClusterDeployment:   fakeClusterDeployment,
			ExpectedClusterImageSet:     fakeClusterImageSet,
			ExpectedAgent:               fakeAgent,
		},
		{
			Name: "invalid-NMStateLabelSelector",
			Assets: []asset.WritableAsset{
				&AgentPullSecret{},
				&InfraEnv{Config: &aiv1beta1.InfraEnv{
					Spec: aiv1beta1.InfraEnvSpec{
						NMStateConfigLabelSelector: v1.LabelSelector{
							MatchLabels: map[string]string{
								"missing-label": "missing-label",
							},
						},
					},
				}},
				&NMStateConfig{
					StaticNetworkConfig: fakeStaticNetworkConfig,
					Config:              fakeNMStatConfig,
				},
				&AgentClusterInstall{},
				&ClusterDeployment{},
				&ClusterImageSet{},
				&Agent{},
			},
			ExpectedError: "invalid agent configuration: Spec.NMStateConfigLabelSelector.MatchLabels: Required value: infraEnv and fake-nmState.NMStateConfig labels do not match",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			m := &AgentManifests{}

			fakeParent := asset.Parents{}
			for _, a := range tt.Assets {
				fakeParent.Add(a)
			}

			err := m.Generate(fakeParent)
			if tt.ExpectedError != "" {
				assert.Equal(t, tt.ExpectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
			if tt.ExpectedPullSecret != nil {
				assert.Equal(t, tt.ExpectedPullSecret, m.PullSecret)
			}
			if tt.ExpectedInfraEnv != nil {
				assert.Equal(t, tt.ExpectedInfraEnv, m.InfraEnv)
			}
			if tt.ExpectedStaticNetworkConfig != nil {
				assert.Equal(t, tt.ExpectedStaticNetworkConfig, m.StaticNetworkConfigs)
			}
			if tt.ExpectedNMStateConfig != nil {
				assert.Equal(t, tt.ExpectedNMStateConfig, m.NMStateConfigs)
			}
			if tt.ExpectedClusterDeployment != nil {
				assert.Equal(t, tt.ExpectedClusterDeployment, m.ClusterDeployment)
			}
			if tt.ExpectedClusterImageSet != nil {
				assert.Equal(t, tt.ExpectedClusterImageSet, m.ClusterImageSet)
			}
			if tt.ExpectedAgent != nil {
				assert.Equal(t, tt.ExpectedAgent, m.Agents)
			}
		})
	}
}
