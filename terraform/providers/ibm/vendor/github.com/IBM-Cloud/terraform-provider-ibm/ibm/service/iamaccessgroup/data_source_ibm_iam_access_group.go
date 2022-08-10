// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamaccessgroup

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIAMAccessGroup() *schema.Resource {
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
						"iam_profile_ids": {
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
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return err
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := userDetails.UserAccount
	userManagement, err := meta.(conns.ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()
	res, err := client.ListUsers(accountID)
	if err != nil {
		return err
	}

	iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	start := ""
	allrecs := []iamidentityv1.ServiceID{}
	var pg int64 = 100
	for {
		listServiceIDOptions := iamidentityv1.ListServiceIdsOptions{
			AccountID: &userDetails.UserAccount,
			Pagesize:  &pg,
		}
		if start != "" {
			listServiceIDOptions.Pagetoken = &start
		}

		serviceIDs, resp, err := iamClient.ListServiceIds(&listServiceIDOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing Service Ids %s %s", err, resp)
		}
		start = flex.GetNextIAM(serviceIDs.Next)
		allrecs = append(allrecs, serviceIDs.Serviceids...)
		if start == "" {
			break
		}
	}

	profileStart := ""
	allprofiles := []iamidentityv1.TrustedProfile{}
	var plimit int64 = 100
	for {
		listProfilesOptions := iamidentityv1.ListProfilesOptions{
			AccountID: &userDetails.UserAccount,
			Pagesize:  &plimit,
		}
		if profileStart != "" {
			listProfilesOptions.Pagetoken = &profileStart
		}

		profileIDs, resp, err := iamClient.ListProfiles(&listProfilesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error listing Trusted Profiles %s %s", err, resp)
		}
		profileStart = flex.GetNextIAM(profileIDs.Next)
		allprofiles = append(allprofiles, profileIDs.Profiles...)
		if profileStart == "" {
			break
		}
	}
	offset := int64(0)
	limit := int64(100)
	listAccessGroupOption := iamAccessGroupsClient.NewListAccessGroupsOptions(accountID)
	listAccessGroupOption.Limit = &plimit
	listAccessGroupOption.Offset = &offset
	retreivedGroups, detailedResponse, err := iamAccessGroupsClient.ListAccessGroups(listAccessGroupOption)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving access groups: %s. API Response is: %s", err, detailedResponse)
	}

	if len(retreivedGroups.Groups) == 0 {
		return fmt.Errorf("[ERROR] No access group in account")
	}
	allGroups := retreivedGroups.Groups

	totalGroups := flex.IntValue(retreivedGroups.TotalCount)
	for len(allGroups) < totalGroups {
		offset = offset + limit
		listAccessGroupOption.SetOffset(offset)
		retreivedGroups, detailedResponse, err := iamAccessGroupsClient.ListAccessGroups(listAccessGroupOption)
		if err != nil {
			return fmt.Errorf("[ERROR] Error retrieving access groups: %s. API Response is: %s", err, detailedResponse)
		}

		allGroups = append(allGroups, retreivedGroups.Groups...)
	}
	var agName string
	var matchGroups []iamaccessgroupsv2.Group
	if v, ok := d.GetOk("access_group_name"); ok {
		agName = v.(string)
		for _, grpData := range allGroups {
			if *grpData.Name == agName {
				matchGroups = append(matchGroups, grpData)
			}
		}
	} else {
		matchGroups = allGroups
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
		ibmID, serviceID, profileID := flex.FlattenMembersData(members.Members, res, allrecs, allprofiles)

		grpInstance := map[string]interface{}{
			"id":              grp.ID,
			"name":            grp.Name,
			"description":     grp.Description,
			"ibm_ids":         ibmID,
			"iam_service_ids": serviceID,
			"iam_profile_ids": profileID,
			"rules":           flex.FlattenAccessGroupRules(rules),
		}

		grpMap = append(grpMap, grpInstance)

	}

	d.SetId(accountID)
	d.Set("groups", grpMap)
	d.Set("access_group_name", agName)

	return nil
}
