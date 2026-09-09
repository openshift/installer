/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package entra

import (
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"k8s.io/apimachinery/pkg/types"

	"github.com/Azure/azure-service-operator/v2/internal/identity"
)

// entraClient is a wrapper around generic client to keep a track of secretData used to create it
// and credentialFrom which that secret was retrieved.
type entraClient struct {
	genericClient *msgraphsdk.GraphServiceClient
	credential    *identity.Credential
}

func newEntraClient(
	client *msgraphsdk.GraphServiceClient,
	credential *identity.Credential,
) *entraClient {
	return &entraClient{
		genericClient: client,
		credential:    credential,
	}
}

func (c *entraClient) Client() *msgraphsdk.GraphServiceClient {
	return c.genericClient
}

func (c *entraClient) Credential() *identity.Credential {
	return c.credential
}

func (c *entraClient) CredentialFrom() types.NamespacedName {
	return c.credential.CredentialFrom()
}
