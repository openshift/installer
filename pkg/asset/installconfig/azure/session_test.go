package azure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckCredentials(t *testing.T) {
	cases := []struct {
		name     string
		creds    Credentials
		errorMsg string
	}{
		{
			name:  "valid",
			creds: Credentials{SubscriptionID: "sub", TenantID: "ten", ClientID: "cid", ClientSecret: "sec"},
		},
		{
			name:     "missing subscriptionId",
			creds:    Credentials{TenantID: "ten"},
			errorMsg: "could not retrieve subscriptionId from auth file",
		},
		{
			name:     "missing tenantId",
			creds:    Credentials{SubscriptionID: "sub"},
			errorMsg: "could not retrieve tenantId from auth file",
		},
		{
			name:     "secret without clientId",
			creds:    Credentials{SubscriptionID: "sub", TenantID: "ten", ClientSecret: "sec"},
			errorMsg: "could not retrieve clientId from auth file",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkCredentials(tc.creds)
			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthTypeFromCredentials(t *testing.T) {
	cases := []struct {
		name     string
		creds    Credentials
		expected AuthenticationType
	}{
		{
			name:     "client certificate",
			creds:    Credentials{ClientCertificatePath: "/cert.pfx", ClientSecret: "sec"},
			expected: ClientCertificateAuth,
		},
		{
			name:     "client secret",
			creds:    Credentials{ClientID: "c", ClientSecret: "sec"},
			expected: ClientSecretAuth,
		},
		{
			name:     "managed identity",
			creds:    Credentials{SubscriptionID: "s", TenantID: "t"},
			expected: ManagedIdentityAuth,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, authTypeFromCredentials(&tc.creds))
		})
	}
}

func TestCredentialsFromAzureCLIProfile(t *testing.T) {
	t.Run("uses default subscription", func(t *testing.T) {
		dir := t.TempDir()
		writeProfile(t, dir, map[string]interface{}{
			"subscriptions": []map[string]interface{}{
				{"id": "other-sub", "tenantId": "t1", "isDefault": false},
				{"id": "default-sub", "tenantId": "t2", "isDefault": true},
			},
		})
		t.Setenv("HOME", dir)

		creds, err := credentialsFromAzureCLIProfile()
		assert.NoError(t, err)
		assert.Equal(t, "default-sub", creds.SubscriptionID)
		assert.Equal(t, "t2", creds.TenantID)
		assert.Empty(t, creds.ClientID)
		assert.Empty(t, creds.ClientSecret)
	})

	t.Run("errors when no default subscription", func(t *testing.T) {
		dir := t.TempDir()
		writeProfile(t, dir, map[string]interface{}{
			"subscriptions": []map[string]interface{}{
				{"id": "sub1", "tenantId": "t1", "isDefault": false},
			},
		})
		t.Setenv("HOME", dir)

		_, err := credentialsFromAzureCLIProfile()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no default subscription")
	})

	t.Run("strips UTF-8 BOM", func(t *testing.T) {
		dir := t.TempDir()
		azureDir := filepath.Join(dir, ".azure")
		assert.NoError(t, os.MkdirAll(azureDir, 0700))
		profile := `{"subscriptions":[{"id":"bom-sub","tenantId":"bom-tenant","isDefault":true}]}`
		assert.NoError(t, os.WriteFile(filepath.Join(azureDir, "azureProfile.json"), append([]byte("\xef\xbb\xbf"), profile...), 0600))
		t.Setenv("HOME", dir)

		creds, err := credentialsFromAzureCLIProfile()
		assert.NoError(t, err)
		assert.Equal(t, "bom-sub", creds.SubscriptionID)
		assert.Equal(t, "bom-tenant", creds.TenantID)
	})
}

func TestCredentialsFromFileOrUser(t *testing.T) {
	t.Run("reads credentials file", func(t *testing.T) {
		dir := t.TempDir()
		authFile := filepath.Join(dir, "osServicePrincipal.json")
		data, err := json.Marshal(Credentials{
			SubscriptionID: "sub", TenantID: "ten",
			ClientID: "cid", ClientSecret: "sec",
		})
		assert.NoError(t, err)
		assert.NoError(t, os.WriteFile(authFile, data, 0600))

		origPath := defaultAuthFilePath
		defaultAuthFilePath = authFile
		t.Cleanup(func() { defaultAuthFilePath = origPath })

		result, authType, err := credentialsFromFileOrUser()
		assert.NoError(t, err)
		assert.Equal(t, "sub", result.SubscriptionID)
		assert.Equal(t, ClientSecretAuth, authType)
	})

	t.Run("falls back to Azure CLI profile", func(t *testing.T) {
		dir := t.TempDir()
		origPath := defaultAuthFilePath
		defaultAuthFilePath = filepath.Join(dir, "nonexistent.json")
		t.Cleanup(func() { defaultAuthFilePath = origPath })

		writeProfile(t, dir, map[string]interface{}{
			"subscriptions": []map[string]interface{}{
				{"id": "cli-sub", "tenantId": "cli-ten", "isDefault": true},
			},
		})
		t.Setenv("HOME", dir)

		result, authType, err := credentialsFromFileOrUser()
		assert.NoError(t, err)
		assert.Equal(t, "cli-sub", result.SubscriptionID)
		assert.Equal(t, "cli-ten", result.TenantID)
		assert.Equal(t, AzureCLIAuth, authType)
	})
}

func TestAzureCLICredentialOptions(t *testing.T) {
	cases := []struct {
		name             string
		creds            Credentials
		wantNil          bool
		wantSubscription string
		wantTenantID     string
	}{
		{
			name:             "subscription takes precedence over tenant",
			creds:            Credentials{SubscriptionID: "sub", TenantID: "ten"},
			wantSubscription: "sub",
		},
		{
			name:             "subscription only",
			creds:            Credentials{SubscriptionID: "sub"},
			wantSubscription: "sub",
		},
		{
			name:         "tenant only",
			creds:        Credentials{TenantID: "ten"},
			wantTenantID: "ten",
		},
		{
			name:    "neither uses CLI current account",
			creds:   Credentials{},
			wantNil: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := azureCLICredentialOptions(&tc.creds)
			if tc.wantNil {
				assert.Nil(t, opts)
				return
			}
			assert.NotNil(t, opts)
			assert.Equal(t, tc.wantSubscription, opts.Subscription)
			assert.Equal(t, tc.wantTenantID, opts.TenantID)
		})
	}
}

func writeProfile(t *testing.T, homeDir string, profile interface{}) {
	t.Helper()
	azureDir := filepath.Join(homeDir, ".azure")
	assert.NoError(t, os.MkdirAll(azureDir, 0700))
	data, err := json.Marshal(profile)
	assert.NoError(t, err)
	assert.NoError(t, os.WriteFile(filepath.Join(azureDir, "azureProfile.json"), data, 0600))
}
