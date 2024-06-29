package gencrypto

import (
	"context"
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
			err := authConfigAsset.Generate(context.Background(), parents)

			assert.NoError(t, err)

			assert.Contains(t, authConfigAsset.PrivateKey, "BEGIN EC PRIVATE KEY")
			assert.Contains(t, authConfigAsset.PublicKey, "BEGIN EC PUBLIC KEY")
			assert.NotEmpty(t, authConfigAsset.Token)
		})
	}
}
