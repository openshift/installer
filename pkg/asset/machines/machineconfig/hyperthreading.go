package machineconfig

import (
	"fmt"

	"github.com/clarketm/json"
	igntypes "github.com/coreos/ignition/v2/config/v3_0/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/asset/ignition"
)

// ForHyperthreadingDisabled creates the MachineConfig to disable hyperthreading.
// RHCOS ships with pivot.service that uses the `/etc/pivot/kernel-args` to override the kernel arguments for hosts.
func ForHyperthreadingDisabled(role string) (*mcfgv1.MachineConfig, error) {
	// TODO lorbus: Go back to using igntypes.MaxVersion.String() once spec 3.1 stable is available
	rawIgnitionConfig, err := json.Marshal(igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: "3.0.0",
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{
				ignition.FileFromString("/etc/pivot/kernel-args", "root", 0600, "ADD nosmt"),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshalling 99-%s-disable-hyperthreading ignition config failed: %v", role, err)
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
			Config: runtime.RawExtension{
				Raw: rawIgnitionConfig,
			},
		},
	}, nil
}
