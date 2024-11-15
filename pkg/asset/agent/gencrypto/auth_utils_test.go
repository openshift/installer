package gencrypto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseExpirationFromToken(t *testing.T) {
	publicKey, privateKey, err := keyPairPEM()
	assert.NotEmpty(t, publicKey)
	assert.NotEmpty(t, privateKey)
	assert.NoError(t, err)

	tokenNoExp, err := generateToken("userAuth", privateKey, nil)
	assert.NotEmpty(t, tokenNoExp)
	assert.NoError(t, err)

	expiry := time.Now().UTC().Add(30 * time.Second)
	tokenWithExp, err := generateToken("userAuth", privateKey, &expiry)
	assert.NotEmpty(t, tokenWithExp)
	assert.NoError(t, err)

	cases := []struct {
		name, token, errorMessage string
		expectedErr               bool
	}{
		{
			name:         "exp-claim-not-set",
			token:        tokenNoExp,
			expectedErr:  true,
			errorMessage: "token missing 'exp' claim",
		},
		{
			name:  "exp-claim-set",
			token: tokenWithExp,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err = ParseExpirationFromToken(tc.token)
			if tc.expectedErr {
				assert.EqualError(t, err, tc.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
