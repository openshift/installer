/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	api "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1api20221201"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.Importer = &FlexibleServersConfigurationExtension{}

// Import skips databases that can't be managed by ARM
func (extension *FlexibleServersConfigurationExtension) Import(
	ctx context.Context,
	rsrc genruntime.ImportableResource,
	owner *genruntime.ResourceReference,
	next extensions.ImporterFunc,
) (extensions.ImportResult, error) {
	result, err := next(ctx, rsrc, owner)
	if err != nil {
		return extensions.ImportResult{}, err
	}

	if config, ok := rsrc.(*api.FlexibleServersConfiguration); ok {
		// Skip system defaults
		if config.Spec.Source != nil && *config.Spec.Source == "system-default" {
			return extensions.ImportSkipped("system-defaults don't need to be imported"), nil
		}
	}

	return result, nil
}
