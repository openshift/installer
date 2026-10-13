package machineconfig

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
)

// ForExtraRoutes adds network routes that may be necessary for Disconnected deploy.
func ForExtraRoutes(role string, cidrs []string, network string) (*mcfgv1.MachineConfig, error) {
	serviceUnit, err := createExtraRoutesUnit(cidrs, network)
	if err != nil {
		return nil, err
	}

	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Systemd: igntypes.Systemd{
			Units: []igntypes.Unit{
				{
					Contents: &serviceUnit,
					Name:     fmt.Sprintf("99-%s-routes.service", role),
					Enabled:  ignutil.BoolToPtr(true),
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
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-routes", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}

func createExtraRoutesUnit(destCidrs []string, network string) (string, error) {
	if len(destCidrs) == 0 {
		return "", fmt.Errorf("destination CIDRs provided are empty")
	}
	_, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		return "", err
	}
	gateway, err := cidr.Host(ipnet, 1)
	if err != nil {
		return "", err
	}
	cmds := ""
	for _, cidr := range destCidrs {
		cmds += fmt.Sprintf("ExecStart=/usr/sbin/ip route add %s via %s\n", cidr, gateway)
	}
	unit := `[Unit]
Description=Add routes
After=ovs-configuration.service
[Service]
Type=oneshot
RemainAfterExit=yes
%s
[Install]
WantedBy=network-online.target
`
	return fmt.Sprintf(unit, cmds), nil
}
