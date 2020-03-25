package machines

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

func TestMasterGenerateMachineConfigs(t *testing.T) {
	cases := []struct {
		name                  string
		key                   string
		hyperthreading        types.HyperthreadingMode
		expectedMachineConfig []string
	}{
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
      version: 3.0.0
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,QUREIG5vc210
        mode: 384
        overwrite: true
        path: /etc/pivot/kernel-args
        user:
          name: root
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
      version: 3.0.0
    storage:
      files:
      - contents:
          source: data:text/plain;charset=utf-8;base64,QUREIG5vc210
        mode: 384
        overwrite: true
        path: /etc/pivot/kernel-args
        user:
          name: root
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
      version: 3.0.0
    passwd:
      users:
      - name: core
        sshAuthorizedKeys:
        - 'ssh-rsa: dummy-key'
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
									Zones: []string{"us-east-1a"},
								},
							},
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
