/*
Copyright 2020 The Kubernetes Authors.

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

package resourceskus

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/cache/ttllru"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// Cache loads resource SKUs at the beginning of reconcile to expose
// features available on compute resources. It exposes convenience
// functionality for trawling Azure SKU capabilities. It may be adapted
// to periodically refresh data in the background.
type Cache struct {
	client Client

	// location is the Azure location for which this cache stores sku info.
	// we do lookup once per reconcile for the given cluster/location.
	location string

	// data is the cached sku information from Azure.
	// synchronization required if data is cached across reconcile calls, (i.e., refreshed in background as Runnable via mgr.Add(...))
	data []armcompute.ResourceSKU
}

// Cacher describes the ability to get and to add items to cache.
type Cacher interface {
	Get(key interface{}) (value interface{}, ok bool)
	Add(key interface{}, value interface{}) bool
}

// NewCacheFunc allows for mocking out the underlying client.
type NewCacheFunc func(azure.Authorizer, string) *Cache

var (
	_           Client = &AzureClient{}
	doOnce      sync.Once
	clientCache Cacher
)

// newCache instantiates a cache and initializes its contents.
func newCache(auth azure.Authorizer, location string) (*Cache, error) {
	cli, err := NewClient(auth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create resourceskus client")
	}
	return &Cache{
		client:   cli,
		location: location,
	}, nil
}

// GetCache either creates a new SKUs cache or returns an existing one based on the location + Authorizer HashKey().
func GetCache(auth azure.Authorizer, location string) (*Cache, error) {
	var err error
	doOnce.Do(func() {
		clientCache, err = ttllru.New(128, 24*time.Hour)
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed creating LRU cache for resourceSKUs cache")
	}

	key := location + "_" + auth.HashKey()
	c, ok := clientCache.Get(key)
	if ok {
		return c.(*Cache), nil
	}

	c, err = newCache(auth, location)
	if err != nil {
		return nil, err
	}
	_ = clientCache.Add(key, c)
	return c.(*Cache), nil
}

// NewStaticCache initializes a cache with data and no ability to refresh. Used for testing.
func NewStaticCache(data []armcompute.ResourceSKU, location string) *Cache {
	return &Cache{
		data:     data,
		location: location,
	}
}

func (c *Cache) refresh(ctx context.Context, location string) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "resourceskus.Cache.refresh")
	defer done()

	data, err := c.client.List(ctx, fmt.Sprintf("location eq '%s'", location))
	if err != nil {
		return errors.Wrap(err, "failed to refresh resource sku cache")
	}

	c.data = data

	return nil
}

// Get returns a resource SKU with the provided name and category. It
// returns an error if we could not find a match. We should consider
// enhancing this function to handle restrictions (e.g. SKU not
// supported in region), which is why it returns an error and not a
// boolean.
func (c *Cache) Get(ctx context.Context, name string, kind ResourceType) (SKU, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "resourceskus.Cache.Get")
	defer done()

	if c.data == nil {
		if err := c.refresh(ctx, c.location); err != nil {
			return SKU{}, err
		}
	}

	for _, sku := range c.data {
		if sku.Name != nil && *sku.Name == name {
			return SKU(sku), nil
		}
	}
	return SKU{}, azure.WithTerminalError(fmt.Errorf("resource sku with name '%s' and category '%s' not found in location '%s'", name, string(kind), c.location))
}

// Map invokes a function over all cached values.
func (c *Cache) Map(ctx context.Context, mapFn func(sku SKU)) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "resourceskus.Cache.Map")
	defer done()

	if c.data == nil {
		if err := c.refresh(ctx, c.location); err != nil {
			return err
		}
	}

	for i := range c.data {
		val := SKU(c.data[i])
		mapFn(val)
	}

	return nil
}

// GetZones looks at all virtual machine sizes and returns the unique
// set of zones into which some machine size may deploy. It removes
// restricted virtual machine sizes and duplicates.
func (c *Cache) GetZones(ctx context.Context, location string) ([]string, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "resourceskus.Cache.GetZones")
	defer done()

	var allZones = make(map[string]bool)
	mapFn := func(sku SKU) {
		// Look for VMs only
		if sku.ResourceType != nil && strings.EqualFold(*sku.ResourceType, string(VirtualMachines)) {
			// find matching location
			for _, locationInfo := range sku.LocationInfo {
				if !strings.EqualFold(*locationInfo.Location, location) {
					continue
				}
				// Use map for easy deletion and iteration
				availableZones := make(map[string]bool)

				// add all zones
				for _, zone := range locationInfo.Zones {
					availableZones[*zone] = true
				}

				if sku.Restrictions != nil {
					for _, restriction := range sku.Restrictions {
						// Can't deploy anything in this subscription in this location. Bail out.
						if ptr.Deref(restriction.Type, "") == armcompute.ResourceSKURestrictionsTypeLocation {
							availableZones = nil
							break
						}

						// remove restricted zones
						for _, restrictedZone := range restriction.RestrictionInfo.Zones {
							delete(availableZones, *restrictedZone)
						}
					}
				}

				// add to global list, if any exist. it's okay for the final list to be empty.
				// that means the region may not support AZ yet.
				for zone := range availableZones {
					allZones[zone] = true
				}

				break
			}
		}
	}

	if err := c.Map(ctx, mapFn); err != nil {
		return nil, err
	}

	var zones = make([]string, 0, len(allZones))
	for zone := range allZones {
		zones = append(zones, zone)
	}

	// lexical sort for testing
	sort.Strings(zones)

	return zones, nil
}

// GetZonesWithVMSize returns available zones for a virtual machine size in the given location.
func (c *Cache) GetZonesWithVMSize(ctx context.Context, size, location string) ([]string, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "resourceskus.Cache.GetZonesWithVMSize")
	defer done()

	var allZones = make(map[string]bool)
	mapFn := func(sku SKU) {
		if sku.Name != nil && strings.EqualFold(*sku.Name, size) && sku.ResourceType != nil && strings.EqualFold(*sku.ResourceType, string(VirtualMachines)) {
			// find matching location
			for _, locationInfo := range sku.LocationInfo {
				if !strings.EqualFold(*locationInfo.Location, location) {
					continue
				}
				// Use map for easy deletion and iteration
				availableZones := make(map[string]bool)

				// add all zones
				for _, zone := range locationInfo.Zones {
					availableZones[*zone] = true
				}

				if sku.Restrictions != nil {
					for _, restriction := range sku.Restrictions {
						// Can't deploy anything in this subscription in this location. Bail out.
						if ptr.Deref(restriction.Type, "") == armcompute.ResourceSKURestrictionsTypeLocation {
							availableZones = nil
							break
						}

						// remove restricted zones
						for _, restrictedZone := range restriction.RestrictionInfo.Zones {
							delete(availableZones, *restrictedZone)
						}
					}
				}

				// add to global list, if any exist. it's okay for the final list to be empty.
				// that means the region may not support AZ yet.
				for zone := range availableZones {
					allZones[zone] = true
				}

				break
			}
		}
	}

	if err := c.Map(ctx, mapFn); err != nil {
		return nil, err
	}

	var zones = make([]string, 0, len(allZones))
	for zone := range allZones {
		zones = append(zones, zone)
	}

	// lexical sort for testing
	sort.Strings(zones)

	return zones, nil
}
