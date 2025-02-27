package manifests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/dns"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

const (
	defaultMachineNetwork = "10.0.0.0/16"
	ipv6MachineNetwork    = "2001:db8::/32"
	ipv4APIVIP            = "192.168.1.0"
	ipv4IngressVIP        = "192.168.222.4"
	ipv6APIVIP            = "fe80::1"
	ipv6IngressVIP        = "fe80::2"
)

func TestGenerateInfrastructure(t *testing.T) {
	cases := []struct {
		name                   string
		installConfig          *types.InstallConfig
		expectedInfrastructure *configv1.Infrastructure
		expectedFilesGenerated int
	}{{
		name:          "vanilla aws",
		installConfig: icBuild.build(icBuild.forAWS()),
		expectedInfrastructure: infraBuild.build(
			infraBuild.forPlatform(configv1.AWSPlatformType),
			infraBuild.withAWSPlatformSpec(),
			infraBuild.withAWSPlatformStatus(),
		),
		expectedFilesGenerated: 1,
	}, {
		name: "aws service endpoints",
		installConfig: icBuild.build(
			icBuild.forAWS(),
			icBuild.withAWSServiceEndpoint("service", "https://endpoint"),
		),
		expectedInfrastructure: infraBuild.build(
			infraBuild.forPlatform(configv1.AWSPlatformType),
			infraBuild.withAWSServiceEndpoint("service", "https://endpoint"),
		),
		expectedFilesGenerated: 1,
	}, {
		name: "gcp service endpoints",
		installConfig: icBuild.build(
			icBuild.forGCP(),
			icBuild.withGCPServiceEndpoint(configv1.GCPServiceEndpointNameCompute, "https://googleapis.com/compute/v1/"),
		),
		expectedInfrastructure: infraBuild.build(
			infraBuild.forPlatform(configv1.GCPPlatformType),
			infraBuild.withGCPServiceEndpoint(configv1.GCPServiceEndpointNameCompute, "https://googleapis.com/compute/v1/"),
		),
		expectedFilesGenerated: 2,
	}, {
		name: "azure resource tags",
		installConfig: icBuild.build(
			icBuild.forAzure(),
			icBuild.withResourceTags(map[string]string{"key": "value"}),
		),
		expectedInfrastructure: infraBuild.build(
			infraBuild.forPlatform(configv1.AzurePlatformType),
			infraBuild.withResourceTags([]configv1.AzureResourceTag{{Key: "key", Value: "value"}}),
		),
		expectedFilesGenerated: 1,
	}, {
		name:          "default GCP custom DNS",
		installConfig: icBuild.build(icBuild.forGCP()),
		expectedInfrastructure: infraBuild.build(
			infraBuild.forPlatform(configv1.GCPPlatformType),
			infraBuild.withGCPClusterHostedDNS("Disabled"),
		),
		expectedFilesGenerated: 2,
	}, {
		name: "GCP custom DNS",
		installConfig: icBuild.build(
			icBuild.forGCP(),
			icBuild.withGCPUserProvisionedDNS("Enabled"),
		),
		expectedInfrastructure: infraBuild.build(
			infraBuild.forPlatform(configv1.GCPPlatformType),
			infraBuild.withGCPClusterHostedDNS("Enabled"),
		),
		expectedFilesGenerated: 2,
	},
		{
			name:          "default AWS custom DNS",
			installConfig: icBuild.build(icBuild.forAWS()),
			expectedInfrastructure: infraBuild.build(
				infraBuild.forPlatform(configv1.AWSPlatformType),
				infraBuild.withAWSClusterHostedDNS("Disabled"),
				infraBuild.withAWSPlatformSpec(),
				infraBuild.withAWSPlatformStatus(),
			),
			expectedFilesGenerated: 1,
		}, {
			name: "AWS custom DNS",
			installConfig: icBuild.build(
				icBuild.forAWS(),
				icBuild.withAWSUserProvisionedDNS("Enabled"),
			),
			expectedInfrastructure: infraBuild.build(
				infraBuild.forPlatform(configv1.AWSPlatformType),
				infraBuild.withAWSClusterHostedDNS("Enabled"),
				infraBuild.withAWSPlatformSpec(),
				infraBuild.withAWSPlatformStatus(),
			),
			expectedFilesGenerated: 1,
		}, {
			name: "vsphere with VIPs appended to machine networks",
			installConfig: icBuild.build(
				icBuild.forVSphere(),
				icBuild.withMachineNetwork(defaultMachineNetwork),
				icBuild.withVSphereAPIVIP(ipv4APIVIP),
				icBuild.withVSphereIngressVIP(ipv4IngressVIP),
			),
			expectedInfrastructure: infraBuild.build(
				infraBuild.forPlatform(configv1.VSpherePlatformType),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(defaultMachineNetwork)),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(ipv4APIVIP+"/32")),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(ipv4IngressVIP+"/32")),
				infraBuild.withVSphereAPIVIP(ipv4APIVIP),
				infraBuild.withVSphereIngressVIP(ipv4IngressVIP),
			),
			expectedFilesGenerated: 1,
		}, {
			name: "dual stack vsphere with VIPs appended to machine networks",
			installConfig: icBuild.build(
				icBuild.forVSphere(),
				icBuild.withMachineNetwork(defaultMachineNetwork),
				icBuild.withMachineNetwork(ipv6MachineNetwork),
				icBuild.withVSphereAPIVIP(ipv4APIVIP),
				icBuild.withVSphereIngressVIP(ipv4IngressVIP),
				icBuild.withVSphereAPIVIP(ipv6APIVIP),
				icBuild.withVSphereIngressVIP(ipv6IngressVIP),
			),
			expectedInfrastructure: infraBuild.build(
				infraBuild.forPlatform(configv1.VSpherePlatformType),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(defaultMachineNetwork)),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(ipv6MachineNetwork)),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(ipv4APIVIP+"/32")),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(ipv6APIVIP+"/128")),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(ipv4IngressVIP+"/32")),
				infraBuild.withVSphereMachineNetworkEntry(configv1.CIDR(ipv6IngressVIP+"/128")),
				infraBuild.withVSphereAPIVIP(ipv4APIVIP),
				infraBuild.withVSphereIngressVIP(ipv4IngressVIP),
				infraBuild.withVSphereAPIVIP(ipv6APIVIP),
				infraBuild.withVSphereIngressVIP(ipv6IngressVIP),
			),
			expectedFilesGenerated: 1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(
				&installconfig.ClusterID{
					UUID:    "test-uuid",
					InfraID: "test-infra-id",
				},
				installconfig.MakeAsset(tc.installConfig),
				&CloudProviderConfig{},
				&AdditionalTrustBundleConfig{},
			)
			infraAsset := &Infrastructure{}
			err := infraAsset.Generate(context.Background(), parents)
			if !assert.NoError(t, err, "failed to generate asset") {
				return
			}

			if !assert.Len(t, infraAsset.FileList, tc.expectedFilesGenerated, "did not generate expected amount of files") {
				return
			}
			assert.Equal(t, infraAsset.FileList[tc.expectedFilesGenerated-1].Filename, "manifests/cluster-infrastructure-02-config.yml")
			var actualInfra configv1.Infrastructure
			err = yaml.Unmarshal(infraAsset.FileList[tc.expectedFilesGenerated-1].Data, &actualInfra)
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

func (b icBuildNamespace) forGCP() icOption {
	return func(ic *types.InstallConfig) {
		if ic.Platform.GCP != nil {
			return
		}
		ic.Platform.GCP = &gcptypes.Platform{}
	}
}

func (b icBuildNamespace) forNone() icOption {
	return func(ic *types.InstallConfig) {
		if ic.Platform.None != nil {
			return
		}
		ic.Platform.None = &nonetypes.Platform{}
	}
}

func (b icBuildNamespace) forVSphere() icOption {
	return func(ic *types.InstallConfig) {
		if ic.Platform.VSphere != nil {
			return
		}
		ic.Platform.VSphere = &vspheretypes.Platform{}
	}
}

func (b icBuildNamespace) withMachineNetwork(cidr string) icOption {
	return func(ic *types.InstallConfig) {
		b.forVSphere()(ic)
		if ic.Networking == nil {
			ic.Networking = &types.Networking{}
		}
		mn := types.MachineNetworkEntry{CIDR: *ipnet.MustParseCIDR(cidr)}
		ic.Networking.MachineNetwork = append(ic.Networking.MachineNetwork, mn)
	}
}

func (b icBuildNamespace) withAWSServiceEndpoint(name, url string) icOption {
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

func (b icBuildNamespace) withGCPServiceEndpoint(name configv1.GCPServiceEndpointName, url string) icOption {
	return func(ic *types.InstallConfig) {
		b.forGCP()(ic)
		ic.Platform.GCP.ServiceEndpoints = append(ic.Platform.GCP.ServiceEndpoints, configv1.GCPServiceEndpoint{Name: name, URL: url})
	}
}

func (b icBuildNamespace) withLBType(lbType configv1.AWSLBType) icOption {
	return func(ic *types.InstallConfig) {
		b.forAWS()(ic)
		ic.Platform.AWS.LBType = lbType
	}
}

func (b icBuildNamespace) withGCPUserProvisionedDNS(enabled string) icOption {
	return func(ic *types.InstallConfig) {
		b.forGCP()(ic)
		if enabled == "Enabled" {
			ic.Platform.GCP.UserProvisionedDNS = dns.UserProvisionedDNSEnabled
			ic.FeatureGates = []string{"GCPClusterHostedDNS=true"}
		}
	}
}

func (b icBuildNamespace) withAWSUserProvisionedDNS(enabled string) icOption {
	return func(ic *types.InstallConfig) {
		b.forAWS()(ic)
		if enabled == "Enabled" {
			ic.Platform.AWS.UserProvisionedDNS = dns.UserProvisionedDNSEnabled
			ic.FeatureGates = []string{"AWSClusterHostedDNS=true"}
		}
	}
}

func (b icBuildNamespace) withVSphereAPIVIP(vip string) icOption {
	return func(ic *types.InstallConfig) {
		b.forVSphere()(ic)
		ic.Platform.VSphere.APIVIPs = append(ic.Platform.VSphere.APIVIPs, vip)
	}
}

func (b icBuildNamespace) withVSphereIngressVIP(vip string) icOption {
	return func(ic *types.InstallConfig) {
		b.forVSphere()(ic)
		ic.Platform.VSphere.IngressVIPs = append(ic.Platform.VSphere.IngressVIPs, vip)
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
			InfrastructureName:     "test-infra-id",
			APIServerURL:           "https://api.test-cluster.test-domain:6443",
			APIServerInternalURL:   "https://api-int.test-cluster.test-domain:6443",
			ControlPlaneTopology:   configv1.HighlyAvailableTopologyMode,
			InfrastructureTopology: configv1.HighlyAvailableTopologyMode,
			PlatformStatus:         &configv1.PlatformStatus{},
			CPUPartitioning:        configv1.CPUPartitioningNone,
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
		infra.Status.PlatformStatus.AWS.CloudLoadBalancerConfig = &configv1.CloudLoadBalancerConfig{
			DNSType: configv1.PlatformDefaultDNSType,
		}
	}
}

func (b infraBuildNamespace) withAWSServiceEndpoint(name, url string) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withAWSPlatformSpec()(infra)
		b.withAWSPlatformStatus()(infra)
		endpoint := configv1.AWSServiceEndpoint{Name: name, URL: url}
		infra.Spec.PlatformSpec.AWS.ServiceEndpoints = append(infra.Spec.PlatformSpec.AWS.ServiceEndpoints, endpoint)
		infra.Status.PlatformStatus.AWS.ServiceEndpoints = append(infra.Status.PlatformStatus.AWS.ServiceEndpoints, endpoint)
	}
}

func (b infraBuildNamespace) withGCPServiceEndpoint(name configv1.GCPServiceEndpointName, url string) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withGCPPlatformStatus()(infra)
		endpoint := configv1.GCPServiceEndpoint{Name: name, URL: url}
		infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig = &configv1.CloudLoadBalancerConfig{}
		infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.DNSType = configv1.PlatformDefaultDNSType
		infra.Status.PlatformStatus.GCP.ServiceEndpoints = append(infra.Status.PlatformStatus.GCP.ServiceEndpoints, endpoint)
	}
}

func (b icBuildNamespace) forAzure() icOption {
	return func(ic *types.InstallConfig) {
		if ic.Platform.Azure != nil {
			return
		}
		ic.Platform.Azure = &azuretypes.Platform{}
	}
}

func (b infraBuildNamespace) withAzurePlatformStatus() infraOption {
	return func(infra *configv1.Infrastructure) {
		if infra.Status.PlatformStatus.Azure != nil {
			return
		}
		infra.Status.PlatformStatus.Azure = &configv1.AzurePlatformStatus{
			ResourceGroupName:        infra.Status.InfrastructureName + "-rg",
			NetworkResourceGroupName: infra.Status.InfrastructureName + "-rg",
		}
	}
}

func (b icBuildNamespace) withResourceTags(tags map[string]string) icOption {
	return func(ic *types.InstallConfig) {
		b.forAzure()(ic)
		ic.Platform.Azure.UserTags = tags
	}
}

func (b infraBuildNamespace) withResourceTags(tags []configv1.AzureResourceTag) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withAzurePlatformStatus()(infra)
		infra.Status.PlatformStatus.Azure.ResourceTags = tags
	}
}

func (b infraBuildNamespace) withGCPPlatformStatus() infraOption {
	return func(infra *configv1.Infrastructure) {
		if infra.Status.PlatformStatus.GCP != nil {
			return
		}
		infra.Status.PlatformStatus.GCP = &configv1.GCPPlatformStatus{}
	}
}

func (b infraBuildNamespace) withGCPClusterHostedDNS(enabled string) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withGCPPlatformStatus()(infra)
		infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig = &configv1.CloudLoadBalancerConfig{}
		infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.DNSType = configv1.PlatformDefaultDNSType
		if enabled == "Enabled" {
			infra.Status.PlatformStatus.GCP.CloudLoadBalancerConfig.DNSType = configv1.ClusterHostedDNSType
		}
	}
}

func (b infraBuildNamespace) withAWSClusterHostedDNS(enabled string) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withAWSPlatformStatus()(infra)
		infra.Status.PlatformStatus.AWS.CloudLoadBalancerConfig = &configv1.CloudLoadBalancerConfig{
			DNSType: configv1.PlatformDefaultDNSType,
		}
		if enabled == "Enabled" {
			infra.Status.PlatformStatus.AWS.CloudLoadBalancerConfig.DNSType = configv1.ClusterHostedDNSType
		}
	}
}

func (b infraBuildNamespace) withVSpherePlatformSpec() infraOption {
	return func(infra *configv1.Infrastructure) {
		if infra.Spec.PlatformSpec.VSphere != nil {
			return
		}
		infra.Spec.PlatformSpec.VSphere = &configv1.VSpherePlatformSpec{}
	}
}

func (b infraBuildNamespace) withVSpherePlatformStatus() infraOption {
	return func(infra *configv1.Infrastructure) {
		if infra.Status.PlatformStatus.VSphere != nil {
			return
		}
		infra.Status.PlatformStatus.VSphere = &configv1.VSpherePlatformStatus{}
	}
}

func (b infraBuildNamespace) withVSphereMachineNetworkEntry(cidr configv1.CIDR) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withVSpherePlatformSpec()(infra)
		b.withVSpherePlatformStatus()(infra)
		infra.Spec.PlatformSpec.VSphere.MachineNetworks = append(infra.Spec.PlatformSpec.VSphere.MachineNetworks, cidr)
		infra.Status.PlatformStatus.VSphere.MachineNetworks = append(infra.Status.PlatformStatus.VSphere.MachineNetworks, cidr)
	}
}

func (b infraBuildNamespace) withVSphereAPIVIP(vip string) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withVSpherePlatformSpec()(infra)
		b.withVSpherePlatformStatus()(infra)
		infra.Spec.PlatformSpec.VSphere.APIServerInternalIPs = append(infra.Spec.PlatformSpec.VSphere.APIServerInternalIPs, configv1.IP(vip))
		infra.Status.PlatformStatus.VSphere.APIServerInternalIPs = append(infra.Status.PlatformStatus.VSphere.APIServerInternalIPs, vip)
		if infra.Status.PlatformStatus.VSphere.APIServerInternalIP == "" {
			infra.Status.PlatformStatus.VSphere.APIServerInternalIP = vip
		}
	}
}

func (b infraBuildNamespace) withVSphereIngressVIP(vip string) infraOption {
	return func(infra *configv1.Infrastructure) {
		b.withVSpherePlatformSpec()(infra)
		b.withVSpherePlatformStatus()(infra)
		infra.Spec.PlatformSpec.VSphere.IngressIPs = append(infra.Spec.PlatformSpec.VSphere.IngressIPs, configv1.IP(vip))
		infra.Status.PlatformStatus.VSphere.IngressIPs = append(infra.Status.PlatformStatus.VSphere.IngressIPs, vip)
		if infra.Status.PlatformStatus.VSphere.IngressIP == "" {
			infra.Status.PlatformStatus.VSphere.IngressIP = vip
		}
	}
}
