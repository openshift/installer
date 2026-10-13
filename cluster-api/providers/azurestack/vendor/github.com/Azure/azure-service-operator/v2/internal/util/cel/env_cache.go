/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package cel

import (
	"context"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/cel-go/cel"
	"github.com/jellydator/ttlcache/v3"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	asometrics "github.com/Azure/azure-service-operator/v2/internal/metrics"
)

type envCacheItem struct {
	env *cel.Env
}

// EnvCache caches cel.Env's for a fixed duration.
type EnvCache struct {
	cache   *ttlcache.Cache[string, envCacheItem]
	newEnv  func(resource reflect.Type) (*cel.Env, error)
	metrics asometrics.CEL
	log     logr.Logger
}

// NewEnvCache creates a new EnvCache
func NewEnvCache(
	metrics asometrics.CEL,
	log logr.Logger,
	newEnv func(resource reflect.Type) (*cel.Env, error),
) *EnvCache {
	log = log.WithName("CELEnvCache")

	cache := ttlcache.New[string, envCacheItem](
		ttlcache.WithTTL[string, envCacheItem](24 * time.Hour), // TODO: Configurable?
	)
	cache.OnInsertion(
		func(ctx context.Context, item *ttlcache.Item[string, envCacheItem]) {
			log.V(Debug).Info("Env cache item inserted", "key", item.Key(), "expiry", item.ExpiresAt())
		})
	cache.OnEviction(
		func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[string, envCacheItem]) {
			log.V(Debug).Info("Env cache item evicted", "key", item.Key(), "expiry", item.ExpiresAt(), "reason", reason)
		})

	return &EnvCache{
		cache:   cache,
		newEnv:  newEnv,
		metrics: metrics,
		log:     log,
	}
}

func (c *EnvCache) Start() {
	go c.cache.Start()
}

func (c *EnvCache) Stop() {
	c.cache.Stop()
}

func (c *EnvCache) Get(resource reflect.Type) (*cel.Env, error) {
	key := getTypeImportPath(resource)

	item := c.cache.Get(key)
	if item != nil {
		c.metrics.RecordEnvCacheHit(key)
		// We found what we wanted, return it
		return item.Value().env, nil
	}

	c.metrics.RecordEnvCacheMiss(key)
	env, err := c.newEnv(resource)
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, envCacheItem{env: env}, ttlcache.DefaultTTL)
	return env, nil
}

func coerceList(types []reflect.Type) []any {
	anyTypes := make([]any, 0, len(types))
	for _, t := range types {
		anyTypes = append(anyTypes, t)
	}

	return anyTypes
}
