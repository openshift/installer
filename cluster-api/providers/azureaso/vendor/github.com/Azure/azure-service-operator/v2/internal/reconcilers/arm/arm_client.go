/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package arm

import (
	"k8s.io/apimachinery/pkg/types"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/identity"
)

// armClient is a wrapper around generic client to keep a track of secretData used to create it and credentialFrom which
// that secret was retrieved.
type armClient struct {
	genericClient *genericarmclient.GenericClient
	credential    *identity.Credential
}

func newARMClient(
	client *genericarmclient.GenericClient,
	credential *identity.Credential,
) *armClient {
	return &armClient{
		genericClient: client,
		credential:    credential,
	}
}

func (c *armClient) Client() *genericarmclient.GenericClient {
	return c.genericClient
}

func (c *armClient) Credential() *identity.Credential {
	return c.credential
}

func (c *armClient) CredentialFrom() types.NamespacedName {
	return c.credential.CredentialFrom()
}

func (c *armClient) SubscriptionID() string {
	return c.credential.SubscriptionID()
}

type Connection interface {
	Client() *genericarmclient.GenericClient
	CredentialFrom() types.NamespacedName
	SubscriptionID() string
}
