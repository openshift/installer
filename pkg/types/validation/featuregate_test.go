package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
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
			name: "GCP UserTags is allowed with Feature Gates enabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.TechPreviewNoUpgrade
				c.GCP = validGCPPlatform()
				c.GCP.UserTags = []gcp.UserTag{{ParentID: "a", Key: "b", Value: "c"}}
				return c
			}(),
		},
		{
			name: "GCP UserTags is not allowed without Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.GCP = validGCPPlatform()
				c.GCP.UserTags = []gcp.UserTag{{ParentID: "a", Key: "b", Value: "c"}}
				return c
			}(),
			expected: `^platform.gcp.userTags: Forbidden: this field is protected by the GCPLabelsTags feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "GCP UserLabels is allowed with Feature Gates enabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.TechPreviewNoUpgrade
				c.GCP = validGCPPlatform()
				c.GCP.UserLabels = []gcp.UserLabel{{Key: "a", Value: "b"}}
				return c
			}(),
		},
		{
			name: "GCP UserProvisionedDNS is not allowed without Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.GCP = validGCPPlatform()
				c.GCP.UserProvisionedDNS = gcp.UserProvisionedDNSEnabled
				return c
			}(),
			expected: `^platform.gcp.userProvisionedDNS: Forbidden: this field is protected by the GCPClusterHostedDNS feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "GCP UserLabels is not allowed without Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.GCP = validGCPPlatform()
				c.GCP.UserLabels = []gcp.UserLabel{{Key: "a", Value: "b"}}
				return c
			}(),
			expected: `^platform.gcp.userLabels: Forbidden: this field is protected by the GCPLabelsTags feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
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
