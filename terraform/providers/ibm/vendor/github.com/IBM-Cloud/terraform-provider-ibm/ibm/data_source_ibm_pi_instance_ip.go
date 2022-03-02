// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceIBMPIInstanceIP() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstancesIPRead,
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

			// Computed attributes
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

func dataSourceIBMPIInstancesIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	networkName := d.Get(helpers.PINetworkName).(string)
	powerC := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

	powervmdata, err := powerC.Get(d.Get(helpers.PIInstanceName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	for _, address := range powervmdata.Addresses {
		if address.NetworkName == networkName {
			log.Printf("Printing the ip %s", address.IP)
			d.SetId(address.NetworkID)
			d.Set("ip", address.IP)
			d.Set("network_id", address.NetworkID)
			d.Set("macaddress", address.MacAddress)
			d.Set("external_ip", address.ExternalIP)
			d.Set("type", address.Type)

			IPObject := net.ParseIP(address.IP).To4()

			d.Set("ipoctet", strconv.Itoa(int(IPObject[3])))

			return nil
		}
	}

	return diag.Errorf("failed to find instance ip that belongs to the given network")
}
