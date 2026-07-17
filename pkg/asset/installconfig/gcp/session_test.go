package gcp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCredentialURLs(t *testing.T) {
	buildCreds := func(tokenURL, saImpURL, credSourceURL string) []byte {
		creds := map[string]interface{}{
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
		b, _ := json.Marshal(creds)
		return b
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
			creds: []byte(`{"type": "service_account", "token_url": "https://evil.com/steal"}`),
		},
		{
			name: "custom universe domain with WIF is not supported",
			creds: func() []byte {
				c := map[string]interface{}{
					"type":             "external_account",
					"universe_domain":  "custom.example.com",
					"token_url":        "https://sts.custom.example.com/v1/token",
				}
				b, _ := json.Marshal(c)
				return b
			}(),
			err: `custom universe domain.*not supported.*open an RFE with Red Hat or a GitHub issue`,
		},
		{
			name: "googleapis.com universe domain with WIF is allowed",
			creds: func() []byte {
				c := map[string]interface{}{
					"type":             "external_account",
					"universe_domain":  "googleapis.com",
					"token_url":        "https://sts.googleapis.com/v1/token",
				}
				b, _ := json.Marshal(c)
				return b
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
			name:  "invalid credential_source.url using HTTP",
			creds: buildCreds("https://sts.googleapis.com/v1/token", "", "http://my-oidc-provider.example.com/token"),
			err:   `credential_source\.url.*must use HTTPS`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			err := validateCredentialURLs(test.creds)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, test.err, err.Error())
			}
		})
	}
}
