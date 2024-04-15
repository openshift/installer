/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	api "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1api20221201"
	hub "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1api20221201/storage"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.PreReconciliationChecker = &FlexibleServersDatabaseExtension{}

func (extension *FlexibleServersDatabaseExtension) PreReconcileCheck(
	_ context.Context,
	_ genruntime.MetaObject,
	owner genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
	_ extensions.PreReconcileCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// Check to see if our owning server is ready for the database to be reconciled
	if owner == nil {
		// TODO: query from ARM instead?
		return extensions.ProceedWithReconcile(), nil
	}
	if server, ok := owner.(*hub.FlexibleServer); ok {
		serverState := server.Status.State
		if serverState != nil && flexibleServerStateBlocksReconciliation(*serverState) {
			return extensions.BlockReconcile(
				fmt.Sprintf(
					"Owning FlexibleServer is in provisioning state %q",
					*serverState)), nil
		}
	}

	return extensions.ProceedWithReconcile(), nil
}

var _ extensions.Importer = &FlexibleServersDatabaseExtension{}

// Import skips databases that can't be managed by ARM
func (extension *FlexibleServersDatabaseExtension) Import(
	ctx context.Context,
	rsrc genruntime.ImportableResource,
	owner *genruntime.ResourceReference,
	next extensions.ImporterFunc,
) (extensions.ImportResult, error) {
	if server, ok := rsrc.(*api.FlexibleServersDatabase); ok {
		if server.Spec.AzureName == "azure_maintenance" {
			return extensions.ImportSkipped("azure_maintenance database is not accessible by users"), nil
		}

		if server.Spec.AzureName == "azure_sys" {
			return extensions.ImportSkipped("built in databases cannot be managed by ARM"), nil
		}

		if server.Spec.AzureName == "postgres" {
			return extensions.ImportSkipped("built in databases cannot be managed by ARM"), nil
		}
	}

	return next(ctx, rsrc, owner)
}
