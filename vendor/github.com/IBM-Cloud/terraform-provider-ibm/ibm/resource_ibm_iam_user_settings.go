// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strings"

	v2 "github.com/IBM-Cloud/bluemix-go/api/usermanagement/usermanagementv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	iamUserSettingIamID              = "iam_id"
	iamUserSettingAllowedIPAddresses = "allowed_ip_addresses"
)

func resourceIBMUserSettings() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMUserSettingsCreate,
		Read:     resourceIBMIAMUserSettingsRead,
		Update:   resourceIBMIAMUserSettingsUpdate,
		Delete:   resourceIBMIAMUserSettingsDelete,
		Exists:   resourceIBMIAMUserSettingsExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			iamUserSettingIamID: {
				Description: "User's IAM ID or or email of user",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			iamUserSettingAllowedIPAddresses: {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of allowed IPv4 or IPv6 addresses ",
			},
		},
	}
}

func resourceIBMIAMUserSettingsCreate(d *schema.ResourceData, meta interface{}) error {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()

	userEmail := d.Get(iamUserSettingIamID).(string)

	//Read from Bluemix UserConfig
	accountID, err := getUserAccountID(d, meta)
	if err != nil {
		return err
	}

	iamID, err := getIBMUniqueId(accountID, userEmail, meta)
	if err != nil {
		return err
	}

	UserSettingsPayload := v2.UserSettingOptions{}

	if ip, ok := d.GetOk(iamUserSettingAllowedIPAddresses); ok && ip != nil {
		var ips = make([]string, 0)
		for _, i := range ip.([]interface{}) {
			ips = append(ips, i.(string))
		}
		ipStr := strings.Join(ips, ",")
		UserSettingsPayload.AllowedIPAddresses = ipStr
	}

	_, UserSettingError := client.ManageUserSettings(accountID, iamID, UserSettingsPayload)

	if UserSettingError != nil && !strings.Contains(UserSettingError.Error(), "EmptyResponseBody") {
		return fmt.Errorf("Error occured during user settings: %s", UserSettingError)
	}

	d.SetId(userEmail)

	return resourceIBMIAMUserSettingsRead(d, meta)
}

func getUserAccountID(d *schema.ResourceData, meta interface{}) (string, error) {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return "", err
	}
	return userDetails.userAccount, nil
}

func resourceIBMIAMUserSettingsRead(d *schema.ResourceData, meta interface{}) error {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()

	accountID, err := getUserAccountID(d, meta)
	if err != nil {
		return err
	}

	iamID, err := getIBMUniqueId(accountID, d.Id(), meta)
	if err != nil {
		return err
	}

	UserSettings, UserSettingError := client.GetUserSettings(accountID, iamID)
	if UserSettingError != nil {
		return UserSettingError
	}

	iplist := strings.Split(UserSettings.AllowedIPAddresses, ",")
	d.Set(iamUserSettingAllowedIPAddresses, iplist)

	return nil

}

func resourceIBMIAMUserSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	// validate change
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()

	accountID, err := getUserAccountID(d, meta)
	if err != nil {
		return err
	}

	iamID, err := getIBMUniqueId(accountID, d.Id(), meta)
	if err != nil {
		return err
	}

	hasChanged := false

	userSettingPayload := v2.UserSettingOptions{}

	if d.HasChange(iamUserSettingAllowedIPAddresses) {
		if ip, ok := d.GetOk(iamUserSettingAllowedIPAddresses); ok && ip != nil {
			var ips = make([]string, 0)
			for _, i := range ip.([]interface{}) {
				ips = append(ips, i.(string))
			}
			ipStr := strings.Join(ips, ",")
			userSettingPayload.AllowedIPAddresses = ipStr
		}
		hasChanged = true
	}

	if hasChanged {
		_, UserSettingError := client.ManageUserSettings(accountID, iamID, userSettingPayload)
		if UserSettingError != nil && !strings.Contains(UserSettingError.Error(), "EmptyResponseBody") {
			return fmt.Errorf("Error occured during user settings: %s", UserSettingError)
		}
	}

	return resourceIBMIAMUserSettingsRead(d, meta)
}

func resourceIBMIAMUserSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return err
	}
	client := userManagement.UserInvite()

	accountID, err := getUserAccountID(d, meta)
	if err != nil {
		return err
	}

	iamID, err := getIBMUniqueId(accountID, d.Id(), meta)
	if err != nil {
		return err
	}

	userSettingPayload := v2.UserSettingOptions{}

	_, UserSettingError := client.ManageUserSettings(accountID, iamID, userSettingPayload)
	if UserSettingError != nil && !strings.Contains(UserSettingError.Error(), "EmptyResponseBody") {
		return fmt.Errorf("Error occured during user settings: %s", UserSettingError)
	}

	return nil
}

func resourceIBMIAMUserSettingsExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return false, err
	}
	client := userManagement.UserInvite()

	accountID, err := getUserAccountID(d, meta)
	if err != nil {
		return false, err
	}

	iamID, err := getIBMUniqueId(accountID, d.Id(), meta)
	if err != nil {
		return false, err
	}

	_, settingErr := client.GetUserSettings(accountID, iamID)

	if settingErr != nil {
		return false, settingErr
	}
	return true, nil
}
