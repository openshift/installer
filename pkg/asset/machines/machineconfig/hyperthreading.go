package machineconfig

import (
	"fmt"

	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ForHyperthreadingDisabled creates the MachineConfig to disable hyperthreading.
// RHCOS ships with pivot.service that uses the `/etc/pivot/kernel-args` to override the kernel arguments for hosts.
func ForHyperthreadingDisabled(role string) *mcfgv1.MachineConfig {
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
			KernelArguments: []string{
				"nosmt",
			},
		},
	}
}
