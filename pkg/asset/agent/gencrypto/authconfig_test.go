package gencrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthConfig_Generate(t *testing.T) {
	cases := []struct {
		name string
	}{
		{
			name: "generate-public-private-keys",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			authConfigAsset := &AuthConfig{}
			err := authConfigAsset.Generate(nil)

			assert.NoError(t, err)

			assert.Contains(t, authConfigAsset.PrivateKey, "BEGIN EC PRIVATE KEY")
			assert.Contains(t, authConfigAsset.PublicKey, "BEGIN EC PUBLIC KEY")
		})
	}
}
