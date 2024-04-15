/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"context"
	"strings"

	api "github.com/Azure/azure-service-operator/v2/api/authorization/v1api20220401"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.Importer = &RoleAssignmentExtension{}

func (extension *RoleAssignmentExtension) Import(
	ctx context.Context,
	rsrc genruntime.ImportableResource,
	owner *genruntime.ResourceReference,
	next extensions.ImporterFunc,
) (extensions.ImportResult, error) {
	result, err := next(ctx, rsrc, owner)
	if err != nil {
		return extensions.ImportResult{}, err
	}

	if assignment, ok := rsrc.(*api.RoleAssignment); ok {
		// Check to see whether this role assignment is inherited or not
		// (we can tell by looking at the scope of the assignment)
		if assignment.Status.Scope != nil {
			if !strings.EqualFold(owner.ARMID, *assignment.Status.Scope) {
				// Scope isn't our owner, so it's inherited from further up and should not be imported
				return extensions.ImportSkipped("role assignment is inherited"), nil
			}
		}
	}

	return result, nil
}
