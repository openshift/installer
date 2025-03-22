// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
)

func DataSourceIBMIAMAccountSettingsExternalInteraction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAccountSettingsExternalInteractionGet,

		Schema: map[string]*schema.Schema{
			"external_account_identity_interaction": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "How external accounts can interact in relation to the requested account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity_types": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The settings for each identity type.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "The core set of properties associated with a user identity type.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"state": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The state of the user identity type.",
												},
												"external_allowed_accounts": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of accounts that the state applies to for the user identity type.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"service_id": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "The core set of properties associated with a serviceId identity type.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"state": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The state of the serviceId identity type.",
												},
												"external_allowed_accounts": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of accounts that the state applies to for the serviceId identity type.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"service": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "The core set of properties associated with a service identity type.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"state": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The state of the service identity type.",
												},
												"external_allowed_accounts": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of accounts that the state applies to for the service identity type.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"accept_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Language code for translations* `default` - English* `de` -  German (Standard)* `en` - English* `es` - Spanish (Spain)* `fr` - French (Standard)* `it` - Italian (Standard)* `ja` - Japanese* `ko` - Korean* `pt-br` - Portuguese (Brazil)* `zh-cn` - Chinese (Simplified, PRC)* `zh-tw` - (Chinese, Taiwan).",
			},
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account ID that the Access Management Account Settings belong to.",
			},
		},
	}
}

func dataSourceIBMAccountSettingsExternalInteractionGet(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_iam_account_settings_external_interaction", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	var accountID string

	if _, ok := d.GetOk("account_id"); ok {
		accountID = d.Get("account_id").(string)
	}

	getSettingsOptions := &iampolicymanagementv1.GetSettingsOptions{
		AccountID: &accountID,
	}

	if _, ok := d.GetOk("accept_language"); ok {
		getSettingsOptions.SetAcceptLanguage(d.Get("accept_language").(string))
	}

	amAccountSettings, _, err := iamPolicyManagementClient.GetSettingsWithContext(context, getSettingsOptions)

	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSettingsWithContext failed: %s", err.Error()), "ibm_iam_account_settings_external_interaction", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.Set("external_account_identity_interaction", flex.FlattenAMSettingsExternalIdentityInteraction(amAccountSettings))

	d.SetId(accountID)
	return nil
}
