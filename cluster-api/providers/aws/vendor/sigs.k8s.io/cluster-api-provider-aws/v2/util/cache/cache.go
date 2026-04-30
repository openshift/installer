/*
Copyright 2025 The Kubernetes Authors.

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

// Package cache implements caching helper functions.
package cache

import (
	"time"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	capicache "sigs.k8s.io/cluster-api/util/cache"
)

// InstanceTypeArchitectureCacheEntry caches DescribeInstanceTypes results since they are not expected to change.
type InstanceTypeArchitectureCacheEntry struct {
	InstanceType ec2types.InstanceType
	Architecture string
}

// Key returns the cache key of a InstanceTypeArchitectureCacheEntry.
func (e InstanceTypeArchitectureCacheEntry) Key() string {
	return string(e.InstanceType)
}

// InstanceTypeArchitectureCache stores InstanceTypeArchitectureCacheEntry items.
type InstanceTypeArchitectureCache = capicache.Cache[InstanceTypeArchitectureCacheEntry]

var (
	// InstanceTypeArchitectureCacheSingleton is the singleton cache for InstanceTypeArchitectureCacheEntry items.
	// It should be used in all relevant controllers (and possibly disabled for unit tests).
	InstanceTypeArchitectureCacheSingleton InstanceTypeArchitectureCache = capicache.New[InstanceTypeArchitectureCacheEntry](2 * time.Hour)
)
