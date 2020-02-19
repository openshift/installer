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

func dataSourceOvirtUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtUsersRead,
		Schema: map[string]*schema.Schema{
			"search": dataSourceSearchSchemaWith(
				"max", "criteria", "case_sensitive"),
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},

			"users": {
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
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"authz_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvirtUsersRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	usersReq := conn.SystemService().UsersService().List()

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
			usersReq.Search(searchCriteria.(string))
		}
		if searchMaxOK {
			maxInt, err := strconv.ParseInt(searchMax.(string), 10, 64)
			if err != nil || maxInt < 1 {
				return fmt.Errorf("search.max must be a positive int")
			}
			usersReq.Max(maxInt)
		}
		if searchCaseSensitiveOK {
			csBool, err := strconv.ParseBool(searchCaseSensitive.(string))
			if err != nil {
				return fmt.Errorf("search.case_sensitive must be true or false")
			}
			usersReq.CaseSensitive(csBool)
		}
	}
	usersResp, err := usersReq.Send()
	if err != nil {
		return err
	}
	users, ok := usersResp.Users()
	if !ok || len(users.Slice()) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	var filteredUsers []*ovirtsdk4.User
	if nameRegexOK {
		r := regexp.MustCompile(nameRegex.(string))
		for _, c := range users.Slice() {
			if r.MatchString(c.MustName()) {
				filteredUsers = append(filteredUsers, c)
			}
		}
	} else {
		filteredUsers = users.Slice()[:]
	}

	if len(filteredUsers) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return usersDescriptionAttributes(d, filteredUsers, meta)
}

func usersDescriptionAttributes(d *schema.ResourceData, users []*ovirtsdk4.User, meta interface{}) error {
	var s []map[string]interface{}
	for _, v := range users {
		mapping := map[string]interface{}{
			"id":         v.MustId(),
			"name":       v.MustName(),
			"namespace":  v.MustNamespace(),
			"principal":  v.MustPrincipal(),
			"user_name":  v.MustUserName(),
			"authz_name": v.MustDomain().MustName(),
		}
		s = append(s, mapping)
	}

	d.SetId(resource.UniqueId())
	return d.Set("users", s)
}
