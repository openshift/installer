/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	api "github.com/Azure/azure-service-operator/v2/api/authorization/v1api20220401"
	storage "github.com/Azure/azure-service-operator/v2/api/authorization/v1api20220401/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/reflecthelpers"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
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

	// If this cast doesn't compile, update the `api` import to reference the now latest
	// stable version of the authorization group (this will happen when we import a new
	// API version in the generator.)
	if assignment, ok := rsrc.(*api.RoleAssignment); ok {
		// Check to see whether this role assignment is inherited or not
		// (we can tell by looking at the scope of the assignment)
		if assignment.Status.Scope != nil && owner != nil {
			if !strings.EqualFold(owner.ARMID, *assignment.Status.Scope) {
				// Scope isn't our owner, so it's inherited from further up and should not be imported
				return extensions.ImportSkipped("role assignment is inherited"), nil
			}
		}
	}

	return result, nil
}

var _ extensions.ARMResourceModifier = &RoleAssignmentExtension{}

func (extension *RoleAssignmentExtension) ModifyARMResource(
	ctx context.Context,
	armClient *genericarmclient.GenericClient,
	armObj genruntime.ARMResource,
	obj genruntime.ARMMetaObject,
	kubeClient kubeclient.Client,
	resolver *resolver.Resolver,
	log logr.Logger,
) (genruntime.ARMResource, error) {
	ra, ok := obj.(*storage.RoleAssignment)
	if !ok {
		return nil, eris.Errorf(
			"Cannot run RoleAssignmentExtension.ModifyARMResource() with unexpected resource type %T",
			obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not been updated to match
	var _ conversion.Hub = ra

	// If the specified role definition uses a well known name, look it up
	roleDefinitionName := ra.Spec.RoleDefinitionReference.WellKnownName
	if roleDefinitionName != "" {
		roleDefinitionId, err := resolveBuiltInRoleDefinition(ctx, roleDefinitionName, armObj, armClient)
		if err != nil {
			return nil, eris.Wrapf(err, "resolving built in role definition %q", roleDefinitionName)
		}

		log.V(1).Info("Resolved built-in role", "roleName", roleDefinitionName, "roleId", roleDefinitionId)

		spec := armObj.Spec()
		err = reflecthelpers.SetProperty(spec, "Properties.RoleDefinitionId", &roleDefinitionId)
		if err != nil {
			return nil, eris.Wrapf(err, "error setting RoleDefinitionId to %s", roleDefinitionId)
		}

		return armObj, nil
	}

	return armObj, nil
}

// resolveBuiltInRoleDefinition looks up the ARM ID for a built-in role definition
// given its well-known name. If it cannot be resolved, the original name is returned.
// roleDefinitionName is the well-known name of the role definition to resolve.
// armObj is the ARM resource for which the role definition is being resolved
// (required because we need the subscription ID).
// armClient is the ARM client to use for any necessary API calls to load the predefined role definitions.
func resolveBuiltInRoleDefinition(
	ctx context.Context,
	roleDefinitionName string,
	armObj genruntime.ARMResource,
	armClient *genericarmclient.GenericClient,
) (string, error) {
	// Load the built-in role definitions
	defs, err := builtInRoleDefinitions(ctx, armClient)
	if err != nil {
		return "", eris.Wrapf(err, "loading built-in role definitions to resolve %q", roleDefinitionName)
	}

	// Look up the role ID by well-known name (case insensitive)
	roleId, ok := defs[strings.ToLower(roleDefinitionName)]
	if !ok {
		// If we can't resolve, return an error
		return "", eris.Errorf(
			"unable to find Azure built-in role-definition with well-known name %q (see %s)",
			roleDefinitionName,
			"https://learn.microsoft.com/azure/role-based-access-control/built-in-roles")
	}

	// We need the subscription ID from the resource to construct the ARM ID for a well known role
	roleARMId, err := arm.ParseResourceID(armObj.GetID())
	if err != nil {
		return "", eris.Wrapf(err, "failed to parse the ARM ResourceId for %s", armObj.GetID())
	}

	armID := fmt.Sprintf(
		"/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s",
		roleARMId.SubscriptionID,
		roleId)

	return armID, nil
}

// builtInRoleDefinitionsCacheSize is the initial size of the built-in role definitions cache
const builtInRoleDefinitionsCacheSize = 800

var (
	// builtInRoleDefinitionsCache is a cache of all known built-in role definitions, keyed by
	// lower case role name. At the time of writing there are ~700 such roles, so we
	// pre-initialize the map to sufficient capacity to accommodate those and some modest growth.
	builtInRoleDefinitionsCache map[string]string = make(map[string]string, builtInRoleDefinitionsCacheSize)

	// builtInRoleDefinitionsLock protects access to builtInRoleDefinitionsCache
	builtInRoleDefinitionsLock sync.RWMutex

	// builtInRoleDefinitionsCachingDisabled is set to true to disable caching of built-in role definitions for testing.
	builtInRoleDefinitionsCachingDisabled bool = false
)

// DisableBuiltInRoleDefinitionsCaching disables caching of built-in role definitions.
// This is intended for use in tests only.
func DisableBuiltInRoleDefinitionsCaching() {
	builtInRoleDefinitionsLock.Lock()
	defer builtInRoleDefinitionsLock.Unlock()

	builtInRoleDefinitionsCachingDisabled = true
	clear(builtInRoleDefinitionsCache)
}

// builtInRoleDefinitions returns a map of built-in role definitions, keyed by lower case role name.
// armClient is the ARM client to use for any necessary API calls to load the predefined role definitions.
func builtInRoleDefinitions(
	ctx context.Context,
	armClient *genericarmclient.GenericClient,
) (map[string]string, error) {
	builtInRoleDefinitionsLock.Lock()
	defer builtInRoleDefinitionsLock.Unlock()

	// Short circuit if already loaded
	if len(builtInRoleDefinitionsCache) > 0 {
		return builtInRoleDefinitionsCache, nil
	}

	result := make(map[string]string, builtInRoleDefinitionsCacheSize)

	cl, err := armauthorization.NewRoleDefinitionsClient(armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, eris.Wrap(err, "creating client to load built-in role definitions")
	}

	pager := cl.NewListPager("/", nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			clear(builtInRoleDefinitionsCache) // ensure we'll try again next time
			return nil, eris.Wrap(err, "loading built-in role definitions")
		}

		for _, role := range page.Value {
			if *role.Properties.RoleType == "BuiltInRole" {
				result[strings.ToLower(*role.Properties.RoleName)] = *role.Name
			}
		}
	}

	// Cache the loaded definitions unless caching is disabled for testing
	// (we acquire the lock at the start of the function, so this is safe)
	if !builtInRoleDefinitionsCachingDisabled {
		// Cache the loaded definitions
		builtInRoleDefinitionsCache = result
	}

	return result, nil
}
