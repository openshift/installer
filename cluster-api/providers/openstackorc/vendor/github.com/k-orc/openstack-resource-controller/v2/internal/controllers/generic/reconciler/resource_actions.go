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
	"iter"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/logging"
	orcerrors "github.com/k-orc/openstack-resource-controller/v2/internal/util/errors"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/finalizers"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
)

const (
	// The time to wait between checking if a delete was successful
	deletePollingPeriod = 1 * time.Second

	// The time to wait before reconciling again when we are waiting for some change in OpenStack
	externalUpdatePollingPeriod = 15 * time.Second
)

func GetOrCreateOSResource[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	}, orcObjectT any,
	resourceSpecT any, filterT any,
	osResourceT any,
](
	ctx context.Context, log logr.Logger, controller ResourceController,
	objAdapter interfaces.APIObjectAdapter[orcObjectPT, resourceSpecT, filterT],
	actuator interfaces.CreateResourceActuator[orcObjectPT, orcObjectT, filterT, osResourceT],
) (*osResourceT, progress.ReconcileStatus) {
	k8sClient := controller.GetK8sClient()

	finalizer := orcstrings.GetFinalizerName(controller.GetName())
	if !controllerutil.ContainsFinalizer(objAdapter.GetObject(), finalizer) {
		patch := finalizers.SetFinalizerPatch(objAdapter.GetObject(), finalizer)
		if err := k8sClient.Patch(ctx, objAdapter.GetObject(), patch, client.ForceOwnership, orcstrings.GetSSAFieldOwnerWithTxn(controller.GetName(), orcstrings.SSATransactionFinalizer)); err != nil {
			return nil, progress.WrapError(fmt.Errorf("setting finalizer: %w", err))
		}
	}

	if resourceID := objAdapter.GetStatusID(); resourceID != nil {
		osResource, reconcileStatus := actuator.GetOSResourceByID(ctx, *resourceID)
		if needsReschedule, err := reconcileStatus.NeedsReschedule(); needsReschedule {
			if orcerrors.IsNotFound(err) {
				// An OpenStack resource we previously referenced has been deleted unexpectedly. We can't recover from this.
				return osResource, progress.WrapError(
					orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "resource has been deleted from OpenStack"))
			} else {
				return osResource, reconcileStatus
			}
		}
		if osResource != nil {
			log.V(logging.Verbose).Info("Got existing OpenStack resource", "ID", actuator.GetResourceID(osResource))
		}
		return osResource, nil
	}

	// Import by ID
	if resourceID := objAdapter.GetImportID(); resourceID != nil {
		osResource, reconcileStatus := actuator.GetOSResourceByID(ctx, *resourceID)
		if needsReschedule, err := reconcileStatus.NeedsReschedule(); needsReschedule {
			if orcerrors.IsNotFound(err) {
				// We assume that a resource imported by ID must already exist. It's a terminal error if it doesn't.
				return osResource, progress.WrapError(
					orcerrors.Terminal(orcv1alpha1.ConditionReasonUnrecoverableError, "referenced resource does not exist in OpenStack"))
			} else {
				return osResource, reconcileStatus
			}
		}
		if osResource != nil {
			log.V(logging.Verbose).Info("Imported existing OpenStack resource by ID", "ID", actuator.GetResourceID(osResource))
		}
		return osResource, nil
	}

	// Import by filter
	if filter := objAdapter.GetImportFilter(); filter != nil {
		resourceIter, reconcileStatus := actuator.ListOSResourcesForImport(ctx, objAdapter.GetObject(), *filter)
		if needsReschedule, _ := reconcileStatus.NeedsReschedule(); needsReschedule {
			return nil, reconcileStatus
		}

		osResource, err := atMostOne(resourceIter, orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "found more than one matching OpenStack resource during import"))
		if err != nil {
			return nil, progress.WrapError(err)
		}

		if osResource == nil {
			return nil, progress.WaitingOnOpenStack(progress.WaitingOnCreation, externalUpdatePollingPeriod)
		}
		return osResource, reconcileStatus
	}

	// Create
	if objAdapter.GetManagementPolicy() == orcv1alpha1.ManagementPolicyUnmanaged {
		// We never create an unmanaged resource
		// API validation should have ensured that one of the above functions returned
		return nil, progress.WrapError(
			orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "Not creating unmanaged resource"))
	}

	osResource, err := getResourceForAdoption(ctx, actuator, objAdapter)
	if err != nil {
		return nil, progress.WrapError(err)
	}
	if osResource != nil {
		log.V(logging.Info).Info("Adopted previously created resource")
		return osResource, nil
	}

	log.V(logging.Info).Info("Creating resource")
	return actuator.CreateResource(ctx, objAdapter.GetObject())
}

func DeleteResource[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	}, orcObjectT any,
	resourceSpecT any, filterT any,
	osResourceT any,
](
	ctx context.Context, log logr.Logger, controller ResourceController,
	objAdapter interfaces.APIObjectAdapter[orcObjectPT, resourceSpecT, filterT],
	actuator interfaces.DeleteResourceActuator[orcObjectPT, orcObjectT, osResourceT],
) (bool, *osResourceT, progress.ReconcileStatus) {
	var osResource *osResourceT

	// We always fetch the resource by ID so we can continue to report status even when waiting for a finalizer
	statusID := objAdapter.GetStatusID()
	if statusID != nil {
		var getOSResourceRS progress.ReconcileStatus
		osResource, getOSResourceRS = actuator.GetOSResourceByID(ctx, *statusID)
		if needsReschedule, err := getOSResourceRS.NeedsReschedule(); needsReschedule {
			// If there's no error it means GetOSResourceByID is waiting on something
			// NotFound is our success condition, handled below
			if err == nil || !orcerrors.IsNotFound(err) {
				return false, osResource, getOSResourceRS
			}

			// Gophercloud can return an empty non-nil object when returning errors,
			// which will confuse us below.
			osResource = nil
		}
	}

	finalizer := orcstrings.GetFinalizerName(controller.GetName())

	var reconcileStatus progress.ReconcileStatus
	var foundFinalizer bool
	for _, f := range objAdapter.GetFinalizers() {
		if f == finalizer {
			foundFinalizer = true
		} else {
			reconcileStatus = reconcileStatus.WaitingOnFinalizer(f)
		}
	}

	// Cleanup not required if our finalizer is not present
	if !foundFinalizer {
		return true, osResource, reconcileStatus
	}

	if needsReschedule, err := reconcileStatus.NeedsReschedule(); needsReschedule {
		if err == nil {
			log.V(logging.Verbose).Info("Deferring resource cleanup due to remaining external finalizers")
		}
		return false, osResource, reconcileStatus
	}

	removeFinalizer := func(reconcileStatus progress.ReconcileStatus) progress.ReconcileStatus {
		if err := controller.GetK8sClient().Patch(ctx, objAdapter.GetObject(), finalizers.RemoveFinalizerPatch(objAdapter.GetObject()), orcstrings.GetSSAFieldOwnerWithTxn(controller.GetName(), orcstrings.SSATransactionFinalizer)); err != nil {
			return reconcileStatus.WithError(fmt.Errorf("removing finalizer: %w", err))
		}
		return reconcileStatus
	}

	// We won't delete the resource for an unmanaged object, or if onDelete is detach
	managementPolicy := objAdapter.GetManagementPolicy()
	managedOptions := objAdapter.GetManagedOptions()
	if managementPolicy == orcv1alpha1.ManagementPolicyUnmanaged || managedOptions.GetOnDelete() == orcv1alpha1.OnDeleteDetach {
		logPolicy := []any{"managementPolicy", managementPolicy}
		if managementPolicy == orcv1alpha1.ManagementPolicyManaged {
			logPolicy = append(logPolicy, "onDelete", managedOptions.GetOnDelete())
		}
		log.V(logging.Verbose).Info("Not deleting OpenStack resource due to policy", logPolicy...)
		return true, osResource, removeFinalizer(reconcileStatus)
	}

	// If status.ID was not set, we still need to check if there's an orphaned object.
	if osResource == nil && statusID == nil {
		var err error
		osResource, err = getResourceForAdoption(ctx, actuator, objAdapter)
		if err != nil {
			return false, osResource, reconcileStatus.WithError(err)
		}
	}

	if osResource == nil {
		log.V(logging.Info).Info("Resource deletion confirmed")

		return true, osResource, removeFinalizer(reconcileStatus)
	}

	log.V(logging.Info).Info("Deleting OpenStack resource")
	deleteRS := actuator.DeleteResource(ctx, objAdapter.GetObject(), osResource)
	if needsReschedule, _ := deleteRS.NeedsReschedule(); needsReschedule {
		return false, osResource, deleteRS.WithReconcileStatus(reconcileStatus)
	}

	// We still need to poll for the deletion of the OpenStack resource
	return false, osResource, reconcileStatus.WaitingOnOpenStack(progress.WaitingOnDeletion, deletePollingPeriod)
}

func atMostOne[osResourceT any](resourceIter iter.Seq2[*osResourceT, error], multipleErr error) (*osResourceT, error) {
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

func getResourceForAdoption[
	orcObjectPT interface {
		*orcObjectT
		client.Object
		orcv1alpha1.ObjectWithConditions
	}, orcObjectT any,
	resourceSpecT any, filterT any,
	osResourceT any,
](
	ctx context.Context,
	actuator interfaces.BaseResourceActuator[orcObjectPT, orcObjectT, osResourceT],
	objAdapter interfaces.APIObjectAdapter[orcObjectPT, resourceSpecT, filterT],
) (*osResourceT, error) {
	resourceIter, canAdopt := actuator.ListOSResourcesForAdoption(ctx, objAdapter.GetObject())
	if !canAdopt {
		return nil, nil
	}

	return atMostOne(resourceIter, orcerrors.Terminal(orcv1alpha1.ConditionReasonInvalidConfiguration, "found more than one matching OpenStack resource during adoption"))
}
