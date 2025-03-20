package agent

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
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
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
  networkType: OVNKubernetes
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
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: [platform.baremetal.apiVIPs: Required value: must specify at least one VIP for the API, platform.baremetal.apiVIPs: Required value: must specify VIP for API, when VIP for ingress is set]",
		},
		{
			name: "ingressVIP missing and deprecated vSphere credentials are present",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
    vCenter: test.vcenter.com
    username: testuser
    password: testpassword
    datacenter: testDatacenter
    defaultDatastore: testDatastore
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.ingressVIPs: Required value: must specify VIP for ingress, when VIP for API is set, platform.vsphere.ingressVIPs: Required value: must specify at least one VIP for the Ingress, Platform.VSphere.failureDomains.topology.folder: Required value: must specify a folder for agent-based installs]`,
		},
		{
			name: "apiVIPs are missing for vsphere platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    ingressVips:
      - 192.168.122.11
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.apiVIPs: Required value: must specify at least one VIP for the API, platform.vsphere.apiVIPs: Required value: must specify VIP for API, when VIP for ingress is set]`,
		},
		{
			name: "invalid IP values of api and ingress VIPs for vsphere platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.01
    ingressVips:
      - 192.168.122.256
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.apiVIPs: Invalid value: "192.168.122.01": "192.168.122.01" is not a valid IP, platform.vsphere.apiVIPs: Invalid value: "192.168.122.01": IP expected to be in one of the machine networks: 192.168.122.0/23, platform.vsphere.ingressVIPs: Invalid value: "192.168.122.256": "192.168.122.256" is not a valid IP, platform.vsphere.ingressVIPs: Invalid value: "192.168.122.256": IP expected to be in one of the machine networks: 192.168.122.0/23]`,
		},
		{
			name: "api and ingressVIP's are missing for vsphere platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere: {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.apiVIPs: Required value: must specify at least one VIP for the API, platform.vsphere.ingressVIPs: Required value: must specify at least one VIP for the Ingress]`,
		},
		{
			name: "ingress and apiVip's must be different when loadbalancer type is not usermanaged for vsphere platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
    ingressVips:
      - 192.168.122.10
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: platform.vsphere.apiVIPs: Invalid value: "192.168.122.10": VIP for API must not be one of the Ingress VIPs`,
		},
		{
			name: "ingress and apiVip's must be from machine network CIDR when loadbalancer type is not usermanaged for vsphere platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 10.0.0.1
    ingressVips:
      - 10.0.0.2
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.apiVIPs: Invalid value: "10.0.0.1": IP expected to be in one of the machine networks: 192.168.122.0/23, platform.vsphere.ingressVIPs: Invalid value: "10.0.0.2": IP expected to be in one of the machine networks: 192.168.122.0/23]`,
		},
		{
			name: "one must be ipv4 and other must be ipv6 when dual api and ingressVIP's are provided for vsphere platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
      - 192.168.122.11
    ingressVips:
      - 192.168.122.12
      - 192.168.122.13
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.apiVIPs: Invalid value: []string{"192.168.122.10", "192.168.122.11"}: If two API VIPs are given, one must be an IPv4 address, the other an IPv6, platform.vsphere.ingressVIPs: Invalid value: []string{"192.168.122.12", "192.168.122.13"}: If two Ingress VIPs are given, one must be an IPv4 address, the other an IPv6]`,
		},
		{
			name: "api and ingress vips must belong to primary machine network's family for dual stack ipv4/v6",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 2001:db8:1234:1::/120
  - cidr: 192.168.122.0/23
  serviceNetwork:
  - 2001:db8:5678::/108
  - 192.168.112.0/23
  clusterNetwork:
  - cidr: 2001:db8:abcd::/48
    hostPrefix: 64
  - cidr: 10.128.0.0/14
    hostPrefix: 23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
    ingressVips:
      - 192.168.122.11
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.apiVIPs: Invalid value: "192.168.122.10": VIP for the API must be of the same IP family with machine network's primary IP Family for dual-stack IPv4/IPv6, platform.vsphere.ingressVIPs: Invalid value: "192.168.122.11": VIP for the Ingress must be of the same IP family with machine network's primary IP Family for dual-stack IPv4/IPv6]`,
		},
		{
			name: "ingressVIP missing and vcenter vSphere credentials are present",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
    vcenters:
    - server: test.vcenter.com
      user: testuser
      password: testpassword
      datacenters:
      - testDatacenter
    failureDomains:
    - name: testFailuredomain
      server: test.vcenter.com
      zone: testZone
      region: testRegion
      topology:
        computeCluster: "/testDatacenter/host/testcluster"
        datacenter: testDatacenter
        datastore: "/testDatacenter/datastore/testDatastore"
        folder: "/testDatacenter/vm/testFolder"
        networks:
        - testNetwork
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.ingressVIPs: Required value: must specify VIP for ingress, when VIP for API is set, platform.vsphere.ingressVIPs: Required value: must specify at least one VIP for the Ingress]`,
		},
		{
			name: "vcenter vSphere credentials are present but failureDomain server does not match",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
    ingressVips:
      - 192.168.122.11
    vcenters:
    - server: test.vcenter.com
      user: testuser
      password: testpassword
      datacenters:
      - testDatacenter
    failureDomains:
    - name: testFailuredomain
      server: diff1.vcenter.com
      zone: testZone
      region: testRegion
      topology:
        computeCluster: "/testDatacenter/host/testcluster"
        datacenter: testDatacenter
        datastore: "/testDatacenter/datastore/testDatastore"
        folder: "/testDatacenter/vm/testFolder"
        networks:
        - testNetwork
    - name: testFailuredomain2
      server: diff2.vcenter.com
      zone: testZone2
      region: testRegion2
      topology:
        computeCluster: "/testDatacenter2/host/testcluster2"
        datacenter: testDatacenter2
        datastore: "/testDatacenter2/datastore/testDatastore2"
        folder: "/testDatacenter2/vm/testFolder"
        networks:
        - testNetwork2
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.failureDomains.server: Invalid value: "diff1.vcenter.com": server does not exist in vcenters, platform.vsphere.failureDomains.server: Invalid value: "diff2.vcenter.com": server does not exist in vcenters]`,
		},
		{
			name: "All required vSphere fields must be entered if some of them are entered - deprecated fields",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
    ingressVips:
      - 192.168.122.11
    vCenter: test.vcenter.com
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [Platform.VSphere.username: Required value: All credential fields are required if any one is specified, Platform.VSphere.password: Required value: All credential fields are required if any one is specified, Platform.VSphere.datacenter: Required value: All credential fields are required if any one is specified, Platform.VSphere.failureDomains.topology.folder: Required value: must specify a folder for agent-based installs]`,
		},
		{
			name: "All required vSphere fields must be entered if some of them are entered - vcenter fields",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
    ingressVips:
      - 192.168.122.11
    vcenters:
    - server: test.vcenter.com
      user: testuser
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [Platform.VSphere.password: Required value: All credential fields are required if any one is specified, Platform.VSphere.datacenter: Required value: All credential fields are required if any one is specified, Platform.VSphere.failureDomains.topology.folder: Required value: must specify a folder for agent-based installs]`,
		},
		{
			name: "ingressVIP missing for vSphere, credentials not provided and should not flag error",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [platform.vsphere.ingressVIPs: Required value: must specify VIP for ingress, when VIP for API is set, platform.vsphere.ingressVIPs: Required value: must specify at least one VIP for the Ingress]`,
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
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
			expectedError: "invalid install-config configuration: [controlPlane.fencing.credentials: Forbidden: there should be exactly two fencing credentials to support the two node cluster, instead 0 credentials were found, ControlPlane.Replicas: Invalid value: 2: ControlPlane.Replicas can only be set to 5, 4, 3, or 1. Found 2]",
		},
		{
			name: "invalid platform for SNO cluster",
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
  aws:
    region: us-east-1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: [Platform: Unsupported value: \"aws\": supported values: \"baremetal\", \"vsphere\", \"none\", \"external\", Platform: Invalid value: \"aws\": Only platform none and external supports 1 ControlPlane and 0 Compute nodes]",
		},
		{
			name: "invalid platform.baremetal for architecture ppc64le",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
compute:
  - architecture: ppc64le
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
controlPlane:
  architecture: ppc64le
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 3
platform:
  baremetal:
    apiVIP: 192.168.122.10
    ingressVIP: 192.168.122.11
    hosts:
    - name: host1
      bootMACAddress: 52:54:01:aa:aa:a1
    - name: host2
      bootMACAddress: 52:54:01:bb:bb:b1
    - name: host3
      bootMACAddress: 52:54:01:cc:cc:c1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: Platform: Invalid value: \"baremetal\": CPU architecture \"ppc64le\" only supports platform \"none\".",
		},
		{
			name: "invalid platform.baremetal for architecture s390x",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
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
  replicas: 3
platform:
  baremetal:
    apiVIP: 192.168.122.10
    ingressVIP: 192.168.122.11
    hosts:
    - name: host1
      bootMACAddress: 52:54:01:aa:aa:a1
    - name: host2
      bootMACAddress: 52:54:01:bb:bb:b1
    - name: host3
      bootMACAddress: 52:54:01:cc:cc:c1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: Platform: Invalid value: \"baremetal\": CPU architecture \"s390x\" only supports platform \"none\".",
		},
		{
			name: "generic platformName for external platform",
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
  external:
   platformName: some-cloud-provider
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
				Platform: types.Platform{
					External: &external.Platform{
						PlatformName:           "some-cloud-provider",
						CloudControllerManager: "",
					},
				},
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
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
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
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
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
		{
			name: "valid configuration control plane replicas set to 5",
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
  replicas: 5
platform:
  none : {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
					Replicas:       pointer.Int64(5),
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
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
		{
			name: "valid configuration control plane replicas set to 4",
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
  replicas: 4
platform:
  none : {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
					Replicas:       pointer.Int64(4),
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
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
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
  networkType: OVNKubernetes
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
    AdditionalNTPServers:
        - "10.0.1.1"
        - "10.0.1.2"
    hosts:
      - name: host1
        bootMACAddress: 52:54:01:aa:aa:a1
      - name: host2
        bootMACAddress: 52:54:01:bb:bb:b1
      - name: host3
        bootMACAddress: 52:54:01:cc:cc:c1
      - name: host4
        bootMACAddress: 52:54:01:dd:dd:d1
      - name: host5
        bootMACAddress: 52:54:01:ee:ee:e1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
						AdditionalNTPServers:           []string{"10.0.1.1", "10.0.1.2"},
					},
				},
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
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
  networkType: OVNKubernetes
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
    folder: testFolder
    cluster: testCluster
    apiVIP: 192.168.122.10
    ingressVIPs: 
      - 192.168.122.11
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
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
						DeprecatedFolder:           "testFolder",
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
								ResourcePool:   "/testDataCenter/host/testCluster/Resources",
								Folder:         "/testDataCenter/vm/testFolder",
							},
						}},
					},
				},
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
		{
			name: "provisioningNetwork invalid for baremetal cluster",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
  serviceNetwork:
  - 172.30.0.0/16
compute:
  - architecture: amd64
    name: worker
    replicas: 0
controlPlane:
  architecture: amd64
  name: master
  replicas: 3
platform:
  baremetal:
    provisioningNetwork: "UNMANAGED"
    ingressVIPs:
      - 192.168.122.11
    apiVIPs:
      - 192.168.122.10
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: platform.baremetal.provisioningNetwork: Unsupported value: \"UNMANAGED\": supported values: \"Disabled\", \"Managed\", \"Unmanaged\"",
		},
		{
			name: "Provisioning validation failures for baremetal cluster",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 192.168.122.0/23
  serviceNetwork:
  - 172.30.0.0/16
compute:
  - architecture: amd64
    name: worker
    replicas: 0
controlPlane:
  architecture: amd64
  name: master
  replicas: 3
platform:
  baremetal:
    ingressVIPs:
      - 192.168.122.11
    apiVIPs:
      - 192.168.122.10
    clusterProvisioningIP: "172.22.0.11"
    provisioningNetwork: "Managed"
    provisioningMACAddress: "52:54:00:6e:3b:02"
    provisioningNetworkInterface: "eth11"
    provisioningDHCPExternal: true
    provisioningDHCPRange: 172.22.0.10,172.22.0.254
    hosts:
      - name: host1
        bootMACAddress: 52:54:01:aa:aa:a1
        bmc:
          username: "admin"
          password: "password"
          address: "redfish+http://10.10.10.1:8000/redfish/v1/Systems/1234"
      - name: host2
        bootMACAddress: 52:54:01:bb:bb:b1
        bmc:
          username: "admin"
          password: "password"
          address: "redfish+http://10.10.10.2:8000/redfish/v1/Systems/1234"
      - name: host3
        bootMACAddress: 52:54:01:cc:cc:c1
        bmc:
          username: "admin"
          password: "password"
          address: "redfish+http://10.10.10.2:8000/redfish/v1/Systems/1234"
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: [Platform.BareMetal.clusterProvisioningIP: Invalid value: "172.22.0.11": "172.22.0.11" overlaps with the allocated DHCP range, Platform.BareMetal.hosts[2].BMC.Address: Duplicate value: "redfish+http://10.10.10.2:8000/redfish/v1/Systems/1234"]`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(InstallConfigFilename).
				Return(
					&asset.File{
						Filename: InstallConfigFilename,
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
