package machineconfig

import (
	"fmt"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset/ignition"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
)

// ForHyperthreadingDisabled creates the MachineConfig to disable hyperthreading.
// RHCOS ships with pivot.service to override the kernel arguments for hosts.
func ForHyperthreadingDisabled(role string) (*mcfgv1.MachineConfig, error) {
	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
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
			Name: fmt.Sprintf("99-%s-disable-hyperthreading", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config:          rawExt,
			KernelArguments: []string{"nosmt"},
		},
	}, nil
}
