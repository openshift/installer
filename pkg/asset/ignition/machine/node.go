package machine

import (
	"fmt"
	"net"
	"net/url"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// pointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func pointerIgnitionConfig(installConfig *types.InstallConfig, rootCA []byte, role string) *ignition.Config {
	var ignitionHost string
	// Default platform independent ignitionHost
	ignitionHost = fmt.Sprintf("api-int.%s:22623", installConfig.ClusterDomain())
	// Update ignitionHost as necessary for platform
	switch installConfig.Platform.Name() {
	case baremetaltypes.Name:
		// Baremetal needs to point directly at the VIP because we don't have a
		// way to configure DNS before Ignition runs.
		ignitionHost = net.JoinHostPort(installConfig.BareMetal.APIVIP, "22623")
	case openstacktypes.Name:
		ignitionHost = net.JoinHostPort(installConfig.OpenStack.APIVIP, "22623")
	case ovirttypes.Name:
		ignitionHost = net.JoinHostPort(installConfig.Ovirt.APIVIP, "22623")
	case vspheretypes.Name:
		if installConfig.VSphere.APIVIP != "" {
			ignitionHost = net.JoinHostPort(installConfig.VSphere.APIVIP, "22623")
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
					CertificateAuthorities: []ignition.CaReference{{
						Source: dataurl.EncodeBytes(rootCA),
					}},
				},
			},
		},
	}
}
