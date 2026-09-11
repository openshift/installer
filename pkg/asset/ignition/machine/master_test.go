package machine

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/openstack"
)

// TestPointerIgnitionConfigOpenStack tests that the ignition host is set
// correctly for OpenStack depending on the LoadBalancer type and DNSRecordsType.
func TestPointerIgnitionConfigOpenStack(t *testing.T) {
	cases := []struct {
		name           string
		lbType         v1.PlatformLoadBalancerType
		dnsRecordsType v1.DNSRecordsType
		expectedHost   string
	}{
		{
			name:           "user-managed LB with external DNS uses api-int FQDN",
			lbType:         v1.LoadBalancerTypeUserManaged,
			dnsRecordsType: v1.DNSRecordsTypeExternal,
			expectedHost:   "api-int.test-cluster.test-domain:22623",
		},
		{
			name:           "openshift-managed LB with internal DNS uses VIP",
			lbType:         v1.LoadBalancerTypeOpenShiftManagedDefault,
			dnsRecordsType: v1.DNSRecordsTypeInternal,
			expectedHost:   "1.2.3.4:22623",
		},
		{
			name:           "user-managed LB with internal DNS uses VIP",
			lbType:         v1.LoadBalancerTypeUserManaged,
			dnsRecordsType: v1.DNSRecordsTypeInternal,
			expectedHost:   "1.2.3.4:22623",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ic := &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cluster",
				},
				BaseDomain: "test-domain",
				Networking: &types.Networking{
					ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
				},
				Platform: types.Platform{
					OpenStack: &openstack.Platform{
						APIVIPs: []string{"1.2.3.4"},
						LoadBalancer: &v1.OpenStackPlatformLoadBalancer{
							Type: tc.lbType,
						},
						DNSRecordsType: tc.dnsRecordsType,
					},
				},
			}

			cfg := pointerIgnitionConfig(ic, []byte("fake-ca"), "master")

			require.NotEmpty(t, cfg.Ignition.Config.Merge, "expected at least one merge entry")
			src := cfg.Ignition.Config.Merge[0].Source
			require.NotNil(t, src, "expected merge source to be set")

			u, err := url.Parse(*src)
			require.NoError(t, err, "unexpected error parsing ignition source URL")

			assert.Equal(t, tc.expectedHost, u.Host)
		})
	}
}

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
		},
	)

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
