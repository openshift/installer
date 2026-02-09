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
	icaws "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

func TestWorkerGenerate(t *testing.T) {
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
    machineconfiguration.openshift.io/role: worker
  name: 99-worker-ssh
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
    machineconfiguration.openshift.io/role: worker
  name: 99-worker-disable-hyperthreading
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
    machineconfiguration.openshift.io/role: worker
  name: 99-worker-disable-hyperthreading
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
    machineconfiguration.openshift.io/role: worker
  name: 99-worker-ssh
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
			cfg := &types.InstallConfig{
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
				Compute: []types.MachinePool{
					{
						Replicas:       ptr.To(int64(1)),
						Hyperthreading: tc.hyperthreading,
						Platform: types.MachinePoolPlatform{
							AWS: &awstypes.MachinePool{
								Zones:        []string{"us-east-1a"},
								InstanceType: "m5.large",
							},
						},
					},
				},
			}
			icAsset := installconfig.MakeAsset(cfg)
			icAsset.AWS = icaws.NewMetadata(cfg.Platform.AWS.Region, cfg.Platform.AWS.VPC.Subnets, nil)
			parents.Add(
				&installconfig.ClusterID{
					UUID:    "test-uuid",
					InfraID: "test-infra-id",
				},
				icAsset,
				rhcos.MakeAsset("test-image"),
				(*rhcos.Release)(ptr.To("412.86.202208101040-0")),
				&machine.Worker{
					File: &asset.File{
						Filename: "worker-ignition",
						Data:     []byte("test-ignition"),
					},
				},
			)
			worker := &Worker{}
			if err := worker.Generate(context.Background(), parents); err != nil {
				t.Fatalf("failed to generate worker machines: %v", err)
			}
			expectedLen := len(tc.expectedMachineConfig)
			if assert.Equal(t, expectedLen, len(worker.MachineConfigFiles)) {
				for i := 0; i < expectedLen; i++ {
					assert.Equal(t, tc.expectedMachineConfig[i], string(worker.MachineConfigFiles[i].Data), "unexepcted machine config contents")
				}
			} else {
				assert.Equal(t, 0, len(worker.MachineConfigFiles), "expected no machine config files")
			}
		})
	}
}

func TestComputeIsNotModified(t *testing.T) {
	parents := asset.Parents{}
	cfg := &types.InstallConfig{
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
		Compute: []types.MachinePool{
			{
				Replicas:       ptr.To(int64(1)),
				Hyperthreading: types.HyperthreadingDisabled,
				Platform: types.MachinePoolPlatform{
					AWS: &awstypes.MachinePool{
						Zones:        []string{"us-east-1a"},
						InstanceType: "",
					},
				},
			},
		},
	}
	installConfig := installconfig.MakeAsset(cfg)
	installConfig.AWS = icaws.NewMetadata(cfg.Platform.AWS.Region, cfg.Platform.AWS.VPC.Subnets, nil)

	parents.Add(
		&installconfig.ClusterID{
			UUID:    "test-uuid",
			InfraID: "test-infra-id",
		},
		installConfig,
		rhcos.MakeAsset("test-image"),
		(*rhcos.Release)(ptr.To("412.86.202208101040-0")),
		&machine.Worker{
			File: &asset.File{
				Filename: "worker-ignition",
				Data:     []byte("test-ignition"),
			},
		},
	)
	worker := &Worker{}
	if err := worker.Generate(context.Background(), parents); err != nil {
		t.Fatalf("failed to generate master machines: %v", err)
	}

	if installConfig.Config.Compute[0].Platform.AWS.Type != "" {
		t.Fatalf("compute in the install config has been modified")
	}
}

func TestDefaultAWSMachinePoolPlatform(t *testing.T) {
	type testCase struct {
		name                string
		poolName            string
		expectedMachinePool awstypes.MachinePool
		assert              func(tc *testCase)
	}
	cases := []testCase{
		{
			name:     "default EBS type for compute pool",
			poolName: types.MachinePoolComputeRoleName,
			expectedMachinePool: awstypes.MachinePool{
				EC2RootVolume: awstypes.EC2RootVolume{
					Type: awstypes.VolumeTypeGp3,
					Size: decimalRootVolumeSize,
				},
			},
			assert: func(tc *testCase) {
				mp := defaultAWSMachinePoolPlatform(tc.poolName)
				want := tc.expectedMachinePool.EC2RootVolume.Type
				got := mp.EC2RootVolume.Type
				assert.Equal(t, want, got, "unexepcted EBS type")
			},
		},
		{
			name:     "default EBS type for edge pool",
			poolName: types.MachinePoolEdgeRoleName,
			expectedMachinePool: awstypes.MachinePool{
				EC2RootVolume: awstypes.EC2RootVolume{
					Type: awstypes.VolumeTypeGp2,
					Size: decimalRootVolumeSize,
				},
			},
			assert: func(tc *testCase) {
				mp := defaultAWSMachinePoolPlatform(tc.poolName)
				want := tc.expectedMachinePool.EC2RootVolume.Type
				got := mp.EC2RootVolume.Type
				assert.Equal(t, want, got, "unexepcted EBS type")
			},
		},
	}
	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.assert(&tc)
		})
	}
}
