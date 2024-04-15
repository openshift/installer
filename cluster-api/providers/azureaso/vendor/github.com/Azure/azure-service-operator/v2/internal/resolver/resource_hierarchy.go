/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package resolver

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

type ResourceHierarchyRoot string

const (
	ResourceHierarchyRootResourceGroup = ResourceHierarchyRoot("ResourceGroup")
	ResourceHierarchyRootSubscription  = ResourceHierarchyRoot("Subscription")
	ResourceHierarchyRootTenant        = ResourceHierarchyRoot("Tenant")
	ResourceHierarchyRootARMID         = ResourceHierarchyRoot("ARMID")
	ResourceHierarchyRootOverride      = ResourceHierarchyRoot("Override")
)

// If we wanted to type-assert we'd have to solve some circular dependency problems... for now this is ok.
const ResourceGroupKind = "ResourceGroup"
const ResourceGroupGroup = "resources.azure.com"

type ResourceHierarchy []genruntime.ARMMetaObject

// ResourceGroup returns the resource group that the hierarchy is in, or an error if the hierarchy is not rooted
// in a resource group.
func (h ResourceHierarchy) ResourceGroup() (string, error) {
	rootKind := h.rootKind(h)
	if rootKind == ResourceHierarchyRootARMID {
		armIDStr := h[0].Owner().ARMID
		armID, err := arm.ParseResourceID(armIDStr)
		if err != nil {
			return "", err
		}

		if armID.ResourceGroupName == "" {
			return "", errors.Errorf("not rooted by a resource group: %s", armIDStr)
		}

		return armID.ResourceGroupName, nil
	}

	if rootKind != ResourceHierarchyRootResourceGroup {
		return "", errors.Errorf("not rooted by a resource group: %s", rootKind)
	}

	resourceGroup := h[0]
	return resourceGroup.GetName(), nil
}

// Location returns the location root of the hierarchy, or an error
// if the root is not a subscription.
func (h ResourceHierarchy) Location() (string, error) {
	rootKind := h.rootKind(h)
	// We don't support ARM ID rooted for this method because it doesn't really make sense.
	// If there's a need for it in the future we can add it
	if rootKind != ResourceHierarchyRootSubscription {
		return "", errors.Errorf("not rooted in a subscription: %s", rootKind)
	}

	// There's an assumption here that the root resource has a location
	locatable, ok := h[0].(genruntime.LocatableResource)
	if !ok {
		return "", errors.Errorf("root does not implement LocatableResource: %T", h[0])
	}

	return locatable.Location(), nil
}

// AzureName returns the Azure name for use in creating a resource.
func (h ResourceHierarchy) AzureName() string {
	azureNames := h.getAzureNames()

	if len(azureNames) == 0 {
		return ""
	}

	return azureNames[len(azureNames)-1]
}

// TODO: It's a bit awkward that this takes a subscriptionID parameter but does nothing with it in the tenant scope case
// FullyQualifiedARMID returns the fully qualified ARM ID of the resource
func (h ResourceHierarchy) FullyQualifiedARMID(subscriptionID string) (string, error) {
	return h.fullyQualifiedARMIDImpl(subscriptionID, h)
}

func (h ResourceHierarchy) fullyQualifiedARMIDImpl(subscriptionID string, originalHierarchy ResourceHierarchy) (string, error) {
	lastResource := h[len(h)-1]
	lastResourceScope := lastResource.GetResourceScope()

	if lastResourceScope == genruntime.ResourceScopeExtension {
		var parentARMID string
		if lastResource.Owner().IsDirectARMReference() {
			parentARMID = lastResource.Owner().ARMID
		} else {
			hierarchy := h[:len(h)-1]
			var err error
			parentARMID, err = hierarchy.fullyQualifiedARMIDImpl(subscriptionID, h)
			if err != nil {
				return "", err
			}
		}

		provider, types, err := genruntime.GetResourceTypeAndProvider(lastResource)
		if err != nil {
			return "", err
		}
		if len(types) != 1 {
			return "", errors.Errorf("extension resource cannot have more than one resource type, but had type: %s", lastResource.GetType())
		}

		return fmt.Sprintf("%s/providers/%s/%s/%s", parentARMID, provider, types[0], lastResource.AzureName()), nil
	}

	azureNames := h.getAzureNames()

	rootKind := h.rootKind(originalHierarchy)
	switch rootKind {
	case ResourceHierarchyRootSubscription:
		// TODO: This is currently a special case as the only resource like this is ResourceGroup and ResourceGroup itself
		// TODO: is a bit funky because it doesn't have a /providers like everything else does...
		return genericarmclient.MakeResourceGroupID(subscriptionID, azureNames[0]), nil
	case ResourceHierarchyRootResourceGroup:
		rgName := azureNames[0]
		remainingNames := azureNames[1:]
		// The only resource we actually care about for figuring out resource types is the
		// most derived resource
		res := h[len(h)-1]
		provider, resourceTypes, err := genruntime.GetResourceTypeAndProvider(res)
		if err != nil {
			return "", err
		}

		root := h[0]

		err = genruntime.VerifyResourceOwnerARMID(root)
		if err != nil {
			return "", err
		}

		// Safe to do it this way, Claimer makes sure the owner exists and is Ready and will always have an armId annotation before we reach here.
		ownerARMID, err := genruntime.GetAndParseResourceID(root)
		if err != nil {
			return "", err
		}

		// Confirm that the subscription ID the user specified matches the subscription ID we're using from our credential
		if ok := genruntime.CheckARMIDMatchesSubscription(subscriptionID, ownerARMID); !ok {
			return "", core.NewSubscriptionMismatchError(ownerARMID.SubscriptionID, subscriptionID)
		}

		// Ensure that we have the same number of names and types
		if len(remainingNames) != len(resourceTypes) {
			return "", errors.Errorf(
				"could not create fully qualified ARM ID, had %d azureNames and %d resourceTypes. azureNames: %+q resourceTypes: %+q",
				len(remainingNames),
				len(resourceTypes),
				remainingNames,
				resourceTypes)
		}

		// Join them together
		interleaved := genruntime.InterleaveStrSlice(resourceTypes, remainingNames)
		return genericarmclient.MakeResourceGroupScopeARMID(subscriptionID, rgName, provider, interleaved...)
	case ResourceHierarchyRootTenant:
		// The only resource we actually care about for figuring out resource types is the
		// most derived resource
		res := h[len(h)-1]
		provider, resourceTypes, err := genruntime.GetResourceTypeAndProvider(res)
		if err != nil {
			return "", err
		}

		// Ensure that we have the same number of names and types
		if len(azureNames) != len(resourceTypes) {
			return "", errors.Errorf(
				"could not create fully qualified ARM ID, had %d azureNames and %d resourceTypes. azureNames: %+q resourceTypes: %+q",
				len(azureNames),
				len(resourceTypes),
				azureNames,
				resourceTypes)
		}
		// Join them together
		interleaved := genruntime.InterleaveStrSlice(resourceTypes, azureNames)
		return genericarmclient.MakeTenantScopeARMID(provider, interleaved...)
	case ResourceHierarchyRootARMID:
		// TODO: Possibly refactor this huge method into sub-functions?
		// The only resource we actually care about for figuring out resource types is the
		// most derived resource
		res := h[len(h)-1]
		provider, resourceTypes, err := genruntime.GetResourceTypeAndProvider(res)
		if err != nil {
			return "", err
		}

		// We also need the ARMID from the root resource
		root := h[0]
		armIDStr := root.Owner().ARMID // Safe to do this without nil-guards because we already checked elsewhere

		// Trim the trailing slash of the ARM ID if it's there (we'll add it back later)
		armIDStr = strings.TrimRight(armIDStr, "/")

		err = genruntime.VerifyResourceOwnerARMID(root)
		if err != nil {
			return "", err
		}

		ownerARMID, err := arm.ParseResourceID(armIDStr)
		if err != nil {
			return "", err
		}

		// Confirm that the subscription ID the user specified matches the subscription ID we're using from our credential
		if ok := genruntime.CheckARMIDMatchesSubscription(subscriptionID, ownerARMID); !ok {
			return "", core.NewSubscriptionMismatchError(ownerARMID.SubscriptionID, subscriptionID)
		}

		// Rooting to an ARM ID means that some of the resourceTypes may not actually be included explicitly in our
		// hierarchy (because they're instead in the ARM ID itself). We filter these out of resourceTypes by
		// removing types that aren't included in the hierarchy.
		_, rootResourceTypes, err := genruntime.GetResourceTypeAndProvider(root)
		if err != nil {
			return "", err
		}
		resourceTypesIncludedInARMID := rootResourceTypes[:len(rootResourceTypes)-1]
		resourceTypes = resourceTypes[len(resourceTypesIncludedInARMID):]

		// Ensure that we have the same number of names and types
		if len(azureNames) != len(resourceTypes) {
			return "", errors.Errorf("could not create fully qualified ARM ID, had %d azureNames and %d resourceTypes. azureNames: %+q resourceTypes: %+q",
				len(azureNames),
				len(resourceTypes),
				azureNames,
				resourceTypes)
		}

		interleaved := genruntime.InterleaveStrSlice(resourceTypes, azureNames)
		suffix := strings.Join(interleaved, "/")
		// If the root ARM ID already contains the provider ID, we can just append the pairs.
		// If it doesn't, we need to build a full ARM ID by appending the provider as well.
		if strings.Contains(strings.ToLower(armIDStr), strings.ToLower(provider)) {
			return fmt.Sprintf("%s/%s", armIDStr, suffix), nil
		} else {
			return fmt.Sprintf("%s/providers/%s/%s", armIDStr, provider, suffix), nil
		}
	case ResourceHierarchyRootOverride:
		// Find the resource that has the override and start building the ID from there:
		idFragment, idx := h.getChildResourceIDOverride()
		if idx == -1 {
			return "", errors.Errorf("resource had root kind %q, but had no child resource ID override", rootKind)
		}
		if idx != 0 {
			return "", errors.Errorf("resource had root kind %q, but child resource override was not at index 0. Instead at index %d", rootKind, idx)
		}
		return idFragment, nil
		// TODO: if we actually need to support this for hierarchies, we could do something like the below
		// Get the type of the resource idx+1
		// Get the type of the resource len(h)-1
		// diff the two
		// compute the remaining names and types
		// append them together and tack them on to the idFragment

	default:
		return "", errors.Errorf("unknown root kind %q", rootKind)
	}
}

// rootKind returns the ResourceHierarchyRoot type of the hierarchy.
// There are 6 cases here:
//  1. The hierarchy is comprised solely of a resource group. This is subscription rooted.
//  2. The hierarchy has multiple entries and roots up to a resource group. This is Resource Group rooted.
//  3. The hierarchy has multiple entries and doesn't root up to a resource group. This is subscription rooted.
//  4. The hierarchy roots up to a tenant scope resource. This is tenant rooted.
//  5. The hierarchy roots up to a resource whose Owner() is an ARMID. This is ARMID rooted.
//  6. The hierarchy contains a resource that sets genruntime.ChildResourceIDOverrideAnnotation. This is
//     "Override" rooted.
func (h ResourceHierarchy) rootKind(originalHierarchy ResourceHierarchy) ResourceHierarchyRoot {
	if len(h) == 0 {
		panic("resource hierarchy cannot be len 0")
	}

	// This is a special kind of root for if genruntime.ChildResourceIDOverrideAnnotation is used
	_, idx := originalHierarchy.getChildResourceIDOverride()
	// Child ID override doesn't apply if it's set on the most derived resource
	if idx != -1 && idx != len(originalHierarchy)-1 {
		return ResourceHierarchyRootOverride
	}

	root := h[0]

	// Check if the root resource is owned by an ARM ID
	rootOwner := root.Owner()
	if rootOwner != nil && rootOwner.IsDirectARMReference() {
		return ResourceHierarchyRootARMID
	}

	scope := root.GetResourceScope()
	if scope == genruntime.ResourceScopeTenant {
		return ResourceHierarchyRootTenant
	}

	if scope == genruntime.ResourceScopeLocation {
		if len(h) == 1 { // Just the location scope resource
			return ResourceHierarchyRootSubscription
		}
		return ResourceHierarchyRootResourceGroup
	}

	return ResourceHierarchyRootSubscription
}

func (h ResourceHierarchy) getAzureNames() []string {
	azureNames := make([]string, 0, len(h))

	for _, res := range h {
		azureNames = append(azureNames, res.AzureName())
	}

	return azureNames
}

// getChildResourceIDOverride returns the child resource ID override and the index at which the override was specified, or
// -1 if there was no childResourceIDOverride
func (h ResourceHierarchy) getChildResourceIDOverride() (string, int) {
	for i, res := range h {
		idFragment, ok := genruntime.GetChildResourceIDOverride(res)
		if ok {
			return idFragment, i
		}
	}

	return "", -1
}
