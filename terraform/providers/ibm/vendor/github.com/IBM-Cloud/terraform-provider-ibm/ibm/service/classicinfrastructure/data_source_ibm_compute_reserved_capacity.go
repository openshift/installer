// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package classicinfrastructure

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func DataSourceIBMComputeReservedCapacity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMComputeReservedCapacityRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of reserved instance",
			},

			"most_recent": {
				Description: "If true and multiple entries are found, the most recently created reserved capacity is used. " +
					"If false, an error is returned",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dataceneter name",
			},

			"pod": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pod name",
			},

			"instances": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "no of the instances",
			},

			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "flavor of the reserved capacity",
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

func dataSourceIBMComputeReservedCapacityRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess := meta.(conns.ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	name := d.Get("name").(string)
	mostRecent := d.Get("most_recent").(bool)

	grps, err := service.
		Filter(filter.Build(filter.Path("reservedCapacityGroups.name").Eq(name))).
		Mask("id,name,instancesCount,createDate,backendRouter[hostname,datacenter[name]],occupiedInstances[guest[id,domain,hostname]],instances[billingItem[item[keyName]]]").GetReservedCapacityGroups()

	if err != nil {
		return diag.FromErr(fmt.Errorf("[Error] retrieving placement group: %s", err))
	}

	if len(grps) == 0 {
		return diag.FromErr(fmt.Errorf("[Error] No reserved capacity found with name [%s]", name))
	}

	var grp datatypes.Virtual_ReservedCapacityGroup

	if len(grps) > 1 {
		if mostRecent {
			grp = mostRecentReservedCapacity(grps)
		} else {
			return diag.FromErr(fmt.Errorf(
				"[Error] More than one reserved capacity found with name "+
					"matching [%s]. Set 'most_recent' to true in your configuration to force the most recent reserved capacity "+
					"to be used", name))
		}
	} else {
		grp = grps[0]
	}

	d.SetId(fmt.Sprintf("%d", *grp.Id))
	d.Set("name", grp.Name)
	d.Set("datacenter", grp.BackendRouter.Datacenter.Name)
	pod := strings.SplitAfter(*grp.BackendRouter.Hostname, ".")[0]
	r, _ := regexp.Compile("[0-9]{2}")
	pod = "pod" + r.FindString(pod)
	d.Set("pod", pod)
	d.Set("instances", grp.InstancesCount)
	keyName, ok := sl.GrabOk(grp, "Instances.0.BillingItem.Item.KeyName")
	if ok {
		d.Set("flavor", keyName)
	}

	vgs := make([]map[string]interface{}, len(grp.OccupiedInstances))
	for i, vg := range grp.OccupiedInstances {
		if vg.Guest != nil {
			v := make(map[string]interface{})
			v["id"] = *vg.Guest.Id
			v["domain"] = *vg.Guest.Domain
			v["hostname"] = *vg.Guest.Hostname
			vgs[i] = v
		}

	}
	d.Set("virtual_guests", vgs)

	return nil
}

type reservedCapacityGroups []datatypes.Virtual_ReservedCapacityGroup

func (k reservedCapacityGroups) Len() int { return len(k) }

func (k reservedCapacityGroups) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

func (k reservedCapacityGroups) Less(i, j int) bool {
	return k[i].CreateDate.Before(k[j].CreateDate.Time)
}

func mostRecentReservedCapacity(keys reservedCapacityGroups) datatypes.Virtual_ReservedCapacityGroup {
	sortedKeys := keys
	sort.Sort(sortedKeys)
	return sortedKeys[len(sortedKeys)-1]
}
