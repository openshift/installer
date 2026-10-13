/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genericarmclient

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/pkg/errors"
)

// NOTE: We don't use arm.ResourceID here because it's only supposed to be created via arm.ParseResourceID,
// which assumes you have an ID first.

// MakeTenantScopeARMID makes an ARM ID at the tenant scope. This has the format:
// /providers/<provider>/<resourceType>/<resourceName>/...
func MakeTenantScopeARMID(provider string, params ...string) (string, error) {
	if len(params) == 0 {
		return "", errors.New("At least 2 params must be specified")
	}
	if len(params)%2 != 0 {
		return "", errors.New("ARM Id params must come in resourceKind/name pairs")
	}

	suffix := strings.Join(params, "/")

	return fmt.Sprintf("/providers/%s/%s", provider, suffix), nil
}

// MakeSubscriptionScopeARMID makes an ARM ID at the subscription scope. This has the format:
// /subscriptions/00000000-0000-0000-0000-000000000000/providers/<provider>/<resourceType>/<resourceName>/...
func MakeSubscriptionScopeARMID(subscription string, provider string, params ...string) (string, error) {
	if len(params) == 0 {
		return "", errors.New("At least 2 params must be specified")
	}

	if len(params)%2 != 0 {
		return "", errors.New("ARM Id params must come in resourceKind/name pairs")
	}

	suffix := strings.Join(params, "/")

	return fmt.Sprintf("/subscriptions/%s/providers/%s/%s", subscription, provider, suffix), nil
}

// MakeResourceGroupScopeARMID makes an ARM ID at the resource group scope. This has the format:
// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/<rgName>/providers/<provider>/<resourceType>/<resourceName>/...
func MakeResourceGroupScopeARMID(subscription string, resourceGroup string, provider string, params ...string) (string, error) {
	if len(params) == 0 {
		return "", errors.New("At least 2 params must be specified")
	}
	if len(params)%2 != 0 {
		return "", errors.New("ARM Id params must come in resourceKind/name pairs")
	}

	suffix := strings.Join(params, "/")
	return fmt.Sprintf("%s/providers/%s/%s", MakeResourceGroupID(subscription, resourceGroup), provider, suffix), nil
}

// MakeResourceGroupID makes an ARM ID representing a resource group. This has the format:
// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/<rgName>
// This is "special" because there is no provider at all - but that is what creating/getting a resourceGroup expects.
func MakeResourceGroupID(subscription string, resourceGroup string) string {
	return fmt.Sprintf("%s/resourceGroups/%s", MakeSubscriptionID(subscription), resourceGroup)
}

// MakeSubscriptionID makes an ARM ID representing a subscription. This has the format:
// /subscriptions/00000000-0000-0000-0000-000000000000
// This is "special" because there is no provider at all
func MakeSubscriptionID(subscription string) string {
	return fmt.Sprintf("/subscriptions/%s", subscription)
}

// GetSubscription uses resource ID to extract and return subscription ID out of it.
func GetSubscription(path string) (string, error) {
	id, err := arm.ParseResourceID(path)
	if err != nil {
		return "", err
	}
	return id.SubscriptionID, nil
}
