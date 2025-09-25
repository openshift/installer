package utils

import (
	"github.com/openshift-online/ocm-common/pkg/log"
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
}

// ASCII codes of important characters:
const (
	aCode    = 97
	zCode    = 122
	zeroCode = 48
	nineCode = 57

	// Number of letters and digits:
	letterCount = zCode - aCode + 1
	digitCount  = nineCode - zeroCode + 1
)

func RandomLabel(size int) string {
	value := r.Int()
	chars := make([]byte, size)
	for size > 0 {
		size--
		if size%2 == 0 {
			chars[size] = byte(aCode + value%letterCount)
			value = value / letterCount
		} else {
			chars[size] = byte(zeroCode + value%digitCount)
			value = value / digitCount
		}
	}
	return string(chars)
}

func Truncate(s string, truncateLength int) string {
	if len(s) > truncateLength {
		s = s[0:truncateLength]
	}
	return s
}

func GeneratePassword(length int) string {
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	special := "!#$^&*()-_=+{}|;:,.<>?/~`"
	allChars := lowercase + uppercase + digits + special

	var password []rune

	password = append(password, rune(lowercase[randInt(len(lowercase))]))
	password = append(password, rune(uppercase[randInt(len(uppercase))]))
	password = append(password, rune(digits[randInt(len(digits))]))
	password = append(password, rune(special[randInt(len(special))]))

	for len(password) < length {
		password = append(password, rune(allChars[randInt(len(allChars))]))
	}

	shuffleStrings(password)
	log.LogInfo("Generate squid password finished.")
	return string(password)
}

func randInt(max int) int {
	return rand.Intn(max)
}

func shuffleStrings(s []rune) {
	for i := len(s) - 1; i > 0; i-- {
		j := randInt(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}
