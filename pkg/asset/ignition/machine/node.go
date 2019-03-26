package machine

import (
	"fmt"
	"net/url"

	ignition "github.com/coreos/ignition/config/v3_0/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
)

// pointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func pointerIgnitionConfig(installConfig *types.InstallConfig, rootCA []byte, role string) *ignition.Config {
	source := func() *url.URL {
		return &url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("api.%s:22623", installConfig.ClusterDomain()),
			Path:   fmt.Sprintf("/config/%s", role),
		}
	}().String()
	return &ignition.Config{
		Ignition: ignition.Ignition{
			Version: ignition.MaxVersion.String(),
			Config: ignition.IgnitionConfig{
				Merge: []ignition.ConfigReference{{
					Source: &source,
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
