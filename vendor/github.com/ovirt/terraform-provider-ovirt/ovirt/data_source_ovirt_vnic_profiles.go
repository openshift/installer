// Copyright (C) 2017 Battelle Memorial Institute
// Copyright (C) 2018 Boris Manojlovic
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func dataSourceOvirtVNicProfiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOvirtVNicProfilesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"vnic_profiles": {
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
						"network_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"migratable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"port_mirroring": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOvirtVNicProfilesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)
	profilesService := conn.SystemService().VnicProfilesService()
	pfsResp, _ := profilesService.List().Send()
	pfSlice, _ := pfsResp.Profiles()
	var s []map[string]interface{}
	rexName := regexp.MustCompile(d.Get("name_regex").(string))
	networkID := d.Get("network_id")

	for _, pf := range pfSlice.Slice() {
		if rexName.FindString(pf.MustName()) != "" &&
			networkID == pf.MustNetwork().MustId() {
			mapping := map[string]interface{}{
				"id":             pf.MustId(),
				"name":           pf.MustName(),
				"network_id":     pf.MustNetwork().MustId(),
				"port_mirroring": pf.MustPortMirroring(),
			}
			if migratable, ok := pf.Migratable(); ok {
				mapping["migratable"] = migratable
			}
			s = append(s, mapping)
		}
	}
	if len(s) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}
	d.SetId(resource.UniqueId())
	err := d.Set("vnic_profiles", s)
	return err
}
