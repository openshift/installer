// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
)

func dataSourceIBMComputePlacementGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMComputePlacementGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"datacenter": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"pod": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"rule": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_guests": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMComputePlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	name := d.Get("name").(string)

	groups, err := service.
		Filter(filter.Build(filter.Path("placementGroup.name").Eq(name))).
		Mask("id,name,rule[name],guests[id,domain,hostname],backendRouter[hostname,datacenter[name]]").GetPlacementGroups()

	if err != nil {
		return fmt.Errorf("Error retrieving placement group: %s", err)
	}

	grps := []datatypes.Virtual_PlacementGroup{}
	for _, g := range groups {
		if name == *g.Name {
			grps = append(grps, g)

		}
	}

	if len(grps) == 0 {
		return fmt.Errorf("No placement group found with name [%s]", name)
	}

	var grp datatypes.Virtual_PlacementGroup

	grp = grps[0]

	d.SetId(fmt.Sprintf("%d", *grp.Id))
	d.Set("name", grp.Name)
	d.Set("datacenter", grp.BackendRouter.Datacenter.Name)
	pod := strings.SplitAfter(*grp.BackendRouter.Hostname, ".")[0]
	r, _ := regexp.Compile("[0-9]{2}")
	pod = "pod" + r.FindString(pod)
	d.Set("pod", pod)
	d.Set("rule", grp.Rule.Name)

	vgs := make([]map[string]interface{}, len(grp.Guests))
	for i, vg := range grp.Guests {
		v := make(map[string]interface{})
		v["id"] = *vg.Id
		v["domain"] = *vg.Domain
		v["hostname"] = *vg.Hostname
		vgs[i] = v
	}
	d.Set("virtual_guests", vgs)

	return nil
}
