/*
Copyright The ORC Authors.

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

package roleassignment

import (
	"context"
	"fmt"
	"iter"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/reconciler"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/resync"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/status"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	"github.com/k-orc/openstack-resource-controller/v2/internal/scope"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/finalizers"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
)

const (
	// The time to wait before reconciling again when we are waiting for some change in OpenStack
	externalUpdatePollingPeriod = 15 * time.Second
)

// roleassignmentReconciler reconciles RoleAssignment objects.
// Unlike other ORC resources, role assignments are relationships (not resources with IDs),
// so this uses a custom reconciler instead of the generic framework.
type roleassignmentReconciler struct {
	client              client.Client
	scopeFactory        scope.Factory
	defaultResyncPeriod time.Duration

	statusWriter roleassignmentStatusWriter
}

func (r *roleassignmentReconciler) GetName() string                { return controllerName }
func (r *roleassignmentReconciler) GetK8sClient() client.Client    { return r.client }
func (r *roleassignmentReconciler) GetScopeFactory() scope.Factory { return r.scopeFactory }

// Reconcile is the main entry point for reconciliation.
// It fetches the RoleAssignment object and routes to either reconcileNormal or reconcileDelete.
func (r *roleassignmentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	orcObject := new(orcObjectT)
	err := r.client.Get(ctx, req.NamespacedName, orcObject)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Object deleted, nothing to do
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	log := ctrl.LoggerFrom(ctx)

	// Check if object is being deleted
	if !orcObject.GetDeletionTimestamp().IsZero() {
		return r.reconcileDelete(ctx, orcObject).Return(log)
	}

	return r.reconcileNormal(ctx, orcObject).Return(log)
}

func hasRoleAssignmentComponents(statusResource *orcv1alpha1.RoleAssignmentResourceStatus) bool {
	return statusResource != nil &&
		statusResource.RoleID != "" &&
		(statusResource.UserID != "" || statusResource.GroupID != "") &&
		(statusResource.ProjectID != "" || statusResource.DomainID != "")
}

// reconcileNormal handles the normal reconciliation flow:
// 1. Check if we should reconcile (based on Progressing condition)
// 2. Create actuator (OpenStack client)
// 3. Get or create the role assignment
// 4. Update status
func (r *roleassignmentReconciler) reconcileNormal(ctx context.Context, orcObject orcObjectPT) (reconcileStatus progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)
	effectiveResyncPeriod := resync.DetermineResyncPeriod(orcObject.Spec.ResyncPeriod, r.defaultResyncPeriod)

	// Check if we should skip reconciliation
	if !reconciler.ShouldReconcile(orcObject, orcObject.Status.LastSyncTime, effectiveResyncPeriod) {
		log.V(logging.Verbose).Info("Status is up to date: not reconciling")
		if remaining := resync.RemainingUntilNextSync(orcObject.Status.LastSyncTime, effectiveResyncPeriod); remaining > 0 {
			return reconcileStatus.WithRequeue(remaining)
		}
		return reconcileStatus
	}

	log.V(logging.Verbose).Info("Reconciling role assignment")

	var osResource *osResourceT

	// Ensure we always update status at the end
	defer func() {
		reconcileStatus = reconcileStatus.WithReconcileStatus(
			status.UpdateStatus(ctx, r, r.statusWriter, orcObject, osResource, reconcileStatus))
	}()

	// Phase 3: Add finalizer if not present
	if !controllerutil.ContainsFinalizer(orcObject, finalizer) {
		patch := finalizers.SetFinalizerPatch(orcObject, finalizer)
		if err := r.client.Patch(ctx, orcObject, patch, client.ForceOwnership, orcstrings.GetSSAFieldOwnerWithTxn(controllerName, orcstrings.SSATransactionFinalizer)); err != nil {
			return progress.WrapError(fmt.Errorf("setting finalizer: %w", err))
		}
	}

	// Phase 3: Create actuator
	actuator, actuatorRS := r.newActuator(ctx, orcObject)
	if needsReschedule, err := actuatorRS.NeedsReschedule(); needsReschedule {
		if err == nil {
			log.V(logging.Verbose).Info("Waiting on events before creation")
		}
		return actuatorRS.WithReconcileStatus(reconcileStatus)
	}

	// Phase 4: Check if role assignment exists using Status.Resource components
	if orcObject.Status.Resource != nil {
		statusResource := orcObject.Status.Resource
		// If we have all components in status, try to fetch the role assignment
		if hasRoleAssignmentComponents(statusResource) {
			osResource, getRS := actuator.GetResourceByComponents(
				ctx,
				statusResource.RoleID,
				statusResource.UserID,
				statusResource.GroupID,
				statusResource.ProjectID,
				statusResource.DomainID,
			)
			if needsReschedule, _ := getRS.NeedsReschedule(); needsReschedule {
				return getRS.WithReconcileStatus(reconcileStatus)
			}

			if osResource != nil {
				log.V(logging.Verbose).Info("Got existing role assignment")
			} else {
				// Status was fully populated but the resource no longer exists in
				// OpenStack. GetResourceByComponents uses a LIST query which returns
				// (nil, nil) for empty results rather than a 404 error, so we detect
				// deletion here.
				if orcObject.Spec.ManagementPolicy == orcv1alpha1.ManagementPolicyUnmanaged {
					return progress.WrapError(
						orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "role assignment has been deleted from OpenStack"))
				}
				log.V(logging.Info).Info("Role assignment was deleted externally; will recreate")
			}
		}
	}

	// Phase 5: Import by filter
	if osResource == nil {
		if importSpec := orcObject.Spec.Import; importSpec != nil {
			if filter := importSpec.Filter; filter != nil {
				resourceIter, importRS := actuator.ListOSResourcesForImport(ctx, orcObject, *filter)
				if needsReschedule, _ := importRS.NeedsReschedule(); needsReschedule {
					return importRS.WithReconcileStatus(reconcileStatus)
				}

				var err error
				osResource, err = atMostOne(resourceIter,
					orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration,
						"found more than one matching OpenStack resource during import"))
				if err != nil {
					return progress.WrapError(err)
				}

				if osResource == nil {
					return progress.WaitingOnOpenStack(progress.WaitingOnCreation, externalUpdatePollingPeriod)
				}

				log.V(logging.Info).Info("Imported role assignment")
			}
		}
	}

	// Phase 6: Adoption - check for existing resource before creating
	if osResource == nil {
		if orcObject.Spec.ManagementPolicy == orcv1alpha1.ManagementPolicyUnmanaged {
			// We never create an unmanaged resource
			// API validation should have ensured that one of the above functions returned
			return progress.WrapError(
				orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Not creating unmanaged resource"))
		}

		if resourceIter, canAdopt := actuator.ListOSResourcesForAdoption(ctx, orcObject); canAdopt {
			var err error
			osResource, err = atMostOne(resourceIter,
				orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration,
					"found more than one matching OpenStack resource during adoption"))
			if err != nil {
				return progress.WrapError(err)
			}
			if osResource != nil {
				log.V(logging.Info).Info("Adopted previously created resource")
			}
		}
	}

	// Phase 7: Fetch dependencies and create role assignment
	if osResource == nil {
		log.V(logging.Info).Info("Creating resource")
		var createRS progress.ReconcileStatus
		osResource, createRS = actuator.CreateResource(ctx, orcObject)
		if needsReschedule, err := createRS.NeedsReschedule(); needsReschedule {
			if err == nil {
				log.V(logging.Verbose).Info("Waiting on dependencies or creation")
			}
			return createRS.WithReconcileStatus(reconcileStatus)
		}

		if osResource == nil {
			return reconcileStatus.WithError(fmt.Errorf("osResource is not set, but no wait events or error"))
		}

		log.V(logging.Info).Info("Role assignment created")
	}

	if resync.ShouldScheduleResync(effectiveResyncPeriod, reconcileStatus) {
		reconcileStatus = reconcileStatus.WithRequeue(resync.CalculateJitteredDuration(effectiveResyncPeriod))
	}
	return reconcileStatus
}

// atMostOne returns the first element from the iterator, or nil if it's empty.
// It returns multipleErr if the iterator yields more than one element.
func atMostOne(resourceIter iter.Seq2[*osResourceT, error], multipleErr error) (*osResourceT, error) {
	next, stop := iter.Pull2(resourceIter)
	defer stop()

	// Try to fetch the first result
	osResource, err, ok := next()
	if err != nil {
		return nil, err
	} else if !ok {
		// No first result
		return nil, nil
	}

	// Check that there are no other results
	_, err, ok = next()
	if err != nil {
		return nil, err
	} else if ok {
		return nil, multipleErr
	}

	return osResource, nil
}

// reconcileDelete handles deletion of the RoleAssignment:
// 1. Check finalizer
// 2. Fetch the role assignment (using Status.Resource components)
// 3. Check management policy
// 4. Delete from OpenStack
// 5. Remove finalizer
func (r *roleassignmentReconciler) reconcileDelete(ctx context.Context, orcObject orcObjectPT) (reconcileStatus progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)
	log.V(logging.Verbose).Info("Reconciling role assignment delete")

	var osResource *osResourceT
	deleted := false

	// Update status unless we've removed the finalizer
	defer func() {
		if !deleted {
			reconcileStatus = reconcileStatus.WithReconcileStatus(
				status.UpdateStatus(ctx, r, r.statusWriter, orcObject, osResource, reconcileStatus))
		}
	}()

	// Check if our finalizer is present
	var foundFinalizer bool
	for _, f := range orcObject.GetFinalizers() {
		if f == finalizer {
			foundFinalizer = true
		} else {
			reconcileStatus = reconcileStatus.WaitingOnFinalizer(f)
		}
	}

	// Cleanup not required if our finalizer is not present
	if !foundFinalizer {
		return reconcileStatus
	}

	if needsReschedule, err := reconcileStatus.NeedsReschedule(); needsReschedule {
		if err == nil {
			log.V(logging.Verbose).Info("Deferring resource cleanup due to remaining external finalizers")
		}
		return reconcileStatus
	}

	removeFinalizer := func(reconcileStatus progress.ReconcileStatus) progress.ReconcileStatus {
		if err := r.client.Patch(ctx, orcObject, finalizers.RemoveFinalizerPatch(orcObject), orcstrings.GetSSAFieldOwnerWithTxn(controllerName, orcstrings.SSATransactionFinalizer)); err != nil {
			return reconcileStatus.WithError(fmt.Errorf("removing finalizer: %w", err))
		}
		deleted = true
		return reconcileStatus
	}

	// Check management policy
	managementPolicy := orcObject.Spec.ManagementPolicy
	managedOptions := orcObject.Spec.ManagedOptions
	if managementPolicy == orcv1alpha1.ManagementPolicyUnmanaged || managedOptions.GetOnDelete() == orcv1alpha1.OnDeleteDetach {
		logPolicy := []any{"managementPolicy", managementPolicy}
		if managementPolicy == orcv1alpha1.ManagementPolicyManaged {
			logPolicy = append(logPolicy, "onDelete", managedOptions.GetOnDelete())
		}
		log.V(logging.Verbose).Info("Not deleting OpenStack resource due to policy", logPolicy...)
		return removeFinalizer(reconcileStatus)
	}

	// Create actuator for OpenStack operations
	actuator, actuatorRS := r.newActuator(ctx, orcObject)
	if needsReschedule, err := actuatorRS.NeedsReschedule(); needsReschedule {
		if err == nil {
			log.V(logging.Verbose).Info("Waiting on events before deletion")
		}
		return actuatorRS.WithReconcileStatus(reconcileStatus)
	}

	// Fetch the role assignment using Status.Resource components
	if orcObject.Status.Resource != nil {
		statusResource := orcObject.Status.Resource
		if statusResource.RoleID != "" &&
			(statusResource.UserID != "" || statusResource.GroupID != "") &&
			(statusResource.ProjectID != "" || statusResource.DomainID != "") {

			var getRS progress.ReconcileStatus
			osResource, getRS = actuator.GetResourceByComponents(
				ctx,
				statusResource.RoleID,
				statusResource.UserID,
				statusResource.GroupID,
				statusResource.ProjectID,
				statusResource.DomainID,
			)
			if needsReschedule, err := getRS.NeedsReschedule(); needsReschedule {
				// NotFound is our success condition for delete
				if err == nil || !orcerrors.IsNotFound(err) {
					return getRS.WithReconcileStatus(reconcileStatus)
				}
				osResource = nil
			}
		}
	}

	// If status was never populated, check for orphaned resources via adoption
	if osResource == nil && orcObject.Status.Resource == nil {
		resourceIter, canAdopt := actuator.ListOSResourcesForAdoption(ctx, orcObject)
		if canAdopt {
			var err error
			osResource, err = atMostOne(resourceIter,
				orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration,
					"found more than one matching OpenStack resource during adoption"))
			if err != nil {
				return reconcileStatus.WithError(err)
			}
		}
	}

	if osResource == nil {
		log.V(logging.Info).Info("Role assignment deletion confirmed")
		return removeFinalizer(reconcileStatus)
	}

	log.V(logging.Info).Info("Deleting role assignment from OpenStack")
	deleteRS := actuator.DeleteResource(ctx, orcObject, osResource)
	if needsReschedule, _ := deleteRS.NeedsReschedule(); needsReschedule {
		return deleteRS.WithReconcileStatus(reconcileStatus)
	}

	log.V(logging.Info).Info("Role assignment deletion confirmed")
	return removeFinalizer(reconcileStatus)
}

// newActuator creates a roleassignmentActuator with OpenStack client setup.
func (r *roleassignmentReconciler) newActuator(ctx context.Context, orcObject orcObjectPT) (roleassignmentActuator, progress.ReconcileStatus) {
	log := ctrl.LoggerFrom(ctx)

	// Ensure credential secrets exist and have our finalizer
	_, reconcileStatus := credentialsDependency.GetDependencies(ctx, r.client, orcObject, func(*corev1.Secret) bool { return true })
	if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
		return roleassignmentActuator{}, reconcileStatus
	}

	clientScope, err := r.scopeFactory.NewClientScopeFromObject(ctx, r.client, log, orcObject)
	if err != nil {
		return roleassignmentActuator{}, progress.WrapError(err)
	}
	osClient, err := clientScope.NewRoleAssignmentClient()
	if err != nil {
		return roleassignmentActuator{}, progress.WrapError(err)
	}

	return roleassignmentActuator{
		osClient:  osClient,
		k8sClient: r.client,
	}, nil
}
