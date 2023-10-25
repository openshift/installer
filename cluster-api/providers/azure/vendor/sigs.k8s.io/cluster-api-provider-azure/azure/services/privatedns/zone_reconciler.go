/*
Copyright 2022 The Kubernetes Authors.

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

package privatedns

import (
	"context"

	"github.com/pkg/errors"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

func (s *Service) reconcileZone(ctx context.Context, zoneSpec azure.ResourceSpecGetter) (managed bool, err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "privatedns.Service.reconcileZone")
	defer done()

	managed, err = s.IsManaged(ctx)
	if err != nil {
		if azure.ResourceNotFound(err) {
			managed = true
		} else {
			return managed, err
		}
	}

	if !managed {
		log.V(1).Info("Skipping reconciliation of unmanaged private DNS zone", "private DNS", zoneSpec.ResourceName())
		return managed, nil
	}

	_, err = s.zoneReconciler.CreateOrUpdateResource(ctx, zoneSpec, serviceName)
	return managed, err
}

func (s *Service) deleteZone(ctx context.Context, zoneSpec azure.ResourceSpecGetter) (managed bool, err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "privatedns.Service.deleteZone")
	defer done()

	// Skip deleting the private DNS zone when it's not managed by capz.
	isManaged, err := s.IsManaged(ctx)
	if err != nil {
		if azure.ResourceNotFound(err) {
			// already deleted or doesn't exist, cleanup status and return.
			s.Scope.DeleteLongRunningOperationState(zoneSpec.ResourceName(), serviceName, infrav1.DeleteFuture)
			return managed, nil
		}
		return managed, errors.Wrapf(err, "could not get private DNS zone state of %s in resource group %s", zoneSpec.ResourceName(), zoneSpec.ResourceGroupName())
	}

	if !isManaged {
		log.V(1).Info("Skipping deletion of unmanaged private DNS zone", "private DNS", zoneSpec.ResourceName())
		return managed, nil
	}

	// if we reach here, it means that this vnet link is managed by capz.
	managed = true

	// Delete the private DNS zone, which also deletes all records
	err = s.zoneReconciler.DeleteResource(ctx, zoneSpec, serviceName)
	return managed, err
}
