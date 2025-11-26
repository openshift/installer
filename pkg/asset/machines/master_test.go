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
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
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
				Replicas: ptr.To[int64](3),
				Platform: types.MachinePoolPlatform{
					BareMetal: &baremetal.MachinePool{},
				},
			},
			Compute: []types.MachinePool{
				{
					Replicas: ptr.To[int64](1),
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

type SecretFile struct {
	Filename string
	Name     string
	Data     string
}

func TestFencingCredentialsSecretsFile(t *testing.T) {
	cases := []struct {
		name                          string
		credentials                   []*types.Credential
		fencingCredentialsSecretFiles []*SecretFile
		err                           error
	}{
		{
			name: "valid credentials",
			credentials: []*types.Credential{
				{
					HostName:                "control-plane-0",
					Address:                 "https://redfish/v1/1",
					Username:                "user1",
					Password:                "pass1",
					CertificateVerification: "Enabled",
				},
				{
					HostName:                "control-plane-1",
					Address:                 "http://redfish/v1/2",
					Username:                "user2",
					Password:                "pass2",
					CertificateVerification: "Disabled",
				},
			},
			fencingCredentialsSecretFiles: []*SecretFile{
				{
					Filename: "openshift/99_openshift-etcd_fencing-credentials-secrets-0.yaml",
					Name:     "fencing-credentials-control-plane-0",
					Data:     "map[address:[104 116 116 112 115 58 47 47 114 101 100 102 105 115 104 47 118 49 47 49] certificateVerification:[69 110 97 98 108 101 100] password:[112 97 115 115 49] username:[117 115 101 114 49]]",
				},
				{
					Filename: "openshift/99_openshift-etcd_fencing-credentials-secrets-1.yaml",
					Name:     "fencing-credentials-control-plane-1",
					Data:     "map[address:[104 116 116 112 58 47 47 114 101 100 102 105 115 104 47 118 49 47 50] certificateVerification:[68 105 115 97 98 108 101 100] password:[112 97 115 115 50] username:[117 115 101 114 50]]",
				},
			},
			err: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			installConfig := installconfig.MakeAsset(
				&types.InstallConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-cluster",
					},
					Platform: types.Platform{
						None: &none.Platform{},
					},
					ControlPlane: &types.MachinePool{
						Replicas: ptr.To(int64(2)),
						Fencing: &types.Fencing{
							Credentials: tc.credentials,
						},
					},
					Compute: []types.MachinePool{
						{
							Replicas: ptr.To(int64(0)),
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
				(*rhcos.Release)(ptr.To("412.86.202208101040-0")),
				&machine.Master{
					File: &asset.File{
						Filename: "master-ignition",
						Data:     []byte("test-ignition"),
					},
				},
			)

			master := &Master{}
			if tc.err != nil {
				assert.Error(t, tc.err, master.Generate(context.Background(), parents))
				assert.Len(t, master.FencingCredentialsSecretFiles, len(tc.fencingCredentialsSecretFiles))
				return
			}

			assert.NoError(t, master.Generate(context.Background(), parents))
			assert.Len(t, master.FencingCredentialsSecretFiles, len(tc.fencingCredentialsSecretFiles))
			for i, secret := range master.FencingCredentialsSecretFiles {
				verifySecret(t, secret, tc.fencingCredentialsSecretFiles[i].Filename, tc.fencingCredentialsSecretFiles[i].Name, tc.fencingCredentialsSecretFiles[i].Data)
			}
			verifyManifestOwnership(t, master)
		})
	}
}

func TestSupportedSNOPlatforms(t *testing.T) {
	cases := []struct {
		name                string
		platform            types.Platform
		bootstrapInPlace    *types.BootstrapInPlace
		machinePoolPlatform types.MachinePoolPlatform
		expectedError       bool
	}{
		{
			name: "AWS platform",
			platform: types.Platform{
				AWS: &awstypes.Platform{
					Region: "us-east-1",
					DefaultMachinePlatform: &awstypes.MachinePool{
						InstanceType: "TEST_INSTANCE_TYPE",
					},
				},
			},
			bootstrapInPlace: &types.BootstrapInPlace{},
			machinePoolPlatform: types.MachinePoolPlatform{
				AWS: &awstypes.MachinePool{
					Zones: []string{"us-east-1a"},
				},
			},
			expectedError: false,
		},
		{
			name: "None platform",
			platform: types.Platform{
				None: &none.Platform{},
			},
			machinePoolPlatform: types.MachinePoolPlatform{},
			expectedError:       false,
		},
		{
			name: "IBMCloud platform",
			platform: types.Platform{
				IBMCloud: &ibmcloudtypes.Platform{},
			},
			bootstrapInPlace: &types.BootstrapInPlace{},
			machinePoolPlatform: types.MachinePoolPlatform{
				IBMCloud: &ibmcloudtypes.MachinePool{
					Zones: []string{"us-east-1"},
				},
			},
			expectedError: false,
		},
		{
			name: "External platform with bootstrapInPlace",
			platform: types.Platform{
				External: &external.Platform{},
			},
			bootstrapInPlace: &types.BootstrapInPlace{
				InstallationDisk: "/dev/nvme0n1",
			},
			machinePoolPlatform: types.MachinePoolPlatform{},
			expectedError:       false,
		},
		{
			name: "External platform with empty bootstrapInPlace",
			platform: types.Platform{
				External: &external.Platform{},
			},
			bootstrapInPlace:    &types.BootstrapInPlace{},
			machinePoolPlatform: types.MachinePoolPlatform{},
			expectedError:       true,
		},
		{
			name: "External platform without bootstrapInPlace",
			platform: types.Platform{
				External: &external.Platform{},
			},
			machinePoolPlatform: types.MachinePoolPlatform{},
			expectedError:       true,
		},
		{
			name: "Nutanix platform",
			platform: types.Platform{
				Nutanix: &nutanix.Platform{},
			},
			machinePoolPlatform: types.MachinePoolPlatform{},
			expectedError:       true,
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
						SSHKey:     "ssh-rsa: dummy-key",
						BaseDomain: "test-domain",
						Platform:   tc.platform,
						ControlPlane: &types.MachinePool{
							Hyperthreading: types.HyperthreadingEnabled,
							Replicas:       ptr.To[int64](1),
							Platform:       tc.machinePoolPlatform,
						},
						BootstrapInPlace: tc.bootstrapInPlace,
					}),

				rhcos.MakeAsset("test-image"),
				(*rhcos.Release)(ptr.To[string]("412.86.202208101040-0")),
				&machine.Master{
					File: &asset.File{
						Filename: "master-ignition",
						Data:     []byte("test-ignition"),
					},
				},
			)
			snoErrorRegex := fmt.Sprintf("this install method does not support Single Node installation on platform %s", tc.platform.Name())
			master := &Master{}
			err := master.Generate(context.Background(), parents)
			if tc.expectedError {
				if assert.Error(t, err) {
					assert.Regexp(t, snoErrorRegex, err.Error())
				}
			} else if err != nil {
				// Did not expect a SNO validation error but received a different error
				assert.NotRegexp(t, snoErrorRegex, err.Error())
			}
		})
	}
}
