/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package arm

import (
	"context"
	"net/http"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/identity"
	"github.com/Azure/azure-service-operator/v2/internal/metrics"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// ARMClientCache is a cache for armClients to hold multiple credential clients and global credential client.
type ARMClientCache struct {
	lock sync.Mutex
	// clients allows quick lookup of an armClient for each namespace
	clients            map[string]*armClient
	cloudConfig        cloud.Configuration
	credentialProvider identity.CredentialProvider
	kubeClient         kubeclient.Client
	httpClient         *http.Client
	armMetrics         *metrics.ARMClientMetrics
}

func NewARMClientCache(
	credentialProvider identity.CredentialProvider,
	kubeClient kubeclient.Client,
	configuration cloud.Configuration,
	httpClient *http.Client,
	armMetrics *metrics.ARMClientMetrics) *ARMClientCache {

	return &ARMClientCache{
		lock:               sync.Mutex{},
		clients:            make(map[string]*armClient),
		cloudConfig:        configuration,
		kubeClient:         kubeClient,
		credentialProvider: credentialProvider,
		httpClient:         httpClient,
		armMetrics:         armMetrics,
	}
}

func (c *ARMClientCache) register(client *armClient) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.clients[client.credential.CredentialFrom().String()] = client
}

func (c *ARMClientCache) lookup(key string) (*armClient, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	client, ok := c.clients[key]
	return client, ok
}

// GetConnection finds and returns connection details to be used for a given resource
func (c *ARMClientCache) GetConnection(ctx context.Context, obj genruntime.ARMMetaObject) (Connection, error) {
	cred, err := c.credentialProvider.GetCredential(ctx, obj)
	if err != nil {
		return nil, err
	}

	client, err := c.getARMClientFromCredential(cred)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *ARMClientCache) getARMClientFromCredential(cred *identity.Credential) (*armClient, error) {
	client, ok := c.lookup(cred.CredentialFrom().String())

	if ok && cred.SecretsEqual(client.credential) {
		return client, nil
	}

	options := &genericarmclient.GenericClientOptions{
		HttpClient: c.httpClient,
		Metrics:    c.armMetrics,
	}
	newClient, err := genericarmclient.NewGenericClient(c.cloudConfig, cred.TokenCredential(), options)
	if err != nil {
		return nil, err
	}

	armClient := newARMClient(newClient, cred)
	c.register(armClient)
	return armClient, nil
}
