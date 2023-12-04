// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSccAccountLocations() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "ibm_scc_account_locations is no longer supported",
	}
}
