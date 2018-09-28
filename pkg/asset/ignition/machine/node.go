package machine

import (
	"encoding/json"
	"fmt"
	"net/url"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/types"
)

// pointerIgnitionConfig generates a config which references the remote config
// served by the machine config server.
func pointerIgnitionConfig(installConfig *types.InstallConfig, rootCA []byte, role string, query string) []byte {
	data, err := json.Marshal(ignition.Config{
		Ignition: ignition.Ignition{
			Version: ignition.MaxVersion.String(),
			Config: ignition.IgnitionConfig{
				Append: []ignition.ConfigReference{{
					Source: func() *url.URL {
						return &url.URL{
							Scheme:   "https",
							Host:     fmt.Sprintf("%s-api.%s:49500", installConfig.Name, installConfig.BaseDomain),
							Path:     fmt.Sprintf("/config/%s", role),
							RawQuery: query,
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
		// XXX: Remove this once MCO supports injecting SSH keys.
		Passwd: ignition.Passwd{
			Users: []ignition.PasswdUser{{
				Name:              "core",
				SSHAuthorizedKeys: []ignition.SSHAuthorizedKey{ignition.SSHAuthorizedKey(installConfig.Admin.SSHKey)},
			}},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal pointer Ignition config: %v", err))
	}
	return data
}
