package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	utilsslice "k8s.io/utils/strings/slices"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestConvertInstallConfig(t *testing.T) {
	cases := []struct {
		name          string
		config        *types.InstallConfig
		expected      *types.InstallConfig
		expectedError string
	}{
		{
			name: "empty",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
			},
		},
		{
			name: "empty networking",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{},
			},
		},
		{
			// all deprecated fields
			name: "old networking",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1beta3",
				},
				Networking: &types.Networking{
					DeprecatedMachineCIDR: ipnet.MustParseCIDR("1.1.1.1/24"),
					DeprecatedType:        "foo",
					DeprecatedServiceCIDR: ipnet.MustParseCIDR("1.2.3.4/32"),
					DeprecatedClusterNetworks: []types.ClusterNetworkEntry{
						{
							CIDR: *ipnet.MustParseCIDR("1.2.3.5/32"),

							DeprecatedHostSubnetLength: 8,
						},
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					NetworkType: "foo",
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("1.1.1.1/24")},
					},
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/32")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR: *ipnet.MustParseCIDR("1.2.3.5/32"),

							HostPrefix:                 24,
							DeprecatedHostSubnetLength: 8,
						},
					},

					// deprecated fields are preserved
					DeprecatedType:        "foo",
					DeprecatedMachineCIDR: ipnet.MustParseCIDR("1.1.1.1/24"),
					DeprecatedServiceCIDR: ipnet.MustParseCIDR("1.2.3.4/32"),
					DeprecatedClusterNetworks: []types.ClusterNetworkEntry{
						{
							CIDR: *ipnet.MustParseCIDR("1.2.3.5/32"),

							HostPrefix:                 24,
							DeprecatedHostSubnetLength: 8,
						},
					},
				},
			},
		},
		{
			name: "new networking",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("1.1.1.1/24")},
					},
					NetworkType:    "foo",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/32")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("1.2.3.5/32"),
							HostPrefix: 24,
						},
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("1.1.1.1/24")},
					},
					NetworkType:    "foo",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/32")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("1.2.3.5/32"),
							HostPrefix: 24,
						},
					},
				},
			},
		},
		{
			name: "empty APIVersion",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "",
				},
			},
			expectedError: "no version was provided",
		},
		{
			name: "deprecated OpenShiftSDN spelling",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					NetworkType: "OpenshiftSDN",
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Networking: &types.Networking{
					NetworkType: "OpenShiftSDN",
				},
			},
		},
		{
			name: "deprecated OpenStack LbFloatingIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedLbFloatingIP: "10.0.109.1",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedLbFloatingIP: "10.0.109.1",
						APIFloatingIP:          "10.0.109.1",
					},
				},
			},
		},
		{
			name: "deprecated OpenStack LbFloatingIP is the same as APIFloatingIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedLbFloatingIP: "10.0.109.1",
						APIFloatingIP:          "10.0.109.1",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedLbFloatingIP: "10.0.109.1",
						APIFloatingIP:          "10.0.109.1",
					},
				},
			},
		},
		{
			name: "deprecated OpenStack LbFloatingIP with APIFloatingIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedLbFloatingIP: "10.0.109.1",
						APIFloatingIP:          "10.0.109.2",
					},
				},
			},
			expectedError: "cannot specify lbFloatingIP and apiFloatingIP together",
		},

		// BareMetal platform conversions
		{
			name: "baremetal external DHCP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedProvisioningDHCPExternal: true,
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedProvisioningDHCPExternal: true,
						ProvisioningNetwork:                "Unmanaged",
					},
				},
			},
		},
		{
			name: "baremetal provisioningHostIP -> clusterProvisioningIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedProvisioningHostIP: "172.22.0.3",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						ClusterProvisioningIP:        "172.22.0.3",
						DeprecatedProvisioningHostIP: "172.22.0.3",
					},
				},
			},
		},
		{
			name: "baremetal provisioningHostIP mismatch clusterProvisioningIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						ClusterProvisioningIP:        "172.22.0.4",
						DeprecatedProvisioningHostIP: "172.22.0.3",
					},
				},
			},
			expectedError: "provisioningHostIP is deprecated; only clusterProvisioningIP needs to be specified",
		},
		{
			name: "baremetal deprecated apiVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
						APIVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "baremetal deprecated ingressVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
						IngressVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "empty OpenStack platform for controlPlane",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{OpenStack: &openstack.Platform{}},
				ControlPlane: &types.MachinePool{
					Platform: types.MachinePoolPlatform{},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{OpenStack: &openstack.Platform{}},
				ControlPlane: &types.MachinePool{
					Platform: types.MachinePoolPlatform{},
				},
			},
		},
		{
			name: "empty OpenStack platform for compute",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{OpenStack: &openstack.Platform{}},
				Compute: []types.MachinePool{
					{
						Platform: types.MachinePoolPlatform{},
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{OpenStack: &openstack.Platform{}},
				Compute: []types.MachinePool{
					{
						Platform: types.MachinePoolPlatform{},
					},
				},
			},
		},
		{
			name: "deprecated OpenStack computeFlavor",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedFlavorName: "big-flavor",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedFlavorName: "big-flavor",
						DefaultMachinePlatform: &openstack.MachinePool{
							FlavorName: "big-flavor",
						},
					},
				},
			},
		},
		{
			name: "deprecated OpenStack computeFlavor with type in defaultMachinePlatform",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedFlavorName: "flavor1",
						DefaultMachinePlatform: &openstack.MachinePool{
							FlavorName: "flavor2",
						},
					},
				},
			},
			expectedError: "cannot specify computeFlavor and type in defaultMachinePlatform together",
		},
		{
			name: "deprecated OpenStack controlPlane with type in rootVolume",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{OpenStack: &openstack.Platform{}},
				ControlPlane: &types.MachinePool{
					Platform: types.MachinePoolPlatform{
						OpenStack: &openstack.MachinePool{
							RootVolume: &openstack.RootVolume{
								DeprecatedType: "fast",
							},
						},
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{OpenStack: &openstack.Platform{}},
				ControlPlane: &types.MachinePool{
					Platform: types.MachinePoolPlatform{
						OpenStack: &openstack.MachinePool{
							RootVolume: &openstack.RootVolume{
								Types: []string{"fast"},
							},
						},
					},
				},
			},
		},
		{
			name: "openstack deprecated apiVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
						APIVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "openstack deprecated ingressVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
						IngressVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "vsphere deprecated apiVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					VSphere: &vsphere.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					VSphere: &vsphere.Platform{
						VCenters: []vsphere.VCenter{{
							Server:      "",
							Port:        443,
							Username:    "",
							Password:    "",
							Datacenters: nil,
						}},
						FailureDomains: []vsphere.FailureDomain{{
							Name:   "generated-failure-domain",
							Region: "generated-region",
							Zone:   "generated-zone",
							Server: "",
							Topology: vsphere.Topology{
								Datacenter:     "",
								ComputeCluster: "",
								Networks:       []string{""},
								Datastore:      "",
								ResourcePool:   "",
								Folder:         "",
							},
						}},
						DeprecatedAPIVIP: "1.2.3.4",
						APIVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "vsphere deprecated ingressVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					VSphere: &vsphere.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					VSphere: &vsphere.Platform{
						VCenters: []vsphere.VCenter{{
							Server:      "",
							Port:        443,
							Username:    "",
							Password:    "",
							Datacenters: nil,
						}},
						FailureDomains: []vsphere.FailureDomain{{
							Name:   "generated-failure-domain",
							Region: "generated-region",
							Zone:   "generated-zone",
							Server: "",
							Topology: vsphere.Topology{
								Datacenter:     "",
								ComputeCluster: "",
								Networks:       []string{""},
								Datastore:      "",
								ResourcePool:   "",
								Folder:         "",
							},
						}},
						DeprecatedIngressVIP: "1.2.3.4",
						IngressVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "ovirt deprecated apiVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					Ovirt: &ovirt.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					Ovirt: &ovirt.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
						APIVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "ovirt deprecated ingressVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					Ovirt: &ovirt.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					Ovirt: &ovirt.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
						IngressVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "nutanix deprecated apiVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					Nutanix: &nutanix.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					Nutanix: &nutanix.Platform{
						DeprecatedAPIVIP: "1.2.3.4",
						APIVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "nutanix deprecated ingressVIP",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					Nutanix: &nutanix.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					Nutanix: &nutanix.Platform{
						DeprecatedIngressVIP: "1.2.3.4",
						IngressVIPs:          []string{"1.2.3.4"},
					},
				},
			},
		},
		{
			name: "aws deprecated platform amiID",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					AWS: &aws.Platform{
						AMIID: "deprec-id",
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					AWS: &aws.Platform{
						AMIID: "deprec-id",
						DefaultMachinePlatform: &aws.MachinePool{
							AMIID: "deprec-id",
						},
					},
				},
			},
		},
		{
			name: "aws deprecated subnets",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					AWS: &aws.Platform{
						DeprecatedSubnets: []string{"subnet-01234567890abcdef", "subnet-abcdef01234567890"},
					},
				},
			},
			expected: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					AWS: &aws.Platform{
						DeprecatedSubnets: []string{"subnet-01234567890abcdef", "subnet-abcdef01234567890"},
						VPC: aws.VPC{
							Subnets: []aws.Subnet{
								{ID: "subnet-01234567890abcdef"},
								{ID: "subnet-abcdef01234567890"},
							},
						},
					},
				},
			},
		},
		{
			name: "aws deprecated subnets with vpc.subnets",
			config: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				Platform: types.Platform{
					AWS: &aws.Platform{
						DeprecatedSubnets: []string{"subnet-01234567890abcdef", "subnet-abcdef01234567890"},
						VPC: aws.VPC{
							Subnets: []aws.Subnet{
								{ID: "subnet-01234567890abcdef"},
								{ID: "subnet-abcdef01234567890"},
							},
						},
					},
				},
			},
			expectedError: `Forbidden: cannot specify platform\.aws.subnets and platform\.aws\.vpc.subnets together`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ConvertInstallConfig(tc.config)
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, tc.config, "unexpected install config")
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}

func Test_upconvertVIPs(t *testing.T) {
	tests := []struct {
		name     string
		vips     []string
		oldVIP   string
		wantErr  bool
		wantVIPs []string
	}{
		{
			name:     "should return error if both fields are set and all are unique",
			vips:     []string{"foo", "bar"},
			oldVIP:   "baz",
			wantErr:  true,
			wantVIPs: []string{},
		},
		{
			name:     "should set VIPs if old VIPs is set",
			vips:     []string{},
			oldVIP:   "baz",
			wantErr:  false,
			wantVIPs: []string{"baz"},
		},
		{
			name:     "should return VIPs if only new VIPs is set",
			vips:     []string{"foo", "bar"},
			oldVIP:   "",
			wantErr:  false,
			wantVIPs: []string{"foo", "bar"},
		},
		{
			name:     "should return no error if old VIP is contained in new VIPs",
			vips:     []string{"foo", "bar"},
			oldVIP:   "bar",
			wantErr:  false,
			wantVIPs: []string{"foo", "bar"},
		},
		{
			name:     "should not fail on nil",
			vips:     nil,
			oldVIP:   "",
			wantErr:  false,
			wantVIPs: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := upconvertVIP(&tt.vips, tt.oldVIP, "new", "old", field.NewPath("test")); (err != nil) != tt.wantErr {
				t.Errorf("upconvertVIP() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if !tt.wantErr {
					for _, wantVIP := range tt.wantVIPs {
						if !utilsslice.Contains(tt.vips, wantVIP) {
							t.Errorf("upconvertVIP() didn't update VIPs field (expected \"%v\" to contain \"%s\")", tt.vips, wantVIP)
						}
					}
				}
			}
		})
	}
}
