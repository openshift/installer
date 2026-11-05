// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package annotations

const PerResourceSecret = "serviceoperator.azure.com/credential-from"

// AzureNameFromConfig is used to store the resolved Azure name from a ConfigMap for resources
// which support resolving their Azure name from a ConfigMap.
const AzureNameFromConfig = "serviceoperator.azure.com/azure-name-from-config"
