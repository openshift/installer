// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customizations

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	network "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701/storage"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PostReconciliationChecker = &PrivateEndpointExtension{}

func (extension *PrivateEndpointExtension) PostReconcileCheck(
	_ context.Context,
	obj genruntime.MetaObject,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
	_ extensions.PostReconcileCheckFunc) (extensions.PostReconcileCheckResult, error) {

	endpoint, ok := obj.(*network.PrivateEndpoint)
	if !ok {
		return extensions.PostReconcileCheckResult{},
			errors.Errorf("cannot run on unknown resource type %T, expected *network.PrivateEndpoint", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = endpoint

	var reqApprovals []string
	// We want to check `ManualPrivateLinkServiceConnections` as these are the ones which are not auto-approved.
	if connections := endpoint.Status.ManualPrivateLinkServiceConnections; connections != nil {
		for _, connection := range connections {
			if *connection.PrivateLinkServiceConnectionState.Status != "Approved" {
				reqApprovals = append(reqApprovals, *connection.Id)
			}
		}
	}

	if len(reqApprovals) > 0 {
		// Returns 'conditions.NewReadyConditionImpactingError' error
		return extensions.PostReconcileCheckResultFailure(
			fmt.Sprintf(
				"Private connection(s) '%q' to the PrivateEndpoint requires approval",
				reqApprovals)), nil
	}

	return extensions.PostReconcileCheckResultSuccess(), nil
}
