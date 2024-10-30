// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudfoundry

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppDomainPrivate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMAppDomainPrivateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the private domain",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		DeprecationMessage: "This service is deprecated.",
	}
}

func dataSourceIBMAppDomainPrivateRead(d *schema.ResourceData, meta interface{}) error {
	cfAPI, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	domainName := d.Get("name").(string)
	prdomain, err := cfAPI.PrivateDomains().FindByName(domainName)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving domain: %s", err)
	}
	d.SetId(prdomain.GUID)
	return nil

}
