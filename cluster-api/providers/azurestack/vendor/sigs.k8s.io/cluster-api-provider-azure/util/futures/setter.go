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

package futures

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// Setter interface defines methods that an object should implement in order to
// use the futures package for setting futures.
type Setter interface {
	Getter
	SetFutures(infrav1.Futures)
}

// Set sets the given future.
//
// NOTE: If a future already exists, we update it.
func Set(to Setter, future *infrav1.Future) {
	if to == nil || future == nil {
		return
	}

	// Check if the new future already exists, and update it if it does.
	futures := to.GetFutures()
	exists := false
	for i, f := range futures {
		if f.Name == future.Name && f.ServiceName == future.ServiceName {
			exists = true
			futures[i] = *future
			break
		}
	}

	// If the future does not exist, add it.
	if !exists {
		futures = append(futures, *future)
	}

	to.SetFutures(futures)
}

// Delete deletes the specified future.
func Delete(to Setter, name, service, futureType string) {
	if to == nil || name == "" || service == "" || futureType == "" {
		return
	}

	futures := to.GetFutures()
	for i, f := range futures {
		if f.Name == name && f.ServiceName == service && f.Type == futureType {
			futures = append(futures[:i], futures[i+1:]...)
			break
		}
	}

	to.SetFutures(futures)
}
