/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package reconcilers

import (
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/Azure/azure-service-operator/v2/internal/util/predicates"
	"github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
)

// ARMReconcilerAnnotationChangedPredicate creates a predicate that emits events when annotations
// interesting to the generic ARM reconciler are changed
func ARMReconcilerAnnotationChangedPredicate() predicate.Predicate {
	return predicates.MakeSelectAnnotationChangedPredicate(
		map[string]predicates.HasAnnotationChanged{
			annotations.ReconcilePolicy: HasReconcilePolicyAnnotationChanged,
		})
}

// ARMPerResourceSecretAnnotationChangedPredicate creates a predicate that emits events when annotations
// interesting to the generic ARM reconciler are changed
func ARMPerResourceSecretAnnotationChangedPredicate() predicate.Predicate {
	return predicates.MakeSelectAnnotationChangedPredicate(
		map[string]predicates.HasAnnotationChanged{
			annotations.PerResourceSecret: HasAnnotationChanged,
		})
}

// HasAnnotationChanged returns true if the annotation has changed
func HasAnnotationChanged(old *string, new *string) bool {
	return !reflect.DeepEqual(old, new)
}
