// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package mqcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func DataSourceIbmMqcloudQueueManagerOptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudQueueManagerOptionsRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"locations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of deployment locations.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sizes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of queue manager sizes.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of queue manager versions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"latest_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest Queue manager version.",
			},
		},
	}
}

func dataSourceIbmMqcloudQueueManagerOptionsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_queue_manager_options", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Read Queue Manager Options failed: %s", err.Error()), "(Data) ibm_mqcloud_queue_manager_options", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getOptionsOptions := &mqcloudv1.GetOptionsOptions{}

	getOptionsOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))

	configurationOptions, _, err := mqcloudClient.GetOptionsWithContext(context, getOptionsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetOptionsWithContext failed: %s", err.Error()), "(Data) ibm_mqcloud_queue_manager_options", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmMqcloudQueueManagerOptionsID(d))

	if err = d.Set("latest_version", configurationOptions.LatestVersion); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting latest_version: %s", err), "(Data) ibm_mqcloud_queue_manager_options", "read", "set-latest_version").GetDiag()
	}

	if err = d.Set("locations", configurationOptions.Locations); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting locations: %s", err), "(Data) ibm_mqcloud_queue_manager_options", "read", "set-locations").GetDiag()
	}

	if err = d.Set("versions", configurationOptions.Versions); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting versions: %s", err), "(Data) ibm_mqcloud_queue_manager_options", "read", "set-versions").GetDiag()
	}

	if err = d.Set("sizes", configurationOptions.Sizes); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting sizes: %s", err), "(Data) ibm_mqcloud_queue_manager_options", "read", "set-sizes").GetDiag()
	}

	return nil
}

// dataSourceIbmMqcloudQueueManagerOptionsID returns a reasonable ID for the list.
func dataSourceIbmMqcloudQueueManagerOptionsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
