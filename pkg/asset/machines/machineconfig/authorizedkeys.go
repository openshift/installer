package machineconfig

import (
	"fmt"

	ignv3_0types "github.com/coreos/ignition/v2/config/v3_0/types"
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
			Config: ignv3_0types.Config{
				Ignition: ignv3_0types.Ignition{
					Version: ignv3_0types.MaxVersion.String(),
				},
				Passwd: ignv3_0types.Passwd{
					Users: []ignv3_0types.PasswdUser{{
						Name: "core", SSHAuthorizedKeys: []ignv3_0types.SSHAuthorizedKey{ignv3_0types.SSHAuthorizedKey(key)},
					}},
				},
			},
		},
	}
}
