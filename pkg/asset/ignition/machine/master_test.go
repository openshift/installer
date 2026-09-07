package machine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
)

// TestMasterGenerate tests generating the master asset.
func TestMasterGenerate(t *testing.T) {
	installConfig := installconfig.MakeAsset(
		&types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			BaseDomain: "test-domain",
			Networking: &types.Networking{
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
			},
			Platform: types.Platform{
				AWS: &aws.Platform{
					Region: "us-east",
				},
			},
			ControlPlane: &types.MachinePool{
				Name:     "master",
				Replicas: ptr.To[int64](3),
			},
		})

	rootCAParents := asset.Parents{}
	rootCAParents.Add(&tls.SignerKeyParams{})
	rootCA := &tls.RootCA{}
	err := rootCA.Generate(context.Background(), rootCAParents)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	master := &Master{}
	err = master.Generate(context.Background(), parents)
	assert.NoError(t, err, "unexpected error generating master asset")
	expectedIgnitionConfigNames := []string{
		"master.ign",
	}
	actualFiles := master.Files()
	actualIgnitionConfigNames := make([]string, len(actualFiles))
	for i, f := range actualFiles {
		actualIgnitionConfigNames[i] = f.Filename
	}
	assert.Equal(t, expectedIgnitionConfigNames, actualIgnitionConfigNames, "unexpected names for master ignition configs")
}

func TestMasterGenerateAzureVarPartition(t *testing.T) {
	rootCAParents := asset.Parents{}
	rootCAParents.Add(&tls.SignerKeyParams{})
	rootCA := &tls.RootCA{}
	err := rootCA.Generate(context.Background(), rootCAParents)
	assert.NoError(t, err)

	hasVarPartition := func(m *Master) bool {
		for _, disk := range m.Config.Storage.Disks {
			for _, p := range disk.Partitions {
				if p.Label != nil && *p.Label == "var" {
					return true
				}
			}
		}
		return false
	}

	t.Run("Azure without /var DiskSetup appends var partition", func(t *testing.T) {
		ic := installconfig.MakeAsset(&types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
			BaseDomain: "test.com",
			Networking: &types.Networking{
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
			},
			Platform: types.Platform{
				Azure: &azure.Platform{Region: "eastus", CloudName: azure.PublicCloud},
			},
			ControlPlane: &types.MachinePool{
				Name:     "master",
				Replicas: ptr.To[int64](3),
			},
		})
		parents := asset.Parents{}
		parents.Add(ic, rootCA)
		master := &Master{}
		err := master.Generate(context.Background(), parents)
		assert.NoError(t, err)
		assert.True(t, hasVarPartition(master), "expected /var partition on Azure without DiskSetup")
	})

	t.Run("Azure with /var UserDefined DiskSetup skips var partition", func(t *testing.T) {
		ic := installconfig.MakeAsset(&types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
			BaseDomain: "test.com",
			Networking: &types.Networking{
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
			},
			Platform: types.Platform{
				Azure: &azure.Platform{Region: "eastus", CloudName: azure.PublicCloud},
			},
			ControlPlane: &types.MachinePool{
				Name:     "master",
				Replicas: ptr.To[int64](3),
				DiskSetup: []types.Disk{
					{
						Type: types.UserDefined,
						UserDefined: &types.DiskUserDefined{
							PlatformDiskID: "var-disk",
							MountPath:      "/var",
						},
					},
				},
			},
		})
		parents := asset.Parents{}
		parents.Add(ic, rootCA)
		master := &Master{}
		err := master.Generate(context.Background(), parents)
		assert.NoError(t, err)
		assert.False(t, hasVarPartition(master), "expected no /var partition when DiskSetup claims /var")
	})

	t.Run("Azure with non-var DiskSetup still appends var partition", func(t *testing.T) {
		ic := installconfig.MakeAsset(&types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
			BaseDomain: "test.com",
			Networking: &types.Networking{
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
			},
			Platform: types.Platform{
				Azure: &azure.Platform{Region: "eastus", CloudName: azure.PublicCloud},
			},
			ControlPlane: &types.MachinePool{
				Name:     "master",
				Replicas: ptr.To[int64](3),
				DiskSetup: []types.Disk{
					{
						Type: types.Etcd,
						Etcd: &types.DiskEtcd{
							PlatformDiskID: "etcd-disk",
						},
					},
				},
			},
		})
		parents := asset.Parents{}
		parents.Add(ic, rootCA)
		master := &Master{}
		err := master.Generate(context.Background(), parents)
		assert.NoError(t, err)
		assert.True(t, hasVarPartition(master), "expected /var partition when DiskSetup does not claim /var")
	})
}
