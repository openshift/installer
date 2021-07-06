// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMIAMUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIAMUsersRead,

		Schema: map[string]*schema.Schema{

			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User's Profiles",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"iam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User's IAM ID or or email of user",
						},

						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user ID used for login. ",
						},

						"realm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The realm of the user. ",
						},

						"first_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The first name of the user. ",
						},

						"last_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last name of the user. ",
						},

						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the user. Possible values are PROCESSING, PENDING, ACTIVE, DISABLED_CLASSIC_INFRASTRUCTURE, and VPN_ONLY. ",
						},

						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The email of the user. ",
						},

						"phonenumber": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The phone for the user.",
						},

						"alt_phonenumber": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alternative phone number of the user. ",
						},

						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An alphanumeric value identifying the account ID. ",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIAMUsersRead(d *schema.ResourceData, meta interface{}) error {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := userDetails.userAccount

	if err != nil {
		return err
	}

	res, err := client.ListUsers(accountID)
	if err != nil {
		return err
	}

	profileList := make([]interface{}, 0)

	for _, userInfo := range res {
		if userInfo.State == "ACTIVE" {

			user := map[string]interface{}{
				"iam_id":          userInfo.IamID,
				"user_id":         userInfo.UserID,
				"realm":           userInfo.Realm,
				"first_name":      userInfo.Firstname,
				"last_name":       userInfo.Lastname,
				"state":           userInfo.State,
				"email":           userInfo.Email,
				"phonenumber":     userInfo.Phonenumber,
				"alt_phonenumber": userInfo.Altphonenumber,
				"account_id":      userInfo.AccountID,
			}

			profileList = append(profileList, user)
		}
	}

	d.SetId(accountID)
	d.Set("users", profileList)

	return nil
}
