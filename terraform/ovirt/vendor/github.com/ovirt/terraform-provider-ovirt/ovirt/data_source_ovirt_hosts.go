// Copyright (C) 2019 Joey Ma <majunjiev@gmail.com>
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

func dataSourceOvirtHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtHostsRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"hosts": {
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
							Computed: true},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvirtHostsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	hostsReq := conn.SystemService().HostsService().List()

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
			hostsReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			hostsReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			hostsReq.CaseSensitive(csBool)
		}
	}
	hostsResp, err := hostsReq.Send()
	if err != nil {
		return err
	}
	hosts, ok := hostsResp.Hosts()
	if !ok || len(hosts.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredHosts []*ovirtsdk4.Host
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, h := range hosts.Slice() {
			if r.MatchString(h.MustName()) {
				filteredHosts = append(filteredHosts, h)
			}
		}
	} else {
		filteredHosts = hosts.Slice()[:]
	}

	if len(filteredHosts) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return hostsDescriptionAttributes(d, filteredHosts, meta)
}

func hostsDescriptionAttributes(d *schema.ResourceData, hosts []*ovirtsdk4.Host, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range hosts {
		desc, ok := v.Description()
		if !ok {
			desc = ""
		}
		mapping := map[string]interface{}{
			"id":          v.MustId(),
			"name":        v.MustName(),
			"cluster_id":  v.MustCluster().MustId(),
			"address":     v.MustAddress(),
			"description": desc,
		}
		s = append(s, mapping)
	}

	d.SetId(resource.UniqueId())
	return d.Set("hosts", s)
}
