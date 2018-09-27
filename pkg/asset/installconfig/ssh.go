package installconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
)

const (
	none = "<none>"
)

type sshPublicKey struct {
	key string
}

var _ asset.Asset = (*sshPublicKey)(nil)

// Dependencies returns no dependencies.
func (a *sshPublicKey) Dependencies() []asset.Asset {
	return nil
}

func readSSHKey(path string) (string, error) {
	keyAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	key := string(keyAsBytes)

	err = validateOpenSSHPublicKey(key)
	if err != nil {
		return "", err
	}

	return key, nil
}

// Generate generates the SSH public key asset.
func (a *sshPublicKey) Generate(asset.Parents) error {
	if value, ok := os.LookupEnv("OPENSHIFT_INSTALL_SSH_PUB_KEY"); ok {
		if value != "" {
			if err := validateOpenSSHPublicKey(value); err != nil {
				return errors.Wrap(err, "failed to validate public key")
			}
		}
		a.key = value
		return nil
	}

	pubKeys := map[string]string{}
	if path, ok := os.LookupEnv("OPENSHIFT_INSTALL_SSH_PUB_KEY_PATH"); ok {
		key, err := readSSHKey(path)
		if err != nil {
			return errors.Wrap(err, "failed to read public key file")
		}
		pubKeys[path] = key
	} else {
		pubKeys[none] = ""
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
	}

	if len(pubKeys) == 1 {
		for _, value := range pubKeys {
			a.key = value
		}
		return nil
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
		return errors.Wrap(err, "failed UserInput for SSH public key")
	}

	a.key = pubKeys[path]
	return nil
}

// Name returns the human-friendly name of the asset.
func (a sshPublicKey) Name() string {
	return "SSH Key"
}

// validateOpenSSHPublicKey checks if the given string is a valid OpenSSH public key and returns an error if not.
// Ignores leading and trailing whitespace.
func validateOpenSSHPublicKey(v string) error {
	trimmed := strings.TrimSpace(v)

	// Don't let users hang themselves
	if isMatch(`-BEGIN [\w-]+ PRIVATE KEY-`, trimmed) {
		return errors.New("invalid SSH public key (appears to be a private key)")
	}

	if strings.Contains(trimmed, "\n") {
		return errors.New("invalid SSH public key (should not contain any newline characters)")
	}

	invalidError := errors.New("invalid SSH public key")

	keyParts := regexp.MustCompile(`\s+`).Split(trimmed, -1)
	if len(keyParts) < 2 {
		return invalidError
	}

	keyType := keyParts[0]
	keyBase64 := keyParts[1]
	if !isMatch(`^[\w-]+$`, keyType) || !isMatch(`^[A-Za-z0-9+\/]+={0,2}$`, keyBase64) {
		return invalidError
	}

	return nil
}

func isMatch(re string, v string) bool {
	return regexp.MustCompile(re).MatchString(v)
}
