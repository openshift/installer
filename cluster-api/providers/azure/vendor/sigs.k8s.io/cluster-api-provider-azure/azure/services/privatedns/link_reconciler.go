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

func (s *Service) reconcileLinks(ctx context.Context, links []azure.ResourceSpecGetter) (managed bool, err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "privatedns.Service.reconcileLinks")
	defer done()

	var resErr error

	// We go through the list of links to reconcile each one, independently of the result of the previous one.
	// If multiple errors occur, we return the most pressing one.
	// Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (i.e. error creating) -> operationNotDoneError (i.e. creating in progress) -> no error (i.e. created)
	for _, linkSpec := range links {
		isLinkManaged, err := s.isVnetLinkManaged(ctx, linkSpec)
		if err != nil {
			if azure.ResourceNotFound(err) {
				isLinkManaged = true
			} else {
				return managed, err
			}
		}

		if !isLinkManaged {
			log.V(2).Info("Skipping vnet link reconciliation for unmanaged vnet link", "vnet link",
				linkSpec.ResourceName(), "private dns zone", linkSpec.OwnerResourceName())
			continue
		}

		// we consider VnetLinks as managed if at least of the links is managed.
		managed = true
		if _, err := s.vnetLinkReconciler.CreateOrUpdateResource(ctx, linkSpec, serviceName); err != nil {
			if !azure.IsOperationNotDoneError(err) || resErr == nil {
				resErr = err
			}
		}
	}

	return managed, resErr
}

func (s *Service) deleteLinks(ctx context.Context, links []azure.ResourceSpecGetter) (managed bool, err error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "privatedns.Service.deleteLinks")
	defer done()

	var resErr error

	// We go through the list of links to delete each one, independently of the result of the previous one.
	// If multiple errors occur, we return the most pressing one.
	// Order of precedence (highest -> lowest) is: error that is not an operationNotDoneError (i.e. error creating) -> operationNotDoneError (i.e. creating in progress) -> no error (i.e. created)
	for _, linkSpec := range links {
		// If the virtual network link is not managed by capz, skip its reconciliation
		isVnetLinkManaged, err := s.isVnetLinkManaged(ctx, linkSpec)
		if err != nil {
			if azure.ResourceNotFound(err) {
				// already deleted or doesn't exist, cleanup status and return.
				s.Scope.DeleteLongRunningOperationState(linkSpec.ResourceName(), serviceName, infrav1.DeleteFuture)
				continue
			}
			return managed, errors.Wrapf(err, "could not get vnet link state of %s in resource group %s",
				linkSpec.OwnerResourceName(), linkSpec.ResourceGroupName())
		}

		if !isVnetLinkManaged {
			log.V(2).Info("Skipping vnet link deletion for unmanaged vnet link", "vnet link",
				linkSpec.ResourceName(), "private dns zone", linkSpec.OwnerResourceName())
			continue
		}

		// if we reach here, it means that this vnet link is managed by capz.
		managed = true

		if err := s.vnetLinkReconciler.DeleteResource(ctx, linkSpec, serviceName); err != nil {
			if !azure.IsOperationNotDoneError(err) || resErr == nil {
				resErr = err
			}
		}
	}

	return managed, resErr
}
