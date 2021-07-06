// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMAccountRead,

		Schema: map[string]*schema.Schema{
			"org_guid": {
				Description: "The guid of the org",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAccountRead(d *schema.ResourceData, meta interface{}) error {
	bmxSess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	accClient, err := meta.(ClientSession).BluemixAcccountAPI()
	if err != nil {
		return err
	}
	orgGUID := d.Get("org_guid").(string)
	account, err := accClient.Accounts().FindByOrg(orgGUID, bmxSess.Config.Region)
	if err != nil {
		return fmt.Errorf("Error retrieving organisation: %s", err)
	}

	accountv1Client, err := meta.(ClientSession).BluemixAcccountv1API()
	if err != nil {
		return err
	}
	accountUsers, err := accountv1Client.Accounts().GetAccountUsers(account.GUID)
	if err != nil {
		return fmt.Errorf("Error retrieving users in account: %s", err)
	}
	accountUsersMap := make([]map[string]string, 0, len(accountUsers))
	for _, user := range accountUsers {
		accountUser := make(map[string]string)
		accountUser["id"] = user.Id
		accountUser["email"] = user.Email
		accountUser["state"] = user.State
		accountUser["role"] = user.Role
		accountUsersMap = append(accountUsersMap, accountUser)
	}

	d.SetId(account.GUID)
	d.Set("account_users", accountUsersMap)
	return nil
}
