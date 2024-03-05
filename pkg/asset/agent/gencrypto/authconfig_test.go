package gencrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/common"
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
			parents := asset.Parents{}
			parents.Add(&common.InfraEnvID{})

			authConfigAsset := &AuthConfig{}
			err := authConfigAsset.Generate(parents)

			assert.NoError(t, err)

			assert.NotEmpty(t, authConfigAsset.PrivateKey)
			assert.NotEmpty(t, authConfigAsset.PublicKey)
			assert.NotEmpty(t, authConfigAsset.Token)
		})
	}
}
