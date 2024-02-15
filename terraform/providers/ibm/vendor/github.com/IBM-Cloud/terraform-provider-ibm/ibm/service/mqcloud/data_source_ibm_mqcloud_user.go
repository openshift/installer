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

func DataSourceIbmMqcloudUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudUserRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQ on Cloud service instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The shortname of the user that will be used as the IBM MQ administrator in interactions with a queue manager for this service instance.",
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the user which was allocated on creation, and can be used for delete calls.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The shortname of the user that will be used as the IBM MQ administrator in interactions with a queue manager for this service instance.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The email of the user.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for the user details.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmMqcloudUserRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Read User failed %s", err))
	}

	serviceInstanceGuid := d.Get("service_instance_guid").(string)

	// Support for pagination
	offset := int64(0)
	limit := int64(25)
	allItems := []mqcloudv1.UserDetails{}

	for {
		listUsersOptions := &mqcloudv1.ListUsersOptions{
			ServiceInstanceGuid: &serviceInstanceGuid,
			Limit:               &limit,
			Offset:              &offset,
		}

		result, response, err := mqcloudClient.ListUsersWithContext(context, listUsersOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Getting Users %s\n%s", err, response))
		}
		if result == nil {
			return diag.FromErr(fmt.Errorf("List Users returned nil"))
		}

		allItems = append(allItems, result.Users...)

		// Check if the number of returned records is less than the limit
		if int64(len(result.Users)) < limit {
			break
		}

		offset += limit
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchUsers []mqcloudv1.UserDetails
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range allItems {
			if *data.Name == name {
				matchUsers = append(matchUsers, data)
			}
		}
	} else {
		matchUsers = allItems
	}

	allItems = matchUsers

	if suppliedFilter {
		if len(allItems) == 0 {
			return diag.FromErr(fmt.Errorf("No User found with name: \"%s\"", name))
		}
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmMqcloudUserID(d))
	}

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelItem := modelItem
		modelMap, err := dataSourceIbmMqcloudUserUserDetailsToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("users", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting users: %s", err))
	}

	return nil
}

// dataSourceIbmMqcloudUserID returns a reasonable ID for the list.
func dataSourceIbmMqcloudUserID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmMqcloudUserUserDetailsToMap(model *mqcloudv1.UserDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["email"] = model.Email
	modelMap["href"] = model.Href
	return modelMap, nil
}
