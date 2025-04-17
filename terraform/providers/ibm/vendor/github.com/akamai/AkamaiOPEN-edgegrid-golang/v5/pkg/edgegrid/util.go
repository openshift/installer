package edgegrid

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"regexp"
)

var whitespaceRegexp = regexp.MustCompile("\\s{2,}")

func stringMinifier(in string) string {
	return whitespaceRegexp.ReplaceAllString(in, " ")
}

func createSignature(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
