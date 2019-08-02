package machine

import (
	"fmt"
	"net/url"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// pointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func pointerIgnitionConfig(installConfig *types.InstallConfig, rootCA []byte, role string) *ignition.Config {
	var ignitionHost string
	switch installConfig.Platform.Name() {
	case baremetaltypes.Name:
		// Baremetal needs to point directly at the VIP because we don't have a
		// way to configure DNS before Ignition runs.
		ignitionHost = fmt.Sprintf("%s:22623", installConfig.BareMetal.APIVIP)
	case gcptypes.Name:
		// For GCP we are using the public IP endpoint for ignition.
		ignitionHost = fmt.Sprintf("api.%s:22623", installConfig.ClusterDomain())
	default:
		ignitionHost = fmt.Sprintf("api-int.%s:22623", installConfig.ClusterDomain())
	}

	return &ignition.Config{
		Ignition: ignition.Ignition{
			Version: ignition.MaxVersion.String(),
			Config: ignition.IgnitionConfig{
				Append: []ignition.ConfigReference{{
					Source: func() *url.URL {
						return &url.URL{
							Scheme: "https",
							Host:   ignitionHost,
							Path:   fmt.Sprintf("/config/%s", role),
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
