package providers

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath                         = "lbaas"
	resourcePath                     = "providers"
	flavorCapabilitiesPath           = "flavor_capabilities"
	availabilityZoneCapabilitiesPath = "availability_zone_capabilities"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func flavorCapabilitiesURL(c *gophercloud.ServiceClient, provider string) string {
	return c.ServiceURL(rootPath, resourcePath, provider, flavorCapabilitiesPath)
}

func availabilityZoneCapabilitiesURL(c *gophercloud.ServiceClient, provider string) string {
	return c.ServiceURL(rootPath, resourcePath, provider, availabilityZoneCapabilitiesPath)
}
