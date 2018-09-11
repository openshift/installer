package ignition

import (
	"encoding/json"
	"fmt"
	"net/url"

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

// fileFromAsset creates an ignition-config file with the contents from the
// specified index in the specified asset state.
func fileFromAsset(path string, mode int, assetState *asset.State, contentIndex int) ignition.File {
	return fileFromBytes(path, mode, assetState.Contents[contentIndex].Data)
}

// fileFromAsset creates an ignition-config file with the given contents.
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
			Config: ignition.IgnitionConfig{
				Append: []ignition.ConfigReference{{
					Source: func() *url.URL {
						return &url.URL{
							Scheme:   "https",
							Host:     fmt.Sprintf("%s-tnc.%s:49500", installConfig.Name, installConfig.BaseDomain),
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
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal pointer Ignition config: %v", err))
	}
	return data
}
