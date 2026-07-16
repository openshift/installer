/*
Copyright 2025 The ORC Authors.

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

package status

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/interfaces"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/applyconfigs"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
)

func SetStatusID[
	orcObjectPT interface {
		client.Object
		orcv1alpha1.ObjectWithConditions
	},
	objectApplyPT interfaces.ORCApplyConfig[objectApplyPT, statusApplyPT],
	statusApplyPT interface {
		*statusApplyT
		interfaces.ORCStatusApplyConfigWithID[statusApplyPT]
	},
	statusApplyT any,
	osResourcePT any,
](
	ctx context.Context,
	controller interfaces.ResourceController,
	orcObject orcObjectPT,
	resourceID string,
	statusWriter interfaces.ResourceStatusWriter[orcObjectPT, osResourcePT, objectApplyPT, statusApplyPT],
) error {
	var status statusApplyPT = new(statusApplyT)
	status.WithID(resourceID)

	applyConfig := statusWriter.GetApplyConfig(orcObject.GetName(), orcObject.GetNamespace()).
		WithUID(orcObject.GetUID()).
		WithStatus(status)

	return controller.GetK8sClient().Status().Patch(ctx, orcObject, applyconfigs.Patch(types.MergePatchType, applyConfig))
}

// ClearStatusID clears the status.id field of an ORC object using a JSON merge
// patch. This is necessary when an externally deleted managed resource is
// detected: clearing the ID allows the next reconciliation to enter the
// standard creation path and assign a new ID after the resource is recreated.
//
// A JSON merge patch with an explicit null value is required because the
// generated apply configuration types use omitempty on the ID field, meaning a
// nil pointer would simply omit the field rather than clear it.
func ClearStatusID(ctx context.Context, controller interfaces.ResourceController, orcObject client.Object) error {
	patch := client.RawPatch(types.MergePatchType, []byte(`{"status":{"id":null}}`))
	return controller.GetK8sClient().Status().Patch(ctx, orcObject, patch)
}

// shouldSetLastSyncTime reports whether lastSyncTime should be set on a status
// update. It returns true only when the reconciliation completed successfully:
// the reconcileStatus contains neither errors nor progress messages. A requeue
// alone (e.g., for a periodic resync) does not prevent the update.
func shouldSetLastSyncTime(reconcileStatus progress.ReconcileStatus) bool {
	needsReschedule, _ := reconcileStatus.NeedsReschedule()
	return !needsReschedule
}

func UpdateStatus[
	orcObjectPT interface {
		client.Object
		orcv1alpha1.ObjectWithConditions
	},
	osResourcePT *osResourceT,
	objectApplyPT interfaces.ORCApplyConfig[objectApplyPT, statusApplyPT],
	statusApplyPT interface {
		interfaces.ORCStatusApplyConfigWithLastSyncTime[statusApplyPT]
		*statusApply
	},
	statusApply any,
	osResourceT any,
](
	ctx context.Context,
	controller interfaces.ResourceController,
	statusWriter interfaces.ResourceStatusWriter[orcObjectPT, osResourcePT, objectApplyPT, statusApplyPT],
	orcObject orcObjectPT, osResource osResourcePT,
	reconcileStatus progress.ReconcileStatus,
) progress.ReconcileStatus {
	log := ctrl.LoggerFrom(ctx)
	now := metav1.NewTime(time.Now())

	// Create a new apply configuration for this status transaction
	var applyConfigStatus statusApplyPT = new(statusApply)
	applyConfig := statusWriter.GetApplyConfig(orcObject.GetName(), orcObject.GetNamespace()).
		WithStatus(applyConfigStatus)

	// Write resource status to the apply configuration
	if osResource != nil {
		statusWriter.ApplyResourceStatus(log, osResource, applyConfigStatus)
	}

	// Set common conditions
	available, availableReconcileStatus := statusWriter.ResourceAvailableStatus(orcObject, osResource)
	reconcileStatus = reconcileStatus.WithReconcileStatus(availableReconcileStatus)
	SetCommonConditions(orcObject, applyConfigStatus, available, reconcileStatus, now)

	// Set lastSyncTime only on successful reconciliation: no errors and no
	// progress messages indicate that the controller successfully fetched the
	// resource state from OpenStack.
	if shouldSetLastSyncTime(reconcileStatus) {
		applyConfigStatus.WithLastSyncTime(now)
	}

	// Patch orcObject with the status transaction
	k8sClient := controller.GetK8sClient()
	ssaFieldOwner := orcstrings.GetSSAFieldOwnerWithTxn(controller.GetName(), orcstrings.SSATransactionStatus)
	return reconcileStatus.
		WithError(k8sClient.Status().Patch(ctx, orcObject, applyconfigs.Patch(types.ApplyPatchType, applyConfig), client.ForceOwnership, ssaFieldOwner))
}
