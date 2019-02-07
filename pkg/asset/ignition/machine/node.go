package machine

import (
	"fmt"
	"net/url"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
)

// pointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func pointerIgnitionConfig(installConfig *types.InstallConfig, rootCA []byte, role types.MachineRole) *ignition.Config {
	return &ignition.Config{
		Ignition: ignition.Ignition{
			Version: ignition.MaxVersion.String(),
			Config: ignition.IgnitionConfig{
				Append: []ignition.ConfigReference{{
					Source: func() *url.URL {
						return &url.URL{
							Scheme: "https",
							Host:   fmt.Sprintf("%s-api.%s:22623", installConfig.ObjectMeta.Name, installConfig.BaseDomain),
							Path:   fmt.Sprintf("/config/%s", machineConfigOperatorMachineRole(role)),
						}
					}().String(),
				}},
			},
			Security: ignition.Security{
				TLS: ignition.TLS{
					CertificateAuthorities: []ignition.CaReference{{
						Source: dataurl.EncodeBytes(rootCA),
					}},
				},
			},
		},
	}
}

func machineConfigOperatorMachineRole(role types.MachineRole) string {
	switch role {
	case types.ControlPlaneMachineRole:
		return "master"
	case types.ComputeMachineRole:
		return "worker"
	default:
		return ""
	}
}
