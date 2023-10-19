/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package reconcilers

import (
	"reflect"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/Azure/azure-service-operator/v2/internal/util/predicates"
	"github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
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

// GetReconcilePolicy gets the reconcile-policy from the ReconcilePolicy
func GetReconcilePolicy(obj genruntime.MetaObject, log logr.Logger) annotations.ReconcilePolicyValue {
	policyStr := obj.GetAnnotations()[annotations.ReconcilePolicy]
	policy, err := ParseReconcilePolicy(policyStr)
	if err != nil {
		log.Error(
			err,
			"failed to get reconcile policy. Applying default policy instead",
			"chosenPolicy", policy,
			"policyAnnotation", policyStr)
	}

	return policy
}

// HasAnnotationChanged returns true if the annotation has changed
func HasAnnotationChanged(old *string, new *string) bool {
	return !reflect.DeepEqual(old, new)
}
