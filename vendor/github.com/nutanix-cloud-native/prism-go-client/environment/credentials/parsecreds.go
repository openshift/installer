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
	// TODO only single API endpoint supported
	for _, cred := range creds.Credentials {
		switch cred.Type {
		case BasicAuthCredentialType:
			basicAuthCreds := BasicAuthCredential{}
			if err := json.Unmarshal(cred.Data, &basicAuthCreds); err != nil {
				return nil, fmt.Errorf("failed to unmarshal the basic-auth data. %w", err)
			}
			pc := basicAuthCreds.PrismCentral
			if pc.Username == "" || pc.Password == "" {
				return nil, fmt.Errorf("the PrismCentral credentials data is not set")
			}
			return &types.ApiCredentials{
				Username: pc.Username,
				Password: pc.Password,
			}, nil
		default:
			return nil, fmt.Errorf("unsupported credentials type: %v", cred.Type)
		}
	}
	return nil, fmt.Errorf("no Prism credentials")
}
