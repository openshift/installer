package logs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cloudEndpoint     = "cloud.ibm.com"
	testCloudEndpoint = "test.cloud.ibm.com"
)

func getLogsInstanceRegion(originalClient *logsv0.LogsV0, d *schema.ResourceData) string {
	_, ok := d.GetOk("region")
	if ok {
		return d.Get("region").(string)
	}
	// extract region from base URL (provider config)
	// base url is like that : "https://api.private.eu-gb.logs.cloud.ibm.com"
	baseUrl := originalClient.Service.GetServiceURL()
	u := strings.Replace(baseUrl, "private.", "", 1)

	return strings.Split(u, ".")[1]
}

// Clone the base logs client and set the API endpoint per the instance
func getLogsInstanceEndpointType(originalClient *logsv0.LogsV0, d *schema.ResourceData) string {
	_, ok := d.GetOk("endpoint_type")
	if ok {
		return d.Get("endpoint_type").(string)
	}
	baseUrl := originalClient.Service.GetServiceURL()
	if strings.Contains(baseUrl, "private.") {
		return "private"
	}

	return "public"
}

// <instance_id>.api.eu-gb.logs.test.cloud.ibm.com
// Clone the base logs client and set the API endpoint per the instance
func getClientWithLogsInstanceEndpoint(originalClient *logsv0.LogsV0, instanceId string, region string, endpointType string) *logsv0.LogsV0 {
	// build the api endpoint
	domain := cloudEndpoint
	if strings.Contains(os.Getenv("IBMCLOUD_IAM_API_ENDPOINT"), "test") {
		domain = testCloudEndpoint
	}
	// getting originalConfigServiceURL to not miss filemap precedence from the url constructed in config.go file
	originalConfigServiceURL := originalClient.GetServiceURL()

	log.Printf("Service URL from the config.go file %s", originalConfigServiceURL)

	var endpoint string
	if endpointType == "private" {
		if strings.Contains(originalConfigServiceURL, fmt.Sprintf("https://%s.api.private.%s.logs.%s", instanceId, region, domain)) {
			endpoint = originalConfigServiceURL
		} else {
			endpoint = fmt.Sprintf("https://%s.api.private.%s.logs.%s:3443", instanceId, region, domain)
		}
	} else {
		if strings.Contains(originalConfigServiceURL, fmt.Sprintf("https://%s.api.%s.logs.%s", instanceId, region, domain)) {
			endpoint = originalConfigServiceURL
		} else {
			endpoint = fmt.Sprintf("https://%s.api.%s.logs.%s", instanceId, region, domain)
		}
	}
	// clone the client and set endpoint
	newClient := &logsv0.LogsV0{
		Service: originalClient.Service.Clone(),
	}

	endpoint = conns.EnvFallBack([]string{"IBMCLOUD_LOGS_API_ENDPOINT"}, endpoint)

	log.Printf("Constructing client with new service URL %s", endpoint)

	newClient.Service.SetServiceURL(endpoint)

	return newClient
}

// Add the fields needed for building the instance endpoint to the given schema
func AddLogsInstanceFields(resource *schema.Resource) *schema.Resource {
	resource.Schema["instance_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The ID of the logs instance.",
	}
	resource.Schema["region"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "The region of the logs instance.",
	}
	resource.Schema["endpoint_type"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "public or private.",
	}

	return resource
}

func updateClientURLWithInstanceEndpoint(id string, logsClient *logsv0.LogsV0, d *schema.ResourceData) (*logsv0.LogsV0, string, string, string, error) {

	idList, err := flex.IdParts(id)
	if err != nil || len(idList) < 2 {
		return logsClient, "", "", "", fmt.Errorf("Invalid Id %s. Error: %s", id, err)
	}

	region := idList[0]
	instanceId := idList[1]
	var resourceId string
	if len(idList) > 2 {
		resourceId = idList[2]
	}

	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	return logsClient, region, instanceId, resourceId, nil
}
