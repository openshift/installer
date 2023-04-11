package validation

import (
	"fmt"
	"k8s.io/utils/pointer"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenters: []vsphere.VCenter{
			{
				Server:   "test-vcenter",
				Port:     443,
				Username: "test-username",
				Password: "test-password",
				Datacenters: []string{
					"test-datacenter",
				},
			},
		},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name:   "test-east-1a",
				Region: "test-east",
				Zone:   "test-east-1a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/test-cluster",
					Datastore:      "/test-datacenter/datastore/test-datastore",
					Networks:       []string{"test-portgroup"},
					ResourcePool:   "/test-datacenter/host/test-cluster/Resources/test-resourcepool",
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
			{
				Name:   "test-east-2a",
				Region: "test-east",
				Zone:   "test-east-2a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/test-cluster",
					Datastore:      "/test-datacenter/datastore/test-datastore",
					Networks:       []string{"test-portgroup"},
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
		},
	}
}

func validHosts() []*vsphere.Host {
	return []*vsphere.Host{
		{
			Role: "bootstrap",
			NetworkDevice: &vsphere.NetworkDeviceSpec{
				IPAddrs: []string{
					"192.168.101.240/24",
				},
				Gateway4: "192.168.101.1",
				Nameservers: []string{
					"192.168.101.2",
				},
			},
		},
		{
			Role:          "control-plane",
			FailureDomain: "test-east-1a",
			NetworkDevice: &vsphere.NetworkDeviceSpec{
				IPAddrs: []string{
					"192.168.101.241/24",
				},
				Gateway4: "192.168.101.1",
				Nameservers: []string{
					"192.168.101.2",
				},
			},
		},
		{
			Role:          "control-plane",
			FailureDomain: "test-east-2a",
			NetworkDevice: &vsphere.NetworkDeviceSpec{
				IPAddrs: []string{
					"192.168.101.242/24",
				},
				Gateway4: "192.168.101.1",
				Nameservers: []string{
					"192.168.101.2",
				},
			},
		},
		{
			Role:          "control-plane",
			FailureDomain: "test-east-1a",
			NetworkDevice: &vsphere.NetworkDeviceSpec{
				IPAddrs: []string{
					"192.168.101.243/24",
				},
				Gateway4: "192.168.101.1",
				Nameservers: []string{
					"192.168.101.2",
				},
			},
		},
		{
			Role:          "compute",
			FailureDomain: "test-east-1a",
			NetworkDevice: &vsphere.NetworkDeviceSpec{
				IPAddrs: []string{
					"192.168.101.244/24",
				},
				Gateway4: "192.168.101.1",
				Nameservers: []string{
					"192.168.101.2",
				},
			},
		},
		{
			Role:          "compute",
			FailureDomain: "test-east-2a",
			NetworkDevice: &vsphere.NetworkDeviceSpec{
				IPAddrs: []string{
					"192.168.101.245/24",
				},
				Gateway4: "192.168.101.1",
				Nameservers: []string{
					"192.168.101.2",
				},
			},
		},
		{
			Role:          "compute",
			FailureDomain: "test-east-1a",
			NetworkDevice: &vsphere.NetworkDeviceSpec{
				IPAddrs: []string{
					"192.168.101.246/24",
				},
				Gateway4: "192.168.101.1",
				Nameservers: []string{
					"192.168.101.2",
				},
			},
		},
	}
}

func validStaticIPInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		FeatureSet: configv1.TechPreviewNoUpgrade,
		ControlPlane: &types.MachinePool{
			Name:     "master",
			Replicas: pointer.Int64(3),
		},
		Compute: []types.MachinePool{
			{
				Name:     "worker",
				Replicas: pointer.Int64(3),
			},
		},
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		config        *types.InstallConfig
		platform      *vsphere.Platform
		expectedError string
	}{
		{
			name: "Valid diskType",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.DiskType = "eagerZeroedThick"
				return p
			}(),
		},
		{
			name: "Invalid diskType",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.DiskType = "invalidDiskType"
				return p
			}(),
			expectedError: `^test-path\.diskType: Invalid value: "invalidDiskType": diskType must be one of \[eagerZeroedThick thick thin\]$`,
		},

		{
			name:     "Valid Multi-zone platform",
			platform: validPlatform(),
		},
		{
			name: "Multi-zone platform missing failureDomains",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains = make([]vsphere.FailureDomain, 0)
				return p
			}(),
			expectedError: `^test-path.failureDomains: Required value: must be defined`,
		},
		{
			name: "Multi-zone platform vCenter missing server",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Server = ""
				return p
			}(),
			expectedError: `test-path\.vcenters\.server: Required value: must be the domain name or IP address of the vCenter(.*)`,
		},
		{
			name: "Multi-zone platform more than one vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters = append(p.VCenters, vsphere.VCenter{
					Server: "additional-vcenter",
				})
				return p
			}(),
			expectedError: `^test-path\.vcenters: Too many: 2: must have at most 1 items`,
		},
		{
			name: "Multi-zone platform Capital letters in vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Server = "tEsT-vCenter"
				return p
			}(),
			expectedError: `(.*)test-path\.vcenters.server: Invalid value: "tEsT-vCenter": must be the domain name or IP address of the vCenter`,
		},
		{
			name: "Multi-zone missing username",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Username = ""
				return p
			}(),
			expectedError: `^test-path\.vcenters.username: Required value: must specify the username$`,
		},
		{
			name: "Multi-zone missing password",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Password = ""
				return p
			}(),
			expectedError: `^test-path\.vcenters.password: Required value: must specify the password$`,
		},
		{
			name: "Multi-zone missing datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Datacenters = []string{}
				return p
			}(),
			expectedError: `^test-path\.vcenters.datacenters: Required value: must specify at least one datacenter$`,
		},
		{
			name: "Multi-zone platform wrong vCenter name in failureDomain zone",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Server = "bad-vcenter"
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.server: Invalid value: "bad-vcenter": server does not exist in vcenters`,
		},
		{
			name: "Multi-zone platform failure domain topology cluster relative path",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Topology.ComputeCluster = "incomplete-path"
				p.FailureDomains[0].Topology.ResourcePool = "/test-datacenter/host/incomplete-path/Resources/test-resourcepool"
				return p
			}(),
			expectedError: `(.*)test-path\.failureDomains\.topology\.computeCluster: Invalid value: "incomplete-path": full path of compute cluster must be provided in format /<datacenter>/host/<cluster>`,
		},
		{
			name: "Multi-zone platform datacenter in failure domain topology doesn't match cluster datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Topology.ComputeCluster = "/other-datacenter/host/cluster"
				return p
			}(),
			expectedError: `^test-path.failureDomains.topology.computeCluster: Invalid value: "/other-datacenter/host/cluster": compute cluster must be in datacenter test-datacenter`,
		},
		{
			name: "Multi-zone platform failureDomain missing name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Name = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform failureDomain region missing name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Region = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.region: Required value: must specify region tag value`,
		},
		{
			name: "Multi-zone platform failureDomain zone missing name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Name = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform failureDomain zone missing tag category",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Zone = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.zone: Required value: must specify zone tag value`,
		},
		{
			name:     "forbidden load balancer field",
			platform: validPlatform(),
			config: &types.InstallConfig{
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.VSpherePlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
						}
						return p
					}(),
				},
			},
			expectedError: `^test-path\.loadBalancer: Forbidden: load balancer is not supported in this feature set`,
		},
		{
			name:     "allowed load balancer field with OpenShift managed default",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.VSpherePlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
						}
						return p
					}(),
				},
			},
		},
		{
			name:     "allowed load balancer field with user-managed",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.VSpherePlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeUserManaged,
						}
						return p
					}(),
				},
			},
		},
		{
			name:     "allowed load balancer field invalid type",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.VSpherePlatformLoadBalancer{
							Type: "FooBar",
						}
						return p
					}(),
				},
			},
			expectedError: `^test-path\.loadBalancer.type: Invalid value: "FooBar": invalid load balancer type`,
		},
		{
			name: "Static IP - valid",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				return p
			}(),
			config: validStaticIPInstallConfig(),
		},
		{
			name: "Static IP - no hosts configured",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				return p
			}(),
			config: validStaticIPInstallConfig(),
		},
		{
			name: "Static IP - invalid Role",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].Role = "crazy-uncle"
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `test-path.hosts.role: Invalid value: "crazy-uncle": role must be one of \[bootstrap compute control-plane]`,
		},
		{
			name: "Static IP - invalid FailureDomain",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].FailureDomain = "north-pole"
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.failureDomain: Invalid value: "north-pole": failure domain not found$`,
		},
		{
			name: "Static IP - missing NetworkDevice",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice = nil
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.networkDevice: Required value: must specify networkDevice configuration$`,
		},
		{
			name: "Static IP - missing IP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.IPAddrs = nil
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.ipAddrs: Required value: must specify a IP$`,
		},
		{
			name: "Static IP - invalid IP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.IPAddrs[0] = "86.7.5.309/24"
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.ipAddrs: Invalid value: "86.7.5.309/24": invalid CIDR address: 86.7.5.309/24$`,
		},
		{
			name: "Static IP - invalid IP blank",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.IPAddrs[0] = ""
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.ipAddrs: Required value: must specify a IP address with CIDR$`,
		},
		{
			name: "Static IP - invalid IP CIDR",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.IPAddrs[0] = "86.7.5.309/55"
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.ipAddrs: Invalid value: "86.7.5.309/55": invalid CIDR address: 86.7.5.309/55$`,
		},
		{
			name: "Static IP - invalid IP missing CIDR",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.IPAddrs[0] = "86.7.5.309"
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.ipAddrs: Invalid value: "86.7.5.309": invalid CIDR address: 86.7.5.309$`,
		},
		{
			name: "Static IP - valid Gateway4 IP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.Gateway4 = "192.168.100.125"
				return p
			}(),
			config: validStaticIPInstallConfig(),
		},
		{
			name: "Static IP - invalid Gateway4 IP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.Gateway4 = "86.7.5.309"
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.gateway4: Invalid value: "86.7.5.309": "86.7.5.309" is not a valid IP$`,
		},
		{
			name: "Static IP - valid Gateway6 IP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.Gateway6 = "2001:db8:3333:4444:5555:6666:7777:8888"
				return p
			}(),
			config: validStaticIPInstallConfig(),
		},
		{
			name: "Static IP - invalid Gateway6 IP",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.Gateway6 = "8888:666:7777:5555:3333:0000:9999:JENNY"
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.gateway6: Invalid value: "8888:666:7777:5555:3333:0000:9999:JENNY": "8888:666:7777:5555:3333:0000:9999:JENNY" is not a valid IP$`,
		},
		{
			name: "Static IP - More than 3 nameservers",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				p.Hosts[1].NetworkDevice.Nameservers = []string{"86.75.30.9", "86.75.30.8", "86.75.30.7", "86.75.30.6"}
				return p
			}(),
			config:        validStaticIPInstallConfig(),
			expectedError: `^test-path.hosts.nameservers: Too many: 4: must have at most 3 items$`,
		},
		{
			name: "Static IP - Not enough control-planes",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				return p
			}(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.Hosts = validHosts()
						return p
					}(),
				},
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64(4),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: pointer.Int64(3),
					},
				},
			},
			expectedError: `^test-path.hosts: Invalid value: "control-plane": not enough hosts found \(3\) to support all the configured ControlPlane replicas \(4\)$`,
		},
		{
			name: "Static IP - Too many control-planes",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				return p
			}(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.Hosts = validHosts()
						return p
					}(),
				},
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64(2),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: pointer.Int64(3),
					},
				},
			},
		},
		{
			name: "Static IP - Not enough workers",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				return p
			}(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.Hosts = validHosts()
						return p
					}(),
				},
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64(3),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: pointer.Int64(4),
					},
				},
			},
			expectedError: `^test-path.hosts: Invalid value: "control-plane": not enough hosts found \(3\) to support all the configured Compute replicas \(4\)$`,
		},
		{
			name: "Static IP - Too many workers",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				return p
			}(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.Hosts = validHosts()
						return p
					}(),
				},
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64(3),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: pointer.Int64(2),
					},
				},
			},
		},
		{
			name: "Static IP - Not enough control-plane and workers",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Hosts = validHosts()
				return p
			}(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.Hosts = validHosts()
						return p
					}(),
				},
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64(4),
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker",
						Replicas: pointer.Int64(4),
					},
				},
			},
			expectedError: `^\[test-path.hosts: Invalid value: "control-plane": not enough hosts found \(3\) to support all the configured ControlPlane replicas \(4\), test-path.hosts: Invalid value: "control-plane": not enough hosts found \(3\) to support all the configured Compute replicas \(4\)]$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
				tc.config.VSphere = tc.platform
			}
			if tc.config.VSphere == nil {
				fmt.Printf("Setting vsphere: %v", tc.platform)
				tc.config.VSphere = tc.platform
			}
			err := ValidatePlatform(tc.config.VSphere, false, field.NewPath("test-path"), tc.config).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, regexp.MustCompile(tc.expectedError), err)
			}
		})
	}
}

type installConfigBuilder struct {
	types.InstallConfig
}

func installConfig() *installConfigBuilder {
	return &installConfigBuilder{
		InstallConfig: types.InstallConfig{},
	}
}

func (icb *installConfigBuilder) build() *types.InstallConfig {
	return &icb.InstallConfig
}
