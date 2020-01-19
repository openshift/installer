// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
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

func dataSourceOvirtMacPools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtMacPoolsRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchemaWith("max"),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"mac_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow_duplicates": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceOvirtMacPoolsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	req := conn.SystemService().MacPoolsService().List()

	search, searchOK := d.GetOk("search")
	if searchOK {
		if v, ok := search.(map[string]interface{})["max"]; ok {
			maxInt, err := strconv.ParseInt(v.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			req.Max(maxInt)
		}
	}

	resp, err := req.Send()
	if err != nil {
		return err
	}
	pools, ok := resp.Pools()
	if !ok || len(pools.Slice()) == 0 {
		return fmt.Errorf("no mac pool exists")
	}

	var filteredMacPools []*ovirtsdk4.MacPool

	nameRegex, nameRegexOK := d.GetOk("name_regex")
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, c := range pools.Slice() {
			if r.MatchString(c.MustName()) {
				filteredMacPools = append(filteredMacPools, c)
			}
		}
	} else {
		filteredMacPools = pools.Slice()[:]
	}

	if len(filteredMacPools) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return macPoolsDescriptionAttributes(d, filteredMacPools, meta)
}

func macPoolsDescriptionAttributes(d *schema.ResourceData, macpools []*ovirtsdk4.MacPool, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range macpools {
		mapping := map[string]interface{}{
			"id":               v.MustId(),
			"name":             v.MustName(),
			"allow_duplicates": v.MustAllowDuplicates(),
			"ranges":           flattenMacPoolRanges(v.MustRanges()),
		}
		if desc, ok := v.Description(); ok {
			mapping["description"] = desc
		}
		s = append(s, mapping)
	}

	d.SetId(resource.UniqueId())
	return d.Set("mac_pools", s)
}
