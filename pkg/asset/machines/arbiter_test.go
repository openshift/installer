package machines

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/none"
)

func TestArbiterGenerateMachineConfigs(t *testing.T) {
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
  labels:
    machineconfiguration.openshift.io/role: arbiter
  name: 99-arbiter-ssh
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
  labels:
    machineconfiguration.openshift.io/role: arbiter
  name: 99-arbiter-disable-hyperthreading
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
  labels:
    machineconfiguration.openshift.io/role: arbiter
  name: 99-arbiter-disable-hyperthreading
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
  labels:
    machineconfiguration.openshift.io/role: arbiter
  name: 99-arbiter-ssh
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
							BareMetal: &baremetal.Platform{},
						},
						Arbiter: &types.MachinePool{
							Hyperthreading: tc.hyperthreading,
							Replicas:       ptr.To(int64(1)),
							Platform: types.MachinePoolPlatform{
								BareMetal: &baremetal.MachinePool{},
							},
						},
					}),
				rhcos.MakeAsset("test-image"),
				(*rhcos.Release)(ptr.To("412.86.202208101040-0")),
				&machine.Arbiter{
					File: &asset.File{
						Filename: "arbiter-ignition",
						Data:     []byte("test-ignition"),
					},
				},
			)
			arbiter := &Arbiter{}
			if err := arbiter.Generate(context.Background(), parents); err != nil {
				t.Fatalf("failed to generate arbiter machines: %v", err)
			}
			expectedLen := len(tc.expectedMachineConfig)
			if assert.Equal(t, expectedLen, len(arbiter.MachineConfigFiles)) {
				for i := 0; i < expectedLen; i++ {
					assert.Equal(t, tc.expectedMachineConfig[i], string(arbiter.MachineConfigFiles[i].Data), "unexepcted machine config contents")
				}
			} else {
				assert.Equal(t, 0, len(arbiter.MachineConfigFiles), "expected no machine config files")
			}
		})
	}
}

func TestArbiterInstallOnlyForBaremetal(t *testing.T) {
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
			Arbiter: &types.MachinePool{
				Hyperthreading: types.HyperthreadingDisabled,
				Replicas:       ptr.To(int64(1)),
				Platform: types.MachinePoolPlatform{
					BareMetal: &baremetal.MachinePool{},
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
		(*rhcos.Release)(ptr.To("412.86.202208101040-0")),
		&machine.Arbiter{
			File: &asset.File{
				Filename: "arbiter-ignition",
				Data:     []byte("test-ignition"),
			},
		},
	)
	arbiter := &Arbiter{}
	err := arbiter.Generate(context.Background(), parents)
	assert.NotNil(t, err, "expected arbiter generate to fail for non baremetal platforms")
	assert.Contains(t, err.Error(), "only BareMetal and None platforms are supported for Arbiter deployments")
}

func TestArbiterInstallNonePlatform(t *testing.T) {
	parents := asset.Parents{}
	installConfig := installconfig.MakeAsset(
		&types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			SSHKey:     "ssh-rsa: dummy-key",
			BaseDomain: "test-domain",
			Platform: types.Platform{
				None: &none.Platform{},
			},
			Arbiter: &types.MachinePool{
				Hyperthreading: types.HyperthreadingDisabled,
				Replicas:       ptr.To(int64(1)),
			},
		})

	parents.Add(
		&installconfig.ClusterID{
			UUID:    "test-uuid",
			InfraID: "test-infra-id",
		},
		installConfig,
		rhcos.MakeAsset("test-image"),
		(*rhcos.Release)(ptr.To("412.86.202208101040-0")),
		&machine.Arbiter{
			File: &asset.File{
				Filename: "arbiter-ignition",
				Data:     []byte("test-ignition"),
			},
		},
	)
	arbiter := &Arbiter{}
	err := arbiter.Generate(context.Background(), parents)
	assert.Nil(t, err, "expected arbiter generate to succeed for None platform")
	// For None platform, no Machine API objects or BareMetalHost CRs should be generated
	assert.Equal(t, 0, len(arbiter.MachineFiles), "expected no machine files for None platform")
	assert.Equal(t, 0, len(arbiter.HostFiles), "expected no host files for None platform")
	assert.Equal(t, 0, len(arbiter.SecretFiles), "expected no secret files for None platform")
	// MachineConfigs should still be generated (hyperthreading disabled + SSH key)
	assert.Greater(t, len(arbiter.MachineConfigFiles), 0, "expected machine config files for None platform")
}

func TestArbiterIsNotModified(t *testing.T) {
	parents := asset.Parents{}
	installConfig := installconfig.MakeAsset(
		&types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			SSHKey:     "ssh-rsa: dummy-key",
			BaseDomain: "test-domain",
			Platform: types.Platform{
				BareMetal: &baremetal.Platform{
					ClusterProvisioningIP:  "127.0.0.1",
					DefaultMachinePlatform: &baremetal.MachinePool{},
				},
			},
			Arbiter: &types.MachinePool{
				Hyperthreading: types.HyperthreadingDisabled,
				Replicas:       ptr.To(int64(1)),
				Platform: types.MachinePoolPlatform{
					BareMetal: &baremetal.MachinePool{},
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
		(*rhcos.Release)(ptr.To("412.86.202208101040-0")),
		&machine.Arbiter{
			File: &asset.File{
				Filename: "arbiter-ignition",
				Data:     []byte("test-ignition"),
			},
		},
	)
	arbiter := &Arbiter{}
	if err := arbiter.Generate(context.Background(), parents); err != nil {
		t.Fatalf("failed to generate arbiter machines: %v", err)
	}

	if installConfig.Config.Arbiter.Platform.BareMetal == nil {
		t.Fatalf("arbiter in the install config has been modified")
	}
}

func TestArbiterBaremetalGeneratedAssetFiles(t *testing.T) {
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
							Name: "arbiter-0",
							Role: "arbiter",
							BMC: baremetal.BMC{
								Username: "usr-0",
								Password: "pwd-0",
							},
							NetworkConfig: networkConfig("interfaces:"),
						},
					},
				},
			},
			Arbiter: &types.MachinePool{
				Replicas: ptr.To(int64(1)),
				Platform: types.MachinePoolPlatform{
					BareMetal: &baremetal.MachinePool{},
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
		(*rhcos.Release)(ptr.To("412.86.202208101040-0")),
		&machine.Arbiter{
			File: &asset.File{
				Filename: "arbiter-ignition",
				Data:     []byte("test-ignition"),
			},
		},
	)
	arbiter := &Arbiter{}
	assert.NoError(t, arbiter.Generate(context.Background(), parents))

	assert.Len(t, arbiter.HostFiles, 1)
	verifyHost(t, arbiter.HostFiles[0], "openshift/99_openshift-cluster-api_arbiter_hosts-0.yaml", "arbiter-0")

	assert.Len(t, arbiter.SecretFiles, 1)
	verifySecret(t, arbiter.SecretFiles[0], "openshift/99_openshift-cluster-api_arbiter_host-bmc-secrets-0.yaml", "arbiter-0-bmc-secret", "map[password:[112 119 100 45 48] username:[117 115 114 45 48]]")

	assert.Len(t, arbiter.NetworkConfigSecretFiles, 1)
	verifySecret(t, arbiter.NetworkConfigSecretFiles[0], "openshift/99_openshift-cluster-api_arbiter_host-network-config-secrets-0.yaml", "arbiter-0-network-config-secret", "map[nmstate:[105 110 116 101 114 102 97 99 101 115 58 32 110 117 108 108 10]]")
}
