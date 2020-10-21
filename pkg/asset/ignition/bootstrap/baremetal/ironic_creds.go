package baremetal

import (
	"crypto/rand"
	"math/big"

	"github.com/openshift/installer/pkg/asset"
)

// IronicCreds is the asset for the ironic user credentials
type IronicCreds struct {
	Username string
	Password string
}

var _ asset.Asset = (*IronicCreds)(nil)

// Dependencies returns no dependencies.
func (a *IronicCreds) Dependencies() []asset.Asset {
	return nil
}

// Generate the ironic password
func (a *IronicCreds) Generate(asset.Parents) error {
	pw, err := generateRandomPassword()
	if err != nil {
		return err
	}
	a.Username = "bootstrap-user"
	a.Password = pw
	return nil
}

func generateRandomPassword() (string, error) {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 16
	buf := make([]rune, length)
	numChars := big.NewInt(int64(len(chars)))
	for i := range buf {
		c, err := rand.Int(rand.Reader, numChars)
		if err != nil {
			return "", err
		}
		buf[i] = chars[c.Uint64()]
	}
	return string(buf), nil
}

// Name returns the human-friendly name of the asset.
func (a *IronicCreds) Name() string {
	return "Ironic bootstrap credentials"
}
