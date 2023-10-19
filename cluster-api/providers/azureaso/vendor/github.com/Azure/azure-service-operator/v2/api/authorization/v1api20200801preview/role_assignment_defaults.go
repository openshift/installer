/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package v1api20200801preview

import (
	"github.com/Azure/azure-service-operator/v2/internal/util/randextensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

var _ genruntime.Defaulter = &RoleAssignment{}

func (assignment *RoleAssignment) CustomDefault() {
	assignment.defaultAzureName()
}

// defaultAzureName performs special AzureName defaulting for RoleAssignment by generating a stable GUID
// based on the assignment name.
// We generate the UUID using UUIDv5 with a seed string based on the group, kind, namespace and name.
// We include the namespace and name to ensure no two RoleAssignments in the same cluster can end up
// with the same UUID.
// We include the group and kind to ensure that different kinds of resources get different UUIDs. This isn't
// entirely required by Azure, but it makes sense to avoid collisions between two resources of different types
// even if they have the same namespace and name.
// We include the owner group, kind, and name to avoid collisions between resources with the same name in different
// clusters that actually point to different Azure resources.
// In the rare case users have multiple ASO instances with resources in the same namespace in each cluster
// having the same name but not actually pointing to the same Azure resource (maybe in a different subscription?)
// they can avoid name conflicts by explicitly specifying AzureName for their RoleAssignment.
// See https://learn.microsoft.com/en-us/azure/role-based-access-control/role-assignments#name for details about
// RoleAssignment name restrictions. Of note is that the UUID must be unique in the AAD tenant (subscription
// uniqueness isn't sufficient). A UUID conflict will result in the following error:
//
//	(RoleAssignmentUpdateNotPermitted) Tenant ID, application ID, principal ID, and scope are not allowed to be updated.
//	Code: RoleAssignmentUpdateNotPermitted
//	Message: Tenant ID, application ID, principal ID, and scope are not allowed to be updated.
func (assignment *RoleAssignment) defaultAzureName() {
	// If owner is not set we can't default AzureName, but the request will be rejected anyway for lack of owner.
	if assignment.Spec.Owner == nil {
		return
	}

	if assignment.AzureName() == "" {
		ownerGK := assignment.Owner().GroupKind()
		gk := assignment.GroupVersionKind().GroupKind()
		assignment.Spec.AzureName = randextensions.MakeUUIDName(
			ownerGK,
			assignment.Spec.Owner.Name,
			gk,
			assignment.Namespace,
			assignment.Name)
	}
}
