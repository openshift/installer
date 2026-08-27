package tls

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"

	jose "github.com/go-jose/go-jose/v4"
)

// BuildJSONWebKeySet builds a JWKS document for the given public key.
// The key ID is derived by hashing the PKIX DER encoding of the public key with SHA-256, matching the method used by kube-apiserver.
// Reference: https://github.com/kubernetes/kubernetes/blob/0f140bf1eeaf63c155f5eba1db8db9b5d52d5467/pkg/serviceaccount/jwt.go#L89-L111.
func BuildJSONWebKeySet(pubKey *rsa.PublicKey) ([]byte, string, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to serialize public key to DER format: %w", err)
	}

	hasher := crypto.SHA256.New()
	hasher.Write(derBytes) //nolint:errcheck
	keyID := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	jwk := jose.JSONWebKey{
		Key:       pubKey,
		KeyID:     keyID,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}

	jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}}
	data, err := json.Marshal(jwks)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JWKS: %w", err)
	}
	return data, keyID, nil
}

// BuildDiscoveryDocument builds an OIDC discovery document for the given issuer URL.
func BuildDiscoveryDocument(issuerURL string) string {
	return fmt.Sprintf(`{
  "issuer": "%s",
  "jwks_uri": "%s/keys.json",
  "response_types_supported": ["id_token"],
  "subject_types_supported": ["public"],
  "id_token_signing_alg_values_supported": ["RS256"],
  "claims_supported": ["aud", "exp", "sub", "iat", "iss"]
}`, issuerURL, issuerURL)
}
