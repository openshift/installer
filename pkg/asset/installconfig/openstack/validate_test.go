package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/types"
)

func TestValidateForProvisioning(t *testing.T) {
	cases := []struct {
		name           string
		installConfig  *types.InstallConfig
		expectedErrMsg string
	}{
		{
			name: "three-node control plane",
			installConfig: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Replicas: pointer.Int64Ptr(3),
				},
			},
			expectedErrMsg: "",
		}, {
			name: "five-node control plane",
			installConfig: &types.InstallConfig{
				ControlPlane: &types.MachinePool{
					Replicas: pointer.Int64Ptr(5),
				},
			},
			expectedErrMsg: `controlPlane.replicas: Invalid value: 5: control plane must be exactly three nodes when provisioning on OpenStack`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateForProvisioning(tc.installConfig)
			if tc.expectedErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedErrMsg, err)
			}
		})
	}
}
