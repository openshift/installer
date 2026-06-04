package azure

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/azure"
)

// Service names for API version overrides.
const (
	ServiceStorage = "storage"
	ServiceCompute = "compute"
	ServiceDNS     = "dns"
	ServiceNetwork = "network"
)

// Environment variable names for API version overrides.
const (
	envAPIVersionStorage = "AZURE_API_VERSION_STORAGE"
	envAPIVersionCompute = "AZURE_API_VERSION_COMPUTE"
	envAPIVersionDNS     = "AZURE_API_VERSION_DNS"
	envAPIVersionNetwork = "AZURE_API_VERSION_NETWORK"
)

// Default API versions for Azure Stack Hub compatibility.
// These versions are known to work with Azure Stack Hub's older API surface.
const (
	// DefaultStackStorageAPIVersion is the API version for storage operations on Azure Stack Hub.
	DefaultStackStorageAPIVersion = "2019-06-01"

	// DefaultStackComputeAPIVersion is the API version for compute operations on Azure Stack Hub.
	DefaultStackComputeAPIVersion = "2020-06-01"

	// DefaultStackDNSAPIVersion is the API version for DNS operations on Azure Stack Hub.
	DefaultStackDNSAPIVersion = "2018-05-01"

	// DefaultStackNetworkAPIVersion is the API version for network operations on Azure Stack Hub.
	DefaultStackNetworkAPIVersion = "2019-11-01"
)

// ClientConfig holds configuration for Azure SDK clients including
// cloud configuration and per-service API version overrides.
type ClientConfig struct {
	// CloudConfig contains the ARM endpoint and authentication configuration.
	CloudConfig cloud.Configuration

	// CloudName identifies the Azure cloud environment.
	CloudName azure.CloudEnvironment

	// APIVersionOverrides allows per-service API version overrides.
	// Key: service name (ServiceStorage, ServiceCompute, ServiceDNS, ServiceNetwork)
	// Value: API version string (e.g., "2019-06-01")
	APIVersionOverrides map[string]string
}

// NewClientConfig creates a new ClientConfig with the given cloud configuration
// and cloud name. It automatically loads API version overrides from environment
// variables and applies Azure Stack Hub defaults when appropriate.
func NewClientConfig(cloudConfig cloud.Configuration, cloudName azure.CloudEnvironment) *ClientConfig {
	config := &ClientConfig{
		CloudConfig:         cloudConfig,
		CloudName:           cloudName,
		APIVersionOverrides: make(map[string]string),
	}

	// Apply Azure Stack Hub defaults first
	if cloudName == azure.StackCloud {
		config.APIVersionOverrides[ServiceStorage] = DefaultStackStorageAPIVersion
		config.APIVersionOverrides[ServiceCompute] = DefaultStackComputeAPIVersion
		config.APIVersionOverrides[ServiceDNS] = DefaultStackDNSAPIVersion
		config.APIVersionOverrides[ServiceNetwork] = DefaultStackNetworkAPIVersion
	}

	// Load overrides from environment variables (these take precedence)
	config.loadAPIVersionOverridesFromEnv()

	return config
}

// loadAPIVersionOverridesFromEnv loads API version overrides from environment variables.
// Environment variables take precedence over defaults (including Azure Stack Hub defaults).
func (c *ClientConfig) loadAPIVersionOverridesFromEnv() {
	envMappings := map[string]string{
		envAPIVersionStorage: ServiceStorage,
		envAPIVersionCompute: ServiceCompute,
		envAPIVersionDNS:     ServiceDNS,
		envAPIVersionNetwork: ServiceNetwork,
	}

	for envVar, service := range envMappings {
		if value := os.Getenv(envVar); value != "" {
			logrus.Debugf("Using API version override from %s: %s", envVar, value)
			c.APIVersionOverrides[service] = value
		}
	}
}

// GetAPIVersion returns the API version for the given service.
// Returns empty string if no override is configured (SDK default will be used).
func (c *ClientConfig) GetAPIVersion(service string) string {
	if c.APIVersionOverrides == nil {
		return ""
	}
	return c.APIVersionOverrides[service]
}

// ClientOptions returns arm.ClientOptions configured for the given service.
// If an API version override is configured for the service, it will be applied.
func (c *ClientConfig) ClientOptions(service string) *arm.ClientOptions {
	opts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.CloudConfig,
		},
	}

	if apiVersion := c.GetAPIVersion(service); apiVersion != "" {
		opts.ClientOptions.APIVersion = apiVersion
	}

	return opts
}

// DefaultClientOptions returns arm.ClientOptions without any API version override.
// Use this for services that should use the SDK default API version.
func (c *ClientConfig) DefaultClientOptions() *arm.ClientOptions {
	return &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: c.CloudConfig,
		},
	}
}

// IsStackCloud returns true if the cloud environment is Azure Stack Hub.
func (c *ClientConfig) IsStackCloud() bool {
	return c.CloudName == azure.StackCloud
}
