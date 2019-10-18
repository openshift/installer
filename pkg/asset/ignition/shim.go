package ignition

import (
	"fmt"
	"net/url"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// PointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func PointerIgnitionConfig(installConfig *types.InstallConfig, rootCA []byte, role string) *ignition.Config {
	var ignitionHost string
	CAReferences := []ignition.CaReference{}

	// TODO(egarcia): Move this logic to the master/worker/bootstrap ignition config generation code
	// and add parameters to service it
	if role == "bootstrap" {
		CAReferences = append(CAReferences, ignition.CaReference{
			Source: installConfig.AdditionalTrustBundle,
		})
	} else {
		CAReferences = append(CAReferences, ignition.CaReference{
			Source: dataurl.EncodeBytes(rootCA),
		})

		switch installConfig.Platform.Name() {
		case baremetaltypes.Name:
			// Baremetal needs to point directly at the VIP because we don't have a
			// way to configure DNS before Ignition runs.
			ignitionHost = fmt.Sprintf("%s:22623", installConfig.BareMetal.APIVIP)
		case openstacktypes.Name:
			apiVIP, err := openstackdefaults.APIVIP(installConfig.Networking)
			if err == nil {
				ignitionHost = fmt.Sprintf("%s:22623", apiVIP.String())
			} else {
				ignitionHost = fmt.Sprintf("api-int.%s:22623", installConfig.ClusterDomain())
			}
		default:
			ignitionHost = fmt.Sprintf("api-int.%s:22623", installConfig.ClusterDomain())
		}
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
					CertificateAuthorities: CAReferences,
				},
			},
		},
	}

}
