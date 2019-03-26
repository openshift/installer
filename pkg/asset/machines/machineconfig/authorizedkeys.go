package machineconfig

import (
	"fmt"

	ignition "github.com/coreos/ignition/config/v3_0/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ForAuthorizedKeys creates the MachineConfig to set the authorized key for `core` user.
func ForAuthorizedKeys(key string, role string) *mcfgv1.MachineConfig {
	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-ssh", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: ignition.Config{
				Ignition: ignition.Ignition{
					Version: ignition.MaxVersion.String(),
				},
				Passwd: ignition.Passwd{
					Users: []ignition.PasswdUser{{
						Name: "core", SSHAuthorizedKeys: []ignition.SSHAuthorizedKey{ignition.SSHAuthorizedKey(key)},
					}},
				},
			},
		},
	}
}
