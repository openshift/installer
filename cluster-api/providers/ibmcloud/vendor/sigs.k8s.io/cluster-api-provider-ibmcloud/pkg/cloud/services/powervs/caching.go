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

package powervs

import (
	"time"

	"k8s.io/client-go/tools/cache"
)

// CacheTTL is duration of time to store the vm ip in cache
// Currently the default sync period is 10 minutes that means every 10 minutes
// there will be a reconciliation, So setting cache timeout to 20 minutes so the cache updates will happen
// once in 2 reconciliations.
const CacheTTL = time.Duration(20) * time.Minute

// VMip holds the vm name and corresponding dhcp ip used to cache the dhcp ip.
type VMip struct {
	Name string
	IP   string
}

// CacheKeyFunc defines the key function required in TTLStore.
func CacheKeyFunc(obj interface{}) (string, error) {
	return obj.(VMip).Name, nil
}

// InitialiseDHCPCacheStore returns a new cache store.
func InitialiseDHCPCacheStore() cache.Store {
	return cache.NewTTLStore(CacheKeyFunc, CacheTTL)
}
