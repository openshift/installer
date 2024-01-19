// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSccNotificationSettings() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "ibm_scc_account_notification_settings is no longer supported",
	}
}
