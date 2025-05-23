package gencrypto

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

var (
	assetStore1Token = map[string]string{
		agentAuthKey: "",
	}
	assetStore3Tokens = map[string]string{
		agentAuthKey:   "",
		userAuthKey:    "",
		watcherAuthKey: "",
	}
)

const (
	//nolint:gosec // no sensitive info
	encodedPubKeyPEM = "LS0tLS1CRUdJTiBFQyBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFQnR0VjU4NGpHWUo3NzkwM1FJeWlYN2k2TFBHQwpQdW5QSlF6dG5PTEtHa1k3WmFSeS94UC9ITm1hL0phTDVxV0RBdnJ1b01oWGljbVN3azdqbGlUcnZBPT0KLS0tLS1FTkQgRUMgUFVCTElDIEtFWS0tLS0tCg==" // notsecret
	//nolint:gosec // no sensitive info
	privateKey = "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIOKe7Km8GggBzE5suav6P66YNa628UKgCCDAQXdc0B+1oAoGCCqGSM49\nAwEHoUQDQgAEBttV584jGYJ77903QIyiX7i6LPGCPunPJQztnOLKGkY7ZaRy/xP/\nHNma/JaL5qWDAvruoMhXicmSwk7jliTrvA==\n-----END EC PRIVATE KEY-----\n" // notsecret
)

func TestAuthConfig_Generate(t *testing.T) {
	cases := []struct {
		name                       string
		tokenExpired               bool
		workflow                   workflow.AgentWorkflowType
		clusterSecret              func(t *testing.T) runtime.Object
		expectedNumberOfAuthTokens int
	}{
		{
			name:                       "generate-public-key-and-token-for-install-workflow",
			workflow:                   workflow.AgentWorkflowTypeInstall,
			expectedNumberOfAuthTokens: 3,
		},
		{
			name:                       "add-nodes-with-only-one-valid-agent-auth-token",
			tokenExpired:               false,
			workflow:                   workflow.AgentWorkflowTypeAddNodes,
			clusterSecret:              secretWithOnly1Token(time.Now().UTC().Add(48 * time.Hour)),
			expectedNumberOfAuthTokens: 1,
		},
		{
			name:                       "add-nodes-with-three-expired-tokens",
			tokenExpired:               true,
			workflow:                   workflow.AgentWorkflowTypeAddNodes,
			clusterSecret:              secretWith3Tokens(time.Now().UTC().Add(-48 * time.Hour)),
			expectedNumberOfAuthTokens: 3,
		},
		{
			name:                       "add-nodes-with-three-valid-tokens",
			tokenExpired:               false,
			workflow:                   workflow.AgentWorkflowTypeAddNodes,
			clusterSecret:              secretWith3Tokens(time.Now().UTC().Add(48 * time.Hour)),
			expectedNumberOfAuthTokens: 3,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			agentWorkflow := &workflow.AgentWorkflow{Workflow: tc.workflow}
			parents := asset.Parents{}
			parents.Add(agentWorkflow)

			var fakeClient *fake.Clientset
			var authConfigAsset = &AuthConfig{}

			if tc.workflow == workflow.AgentWorkflowTypeAddNodes {
				addNodesConfig := &joiner.AddNodesConfig{}
				parents.Add(addNodesConfig)

				var secret runtime.Object
				if tc.clusterSecret != nil {
					secret = tc.clusterSecret(t)
				}
				fakeClient = fake.NewSimpleClientset(secret)
				authConfigAsset = &AuthConfig{
					Client: fakeClient,
				}
			}

			err := authConfigAsset.Generate(context.Background(), parents)
			assert.NoError(t, err)

			assert.NotEmpty(t, authConfigAsset.PublicKey)
			assert.NotEmpty(t, authConfigAsset.AgentAuthToken)

			if tc.expectedNumberOfAuthTokens == 3 {
				assert.NotEmpty(t, authConfigAsset.UserAuthToken)
				assert.NotEmpty(t, authConfigAsset.WatcherAuthToken)

				// All the 3 tokens should be unique
				assert.NotEqual(t, authConfigAsset.AgentAuthToken, authConfigAsset.UserAuthToken)
				assert.NotEqual(t, authConfigAsset.AgentAuthToken, authConfigAsset.WatcherAuthToken)
				assert.NotEqual(t, authConfigAsset.UserAuthToken, authConfigAsset.WatcherAuthToken)

				// Verify each token is signed with the correct persona
				assertTokenPersona(t, authConfigAsset.AgentAuthToken, agentPersona)
				assertTokenPersona(t, authConfigAsset.UserAuthToken, userPersona)
				assertTokenPersona(t, authConfigAsset.WatcherAuthToken, watcherPersona)
			}

			if tc.workflow == workflow.AgentWorkflowTypeAddNodes {
				// Retrieve the cluster secret and verify there are no errors
				clusterSecret, err := fakeClient.CoreV1().Secrets(authTokenSecretNamespace).Get(context.Background(), authTokenSecretName, metav1.GetOptions{})
				assert.NoError(t, err)

				var secretTokens map[string]string
				if tc.tokenExpired {
					// Verify secret annotations are updated when tokens are expired.
					assertAnnotations(t, clusterSecret, false)

					if tc.expectedNumberOfAuthTokens == 3 {
						// Ensure all 3 tokens (agentAuthToken, userAuthToken, watcherAuthToken) are present for OCP 4.18+.
						secretTokens = verifyAllTokensPresent(t, clusterSecret.Data, agentAuthKey, userAuthKey, watcherAuthKey)
						// Verify expired tokens in the cluster are replaced with tokens from the asset store.
						assertTokensConsistency(t, assetStore3Tokens, secretTokens, false)
					} else {
						// Ensure only agentAuthToken is present in the secret for OCP 4.17.
						secretTokens = verifyAllTokensPresent(t, clusterSecret.Data, agentAuthKey)
						// Verify expired tokens in the cluster are replaced with tokens from the asset store.
						assertTokensConsistency(t, assetStore1Token, secretTokens, false)
					}
				} else {
					// Verify secret annotations are unchanged when tokens are valid.
					assertAnnotations(t, clusterSecret, true)

					if tc.expectedNumberOfAuthTokens == 3 {
						// Ensure all 3 tokens (agentAuthToken, userAuthToken, watcherAuthToken) are present for OCP 4.18+.
						secretTokens = verifyAllTokensPresent(t, clusterSecret.Data, agentAuthKey, userAuthKey, watcherAuthKey)
						// Verify valid tokens in the cluster match those in the asset store.
						assertTokensConsistency(t, assetStore3Tokens, secretTokens, true)
					} else {
						// Ensure only agentAuthToken is present in the secret for OCP 4.17.
						secretTokens = verifyAllTokensPresent(t, clusterSecret.Data, agentAuthKey)
						// Verify valid tokens in the cluster match those in the asset store.
						assertTokensConsistency(t, assetStore1Token, secretTokens, true)
					}
				}
			}
		})
	}
}

func secretWithOnly1Token(expiry time.Time) func(t *testing.T) runtime.Object {
	return func(t *testing.T) runtime.Object {
		t.Helper()
		agentAuthToken, err := generateToken(agentPersona, privateKey, &expiry)
		assert.NoError(t, err)
		assetStore1Token[agentAuthKey] = agentAuthToken

		return &corev1.Secret{
			ObjectMeta: generateAuthTokenObjectMeta(),
			Data: map[string][]byte{
				agentAuthKey:           []byte(agentAuthToken),
				authTokenPublicDataKey: []byte(encodedPubKeyPEM),
			},
		}
	}
}

func secretWith3Tokens(expiry time.Time) func(t *testing.T) runtime.Object {
	return func(t *testing.T) runtime.Object {
		t.Helper()

		agentAuthToken, err := generateToken(agentPersona, privateKey, &expiry)
		assert.NoError(t, err)
		assetStore3Tokens[agentAuthKey] = agentAuthToken

		userAuthToken, err := generateToken(userPersona, privateKey, &expiry)
		assert.NoError(t, err)
		assetStore3Tokens[userAuthKey] = userAuthToken

		watcherAuthToken, err := generateToken(watcherPersona, privateKey, &expiry)
		assert.NoError(t, err)
		assetStore3Tokens[watcherAuthKey] = watcherAuthToken

		return &corev1.Secret{
			ObjectMeta: generateAuthTokenObjectMeta(),
			Data: map[string][]byte{
				agentAuthKey:           []byte(agentAuthToken),
				userAuthKey:            []byte(userAuthToken),
				watcherAuthKey:         []byte(watcherAuthToken),
				authTokenPublicDataKey: []byte(encodedPubKeyPEM),
			},
		}
	}
}

// generateAuthTokenObjectMeta creates and returns the ObjectMeta for the authentication token secret,
// including its name, namespace, and annotations for "updatedAt" and "expiresAt" with predefined values.
func generateAuthTokenObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      authTokenSecretName,
		Namespace: authTokenSecretNamespace,
		Annotations: map[string]string{
			"updatedAt": "some-time-before",
			"expiresAt": "some-time-in-future",
		},
	}
}

// verifyAllTokensPresent checks that the specified tokens are present in the secret and are not empty.
// Returns a map of token keys to their corresponding string values.
func verifyAllTokensPresent(t *testing.T, secretData map[string][]byte, keys ...string) map[string]string {
	t.Helper()
	tokenMap := make(map[string]string)
	for _, key := range keys {
		token, exists := secretData[key]
		assert.True(t, exists)
		assert.NotEmpty(t, token)
		tokenMap[key] = string(token)
	}
	return tokenMap
}

// assertTokensConsistency compares the tokens in the cluster secret with those in the asset store.
// When shouldMatch is true, it ensures the tokens match exactly.
// When shouldMatch is false, it ensures the tokens are not identical.
func assertTokensConsistency(t *testing.T, assetStoreTokens map[string]string, secretTokens map[string]string, shouldMatch bool) {
	t.Helper()
	if shouldMatch {
		assert.Equal(t, assetStoreTokens, secretTokens)
	} else {
		assert.NotEqual(t, assetStoreTokens, secretTokens)
	}
}

// assertTokenPersona ensures that the token is signed with the expected persona (e.g., agentAuth).
func assertTokenPersona(t *testing.T, token string, expectedPersona string) {
	t.Helper()
	claims, err := ParseToken(token)
	assert.NoError(t, err)
	persona, ok := claims["auth_scheme"].(string)
	assert.True(t, ok)
	assert.Equal(t, expectedPersona, persona)
}

// assertAnnotations validates whether the secret annotations are up-to-date based on token validity.
// If tokensValid is true, it checks the annotations match expected values.
// If tokensValid is false, it checks the annotations differ from previous values.
func assertAnnotations(t *testing.T, clusterSecret *corev1.Secret, tokensValid bool) {
	t.Helper()
	if tokensValid {
		assert.Equal(t, clusterSecret.Annotations["updatedAt"], "some-time-before")
		assert.Equal(t, clusterSecret.Annotations["expiresAt"], "some-time-in-future")
	} else {
		assert.NotEqual(t, clusterSecret.Annotations["updatedAt"], "some-time-before")
		assert.NotEqual(t, clusterSecret.Annotations["expiresAt"], "some-time-in-future")
	}
}
