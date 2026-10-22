package machines

import (
	"encoding/json"
	"testing"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

func TestNodeDiskSetupAzureDevice(t *testing.T) {
	ic := installconfig.MakeAsset(&types.InstallConfig{
		FeatureSet: configv1.TechPreviewNoUpgrade,
		Platform:   types.Platform{Azure: &azuretypes.Platform{}},
	})

	cases := []struct {
		name      string
		role      string
		diskSetup types.Disk
		lun       int32
		device    string
	}{
		{
			name:      "etcd disk on the first LUN",
			role:      "master",
			diskSetup: types.Disk{Type: types.Etcd, Etcd: &types.DiskEtcd{PlatformDiskID: "etcd"}},
			lun:       0,
			device:    "/dev/disk/azure/data/by-lun/0",
		},
		{
			name:      "etcd disk after a LUN gap",
			role:      "master",
			diskSetup: types.Disk{Type: types.Etcd, Etcd: &types.DiskEtcd{PlatformDiskID: "etcd"}},
			lun:       5,
			device:    "/dev/disk/azure/data/by-lun/5",
		},
		{
			name: "user-defined disk on the last LUN",
			role: "worker",
			diskSetup: types.Disk{Type: types.UserDefined, UserDefined: &types.DiskUserDefined{
				PlatformDiskID: "data", MountPath: "/var/lib/containers",
			}},
			lun:    63,
			device: "/dev/disk/azure/data/by-lun/63",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mc, err := NodeDiskSetup(ic, tc.role, tc.diskSetup, v1beta1.DataDisk{Lun: ptr.To(tc.lun)})
			require.NoError(t, err)
			require.NotNil(t, mc)

			var ign igntypes.Config
			require.NoError(t, json.Unmarshal(mc.Spec.Config.Raw, &ign))
			require.Len(t, ign.Storage.Disks, 1)
			assert.Equal(t, tc.device, ign.Storage.Disks[0].Device)
		})
	}
}
