package machine

import (
	"encoding/json"
	"testing"

	igntypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/openstack"
)

// TestMasterGenerate tests generating the master asset.
func TestMasterGenerate(t *testing.T) {
	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
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
				Replicas: pointer.Int64Ptr(3),
			},
		},
	}

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	master := &Master{}
	err = master.Generate(parents)
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

func TestMasterGenerateCiscoAci(t *testing.T) {
	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			TypeMeta: metav1.TypeMeta{
					APIVersion: types.InstallConfigVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			BaseDomain: "test-domain",
			Networking: &types.Networking{
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
				MachineNetwork: []types.MachineNetworkEntry{
					{
						CIDR: *ipnet.MustParseCIDR("1.2.3.0/5"),
					},
				},
				NetworkType: "CiscoAci",
			},
			Platform: types.Platform{
				OpenStack: &openstack.Platform{
					Region: "us-east",
					AciNetExt: openstack.AciNetExtStruct{
						InfraVLAN: "4094",
						ServiceVLAN: "1022",
						KubeApiVLAN: "1021",
						NeutronCIDR: ipnet.MustParseCIDR("5.6.7.8/5"),
						InstallerHostSubnet: "9.10.11.12/10",
						Mtu: "1500",
					},
				},
			},
			ControlPlane: &types.MachinePool{
				Name:     "master",
				Replicas: pointer.Int64Ptr(3),
			},
		},
	}

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	master := &Master{}
	err = master.Generate(parents)
	assert.NoError(t, err, "unexpected error generating master asset")
	expectedIgnitionConfigNames := []string{
		"master.ign",
	}

	actualFiles := master.Files()
	actualIgnitionConfigNames := make([]string, len(actualFiles))
	var masterFileData []byte
	for i, f := range actualFiles {
		actualIgnitionConfigNames[i] = f.Filename
		masterFileData = f.Data
	}
	assert.Equal(t, expectedIgnitionConfigNames, actualIgnitionConfigNames, "unexpected names for master ignition configs")

	var ignConfig *igntypes.Config
	json.Unmarshal(masterFileData, &ignConfig)

	ignition.CheckIgnitionFiles(t, ignConfig)
	ignition.CheckSystemdUnitFiles(t, ignConfig.Systemd.Units)
}
