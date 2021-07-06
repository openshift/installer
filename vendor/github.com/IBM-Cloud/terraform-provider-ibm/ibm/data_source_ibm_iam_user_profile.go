// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMIAMUserProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIAMUserProfileRead,

		Schema: map[string]*schema.Schema{

			"iam_id": {
				Description: "User's IAM ID or or email of user",
				Type:        schema.TypeString,
				Required:    true,
			},

			"allowed_ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of allowed IPv4 or IPv6 addresses ",
			},

			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID used for login. ",
			},

			"firstname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The first name of the user. ",
			},

			"lastname": {
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

			"altphonenumber": {
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
	}
}

func dataSourceIBMIAMUserProfileRead(d *schema.ResourceData, meta interface{}) error {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()

	userEmail := d.Get("iam_id").(string)

	accountID, err := getUserAccountID(d, meta)
	if err != nil {
		return err
	}

	iamID, err := getIBMUniqueId(accountID, userEmail, meta)
	if err != nil {
		return err
	}

	userInfo, error := client.GetUserProfile(accountID, iamID)
	if error != nil {
		return error
	}

	d.Set("user_id", userInfo.UserID)
	d.Set("firstname", userInfo.Firstname)
	d.Set("lastname", userInfo.Lastname)
	d.Set("state", userInfo.State)
	d.Set("email", userInfo.Email)
	d.Set("phonenumber", userInfo.Phonenumber)
	d.Set("altphonenumber", userInfo.Altphonenumber)
	d.Set("account_id", userInfo.AccountID)

	UserSettings, UserSettingError := client.GetUserSettings(accountID, iamID)
	if UserSettingError != nil {
		return UserSettingError
	}

	iplist := strings.Split(UserSettings.AllowedIPAddresses, ",")
	d.Set("allowed_ip_addresses", iplist)
	d.SetId(userEmail)

	return nil
}
