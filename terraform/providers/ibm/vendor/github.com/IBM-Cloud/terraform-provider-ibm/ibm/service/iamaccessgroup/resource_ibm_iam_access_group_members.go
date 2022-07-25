// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamaccessgroup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMAccessGroupMembers() *schema.Resource {
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
				Set:      flex.ResourceIBMVPCHash,
			},

			"iam_service_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"iam_profile_ids": {
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
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	grpID := d.Get("access_group_id").(string)

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := userDetails.UserAccount

	var userids, serviceids, profileids []string

	users := flex.ExpandStringList(d.Get("ibm_ids").(*schema.Set).List())
	services := flex.ExpandStringList(d.Get("iam_service_ids").(*schema.Set).List())
	profiles := flex.ExpandStringList(d.Get("iam_profile_ids").(*schema.Set).List())

	if len(users) == 0 && len(services) == 0 && len(profiles) == 0 {
		return diag.FromErr(fmt.Errorf("ERROR] Provide either `ibm_ids` or `iam_service_ids` or `iam_profile_ids`"))

	}

	userids, err = flex.FlattenUserIds(accountID, users, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceids, err = FlattenServiceIds(services, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	profileids, err = FlattenProfileIds(profiles, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	members := prepareMemberAddRequest(iamAccessGroupsClient, userids, serviceids, profileids)

	addMembersToAccessGroupOptions := iamAccessGroupsClient.NewAddMembersToAccessGroupOptions(grpID)
	addMembersToAccessGroupOptions.SetMembers(members)
	membership, detailResponse, err := iamAccessGroupsClient.AddMembersToAccessGroup(addMembersToAccessGroupOptions)
	if err != nil || membership == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error adding members to group(%s). API response: %s", grpID, detailResponse))
	}

	d.SetId(fmt.Sprintf("%s/%s", grpID, time.Now().UTC().String()))

	return resourceIBMIAMAccessGroupMembersRead(context, d, meta)
}

func resourceIBMIAMAccessGroupMembersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
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
	totalMembers := flex.IntValue(members.TotalCount)
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

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := userDetails.UserAccount

	userManagement, err := meta.(conns.ClientSession).UserManagementAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	client := userManagement.UserInvite()
	res, err := client.ListUsers(accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
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
			return diag.FromErr(fmt.Errorf("[ERROR] Error listing Service Ids %s %s", err, resp))
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
			return diag.FromErr(fmt.Errorf("[ERROR] Error listing Trusted Profiles %s %s", err, resp))
		}
		profileStart = flex.GetNextIAM(profileIDs.Next)
		allprofiles = append(allprofiles, profileIDs.Profiles...)
		if profileStart == "" {
			break
		}
	}

	d.Set("members", flex.FlattenAccessGroupMembers(allMembers, res, allrecs))
	ibmID, serviceID, profileID := flex.FlattenMembersData(allMembers, res, allrecs, allprofiles)
	if len(ibmID) > 0 {
		d.Set("ibm_ids", ibmID)
	}
	if len(serviceID) > 0 {
		d.Set("iam_service_ids", serviceID)
	}
	if len(profileID) > 0 {
		d.Set("iam_profile_ids", profileID)
	}
	return nil
}

func resourceIBMIAMAccessGroupMembersUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	grpID := parts[0]

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	accountID := userDetails.UserAccount

	var removeUsers, addUsers, removeServiceids, addServiceids, removeProfileids, addProfileids []string
	o, n := d.GetChange("ibm_ids")
	ou := o.(*schema.Set)
	nu := n.(*schema.Set)

	removeUsers = flex.ExpandStringList(ou.Difference(nu).List())
	addUsers = flex.ExpandStringList(nu.Difference(ou).List())

	os, ns := d.GetChange("iam_service_ids")
	osi := os.(*schema.Set)
	nsi := ns.(*schema.Set)

	removeServiceids = flex.ExpandStringList(osi.Difference(nsi).List())
	addServiceids = flex.ExpandStringList(nsi.Difference(osi).List())

	op, np := d.GetChange("iam_profile_ids")
	opi := op.(*schema.Set)
	npi := np.(*schema.Set)

	removeProfileids = flex.ExpandStringList(opi.Difference(npi).List())
	addProfileids = flex.ExpandStringList(npi.Difference(opi).List())

	if len(addUsers) > 0 || len(addServiceids) > 0 || len(addProfileids) > 0 && !d.IsNewResource() {
		var userids, serviceids, profileids []string
		userids, err = flex.FlattenUserIds(accountID, addUsers, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		serviceids, err = FlattenServiceIds(addServiceids, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		profileids, err = FlattenProfileIds(addProfileids, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		members := prepareMemberAddRequest(iamAccessGroupsClient, userids, serviceids, profileids)

		addMembersToAccessGroupOptions := iamAccessGroupsClient.NewAddMembersToAccessGroupOptions(grpID)
		addMembersToAccessGroupOptions.SetMembers(members)
		membership, detailResponse, err := iamAccessGroupsClient.AddMembersToAccessGroup(addMembersToAccessGroupOptions)
		if err != nil || membership == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating members to group(%s). API response: %s", grpID, detailResponse))
		}

	}
	if len(removeUsers) > 0 || len(removeServiceids) > 0 || len(removeProfileids) > 0 && !d.IsNewResource() {
		iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
		if err != nil {
			return diag.FromErr(err)
		}
		for _, u := range removeUsers {
			ibmUniqueId, err := flex.GetIBMUniqueId(accountID, u, meta)
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
				return diag.FromErr(fmt.Errorf("[ERROR] Error removing members to group(%s). API Response: %s", grpID, detailResponse))
			}

		}

		for _, p := range removeProfileids {
			getProfileOptions := iamidentityv1.GetProfileOptions{
				ProfileID: &p,
			}
			profileID, resp, err := iamClient.GetProfile(&getProfileOptions)
			if err != nil || profileID == nil {
				return diag.FromErr(fmt.Errorf("ERROR] Error Getting Profile Ids %s %s", err, resp))
			}
			removeMembersFromAccessGroupOptions := iamAccessGroupsClient.NewRemoveMemberFromAccessGroupOptions(grpID, *profileID.IamID)
			detailResponse, err := iamAccessGroupsClient.RemoveMemberFromAccessGroup(removeMembersFromAccessGroupOptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error removing members to group(%s). API Response: %s", grpID, detailResponse))
			}

		}
	}

	return resourceIBMIAMAccessGroupMembersRead(context, d, meta)

}

func resourceIBMIAMAccessGroupMembersDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	grpID := parts[0]

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	users := flex.ExpandStringList(d.Get("ibm_ids").(*schema.Set).List())

	for _, name := range users {

		ibmUniqueID, err := flex.GetIBMUniqueId(userDetails.UserAccount, name, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		removeMembersFromAccessGroupOptions := iamAccessGroupsClient.NewRemoveMemberFromAccessGroupOptions(grpID, ibmUniqueID)
		_, err = iamAccessGroupsClient.RemoveMemberFromAccessGroup(removeMembersFromAccessGroupOptions)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	services := flex.ExpandStringList(d.Get("iam_service_ids").(*schema.Set).List())

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

	profiles := flex.ExpandStringList(d.Get("iam_profile_ids").(*schema.Set).List())

	for _, id := range profiles {
		profileID, err := getProfileID(id, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		removeMembersFromAccessGroupOptions := &iamaccessgroupsv2.RemoveMemberFromAccessGroupOptions{
			AccessGroupID: &grpID,
			IamID:         profileID.IamID,
		}
		_, err = iamAccessGroupsClient.RemoveMemberFromAccessGroup(removeMembersFromAccessGroupOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return nil
}

func prepareMemberAddRequest(iamAccessGroupsClient *iamaccessgroupsv2.IamAccessGroupsV2, userIds, serviceIds, profileIds []string) (members []iamaccessgroupsv2.AddGroupMembersRequestMembersItem) {
	members = make([]iamaccessgroupsv2.AddGroupMembersRequestMembersItem, len(userIds)+len(serviceIds)+len(profileIds))
	var i = 0
	userType := "user"
	serviceType := "service"
	profileType := "profile"

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

	for _, id := range profileIds {
		membersItem, err := iamAccessGroupsClient.NewAddGroupMembersRequestMembersItem(id, profileType)
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
	iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
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

func FlattenServiceIds(services []string, meta interface{}) ([]string, error) {
	serviceids := make([]string, len(services))
	for i, id := range services {
		serviceID, err := getServiceID(id, meta)
		if err != nil {
			return nil, err
		}
		serviceids[i] = *serviceID.IamID
	}
	return serviceids, nil
}

func FlattenProfileIds(profiles []string, meta interface{}) ([]string, error) {
	profileids := make([]string, len(profiles))
	for i, id := range profiles {
		profileID, err := getProfileID(id, meta)
		if err != nil {
			return nil, err
		}
		profileids[i] = *profileID.IamID
	}
	return profileids, nil
}

func getProfileID(id string, meta interface{}) (iamidentityv1.TrustedProfile, error) {
	profileids := iamidentityv1.TrustedProfile{}
	iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return profileids, err
	}
	getProfileOptions := iamidentityv1.GetProfileOptions{
		ProfileID: &id,
	}
	profileID, resp, err := iamClient.GetProfile(&getProfileOptions)
	if err != nil || profileID == nil {
		return profileids, fmt.Errorf("ERROR] Error Getting Profile Ids %s %s", err, resp)
	}
	return *profileID, nil
}
