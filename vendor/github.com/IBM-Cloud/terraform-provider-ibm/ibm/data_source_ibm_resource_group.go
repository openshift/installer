// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/managementv2"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMResourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMResourceGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:   "Resource group name",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"is_default"},
			},
			"is_default": {
				Description:   "Default Resource group",
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
		},
	}
}

func dataSourceIBMResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	rsManagementAPI, err := meta.(ClientSession).ResourceManagementAPIv2()
	if err != nil {
		return err
	}
	rsGroup := rsManagementAPI.ResourceGroup()

	var defaultGrp bool
	if group, ok := d.GetOk("is_default"); ok {
		defaultGrp = group.(bool)
	}
	var name string
	if n, ok := d.GetOk("name"); ok {
		name = n.(string)
	}
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	accountID := userDetails.userAccount
	var grp []models.ResourceGroupv2
	if defaultGrp {
		resourceGroupQuery := managementv2.ResourceGroupQuery{
			Default:   true,
			AccountID: accountID,
		}

		grp, err = rsGroup.List(&resourceGroupQuery)

		if err != nil {
			return fmt.Errorf("Error retrieving default resource group: %s", err)
		}
		d.SetId(grp[0].ID)

	} else if name != "" {
		resourceGroupQuery := &managementv2.ResourceGroupQuery{
			AccountID: accountID,
		}
		grp, err := rsGroup.FindByName(resourceGroupQuery, name)
		if err != nil {
			return fmt.Errorf("Error retrieving resource group %s: %s", name, err)
		}
		d.SetId(grp[0].ID)

	} else {
		return fmt.Errorf("Missing required properties. Need a resource group name, or the is_default true")
	}

	return nil
}
