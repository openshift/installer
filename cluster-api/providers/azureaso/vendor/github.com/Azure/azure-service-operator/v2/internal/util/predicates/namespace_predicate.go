/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// MakeNamespacePredicate creates a predicate that matches resources in a specific namespace only
func MakeNamespacePredicate(namespaces ...string) predicate.Predicate {
	return makeObjectPredicate(func(_ client.Object, objectNew client.Object) bool {
		if objectNew == nil {
			return false
		}

		for _, ns := range namespaces {
			if objectNew.GetNamespace() == ns {
				return true
			}
		}

		return false
	})
}
