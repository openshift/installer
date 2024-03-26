package manifests

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/validate"
)

type SSHPublicKey struct {
	Key *string
}

var _ asset.Asset = (*SSHPublicKey)(nil)

// Dependencies returns no dependencies.
func (a *SSHPublicKey) Dependencies() []asset.Asset {
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
func (a *SSHPublicKey) Generate(asset.Parents) error {
	home := os.Getenv("HOME")
	if home != "" {
		key, err := readSSHKey(filepath.Join(home, ".ssh", "id_rsa.pub"))
		if err == nil {
			a.Key = &key
		}
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a SSHPublicKey) Name() string {
	return "SSH Key"
}
