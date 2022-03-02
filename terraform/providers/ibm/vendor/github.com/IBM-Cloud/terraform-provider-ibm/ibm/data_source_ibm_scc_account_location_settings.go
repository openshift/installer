// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/adminserviceapiv1"
)

func dataSourceIbmSccAccountLocationSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSccAccountLocationSettingsRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The programatic ID of the location that you want to work in.",
			},
		},
	}
}

func dataSourceIbmSccAccountLocationSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminServiceApiClient, err := meta.(ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions := &adminserviceapiv1.GetSettingsOptions{}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions.SetAccountID(userDetails.userAccount)

	locationSettings, response, err := adminServiceApiClient.GetSettingsWithContext(context, getSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId(*locationSettings.Location.ID)
	if err = d.Set("id", locationSettings.Location.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}

	return nil
}
