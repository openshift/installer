// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
)

func dataSourceIBMNetworkVlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMNetworkVlanRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"number": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"router_hostname": {
				Type:     schema.TypeString,
				Optional: true,
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
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMNetworkVlanRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	name := d.Get("name").(string)
	number := d.Get("number").(int)
	routerHostname := d.Get("router_hostname").(string)
	var vlan *datatypes.Network_Vlan
	var err error

	if number != 0 && routerHostname != "" {
		// Got vlan number and router, get vlan, and compute name
		vlan, err = getVlan(number, routerHostname, name, meta)
		if err != nil {
			return err
		}

		d.SetId(fmt.Sprintf("%d", *vlan.Id))
		if vlan.Name != nil {
			d.Set("name", *vlan.Name)
		}
	} else if name != "" {
		// Got name, get vlan, and compute router hostname and vlan number
		filters := filter.New(filter.Path("networkVlans.name").Eq(name))
		if number != 0 {
			filters = append(filters, filter.Path("networkVlans.vlanNumber").Eq(number))
		}
		networkVlans, err := service.
			Mask("id,vlanNumber,name,primaryRouter[hostname],subnets[networkIdentifier,cidr,subnetType,id,gateway],virtualGuests[id,domain,hostname]").
			Filter(
				filter.Build(
					filters...,
				),
			).
			GetNetworkVlans()
		if err != nil {
			return fmt.Errorf("Error obtaining VLAN id: %s", err)
		} else if len(networkVlans) == 0 {
			return fmt.Errorf("No VLAN was found with the name '%s'", name)
		}

		vlan = &networkVlans[0]
		d.SetId(fmt.Sprintf("%d", *vlan.Id))
		d.Set("number", *vlan.VlanNumber)

		if vlan.PrimaryRouter != nil && vlan.PrimaryRouter.Hostname != nil {
			d.Set("router_hostname", *vlan.PrimaryRouter.Hostname)
		}
	} else {
		return errors.New("missing required properties. Need a VLAN name, or the VLAN's number and router hostname")
	}

	// Get subnets in cidr format for display
	sbns := make([]map[string]interface{}, len(vlan.Subnets))
	for i, elem := range vlan.Subnets {
		subnet := make(map[string]interface{})
		subnet["subnet"] = fmt.Sprintf("%s/%s", *elem.NetworkIdentifier, strconv.Itoa(*elem.Cidr))
		subnet["subnet_type"] = *elem.SubnetType
		subnet["subnet_size"] = 1 << (uint)(32-*elem.Cidr)
		subnet["cidr"] = *elem.Cidr
		subnet["id"] = *elem.Id
		if elem.Gateway != nil {
			subnet["gateway"] = *elem.Gateway
		}
		sbns[i] = subnet

	}
	d.Set("subnets", sbns)

	vgs := make([]map[string]interface{}, len(vlan.VirtualGuests))
	for i, vg := range vlan.VirtualGuests {
		v := make(map[string]interface{})
		v["id"] = *vg.Id
		v["domain"] = *vg.Domain
		v["hostname"] = *vg.Hostname
		vgs[i] = v
	}
	d.Set("virtual_guests", vgs)
	return nil
}

func getVlan(vlanNumber int, primaryRouterHostname string, name string, meta interface{}) (*datatypes.Network_Vlan, error) {
	service := services.GetAccountService(meta.(ClientSession).SoftLayerSession())

	filters := filter.New(filter.Path("networkVlans.primaryRouter.hostname").Eq(primaryRouterHostname),
		filter.Path("networkVlans.vlanNumber").Eq(vlanNumber))
	if name != "" {
		filters = append(filters, filter.Path("networkVlans.name").Eq(name))
	}
	networkVlans, err := service.
		Mask("id,name,subnets[networkIdentifier,cidr,subnetType,id,gateway],virtualGuests[id,domain,hostname]").
		Filter(
			filter.Build(
				filters...,
			),
		).
		GetNetworkVlans()

	if err != nil {
		return &datatypes.Network_Vlan{}, fmt.Errorf("Error looking up Vlan: %s", err)
	}

	if len(networkVlans) < 1 {
		return &datatypes.Network_Vlan{}, fmt.Errorf(
			"Unable to locate a vlan matching the provided router hostname and vlan number: %s/%d",
			primaryRouterHostname,
			vlanNumber)
	}

	return &networkVlans[0], nil
}
