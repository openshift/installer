/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package cel

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/cel-go/cel"
	"github.com/jellydator/ttlcache/v3"
	"github.com/pkg/errors"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	asometrics "github.com/Azure/azure-service-operator/v2/internal/metrics"
)

type ProgramCacher interface {
	Start()
	Stop()
	Get(resource reflect.Type, expression string) (*CompilationResult, error)
}

type programCacheItem struct {
	result *CompilationResult
	err    error
}

type ProgramCache struct {
	// key is <env-key>-<expression>. The key must contain a string uniquely identifying the env
	// as the same expression may have different meanings depending on the env.
	// Note that this is expressly safe to cache as
	// per https://github.com/google/cel-go?tab=readme-ov-file#parse-and-check.
	cache    *ttlcache.Cache[string, *programCacheItem]
	envCache *EnvCache
	compile  func(env *cel.Env, expression string) (*CompilationResult, error)

	metrics asometrics.CEL
	log     logr.Logger
}

var _ ProgramCacher = &ProgramCache{}

// NewProgramCache starts the program cache
func NewProgramCache(
	envCache *EnvCache,
	metrics asometrics.CEL,
	log logr.Logger,
	compile func(env *cel.Env, expression string) (*CompilationResult, error),
) *ProgramCache {
	log = log.WithName("CELProgramCache")

	cache := ttlcache.New[string, *programCacheItem](
		ttlcache.WithTTL[string, *programCacheItem](24 * time.Hour),
	)
	cache.OnInsertion(
		func(ctx context.Context, item *ttlcache.Item[string, *programCacheItem]) {
			log.V(Debug).Info("Program cache item inserted", "key", item.Key(), "expiry", item.ExpiresAt())
		})
	cache.OnEviction(
		func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[string, *programCacheItem]) {
			log.V(Debug).Info("Program cache item evicted", "key", item.Key(), "expiry", item.ExpiresAt(), "reason", reason)
		})
	return &ProgramCache{
		cache:    cache,
		envCache: envCache,
		compile:  compile,
		metrics:  metrics,
		log:      log,
	}
}

func (c *ProgramCache) Start() {
	go c.envCache.Start()
	go c.cache.Start()
}

func (c *ProgramCache) Stop() {
	c.cache.Stop()
	c.envCache.Stop()
}

func (c *ProgramCache) Get(resource reflect.Type, expression string) (*CompilationResult, error) {
	envKey := getTypeImportPath(resource)

	env, err := c.envCache.Get(resource)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CEL env")
	}
	key := fmt.Sprintf("%s-%s", envKey, expression)

	item := c.cache.Get(key)
	if item != nil {
		c.metrics.RecordProgramCacheHits(envKey)

		val := item.Value()
		// Error was cached
		if val.err != nil {
			return nil, val.err
		}

		// We found what we wanted, return it
		return val.result, nil
	}

	c.metrics.RecordEnvCacheMiss(envKey)

	start := time.Now()
	result, err := c.compile(env, expression)
	duration := time.Since(start)
	c.metrics.RecordCompilationTime(envKey, duration)

	cacheItem := &programCacheItem{
		result: result,
		err:    err,
	}
	c.cache.Set(key, cacheItem, ttlcache.DefaultTTL)

	if err != nil {
		return nil, err
	}

	return result, nil
}
