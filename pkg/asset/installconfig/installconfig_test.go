package installconfig

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/none"
)

func TestInstallConfigGenerate_FillsInDefaults(t *testing.T) {
	sshPublicKey := &sshPublicKey{}
	baseDomain := &baseDomain{"test-domain", types.ExternalPublishingStrategy}
	clusterName := &clusterName{"test-cluster"}
	pullSecret := &pullSecret{`{"auths":{"example.com":{"auth":"authorization value"}}}`}
	platform := &platform{
		Platform: types.Platform{None: &none.Platform{}},
	}
	installConfig := &InstallConfig{}
	parents := asset.Parents{}
	parents.Add(
		sshPublicKey,
		baseDomain,
		clusterName,
		pullSecret,
		platform,
	)
	if err := installConfig.Generate(context.Background(), parents); err != nil {
		t.Errorf("unexpected error generating install config: %v", err)
	}
	expected := &types.InstallConfig{
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
			Replicas:       pointer.Int64Ptr(3),
			Hyperthreading: types.HyperthreadingEnabled,
			Architecture:   types.ArchitectureAMD64,
		},
		Compute: []types.MachinePool{
			{
				Name:           "worker",
				Replicas:       pointer.Int64Ptr(3),
				Hyperthreading: types.HyperthreadingEnabled,
				Architecture:   types.ArchitectureAMD64,
			},
		},
		Platform: types.Platform{
			None: &none.Platform{},
		},
		PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
		Publish:    types.ExternalPublishingStrategy,
	}
	assert.Equal(t, expected, installConfig.Config, "unexpected config generated")
}

func TestInstallConfigLoad(t *testing.T) {
	cases := []struct {
		name           string
		data           string
		fetchError     error
		expectedFound  bool
		expectedError  bool
		expectedConfig *types.InstallConfig
	}{
		{
			name: "valid InstallConfig",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
baseDomain: test-domain
platform:
  none: {}
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
					Replicas:       pointer.Int64Ptr(3),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       pointer.Int64Ptr(3),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform: types.Platform{
					None: &none.Platform{},
				},
				PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
		{
			name: "invalid InstallConfig",
			data: `
metadata:
  name: test-cluster
`,
			expectedError: true,
		},
		{
			name:          "empty",
			data:          "",
			expectedError: true,
		},
		{
			name:          "not yaml",
			data:          "This is not yaml.",
			expectedError: true,
		},
		{
			name:       "file not found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error fetching file",
			fetchError:    fmt.Errorf("fetch failed"),
			expectedError: true,
		},
		{
			name: "unknown field",
			data: `
apiVersion: v1
metadata:
  name: test-cluster
additionalTrustBundlePolicy: Proxyonly
baseDomain: test-domain
platform:
  none: {}
pullSecret: "{\"auths\":{\"example.com\":{\"auth\":\"authorization value\"}}}"
wrong_key: wrong_value 
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
					Replicas:       pointer.Int64Ptr(3),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       pointer.Int64Ptr(3),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform: types.Platform{
					None: &none.Platform{},
				},
				PullSecret: `{"auths":{"example.com":{"auth":"authorization value"}}}`,
				Publish:    types.ExternalPublishingStrategy,
			},
		},
		{
			name: "old valid InstallConfig",
			data: `
apiVersion: v1beta3
metadata:
  name: test-cluster
baseDomain: test-domain
platform:
  none: {}
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
					Replicas:       pointer.Int64Ptr(3),
					Hyperthreading: types.HyperthreadingEnabled,
					Architecture:   types.ArchitectureAMD64,
				},
				Compute: []types.MachinePool{
					{
						Name:           "worker",
						Replicas:       pointer.Int64Ptr(3),
						Hyperthreading: types.HyperthreadingEnabled,
						Architecture:   types.ArchitectureAMD64,
					},
				},
				Platform: types.Platform{
					None: &none.Platform{},
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
				)

			ic := &InstallConfig{}
			found, err := ic.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError {
				assert.Error(t, err, "expected error from Load")
			} else {
				assert.NoError(t, err, "unexpected error from Load")
			}
			if tc.expectedFound {
				assert.Equal(t, tc.expectedConfig, ic.Config, "unexpected Config in InstallConfig")
			}
		})
	}
}
