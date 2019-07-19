package machineconfig

import (
	"fmt"

	ignv2_2types "github.com/coreos/ignition/config/v2_2/types"
	"github.com/openshift/installer/pkg/asset/ignition"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ForAdditionalTrustBundle creates the MachineConfig to set the trusted certificate bundles.
func ForAdditionalTrustBundle(certificate string, role string) *mcfgv1.MachineConfig {
	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-additionaltrustbundle", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: ignv2_2types.Config{
				Ignition: ignv2_2types.Ignition{
					Version: ignv2_2types.MaxVersion.String(),
				},
				Storage: ignv2_2types.Storage{
					Files: []ignv2_2types.File{
						ignition.FileFromString("/etc/pki/ca-trust/source/anchors/ca.crt", "root", 0600, certificate),
					},
				},
			},
		},
	}
}
