package clusterapi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWIFPoolName(t *testing.T) {
	assert.Equal(t, "test-abc-wif-pool", WIFPoolName("test-abc"))
}

func TestWIFProviderName(t *testing.T) {
	assert.Equal(t, "test-abc-oidc-provider", WIFProviderName("test-abc"))
}

func TestOIDCBucketName(t *testing.T) {
	assert.Equal(t, "test-abc-oidc", OIDCBucketName("test-abc"))
}

func TestOIDCIssuerURL(t *testing.T) {
	assert.Equal(t, "https://storage.googleapis.com/test-abc-oidc", OIDCIssuerURL("test-abc"))
}

func TestBuildAudienceURI(t *testing.T) {
	uri := BuildAudienceURI("123456789", "my-pool", "my-provider")
	expected := "//iam.googleapis.com/projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider"
	assert.Equal(t, expected, uri)
}

func TestGenerateOIDCDiscoveryDoc(t *testing.T) {
	issuerURL := "https://storage.googleapis.com/test-oidc"
	doc, err := generateOIDCDiscoveryDoc(issuerURL)
	if !assert.NoError(t, err) {
		return
	}

	var parsed map[string]any
	err = json.Unmarshal(doc, &parsed)
	assert.NoError(t, err)

	assert.Equal(t, issuerURL, parsed["issuer"])
	assert.Equal(t, issuerURL+"/keys.json", parsed["jwks_uri"])
}

func TestGenerateJWKS(t *testing.T) {
	// Generate a fresh RSA key pair for testing
	testPEM := []byte("-----BEGIN PUBLIC KEY-----\n" +
		"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqyJvLlIIFOzh7gXbWFjC\n" +
		"8yAPEWXte+F1+nZa7k0XetvlSC1X6z6hnGqdUJCNh4UDMwhnfeDO4f8fsKx2gtEU\n" +
		"GJ0mruweWuOeBJguhR1DI5LitCOadh660lKc1PqNCgNVGvkInHi2IdZSAnBcPJTp\n" +
		"MuNQ+5+tbmJ+a8y2GMtxfEVIZj/AQ5mIJ+ItjDp4x10ePvC/THVRDIwYf96jNhnq\n" +
		"ijPsEE++qYchLV7aOuJ4s7F4DkJFEFh0pH1POltn9aiXBdsKDCk0/fmQqjl8p5n2\n" +
		"5hnS6EHdzdtgurZkxpF4wcHcGpchBxHWaRC9RXVbjHI0qdery6AtdjUjZT8UBxdN\n" +
		"GwIDAQAB\n" +
		"-----END PUBLIC KEY-----\n")

	jwksBytes, err := GenerateJWKS(testPEM)
	if !assert.NoError(t, err) {
		return
	}

	var jwks map[string]any
	err = json.Unmarshal(jwksBytes, &jwks)
	if !assert.NoError(t, err) {
		return
	}

	keys, ok := jwks["keys"].([]any)
	if !assert.True(t, ok) || !assert.Len(t, keys, 1) {
		return
	}

	key := keys[0].(map[string]any)
	assert.Equal(t, "RSA", key["kty"])
	assert.Equal(t, "RS256", key["alg"])
	assert.Equal(t, "sig", key["use"])
	assert.NotEmpty(t, key["kid"])
	assert.NotEmpty(t, key["n"])
	assert.NotEmpty(t, key["e"])
}
