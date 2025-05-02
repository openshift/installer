// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileRead,

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile to get.",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile",
					"profile_id"),
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the trusted profile details object. You need to specify this value when updating the trusted profile to avoid stale updates.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud Resource Name of the item. Example Cloud Resource Name: 'crn:v1:bluemix:public:iam-identity:us-south:a/myaccount::profile:Profile-94497d0d-2ac3-41bf-a993-a49d1b14627c'.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the trusted profile. The name is checked for uniqueness. Therefore trusted profiles with the same names can not exist in the same account.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The optional description of the trusted profile. The 'description' property is only available if a description was provided during a create of a trusted profile.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the creation date in ISO format.",
			},
			"modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the last modification date in ISO format.",
			},
			"iam_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The iam_id of this trusted profile.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the account that this trusted profile belong to.",
			},
			"ims_account_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IMS acount ID of the trusted profile.",
			},
			"ims_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IMS user ID of the trusted profile.",
			},
			"history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History of the trusted profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action of the history entry.",
						},
						"params": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Params of the history entry.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message which summarizes the executed action.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMIamTrustedProfileValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:id"},
			Required:                   true})

	iBMIamTrustedProfileValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile", Schema: validateSchema}
	return &iBMIamTrustedProfileValidator
}

func dataSourceIBMIamTrustedProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileOptions := &iamidentityv1.GetProfileOptions{}

	getProfileOptions.SetProfileID(d.Get("profile_id").(string))

	trustedProfile, response, err := iamIdentityClient.GetProfile(getProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProfile failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProfile failed %s\n%s", err, response))
	}

	d.SetId(*trustedProfile.ID)

	if err = d.Set("entity_tag", trustedProfile.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting entity_tag: %s", err))
	}
	if err = d.Set("crn", trustedProfile.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("name", trustedProfile.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("description", trustedProfile.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(trustedProfile.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", flex.DateTimeToString(trustedProfile.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting modified_at: %s", err))
	}
	if err = d.Set("iam_id", trustedProfile.IamID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting iam_id: %s", err))
	}
	if err = d.Set("account_id", trustedProfile.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting account_id: %s", err))
	}
	if err = d.Set("ims_account_id", flex.IntValue(trustedProfile.ImsAccountID)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting ims_account_id: %s", err))
	}
	if err = d.Set("ims_user_id", flex.IntValue(trustedProfile.ImsUserID)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting ims_user_id: %s", err))
	}

	err = d.Set("history", dataSourceTrustedProfileFlattenHistory(trustedProfile.History))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting history %s", err))
	}

	return nil
}

func dataSourceTrustedProfileFlattenHistory(result []iamidentityv1.EnityHistoryRecord) (history []map[string]interface{}) {
	for _, historyItem := range result {
		history = append(history, dataSourceTrustedProfileHistoryToMap(historyItem))
	}

	return history
}

func dataSourceTrustedProfileHistoryToMap(historyItem iamidentityv1.EnityHistoryRecord) (historyMap map[string]interface{}) {
	historyMap = map[string]interface{}{}

	if historyItem.Timestamp != nil {
		historyMap["timestamp"] = historyItem.Timestamp
	}
	if historyItem.IamID != nil {
		historyMap["iam_id"] = historyItem.IamID
	}
	if historyItem.IamIDAccount != nil {
		historyMap["iam_id_account"] = historyItem.IamIDAccount
	}
	if historyItem.Action != nil {
		historyMap["action"] = historyItem.Action
	}
	if historyItem.Params != nil {
		historyMap["params"] = historyItem.Params
	}
	if historyItem.Message != nil {
		historyMap["message"] = historyItem.Message
	}

	return historyMap
}
