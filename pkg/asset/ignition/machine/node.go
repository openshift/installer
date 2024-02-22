package machine

import (
	"fmt"
	"net"
	"net/url"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

const directory = "openshift"

// pointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func pointerIgnitionConfig(installConfig *types.InstallConfig, rootCA []byte, role string) *igntypes.Config {
	var ignitionHost string
	// Default platform independent ignitionHost
	ignitionHost = fmt.Sprintf("api-int.%s:22623", installConfig.ClusterDomain())
	// Update ignitionHost as necessary for platform
	switch installConfig.Platform.Name() {
	case baremetaltypes.Name:
		// Baremetal needs to point directly at the VIP because we don't have a
		// way to configure DNS before Ignition runs.
		ignitionHost = net.JoinHostPort(installConfig.BareMetal.APIVIPs[0], "22623")
	case nutanixtypes.Name:
		if len(installConfig.Nutanix.APIVIPs) > 0 {
			ignitionHost = net.JoinHostPort(installConfig.Nutanix.APIVIPs[0], "22623")
		}
	case openstacktypes.Name:
		ignitionHost = net.JoinHostPort(installConfig.OpenStack.APIVIPs[0], "22623")
	case ovirttypes.Name:
		ignitionHost = net.JoinHostPort(installConfig.Ovirt.APIVIPs[0], "22623")
	case vspheretypes.Name:
		if len(installConfig.VSphere.APIVIPs) > 0 {
			ignitionHost = net.JoinHostPort(installConfig.VSphere.APIVIPs[0], "22623")
		}
	}
	return &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Config: igntypes.IgnitionConfig{
				Merge: []igntypes.Resource{{
					Source: ignutil.StrToPtr(func() *url.URL {
						return &url.URL{
							Scheme: "https",
							Host:   ignitionHost,
							Path:   fmt.Sprintf("/config/%s", role),
						}
					}().String()),
				}},
			},
			Security: igntypes.Security{
				TLS: igntypes.TLS{
					CertificateAuthorities: []igntypes.Resource{{
						Source: ignutil.StrToPtr(dataurl.EncodeBytes(rootCA)),
					}},
				},
			},
		},
	}
}

// generatePointerMachineConfig generates a machineconfig when a user customizes
// the pointer ignition file manually in an IPI deployment
func generatePointerMachineConfig(config igntypes.Config, role string) (*mcfgv1.MachineConfig, error) {
	// Remove the merge section from the pointer config
	config.Ignition.Config.Merge = nil

	rawExt, err := ignition.ConvertToRawExtension(config)
	if err != nil {
		return nil, err
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-installer-ignition-%s", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}
