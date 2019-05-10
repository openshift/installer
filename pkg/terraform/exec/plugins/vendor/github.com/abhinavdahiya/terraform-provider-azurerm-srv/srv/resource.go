package srv

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// resourceID represents a parsed long-form Azure Resource Manager ID
// with the Subscription ID, Resource Group and the Provider as top-
// level fields, and other key-value pairs available via a map in the
// Path field.
type resourceID struct {
	SubscriptionID string
	ResourceGroup  string
	Provider       string
	Path           map[string]string
}

// parseAzureResourceID converts a long-form Azure Resource Manager ID
// into a ResourceID. We make assumptions about the structure of URLs,
// which is obviously not good, but the best thing available given the
// SDK.
func parseAzureResourceID(id string) (*resourceID, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure Id: %s", err)
	}

	path := idURL.Path

	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	components := strings.Split(path, "/")

	// We should have an even number of key-value pairs.
	if len(components)%2 != 0 {
		return nil, fmt.Errorf("The number of path segments is not divisible by 2 in %q", path)
	}

	var subscriptionID string

	// Put the constituent key-value pairs into a map
	componentMap := make(map[string]string, len(components)/2)
	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}

		// Catch the subscriptionID before it can be overwritten by another "subscriptions"
		// value in the ID which is the case for the Service Bus subscription resource
		if key == "subscriptions" && subscriptionID == "" {
			subscriptionID = value
		} else {
			componentMap[key] = value
		}
	}

	// Build up a resourceID from the map
	idObj := &resourceID{}
	idObj.Path = componentMap

	if subscriptionID != "" {
		idObj.SubscriptionID = subscriptionID
	} else {
		return nil, fmt.Errorf("No subscription ID found in: %q", path)
	}

	if resourceGroup, ok := componentMap["resourceGroups"]; ok {
		idObj.ResourceGroup = resourceGroup
		delete(componentMap, "resourceGroups")
	} else {
		// Some Azure APIs are weird and provide things in lower case...
		// However it's not clear whether the casing of other elements in the URI
		// matter, so we explicitly look for that case here.
		if resourceGroup, ok := componentMap["resourcegroups"]; ok {
			idObj.ResourceGroup = resourceGroup
			delete(componentMap, "resourcegroups")
		} else {
			return nil, fmt.Errorf("No resource group name found in: %q", path)
		}
	}

	// It is OK not to have a provider in the case of a resource group
	if provider, ok := componentMap["providers"]; ok {
		idObj.Provider = provider
		delete(componentMap, "providers")
	}

	return idObj, nil
}

func composeAzureResourceID(idObj *resourceID) (id string, err error) {
	if idObj.SubscriptionID == "" || idObj.ResourceGroup == "" {
		return "", fmt.Errorf("SubscriptionID and ResourceGroup cannot be empty")
	}

	id = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", idObj.SubscriptionID, idObj.ResourceGroup)

	if idObj.Provider != "" {
		if len(idObj.Path) < 1 {
			return "", fmt.Errorf("ResourceID.Path should have at least one item when ResourceID.Provider is specified")
		}

		id += fmt.Sprintf("/providers/%s", idObj.Provider)

		// sort the path keys so our output is deterministic
		var pathKeys []string
		for k := range idObj.Path {
			pathKeys = append(pathKeys, k)
		}
		sort.Strings(pathKeys)

		for _, k := range pathKeys {
			v := idObj.Path[k]
			if k == "" || v == "" {
				return "", fmt.Errorf("resourceID.Path cannot contain empty strings")
			}
			id += fmt.Sprintf("/%s/%s", k, v)
		}
	}

	return
}
