/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

// AzureNameFromConfigProvider is implemented by resource Spec types that support resolving their
// AzureName from a ConfigMap at reconciliation time.
type AzureNameFromConfigProvider interface {
	// GetAzureNameFromConfig returns the ConfigMapReference to resolve the Azure name from, or nil if not set.
	GetAzureNameFromConfig() *ConfigMapReference
}
