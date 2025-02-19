package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestFeatureGates(t *testing.T) {
	cases := []struct {
		name          string
		installConfig *types.InstallConfig
		expected      string
	}{
		{
			name: "GCP UserProvisionedDNS is not allowed without Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.GCP = validGCPPlatform()
				c.GCP.UserProvisionedDNS = dns.UserProvisionedDNSEnabled
				return c
			}(),
			expected: `^platform.gcp.userProvisionedDNS: Forbidden: this field is protected by the GCPClusterHostedDNS feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "GCP Custom API Endpoints is not allowed without Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.GCP = validGCPPlatform()
				c.GCP.ServiceEndpoints = []gcp.ServiceEndpoint{
					{
						Name: gcp.ComputeServiceName,
						URL:  "https://compute.googleapis.com",
					},
				}
				return c
			}(),
			expected: `^platform.gcp.serviceEndpoints: Forbidden: this field is protected by the GCPCustomAPIEndpoints feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "AWS UserProvisionedDNS is not allowed without Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.AWS = validAWSPlatform()
				c.AWS.UserProvisionedDNS = dns.UserProvisionedDNSEnabled
				return c
			}(),
			expected: `^platform.aws.userProvisionedDNS: Forbidden: this field is protected by the AWSClusterHostedDNS feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "vSphere hosts is allowed with Feature Gates enabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.Default
				c.VSphere = validVSpherePlatform()
				c.VSphere.Hosts = []*vsphere.Host{{Role: "test"}}
				return c
			}(),
		},
		{
			name: "vSphere hosts is allowed with custom Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.CustomNoUpgrade
				c.FeatureGates = []string{"VSphereStaticIPs=true"}
				c.VSphere = validVSpherePlatform()
				c.VSphere.Hosts = []*vsphere.Host{{Role: "test"}}
				return c
			}(),
		},
		{
			name: "vSphere hosts is not allowed with custom Feature Gate disabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.CustomNoUpgrade
				c.FeatureGates = []string{"VSphereStaticIPs=false"}
				c.VSphere = validVSpherePlatform()
				c.VSphere.Hosts = []*vsphere.Host{{Role: "test"}}
				return c
			}(),
			expected: `^platform.vsphere.hosts: Forbidden: this field is protected by the VSphereStaticIPs feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "vSphere one vcenter is allowed with default Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.Default
				c.VSphere = validVSpherePlatform()
				c.VSphere.Hosts = []*vsphere.Host{{Role: "test"}}
				return c
			}(),
		},
		{
			name: "vSphere two vcenters is not allowed with Feature Gates disabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.CustomNoUpgrade
				c.FeatureGates = []string{"VSphereMultiVCenters=false"}
				c.VSphere = validVSpherePlatform()
				c.VSphere.VCenters = append(c.VSphere.VCenters, vsphere.VCenter{Server: "additional-vcenter"})
				return c
			}(),
			expected: `^platform.vsphere.vcenters: Forbidden: this field is protected by the VSphereMultiVCenters feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set`,
		},
		{
			name: "vSphere two vcenters is allowed with custom Feature Gate enabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.CustomNoUpgrade
				c.FeatureGates = []string{"VSphereMultiVCenters=true"}
				c.VSphere = validVSpherePlatform()
				c.VSphere.VCenters = append(c.VSphere.VCenters, vsphere.VCenter{Server: "additional-vcenter"})
				return c
			}(),
		},
		{
			name: "vSphere two vcenters is allowed with TechPreview Feature Set",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.TechPreviewNoUpgrade
				c.VSphere = validVSpherePlatform()
				c.VSphere.VCenters = append(c.VSphere.VCenters, vsphere.VCenter{Server: "Number2"})
				return c
			}(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateFeatureSet(tc.installConfig).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
