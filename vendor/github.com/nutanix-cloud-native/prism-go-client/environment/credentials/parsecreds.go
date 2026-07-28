package credentials

import (
	"encoding/json"
	"fmt"

	"github.com/nutanix-cloud-native/prism-go-client/environment/types"
)

func ParseCredentials(credsData []byte) (*types.ApiCredentials, error) {
	creds := &NutanixCredentials{}
	err := json.Unmarshal(credsData, &creds.Credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the credentials data. %w", err)
	}

	credentials := &types.ApiCredentials{}
	// TODO only single API endpoint supported
	if len(creds.Credentials) == 0 {
		return nil, fmt.Errorf("no credentials found")
	}

	// Do not error out for multiple API endpoint credentials, just use the first one as existing code does,
	// to avoid breaking changes in brownfield scenarios, if there are secret/configmaps/local files
	// with multiple credentials.

	cred := creds.Credentials[0]
	switch cred.Type {
	case BasicAuthCredentialType:
		basicAuthCreds := BasicAuthCredential{}
		if err := json.Unmarshal(cred.Data, &basicAuthCreds); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the basic-auth data. %w", err)
		}

		credentials.Username = basicAuthCreds.PrismCentral.Username
		credentials.Password = basicAuthCreds.PrismCentral.Password
	case APIKeyCredentialType:
		apiKeyCreds := APIKeyCredential{}
		if err := json.Unmarshal(cred.Data, &apiKeyCreds); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the api key data. %w", err)
		}

		credentials.APIKey = string(apiKeyCreds.PrismCentral.APIKey)
	default:
		return nil, fmt.Errorf("unsupported credentials type: %v", cred.Type)
	}

	// run validation for safety check
	if err := credentials.Validate(); err != nil {
		return nil, err
	}

	return credentials, nil
}
