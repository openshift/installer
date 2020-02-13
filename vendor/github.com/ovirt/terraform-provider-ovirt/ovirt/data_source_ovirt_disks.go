// Copyright (C) 2017 Battelle Memorial Institute
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

func dataSourceOvirtDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtDisksRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			// Computed
			"disks": {
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
						"alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"shareable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sparse": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvirtDisksRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	disksReq := conn.SystemService().DisksService().List()

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
			disksReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			disksReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			disksReq.CaseSensitive(csBool)
		}
	}

	disksResp, err := disksReq.Send()
	if err != nil {
		return err
	}
	disks, ok := disksResp.Disks()
	if !ok || len(disks.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredDisks []*ovirtsdk4.Disk
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, disk := range disks.Slice() {
			if r.MatchString(disk.MustName()) {
				filteredDisks = append(filteredDisks, disk)
			}
		}
	} else {
		filteredDisks = disks.Slice()[:]
	}

	if len(filteredDisks) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return disksDescriptionAttributes(d, filteredDisks, meta)

}

func disksDescriptionAttributes(d *schema.ResourceData, disks []*ovirtsdk4.Disk, meta interface{}) error {
	var s []map[string]interface{}

	for _, v := range disks {
		mapping := map[string]interface{}{
			"id":     v.MustId(),
			"name":   v.MustName(),
			"format": v.MustFormat(),
			"size":   v.MustProvisionedSize(),
		}
		if sds, ok := v.StorageDomains(); ok {
			if len(sds.Slice()) > 0 {
				mapping["storage_domain_id"] = sds.Slice()[0].MustId()
			}
		}
		if quota, ok := v.Quota(); ok {
			mapping["quota_id"] = quota.MustId()
		}
		if alias, ok := v.Alias(); ok {
			mapping["alias"] = alias
		}
		if shareable, ok := v.Shareable(); ok {
			mapping["shareable"] = shareable
		}
		if sparse, ok := v.Sparse(); ok {
			mapping["sparse"] = sparse
		}

		s = append(s, mapping)
	}
	d.SetId(resource.UniqueId())
	return d.Set("disks", s)
}
