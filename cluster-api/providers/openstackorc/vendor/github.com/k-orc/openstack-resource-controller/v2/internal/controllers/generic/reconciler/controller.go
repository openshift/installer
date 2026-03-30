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

package reconciler

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/status"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
)

type ResourceController interface {
	GetName() string

	GetK8sClient() client.Client
	GetScopeFactory() scope.Factory
}

func NewController[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	}, orcObjectT any,
	resourceSpecT any, filterT any,
	objectApplyPT interfaces.ORCApplyConfig[objectApplyPT, statusApplyPT],
	statusApplyPT interface {
		*statusApplyT
		interfaces.ORCStatusApplyConfig[statusApplyPT]
	}, statusApplyT any,
	osResourceT any,
](
	name string, k8sClient client.Client, scopeFactory scope.Factory,
	helperFactory interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT],
	statusWriter interfaces.ResourceStatusWriter[orcObjectPT, *osResourceT, objectApplyPT, statusApplyPT],
) Controller[orcObjectPT, orcObjectT, resourceSpecT, filterT, objectApplyPT, statusApplyPT, statusApplyT, osResourceT] {
	return Controller[orcObjectPT, orcObjectT, resourceSpecT, filterT, objectApplyPT, statusApplyPT, statusApplyT, osResourceT]{
		name:          name,
		client:        k8sClient,
		scopeFactory:  scopeFactory,
		helperFactory: helperFactory,
		statusWriter:  statusWriter,
	}
}

type Controller[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	},
	orcObjectT any,
	resourceSpecT any,
	filterT any,
	objectApplyPT interfaces.ORCApplyConfig[objectApplyPT, statusApplyPT],
	statusApplyPT interface {
		*statusApplyT
		interfaces.ORCStatusApplyConfig[statusApplyPT]
	},
	statusApplyT any,
	osResourceT any,
] struct {
	name         string
	client       client.Client
	scopeFactory scope.Factory

	helperFactory interfaces.ResourceHelperFactory[orcObjectPT, orcObjectT, resourceSpecT, filterT, osResourceT]
	statusWriter  interfaces.ResourceStatusWriter[orcObjectPT, *osResourceT, objectApplyPT, statusApplyPT]
}

func (c *Controller[_, _, _, _, _, _, _, _]) GetName() string {
	return c.name
}

func (c *Controller[_, _, _, _, _, _, _, _]) GetK8sClient() client.Client {
	return c.client
}

func (c *Controller[_, _, _, _, _, _, _, _]) GetScopeFactory() scope.Factory {
	return c.scopeFactory
}

func (c *Controller[
	orcObjectPT, orcObjectT,
	resourceSpecT, filterT,
	objectApplyPT,
	statusApplyPT, statusApplyT,
	osResourceT,
]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var orcObject orcObjectPT = new(orcObjectT)
	err := c.client.Get(ctx, req.NamespacedName, orcObject)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	adapter := c.helperFactory.NewAPIObjectAdapter(orcObject)

	log := ctrl.LoggerFrom(ctx)
	if !orcObject.GetDeletionTimestamp().IsZero() {
		return c.reconcileDelete(ctx, adapter).Return(log)
	}

	return c.reconcileNormal(ctx, adapter).Return(log)
}

// shouldReconcile filters events when the object status is up to date, and its
// status indicates that no further reconciliation is required.
//
// Specifically it looks at the Progressing condition. It has the following behaviour:
// - Progressing condition is not present -> reconcile
// - Progressing condition is present and True -> reconcile
// - Progressing condition is present and False, but observedGeneration is old -> reconcile
// - Progressing condition is false and observedGeneration is up to date -> do not reconcile
//
// If shouldReconcile is preventing an object from being reconciled which should
// be reconciled, consider if that object's actuator is correctly returning a
// ProgressStatus indicating that the reconciliation should continue.
func shouldReconcile(obj orcv1alpha1.ObjectWithConditions) bool {
	progressing := meta.FindStatusCondition(obj.GetConditions(), orcv1alpha1.ConditionProgressing)
	if progressing == nil {
		return true
	}

	if progressing.Status == metav1.ConditionTrue {
		return true
	}

	return progressing.ObservedGeneration != obj.GetGeneration()
}

func (c *Controller[
	orcObjectPT, orcObjectT,
	resourceSpecT, filterT,
	objectApplyPT,
	statusApplyPT, statusApplyT,
	osResourceT,
]) reconcileNormal(ctx context.Context, objAdapter interfaces.APIObjectAdapter[orcObjectPT, resourceSpecT, filterT]) (reconcileStatus progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// We do this here rather than in a predicate because predicates only cover
	// a single watch. Doing it here means we cover all sources of
	// reconciliation, including our dependencies.
	if !shouldReconcile(objAdapter.GetObject()) {
		log.V(logging.Verbose).Info("Status is up to date: not reconciling")
		return reconcileStatus
	}

	log.V(logging.Verbose).Info("Reconciling resource")

	var osResource *osResourceT

	// Ensure we always update status
	defer func() {
		reconcileStatus = reconcileStatus.WithReconcileStatus(
			status.UpdateStatus(ctx, c, c.statusWriter, objAdapter.GetObject(), osResource, reconcileStatus))
	}()

	actuator, actuatorRS := c.helperFactory.NewCreateActuator(ctx, objAdapter.GetObject(), c)
	if needsReschedule, err := actuatorRS.NeedsReschedule(); needsReschedule {
		if err == nil {
			log.V(logging.Verbose).Info("Waiting on events before creation")
		}
		return actuatorRS.WithReconcileStatus(reconcileStatus)
	}

	osResource, getOSResourceRS := GetOrCreateOSResource(ctx, log, c, objAdapter, actuator)
	if needsReschedule, err := getOSResourceRS.NeedsReschedule(); needsReschedule {
		if err == nil {
			log.V(logging.Verbose).Info("Waiting on events before creation")
		}
		return getOSResourceRS.WithReconcileStatus(reconcileStatus)
	}

	if osResource == nil {
		// Programming error: if we don't have a resource we should either have an error or be waiting on something
		return reconcileStatus.WithError(fmt.Errorf("oResource is not set, but no wait events or error"))
	}

	if objAdapter.GetStatusID() == nil {
		resourceID := actuator.GetResourceID(osResource)
		if err := status.SetStatusID(ctx, c, objAdapter.GetObject(), resourceID, c.statusWriter); err != nil {
			return reconcileStatus.WithError(err)
		}
	}

	log = log.WithValues("ID", actuator.GetResourceID(osResource))
	log.V(logging.Debug).Info("Got resource")
	ctx = ctrl.LoggerInto(ctx, log)

	if objAdapter.GetManagementPolicy() == orcv1alpha1.ManagementPolicyManaged {
		if reconciler, ok := actuator.(interfaces.ReconcileResourceActuator[orcObjectPT, osResourceT]); ok {
			// We deliberately execute all reconcilers returned by GetResourceReconcilers, even if it returns an error.
			reconcilers, getReconcilersRS := reconciler.GetResourceReconcilers(ctx, objAdapter.GetObject(), osResource, c)
			reconcileStatus = getReconcilersRS.WithReconcileStatus(reconcileStatus)

			// We execute all returned updaters, even if some return errors
			for _, updater := range reconcilers {
				updaterRS := updater(ctx, objAdapter.GetObject(), osResource)
				reconcileStatus = updaterRS.WithReconcileStatus(reconcileStatus)
			}
		}
	}

	return reconcileStatus
}

func (c *Controller[
	orcObjectPT, orcObjectT,
	resourceSpecT,
	filterT,
	objectApplyPT,
	statusApplyPT, statusApplyT,
	osResourceT,
]) reconcileDelete(ctx context.Context, objAdapter interfaces.APIObjectAdapter[orcObjectPT, resourceSpecT, filterT]) (reconcileStatus progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)
	log.V(logging.Verbose).Info("Reconciling OpenStack resource delete")

	var osResource *osResourceT

	deleted := false
	defer func() {
		// No point updating status after removing the finalizer
		if !deleted {
			reconcileStatus = reconcileStatus.WithReconcileStatus(
				status.UpdateStatus(ctx, c, c.statusWriter, objAdapter.GetObject(), osResource, reconcileStatus))
		}
	}()

	actuator, reconcileStatus := c.helperFactory.NewDeleteActuator(ctx, objAdapter.GetObject(), c)
	if needsReschedule, err := reconcileStatus.NeedsReschedule(); needsReschedule {
		if err == nil {
			log.V(logging.Verbose).Info("Waiting on events before deletion")
		}
		return reconcileStatus
	}

	deleted, osResource, reconcileStatus = DeleteResource(ctx, log, c, objAdapter, actuator)
	if needsReschedule, err := reconcileStatus.NeedsReschedule(); needsReschedule && err == nil {
		log.V(logging.Verbose).Info("Waiting on events before deletion")
	}
	return reconcileStatus
}
