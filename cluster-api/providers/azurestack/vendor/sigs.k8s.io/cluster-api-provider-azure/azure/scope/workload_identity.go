/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scope

import (
	"os"
	"strings"
)

/*

For workload identity to work we need the following.

|-----------------------------------------------------------------------------------|
|AZURE_FEDERATED_TOKEN_FILE | The path of the projected service account token file. |
|-----------------------------------------------------------------------------------|

With the current implementation, AZURE_CLIENT_ID and AZURE_TENANT_ID are read via AzureClusterIdentity.

AZURE_FEDERATED_TOKEN_FILE is the path of the projected service account token which is by default
"/var/run/secrets/azure/tokens/azure-identity-token".
The path can be overridden by setting "AZURE_FEDERATED_TOKEN_FILE" env variable.

*/

const (
	// azureFederatedTokenFileEnvKey is the env key for AZURE_FEDERATED_TOKEN_FILE.
	azureFederatedTokenFileEnvKey = "AZURE_FEDERATED_TOKEN_FILE"
	// azureTokenFilePath is the path of the projected token.
	azureTokenFilePath = "/var/run/secrets/azure/tokens/azure-identity-token" // #nosec G101
)

// GetProjectedTokenPath return projected token file path from the env variable.
func GetProjectedTokenPath() string {
	tokenPath := strings.TrimSpace(os.Getenv(azureFederatedTokenFileEnvKey))
	if tokenPath == "" {
		return azureTokenFilePath
	}
	return tokenPath
}
