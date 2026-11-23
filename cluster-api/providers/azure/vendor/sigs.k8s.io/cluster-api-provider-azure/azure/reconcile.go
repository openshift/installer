/*
Copyright The Kubernetes Authors.

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

package azure

import (
	"context"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// ResourceReconciler creates, updates, and deletes individual Azure resources.
type ResourceReconciler interface {
	CreateOrUpdateResource(ctx context.Context, spec ResourceSpecGetter, serviceName string) (any, error)
	DeleteResource(ctx context.Context, spec ResourceSpecGetter, serviceName string) error
}

// ReconcileAll creates or updates each resource spec independently, returning the
// most critical error. Error precedence (highest to lowest): real error >
// OperationNotDoneError > nil. The scope's put status condition is updated when done.
func ReconcileAll(
	ctx context.Context,
	reconciler ResourceReconciler,
	updater AsyncStatusUpdater,
	specs []ResourceSpecGetter,
	serviceName string,
	conditionType clusterv1beta1.ConditionType,
) error {
	if len(specs) == 0 {
		return nil
	}

	var result error
	for _, spec := range specs {
		if _, err := reconciler.CreateOrUpdateResource(ctx, spec, serviceName); err != nil {
			if !IsOperationNotDoneError(err) || result == nil {
				result = err
			}
		}
	}

	updater.UpdatePutStatus(conditionType, serviceName, result)
	return result
}

// DeleteAll deletes each resource spec independently, returning the most critical
// error. Error precedence (highest to lowest): real error >
// OperationNotDoneError > nil. The scope's delete status condition is updated when done.
func DeleteAll(
	ctx context.Context,
	reconciler ResourceReconciler,
	updater AsyncStatusUpdater,
	specs []ResourceSpecGetter,
	serviceName string,
	conditionType clusterv1beta1.ConditionType,
) error {
	if len(specs) == 0 {
		return nil
	}

	var result error
	for _, spec := range specs {
		if err := reconciler.DeleteResource(ctx, spec, serviceName); err != nil {
			if !IsOperationNotDoneError(err) || result == nil {
				result = err
			}
		}
	}

	updater.UpdateDeleteStatus(conditionType, serviceName, result)
	return result
}
