/*
Copyright 2021 The Kubernetes Authors.

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

package coalescing

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/util/cache/ttllru"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type (
	// ReconcileCache uses and underlying time to live last recently used cache to track high frequency requests.
	// A reconciler should call ShouldProcess to determine if the key has expired. If the key has expired, a zero value
	// time.Time and true is returned. If the key has not expired, the expiration and false is returned. Upon successful
	// reconciliation a reconciler should call Reconciled to update the cache expiry.
	ReconcileCache struct {
		lastSuccessfulReconciliationCache ttllru.PeekingCacher
	}

	// ReconcileCacher describes an interface for determining if a request should be reconciled through a call to
	// ShouldProcess and if ok, reset the cool down through a call to Reconciled.
	ReconcileCacher interface {
		ShouldProcess(key string) (expiration time.Time, ok bool)
		Reconciled(key string)
	}

	// reconciler is the caching reconciler middleware that uses the cache.
	reconciler struct {
		upstream reconcile.Reconciler
		cache    ReconcileCacher
		log      logr.Logger
	}
)

// NewRequestCache creates a new instance of a ReconcileCache given a specified window of expiration.
func NewRequestCache(window time.Duration) (*ReconcileCache, error) {
	cache, err := ttllru.New(1024, window)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build ttllru cache")
	}

	return &ReconcileCache{
		lastSuccessfulReconciliationCache: cache,
	}, nil
}

// ShouldProcess determines if the key has expired. If the key has expired, a zero value
// time.Time and true is returned. If the key has not expired, the expiration and false is returned.
func (cache *ReconcileCache) ShouldProcess(key string) (time.Time, bool) {
	_, expiration, ok := cache.lastSuccessfulReconciliationCache.Peek(key)
	return expiration, !ok
}

// Reconciled updates the cache expiry for a given key.
func (cache *ReconcileCache) Reconciled(key string) {
	cache.lastSuccessfulReconciliationCache.Add(key, nil)
}

// NewReconciler returns a reconcile wrapper that will delay new reconcile.Requests
// after the cache expiry of the request string key.
// A successful reconciliation is defined as one where no error is returned.
func NewReconciler(upstream reconcile.Reconciler, cache ReconcileCacher, log logr.Logger) reconcile.Reconciler {
	return &reconciler{
		upstream: upstream,
		cache:    cache,
		log:      log.WithName("CoalescingReconciler"),
	}
}

// Reconcile sends a request to the upstream reconciler if the request is outside of the debounce window.
func (rc *reconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "controllers.reconciler.Reconcile",
		tele.KVP("namespace", r.Namespace),
		tele.KVP("name", r.Name),
	)
	defer done()

	log = log.WithValues("request", r.String())

	if expiration, ok := rc.cache.ShouldProcess(r.String()); !ok {
		log.V(4).Info("not processing", "expiration", expiration, "timeUntil", time.Until(expiration))
		var requeueAfter = time.Until(expiration)
		if requeueAfter < 1*time.Second {
			requeueAfter = 1 * time.Second
		}
		return reconcile.Result{RequeueAfter: requeueAfter}, nil
	}

	log.V(4).Info("processing")
	result, err := rc.upstream.Reconcile(ctx, r)
	if err != nil {
		log.V(4).Info("not successful")
		return result, err
	}

	log.V(4).Info("successful")
	rc.cache.Reconciled(r.String())
	return result, nil
}
