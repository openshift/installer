/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scope

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"k8s.io/apimachinery/pkg/util/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
)

// NewFactory creates the default scope factory. It generates service clients which make OpenStack API calls against a running cloud.
func NewFactory(maxCacheSize int) Factory {
	var c *cache.LRUExpireCache
	if maxCacheSize > 0 {
		c = cache.NewLRUExpireCache(maxCacheSize)
	}
	return &providerScopeFactory{
		clientCache: c,
	}
}

// Factory instantiates a new Scope using credentials from an IdentityRefProvider.
type Factory interface {
	// NewClientScopeFromObject creates a new scope from the first object which returns an OpenStackIdentityRef
	NewClientScopeFromObject(ctx context.Context, ctrlClient client.Client, defaultCACert []byte, logger logr.Logger, objects ...infrav1.IdentityRefProvider) (Scope, error)
}

// Scope contains arguments common to most operations.
type Scope interface {
	NewComputeClient() (clients.ComputeClient, error)
	NewVolumeClient() (clients.VolumeClient, error)
	NewImageClient() (clients.ImageClient, error)
	NewNetworkClient() (clients.NetworkClient, error)
	NewLbClient() (clients.LbClient, error)
	ProjectID() string
	ExtractToken() (*tokens.Token, error)
}

// WithLogger extends Scope with a logger.
type WithLogger struct {
	Scope

	logger logr.Logger
}

func NewWithLogger(scope Scope, logger logr.Logger) *WithLogger {
	return &WithLogger{
		Scope:  scope,
		logger: logger,
	}
}

func (s *WithLogger) Logger() logr.Logger {
	return s.logger
}
