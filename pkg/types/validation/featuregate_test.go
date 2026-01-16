package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestFeatureGates(t *testing.T) {
	cases := []struct {
		name          string
		installConfig *types.InstallConfig
		expected      string
	}{
		{
			name: "AWS UserProvisionedDNS is not allowed without Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.AWS = validAWSPlatform()
				c.AWS.UserProvisionedDNS = dns.UserProvisionedDNSEnabled
				return c
			}(),
			expected: `^platform.aws.userProvisionedDNS: Forbidden: this field is protected by the AWSClusterHostedDNSInstall feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "Azure UserProvisionedDNS is not allowed without Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.AWS = nil // validInstallConfig defaults to AWS
				c.Azure = &azure.Platform{}
				c.Azure.UserProvisionedDNS = dns.UserProvisionedDNSEnabled
				return c
			}(),
			expected: `^platform.azure.userProvisionedDNS: Forbidden: this field is protected by the AzureClusterHostedDNSInstall feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "vSphere hosts is allowed with Feature Gates enabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.Default
				c.AWS = nil // validInstallConfig defaults to AWS
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
				c.VSphere = validVSpherePlatform()
				c.VSphere.Hosts = []*vsphere.Host{{Role: "test"}}
				return c
			}(),
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
			name: "vSphere two vcenters is allowed with default Feature Gates",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.Default
				c.VSphere = validVSpherePlatform()
				c.VSphere.VCenters = append(c.VSphere.VCenters, vsphere.VCenter{Server: "additional-vcenter"})
				return c
			}(),
		},
		{
			name: "vSphere two vcenters is allowed with custom Feature Gate enabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.CustomNoUpgrade
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
		{
			name: "Azure user-assigned identities (control plane) > 1 requires MachineAPIMigration feature gate",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.AWS = nil // validInstallConfig defaults to AWS
				c.Azure = &azure.Platform{}
				c.ControlPlane.Platform.Azure = &azure.MachinePool{
					Identity: &azure.VMIdentity{
						Type: azure.VMIdentityUserAssigned,
						UserAssignedIdentities: []azure.UserAssignedIdentity{
							{
								Name:          "first-identity",
								Subscription:  "my-subscription",
								ResourceGroup: "my-resource-group",
							},
							{
								Name:          "second-identity",
								Subscription:  "my-subscription",
								ResourceGroup: "my-resource-group",
							},
						},
					},
				}
				return c
			}(),
			expected: `^controlPlane.azure.identity.userAssignedIdentities: Forbidden: this field is protected by the MachineAPIMigration feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set`,
		},
		{
			name: "Azure user-assigned identities (default machine platform) > 1 requires MachineAPIMigration feature gate",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.AWS = nil // validInstallConfig defaults to AWS
				c.Azure = &azure.Platform{}
				c.Azure.DefaultMachinePlatform = &azure.MachinePool{
					Identity: &azure.VMIdentity{
						Type: azure.VMIdentityUserAssigned,
						UserAssignedIdentities: []azure.UserAssignedIdentity{
							{
								Name:          "first-identity",
								Subscription:  "my-subscription",
								ResourceGroup: "my-resource-group",
							},
							{
								Name:          "second-identity",
								Subscription:  "my-subscription",
								ResourceGroup: "my-resource-group",
							},
						},
					},
				}
				return c
			}(),
			expected: `^platform.azure.defaultMachinePlatform.identity.userAssignedIdentities: Forbidden: this field is protected by the MachineAPIMigration feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set`,
		},
		{
			name: "Azure user-assigned identities (control plane) == 1 does not require feature gate",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.AWS = nil // validInstallConfig defaults to AWS
				c.Azure = &azure.Platform{}
				c.ControlPlane.Platform.Azure = &azure.MachinePool{
					Identity: &azure.VMIdentity{
						Type: azure.VMIdentityUserAssigned,
						UserAssignedIdentities: []azure.UserAssignedIdentity{
							{
								Name:          "solo-bolo!",
								Subscription:  "my-subscription",
								ResourceGroup: "my-resource-group",
							},
						},
					},
				}
				return c
			}(),
		},
		{
			name: "FencingCredentials is not allowed with Feature Gates disabled",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.ControlPlane.Fencing = &types.Fencing{Credentials: []*types.Credential{{HostName: "host1"}, {HostName: "host2"}}}
				return c
			}(),
			expected: `^platform.none.fencingCredentials: Forbidden: this field is protected by the DualReplica feature gate which must be enabled through either the TechPreviewNoUpgrade or CustomNoUpgrade feature set$`,
		},
		{
			name: "FencingCredentials is allowed with DevPreviewNoUpgrade Feature Set",
			installConfig: func() *types.InstallConfig {
				c := validInstallConfig()
				c.FeatureSet = v1.DevPreviewNoUpgrade
				c.ControlPlane.Fencing = &types.Fencing{Credentials: []*types.Credential{{HostName: "host1"}, {HostName: "host2"}}}
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
