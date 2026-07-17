package gcp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return b
}

func TestValidateCredentialURLs(t *testing.T) {
	buildCreds := func(tokenURL, saImpURL, credSourceURL string) []byte {
		creds := map[string]interface{}{ //nolint:gosec
			"type":               "external_account",
			"audience":           "//iam.googleapis.com/projects/123/locations/global/workloadIdentityPools/pool/providers/provider",
			"subject_token_type": "urn:ietf:params:oauth:token-type:jwt",
		}
		if tokenURL != "" {
			creds["token_url"] = tokenURL
		}
		if saImpURL != "" {
			creds["service_account_impersonation_url"] = saImpURL
		}
		if credSourceURL != "" {
			creds["credential_source"] = map[string]interface{}{
				"url": credSourceURL,
			}
		}
		return mustMarshal(t, creds)
	}

	cases := []struct {
		name  string
		creds []byte
		err   string
	}{
		{
			name:  "valid standard WIF credentials",
			creds: buildCreds("https://sts.googleapis.com/v1/token", "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/sa@proj.iam.gserviceaccount.com:generateAccessToken", "https://my-oidc-provider.example.com/token"),
		},
		{
			name:  "all URL fields empty uses safe defaults",
			creds: buildCreds("", "", ""),
		},
		{
			name:  "valid credential_source.url to external OIDC provider",
			creds: buildCreds("https://sts.googleapis.com/v1/token", "", "https://login.microsoftonline.com/tenant/oauth2/v2.0/token"),
		},
		{
			name:  "non-external_account type is not validated",
			creds: []byte(`{"type": "service_account", "token_url": "https://bad-host.example.com/steal"}`),
		},
		{
			name: "custom universe domain with WIF is not supported",
			creds: func() []byte {
				c := map[string]interface{}{ //nolint:gosec
					"type":            "external_account",
					"universe_domain": "custom.example.com",
					"token_url":       "https://sts.custom.example.com/v1/token",
				}
				return mustMarshal(t, c)
			}(),
			err: `custom universe domain.*not supported.*open an RFE with Red Hat or a GitHub issue`,
		},
		{
			name: "googleapis.com universe domain with WIF is allowed",
			creds: func() []byte {
				c := map[string]interface{}{ //nolint:gosec
					"type":            "external_account",
					"universe_domain": "googleapis.com",
					"token_url":       "https://sts.googleapis.com/v1/token",
				}
				return mustMarshal(t, c)
			}(),
		},
		{
			name:  "invalid token_url pointing to unrelated host",
			creds: buildCreds("https://bad-host.example.com/v1/token", "", ""),
			err:   `token_url.*must equal https://sts\.googleapis\.com/v1/token`,
		},
		{
			name:  "invalid token_url with extra path",
			creds: buildCreds("https://sts.googleapis.com/v1/token/extra", "", ""),
			err:   `token_url.*must equal https://sts\.googleapis\.com/v1/token`,
		},
		{
			name:  "invalid token_url with custom universe domain",
			creds: buildCreds("https://sts.custom-domain.example.com/v1/token", "", ""),
			err:   `token_url.*must equal https://sts\.googleapis\.com/v1/token`,
		},
		{
			name:  "invalid service_account_impersonation_url pointing to unrelated host",
			creds: buildCreds("https://sts.googleapis.com/v1/token", "https://bad-host.example.com/capture", ""),
			err:   `service_account_impersonation_url.*must begin with https://iamcredentials\.googleapis\.com/`,
		},
		{
			name: "external_account_authorized_user with invalid token_url",
			creds: func() []byte {
				c := map[string]interface{}{ //nolint:gosec
					"type":      "external_account_authorized_user",
					"token_url": "https://bad-host.example.com/token",
				}
				return mustMarshal(t, c)
			}(),
			err: `token_url.*must equal https://sts\.googleapis\.com/v1/token`,
		},
		{
			name: "external_account_authorized_user with valid token_url",
			creds: func() []byte {
				c := map[string]interface{}{ //nolint:gosec
					"type":      "external_account_authorized_user",
					"token_url": "https://sts.googleapis.com/v1/token",
				}
				return mustMarshal(t, c)
			}(),
		},
		{
			name:  "invalid credential_source.url using HTTP",
			creds: buildCreds("https://sts.googleapis.com/v1/token", "", "http://my-oidc-provider.example.com/token"),
			err:   `credential_source\.url.*must use HTTPS`,
		},
		{
			name:  "invalid credential_source.url using file URI scheme",
			creds: buildCreds("https://sts.googleapis.com/v1/token", "", "file:///etc/shadow"),
			err:   `credential_source\.url.*must use HTTPS`,
		},
		{
			name:  "invalid JSON returns error",
			creds: []byte(`{not valid json`),
			err:   `failed to parse credentials JSON`,
		},
		{
			name: "valid credential_source.file",
			creds: func() []byte {
				c := map[string]interface{}{ //nolint:gosec
					"type":      "external_account",
					"token_url": "https://sts.googleapis.com/v1/token",
					"credential_source": map[string]interface{}{
						"file": "/var/run/secrets/openshift/serviceaccount/token",
					},
				}
				return mustMarshal(t, c)
			}(),
		},
		{
			name: "invalid credential_source.file with relative path",
			creds: func() []byte {
				c := map[string]interface{}{ //nolint:gosec
					"type":      "external_account",
					"token_url": "https://sts.googleapis.com/v1/token",
					"credential_source": map[string]interface{}{
						"file": "../../../etc/shadow",
					},
				}
				return mustMarshal(t, c)
			}(),
			err: `credential_source\.file.*must be an absolute path`,
		},
		{
			name: "invalid credential_source.file with path traversal",
			creds: func() []byte {
				c := map[string]interface{}{ //nolint:gosec
					"type":      "external_account",
					"token_url": "https://sts.googleapis.com/v1/token",
					"credential_source": map[string]interface{}{
						"file": "/var/run/../../etc/shadow",
					},
				}
				return mustMarshal(t, c)
			}(),
			err: `credential_source\.file.*must not contain path traversal`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			err := validateCredentialURLs(test.creds)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Regexp(t, test.err, err.Error())
			}
		})
	}
}
