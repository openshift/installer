package ignition

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

const (
	// keyCertAssetKeyIndex is the index of the private key in a key-pair asset.
	keyCertAssetKeyIndex = 0
	// keyCertAssetCrtIndex is the index of the public key in a key-pair asset.
	keyCertAssetCrtIndex = 1
)

// fileFromContents creates an ignition-config file with the contents from the
// specified index in the specified asset state.
func filesFromContents(pathPrefix string, mode int, contents []asset.Content) []ignition.File {
	var files []ignition.File
	for _, c := range contents {
		files = append(files, fileFromBytes(filepath.Join(pathPrefix, c.Name), mode, c.Data))
	}
	return files
}

// fileFromString creates an ignition-config file with the given contents.
func fileFromString(path string, mode int, contents string) ignition.File {
	return fileFromBytes(path, mode, []byte(contents))
}

// fileFromBytes creates an ignition-config file with the given contents.
func fileFromBytes(path string, mode int, contents []byte) ignition.File {
	return ignition.File{
		Node: ignition.Node{
			Filesystem: "root",
			Path:       path,
		},
		FileEmbedded1: ignition.FileEmbedded1{
			Mode: &mode,
			Contents: ignition.FileContents{
				Source: dataurl.EncodeBytes(contents),
			},
		},
	}
}

// masterCount determines the number of master nodes from the install config,
// defaulting to one if it is unspecified.
func masterCount(installConfig *types.InstallConfig) int {
	for _, m := range installConfig.Machines {
		if m.Name == "master" && m.Replicas != nil {
			return int(*m.Replicas)
		}
	}
	return 1
}

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
