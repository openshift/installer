package machineconfig

import (
	"fmt"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
)

// ForCustomNTP lays down chrony.conf with given NTP server.
func ForCustomNTP(role string, servers []string) (*mcfgv1.MachineConfig, error) {
	chronyConf, err := createChronyConf(servers)
	if err != nil {
		return nil, err
	}

	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{},
		},
	}
	ignConfig.Storage.Files = append(ignConfig.Storage.Files, ignition.FileFromString("/etc/chrony.conf", "root", 0644, chronyConf))

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
			Name: fmt.Sprintf("99-%s-chrony", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}

func createChronyConf(servers []string) (string, error) {
	lines := []string{}

	for _, server := range servers {
		lines = append(lines, fmt.Sprintf("server %s iburst", server))
	}

	unit := `%s
driftfile /var/lib/chrony/drift
makestep 1.0 3
rtcsync
logdir /var/log/chrony`
	return fmt.Sprintf(unit, strings.Join(lines, "\n")), nil
}
