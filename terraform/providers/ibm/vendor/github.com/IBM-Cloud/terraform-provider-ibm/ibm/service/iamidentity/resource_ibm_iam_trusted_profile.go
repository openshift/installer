// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func ResourceIBMIAMTrustedProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIamTrustedProfileCreate,
		ReadContext:   resourceIBMIamTrustedProfileRead,
		UpdateContext: resourceIBMIamTrustedProfileUpdate,
		DeleteContext: resourceIBMIamTrustedProfileDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the trusted profile. The name is checked for uniqueness. Therefore trusted profiles with the same names can not exist in the same account.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID of the trusted profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional description of the trusted profile. The 'description' property is only available if a description was provided during creation of trusted profile.",
			},
			"profile_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of this trusted profile.",
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
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template id the profile was created from.",
			},
			"assignment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of assignment that assigned the template.",
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
							Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceIBMIamTrustedProfileCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createProfileOptions := &iamidentityv1.CreateProfileOptions{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := userDetails.UserAccount

	createProfileOptions.SetName(d.Get("name").(string))
	createProfileOptions.SetAccountID(accountID)
	if _, ok := d.GetOk("description"); ok {
		createProfileOptions.SetDescription(d.Get("description").(string))
	}

	trustedProfile, response, err := iamIdentityClient.CreateProfile(createProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProfileWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateProfileWithContext failed %s\n%s", err, response))
	}

	d.SetId(*trustedProfile.ID)

	return resourceIBMIamTrustedProfileRead(context, d, meta)
}

func resourceIBMIamTrustedProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileOptions := &iamidentityv1.GetProfileOptions{}

	getProfileOptions.SetProfileID(d.Id())

	trustedProfile, response, err := iamIdentityClient.GetProfile(getProfileOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProfile failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProfile failed %s\n%s", err, response))
	}

	if err = d.Set("name", trustedProfile.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("account_id", trustedProfile.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting account_id: %s", err))
	}
	if err = d.Set("description", trustedProfile.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}
	if err = d.Set("profile_id", trustedProfile.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting id: %s", err))
	}
	if err = d.Set("entity_tag", trustedProfile.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting entity_tag: %s", err))
	}
	if err = d.Set("crn", trustedProfile.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
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
	if err = d.Set("ims_account_id", flex.IntValue(trustedProfile.ImsAccountID)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting ims_account_id: %s", err))
	}
	if err = d.Set("ims_user_id", flex.IntValue(trustedProfile.ImsUserID)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting ims_user_id: %s", err))
	}
	if err = d.Set("template_id", trustedProfile.TemplateID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting template_id: %s", err))
	}
	if err = d.Set("assignment_id", trustedProfile.AssignmentID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting assignment_id: %s", err))
	}

	history := []map[string]interface{}{}
	if trustedProfile.History != nil {
		for _, historyItem := range trustedProfile.History {
			historyItemMap := resourceIBMIamTrustedProfileEnityHistoryRecordToMap(historyItem)
			history = append(history, historyItemMap)
		}
	}
	if err = d.Set("history", history); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting history: %s", err))
	}

	return nil
}

func resourceIBMIamTrustedProfileResponseContextToMap(responseContext iamidentityv1.ResponseContext) map[string]interface{} {
	responseContextMap := map[string]interface{}{}

	if responseContext.TransactionID != nil {
		responseContextMap["transaction_id"] = responseContext.TransactionID
	}
	if responseContext.Operation != nil {
		responseContextMap["operation"] = responseContext.Operation
	}
	if responseContext.UserAgent != nil {
		responseContextMap["user_agent"] = responseContext.UserAgent
	}
	if responseContext.URL != nil {
		responseContextMap["url"] = responseContext.URL
	}
	if responseContext.InstanceID != nil {
		responseContextMap["instance_id"] = responseContext.InstanceID
	}
	if responseContext.ThreadID != nil {
		responseContextMap["thread_id"] = responseContext.ThreadID
	}
	if responseContext.Host != nil {
		responseContextMap["host"] = responseContext.Host
	}
	if responseContext.StartTime != nil {
		responseContextMap["start_time"] = responseContext.StartTime
	}
	if responseContext.EndTime != nil {
		responseContextMap["end_time"] = responseContext.EndTime
	}
	if responseContext.ElapsedTime != nil {
		responseContextMap["elapsed_time"] = responseContext.ElapsedTime
	}
	if responseContext.ClusterName != nil {
		responseContextMap["cluster_name"] = responseContext.ClusterName
	}

	return responseContextMap
}

func resourceIBMIamTrustedProfileEnityHistoryRecordToMap(enityHistoryRecord iamidentityv1.EnityHistoryRecord) map[string]interface{} {
	enityHistoryRecordMap := map[string]interface{}{}

	enityHistoryRecordMap["timestamp"] = enityHistoryRecord.Timestamp
	enityHistoryRecordMap["iam_id"] = enityHistoryRecord.IamID
	enityHistoryRecordMap["iam_id_account"] = enityHistoryRecord.IamIDAccount
	enityHistoryRecordMap["action"] = enityHistoryRecord.Action
	enityHistoryRecordMap["params"] = enityHistoryRecord.Params
	enityHistoryRecordMap["message"] = enityHistoryRecord.Message

	return enityHistoryRecordMap
}

func resourceIBMIamTrustedProfileUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateProfileOptions := &iamidentityv1.UpdateProfileOptions{}

	updateProfileOptions.SetIfMatch("*")
	updateProfileOptions.SetProfileID(d.Id())
	updateProfileOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("description"); ok {
		updateProfileOptions.SetDescription(d.Get("description").(string))
	}

	_, response, err := iamIdentityClient.UpdateProfile(updateProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateProfile failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("UpdateProfile failed %s\n%s", err, response))
	}

	return resourceIBMIamTrustedProfileRead(context, d, meta)
}

func resourceIBMIamTrustedProfileDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProfileOptions := &iamidentityv1.DeleteProfileOptions{}

	deleteProfileOptions.SetProfileID(d.Id())

	response, err := iamIdentityClient.DeleteProfile(deleteProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteProfile failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteProfile failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
