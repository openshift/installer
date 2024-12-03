// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cloudfoundry

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppDomainShared() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMAppDomainSharedRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "The name of the shared domain",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateDomainName,
			},
		},
		DeprecationMessage: "This service is deprecated.",
	}
}

func dataSourceIBMAppDomainSharedRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(conns.ClientSession).MccpAPI()
	if err != nil {
		return err
	}
	domainName := d.Get("name").(string)
	shdomain, err := cfClient.SharedDomains().FindByName(domainName)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving shared domain: %s", err)
	}
	d.SetId(shdomain.GUID)
	return nil

}
