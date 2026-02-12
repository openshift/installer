package machineconfig

import (
	"encoding/base64"
	"fmt"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	ignutil "github.com/coreos/ignition/v2/config/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/types"
)

// ForNetworkConfig creates the MachineConfig to configure network settings.
func ForNetworkConfig(role string, configs []types.HostConfigEntry) (*mcfgv1.MachineConfig, error) {
	files := []igntypes.File{}
	for _, config := range configs{
		yamlNetworkConfig, err := yaml.Marshal(config.NetworkConfig)
		if err != nil {
			return nil, err
		}
		encoded := base64.StdEncoding.EncodeToString([]byte(yamlNetworkConfig))
		source := fmt.Sprintf("data:text/plain;charset=utf-8;base64,%s", encoded)
		file := igntypes.File{
			Node: igntypes.Node{
				Path: fmt.Sprintf("/etc/nmstate/openshift/%s.yml", config.Hostname),
			},
			FileEmbedded1: igntypes.FileEmbedded1{
				Contents: igntypes.Resource{
					Source: ignutil.StrToPtr(source),
				},
				Mode: ignutil.IntToPtr(0644),
			},
		}
		files = append(files, file)
	}

	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Storage: igntypes.Storage{
			Files: files,
		},
	}

	rawExt, err := ignition.ConvertToRawExtension(ignConfig)
	if err != nil {
		return nil, err
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-network-config-%s", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}
