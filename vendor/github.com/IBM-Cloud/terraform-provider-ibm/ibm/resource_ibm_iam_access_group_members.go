// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/iamuum/iamuumv2"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMIAMAccessGroupMembers() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMAccessGroupMembersCreate,
		Read:     resourceIBMIAMAccessGroupMembersRead,
		Update:   resourceIBMIAMAccessGroupMembersUpdate,
		Delete:   resourceIBMIAMAccessGroupMembersDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the access group",
			},

			"ibm_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceIBMIAMAccessGroupMembersCreate(d *schema.ResourceData, meta interface{}) error {
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}

	grpID := d.Get("access_group_id").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := userDetails.userAccount

	var userids, serviceids []string

	users := expandStringList(d.Get("ibm_ids").(*schema.Set).List())
	services := expandStringList(d.Get("iam_service_ids").(*schema.Set).List())

	if len(users) == 0 && len(services) == 0 {
		return fmt.Errorf("ERROR] Provide either `ibm_ids` or `iam_service_ids`")

	}

	userids, err = flattenUserIds(accountID, users, meta)
	if err != nil {
		return err
	}

	serviceids, err = flattenServiceIds(services, meta)
	if err != nil {
		return err
	}

	request := prepareMemberAddRequest(userids, serviceids)

	_, err = iamuumClient.AccessGroupMember().Add(grpID, request)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", grpID, time.Now().UTC().String()))

	return resourceIBMIAMAccessGroupMembersRead(d, meta)
}

func resourceIBMIAMAccessGroupMembersRead(d *schema.ResourceData, meta interface{}) error {
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	grpID := parts[0]

	members, err := iamuumClient.AccessGroupMember().List(grpID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving access group members: %s", err)
	}

	d.Set("access_group_id", grpID)

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

	d.Set("members", flattenAccessGroupMembers(members, res, allrecs))

	return nil
}

func resourceIBMIAMAccessGroupMembersUpdate(d *schema.ResourceData, meta interface{}) error {

	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	grpID := parts[0]

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
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
			return err
		}

		serviceids, err = flattenServiceIds(addServiceids, meta)
		if err != nil {
			return err
		}
		request := prepareMemberAddRequest(userids, serviceids)

		_, err = iamuumClient.AccessGroupMember().Add(grpID, request)
		if err != nil {
			return err
		}

	}
	if len(removeUsers) > 0 || len(removeServiceids) > 0 && !d.IsNewResource() {
		iamClient, err := meta.(ClientSession).IAMIdentityV1API()
		if err != nil {
			return err
		}
		for _, u := range removeUsers {
			ibmUniqueId, err := getIBMUniqueId(accountID, u, meta)
			if err != nil {
				return err
			}
			err = iamuumClient.AccessGroupMember().Remove(grpID, ibmUniqueId)
			if err != nil {
				return err
			}

		}

		for _, s := range removeServiceids {
			getServiceIDOptions := iamidentityv1.GetServiceIDOptions{
				ID: &s,
			}
			serviceID, resp, err := iamClient.GetServiceID(&getServiceIDOptions)
			if err != nil {
				return fmt.Errorf("ERROR] Error Getting Service Ids %s %s", err, resp)
			}
			err = iamuumClient.AccessGroupMember().Remove(grpID, *serviceID.IamID)
			if err != nil {
				return err
			}

		}
	}

	return resourceIBMIAMAccessGroupMembersRead(d, meta)

}

func resourceIBMIAMAccessGroupMembersDelete(d *schema.ResourceData, meta interface{}) error {
	iamuumClient, err := meta.(ClientSession).IAMUUMAPIV2()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	grpID := parts[0]

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	users := expandStringList(d.Get("ibm_ids").(*schema.Set).List())

	for _, name := range users {

		ibmUniqueID, err := getIBMUniqueId(userDetails.userAccount, name, meta)
		if err != nil {
			return err
		}
		err = iamuumClient.AccessGroupMember().Remove(grpID, ibmUniqueID)
		if err != nil {
			return err
		}

	}

	services := expandStringList(d.Get("iam_service_ids").(*schema.Set).List())

	for _, id := range services {
		serviceID, err := getServiceID(id, meta)
		if err != nil {
			return err
		}
		err = iamuumClient.AccessGroupMember().Remove(grpID, *serviceID.IamID)
		if err != nil {
			return err
		}
	}

	d.SetId("")

	return nil
}

func prepareMemberAddRequest(userIds, serviceIds []string) (req iamuumv2.AddGroupMemberRequestV2) {
	req.Members = make([]models.AccessGroupMemberV2, len(userIds)+len(serviceIds))
	var i = 0
	for _, id := range userIds {
		req.Members[i] = models.AccessGroupMemberV2{
			ID:   id,
			Type: iamuumv2.AccessGroupMemberUser,
		}
		i++
	}

	for _, id := range serviceIds {
		req.Members[i] = models.AccessGroupMemberV2{
			ID:   id,
			Type: iamuumv2.AccessGroupMemberService,
		}
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
	if err != nil {
		return serviceids, fmt.Errorf("ERROR] Error Getting Service Ids %s %s", err, resp)
	}
	return *serviceID, nil
}
