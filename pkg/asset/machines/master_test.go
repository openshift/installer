package machines

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

func workloadForCpuset(cpuset string) []types.Workload {
	return []types.Workload{
		{
			Name:   types.ManagementWorkload,
			CPUIDs: cpuset,
		},
	}
}

func expectedMachineConfigForWorkloads(role string, workloads []types.Workload) string {
	expectedCrioCfg, _ := machineconfig.CrioWorkloadDropinContents(workloads)
	expectedKubeletCfg, _ := machineconfig.KubeletWorkloadDropinContents(workloads)
	return fmt.Sprintf(`apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  creationTimestamp: null
  labels:
    machineconfiguration.openshift.io/role: %[1]s
  name: 02-%[1]s-workload-partitioning
spec:
  config:
    ignition:
      version: 3.2.0
    storage:
      files:
      - contents:
          source: %s
        mode: 420
        overwrite: true
        path: /etc/crio/crio.conf.d/01-workload-partitioning
        user:
          name: root
      - contents:
          source: %s
        mode: 420
        overwrite: true
        path: /etc/kubernetes/workload-pinning
        user:
          name: root
  extensions: null
  fips: false
  kernelArguments: null
  kernelType: ""
  osImageURL: ""
`, role, dataurl.EncodeBytes([]byte(expectedCrioCfg)), dataurl.EncodeBytes([]byte(expectedKubeletCfg)))
}

func TestMasterGenerateMachineConfigs(t *testing.T) {
	cases := []struct {
		name                  string
		key                   string
		hyperthreading        types.HyperthreadingMode
		workloads             []types.Workload
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
  config:
    ignition:
      version: 3.2.0
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,QUREIG5vc210
        mode: 384
        overwrite: true
        path: /etc/pivot/kernel-args
        user:
          name: root
  extensions: null
  fips: false
  kernelArguments: null
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
  config:
    ignition:
      version: 3.2.0
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,QUREIG5vc210
        mode: 384
        overwrite: true
        path: /etc/pivot/kernel-args
        user:
          name: root
  extensions: null
  fips: false
  kernelArguments: null
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
			name:                  "workload enabled",
			workloads:             workloadForCpuset("3,4,5"),
			expectedMachineConfig: []string{expectedMachineConfigForWorkloads("master", workloadForCpuset("3,4,5"))},
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
				&installconfig.InstallConfig{
					Config: &types.InstallConfig{
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
							Workloads: tc.workloads,
						},
					},
				},
				(*rhcos.Image)(pointer.StringPtr("test-image")),
				&machine.Master{
					File: &asset.File{
						Filename: "master-ignition",
						Data:     []byte("test-ignition"),
					},
				},
			)
			master := &Master{}
			if err := master.Generate(parents); err != nil {
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
	installConfig := installconfig.InstallConfig{
		Config: &types.InstallConfig{
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
		},
	}

	parents.Add(
		&installconfig.ClusterID{
			UUID:    "test-uuid",
			InfraID: "test-infra-id",
		},
		&installConfig,
		(*rhcos.Image)(pointer.StringPtr("test-image")),
		&machine.Master{
			File: &asset.File{
				Filename: "master-ignition",
				Data:     []byte("test-ignition"),
			},
		},
	)
	master := &Master{}
	if err := master.Generate(parents); err != nil {
		t.Fatalf("failed to generate master machines: %v", err)
	}

	if installConfig.Config.ControlPlane.Platform.AWS.Type != "" {
		t.Fatalf("control plance in the install config has been modified")
	}
}
