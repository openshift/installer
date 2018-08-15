package installerassets

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/openshift/installer/pkg/validate"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

func readSSHKey(path string) (string, error) {
	keyAsBytes, err := ioutil.ReadFile(path)
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

func sshDefaulter(ctx context.Context) ([]byte, error) {
	if value, ok := os.LookupEnv("OPENSHIFT_INSTALL_SSH_PUB_KEY"); ok {
		if value != "" {
			if err := validate.SSHPublicKey(value); err != nil {
				return nil, errors.Wrap(err, "failed to validate public key")
			}
		}
		return []byte(value), nil
	}

	none := "<none>"
	pubKeys := map[string]string{}
	if path, ok := os.LookupEnv("OPENSHIFT_INSTALL_SSH_PUB_KEY_PATH"); ok {
		key, err := readSSHKey(path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read public key file")
		}
		pubKeys[path] = key
	} else {
		pubKeys[none] = ""
		home := os.Getenv("HOME")
		if home != "" {
			paths, err := filepath.Glob(filepath.Join(home, ".ssh", "*.pub"))
			if err != nil {
				return nil, errors.Wrap(err, "failed to glob for public key files")
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
		for _, value := range pubKeys {
			return []byte(value), nil
		}
	}

	var paths []string
	for path := range pubKeys {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var path string
	if err := survey.AskOne(&survey.Select{
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
	}); err != nil {
		return nil, errors.Wrap(err, "failed UserInput for SSH public key")
	}

	return []byte(pubKeys[path]), nil
}

func init() {
	Defaults["ssh.pub"] = sshDefaulter
}
