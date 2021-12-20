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

func dataSourceOvirtClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtClustersRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchema(),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"clusters": {
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
						"datacenter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"networks": ListOfIdName(),
					},
				},
			},
		},
	}
}

func dataSourceOvirtClustersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	clustersReq := conn.SystemService().ClustersService().
		List().Follow("networks")

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
			clustersReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			clustersReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			clustersReq.CaseSensitive(csBool)
		}
	}
	clustersResp, err := clustersReq.Send()
	if err != nil {
		return err
	}
	clusters, ok := clustersResp.Clusters()
	if !ok || len(clusters.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredClusters []*ovirtsdk4.Cluster
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, c := range clusters.Slice() {
			if r.MatchString(c.MustName()) {
				filteredClusters = append(filteredClusters, c)
			}
		}
	} else {
		filteredClusters = clusters.Slice()[:]
	}

	if len(filteredClusters) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return clustersDescriptionAttributes(d, filteredClusters, meta)
}

func clustersDescriptionAttributes(d *schema.ResourceData, clusters []*ovirtsdk4.Cluster, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range clusters {
		desc, ok := v.Description()
		if !ok {
			desc = ""
		}

		// local DCs doesn't return a data_center reference - use zero value instead.
		var dcId string
		dc, ok := v.DataCenter()
		if ok {
			dcId = dc.MustId()
		}
		mapping := map[string]interface{}{
			"id":            v.MustId(),
			"name":          v.MustName(),
			"datacenter_id": dcId,
			"description":   desc,
		}
		if slice, ok := v.Networks(); ok {
			networks := make([]map[string]interface{}, 0)
			for _, n := range slice.Slice() {
				networks = append(networks, map[string]interface{}{
					"id":   n.MustId(),
					"name": n.MustName(),
				})
			}
			mapping["networks"] = networks
		}
		s = append(s, mapping)
	}

	d.SetId(resource.UniqueId())
	return d.Set("clusters", s)
}
