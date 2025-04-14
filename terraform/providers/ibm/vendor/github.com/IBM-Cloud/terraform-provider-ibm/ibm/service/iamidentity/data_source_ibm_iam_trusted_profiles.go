// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileListRead,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account ID to query for trusted profiles.",
			},
			"name": {
				Description: "Name of the profile",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profiles",
					"name"),
			},
			"include_history": {
				Description: "Defines if the entity history is included in the response. Default is false",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of trusted profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the unique identifier of the trusted profile. Example:'Profile-94497d0d-2ac3-41bf-a993-a49d1b14627c'.",
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
				},
			},
		},
	}
}

func DataSourceIBMIamTrustedProfilesValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:name"},
			Optional:                   true})

	iBMIamTrustedProfilesValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profiles", Schema: validateSchema}
	return &iBMIamTrustedProfilesValidator
}

func dataSourceIBMIamTrustedProfileListRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []iamidentityv1.TrustedProfile{}
	accountID := userDetails.UserAccount
	for {
		listProfileOptions := &iamidentityv1.ListProfilesOptions{}

		if v, ok := d.GetOk("account_id"); ok {
			listProfileOptions.SetAccountID(v.(string))
		} else {
			listProfileOptions.SetAccountID(accountID)
		}

		if v, ok := d.GetOk("name"); ok {
			listProfileOptions.SetName(v.(string))
		}

		if v, ok := d.GetOk("include_history"); ok {
			listProfileOptions.SetIncludeHistory(v.(bool))
		}

		listProfileOptions.SetPagesize(int64(100))

		if start != "" {
			listProfileOptions.Pagetoken = &start
		}

		trustedProfiles, response, err := iamIdentityClient.ListProfiles(listProfileOptions)
		if err != nil {
			log.Printf("[DEBUG] ListProfile failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListProfile failed %s\n%s", err, response))
		}
		start = flex.GetNextIAM(trustedProfiles.Next)
		allrecs = append(allrecs, trustedProfiles.Profiles...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIamTrustedProfileListID(d))

	d.Set("profiles", dataSourceTrustedProfilesListFlattenProfiles(allrecs))

	return nil
}

// dataSourceIBMIamTrustedProfileListID returns a reasonable ID for the list.
func dataSourceIBMIamTrustedProfileListID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceTrustedProfilesListFlattenProfiles(result []iamidentityv1.TrustedProfile) (profiles []map[string]interface{}) {
	for _, profilesItem := range result {
		profiles = append(profiles, dataSourceTrustedProfilesListProfilesToMap(profilesItem))
	}

	return profiles
}

func dataSourceTrustedProfilesListProfilesToMap(profilesItem iamidentityv1.TrustedProfile) (profilesMap map[string]interface{}) {
	profilesMap = map[string]interface{}{}

	if profilesItem.ID != nil {
		profilesMap["id"] = profilesItem.ID
	}
	if profilesItem.EntityTag != nil {
		profilesMap["entity_tag"] = profilesItem.EntityTag
	}
	if profilesItem.CRN != nil {
		profilesMap["crn"] = profilesItem.CRN
	}
	if profilesItem.Name != nil {
		profilesMap["name"] = profilesItem.Name
	}
	if profilesItem.Description != nil {
		profilesMap["description"] = profilesItem.Description
	}
	if profilesItem.CreatedAt != nil {
		profilesMap["created_at"] = profilesItem.CreatedAt.String()
	}
	if profilesItem.ModifiedAt != nil {
		profilesMap["modified_at"] = profilesItem.ModifiedAt.String()
	}
	if profilesItem.IamID != nil {
		profilesMap["iam_id"] = profilesItem.IamID
	}
	if profilesItem.AccountID != nil {
		profilesMap["account_id"] = profilesItem.AccountID
	}
	if profilesItem.ImsAccountID != nil {
		profilesMap["ims_account_id"] = profilesItem.ImsAccountID
	}
	if profilesItem.ImsUserID != nil {
		profilesMap["ims_user_id"] = profilesItem.ImsUserID
	}
	if profilesItem.History != nil {
		historyList := []map[string]interface{}{}
		for _, historyItem := range profilesItem.History {
			historyList = append(historyList, dataSourceTrustedProfilesListProfilesHistoryToMap(historyItem))
		}
		profilesMap["history"] = historyList
	}

	return profilesMap
}

func dataSourceTrustedProfilesListProfilesHistoryToMap(historyItem iamidentityv1.EnityHistoryRecord) (historyMap map[string]interface{}) {
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
