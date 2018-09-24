package installconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/AlecAivazis/survey"

	"github.com/openshift/installer/installer/pkg/validate"
	"github.com/openshift/installer/pkg/asset"
)

const (
	none = "<none>"
)

type sshPublicKey struct{}

// Dependencies returns no dependencies.
func (a *sshPublicKey) Dependencies() []asset.Asset {
	return nil
}

// Generate generates the SSH public key asset.
func (a *sshPublicKey) Generate(map[asset.Asset]*asset.State) (state *asset.State, err error) {
	if value, ok := os.LookupEnv("OPENSHIFT_INSTALL_SSH_PUB_KEY"); ok {
		return &asset.State{
			Contents: []asset.Content{{
				Data: []byte(value),
			}},
		}, nil
	}

	pubKeys := map[string][]byte{
		none: {},
	}
	home := os.Getenv("HOME")
	if home != "" {
		paths, err := filepath.Glob(filepath.Join(home, ".ssh", "*.pub"))
		if err != nil {
			return nil, err
		}

		for _, path := range paths {
			pubKeyBytes, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}
			pubKey := string(pubKeyBytes)

			err = validate.OpenSSHPublicKey(pubKey)
			if err != nil {
				continue
			}

			pubKeys[path] = pubKeyBytes
		}
	}

	var paths []string
	for path := range pubKeys {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var path string
	survey.AskOne(&survey.Select{
		Message: "SSH Public Key",
		Help:    "The SSH key used to access all nodes within the cluster. This is optional.",
		Options: paths,
		Default: none,
	}, &path, nil)

	return &asset.State{
		Contents: []asset.Content{{
			Data: pubKeys[path],
		}},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (a *sshPublicKey) Name() string {
	return "SSH Key"
}
