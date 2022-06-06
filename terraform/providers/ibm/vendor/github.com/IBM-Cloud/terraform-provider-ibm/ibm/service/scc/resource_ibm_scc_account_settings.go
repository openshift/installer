// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v3/adminserviceapiv1"
)

func ResourceIBMSccAccountSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSccAccountSettingsUpdate,
		ReadContext:   resourceIbmSccAccountSettingsRead,
		UpdateContext: resourceIbmSccAccountSettingsUpdate,
		DeleteContext: resourceIbmSccAccountSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"location_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_account_settings", "location_id"),
				Description:  "The programatic ID of the location that you want to work in.",
			},
		},
	}
}

func resourceIbmSccAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminServiceClient, err := meta.(conns.ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOption := &adminserviceapiv1.GetSettingsOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}
	accountID := userDetails.UserAccount

	getAccountSettingsOption.SetAccountID(accountID)

	accountSettings, response, err := adminServiceClient.GetSettingsWithContext(context, getAccountSettingsOption)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		log.Printf("[DEBUG] GetSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId(*accountSettings.Location.ID)

	return nil

}

func resourceIbmSccAccountSettingsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminServiceClient, err := meta.(conns.ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}
	accountID := userDetails.UserAccount

	locationID := d.Get("location_id").(string)
	updateAccountSettingsOption := &adminserviceapiv1.PatchAccountSettingsOptions{}
	updateAccountSettingsOption.SetAccountID(accountID)
	updateAccountSettingsOption.SetLocation(&adminserviceapiv1.LocationID{
		ID: core.StringPtr(locationID),
	})

	_, response, err := adminServiceClient.PatchAccountSettingsWithContext(context, updateAccountSettingsOption)
	if err != nil {
		log.Printf("[DEBUG] PatchAccountSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("PatchAccountSettingsWithContext failed %s\n%s", err, response))
	}

	d.SetId(locationID)

	return resourceIbmSccAccountSettingsRead(context, d, meta)
}

func ResourceIBMSccAccountSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 2)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "location_id",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "us, eu, uk",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_account_settings", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSccAccountSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
