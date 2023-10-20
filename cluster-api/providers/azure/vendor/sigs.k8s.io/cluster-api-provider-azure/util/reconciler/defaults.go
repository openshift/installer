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

package reconciler

import (
	"time"
)

const (
	// DefaultLoopTimeout is the default timeout for a reconcile loop (defaulted to the max ARM template duration).
	DefaultLoopTimeout = 90 * time.Minute
	// DefaultMappingTimeout is the default timeout for a controller request mapping func.
	DefaultMappingTimeout = 60 * time.Second
	// DefaultAzureServiceReconcileTimeout is the default timeout for an Azure service reconcile.
	DefaultAzureServiceReconcileTimeout = 12 * time.Second
	// DefaultAKSServiceReconcileTimeout is the default timeout for an AKS service reconcile.
	DefaultAKSServiceReconcileTimeout = 30 * time.Second
	// DefaultAzureCallTimeout is the default timeout for an Azure request after which an Azure operation is considered long running.
	DefaultAzureCallTimeout = 2 * time.Second
	// DefaultReconcilerRequeue is the default value for the reconcile retry.
	DefaultReconcilerRequeue = 15 * time.Second
	// DefaultHTTP429RetryAfter is a default backoff wait time when we get a HTTP 429 response with no Retry-After data.
	DefaultHTTP429RetryAfter = 1 * time.Minute
)

// DefaultedLoopTimeout will default the timeout if it is zero-valued.
func DefaultedLoopTimeout(timeout time.Duration) time.Duration {
	if timeout <= 0 {
		return DefaultLoopTimeout
	}

	return timeout
}
