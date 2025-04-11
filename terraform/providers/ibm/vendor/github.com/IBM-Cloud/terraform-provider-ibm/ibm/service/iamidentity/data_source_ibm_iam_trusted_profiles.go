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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfilesRead,

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account ID to query for trusted profiles.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the trusted profile to query.",
			},
			"include_history": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of trusted profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
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
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the trusted profile. The name is checked for uniqueness. Therefore trusted profiles with the same names can not exist in the same account.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The optional description of the trusted profile. The 'description' property is only available if a description was provided during a create of a trusted profile.",
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
						"account_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the account that this trusted profile belong to.",
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
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
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
				},
			},
		},
	}
}

func dataSourceIBMIamTrustedProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profiles", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listProfilesOptions := &iamidentityv1.ListProfilesOptions{}

	listProfilesOptions.SetAccountID(d.Get("account_id").(string))
	if _, ok := d.GetOk("name"); ok {
		listProfilesOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("include_history"); ok {
		listProfilesOptions.SetIncludeHistory(d.Get("include_history").(bool))
	}

	trustedProfilesList, _, err := iamIdentityClient.ListProfilesWithContext(context, listProfilesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListProfilesWithContext failed: %s", err.Error()), "(Data) ibm_iam_trusted_profiles", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIamTrustedProfilesID(d))

	profiles := []map[string]interface{}{}
	for _, profilesItem := range trustedProfilesList.Profiles {
		profilesItemMap, err := DataSourceIBMIamTrustedProfilesTrustedProfileToMap(&profilesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profiles", "read", "profiles-to-map").GetDiag()
		}
		profiles = append(profiles, profilesItemMap)
	}
	if err = d.Set("profiles", profiles); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting profiles: %s", err), "(Data) ibm_iam_trusted_profiles", "read", "set-profiles").GetDiag()
	}

	return nil
}

// dataSourceIBMIamTrustedProfilesID returns a reasonable ID for the list.
func dataSourceIBMIamTrustedProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIamTrustedProfilesTrustedProfileToMap(model *iamidentityv1.TrustedProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["entity_tag"] = *model.EntityTag
	modelMap["crn"] = *model.CRN
	modelMap["name"] = *model.Name
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.ModifiedAt != nil {
		modelMap["modified_at"] = model.ModifiedAt.String()
	}
	modelMap["iam_id"] = *model.IamID
	modelMap["account_id"] = *model.AccountID
	if model.TemplateID != nil {
		modelMap["template_id"] = *model.TemplateID
	}
	if model.AssignmentID != nil {
		modelMap["assignment_id"] = *model.AssignmentID
	}
	if model.ImsAccountID != nil {
		modelMap["ims_account_id"] = flex.IntValue(model.ImsAccountID)
	}
	if model.ImsUserID != nil {
		modelMap["ims_user_id"] = flex.IntValue(model.ImsUserID)
	}
	if model.History != nil {
		history := []map[string]interface{}{}
		for _, historyItem := range model.History {
			historyItemMap, err := DataSourceIBMIamTrustedProfilesEnityHistoryRecordToMap(&historyItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			history = append(history, historyItemMap)
		}
		modelMap["history"] = history
	}
	return modelMap, nil
}

func DataSourceIBMIamTrustedProfilesEnityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = *model.Timestamp
	modelMap["iam_id"] = *model.IamID
	modelMap["iam_id_account"] = *model.IamIDAccount
	modelMap["action"] = *model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = *model.Message
	return modelMap, nil
}
