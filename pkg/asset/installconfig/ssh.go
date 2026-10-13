package installconfig

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/validate"
)

const (
	noSSHKey = "<none>"
)

type sshPublicKey struct {
	Key string
}

var _ asset.Asset = (*sshPublicKey)(nil)

// Dependencies returns no dependencies.
func (a *sshPublicKey) Dependencies() []asset.Asset {
	return nil
}

func readSSHKey(path string) (string, error) {
	keyAsBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	key := string(keyAsBytes)

	err = validate.SSHPublicKey(key)
	if err != nil {
		return "", err
	}

	return key, nil
}

// Generate generates the SSH public key asset.
func (a *sshPublicKey) Generate(context.Context, asset.Parents) error {
	pubKeys := map[string]string{
		noSSHKey: "",
	}
	home := os.Getenv("HOME")
	if home != "" {
		paths, err := filepath.Glob(filepath.Join(home, ".ssh", "*.pub"))
		if err != nil {
			return errors.Wrap(err, "failed to glob for public key files")
		}
		for _, path := range paths {
			key, err := readSSHKey(path)
			if err != nil {
				continue
			}
			pubKeys[path] = key
		}
	}

	if len(pubKeys) == 1 {
		for _, value := range pubKeys {
			a.Key = value
		}
		return nil
	}

	var paths []string
	for path := range pubKeys {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var path string
	if err := survey.AskOne(
		&survey.Select{
			Message: "SSH Public Key",
			Help:    "The SSH public key used to access all nodes within the cluster. This is optional.",
			Options: paths,
			Default: noSSHKey,
		},
		&path,
		survey.WithValidator(func(ans interface{}) error {
			choice := ans.(core.OptionAnswer).Value
			i := sort.SearchStrings(paths, choice)
			if i == len(paths) || paths[i] != choice {
				return fmt.Errorf("invalid path %q", choice)
			}
			return nil
		}),
	); err != nil {
		return errors.Wrap(err, "failed UserInput")
	}

	a.Key = pubKeys[path]
	return nil
}

// Name returns the human-friendly name of the asset.
func (a sshPublicKey) Name() string {
	return "SSH Key"
}
