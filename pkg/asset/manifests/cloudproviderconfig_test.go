package manifests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/network"
)

func TestGenerateCloudProviderConfig(t *testing.T) {
	cases := []struct {
		name             string
		installConfig    *types.InstallConfig
		expectedCloudCfg map[string]string
	}{
		{
			name: "aws: empty ipFamily - default IPv4",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withAWSRegion("us-east-1"),
			),
			expectedCloudCfg: map[string]string{
				cloudProviderConfigDataKey: `[Global]
`,
			},
		},
		{
			name: "aws: ipFamily IPv4",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withAWSRegion("us-east-1"),
				icBuild.withAWSIPFamily(network.IPv4),
			),
			expectedCloudCfg: map[string]string{
				cloudProviderConfigDataKey: `[Global]
`,
			},
		},
		{
			name: "aws: ipFamily DualStackIPv4Primary",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withAWSRegion("us-east-1"),
				icBuild.withAWSIPFamily(network.DualStackIPv4Primary),
			),
			expectedCloudCfg: map[string]string{
				cloudProviderConfigDataKey: `[Global]
NodeIPFamilies=ipv4
NodeIPFamilies=ipv6
`,
			},
		},
		{
			name: "aws: ipFamily DualStackIPv6Primary",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withAWSRegion("us-east-1"),
				icBuild.withAWSIPFamily(network.DualStackIPv6Primary),
			),
			expectedCloudCfg: map[string]string{
				cloudProviderConfigDataKey: `[Global]
NodeIPFamilies=ipv6
NodeIPFamilies=ipv4
`,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(
				installconfig.MakeAsset(tc.installConfig),
				&installconfig.ClusterID{
					InfraID: "test-infra-id",
				},
			)

			cloudConfig := &CloudProviderConfig{}

			err := cloudConfig.Generate(context.Background(), parents)
			if !assert.NoError(t, err, "failed to generate asset") {
				return
			}

			if tc.expectedCloudCfg == nil {
				assert.Nil(t, cloudConfig.ConfigMap)
				assert.Len(t, cloudConfig.Files(), 0)
				return
			}

			if !assert.NotNil(t, cloudConfig.ConfigMap) {
				return
			}

			assert.Equal(t, "openshift-config", cloudConfig.ConfigMap.Namespace)
			assert.Equal(t, "cloud-provider-config", cloudConfig.ConfigMap.Name)
			assert.Equal(t, tc.expectedCloudCfg, cloudConfig.ConfigMap.Data)

			if !assert.Len(t, cloudConfig.Files(), 1, "expected 1 file to be generated") {
				return
			}

			assert.Equal(t, cloudConfig.Files()[0].Filename, cloudProviderConfigFileName)

			var actualCloudConfig corev1.ConfigMap
			err = yaml.Unmarshal(cloudConfig.Files()[0].Data, &actualCloudConfig)
			if !assert.NoError(t, err, "failed to unmarshal cloud-provider-config manifest") {
				return
			}

			assert.Equal(t, tc.expectedCloudCfg, actualCloudConfig.Data)
		})
	}
}
