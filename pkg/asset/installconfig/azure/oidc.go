package azure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// OIDCDiscoveryDocument represents the OpenID Connect discovery document.
type OIDCDiscoveryDocument struct {
	Issuer  string `json:"issuer"`
	JWKSURI string `json:"jwks_uri"`
}

// ExtractIssuerFromTokenFile reads a JWT token file and extracts the
// issuer (iss) claim without signature verification.
func ExtractIssuerFromTokenFile(tokenFilePath string) (string, error) {
	data, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file %s: %w", tokenFilePath, err)
	}

	parts := strings.Split(strings.TrimSpace(string(data)), ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid JWT format in %s: expected at least 2 parts, got %d", tokenFilePath, len(parts))
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	var claims struct {
		Issuer string `json:"iss"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	if claims.Issuer == "" {
		return "", fmt.Errorf("JWT in %s has no issuer (iss) claim", tokenFilePath)
	}

	return claims.Issuer, nil
}

// ValidateOIDCIssuer checks that the OIDC issuer URL has a valid
// discovery document and a reachable JWKS endpoint.
func ValidateOIDCIssuer(ctx context.Context, issuerURL string) error {
	discoveryURL := strings.TrimSuffix(issuerURL, "/") + "/.well-known/openid-configuration"

	client := &http.Client{Timeout: 30 * time.Second}

	doc, err := fetchJSON[OIDCDiscoveryDocument](ctx, client, discoveryURL)
	if err != nil {
		return fmt.Errorf("failed to fetch OIDC discovery document from %s: %w", discoveryURL, err)
	}

	if doc.Issuer == "" {
		return fmt.Errorf("OIDC discovery document at %s has no issuer field", discoveryURL)
	}

	if doc.Issuer != issuerURL {
		return fmt.Errorf("OIDC issuer mismatch: token has %q but discovery document has %q", issuerURL, doc.Issuer)
	}

	if doc.JWKSURI == "" {
		return fmt.Errorf("OIDC discovery document at %s has no jwks_uri field", discoveryURL)
	}

	if err := checkJWKS(ctx, client, doc.JWKSURI); err != nil {
		return fmt.Errorf("failed to validate JWKS at %s: %w", doc.JWKSURI, err)
	}

	return nil
}

// ValidateOIDCFromTokenFile extracts the issuer from the token file
// and validates that its OIDC discovery and JWKS endpoints are reachable.
func ValidateOIDCFromTokenFile(ctx context.Context, tokenFilePath string) (string, error) {
	issuer, err := ExtractIssuerFromTokenFile(tokenFilePath)
	if err != nil {
		return "", err
	}

	if err := ValidateOIDCIssuer(ctx, issuer); err != nil {
		return "", err
	}

	return issuer, nil
}

func fetchJSON[T any](ctx context.Context, client *http.Client, url string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// The URL is the OIDC issuer's discovery endpoint; fetching it is the
	// intended behavior of issuer validation, so the SSRF check is expected.
	resp, err := client.Do(req) //nolint:gosec // G704: issuer discovery fetch is intentional
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &result, nil
}

func checkJWKS(ctx context.Context, client *http.Client, jwksURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwksURL, nil)
	if err != nil {
		return err
	}

	// The jwks_uri comes from the issuer's own discovery document; fetching it
	// is the intended behavior of issuer validation, so the SSRF check is expected.
	resp, err := client.Do(req) //nolint:gosec // G704: JWKS fetch is intentional
	if err != nil {
		return fmt.Errorf("JWKS endpoint is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read JWKS response: %w", err)
	}

	var jwks struct {
		Keys []json.RawMessage `json:"keys"`
	}
	if err := json.Unmarshal(body, &jwks); err != nil {
		return fmt.Errorf("JWKS response is not valid JSON: %w", err)
	}

	if len(jwks.Keys) == 0 {
		return fmt.Errorf("JWKS has no keys")
	}

	return nil
}
