/*
Copyright 2022 The ORC Authors.

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

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
)

// NewFactory creates the default scope factory. It generates service clients which make OpenStack API calls against a running cloud.
func NewFactory(maxCacheSize int, defaultCACert []byte) Factory {
	var c *cache.LRUExpireCache
	if maxCacheSize > 0 {
		c = cache.NewLRUExpireCache(maxCacheSize)
	}
	return &providerScopeFactory{
		clientCache:   c,
		defaultCACert: defaultCACert,
	}
}

// Factory instantiates a new Scope using credentials from an IdentityRefProvider.
type Factory interface {
	// NewClientScopeFromObject creates a new scope from the first object which returns an OpenStackIdentityRef
	NewClientScopeFromObject(ctx context.Context, ctrlClient client.Client, logger logr.Logger, objects ...orcv1alpha1.CloudCredentialsRefProvider) (Scope, error)
}

// Scope contains arguments common to most operations.
type Scope interface {
	NewComputeClient() (osclients.ComputeClient, error)
	NewImageClient() (osclients.ImageClient, error)
	NewNetworkClient() (osclients.NetworkClient, error)
	NewIdentityClient() (osclients.IdentityClient, error)
	NewVolumeClient() (osclients.VolumeClient, error)
	NewVolumeTypeClient() (osclients.VolumeTypeClient, error)
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

type CredentialsProvider interface {
	// Return the name of a secret holding cloud credentials
	GetSecretName() *string

	// Return the name of the cloud to use, which is defined in the associated credentials
	GetCloudName() *string
}
