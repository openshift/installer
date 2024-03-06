package utils

import (
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
