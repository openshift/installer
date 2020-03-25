package machineconfig

import (
	"fmt"

	"github.com/clarketm/json"
	igntypes "github.com/coreos/ignition/v2/config/v3_0/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ForAuthorizedKeys creates the MachineConfig to set the authorized key for `core` user.
func ForAuthorizedKeys(key string, role string) (*mcfgv1.MachineConfig, error) {
	rawIgnitionConfig, err := json.Marshal(igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Passwd: igntypes.Passwd{
			Users: []igntypes.PasswdUser{{
				Name: "core", SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{igntypes.SSHAuthorizedKey(key)},
			}},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshalling 99-%s-ssh ignition config failed: %v", role, err)
	}

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
			Config: runtime.RawExtension{
				Raw: rawIgnitionConfig,
			},
		},
	}, nil
}
