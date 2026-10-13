/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// makeObjectPredicate creates a predicate that executes the specified function for update, create, delete, and generic.
// The first parameter of the function "old" is nil in all cases except for update.
func makeObjectPredicate(f func(objectOld client.Object, objectNew client.Object) bool) objectPredicate {
	return objectPredicate{
		f: f,
	}
}

type objectPredicate struct {
	f func(objectOld client.Object, objectNew client.Object) bool
}

var _ predicate.Predicate = objectPredicate{}

func (p objectPredicate) Create(e event.CreateEvent) bool {
	return p.f(nil, e.Object)
}

func (p objectPredicate) Delete(e event.DeleteEvent) bool {
	return p.f(nil, e.Object)
}

func (p objectPredicate) Update(e event.UpdateEvent) bool {
	return p.f(e.ObjectOld, e.ObjectNew)
}

func (p objectPredicate) Generic(e event.GenericEvent) bool {
	return p.f(nil, e.Object)
}
