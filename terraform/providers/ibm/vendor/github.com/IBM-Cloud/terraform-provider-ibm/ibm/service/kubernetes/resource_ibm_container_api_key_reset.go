// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMContainerAPIKeyReset() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMContainerAPIKeyResetUpdate,
		Read:   resourceIBMContainerAPIKeyResetRead,
		Update: resourceIBMContainerAPIKeyResetUpdate,
		Delete: resourceIBMContainerAPIKeyResetdelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region which api key has to be reset",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of Resource Group",
			},
			"reset_api_key": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Determines if apikey has to be reset or not",
				Default:     1,
			},
		},
	}
}

func resourceIBMContainerAPIKeyResetUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("reset_api_key") {
		apikeyClient, err := meta.(conns.ClientSession).ContainerAPI()
		if err != nil {
			return err
		}
		apikeyAPI := apikeyClient.Apikeys()
		region := d.Get("region").(string)
		targetEnv, err := getClusterTargetHeader(d, meta)
		if err != nil {
			return err
		}
		targetEnv.Region = region
		err = apikeyAPI.ResetApiKey(targetEnv)
		if err != nil {
			return err
		}
		if targetEnv.ResourceGroup == "" {
			defaultRg, err := flex.DefaultResourceGroup(meta)
			if err != nil {
				return err
			}
			targetEnv.ResourceGroup = defaultRg
		}

		d.SetId(fmt.Sprintf("%s/%s", region, targetEnv.ResourceGroup))
	}

	return nil
}
func resourceIBMContainerAPIKeyResetRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceIBMContainerAPIKeyResetdelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
