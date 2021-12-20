// Copyright (C) 2018 Boris Manojlovic
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func networkSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"datacenter_id": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"vlan_id": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"mtu": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func dataSourceOvirtNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtNetworksRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"networks": networkSchema(),
		},
	}
}

func dataSourceOvirtNetworksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	networksReq := conn.SystemService().NetworksService().List()

	search, searchOK := d.GetOk("search")
	nameRegex, nameRegexOK := d.GetOk("name_regex")

	if searchOK {
		searchMap := search.(map[string]interface{})
		searchCriteria, searchCriteriaOK := searchMap["criteria"]
		searchMax, searchMaxOK := searchMap["max"]
		searchCaseSensitive, searchCaseSensitiveOK := searchMap["case_sensitive"]
		if !searchCriteriaOK && !searchMaxOK && !searchCaseSensitiveOK {
			return fmt.Errorf("One of criteria or max or case_sensitive in search must be assigned")
		}

		if searchCriteriaOK {
			networksReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			networksReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			networksReq.CaseSensitive(csBool)
		}
	}

	networksResp, err := networksReq.Send()
	if err != nil {
		return err
	}
	networks, ok := networksResp.Networks()
	if !ok && len(networks.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredNetworks []*ovirtsdk4.Network
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, network := range networks.Slice() {
			if r.MatchString(network.MustName()) {
				filteredNetworks = append(filteredNetworks, network)
			}
		}
	} else {
		filteredNetworks = networks.Slice()[:]
	}

	if len(filteredNetworks) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return networksDecriptionAttributes(d, filteredNetworks, meta)
}

func networksDecriptionAttributes(d *schema.ResourceData, network []*ovirtsdk4.Network, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range network {
		// description is not mandatory and if using MustDescription will fail with nil value
		desc, ok := v.Description()
		if !ok {
			desc = ""
		}
		mapping := map[string]interface{}{
			"id":            v.MustId(),
			"name":          v.MustName(),
			"datacenter_id": v.MustDataCenter().MustId(),
			"description":   desc,
			"mtu":           v.MustMtu(),
		}
		if vlan, ok := v.Vlan(); ok {
			mapping["vlan_id"] = vlan.MustId()
		}
		s = append(s, mapping)
	}
	d.SetId(resource.UniqueId())
	return d.Set("networks", s)
}
