package machineconfig

import (
	"fmt"

	igntypes "github.com/coreos/ignition/config/v2_2/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/types"
)

// ForOSEncryptionPolicy creates the MachineConfig based on OSEncryptionPolicy
func ForOSEncryptionPolicy(policy types.OSEncryptionPolicy, role string) *mcfgv1.MachineConfig {
	machineConfig := mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("50-%s-osencryption", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: igntypes.Config{
				Ignition: igntypes.Ignition{
					Version: igntypes.MaxVersion.String(),
				},
			},
		},
	}
	switch policy {
	case "", types.OSEncryptionPolicyDisabled:
		return nil
	case types.OSEncryptionPolicyTPM2:
		machineConfig.Spec.Config.Storage = igntypes.Storage{
			Files: []igntypes.File{
				ignition.FileFromString("/etc/clevis.json", "root", 0644, "{}\n"),
			},
		}
		return &machineConfig
	default:
		machineConfig.Spec.Config.Storage = igntypes.Storage{
			Files: []igntypes.File{
				ignition.FileFromString("/etc/clevis.json", "root", 0644, string(policy)),
			},
		}
		return &machineConfig
	}
}
