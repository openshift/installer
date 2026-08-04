package powervs

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// IsStagingMode returns true if the installer is running in staging/testing mode.
// This is controlled by the IBMCLOUD_STAGING environment variable.
// Valid values: "true", "1", "yes", "on" (case-insensitive).
func IsStagingMode() bool {
	value := strings.ToLower(os.Getenv("IBMCLOUD_STAGING"))
	isStaging := value == "true" || value == "1" || value == "yes" || value == "on"
	return isStaging
}

// GetIAMURL returns the appropriate IAM URL based on staging mode.
// This is checked at runtime to respect the --staging flag.
func GetIAMURL() string {
	url := "https://iam.cloud.ibm.com"
	if IsStagingMode() {
		url = "https://iam.test.cloud.ibm.com"
	}
	return url
}

// GetResourceControllerURL returns the appropriate Resource Controller URL based on staging mode.
// This is checked at runtime to respect the --staging flag.
func GetResourceControllerURL() string {
	if IsStagingMode() {
		return "https://resource-controller.test.cloud.ibm.com"
	}
	return "https://resource-controller.cloud.ibm.com"
}

// GetCISURL returns the appropriate CIS URL based on staging mode.
// This is checked at runtime to respect the --staging flag.
func GetCISURL() string {
	if IsStagingMode() {
		return "https://api.cis.test.cloud.ibm.com"
	}
	return "https://api.cis.cloud.ibm.com"
}

// GetDNSServicesURL returns the appropriate DNS Services URL based on staging mode.
// This is checked at runtime to respect the --staging flag.
// NOTE: The dnszonesv1 SDK does NOT automatically add /v1/ to the base URL,
// so we must include it in the URL we provide.
func GetDNSServicesURL() string {
	if IsStagingMode() {
		// DNS Services staging endpoint - must include /v1/ path
		return "https://api.dns-svcs.test.cloud.ibm.com/v1"
	}
	return "https://api.dns-svcs.cloud.ibm.com/v1"
}

// Default COSEndpoint for staging is "us-west".
const cosEndpointStagingRegionDefault = "us-west"

// GetCOSEndpoint returns the appropriate COS S3 endpoint based on staging mode and region.
// This returns the S3 data-plane endpoint for bucket/object operations.
// For the control-plane endpoint (service API), use GetCOSControlURL instead.
// Regional endpoints for COS are listed here: https://test.cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-endpoints
// However note that the staging enpoints vary signifcantly to the production
// To see the staging enpoints for your account go to https://test.cloud.ibm.com/objectstorage/overview and click instances to find your instance, and then click endpoints.
func GetCOSEndpoint(region string) string {
	if IsStagingMode() {
		// In staging, determine the endpointRegion from the PowerVS region
		var endpointRegion string
		switch region {
		case "au-syd", "br-sao", "ca-tor", "eu-de", "eu-es", "eu-gb", "jp-osa", "jp-tok", "us-east", "us-south", cosEndpointStagingRegionDefault:
			endpointRegion = cosEndpointStagingRegionDefault
		case "ca-mon", "in-che", "in-mum":
			endpointRegion = "dal10"
		default:
			// Default to us-west for unknown regions
			endpointRegion = cosEndpointStagingRegionDefault
		}
		return "https://s3." + endpointRegion + ".cloud-object-storage.test.appdomain.cloud"
	}
	return "https://s3." + region + ".cloud-object-storage.appdomain.cloud"
}

// GetCOSControlURL returns the appropriate COS control-plane URL based on staging mode.
// This is the service endpoint that returns the map of S3 endpoints (/v2/endpoints).
// For S3 bucket/object operations, use GetCOSEndpoint instead.
func GetCOSControlURL() string {
	if IsStagingMode() {
		return "https://control.test.cloud-object-storage.cloud.ibm.com"
	}
	return "https://control.cloud-object-storage.cloud.ibm.com"
}

// // GetCOSDirectEndpoint returns the appropriate COS S3 direct endpoint based on staging mode, PowerVS region, and COS region.
// // This returns the direct/private S3 endpoint for bucket/object operations from within IBM Cloud.
// // For public endpoints, use GetCOSEndpoint instead.
// //
// // In production, the format is: s3.direct.{cosRegion}.cloud-object-storage.appdomain.cloud
// // In staging, the format is: s3.{powervsRegion}.{geo}.direct.cloud-object-storage.test.appdomain.cloud
// //
// // Parameters:
// //   - powervsRegion: The PowerVS region code (e.g., "dal", "wdc", "lon")
// //   - cosRegion: The COS region (e.g., "us-south", "eu-de") - only used in production
// func GetCOSDirectEndpoint(powervsRegion string, cosRegion string) string {
// 	if IsStagingMode() {
// 		// In staging, determine the geo from the PowerVS region
// 		// US regions: dal, wdc, sjc -> "us"
// 		// EU regions: lon, mad, eu-de -> "eu"
// 		// AP regions: tok, osa, syd -> "ap"
// 		var geo string
// 		switch powervsRegion {
// 		case "dal", "wdc", "sjc":
// 			geo = "us"
// 		case "lon", "mad", "eu-de":
// 			geo = "eu"
// 		case "tok", "osa", "syd":
// 			geo = "ap"
// 		case "sao":
// 			geo = "br"
// 		default:
// 			// Default to us for unknown regions
// 			geo = "us"
// 		}
// 		return "https://s3." + powervsRegion + "." + geo + ".direct.cloud-object-storage.test.appdomain.cloud"
// 	}
// 	return "https://s3.direct." + cosRegion + ".cloud-object-storage.appdomain.cloud"
// }

// GetPowerURL returns the appropriate Power endpoint based on staging mode and region.
func GetPowerURL(zone string, region string) string {
	// region here isnt the same as us-south

	// Map zone → datacenter prefix (test env)
	var zoneToDC = map[string]string{
		// Dallas
		"dal10": "dal",
		"dal12": "dal",
		"dal14": "dal",

		// Washington DC
		"wdc06": "wdc",
		"wdc07": "wdc",

		// Toronto / Montreal / São Paulo
		"tor01": "tor",
		"mon01": "mon",
		"sao01": "sao",
		"sao04": "sao",
		"sao05": "sao",

		// Europe
		"lon04":   "lon",
		"lon06":   "lon",
		"eu-de-1": "eu-de",
		"eu-de-2": "eu-de",
		"mad02":   "mad",
		"mad04":   "mad",

		// APAC
		"syd04": "syd",
		"syd05": "syd",
		"tok04": "tok",
		"osa21": "osa",
		"che01": "che",
		"che02": "che",
		"che03": "che",
	}

	if IsStagingMode() {
		zone = strings.ToLower(zone)
		powerZone, ok := zoneToDC[zone]
		if !ok || powerZone == "" {
			logrus.Warnf("GetPowerURL: zone %q not found in staging zone map, falling back to region %q", zone, region)
			return "https://" + region + ".power-iaas.test.cloud.ibm.com"
		}
		return "https://" + powerZone + ".power-iaas.test.cloud.ibm.com"
	}
	return "https://" + region + ".power-iaas.cloud.ibm.com"
}

// GetVPCURL returns the appropriate VPC endpoint based on staging mode and region.
// Note: VPC staging uses the pattern {region}-stage01.iaasdev.cloud.ibm.com.
// The URL includes /v1 as the VPC API requires it.
func GetVPCURL(region string) string {
	if IsStagingMode() {
		return "https://" + region + "-stage01.iaasdev.cloud.ibm.com/v1"
	}
	return "https://" + region + ".iaas.cloud.ibm.com/v1"
}

// GetResourceManagerURL returns the appropriate Resource Manager endpoint based on staging mode.
// Note that the Resource Manager and Resource Controller endpoints are the same in staging.
func GetResourceManagerURL() string {
	if IsStagingMode() {
		return "https://resource-controller.test.cloud.ibm.com"
	}
	return "https://resource-controller.cloud.ibm.com"
}

// GetTransitGatewayURL returns the appropriate Transit Gateway endpoint based on staging mode.
func GetTransitGatewayURL() string {
	if IsStagingMode() {
		return "https://transit.test.cloud.ibm.com/v1"
	}
	return "https://transit.cloud.ibm.com/v1"
}

// GetCOSResourcePlanID returns the appropriate COS resource plan ID based on staging mode.
// In staging, IBM Cloud accounts use the Lite (free) plan; production uses the Standard plan.
// These IDs can be verified with: ibmcloud catalog service cloud-object-storage.
func GetCOSResourcePlanID() string {
	if IsStagingMode() {
		// Lite plan — used by non-paid staging accounts.
		return "2fdf0c08-2d32-4f46-84b5-32e0c92fffd8"
	}
	// Standard plan — used by production accounts.
	return "1e4e33e4-cfa6-4f12-9016-be594a6d5f87"
}
