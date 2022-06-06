// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package classicinfrastructure

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
)

func DataSourceIBMDNSDomainRegistration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDNSDomainRegistrationRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "A domain registration record's internal identifier",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"name": {
				Description: "The name of the domain registration",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name_servers": {
				Description: "Custom name servers for the domain registration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMDNSDomainRegistrationRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	name := d.Get("name").(string)
	names, err := service.
		Filter(filter.Build(filter.Path("domainRegistrations.name").Eq(name))).
		Mask("id,name").
		GetDomainRegistrations()

	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving domain registration: %s", err)
	}

	if len(names) == 0 {
		return fmt.Errorf("[ERROR] No domain registration found with name [%s]", name)
	}

	log.Printf("names %v\n", names)
	dnsId := *names[0].Id
	log.Printf("Domain Registration Id %d\n", dnsId)

	// Get nameservers for domain
	nService := services.GetDnsDomainRegistrationService(sess)

	// retrieve remote object state
	dns_domain_nameservers, err := nService.Id(dnsId).
		Mask("nameservers.name").
		GetDomainNameservers()

	log.Printf("list %v\n", dns_domain_nameservers)

	ns := make([]string, len(dns_domain_nameservers[0].Nameservers))
	for i, elem := range dns_domain_nameservers[0].Nameservers {
		ns[i] = *elem.Name
	}

	log.Printf("names %v\n", ns)

	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving domain registration nameservers: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", dnsId))
	d.Set("name_servers", ns)
	return nil
}
