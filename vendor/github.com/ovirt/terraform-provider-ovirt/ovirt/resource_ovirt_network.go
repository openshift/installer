// Copyright (C) 2018 Joey Ma <majunjiev@gmail.com>
// All rights reserved.
//
// This software may be modified and distributed under the terms
// of the BSD-2 license.  See the LICENSE file for details.

package ovirt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

func resourceOvirtNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvirtNetworkCreate,
		Read:   resourceOvirtNetworkRead,
		Update: resourceOvirtNetworkUpdate,
		Delete: resourceOvirtNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"datacenter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"mtu": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceOvirtNetworkCreate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*ovirtsdk4.Connection)
	name := d.Get("name").(string)
	datacenterID := d.Get("datacenter_id").(string)

	datacenter, err := ovirtsdk4.NewDataCenterBuilder().Id(datacenterID).Build()
	if err != nil {
		return err
	}
	//Name and datacenter_id are required when create a datacenter
	networkbuilder := ovirtsdk4.NewNetworkBuilder().Name(name).
		DataCenter(datacenter)

	if description, ok := d.GetOkExists("description"); ok {
		networkbuilder = networkbuilder.Description(description.(string))
	}

	if vlanID, ok := d.GetOkExists("vlan_id"); ok {
		vlanIDInt := vlanID.(int)
		vlan, err := ovirtsdk4.NewVlanBuilder().Id(int64(vlanIDInt)).Build()
		if err != nil {
			return err
		}
		networkbuilder = networkbuilder.Vlan(vlan)
	}

	if mtu, ok := d.GetOkExists("mtu"); ok {
		mtuInt := mtu.(int)
		networkbuilder = networkbuilder.Mtu(int64(mtuInt))
	}

	network, err := networkbuilder.Build()
	if err != nil {
		return err
	}

	addResp, err := conn.SystemService().NetworksService().Add().Network(network).Send()
	if err != nil {
		return err
	}

	d.SetId(addResp.MustNetwork().MustId())
	return resourceOvirtNetworkRead(d, meta)

}

func resourceOvirtNetworkUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*ovirtsdk4.Connection)
	networkService := conn.SystemService().NetworksService().NetworkService(d.Id())
	networkBuilder := ovirtsdk4.NewNetworkBuilder()

	if name, ok := d.GetOkExists("name"); ok {
		if d.HasChange("name") {
			networkBuilder.Name(name.(string))
		}
	} else {
		return fmt.Errorf("Network's name does not exist")
	}

	if datacenterID, ok := d.GetOkExists("datacenter_id"); ok {
		if d.HasChange("datacenter_id") {
			datacenter, err := ovirtsdk4.NewDataCenterBuilder().Id(datacenterID.(string)).Build()
			if err != nil {
				return err
			}
			networkBuilder.DataCenter(datacenter)
		}
	} else {
		return fmt.Errorf("Network's 'datacenter_id' does not exist")
	}

	if description, ok := d.GetOkExists("description"); ok && d.HasChange("description") {
		networkBuilder.Description(description.(string))
	}

	if vlanID, ok := d.GetOkExists("vlan_id"); ok && d.HasChange("vlan_id") {
		vlanIDInt := vlanID.(int)
		vlan, err := ovirtsdk4.NewVlanBuilder().Id(int64(vlanIDInt)).Build()
		if err != nil {
			return err
		}
		networkBuilder.Vlan(vlan)
	}

	if mtu, ok := d.GetOkExists("mtu"); ok && d.HasChange("mtu") {
		mtuInt := mtu.(int)
		networkBuilder.Mtu(int64(mtuInt))
	}

	network, err := networkBuilder.Build()
	if err != nil {
		return err
	}

	_, err = networkService.Update().Network(network).Send()

	if err != nil {
		return err
	}

	return resourceOvirtNetworkRead(d, meta)
}

func resourceOvirtNetworkRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*ovirtsdk4.Connection)
	getNetworkResp, err := conn.SystemService().NetworksService().
		NetworkService(d.Id()).Get().Send()
	if err != nil {
		return err
	}

	network, ok := getNetworkResp.Network()
	if !ok {
		d.SetId("")
		return nil
	}

	d.Set("name", network.MustName())

	if datacenter, ok := network.DataCenter(); ok {
		if datacenterID, ok := datacenter.Id(); ok {
			d.Set("datacenter_id", datacenterID)
		} else {
			return fmt.Errorf("Network's datacenter_id does not exist")
		}
	}

	if vlan, ok := network.Vlan(); ok {
		d.Set("vlan_id", vlan.MustId())
	}

	if description, ok := network.Description(); ok {
		d.Set("description", description)
	}

	mtu, ok := network.Mtu()
	if ok {
		d.Set("mtu", mtu)
	}

	return nil
}

func resourceOvirtNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*ovirtsdk4.Connection)

	_, err := conn.SystemService().NetworksService().
		NetworkService(d.Id()).Remove().Send()
	if err != nil {
		return err
	}
	return nil
}

func resourceOvirtNetworkImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*ovirtsdk4.Connection)

	resp, err := conn.SystemService().NetworksService().NetworkService(d.Id()).Get().Send()
	if err != nil {
		return nil, err
	}
	network, ok := resp.Network()
	if !ok {
		d.SetId("")
		return nil, nil
	}
	d.Set("name", network.MustName())
	d.Set("vlan_id", network.MustVlan())
	d.Set("mtu", network.MustMtu())
	if datacenter, ok := network.DataCenter(); ok {
		if datacenterID, ok := datacenter.Id(); ok {
			d.Set("datacenter_id", datacenterID)
		} else {
			return nil, fmt.Errorf("Network's datacenter_id does not exist")
		}
	}
	if description, ok := network.Description(); ok {
		d.Set("description", description)
	}
	return []*schema.ResourceData{d}, nil
}
