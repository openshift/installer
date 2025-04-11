/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genericarmclient

// GenericResource is used as a response type for poller. This is a shortened version of
// https://github.com/Azure/azure-sdk-for-go/blob/5659d929cb5966c1296568cb33410d12e0ee06c6/sdk/resourcemanager/resources/armresources/zz_generated_models.go#L922
type GenericResource struct {
	// The kind of the resource.
	Kind *string `json:"kind,omitempty"`

	// READ-ONLY; Resource ID
	ID *string `json:"id,omitempty"`

	// READ-ONLY; Resource name
	Name *string `json:"name,omitempty"`

	// READ-ONLY; Resource type
	Type *string `json:"type,omitempty"`
}

type GenericDeleteResponse struct {
	// Empty, for extension later
}
