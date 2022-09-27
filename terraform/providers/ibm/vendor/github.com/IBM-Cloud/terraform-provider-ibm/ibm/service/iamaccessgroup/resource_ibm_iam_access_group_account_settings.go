package iamaccessgroup

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMAccessGroupAccountSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceIBMIAMAccessGroupAccountSettingGet,
		UpdateContext: resourceIBMIAMAccessGroupAccountSettingSet,
		CreateContext: resourceIBMIAMAccessGroupAccountSettingSet,
		DeleteContext: resourceIBMIAMAccessGroupAccountSettingUnSet,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"public_access_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Flag to enable/disable public access groups",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the account",
			},
		},
	}
}

func resourceIBMIAMAccessGroupAccountSettingGet(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}
	getAccessGroupOptions := iamAccessGroupsClient.NewGetAccountSettingsOptions(userDetails.UserAccount)
	accountSetting, detailedResponse, err := iamAccessGroupsClient.GetAccountSettings(getAccessGroupOptions)
	if err != nil || accountSetting == nil {
		if detailedResponse != nil && detailedResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving access group account setting: %s. API Response: %s", err, detailedResponse))
	}
	d.SetId(*accountSetting.AccountID)
	d.Set("public_access_enabled", accountSetting.PublicAccessEnabled)
	d.Set("account_id", accountSetting.AccountID)
	return nil
}

func resourceIBMIAMAccessGroupAccountSettingSet(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	updateAccountSettingsOptions := iamAccessGroupsClient.NewUpdateAccountSettingsOptions(userDetails.UserAccount)
	publicAccessEnabled := d.Get("public_access_enabled").(bool)
	updateAccountSettingsOptions.PublicAccessEnabled = core.BoolPtr(publicAccessEnabled)
	accountSetting, detailedResponse, err := iamAccessGroupsClient.UpdateAccountSettings(updateAccountSettingsOptions)
	if err != nil || accountSetting == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error updating access public account setting: %s. API Response: %s", err, detailedResponse))
	}
	d.SetId(*accountSetting.AccountID)
	d.Set("public_access_enabled", *accountSetting.PublicAccessEnabled)
	return resourceIBMIAMAccessGroupAccountSettingGet(context, d, meta)
}

func resourceIBMIAMAccessGroupAccountSettingUnSet(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// DELETE NOT SUPPORTED
	d.SetId("")
	return nil
}
