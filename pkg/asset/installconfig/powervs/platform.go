package powervs

import (
	"os"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/powervs"
)

// Platform collects powervs-specific configuration.
func Platform() (*powervs.Platform, error) {

	bxCli, err := NewBxClient(true)
	if err != nil {
		return nil, err
	}

	var p powervs.Platform

	// @TODO: The way we're using this (a precreated boot image in a Power VS Service instance) doesn't
	// align with the installer's definition of this. We need a new var here and in the install config.
	// This should be done before code cutoff in a followon PR.
	if osOverride := os.Getenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); len(osOverride) != 0 {
		p.ClusterOSImage = osOverride
	}

	p.Region = bxCli.Region
	p.Zone = bxCli.Zone
	p.UserID = bxCli.User.ID
	p.PowerVSResourceGroup = bxCli.PowerVSResourceGroup

	// When running in staging mode, automatically inject the staging service endpoint
	// overrides so that the generated install-config does not require manual edits.
	if IsStagingMode() {
		p.ServiceEndpoints = stagingServiceEndpoints(p.Region, p.Zone)
	}

	return &p, nil
}

// stagingServiceEndpoints builds the full list of service endpoint overrides
// needed for an IBM Cloud staging (test) environment, derived from the given
// PowerVS region and zone.
func stagingServiceEndpoints(region, zone string) []configv1.PowerVSServiceEndpoint {
	var endpoints []configv1.PowerVSServiceEndpoint

	add := func(name, url string) {
		endpoints = append(endpoints, configv1.PowerVSServiceEndpoint{Name: name, URL: url})
	}

	add(string(configv1.IBMCloudServiceIAM), GetIAMURL())
	add(string(configv1.IBMCloudServiceResourceController), GetResourceControllerURL())
	add(string(configv1.IBMCloudServiceResourceManager), GetResourceManagerURL())
	add("CIS", GetCISURL())
	add(string(configv1.IBMCloudServiceDNSServices), GetDNSServicesURL())
	add("Power", GetPowerURL(zone, region))
	add(string(configv1.IBMCloudServiceVPC), GetVPCURL(region))
	add("TransitGateway", GetTransitGatewayURL())
	add(string(configv1.IBMCloudServiceCOS), GetCOSEndpoint(region))

	return endpoints
}
