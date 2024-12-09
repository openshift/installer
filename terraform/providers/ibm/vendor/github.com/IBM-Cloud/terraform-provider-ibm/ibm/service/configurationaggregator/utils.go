package configurationaggregator

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/configuration-aggregator-go-sdk/configurationaggregatorv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cloudEndpoint     = "cloud.ibm.com"
	testCloudEndpoint = "test.cloud.ibm.com"
)

func getConfigurationInstanceRegion(originalClient *configurationaggregatorv1.ConfigurationAggregatorV1, d *schema.ResourceData) string {
	_, ok := d.GetOk("region")
	if ok {
		return d.Get("region").(string)
	}
	baseUrl := originalClient.Service.GetServiceURL()
	url_01 := strings.Split(baseUrl, ".")[0]
	return (strings.Split(url_01, "://")[1])
}

func getClientWithConfigurationInstanceEndpoint(originalClient *configurationaggregatorv1.ConfigurationAggregatorV1, instanceId string, region string) *configurationaggregatorv1.ConfigurationAggregatorV1 {
	// build the api endpoint
	domain := cloudEndpoint
	if strings.Contains(os.Getenv("IBMCLOUD_IAM_API_ENDPOINT"), "test") {
		domain = testCloudEndpoint
	}
	endpoint := fmt.Sprintf("https://%s.apprapp.%s/apprapp/config_aggregator/v1/instances/%s", region, domain, instanceId)

	// clone the client and set endpoint
	newClient := &configurationaggregatorv1.ConfigurationAggregatorV1{
		Service: originalClient.Service.Clone(),
	}

	newClient.Service.SetServiceURL(endpoint)

	return newClient
}

func AddConfigurationAggregatorInstanceFields(resource *schema.Resource) *schema.Resource {
	resource.Schema["instance_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The ID of the configuration aggregator instance.",
	}
	resource.Schema["region"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "The region of the configuration aggregator instance.",
	}
	return resource
}

func updateClientURLWithInstanceEndpoint(id string, configsClient *configurationaggregatorv1.ConfigurationAggregatorV1, d *schema.ResourceData) (*configurationaggregatorv1.ConfigurationAggregatorV1, string, string, error) {

	idList, err := flex.IdParts(id)
	if err != nil || len(idList) < 2 {
		return configsClient, "", "", fmt.Errorf("Invalid Id %s. Error: %s", id, err)
	}

	region := idList[0]
	instanceId := idList[1]

	configsClient = getClientWithConfigurationInstanceEndpoint(configsClient, instanceId, region)

	return configsClient, region, instanceId, nil
}
