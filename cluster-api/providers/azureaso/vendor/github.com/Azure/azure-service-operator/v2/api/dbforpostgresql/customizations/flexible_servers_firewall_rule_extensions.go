/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	postgresql "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v20250801/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PreReconciliationOwnerChecker = &FlexibleServersFirewallRuleExtension{}

func (ext *FlexibleServersFirewallRuleExtension) PreReconcileOwnerCheck(
	ctx context.Context,
	owner genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
	next extensions.PreReconcileOwnerCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// Check to see if our owning server is ready for the database to be reconciled
	if server, ok := owner.(*postgresql.FlexibleServer); ok {
		serverState := server.Status.State
		if serverState != nil && flexibleServerStateBlocksReconciliation(*serverState) {
			return extensions.BlockReconcile(
				fmt.Sprintf(
					"Owning FlexibleServer is in state %q",
					*serverState)), nil
		}
	}

	return next(ctx, owner, resourceResolver, armClient, log)
}
