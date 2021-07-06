// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisIPv4CIDRs = "ipv4_cidrs"
	cisIPv6CIDRs = "ipv6_cidrs"
)

func dataSourceIBMCISIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISIPRead,

		Schema: map[string]*schema.Schema{
			cisIPv4CIDRs: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			cisIPv6CIDRs: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMCISIPRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisIPClientSession()
	if err != nil {
		return err
	}
	opt := cisClient.NewListIpsOptions()
	result, response, err := cisClient.ListIps(opt)
	if err != nil {
		log.Printf("Failed to list IP addresses: %v", response)
		return err
	}

	d.Set(cisIPv4CIDRs, flattenStringList(result.Result.Ipv4Cidrs))
	d.Set(cisIPv6CIDRs, flattenStringList(result.Result.Ipv4Cidrs))
	d.SetId(dataSourceIBMCISIPID(d))
	return nil
}

func dataSourceIBMCISIPID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
