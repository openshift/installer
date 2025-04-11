package clusterapi

import (
	"encoding/json"
	"fmt"
	"net"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/openshiftinstall"
	"github.com/openshift/installer/pkg/types"
)

// injectInstallInfo adds information about the installer and its invoker as a
// ConfigMap to the provided bootstrap Ignition config.
func injectInstallInfo(bootstrap []byte) ([]byte, error) {
	config := &igntypes.Config{}
	if err := json.Unmarshal(bootstrap, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bootstrap Ignition config: %w", err)
	}

	cm, err := openshiftinstall.CreateInstallConfigMap("openshift-install")
	if err != nil {
		return nil, fmt.Errorf("failed to generate openshift-install config: %w", err)
	}

	config.Storage.Files = append(config.Storage.Files, ignition.FileFromString("/opt/openshift/manifests/openshift-install.yaml", "root", 0644, cm))

	ign, err := ignition.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bootstrap Ignition config: %w", err)
	}

	return ign, nil
}

func prioritizeIPv4(config *types.InstallConfig, addresses []string) string {
	if len(addresses) == 0 {
		return ""
	}

	if config.Platform.Name() == "vsphere" {
		for _, a := range addresses {
			ip := net.ParseIP(a)
			if ip.To4() != nil {
				return a
			}
		}
	}

	return addresses[0]
}
