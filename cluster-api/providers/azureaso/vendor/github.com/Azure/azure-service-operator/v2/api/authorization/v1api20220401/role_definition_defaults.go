/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package v1api20220401

import (
	"strings"

	"github.com/Azure/azure-service-operator/v2/internal/util/randextensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

var _ genruntime.Defaulter = &RoleDefinition{}

func (definition *RoleDefinition) CustomDefault() {
	definition.defaultAzureName()
}

// defaultAzureName performs special AzureName defaulting for RoleDefinition by generating a GUID
// based on the Role name.
// The GUID generation technique is picked by the user using OperatorSpec.NamingConvention property if they want
// the GUID to be a stable or random string. Stable being the default.
// Stable GUID algorithm:
// We generate stable the UUID using UUIDv5 with a seed string based on the group+kind of the RoleDefinition and the
// namespace+name it's deployed into.
// We include the namespace and name to ensure no two RoleDefinitions in the same cluster can end up
// with the same UUID.
// We include the group and kind to ensure that different kinds of resources get different UUIDs. This isn't
// entirely required by Azure, but it makes sense to avoid collisions between two resources of different types
// even if they have the same namespace and name.
// In the rare case users have multiple ASO instances with resources in the same namespace in each cluster
// having the same name but not actually pointing to the same Azure resource (maybe in a different subscription?)
// they can avoid name conflicts by explicitly specifying AzureName for their RoleDefinition.
func (definition *RoleDefinition) defaultAzureName() {
	// If owner is not set we can't default AzureName, but the request will be rejected anyway for lack of owner.
	if definition.Spec.Owner == nil {
		return
	}

	if definition.AzureName() == "" {
		if definition.Spec.OperatorSpec != nil &&
			strings.EqualFold(*definition.Spec.OperatorSpec.NamingConvention, "random") {
			definition.Spec.AzureName = randextensions.MakeRandomUUID()
		} else if definition.Spec.OperatorSpec == nil ||
			definition.Spec.OperatorSpec != nil && strings.EqualFold(*definition.Spec.OperatorSpec.NamingConvention, "stable") {
			gk := definition.GroupVersionKind().GroupKind()
			definition.Spec.AzureName = randextensions.MakeUUIDName(
				definition.Name,
				randextensions.MakeUniqueOwnerScopedString(
					definition.Owner(),
					gk,
					definition.Namespace,
					definition.Name))
		}
	}
}
