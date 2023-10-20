/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package reconcilers

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
)

// ParseReconcilePolicy parses the provided reconcile policy.
func ParseReconcilePolicy(policy string) (annotations.ReconcilePolicyValue, error) {
	switch policy {
	case "":
		return annotations.ReconcilePolicyManage, nil
	case string(annotations.ReconcilePolicyManage):
		return annotations.ReconcilePolicyManage, nil
	case string(annotations.ReconcilePolicySkip):
		return annotations.ReconcilePolicySkip, nil
	case string(annotations.ReconcilePolicyDetachOnDelete):
		return annotations.ReconcilePolicyDetachOnDelete, nil
	default:
		// Defaulting to manage.
		return annotations.ReconcilePolicyManage, errors.Errorf("%q is not a known reconcile policy", policy)
	}
}

// HasReconcilePolicyAnnotationChanged returns true if the reconcile-policy annotation has
// changed in a way that needs to trigger a reconcile.
func HasReconcilePolicyAnnotationChanged(old *string, new *string) bool {
	equal := reflect.DeepEqual(old, new)
	if equal {
		// If the annotations are equal there's been no change
		return false
	}

	oldStr := to.Value(old)
	newStr := to.Value(new)

	// We only care about transitions to or from ReconcilePolicySkip. We don't need to
	// trigger an event if ReconcilePolicyDetachOnDelete is added or removed, as that annotation
	// only applies on delete (which we will always run reconcile on).
	return oldStr == string(annotations.ReconcilePolicySkip) || newStr == string(annotations.ReconcilePolicySkip)
}
