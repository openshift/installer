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

package routerinterface

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/progress"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/generic/status"
	"github.com/k-orc/openstack-resource-controller/v2/internal/controllers/port"
	osclients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/util/applyconfigs"
	orcstrings "github.com/k-orc/openstack-resource-controller/v2/internal/util/strings"
	orcapplyconfigv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/pkg/clients/applyconfiguration/api/v1alpha1"
)

func getStatusSummary(osResource *osclients.PortExt) (metav1.ConditionStatus, progress.ReconcileStatus) {
	if osResource == nil {
		return metav1.ConditionFalse, nil
	}

	if osResource.Status == port.PortStatusActive {
		return metav1.ConditionTrue, nil
	}

	// port exists but is not ACTIVE
	return metav1.ConditionFalse, progress.WaitingOnOpenStack(progress.WaitingOnReady, portStatusPollingPeriod)
}

// createStatusUpdate computes a complete status update based on the given
// observed state. This is separated from updateStatus to facilitate unit
// testing, as the version of k8s we currently import does not support patch
// apply in the fake client.
// Needs: https://github.com/kubernetes/kubernetes/pull/125560
func createStatusUpdate(orcObject *orcv1alpha1.RouterInterface, port *osclients.PortExt, reconcileStatus progress.ReconcileStatus, now metav1.Time) (*orcapplyconfigv1alpha1.RouterInterfaceApplyConfiguration, progress.ReconcileStatus) {
	applyConfigStatus := orcapplyconfigv1alpha1.RouterInterfaceStatus()
	applyConfig := orcapplyconfigv1alpha1.RouterInterface(orcObject.Name, orcObject.Namespace).WithStatus(applyConfigStatus)

	// Note that unlike other resources we don't rely on this value to be immutable, so it's not in a separate transaction.
	if port != nil {
		applyConfigStatus.WithID(port.ID)
	}

	isAvailable, statusReconcileStatus := getStatusSummary(port)
	reconcileStatus = reconcileStatus.WithReconcileStatus(statusReconcileStatus)
	status.SetCommonConditions(orcObject, applyConfigStatus, isAvailable, reconcileStatus, now)

	return applyConfig, reconcileStatus
}

// updateStatus computes a complete status based on the given observed state and writes it to status.
func (r *orcRouterInterfaceReconciler) updateStatus(ctx context.Context, orcObject *orcv1alpha1.RouterInterface, port *osclients.PortExt, reconcileStatus progress.ReconcileStatus) progress.ReconcileStatus {
	now := metav1.NewTime(time.Now())

	statusUpdate, reconcileStatus := createStatusUpdate(orcObject, port, reconcileStatus, now)
	return reconcileStatus.WithError(
		r.client.Status().Patch(ctx, orcObject, applyconfigs.Patch(types.ApplyPatchType, statusUpdate), client.ForceOwnership, orcstrings.GetSSAFieldOwnerWithTxn(controllerName, orcstrings.SSATransactionFinalizer)))
}
