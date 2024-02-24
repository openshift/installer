package rosa

import (
	"crypto/rand"
	"fmt"
	"math/big"

	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
)

func handleErr(res *ocmerrors.Error, err error) error {
	msg := res.Reason()
	if msg == "" {
		msg = err.Error()
	}
	// Hack to always display the correct terms and conditions message
	if res.Code() == "CLUSTERS-MGMT-451" {
		msg = "You must accept the Terms and Conditions in order to continue.\n" +
			"Go to https://www.redhat.com/wapps/tnc/ackrequired?site=ocm&event=register\n" +
			"Once you accept the terms, you will need to retry the action that was blocked."
	}
	return fmt.Errorf(msg)
}

// GenerateRandomPassword generates a random password which satisfies OCM requiremts for passwords.
func GenerateRandomPassword() (string, error) {
	const (
		maxPasswordLength = 23
		lowerLetters      = "abcdefghijkmnopqrstuvwxyz"
		upperLetters      = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
		digits            = "23456789"
		all               = lowerLetters + upperLetters + digits
	)
	var password string
	for i := 0; i < maxPasswordLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", err
		}
		newchar := string(all[n.Int64()])
		if password == "" {
			password = newchar
		}
		if i < maxPasswordLength-1 {
			n, err = rand.Int(rand.Reader, big.NewInt(int64(len(password)+1)))
			if err != nil {
				return "", err
			}
			j := n.Int64()
			password = password[0:j] + newchar + password[j:]
		}
	}

	pw := []rune(password)
	for _, replace := range []int{5, 11, 17} {
		pw[replace] = '-'
	}

	return string(pw), nil
}
