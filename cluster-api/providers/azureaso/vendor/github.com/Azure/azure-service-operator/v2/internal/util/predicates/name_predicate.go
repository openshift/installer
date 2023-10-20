/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// MakeNamePredicate creates a predicate that matches resources with a specific name only
func MakeNamePredicate(names ...string) predicate.Predicate {
	return makeObjectPredicate(func(_ client.Object, new client.Object) bool {
		if new == nil {
			return false
		}

		for _, n := range names {
			if new.GetName() == n {
				return true
			}
		}

		return false
	})
}
