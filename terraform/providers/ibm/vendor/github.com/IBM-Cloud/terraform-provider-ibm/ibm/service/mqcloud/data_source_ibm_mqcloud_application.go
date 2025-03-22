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

func DataSourceIbmMqcloudApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudApplicationRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the application - conforming to MQ rules.",
			},
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of applications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the application which was allocated on creation, and can be used for delete calls.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the application - conforming to MQ rules.",
						},
						"create_api_key_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URI to create a new apikey for the application.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this application.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmMqcloudApplicationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_application", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Read Application failed: %s", err.Error()), "(Data) ibm_mqcloud_application", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{}

	listApplicationsOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))

	var pager *mqcloudv1.ApplicationsPager
	pager, err = mqcloudClient.NewApplicationsPager(listApplicationsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_application", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ApplicationsPager.GetAll() failed %s", err), "(Data) ibm_mqcloud_application", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchApplications []mqcloudv1.ApplicationDetails
	var suppliedFilter bool
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range allItems {
			if *data.Name == name {
				matchApplications = append(matchApplications, data)
			}
		}
	} else {
		matchApplications = allItems
	}

	allItems = matchApplications

	if suppliedFilter {
		if len(allItems) == 0 {
			return diag.FromErr(fmt.Errorf("No Application found with name: \"%s\"", name))
		}
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmMqcloudApplicationID(d))
	}

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelItem := modelItem
		modelMap, err := DataSourceIbmMqcloudApplicationApplicationDetailsToMap(&modelItem) // #nosec G601
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_application", "read")
			return tfErr.GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("applications", mapSlice); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting applications %s", err), "(Data) ibm_mqcloud_application", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmMqcloudApplicationID returns a reasonable ID for the list.
func dataSourceIbmMqcloudApplicationID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmMqcloudApplicationApplicationDetailsToMap(model *mqcloudv1.ApplicationDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["create_api_key_uri"] = *model.CreateApiKeyURI
	modelMap["href"] = *model.Href
	return modelMap, nil
}
