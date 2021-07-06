// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"log"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMIAMAccessGroup() *schema.Resource {
	return &schema.Resource{
		Read:     dataIBMIAMAccessGroupRead,
		Exists:   resourceIBMIAMAccessGroupExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the access group",
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the access group",
						},
						"ibm_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"iam_service_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the Rule",
									},
									"expiration": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The expiration in hours",
									},
									"identity_provider": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The realm name or identity proivider url",
									},
									"conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"claim": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"rule_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "id of the rule",
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

func dataIBMIAMAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := userDetails.userAccount
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()
	res, err := client.ListUsers(accountID)
	if err != nil {
		return err
	}

	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return err
	}

	boundTo := crn.New(userDetails.cloudName, userDetails.cloudType)
	boundTo.ScopeType = crn.ScopeAccount
	boundTo.Scope = userDetails.userAccount

	serviceIDs, err := iamClient.ServiceIds().List(boundTo.String())
	if err != nil {
		return err
	}

	retreivedGroups, err := iamuumClient.AccessGroup().List(accountID)
	if err != nil {
		return fmt.Errorf("Error retrieving access groups: %s", err)
	}

	if len(retreivedGroups) == 0 {
		return fmt.Errorf("No access group in account")
	}
	var agName string
	var matchGroups []models.AccessGroupV2
	if v, ok := d.GetOk("access_group_name"); ok {
		agName = v.(string)
		for _, grpData := range retreivedGroups {
			if grpData.Name == agName {
				matchGroups = append(matchGroups, grpData)
			}
		}
	} else {
		matchGroups = retreivedGroups
	}
	if len(matchGroups) == 0 {
		return fmt.Errorf("No Access Groups with name %s in Account", agName)
	}

	grpMap := make([]map[string]interface{}, 0, len(matchGroups))

	for _, grp := range matchGroups {
		members, err := iamuumClient.AccessGroupMember().List(grp.ID)
		if err != nil {
			log.Println("Error retrieving access group members: ", err)
		}
		rules, err := iamuumClient.DynamicRule().List(grp.ID)
		if err != nil {
			log.Println("Error retrieving access group rules: ", err)
		}
		ibmID, serviceID := flattenMembersData(members, res, serviceIDs)

		grpInstance := map[string]interface{}{
			"id":              grp.ID,
			"name":            grp.Name,
			"description":     grp.Description,
			"ibm_ids":         ibmID,
			"iam_service_ids": serviceID,
			"rules":           flattenAccessGroupRules(rules),
		}

		grpMap = append(grpMap, grpInstance)

	}

	d.SetId(accountID)
	d.Set("groups", grpMap)
	d.Set("access_group_name", agName)

	return nil
}
