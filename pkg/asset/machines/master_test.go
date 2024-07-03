package machines

import (
	"context"
	"fmt"
	"testing"

	"github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/baremetal"
)

func TestMasterGenerateMachineConfigs(t *testing.T) {
	cases := []struct {
		name                  string
		key                   string
		hyperthreading        types.HyperthreadingMode
		expectedMachineConfig []string
	}{
		{
			name:           "no key hyperthreading enabled",
			hyperthreading: types.HyperthreadingEnabled,
		},
		{
			name:           "key present hyperthreading enabled",
			key:            "ssh-rsa: dummy-key",
			hyperthreading: types.HyperthreadingEnabled,
			expectedMachineConfig: []string{`apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: master
  name: 99-master-ssh
spec:
  baseOSExtensionsContainerImage: ""
  config:
    ignition:
      version: 3.2.0
    passwd:
      users:
      - name: core
        sshAuthorizedKeys:
        - 'ssh-rsa: dummy-key'
  extensions: null
  fips: false
  kernelArguments: null
  kernelType: ""
  osImageURL: ""
`},
		},
		{
			name:           "no key hyperthreading disabled",
			hyperthreading: types.HyperthreadingDisabled,
			expectedMachineConfig: []string{`apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: master
  name: 99-master-disable-hyperthreading
spec:
  baseOSExtensionsContainerImage: ""
  config:
    ignition:
      version: 3.2.0
  extensions: null
  fips: false
  kernelArguments:
  - nosmt
  - smt-enabled=off
  kernelType: ""
  osImageURL: ""
`},
		},
		{
			name:           "key present hyperthreading disabled",
			key:            "ssh-rsa: dummy-key",
			hyperthreading: types.HyperthreadingDisabled,
			expectedMachineConfig: []string{`apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: master
  name: 99-master-disable-hyperthreading
spec:
  baseOSExtensionsContainerImage: ""
  config:
    ignition:
      version: 3.2.0
  extensions: null
  fips: false
  kernelArguments:
  - nosmt
  - smt-enabled=off
  kernelType: ""
  osImageURL: ""
`, `apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: master
  name: 99-master-ssh
spec:
  baseOSExtensionsContainerImage: ""
  config:
    ignition:
      version: 3.2.0
    passwd:
      users:
      - name: core
        sshAuthorizedKeys:
        - 'ssh-rsa: dummy-key'
  extensions: null
  fips: false
  kernelArguments: null
  kernelType: ""
  osImageURL: ""
`},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(
				&installconfig.ClusterID{
					UUID:    "test-uuid",
					InfraID: "test-infra-id",
				},
				installconfig.MakeAsset(
					&types.InstallConfig{
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-cluster",
						},
						SSHKey:     tc.key,
						BaseDomain: "test-domain",
						Platform: types.Platform{
							AWS: &awstypes.Platform{
								Region: "us-east-1",
							},
						},
						ControlPlane: &types.MachinePool{
							Hyperthreading: tc.hyperthreading,
							Replicas:       pointer.Int64Ptr(1),
							Platform: types.MachinePoolPlatform{
								AWS: &awstypes.MachinePool{
									Zones:        []string{"us-east-1a"},
									InstanceType: "m5.xlarge",
								},
							},
						},
					}),
				rhcos.MakeAsset("test-image"),
				(*rhcos.Release)(pointer.StringPtr("412.86.202208101040-0")),
				&machine.Master{
					File: &asset.File{
						Filename: "master-ignition",
						Data:     []byte("test-ignition"),
					},
				},
			)
			master := &Master{}
			if err := master.Generate(context.Background(), parents); err != nil {
				t.Fatalf("failed to generate master machines: %v", err)
			}
			expectedLen := len(tc.expectedMachineConfig)
			if assert.Equal(t, expectedLen, len(master.MachineConfigFiles)) {
				for i := 0; i < expectedLen; i++ {
					assert.Equal(t, tc.expectedMachineConfig[i], string(master.MachineConfigFiles[i].Data), "unexepcted machine config contents")
				}
			} else {
				assert.Equal(t, 0, len(master.MachineConfigFiles), "expected no machine config files")
			}
		})
	}
}

func TestControlPlaneIsNotModified(t *testing.T) {
	parents := asset.Parents{}
	installConfig := installconfig.MakeAsset(
		&types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			SSHKey:     "ssh-rsa: dummy-key",
			BaseDomain: "test-domain",
			Platform: types.Platform{
				AWS: &awstypes.Platform{
					Region: "us-east-1",
					DefaultMachinePlatform: &awstypes.MachinePool{
						InstanceType: "TEST_INSTANCE_TYPE",
					},
				},
			},
			ControlPlane: &types.MachinePool{
				Hyperthreading: types.HyperthreadingDisabled,
				Replicas:       pointer.Int64Ptr(1),
				Platform: types.MachinePoolPlatform{
					AWS: &awstypes.MachinePool{
						Zones: []string{"us-east-1a"},
					},
				},
			},
		})

	parents.Add(
		&installconfig.ClusterID{
			UUID:    "test-uuid",
			InfraID: "test-infra-id",
		},
		installConfig,
		rhcos.MakeAsset("test-image"),
		(*rhcos.Release)(pointer.StringPtr("412.86.202208101040-0")),
		&machine.Master{
			File: &asset.File{
				Filename: "master-ignition",
				Data:     []byte("test-ignition"),
			},
		},
	)
	master := &Master{}
	if err := master.Generate(context.Background(), parents); err != nil {
		t.Fatalf("failed to generate master machines: %v", err)
	}

	if installConfig.Config.ControlPlane.Platform.AWS.Type != "" {
		t.Fatalf("control plance in the install config has been modified")
	}
}

func TestBaremetalGeneratedAssetFiles(t *testing.T) {
	parents := asset.Parents{}
	installConfig := installconfig.MakeAsset(
		&types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			Platform: types.Platform{
				BareMetal: &baremetal.Platform{
					Hosts: []*baremetal.Host{
						{
							Name: "master-0",
							Role: "master",
							BMC: baremetal.BMC{
								Username: "usr-0",
								Password: "pwd-0",
							},
							NetworkConfig: networkConfig("interfaces:"),
						},
						{
							Name: "worker-0",
							Role: "worker",
							BMC: baremetal.BMC{
								Username: "usr-1",
								Password: "pwd-1",
							},
						},
					},
				},
			},
			ControlPlane: &types.MachinePool{
				Replicas: pointer.Int64Ptr(1),
				Platform: types.MachinePoolPlatform{
					BareMetal: &baremetal.MachinePool{},
				},
			},
			Compute: []types.MachinePool{
				{
					Replicas: pointer.Int64Ptr(1),
					Platform: types.MachinePoolPlatform{},
				},
			},
		})

	parents.Add(
		&installconfig.ClusterID{
			UUID:    "test-uuid",
			InfraID: "test-infra-id",
		},
		installConfig,
		rhcos.MakeAsset("test-image"),
		(*rhcos.Release)(pointer.StringPtr("412.86.202208101040-0")),
		&machine.Master{
			File: &asset.File{
				Filename: "master-ignition",
				Data:     []byte("test-ignition"),
			},
		},
	)
	master := &Master{}
	assert.NoError(t, master.Generate(context.Background(), parents))

	assert.Len(t, master.HostFiles, 2)
	verifyHost(t, master.HostFiles[0], "openshift/99_openshift-cluster-api_hosts-0.yaml", "master-0")
	verifyHost(t, master.HostFiles[1], "openshift/99_openshift-cluster-api_hosts-1.yaml", "worker-0")

	assert.Len(t, master.SecretFiles, 2)
	verifySecret(t, master.SecretFiles[0], "openshift/99_openshift-cluster-api_host-bmc-secrets-0.yaml", "master-0-bmc-secret", "map[password:[112 119 100 45 48] username:[117 115 114 45 48]]")
	verifySecret(t, master.SecretFiles[1], "openshift/99_openshift-cluster-api_host-bmc-secrets-1.yaml", "worker-0-bmc-secret", "map[password:[112 119 100 45 49] username:[117 115 114 45 49]]")

	assert.Len(t, master.NetworkConfigSecretFiles, 1)
	verifySecret(t, master.NetworkConfigSecretFiles[0], "openshift/99_openshift-cluster-api_host-network-config-secrets-0.yaml", "master-0-network-config-secret", "map[nmstate:[105 110 116 101 114 102 97 99 101 115 58 32 110 117 108 108 10]]")
	verifyManifestOwnership(t, master)
}

func verifyHost(t *testing.T, a *asset.File, eFilename, eName string) {
	t.Helper()
	assert.Equal(t, a.Filename, eFilename)
	var host v1alpha1.BareMetalHost
	assert.NoError(t, yaml.Unmarshal(a.Data, &host))
	assert.Equal(t, eName, host.Name)
}

func verifySecret(t *testing.T, a *asset.File, eFilename, eName, eData string) {
	t.Helper()
	assert.Equal(t, a.Filename, eFilename)
	var secret corev1.Secret
	assert.NoError(t, yaml.Unmarshal(a.Data, &secret))
	assert.Equal(t, eName, secret.Name)
	assert.Equal(t, eData, fmt.Sprintf("%v", secret.Data))
}

func verifyManifestOwnership(t *testing.T, a asset.WritableAsset) {
	t.Helper()
	for _, file := range a.Files() {
		assert.True(t, IsMachineManifest(file), file.Filename)
	}
}

func networkConfig(config string) *v1.JSON {
	var nc v1.JSON
	yaml.Unmarshal([]byte(config), &nc)
	return &nc
}
