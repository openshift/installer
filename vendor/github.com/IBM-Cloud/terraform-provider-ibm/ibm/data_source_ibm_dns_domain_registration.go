// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"log"
)

func dataSourceIBMDNSDomainRegistration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDNSDomainRegistrationRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Description: "A domain registration record's internal identifier",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"name": &schema.Schema{
				Description: "The name of the domain registration",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name_servers": &schema.Schema{
				Description: "Custom name servers for the domain registration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMDNSDomainRegistrationRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	name := d.Get("name").(string)
	names, err := service.
		Filter(filter.Build(filter.Path("domainRegistrations.name").Eq(name))).
		Mask("id,name").
		GetDomainRegistrations()

	if err != nil {
		return fmt.Errorf("Error retrieving domain registration: %s", err)
	}

	if len(names) == 0 {
		return fmt.Errorf("No domain registration found with name [%s]", name)
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
		return fmt.Errorf("Error retrieving domain registration nameservers: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", dnsId))
	d.Set("name_servers", ns)
	return nil
}
