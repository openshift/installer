// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"net"
	"strconv"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceIBMPIInstanceIP() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceIBMPIInstancesIPRead,
		Schema: map[string]*schema.Schema{

			helpers.PIInstanceName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Server Name to be used for pvminstances",
				ValidateFunc: validation.NoZeroValues,
			},

			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			helpers.PINetworkName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipoctet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"macaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMPIInstancesIPRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		return err
	}

	checkValidSubnet(d, meta)

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	powerinstancesubnet := d.Get(helpers.PINetworkName).(string)
	powerC := instance.NewIBMPIInstanceClient(sess, powerinstanceid)
	powervmdata, err := powerC.Get(d.Get(helpers.PIInstanceName).(string), powerinstanceid, getTimeOut)

	if err != nil {
		return err
	}

	for i, _ := range powervmdata.Addresses {
		if powervmdata.Addresses[i].NetworkName == powerinstancesubnet {
			log.Printf("Printing the ip %s", powervmdata.Addresses[i].IP)
			d.Set("ip", powervmdata.Addresses[i].IP)
			d.Set("network_id", powervmdata.Addresses[i].NetworkID)
			d.Set("macaddress", powervmdata.Addresses[i].MacAddress)
			d.Set("external_ip", powervmdata.Addresses[i].ExternalIP)
			d.Set("type", powervmdata.Addresses[i].Type)

			IPObject := net.ParseIP(powervmdata.Addresses[i].IP).To4()

			d.Set("ipoctet", strconv.Itoa(int(IPObject[3])))

		}

	}

	return nil

}

func checkValidSubnet(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		return err
	}

	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	powerinstancesubnet := d.Get(helpers.PINetworkName).(string)

	networkC := instance.NewIBMPINetworkClient(sess, powerinstanceid)
	networkdata, err := networkC.Get(powerinstancesubnet, powerinstanceid, getTimeOut)

	if err != nil {
		return err
	}

	d.SetId(*networkdata.NetworkID)

	return nil
}
