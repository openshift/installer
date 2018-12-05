package password

import (
	"crypto/rand"
	"math/big"

	"github.com/openshift/installer/pkg/asset"
	"golang.org/x/crypto/bcrypt"
)

// KubeadminPassword is the asset for the kubeadmin user password
type KubeadminPassword struct {
	Password     string
	PasswordHash []byte
}

var _ asset.Asset = (*KubeadminPassword)(nil)

// Dependencies returns no dependencies.
func (a *KubeadminPassword) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate the kubeadmin password
func (a *KubeadminPassword) Generate(asset.Parents) error {
	err := a.generateRandomPasswordHash(23)
	if err != nil {
		return err
	}
	return nil
}

// generateRandomPasswordHash generates a hash of a random ASCII password
// 5char-5char-5char-5char
func (a *KubeadminPassword) generateRandomPasswordHash(length int) error {
	const (
		lowerLetters = "abcdefghijkmnopqrstuvwxyz"
		upperLetters = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
		digits       = "23456789"
		all          = lowerLetters + upperLetters + digits
	)
	var password string
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return err
		}
		newchar := string(all[n.Int64()])
		if password == "" {
			password = newchar
		}
		if i < length-1 {
			n, err = rand.Int(rand.Reader, big.NewInt(int64(len(password)+1)))
			if err != nil {
				return err
			}
			j := n.Int64()
			password = password[0:j] + newchar + password[j:]
		}
	}
	pw := []rune(password)
	for _, replace := range []int{5, 11, 17} {
		pw[replace] = '-'
	}
	if a.Password == "" {
		a.Password = string(pw)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.PasswordHash = bytes
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *KubeadminPassword) Name() string {
	return "Kubeadmin Password"
}
