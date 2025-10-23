package machineconfig

import (
	"fmt"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/types"
)

// ForDualStackAddresses creates the MachineConfig to tell kernel to configure the IP addresses with DHCP and DHCPV6.
func ForDualStackAddresses(role string, ic *types.InstallConfig) (*mcfgv1.MachineConfig, error) {
	kubeletAddr := "0.0.0.0"
	if ic.InfraStack() == mcfgv1.IPFamiliesDualStackIPv6Primary {
		kubeletAddr = "::"
	}

	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{
				{
					Node: igntypes.Node{
						Path:      "/etc/kubernetes/kubelet-workaround",
						Overwrite: ptr.To(true),
					},
					FileEmbedded1: igntypes.FileEmbedded1{
						Mode: ptr.To(644),
						Contents: igntypes.Resource{
							Source: ptr.To(dataurl.EncodeBytes([]byte(fmt.Sprintf("KUBELET_NODE_IPS=%s", kubeletAddr)))),
						},
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
			Name: fmt.Sprintf("99-dual-stack-%s", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config:          rawExt,
			KernelArguments: []string{"ip=dhcp,dhcp6"},
		},
	}, nil
}
