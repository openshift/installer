// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package config

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

// These are hardcoded because the init function that initializes them in azcore isn't in /cloud it's in /arm which
// we don't import.
const (
	DefaultEndpoint         = "https://management.azure.com"
	DefaultAudience         = "https://management.core.windows.net/"
	DefaultAADAuthorityHost = "https://login.microsoftonline.com/"
)

type Configuration struct {
	AzureAuthorityHost      string
	ResourceManagerEndpoint string
	ResourceManagerAudience string
}

// Cloud returns the cloud the configuration is using
func (v Configuration) Cloud() cloud.Configuration {
	// Special handling if we've got all the defaults just return the official public cloud
	// configuration
	hasDefaultAzureAuthorityHost := v.AzureAuthorityHost == "" || v.AzureAuthorityHost == DefaultAADAuthorityHost
	hasDefaultResourceManagerEndpoint := v.ResourceManagerEndpoint == "" || v.ResourceManagerEndpoint == DefaultEndpoint
	hasDefaultResourceManagerAudience := v.ResourceManagerAudience == "" || v.ResourceManagerAudience == DefaultAudience

	if hasDefaultResourceManagerEndpoint && hasDefaultResourceManagerAudience && hasDefaultAzureAuthorityHost {
		return cloud.AzurePublic
	}

	// We default here too to more easily support empty Values objects
	azureAuthorityHost := v.AzureAuthorityHost
	resourceManagerEndpoint := v.ResourceManagerEndpoint
	resourceManagerAudience := v.ResourceManagerAudience
	if azureAuthorityHost == "" {
		azureAuthorityHost = DefaultAADAuthorityHost
	}
	if resourceManagerAudience == "" {
		resourceManagerAudience = DefaultAudience
	}
	if resourceManagerEndpoint == "" {
		resourceManagerEndpoint = DefaultEndpoint
	}

	return cloud.Configuration{
		ActiveDirectoryAuthorityHost: azureAuthorityHost,
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Endpoint: resourceManagerEndpoint,
				Audience: resourceManagerAudience,
			},
		},
	}
}
