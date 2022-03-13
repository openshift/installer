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

func dataSourceOvirtDataCenters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtDataCentersRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			// Computed
			"datacenters": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"quota_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvirtDataCentersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	dcsReq := conn.SystemService().DataCentersService().List()

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
			dcsReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			dcsReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			dcsReq.CaseSensitive(csBool)
		}
	}
	dcsResp, err := dcsReq.Send()
	if err != nil {
		return err
	}
	dcs, ok := dcsResp.DataCenters()
	if !ok || len(dcs.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredDcs []*ovirtsdk4.DataCenter
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, dc := range dcs.Slice() {
			if r.MatchString(dc.MustName()) {
				filteredDcs = append(filteredDcs, dc)
			}
		}
	} else {
		filteredDcs = dcs.Slice()[:]
	}

	if len(filteredDcs) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return dataCentersDescriptionAttributes(d, filteredDcs, meta)
}

func dataCentersDescriptionAttributes(d *schema.ResourceData, dcs []*ovirtsdk4.DataCenter, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range dcs {
		mapping := map[string]interface{}{
			"id":         v.MustId(),
			"name":       v.MustName(),
			"status":     v.MustStatus(),
			"local":      v.MustLocal(),
			"quota_mode": v.MustQuotaMode(),
		}
		s = append(s, mapping)
	}
	d.SetId(resource.UniqueId())
	return d.Set("datacenters", s)
}
