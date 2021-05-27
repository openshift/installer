// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	//"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func dataSourceIBMPINetwork() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPINetworksRead,
		Schema: map[string]*schema.Schema{

			helpers.PINetworkName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Network Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
			},

			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vlan_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_ip_count": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"used_ip_count": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"used_ip_percent": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMPINetworksRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	networkC := instance.NewIBMPINetworkClient(sess, powerinstanceid)
	networkdata, err := networkC.Get(d.Get(helpers.PINetworkName).(string), powerinstanceid, getTimeOut)

	if err != nil {
		return err
	}

	d.SetId(*networkdata.NetworkID)
	if networkdata.Cidr != nil {
		d.Set("cidr", networkdata.Cidr)
	}
	if networkdata.Type != nil {
		d.Set("type", networkdata.Type)
	}
	if &networkdata.Gateway != nil {
		d.Set("gateway", networkdata.Gateway)
	}
	if networkdata.VlanID != nil {
		d.Set("vlan_id", networkdata.VlanID)
	}
	if networkdata.IPAddressMetrics.Available != nil {
		d.Set("available_ip_count", networkdata.IPAddressMetrics.Available)
	}
	if networkdata.IPAddressMetrics.Used != nil {
		d.Set("used_ip_count", networkdata.IPAddressMetrics.Used)
	}
	if networkdata.IPAddressMetrics.Utilization != nil {
		d.Set("used_ip_percent", networkdata.IPAddressMetrics.Utilization)
	}
	if networkdata.Name != nil {
		d.Set("name", networkdata.Name)
	}

	return nil

}
