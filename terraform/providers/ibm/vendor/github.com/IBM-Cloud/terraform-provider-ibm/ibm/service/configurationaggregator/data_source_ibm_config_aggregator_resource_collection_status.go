// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.92.0-af5c89a5-20240617-153232
 */

package configurationaggregator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/configuration-aggregator-go-sdk/configurationaggregatorv1"
)

func DataSourceIbmConfigAggregatorResourceCollectionStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmConfigAggregatorResourceCollectionStatusRead,

		Schema: map[string]*schema.Schema{
			"last_config_refresh_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp at which the configuration was last refreshed.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the resource collection.",
			},
		},
	}
}

func dataSourceIbmConfigAggregatorResourceCollectionStatusRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationAggregatorClient, err := meta.(conns.ClientSession).ConfigurationAggregatorV1()
	region := getConfigurationInstanceRegion(configurationAggregatorClient, d)
	instanceId := d.Get("instance_id").(string)
	configurationAggregatorClient = getClientWithConfigurationInstanceEndpoint(configurationAggregatorClient, instanceId, region)
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_config_aggregator_resource_collection_status", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getResourceCollectionStatusOptions := &configurationaggregatorv1.GetResourceCollectionStatusOptions{}

	statusResponse, _, err := configurationAggregatorClient.GetResourceCollectionStatusWithContext(context, getResourceCollectionStatusOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetResourceCollectionStatusWithContext failed: %s", err.Error()), "(Data) ibm_config_aggregator_resource_collection_status", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmConfigAggregatorResourceCollectionStatusID(d))

	if err = d.Set("last_config_refresh_time", flex.DateTimeToString(statusResponse.LastConfigRefreshTime)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting last_config_refresh_time: %s", err), "(Data) ibm_config_aggregator_resource_collection_status", "read", "set-last_config_refresh_time").GetDiag()
	}

	if err = d.Set("status", statusResponse.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_config_aggregator_resource_collection_status", "read", "set-status").GetDiag()
	}

	return nil
}

// dataSourceIbmConfigAggregatorResourceCollectionStatusID returns a reasonable ID for the list.
func dataSourceIbmConfigAggregatorResourceCollectionStatusID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
