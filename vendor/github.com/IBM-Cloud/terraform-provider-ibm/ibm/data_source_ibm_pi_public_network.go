// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	//"fmt"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func dataSourceIBMPIPublicNetwork() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPIPublicNetworksRead,
		Schema: map[string]*schema.Schema{

			helpers.PINetworkName: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Network Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
				Deprecated:   "This field is deprectaed.",
			},

			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes

			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
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
		},
	}
}

func dataSourceIBMPIPublicNetworksRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	networkC := instance.NewIBMPINetworkClient(sess, powerinstanceid)
	networkdata, err := networkC.GetPublic(powerinstanceid, getTimeOut)

	if err != nil {
		return err
	}
	if len(networkdata.Networks) < 1 {
		return fmt.Errorf("No Public Network Found in %s", powerinstanceid)
	}
	d.SetId(*networkdata.Networks[0].NetworkID)
	if networkdata.Networks[0].Type != nil {
		d.Set("type", networkdata.Networks[0].Type)
	}
	if networkdata.Networks[0].Name != nil {
		d.Set("name", networkdata.Networks[0].Name)
	}
	if networkdata.Networks[0].VlanID != nil {
		d.Set("vlan_id", networkdata.Networks[0].VlanID)
	}
	if networkdata.Networks[0].NetworkID != nil {
		d.Set("network_id", networkdata.Networks[0].NetworkID)
	}
	d.Set(helpers.PICloudInstanceId, powerinstanceid)

	return nil

}
