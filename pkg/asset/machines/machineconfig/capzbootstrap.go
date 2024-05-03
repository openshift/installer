package machineconfig

import (
	"fmt"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
)

// CAPZBootstrap lays down the bootstrap sentinel file used by CAPZ
// to determine if CAPI bootstrapping is complete:
// https://github.com/kubernetes-sigs/cluster-api-provider-azure/blob/main/azure/defaults.go#L98-L100
// As OpenShift bootstraps without CAPI, we will lay down this file
// manually to satisfy CAPZ.
func ForCAPZBootstrap(role string) (*mcfgv1.MachineConfig, error) {
	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{
				{
					Node: igntypes.Node{
						Path: "/run/cluster-api/bootstrap-success.complete",
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
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-capz-bootstrap", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}
