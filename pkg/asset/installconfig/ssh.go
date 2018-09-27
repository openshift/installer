package installconfig

import (
	"fmt"
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

func readSSHKey(path string) (key []byte, err error) {
	key, err = ioutil.ReadFile(path)
	if err != nil {
		return key, err
	}

	err = validate.OpenSSHPublicKey(string(key))
	if err != nil {
		return key, err
	}

	return key, nil
}

// Generate generates the SSH public key asset.
func (a *sshPublicKey) Generate(map[asset.Asset]*asset.State) (state *asset.State, err error) {
	if value, ok := os.LookupEnv("OPENSHIFT_INSTALL_SSH_PUB_KEY"); ok {
		if value != "" {
			if err := validate.OpenSSHPublicKey(value); err != nil {
				return nil, err
			}
		}
		return &asset.State{
			Contents: []asset.Content{{
				Data: []byte(value),
			}},
		}, nil
	}

	pubKeys := map[string][]byte{}
	if path, ok := os.LookupEnv("OPENSHIFT_INSTALL_SSH_PUB_KEY_PATH"); ok {
		key, err := readSSHKey(path)
		if err != nil {
			return nil, err
		}
		pubKeys[path] = key
	} else {
		pubKeys[none] = []byte{}
		home := os.Getenv("HOME")
		if home != "" {
			paths, err := filepath.Glob(filepath.Join(home, ".ssh", "*.pub"))
			if err != nil {
				return nil, err
			}
			for _, path := range paths {
				key, err := readSSHKey(path)
				if err != nil {
					continue
				}
				pubKeys[path] = key
			}
		}
	}

	if len(pubKeys) == 1 {
		return &asset.State{
			Contents: []asset.Content{{
				Data: []byte{},
			}},
		}, nil
	}

	var paths []string
	for path := range pubKeys {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var path string
	survey.AskOne(&survey.Select{
		Message: "SSH Public Key",
		Help:    "The SSH public key used to access all nodes within the cluster. This is optional.",
		Options: paths,
		Default: none,
	}, &path, func(ans interface{}) error {
		choice := ans.(string)
		i := sort.SearchStrings(paths, choice)
		if i == len(paths) || paths[i] != choice {
			return fmt.Errorf("invalid path %q", choice)
		}
		return nil
	})

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
