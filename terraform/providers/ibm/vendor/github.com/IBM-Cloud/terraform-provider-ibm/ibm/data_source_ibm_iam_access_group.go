// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"log"

	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIBMIAMAccessGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMIAMAccessGroupRead,
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
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
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

	iamClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	start := ""
	allrecs := []iamidentityv1.ServiceID{}
	var pg int64 = 100
	for {
		listServiceIDOptions := iamidentityv1.ListServiceIdsOptions{
			AccountID: &userDetails.userAccount,
			Pagesize:  &pg,
		}
		if start != "" {
			listServiceIDOptions.Pagetoken = &start
		}

		serviceIDs, resp, err := iamClient.ListServiceIds(&listServiceIDOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing Service Ids %s %s", err, resp)
		}
		start = GetNextIAM(serviceIDs.Next)
		allrecs = append(allrecs, serviceIDs.Serviceids...)
		if start == "" {
			break
		}
	}

	listAccessGroupOption := iamAccessGroupsClient.NewListAccessGroupsOptions(accountID)
	retreivedGroups, detailedResponse, err := iamAccessGroupsClient.ListAccessGroups(listAccessGroupOption)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving access groups: %s. API Response is: %s", err, detailedResponse)
	}

	if len(retreivedGroups.Groups) == 0 {
		return fmt.Errorf("[ERROR] No access group in account")
	}
	var agName string
	var matchGroups []iamaccessgroupsv2.Group
	if v, ok := d.GetOk("access_group_name"); ok {
		agName = v.(string)
		for _, grpData := range retreivedGroups.Groups {
			if *grpData.Name == agName {
				matchGroups = append(matchGroups, grpData)
			}
		}
	} else {
		matchGroups = retreivedGroups.Groups
	}
	if len(matchGroups) == 0 {
		return fmt.Errorf("[ERROR] No Access Groups with name %s in Account", agName)
	}

	grpMap := make([]map[string]interface{}, 0, len(matchGroups))

	for _, grp := range matchGroups {
		accessGroupMembersListOptions := iamAccessGroupsClient.NewListAccessGroupMembersOptions(*grp.ID)
		members, detailedResponse, err := iamAccessGroupsClient.ListAccessGroupMembers(accessGroupMembersListOptions)
		if err != nil {
			log.Printf("Error retrieving access group members: %s.API Response: %s", err, detailedResponse)
		}

		accessGroupRulesListOptions := iamAccessGroupsClient.NewListAccessGroupRulesOptions(*grp.ID)
		rules, detailedResponse, err := iamAccessGroupsClient.ListAccessGroupRules(accessGroupRulesListOptions)
		if err != nil {
			log.Printf("Error retrieving access group rules: %s. API Response: %s", err, detailedResponse)
		}
		ibmID, serviceID := flattenMembersData(members.Members, res, allrecs)

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
