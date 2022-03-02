// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIBMIAMAccessGroupMembers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIAMAccessGroupMembersCreate,
		ReadContext:   resourceIBMIAMAccessGroupMembersRead,
		UpdateContext: resourceIBMIAMAccessGroupMembersUpdate,
		DeleteContext: resourceIBMIAMAccessGroupMembersDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the access group",
				ForceNew:    true,
			},

			"ibm_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      resourceIBMVPCHash,
			},

			"iam_service_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIBMIAMAccessGroupMembersCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	grpID := d.Get("access_group_id").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := userDetails.userAccount

	var userids, serviceids []string

	users := expandStringList(d.Get("ibm_ids").(*schema.Set).List())
	services := expandStringList(d.Get("iam_service_ids").(*schema.Set).List())

	if len(users) == 0 && len(services) == 0 {
		return diag.FromErr(fmt.Errorf("ERROR] Provide either `ibm_ids` or `iam_service_ids`"))

	}

	userids, err = flattenUserIds(accountID, users, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceids, err = flattenServiceIds(services, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	members := prepareMemberAddRequest(iamAccessGroupsClient, userids, serviceids)

	addMembersToAccessGroupOptions := iamAccessGroupsClient.NewAddMembersToAccessGroupOptions(grpID)
	addMembersToAccessGroupOptions.SetMembers(members)
	membership, detailResponse, err := iamAccessGroupsClient.AddMembersToAccessGroup(addMembersToAccessGroupOptions)
	if err != nil || membership == nil {
		return diag.FromErr(fmt.Errorf("Error adding members to group(%s). API response: %s", grpID, detailResponse))
	}

	d.SetId(fmt.Sprintf("%s/%s", grpID, time.Now().UTC().String()))

	return resourceIBMIAMAccessGroupMembersRead(context, d, meta)
}

func resourceIBMIAMAccessGroupMembersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	grpID := parts[0]
	listAccessGroupMembersOptions := iamAccessGroupsClient.NewListAccessGroupMembersOptions(grpID)
	offset := int64(0)
	// lets fetch 100 in a single pagination
	limit := int64(100)
	listAccessGroupMembersOptions.SetLimit(limit)
	members, detailedResponse, err := iamAccessGroupsClient.ListAccessGroupMembers(listAccessGroupMembersOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving access group members: %s. API Response: %s", err, detailedResponse))
	}
	allMembers := members.Members
	totalMembers := intValue(members.TotalCount)
	for len(allMembers) < totalMembers {
		offset = offset + limit
		listAccessGroupMembersOptions.SetOffset(offset)
		members, detailedResponse, err = iamAccessGroupsClient.ListAccessGroupMembers(listAccessGroupMembersOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving access group members: %s. API Response: %s", err, detailedResponse))
		}
		allMembers = append(allMembers, members.Members...)
	}

	d.Set("access_group_id", grpID)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := userDetails.userAccount

	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	client := userManagement.UserInvite()
	res, err := client.ListUsers(accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	iamClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
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
			return diag.FromErr(fmt.Errorf("[ERROR] Error listing Service Ids %s %s", err, resp))
		}
		start = GetNextIAM(serviceIDs.Next)
		allrecs = append(allrecs, serviceIDs.Serviceids...)
		if start == "" {
			break
		}
	}

	d.Set("members", flattenAccessGroupMembers(allMembers, res, allrecs))
	ibmID, serviceID := flattenMembersData(allMembers, res, allrecs)
	if len(ibmID) > 0 {
		d.Set("ibm_ids", ibmID)
	}
	if len(serviceID) > 0 {
		d.Set("iam_service_ids", serviceID)
	}
	return nil
}

func resourceIBMIAMAccessGroupMembersUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	grpID := parts[0]

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := userDetails.userAccount

	var removeUsers, addUsers, removeServiceids, addServiceids []string
	o, n := d.GetChange("ibm_ids")
	ou := o.(*schema.Set)
	nu := n.(*schema.Set)

	removeUsers = expandStringList(ou.Difference(nu).List())
	addUsers = expandStringList(nu.Difference(ou).List())

	os, ns := d.GetChange("iam_service_ids")
	osi := os.(*schema.Set)
	nsi := ns.(*schema.Set)

	removeServiceids = expandStringList(osi.Difference(nsi).List())
	addServiceids = expandStringList(nsi.Difference(osi).List())

	if len(addUsers) > 0 || len(addServiceids) > 0 && !d.IsNewResource() {
		var userids, serviceids []string
		userids, err = flattenUserIds(accountID, addUsers, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		serviceids, err = flattenServiceIds(addServiceids, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		members := prepareMemberAddRequest(iamAccessGroupsClient, userids, serviceids)

		addMembersToAccessGroupOptions := iamAccessGroupsClient.NewAddMembersToAccessGroupOptions(grpID)
		addMembersToAccessGroupOptions.SetMembers(members)
		membership, detailResponse, err := iamAccessGroupsClient.AddMembersToAccessGroup(addMembersToAccessGroupOptions)
		if err != nil || membership == nil {
			return diag.FromErr(fmt.Errorf("Error updating members to group(%s). API response: %s", grpID, detailResponse))
		}

	}
	if len(removeUsers) > 0 || len(removeServiceids) > 0 && !d.IsNewResource() {
		iamClient, err := meta.(ClientSession).IAMIdentityV1API()
		if err != nil {
			return diag.FromErr(err)
		}
		for _, u := range removeUsers {
			ibmUniqueId, err := getIBMUniqueId(accountID, u, meta)
			if err != nil {
				return diag.FromErr(err)
			}
			removeMembersFromAccessGroupOptions := iamAccessGroupsClient.NewRemoveMemberFromAccessGroupOptions(grpID, ibmUniqueId)
			_, err = iamAccessGroupsClient.RemoveMemberFromAccessGroup(removeMembersFromAccessGroupOptions)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, s := range removeServiceids {
			getServiceIDOptions := iamidentityv1.GetServiceIDOptions{
				ID: &s,
			}
			serviceID, resp, err := iamClient.GetServiceID(&getServiceIDOptions)
			if err != nil || serviceID == nil {
				return diag.FromErr(fmt.Errorf("ERROR] Error Getting Service Ids %s %s", err, resp))
			}
			removeMembersFromAccessGroupOptions := iamAccessGroupsClient.NewRemoveMemberFromAccessGroupOptions(grpID, *serviceID.IamID)
			detailResponse, err := iamAccessGroupsClient.RemoveMemberFromAccessGroup(removeMembersFromAccessGroupOptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("Error removing members to group(%s). API Response: %s", grpID, detailResponse))
			}

		}
	}

	return resourceIBMIAMAccessGroupMembersRead(context, d, meta)

}

func resourceIBMIAMAccessGroupMembersDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	grpID := parts[0]

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	users := expandStringList(d.Get("ibm_ids").(*schema.Set).List())

	for _, name := range users {

		ibmUniqueID, err := getIBMUniqueId(userDetails.userAccount, name, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		removeMembersFromAccessGroupOptions := iamAccessGroupsClient.NewRemoveMemberFromAccessGroupOptions(grpID, ibmUniqueID)
		_, err = iamAccessGroupsClient.RemoveMemberFromAccessGroup(removeMembersFromAccessGroupOptions)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	services := expandStringList(d.Get("iam_service_ids").(*schema.Set).List())

	for _, id := range services {
		serviceID, err := getServiceID(id, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		removeMembersFromAccessGroupOptions := &iamaccessgroupsv2.RemoveMemberFromAccessGroupOptions{
			AccessGroupID: &grpID,
			IamID:         serviceID.IamID,
		}
		_, err = iamAccessGroupsClient.RemoveMemberFromAccessGroup(removeMembersFromAccessGroupOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return nil
}

func prepareMemberAddRequest(iamAccessGroupsClient *iamaccessgroupsv2.IamAccessGroupsV2, userIds, serviceIds []string) (members []iamaccessgroupsv2.AddGroupMembersRequestMembersItem) {
	members = make([]iamaccessgroupsv2.AddGroupMembersRequestMembersItem, len(userIds)+len(serviceIds))
	var i = 0
	userType := "user"
	serviceType := "service"
	for _, id := range userIds {
		membersItem, err := iamAccessGroupsClient.NewAddGroupMembersRequestMembersItem(id, userType)
		if err != nil {
			log.Printf("Error in preparing membership data. %s", err)
		}
		members[i] = *membersItem
		i++
	}

	for _, id := range serviceIds {
		membersItem, err := iamAccessGroupsClient.NewAddGroupMembersRequestMembersItem(id, serviceType)
		if err != nil || membersItem == nil {
			log.Printf("Error in preparing membership data. %s", err)
		}
		members[i] = *membersItem
		i++
	}
	return
}

func getServiceID(id string, meta interface{}) (iamidentityv1.ServiceID, error) {
	serviceids := iamidentityv1.ServiceID{}
	iamClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return serviceids, err
	}
	getServiceIDOptions := iamidentityv1.GetServiceIDOptions{
		ID: &id,
	}
	serviceID, resp, err := iamClient.GetServiceID(&getServiceIDOptions)
	if err != nil || serviceID == nil {
		return serviceids, fmt.Errorf("ERROR] Error Getting Service Ids %s %s", err, resp)
	}
	return *serviceID, nil
}
