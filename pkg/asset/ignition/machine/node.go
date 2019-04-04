package machine

import (
	"fmt"
	"net/url"

	"github.com/coreos/ignition/config/util"
	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// pointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func pointerIgnitionConfig(installConfig *types.InstallConfig, rootCA []byte, role string) *ignition.Config {
	config := &ignition.Config{
		Ignition: ignition.Ignition{
			Version: ignition.MaxVersion.String(),
			Config: ignition.IgnitionConfig{
				Append: []ignition.ConfigReference{{
					Source: func() *url.URL {
						return &url.URL{
							Scheme: "https",
							Host:   fmt.Sprintf("api.%s:22623", installConfig.ClusterDomain()),
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
	if installConfig.Platform.Name() == azure.Name {
		setAzureHostName(config)
	}
	return config
}

//until hostname is set by afterburn fixed
//and included in the azure image: https://github.com/coreos/afterburn/issues/197
func setAzureHostName(ignitionConfig *ignition.Config) {
	ignitionConfig.Systemd.Units = append(ignitionConfig.Systemd.Units, ignition.Unit{
		Name:     "setazurehostname.service",
		Enabled:  util.BoolToPtr(true),
		Enable:   true,
		Contents: "[Service]\nType=oneshot\nExecStart=curl -H Metadata:true \"http://169.254.169.254/metadata/instance/compute/name?api-version=2017-08-01&format=text\" -o /tmp/hostname\nExecStart=mv -f /tmp/hostname /etc/hostname\n\n[Install]\nWantedBy=multi-user.target",
	})
}
