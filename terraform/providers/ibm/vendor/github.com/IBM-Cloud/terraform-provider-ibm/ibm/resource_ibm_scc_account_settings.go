// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/adminserviceapiv1"
)

func resourceIBMSccAccountSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSccAccountSettingsUpdate,
		ReadContext:   resourceIbmSccAccountSettingsRead,
		UpdateContext: resourceIbmSccAccountSettingsUpdate,
		DeleteContext: resourceIbmSccAccountSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"location_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_scc_account_settings", "location_id"),
				Description:  "The programatic ID of the location that you want to work in.",
			},
		},
	}
}

func resourceIbmSccAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	adminServiceClient, err := meta.(ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOption := &adminserviceapiv1.GetSettingsOptions{}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}
	accountID := userDetails.userAccount

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
	adminServiceClient, err := meta.(ClientSession).AdminServiceApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}
	accountID := userDetails.userAccount

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

func resourceIBMSccAccountSettingsValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 2)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "location_id",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "us, eu, uk",
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_scc_account_settings", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSccAccountSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
