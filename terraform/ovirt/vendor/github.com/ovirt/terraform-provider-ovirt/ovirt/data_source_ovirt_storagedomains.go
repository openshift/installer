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

func dataSourceOvirtStorageDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtStorageDomainsRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"storagedomains": {
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
						"datacenter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceOvirtStorageDomainsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	sdsReq := conn.SystemService().StorageDomainsService().List()

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
			sdsReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			sdsReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			sdsReq.CaseSensitive(csBool)
		}
	}
	sdsResp, err := sdsReq.Send()
	if err != nil {
		return err
	}
	sds, ok := sdsResp.StorageDomains()
	if !ok || len(sds.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredSds []*ovirtsdk4.StorageDomain
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, sd := range sds.Slice() {
			if r.MatchString(sd.MustName()) {
				filteredSds = append(filteredSds, sd)
			}
		}
	} else {
		filteredSds = sds.Slice()[:]
	}

	if len(filteredSds) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return storageDomainsDecriptionAttributes(d, filteredSds, meta)
}

func storageDomainsDecriptionAttributes(d *schema.ResourceData, storagedomains []*ovirtsdk4.StorageDomain, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range storagedomains {
		// description is not mandatory and if using MustDescription will fail with nil value
		mapping := map[string]interface{}{
			"id":   v.MustId(),
			"name": v.MustName(),
			"type": v.MustType(),
		}
		if description, ok := v.Description(); ok {
			mapping["description"] = description
		}
		if dcs, ok := v.DataCenters(); ok {
			mapping["datacenter_id"] = dcs.Slice()[0].MustId()
		}

		if externalStatus, ok := v.ExternalStatus(); ok {
			mapping["external_status"] = externalStatus
		}
		if status, ok := v.Status(); ok {
			mapping["status"] = status
		}
		s = append(s, mapping)
	}
	d.SetId(resource.UniqueId())
	return d.Set("storagedomains", s)
}
