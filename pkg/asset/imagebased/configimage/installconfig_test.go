package configimage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/none"
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
networking:
  networkType: OVNKubernetes
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 1
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
platform:
  aws:
    region: us-east-1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"c3VwZXItc2VjcmV0Cg==\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: Platform: Unsupported value: "aws": supported values: "none"`,
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
			expectedError: "invalid install-config configuration: Compute.Replicas: Required value: Total number of Compute.Replicas must be 0 when ControlPlane.Replicas is 1 for platform none. Found 3",
		},
		{
			name: "incorrect controlplane.replicas set",
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
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
platform:
  none : {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: ControlPlane.Replicas: Required value: Only Single Node OpenShift (SNO) is supported, total number of ControlPlane.Replicas must be 1. Found 2",
		},
		{
			name: "valid dual-stack configuration with two IPv4 machine networks",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 10.0.0.0/16
  - cidr: 10.10.0.0/24
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 1
compute:
  - architecture: amd64
    hyperthreading: Enabled
    name: worker
    platform: {}
    replicas: 0
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
				BaseDomain: "test-domain",
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
						{CIDR: *ipnet.MustParseCIDR("10.10.0.0/24")},
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       ptr.To[int64](1),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       ptr.To[int64](0),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform:   types.Platform{None: &none.Platform{}},
				PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
			},
		},
		{
			name: "valid dual-stack configuration with IPv4 and IPv6 machine networks",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 10.0.0.0/24
  - cidr: 2001:db8::/64
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
				BaseDomain: "test-domain",
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
						{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       ptr.To[int64](1),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       ptr.To[int64](0),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform:   types.Platform{None: &none.Platform{}},
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
			},
		},
		{
			name: "valid single-stack IPv6 configuration",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 2001:db8::/64
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
				BaseDomain: "test-domain",
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("2001:db8::/64")},
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       ptr.To[int64](1),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       ptr.To[int64](0),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform:   types.Platform{None: &none.Platform{}},
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
			},
		},
		{
			name: "invalid configuration with 3 machine networks",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 10.0.0.0/24
  - cidr: 2001:db8::/64
  - cidr: 192.168.1.0/24
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
			expectedFound: false,
			expectedError: "invalid install-config configuration: Networking.MachineNetwork: Too many: 3: must have at most 2 items",
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
				BaseDomain: "test-domain",
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       ptr.To[int64](1),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       ptr.To[int64](0),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform:   types.Platform{None: &none.Platform{}},
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
			},
		},
		{
			name: "valid configuration for none platform for sno with machine network",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  networkType: OVNKubernetes
  machineNetwork:
  - cidr: 10.0.0.0/24
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
				BaseDomain: "test-domain",
				Networking: &types.Networking{
					NetworkType: "OVNKubernetes",
					MachineNetwork: []types.MachineNetworkEntry{
						{CIDR: *ipnet.MustParseCIDR("10.0.0.0/24")},
					},
				},
				ControlPlane: &types.MachinePool{
					Name:           "master",
					Replicas:       ptr.To[int64](1),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       ptr.To[int64](0),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform:   types.Platform{None: &none.Platform{}},
				PullSecret: `{"auths":{"example.com":{"auth":"c3VwZXItc2VjcmV0Cg=="}}}`,
			},
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

			asset := &InstallConfig{}
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
