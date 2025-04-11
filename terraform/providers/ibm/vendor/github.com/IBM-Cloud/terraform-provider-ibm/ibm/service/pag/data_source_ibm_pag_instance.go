// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package pag

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/resourcecontroller"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMPag() *schema.Resource {
	riSchema := resourcecontroller.DataSourceIBMResourceInstance().Schema

	return &schema.Resource{
		Read:   resourcecontroller.DataSourceIBMResourceInstanceRead,
		Schema: riSchema,
	}
}
