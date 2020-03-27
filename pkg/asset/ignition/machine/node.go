package machine

import (
	"fmt"
	"net"
	"net/url"

	igntypes "github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// pointerIgnitionConfigSpecV3 generates a spec v3 ignition config which references the remote config
// served by the machine config server.
func pointerIgnitionConfigSpecV3(installConfig *types.InstallConfig, rootCA []byte, role string) *igntypes.Config {
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

	mergeSourceURL := url.URL{
		Scheme: "https",
		Host:   ignitionHost,
		Path:   fmt.Sprintf("/config/%s", role),
	}
	mergeSource := mergeSourceURL.String()

	return &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Config: igntypes.IgnitionConfig{
				Merge: []igntypes.ConfigReference{{
					Source: &mergeSource,
				}},
			},
			Security: igntypes.Security{
				TLS: igntypes.TLS{
					CertificateAuthorities: []igntypes.CaReference{{
						Source: dataurl.EncodeBytes(rootCA),
					}},
				},
			},
		},
	}
}
