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
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// Getter interface defines methods that an object should implement in order to
// use the futures package for getting long running operation states.
type Getter interface {
	client.Object

	// GetFutures returns the list of long running operation states for an object.
	GetFutures() infrav1.Futures
}

// Get returns the future with the given name, if the future does not exists,
// it returns nil.
func Get(from Getter, name, service, futureType string) *infrav1.Future {
	futures := from.GetFutures()
	if futures == nil {
		return nil
	}

	for _, f := range futures {
		if f.Name == name && f.ServiceName == service && f.Type == futureType {
			return &f
		}
	}
	return nil
}

// Has returns true if a future with the given name exists.
func Has(from Getter, name, service, futureType string) bool {
	return Get(from, name, service, futureType) != nil
}
