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

func DataSourceIbmMqcloudUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudUserRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQaaS service instance.",
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
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_user", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Read User failed: %s", err.Error()), "(Data) ibm_mqcloud_user", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listUsersOptions := &mqcloudv1.ListUsersOptions{}

	listUsersOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))

	var pager *mqcloudv1.UsersPager
	pager, err = mqcloudClient.NewUsersPager(listUsersOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_user", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UsersPager.GetAll() failed %s", err), "(Data) ibm_mqcloud_user", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchUsers []mqcloudv1.UserDetails
	var suppliedFilter bool
	var name string

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
		modelMap, err := DataSourceIbmMqcloudUserUserDetailsToMap(&modelItem) // #nosec G601
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_user", "read")
			return tfErr.GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("users", mapSlice); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting users %s", err), "(Data) ibm_mqcloud_user", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmMqcloudUserID returns a reasonable ID for the list.
func dataSourceIbmMqcloudUserID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmMqcloudUserUserDetailsToMap(model *mqcloudv1.UserDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["email"] = *model.Email
	modelMap["href"] = *model.Href
	return modelMap, nil
}
