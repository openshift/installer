package machine

import (
	"encoding/json"
	"testing"

	igntypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/conversion"
	"github.com/openshift/installer/pkg/types/openstack"
)

// TestWorkerGenerate tests generating the worker asset.
func TestWorkerGenerate(t *testing.T) {
	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			Networking: &types.Networking{
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
			},
			Platform: types.Platform{
				AWS: &aws.Platform{
					Region: "us-east",
				},
			},
		},
	}

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	worker := &Worker{}
	err = worker.Generate(parents)
	assert.NoError(t, err, "unexpected error generating worker asset")

	actualFiles := worker.Files()
	assert.Equal(t, 1, len(actualFiles), "unexpected number of files in worker state")
	assert.Equal(t, "worker.ign", actualFiles[0].Filename, "unexpected name for worker ignition config")
}

func TestWorkerGenerateCiscoAci(t *testing.T) {
	installConfig1 := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
                        ObjectMeta: metav1.ObjectMeta{
                                Name: "test-cluster",
                        },
                        TypeMeta: metav1.TypeMeta{
                                        APIVersion: types.InstallConfigVersion,
                        },
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
						Mtu: "1600",
                                        },
                                },
                        },
		},
	}

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	conversion.ConvertInstallConfig(installConfig1.Config)
	parents := asset.Parents{}
	parents.Add(installConfig1, rootCA)

	worker := &Worker{}
	err = worker.Generate(parents)
	assert.NoError(t, err, "unexpected error generating worker asset")

	actualFiles := worker.Files()
	assert.Equal(t, 1, len(actualFiles), "unexpected number of files in worker state")
	assert.Equal(t, "worker.ign", actualFiles[0].Filename, "unexpected name for worker ignition config")

	var workerFileData []byte
	actualIgnitionConfigNames := make([]string, len(actualFiles))
	for i, f := range actualFiles {
		actualIgnitionConfigNames[i] = f.Filename
		workerFileData = f.Data
	}

	var ignConfig *igntypes.Config
	json.Unmarshal(workerFileData, &ignConfig)

	ignition.CheckIgnitionFiles(t, ignConfig)
	ignition.CheckSystemdUnitFiles(t, ignConfig.Systemd.Units)
}
