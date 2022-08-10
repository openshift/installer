// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/adminserviceapiv1"
)

func DataSourceIBMSccAccountLocationSettings() *schema.Resource {
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
	adminServiceApiClient, err := meta.(conns.ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions := &adminserviceapiv1.GetSettingsOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	getSettingsOptions.SetAccountID(userDetails.UserAccount)

	accountSettings, response, err := adminServiceApiClient.GetSettingsWithContext(context, getSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSettingsWithContext failed %s\n%s", err, response))
	}

	locationSettings := accountSettings.Location
	d.SetId(*accountSettings.Location.ID)

	if err = d.Set("id", locationSettings.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}

	if d.HasChanges() {
		d.SetId(dataSourceIbmSccAccountLocationSettingsID(d))
	}

	return nil
}

func dataSourceIbmSccAccountLocationSettingsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
