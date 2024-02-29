package utils

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const MaxPasswordLength = 23

func GenerateRandomPassword() (string, error) {
	const (
		lowerLetters = "abcdefghijkmnopqrstuvwxyz"
		upperLetters = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
		digits       = "23456789"
		all          = lowerLetters + upperLetters + digits
	)
	var password string
	for i := 0; i < MaxPasswordLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(all))))
		if err != nil {
			return "", err
		}
		newchar := string(all[n.Int64()])
		if password == "" {
			password = newchar
		}
		if i < MaxPasswordLength-1 {
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

// Encrypts the input plain-text using bcrypt which is one of the hashes accepted by HTPasswd IDP
// The same encryption is used in CS as well for hashing HTPasswd IDP payload
func GenerateHTPasswdCompatibleHash(plaintxtPassword string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(plaintxtPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedValue), nil
}
