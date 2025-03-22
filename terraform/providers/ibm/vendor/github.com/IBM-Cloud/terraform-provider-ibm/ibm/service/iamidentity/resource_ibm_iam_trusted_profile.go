// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.98.0-8be2046a-20241205-162752
 */

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the trusted profile. The name is checked for uniqueness. Therefore trusted profiles with the same names can not exist in the same account.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional description of the trusted profile. The 'description' property is only available if a description was provided during a create of a trusted profile.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the account that this trusted profile belong to.",
			},
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the unique identifier of the trusted profile. Example:'Profile-94497d0d-2ac3-41bf-a993-a49d1b14627c'.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the trusted profile details object. You need to specify this value when updating the trusted profile to avoid stale updates.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud Resource Name of the item. Example Cloud Resource Name: 'crn:v1:bluemix:public:iam-identity:us-south:a/myaccount::profile:Profile-94497d0d-2ac3-41bf-a993-a49d1b14627c'.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the creation date in ISO format.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the last modification date in ISO format.",
			},
			"iam_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The iam_id of this trusted profile.",
			},
			"template_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the IAM template that was used to create an enterprise-managed trusted profile in your account. When returned, this indicates that the trusted profile is created from and managed by a template in the root enterprise account.",
			},
			"assignment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the assignment that was used to create an enterprise-managed trusted profile in your account. When returned, this indicates that the trusted profile is created from and managed by a template in the root enterprise account.",
			},
			"ims_account_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IMS acount ID of the trusted profile.",
			},
			"ims_user_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IMS user ID of the trusted profile.",
			},
			"history": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History of the trusted profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action of the history entry.",
						},
						"params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Params of the history entry.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"message": &schema.Schema{
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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

	trustedProfile, _, err := iamIdentityClient.CreateProfileWithContext(context, createProfileOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateProfileWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*trustedProfile.ID)

	return resourceIBMIamTrustedProfileRead(context, d, meta)
}

func resourceIBMIamTrustedProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProfileOptions := &iamidentityv1.GetProfileOptions{}

	getProfileOptions.SetProfileID(d.Id())

	trustedProfile, response, err := iamIdentityClient.GetProfileWithContext(context, getProfileOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProfileWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", trustedProfile.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-name").GetDiag()
	}
	if !core.IsNil(trustedProfile.Description) {
		if err = d.Set("description", trustedProfile.Description); err != nil {
			err = fmt.Errorf("Error setting description: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-description").GetDiag()
		}
	}
	if err = d.Set("account_id", trustedProfile.AccountID); err != nil {
		err = fmt.Errorf("Error setting account_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-account_id").GetDiag()
	}
	if err = d.Set("profile_id", trustedProfile.ID); err != nil {
		err = fmt.Errorf("Error setting id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-id").GetDiag()
	}
	if err = d.Set("entity_tag", trustedProfile.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-entity_tag").GetDiag()
	}
	if err = d.Set("crn", trustedProfile.CRN); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-crn").GetDiag()
	}
	if !core.IsNil(trustedProfile.CreatedAt) {
		if err = d.Set("created_at", flex.DateTimeToString(trustedProfile.CreatedAt)); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-created_at").GetDiag()
		}
	}
	if !core.IsNil(trustedProfile.ModifiedAt) {
		if err = d.Set("modified_at", flex.DateTimeToString(trustedProfile.ModifiedAt)); err != nil {
			err = fmt.Errorf("Error setting modified_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-modified_at").GetDiag()
		}
	}
	if err = d.Set("iam_id", trustedProfile.IamID); err != nil {
		err = fmt.Errorf("Error setting iam_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-iam_id").GetDiag()
	}
	if !core.IsNil(trustedProfile.TemplateID) {
		if err = d.Set("template_id", trustedProfile.TemplateID); err != nil {
			err = fmt.Errorf("Error setting template_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-template_id").GetDiag()
		}
	}
	if !core.IsNil(trustedProfile.AssignmentID) {
		if err = d.Set("assignment_id", trustedProfile.AssignmentID); err != nil {
			err = fmt.Errorf("Error setting assignment_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-assignment_id").GetDiag()
		}
	}
	if !core.IsNil(trustedProfile.ImsAccountID) {
		if err = d.Set("ims_account_id", flex.IntValue(trustedProfile.ImsAccountID)); err != nil {
			err = fmt.Errorf("Error setting ims_account_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-ims_account_id").GetDiag()
		}
	}
	if !core.IsNil(trustedProfile.ImsUserID) {
		if err = d.Set("ims_user_id", flex.IntValue(trustedProfile.ImsUserID)); err != nil {
			err = fmt.Errorf("Error setting ims_user_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-ims_user_id").GetDiag()
		}
	}
	history := []map[string]interface{}{}
	if !core.IsNil(trustedProfile.History) {
		for _, historyItem := range trustedProfile.History {
			historyItemMap, err := ResourceIBMIamTrustedProfileEnityHistoryRecordToMap(&historyItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "history-to-map").GetDiag()
			}
			history = append(history, historyItemMap)
		}
	}
	if err = d.Set("history", history); err != nil {
		err = fmt.Errorf("Error setting history: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "read", "set-history").GetDiag()

	}

	return nil
}

func resourceIBMIamTrustedProfileUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateProfileOptions := &iamidentityv1.UpdateProfileOptions{}

	updateProfileOptions.SetIfMatch("*")
	updateProfileOptions.SetProfileID(d.Id())
	if _, ok := d.GetOk("name"); ok {
		updateProfileOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		updateProfileOptions.SetDescription(d.Get("description").(string))
	}

	_, _, err = iamIdentityClient.UpdateProfileWithContext(context, updateProfileOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateProfileWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMIamTrustedProfileRead(context, d, meta)
}

func resourceIBMIamTrustedProfileDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteProfileOptions := &iamidentityv1.DeleteProfileOptions{}

	deleteProfileOptions.SetProfileID(d.Id())

	_, err = iamIdentityClient.DeleteProfileWithContext(context, deleteProfileOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteProfileWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIamTrustedProfileEnityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = *model.Timestamp
	modelMap["iam_id"] = *model.IamID
	modelMap["iam_id_account"] = *model.IamIDAccount
	modelMap["action"] = *model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = *model.Message
	return modelMap, nil
}
