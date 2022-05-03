// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package classicinfrastructure

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
)

func DataSourceIBMDNSDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDNSDomainRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "A domain record's internal identifier",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"name": {
				Description: "The name of the domain",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	name := d.Get("name").(string)

	names, err := service.
		Filter(filter.Build(filter.Path("domains.name").Eq(name))).
		Mask("id,name").
		GetDomains()

	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving domain: %s", err)
	}

	if len(names) == 0 {
		return fmt.Errorf("[ERROR] No domain found with name [%s]", name)
	}

	d.SetId(fmt.Sprintf("%d", *names[0].Id))
	return nil
}
