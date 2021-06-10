package machineconfig

import (
	"encoding/json"
	"fmt"

	azureenv "github.com/Azure/go-autorest/autorest/azure"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset/ignition"
)

// AzureStack creates a machine config to configure the Kubelet in AzureStack environments.
// For background: https://kubernetes-sigs.github.io/cloud-provider-azure/install/configs/#azure-stack-configuration
// Specifically, this machineconfig: writes the azurestack.json file containing the endpoints and configures
// an environment variable on the kubelet pointing to the file.
func AzureStack(env azureenv.Environment, role string) (*mcfgv1.MachineConfig, error) {
	endpoints, err := json.Marshal(env)
	if err != nil {
		return nil, errors.Wrap(err, "failed marshalling while creating AzureStack machine config")
	}

	dropIn := "[Service]\nEnvironment=\"AZURE_ENVIRONMENT_FILEPATH=/etc/kubernetes/azurestackcloud.json\""

	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{
				ignition.FileFromBytes("/etc/kubernetes/azurestackcloud.json", "root", 0600, endpoints),
			},
		},
		Systemd: igntypes.Systemd{
			Units: []igntypes.Unit{
				{
					Name: "kubelet.service",
					Dropins: []igntypes.Dropin{
						{
							Name:     "10-azurestack.conf",
							Contents: &dropIn,
						},
					},
				},
			},
		},
	}

	rawExt, err := ignition.ConvertToRawExtension(ignConfig)
	if err != nil {
		return nil, err
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-azurestack", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}
