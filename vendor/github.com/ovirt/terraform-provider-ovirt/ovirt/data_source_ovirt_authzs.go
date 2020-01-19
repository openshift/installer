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

func dataSourceOvirtAuthzs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtAuthzsRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchemaWith("max"),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"authzs": {
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
					},
				},
			},
		},
	}
}

func dataSourceOvirtAuthzsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	req := conn.SystemService().DomainsService().List()

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
	authzs, ok := resp.Domains()
	if !ok || len(authzs.Slice()) == 0 {
		return fmt.Errorf("no authzs exists")
	}

	var filteredAuthzs []*ovirtsdk4.Domain

	nameRegex, nameRegexOK := d.GetOk("name_regex")
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, c := range authzs.Slice() {
			if r.MatchString(c.MustName()) {
				filteredAuthzs = append(filteredAuthzs, c)
			}
		}
	} else {
		filteredAuthzs = authzs.Slice()[:]
	}

	if len(filteredAuthzs) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return authzsDescriptionAttributes(d, filteredAuthzs, meta)
}

func authzsDescriptionAttributes(d *schema.ResourceData, authzs []*ovirtsdk4.Domain, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range authzs {
		mapping := map[string]interface{}{
			"id":   v.MustId(),
			"name": v.MustName(),
		}
		s = append(s, mapping)
	}

	d.SetId(resource.UniqueId())
	return d.Set("authzs", s)
}
