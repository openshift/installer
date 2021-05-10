package manifests

import (
	"testing"

	"github.com/ghodss/yaml"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

func TestGenerateInfrastructe(t *testing.T) {
	cases := []struct {
		name                   string
		installConfig          *types.InstallConfig
		expectedInfrastructure *configv1.Infrastructure
	}{{
		name:          "vanilla aws",
		installConfig: icBuild.build(icBuild.forAWS()),
		expectedInfrastructure: infraBuild.build(
			infraBuild.forPlatform(configv1.AWSPlatformType),
			infraBuild.withAWSPlatformSpec(),
			infraBuild.withAWSPlatformStatus(),
		),
	}, {
		name: "service endpoints",
		installConfig: icBuild.build(
			icBuild.forAWS(),
			icBuild.withServiceEndpoint("service", "https://endpoint"),
		),
		expectedInfrastructure: infraBuild.build(
			infraBuild.forPlatform(configv1.AWSPlatformType),
			infraBuild.withServiceEndpoint("service", "https://endpoint"),
		),
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(
				&installconfig.ClusterID{
					UUID:    "test-uuid",
					InfraID: "test-infra-id",
				},
				&installconfig.InstallConfig{Config: tc.installConfig},
				&CloudProviderConfig{},
				&AdditionalTrustBundleConfig{},
			)
			infraAsset := &Infrastructure{}
			err := infraAsset.Generate(parents)
			if !assert.NoError(t, err, "failed to generate asset") {
				return
			}
			if !assert.Len(t, infraAsset.FileList, 1, "expected only one file to be generated") {
				return
			}
			assert.Equal(t, infraAsset.FileList[0].Filename, "manifests/cluster-infrastructure-02-config.yml")
			var actualInfra configv1.Infrastructure
			err = yaml.Unmarshal(infraAsset.FileList[0].Data, &actualInfra)
			if !assert.NoError(t, err, "failed to unmarshal infra manifest") {
				return
			}
			assert.Equal(t, tc.expectedInfrastructure, &actualInfra)
		})
	}
}

type icOption func(*types.InstallConfig)

type icBuildNamespace struct{}

var icBuild icBuildNamespace

func (icBuildNamespace) build(opts ...icOption) *types.InstallConfig {
	ic := &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		BaseDomain:   "test-domain",
		ControlPlane: &types.MachinePool{},
	}
	for _, opt := range opts {
		opt(ic)
	}
	return ic
}

func (b icBuildNamespace) forAWS() icOption {
	return func(ic *types.InstallConfig) {
		if ic.Platform.AWS != nil {
			return
		}
		ic.Platform.AWS = &awstypes.Platform{}
	}
}

func (b icBuildNamespace) withServiceEndpoint(name, url string) icOption {
	return func(ic *types.InstallConfig) {
		b.forAWS()(ic)
		ic.Platform.AWS.ServiceEndpoints = append(
			ic.Platform.AWS.ServiceEndpoints,
			awstypes.ServiceEndpoint{
				Name: name,
				URL:  url,
			},
		)
	}
}

type infraOption func(*configv1.Infrastructure)

type infraBuildNamespace struct{}

var infraBuild infraBuildNamespace

func (b infraBuildNamespace) build(opts ...infraOption) *configv1.Infrastructure {
	infra := &configv1.Infrastructure{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "Infrastructure",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Spec: configv1.InfrastructureSpec{
			PlatformSpec: configv1.PlatformSpec{},
		},
		Status: configv1.InfrastructureStatus{
			InfrastructureName:   "test-infra-id",
			APIServerURL:         "https://api.test-cluster.test-domain:6443",
			APIServerInternalURL: "https://api-int.test-cluster.test-domain:6443",
			PlatformStatus:       &configv1.PlatformStatus{},
		},
	}
	for _, opt := range opts {
		opt(infra)
	}
	return infra
}

func (b infraBuildNamespace) forPlatform(platform configv1.PlatformType) infraOption {
	return func(infra *configv1.Infrastructure) {
		infra.Spec.PlatformSpec.Type = platform
		infra.Status.PlatformStatus.Type = platform
		infra.Status.Platform = platform
	}
}

func (b infraBuildNamespace) withAWSPlatformSpec() infraOption {
	return func(infra *configv1.Infrastructure) {
		if infra.Spec.PlatformSpec.AWS != nil {
			return
		}
		infra.Spec.PlatformSpec.AWS = &configv1.AWSPlatformSpec{}
	}
}

func (b infraBuildNamespace) withAWSPlatformStatus() infraOption {
	return func(infra *configv1.Infrastructure) {
		if infra.Status.PlatformStatus.AWS != nil {
			return
		}
		infra.Status.PlatformStatus.AWS = &configv1.AWSPlatformStatus{}
	}
}

func (b infraBuildNamespace) withServiceEndpoint(name, url string) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withAWSPlatformSpec()(infra)
		b.withAWSPlatformStatus()(infra)
		endpoint := configv1.AWSServiceEndpoint{Name: name, URL: url}
		infra.Spec.PlatformSpec.AWS.ServiceEndpoints = append(infra.Spec.PlatformSpec.AWS.ServiceEndpoints, endpoint)
		infra.Status.PlatformStatus.AWS.ServiceEndpoints = append(infra.Status.PlatformStatus.AWS.ServiceEndpoints, endpoint)
	}
}
