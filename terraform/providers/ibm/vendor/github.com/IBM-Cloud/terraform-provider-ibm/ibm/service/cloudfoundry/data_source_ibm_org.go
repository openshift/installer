// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudfoundry

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMOrg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMOrgRead,

		Schema: map[string]*schema.Schema{
			"org": {
				Description:  "Org name, for example myorg@domain",
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "use name instead",
				ExactlyOneOf: []string{"org", "name"},
			},
			"name": {
				Description:  "Org name, for example myorg@domain",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"org", "name"},
			},
		},
		DeprecationMessage: "This service is deprecated.",
	}
}

func dataSourceIBMOrgRead(d *schema.ResourceData, meta interface{}) error {
	cfAPI, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	orgAPI := cfAPI.Organizations()
	var org string
	if v, ok := d.GetOk("name"); ok {
		org = v.(string)
	}
	if v, ok := d.GetOk("org"); ok {
		org = v.(string)
	}

	orgFields, err := orgAPI.FindByName(org, conns.BluemixRegion)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving organisation: %s", err)
	}
	d.SetId(orgFields.GUID)

	return nil
}
