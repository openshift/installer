// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
)

func dataSourceIBMSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSecurityGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the security group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the security group",
			},
			"most_recent": &schema.Schema{
				Description: "If true and multiple entries are found, the most recently created group is used. " +
					"If false, an error is returned",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceIBMSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	name := d.Get("name").(string)
	mostRecent := d.Get("most_recent").(bool)

	filters := filter.New(filter.Path("securityGroups.name").Eq(name))
	if v, ok := d.GetOk("description"); ok {
		filters = append(filters, filter.Path("securityGroups.description").Eq(v.(string)))
	}

	groups, err := services.GetAccountService(sess).
		Filter(filter.Build(
			filters...,
		)).
		GetSecurityGroups()

	if err != nil {
		return fmt.Errorf("Error retrieving Security group: %s", err)
	}
	if len(groups) == 0 {
		return fmt.Errorf("No security group found with name [%s]", name)
	}

	var sg datatypes.Network_SecurityGroup
	if len(groups) > 1 {
		if mostRecent {
			sg = mostRecentSecurityGroup(groups)
		} else {
			return fmt.Errorf(
				"More than one security group found with name matching [%s]. "+
					"Either set 'most_recent' to true in your "+
					"configuration to force the most security group "+
					"to be used, or ensure that the name and/or description is unique", name)
		}
	} else {
		sg = groups[0]
	}
	d.SetId(fmt.Sprintf("%d", *sg.Id))
	d.Set("description", sg.Description)
	return nil
}

type securityGroups []datatypes.Network_SecurityGroup

func (sgs securityGroups) Len() int { return len(sgs) }

func (sgs securityGroups) Swap(i, j int) { sgs[i], sgs[j] = sgs[j], sgs[i] }

func (sgs securityGroups) Less(i, j int) bool {
	return sgs[i].CreateDate.Before(sgs[j].CreateDate.Time)
}

func mostRecentSecurityGroup(sgs securityGroups) datatypes.Network_SecurityGroup {
	sortedKeys := sgs
	sort.Sort(sortedKeys)
	return sortedKeys[len(sortedKeys)-1]
}
