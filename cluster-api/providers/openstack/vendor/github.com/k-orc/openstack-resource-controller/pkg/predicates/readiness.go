/*
Copyright 2024 The ORC Authors.

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

package predicates

import (
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/api/v1alpha1"
)

type availabilityChanged struct {
	predicate.Funcs
}

// NewBecameAvailable filters events to only those where the object became available.
func NewBecameAvailable(log logr.Logger, specimen orcv1alpha1.ObjectWithConditions) predicate.Predicate {
	// The primary purpose of the specimen argument is to hopefully turn
	// accidental use of this predicate on an object which doesn't implement
	// ObjectWithConditions into a compile error rather than just a log message.

	getObjWithConditions := func(obj client.Object, event string) orcv1alpha1.ObjectWithConditions {
		objWithConditions, ok := obj.(orcv1alpha1.ObjectWithConditions)
		if !ok {
			log.Info("ReadinessChanged got event object which does not implement ObjectWithConditions",
				"got", fmt.Sprintf("%T", obj),
				"expected", fmt.Sprintf("%T", specimen),
				"event", event)
			return nil
		}

		return objWithConditions
	}

	log = log.WithValues("watchKind", fmt.Sprintf("%T", specimen))

	return availabilityChanged{
		predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				log := log.WithValues("name", e.Object.GetName(), "namespace", e.Object.GetNamespace())
				log.V(5).Info("Observed create")

				obj := getObjWithConditions(e.Object, "create")
				if obj == nil {
					return false
				}

				// Only reconcile if the new object is available
				available := orcv1alpha1.IsAvailable(obj)

				return available
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				log := log.WithValues("name", e.ObjectOld.GetName(), "namespace", e.ObjectOld.GetNamespace())
				log.V(5).Info("Observed update")

				oldObj := getObjWithConditions(e.ObjectOld, "update")
				newObj := getObjWithConditions(e.ObjectNew, "update")

				if oldObj == nil || newObj == nil {
					return false
				}

				// Only reconcile if object became available
				return !orcv1alpha1.IsAvailable(oldObj) && orcv1alpha1.IsAvailable(newObj)
			},
		},
	}
}
