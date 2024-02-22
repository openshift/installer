// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package mqcloud

import (
	"context"
	"fmt"

	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func DataSourceIbmMqcloudApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudApplicationRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQ on Cloud service instance.",
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
		return diag.FromErr(err)
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Read Application failed %s", err))
	}

	serviceInstanceGuid := d.Get("service_instance_guid").(string)

	// Support for pagination
	offset := int64(0)
	limit := int64(25)
	allItems := []mqcloudv1.ApplicationDetails{}

	for {
		listApplicationsOptions := &mqcloudv1.ListApplicationsOptions{
			ServiceInstanceGuid: &serviceInstanceGuid,
			Limit:               &limit,
			Offset:              &offset,
		}

		result, response, err := mqcloudClient.ListApplicationsWithContext(context, listApplicationsOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Getting Applications %s\n%s", err, response))
		}
		if result == nil {
			return diag.FromErr(fmt.Errorf("List Applications returned nil"))
		}

		allItems = append(allItems, result.Applications...)

		// Check if the number of returned records is less than the limit
		if int64(len(result.Applications)) < limit {
			break
		}

		offset += limit
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchApplications []mqcloudv1.ApplicationDetails
	var name string
	var suppliedFilter bool

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
		modelMap, err := dataSourceIbmMqcloudApplicationApplicationDetailsToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("applications", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting applications: %s", err))
	}

	return nil
}

// dataSourceIbmMqcloudApplicationID returns a reasonable ID for the list.
func dataSourceIbmMqcloudApplicationID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmMqcloudApplicationApplicationDetailsToMap(model *mqcloudv1.ApplicationDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["create_api_key_uri"] = model.CreateApiKeyURI
	modelMap["href"] = model.Href
	return modelMap, nil
}
