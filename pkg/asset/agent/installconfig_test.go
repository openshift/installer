package agent

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

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
platform:
  aws:
    region: us-east-1
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: `invalid install-config configuration: Platform: Unsupported value: "aws": supported values: "baremetal", "vsphere", "none"`,
		},
		{
			name: "apiVips not set for baremetal platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
platform:
  baremetal:
    hosts:
      - name: host1
        bootMACAddress: 52:54:01:xx:zz:z1
    ingressVips:
      - 192.168.122.11
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: Platform.Baremetal.ApiVips: Required value: apiVips must be set for baremetal platform",
		},
		{
			name: "ingressVips not set for vsphere platform",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
platform:
  vsphere:
    apiVips:
      - 192.168.122.10
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
`,
			expectedFound: false,
			expectedError: "invalid install-config configuration: Platform.VSphere.IngressVips: Required value: ingressVips must be set for vsphere platform",
		},
		{
			name: "invalid configuration for none platform for sno",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
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
			expectedFound: false,
			expectedError: "invalid install-config configuration: [ControlPlane.Replicas: Required value: control plane replicas must be 1 for none platform. Found 3, Compute.Replicas: Required value: total number of worker replicas must be 0 for none platform. Found 2]",
		},
		{
			name: "no compute.replicas set for SNO",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
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
			expectedError: "invalid install-config configuration: Compute.Replicas: Required value: Installing a Single Node Openshift requires explicitly setting compute replicas to zero",
		},
		{
			name: "valid configuration for none platform for sno",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
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
				BaseDomain: "test-domain",
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
					Replicas:       pointer.Int64Ptr(1),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       pointer.Int64Ptr(0),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform:   types.Platform{None: &none.Platform{}},
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
