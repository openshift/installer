//nolint:revive // package name matches the directory; renaming would break test access to internals
package tls

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIRIRegistryCredentialsMergeIntoPullSecret(t *testing.T) {
	existingAuth := map[string]interface{}{
		"auths": map[string]interface{}{
			"quay.io": map[string]interface{}{
				"auth": base64.StdEncoding.EncodeToString([]byte("existing:secret")),
			},
		},
	}
	existingPullSecret, err := json.Marshal(existingAuth)
	if !assert.NoError(t, err) {
		return
	}

	tests := []struct {
		name        string
		credentials *IRIRegistryCredentials
		pullSecret  string
		assertions  func(t *testing.T, merged string, err error)
	}{
		{
			name: "merge credentials into the pull secret",
			credentials: &IRIRegistryCredentials{
				Username: IRIRegistryUsername,
				Password: "test-password",
			},
			pullSecret: string(existingPullSecret),
			assertions: func(t *testing.T, merged string, err error) {
				t.Helper()
				if !assert.NoError(t, err) {
					return
				}

				var mergedSecret map[string]interface{}
				if !assert.NoError(t, json.Unmarshal([]byte(merged), &mergedSecret)) {
					return
				}
				auths, ok := mergedSecret["auths"].(map[string]interface{})
				if !assert.True(t, ok, "merged pull secret should carry auths") {
					return
				}

				// The existing entry must survive the merge.
				quayAuth, ok := auths["quay.io"].(map[string]interface{})
				if assert.True(t, ok, "existing quay.io auth should survive the merge") {
					assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("existing:secret")), quayAuth["auth"])
				}

				// The IRI entry must carry the credentials for the registry
				// behind api-int on the cluster domain.
				iriAuth, ok := auths["api-int.test-cluster.example.com:22625"].(map[string]interface{})
				if assert.True(t, ok, "merged pull secret should carry the IRI registry auth") {
					assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("openshift:test-password")), iriAuth["auth"])
				}
			},
		},
		{
			name:        "leave the pull secret alone without credentials",
			credentials: &IRIRegistryCredentials{},
			pullSecret:  string(existingPullSecret),
			assertions: func(t *testing.T, merged string, err error) {
				t.Helper()
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, string(existingPullSecret), merged, "pull secret should be returned unchanged")
			},
		},
		{
			name: "fail on an unparseable pull secret",
			credentials: &IRIRegistryCredentials{
				Username: IRIRegistryUsername,
				Password: "test-password",
			},
			pullSecret: "not json",
			assertions: func(t *testing.T, merged string, err error) {
				t.Helper()
				assert.Error(t, err, "an unparseable pull secret should fail the merge")
			},
		},
		{
			name: "fail on a pull secret without auths",
			credentials: &IRIRegistryCredentials{
				Username: IRIRegistryUsername,
				Password: "test-password",
			},
			pullSecret: `{"someOtherField": {}}`,
			assertions: func(t *testing.T, merged string, err error) {
				t.Helper()
				assert.Error(t, err, "a pull secret without auths should fail the merge")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged, err := tt.credentials.MergeIntoPullSecret(tt.pullSecret, "test-cluster.example.com")
			tt.assertions(t, merged, err)
		})
	}
}
