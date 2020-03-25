package machineconfig

import (
	"fmt"

	"github.com/clarketm/json"
	igntypes "github.com/coreos/ignition/v2/config/v3_0/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ForFIPSEnabled creates the MachineConfig to enable FIPS.
// See also https://github.com/openshift/machine-config-operator/pull/889
func ForFIPSEnabled(role string) (*mcfgv1.MachineConfig, error) {
	rawIgnitionConfig, err := json.Marshal(igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshalling 99-%s-fips ignition config failed: %v", role, err)
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-fips", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: runtime.RawExtension{
				Raw: rawIgnitionConfig,
			},
			FIPS: true,
		},
	}, nil
}
