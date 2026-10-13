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
		interfaces.ORCStatusApplyConfig[statusApplyPT]
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

func UpdateStatus[
	orcObjectPT interface {
		client.Object
		orcv1alpha1.ObjectWithConditions
	},
	osResourcePT *osResourceT,
	objectApplyPT interfaces.ORCApplyConfig[objectApplyPT, statusApplyPT],
	statusApplyPT interface {
		interfaces.ORCStatusApplyConfig[statusApplyPT]
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

	// Patch orcObject with the status transaction
	k8sClient := controller.GetK8sClient()
	ssaFieldOwner := orcstrings.GetSSAFieldOwnerWithTxn(controller.GetName(), orcstrings.SSATransactionStatus)
	return reconcileStatus.
		WithError(k8sClient.Status().Patch(ctx, orcObject, applyconfigs.Patch(types.ApplyPatchType, applyConfig), client.ForceOwnership, ssaFieldOwner))
}
