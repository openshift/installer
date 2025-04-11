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

func DataSourceIbmMqcloudQueueManager() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudQueueManagerRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A queue manager name conforming to MQ restrictions.",
			},
			"queue_managers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of queue managers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the queue manager which was allocated on creation, and can be used for delete calls.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A queue manager name conforming to MQ restrictions.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A displayable name for the queue manager - limited only in length.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The locations in which the queue manager could be deployed.",
						},
						"size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The queue manager sizes of deployment available.",
						},
						"status_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reference uri to get deployment status of the queue manager.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MQ version of the queue manager.",
						},
						"web_console_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The url through which to access the web console for this queue manager.",
						},
						"rest_api_endpoint_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The url through which to access REST APIs for this queue manager.",
						},
						"administrator_api_endpoint_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The url through which to access the Admin REST APIs for this queue manager.",
						},
						"connection_info_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The uri through which the CDDT for this queue manager can be obtained.",
						},
						"date_created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "RFC3339 formatted UTC date for when the queue manager was created.",
						},
						"upgrade_available": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Describes whether an upgrade is available for this queue manager.",
						},
						"available_upgrade_versions_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The uri through which the available versions to upgrade to can be found for this queue manager.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this queue manager.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmMqcloudQueueManagerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_queue_manager", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Read Queue Manager failed: %s", err.Error()), "(Data) ibm_mqcloud_queue_manager", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listQueueManagersOptions := &mqcloudv1.ListQueueManagersOptions{}

	listQueueManagersOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))

	var pager *mqcloudv1.QueueManagersPager
	pager, err = mqcloudClient.NewQueueManagersPager(listQueueManagersOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_queue_manager", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("QueueManagersPager.GetAll() failed %s", err), "(Data) ibm_mqcloud_queue_manager", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchQueueManagers []mqcloudv1.QueueManagerDetails
	var suppliedFilter bool
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range allItems {
			if *data.Name == name {
				matchQueueManagers = append(matchQueueManagers, data)
			}
		}
	} else {
		matchQueueManagers = allItems
	}

	allItems = matchQueueManagers

	if suppliedFilter {
		if len(allItems) == 0 {
			return diag.FromErr(fmt.Errorf("No Queue Manager found with name: \"%s\"", name))
		}
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmMqcloudQueueManagerID(d))
	}

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelItem := modelItem
		modelMap, err := DataSourceIbmMqcloudQueueManagerQueueManagerDetailsToMap(&modelItem) // #nosec G601
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_queue_manager", "read")
			return tfErr.GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("queue_managers", mapSlice); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting queue_managers %s", err), "(Data) ibm_mqcloud_queue_manager", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmMqcloudQueueManagerID returns a reasonable ID for the list.
func dataSourceIbmMqcloudQueueManagerID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmMqcloudQueueManagerQueueManagerDetailsToMap(model *mqcloudv1.QueueManagerDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["display_name"] = *model.DisplayName
	modelMap["location"] = *model.Location
	modelMap["size"] = *model.Size
	modelMap["status_uri"] = *model.StatusURI
	modelMap["version"] = *model.Version
	modelMap["web_console_url"] = *model.WebConsoleURL
	modelMap["rest_api_endpoint_url"] = *model.RestApiEndpointURL
	modelMap["administrator_api_endpoint_url"] = *model.AdministratorApiEndpointURL
	modelMap["connection_info_uri"] = *model.ConnectionInfoURI
	modelMap["date_created"] = model.DateCreated.String()
	modelMap["upgrade_available"] = *model.UpgradeAvailable
	modelMap["available_upgrade_versions_uri"] = *model.AvailableUpgradeVersionsURI
	modelMap["href"] = *model.Href
	return modelMap, nil
}
