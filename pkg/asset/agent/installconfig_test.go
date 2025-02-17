package agent

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func TestInstallConfigLoad(t *testing.T) {
	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  string
		expectedConfig *types.InstallConfig
	}{
		{
			name: "unsupported platform",
			data: `
apiVersion: v1
metadata:
    name: test-cluster
baseDomain: test-domain
platform:
  aws:
    region: us-east-1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: Platform: Unsupported value: "aws": supported values: "baremetal", "vsphere", "none", "external"`,
		},
		{
			name: "apiVips not set for baremetal Compact platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  networkType: OpenShiftSDN
  machineNetwork:
  - cidr: 192.168.122.0/23
  serviceNetwork:
  - 172.30.0.0/16
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 3
platform:
  baremetal:
    externalMACAddress: "52:54:00:f6:b4:02"
    provisioningMACAddress: "52:54:00:6e:3b:02"
    ingressVIPs:
      - 192.168.122.11
    hosts:
      - name: host1
        bootMACAddress: 52:54:01:aa:aa:a1
      - name: host2
        bootMACAddress: 52:54:01:bb:bb:b1
      - name: host3
        bootMACAddress: 52:54:01:cc:cc:c1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: [platform.baremetal.apiVIPs: Required value: must specify at least one VIP for the API, platform.baremetal.apiVIPs: Required value: must specify VIP for API, when VIP for ingress is set]",
		},
		{
			name: "Required values not set for vsphere platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OpenShiftSDN
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: [platform.vsphere.ingressVIPs: Required value: must specify VIP for ingress, when VIP for API is set, platform.vsphere.ingressVIPs: Required value: must specify at least one VIP for the Ingress]",
		},
		{
			name: "no compute.replicas set for SNO",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 1
platform:
  none : {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: Compute.Replicas: Required value: Total number of Compute.Replicas must be 0 when ControlPlane.Replicas is 1 for platform none or external. Found 3",
		},
		{
			name: "incorrect compute.replicas set",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 2
platform:
  external:
    platformName: oci
    cloudControllerManager: External
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: ControlPlane.Replicas: Invalid value: 2: ControlPlane.Replicas can only be set to 3 or 1. Found 2",
		},
		{
			name: "invalid networkType for SNO cluster",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OpenShiftSDN
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 1
platform:
  none : {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: Networking.NetworkType: Invalid value: \"OpenShiftSDN\": Only OVNKubernetes network type is allowed for Single Node OpenShift (SNO) cluster",
		},
		{
			name: "invalid platform for SNO cluster",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OpenShiftSDN
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 1
platform:
  aws:
    region: us-east-1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: [Platform: Unsupported value: \"aws\": supported values: \"baremetal\", \"vsphere\", \"none\", \"external\", Platform: Invalid value: \"aws\": Only platform none and external supports 1 ControlPlane and 0 Compute nodes]",
		},
		{
			name: "invalid architecture for SNO cluster",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
compute:
  - architecture: s390x
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
controlPlane:
  architecture: s390x
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 1
platform:
  none : {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: [ControlPlane.Architecture: Unsupported value: \"s390x\": supported values: \"amd64\", \"arm64\", Compute[0].Architecture: Unsupported value: \"s390x\": supported values: \"amd64\", \"arm64\"]",
		},
		{
			name: "unsupported platformName for external platform",
			data: `
apiVersion: v1
metadata:
    name: test-cluster
baseDomain: test-domain
platform:
  external:
    platformName: some-cloud-provider
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: Platform.External.PlatformName: Unsupported value: "some-cloud-provider": supported values: "oci"`,
		},
		{
			name: "unsupported CloudControllerManager for external platform",
			data: `
apiVersion: v1
metadata:
    name: test-cluster
baseDomain: test-domain
platform:
  external:
    platformName: oci
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: Platform.External.CloudControllerManager: Invalid value: "": When using external oci platform, Platform.External.CloudControllerManager must be set to External`,
		},
		{
			name: "valid configuration for none platform for sno",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 1
platform:
  none : {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: true,
			expectedConfig: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				AdditionalTrustBundlePolicy: types.PolicyProxyOnly,
				BaseDomain:                  "test-domain",
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					},
					NetworkType:    "OVNKubernetes",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("172.30.0.0/16")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("10.128.0.0/14"),
							HostPrefix: 23,
						},
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       pointer.Int64(1),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       pointer.Int64(0),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform:   types.Platform{None: &none.Platform{}},
				PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
		{
			name: "valid configuration for none platform for HA cluster",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 2
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 3
platform:
  none : {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: true,
			expectedConfig: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				AdditionalTrustBundlePolicy: types.PolicyProxyOnly,
				BaseDomain:                  "test-domain",
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
					},
					NetworkType:    "OVNKubernetes",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("172.30.0.0/16")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("10.128.0.0/14"),
							HostPrefix: 23,
						},
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       pointer.Int64(3),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       pointer.Int64(2),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform:   types.Platform{None: &none.Platform{}},
				PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
		{
			name: "valid configuration for baremetal platform for HA cluster - deprecated and unused fields",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  networkType: OpenShiftSDN
  machineNetwork:
  - cidr: 192.168.122.0/23
  serviceNetwork:
  - 172.30.0.0/16
compute:
  - architecture: amd64
    hyperthreading: Disabled
    name: worker
    platform: {}
    replicas: 2
controlPlane:
  architecture: amd64
  hyperthreading: Disabled
  name: master
  platform: {}
  replicas: 3
platform:
  baremetal:
    libvirtURI: qemu+ssh://root@52.116.73.24/system
    clusterProvisioningIP: "192.168.122.90"
    bootstrapProvisioningIP: "192.168.122.91"
    externalBridge: "somevalue"
    externalMACAddress: "52:54:00:f6:b4:02"
    provisioningNetwork: "Disabled"
    provisioningBridge: br0
    provisioningMACAddress: "52:54:00:6e:3b:02"
    provisioningNetworkInterface: "eth11"
    provisioningDHCPExternal: true
    provisioningDHCPRange: 172.22.0.10,172.22.0.254
    apiVIP: 192.168.122.10
    ingressVIP: 192.168.122.11
    bootstrapOSImage: https://mirror.example.com/images/qemu.qcow2.gz?sha256=a07bd
    clusterOSImage: https://mirror.example.com/images/metal.qcow2.gz?sha256=3b5a8
    bootstrapExternalStaticIP: 192.1168.122.50
    bootstrapExternalStaticGateway: gateway
    hosts:
      - name: host1
        bootMACAddress: 52:54:01:aa:aa:a1
        bmc:
          address: addr
      - name: host2
        bootMACAddress: 52:54:01:bb:bb:b1
      - name: host3
        bootMACAddress: 52:54:01:cc:cc:c1
      - name: host4
        bootMACAddress: 52:54:01:dd:dd:d1
      - name: host5
        bootMACAddress: 52:54:01:ee:ee:e1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: true,
			expectedConfig: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				AdditionalTrustBundlePolicy: types.PolicyProxyOnly,
				BaseDomain:                  "test-domain",
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("192.168.122.0/23")},
					},
					NetworkType:    "OpenShiftSDN",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("172.30.0.0/16")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("10.128.0.0/14"),
							HostPrefix: 23,
						},
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       pointer.Int64(3),
					Hyperthreading: types.HyperthreadingDisabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       pointer.Int64(2),
						Hyperthreading: types.HyperthreadingDisabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						LibvirtURI:                         "qemu+ssh://root@52.116.73.24/system",
						ClusterProvisioningIP:              "192.168.122.90",
						BootstrapProvisioningIP:            "192.168.122.91",
						ExternalBridge:                     "somevalue",
						ExternalMACAddress:                 "52:54:00:f6:b4:02",
						ProvisioningNetwork:                "Disabled",
						ProvisioningBridge:                 "br0",
						ProvisioningMACAddress:             "52:54:00:6e:3b:02",
						ProvisioningDHCPRange:              "172.22.0.10,172.22.0.254",
						DeprecatedProvisioningDHCPExternal: true,
						ProvisioningNetworkCIDR: &ipnet.IPNet{
							IPNet: net.IPNet{
								IP:   []byte("\xc0\xa8\x7a\x00"),
								Mask: []byte("\xff\xff\xfe\x00"),
							},
						},
						ProvisioningNetworkInterface: "eth11",
						Hosts: []*baremetal.Host{
							{
								Name:            "host1",
								BootMACAddress:  "52:54:01:aa:aa:a1",
								BootMode:        "UEFI",
								HardwareProfile: "default",
								BMC:             baremetal.BMC{Address: "addr"},
							},
							{
								Name:            "host2",
								BootMACAddress:  "52:54:01:bb:bb:b1",
								BootMode:        "UEFI",
								HardwareProfile: "default",
							},
							{
								Name:            "host3",
								BootMACAddress:  "52:54:01:cc:cc:c1",
								BootMode:        "UEFI",
								HardwareProfile: "default",
							},
							{
								Name:            "host4",
								BootMACAddress:  "52:54:01:dd:dd:d1",
								BootMode:        "UEFI",
								HardwareProfile: "default",
							},
							{
								Name:            "host5",
								BootMACAddress:  "52:54:01:ee:ee:e1",
								BootMode:        "UEFI",
								HardwareProfile: "default",
							}},
						DeprecatedAPIVIP:               "192.168.122.10",
						APIVIPs:                        []string{"192.168.122.10"},
						DeprecatedIngressVIP:           "192.168.122.11",
						IngressVIPs:                    []string{"192.168.122.11"},
						BootstrapOSImage:               "https://mirror.example.com/images/qemu.qcow2.gz?sha256=a07bd",
						ClusterOSImage:                 "https://mirror.example.com/images/metal.qcow2.gz?sha256=3b5a8",
						BootstrapExternalStaticIP:      "192.1168.122.50",
						BootstrapExternalStaticGateway: "gateway",
					},
				},
				PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
		{
			name: "valid configuration for vsphere platform for compact cluster - deprecated field apiVip",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  networkType: OpenShiftSDN
  machineNetwork:
  - cidr: 192.168.122.0/23
  serviceNetwork: 
  - 172.30.0.0/16
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 3
platform:
  vsphere :
    vcenter: 192.168.122.30
    username: testUsername
    password: testPassword
    datacenter: testDataCenter
    defaultDataStore: testDefaultDataStore
    cluster: testCluster
    apiVIP: 192.168.122.10
    ingressVIPs: 
      - 192.168.122.11
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: true,
			expectedConfig: &types.InstallConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				AdditionalTrustBundlePolicy: types.PolicyProxyOnly,
				BaseDomain:                  "test-domain",
				Networking: &types.Networking{
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("192.168.122.0/23")},
					},
					NetworkType:    "OpenShiftSDN",
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("172.30.0.0/16")},
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       *ipnet.MustParseCIDR("10.128.0.0/14"),
							HostPrefix: 23,
						},
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       pointer.Int64(3),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       pointer.Int64(0),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform: types.Platform{
					VSphere: &vsphere.Platform{
						DeprecatedVCenter:          "192.168.122.30",
						DeprecatedUsername:         "testUsername",
						DeprecatedPassword:         "testPassword",
						DeprecatedDatacenter:       "testDataCenter",
						DeprecatedCluster:          "testCluster",
						DeprecatedDefaultDatastore: "testDefaultDataStore",
						DeprecatedAPIVIP:           "192.168.122.10",
						APIVIPs:                    []string{"192.168.122.10"},
						IngressVIPs:                []string{"192.168.122.11"},
						VCenters: []vsphere.VCenter{{
							Server:      "192.168.122.30",
							Port:        443,
							Username:    "testUsername",
							Password:    "testPassword",
							Datacenters: []string{"testDataCenter"},
						}},
						FailureDomains: []vsphere.FailureDomain{{
							Name:   "generated-failure-domain",
							Region: "generated-region",
							Zone:   "generated-zone",
							Server: "192.168.122.30",
							Topology: vsphere.Topology{
								Datacenter:     "testDataCenter",
								ComputeCluster: "/testDataCenter/host/testCluster",
								Networks:       []string{""},
								Datastore:      "/testDataCenter/datastore/testDefaultDataStore",
								ResourcePool:   "/testDataCenter/host/testCluster//Resources",
								Folder:         "",
							},
						}},
					},
				},
				PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(installConfigFilename).
				Return(
					&asset.File{
						Filename: installConfigFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				).MaxTimes(2)

			asset := &OptionalInstallConfig{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, asset.Config, "unexpected Config in InstallConfig")
			}
		})
	}
}
