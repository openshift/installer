// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIInstanceIP() *schema.Resource {

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
	sess, err := meta.(conns.ClientSession).IBMPISession()
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

	for _, network := range powervmdata.Networks {
		if network.NetworkName == networkName {
			log.Printf("Printing the ip %s", network.IPAddress)
			d.SetId(network.NetworkID)
			d.Set("ip", network.IPAddress)
			d.Set("network_id", network.NetworkID)
			d.Set("macaddress", network.MacAddress)
			d.Set("external_ip", network.ExternalIP)
			d.Set("type", network.Type)

			IPObject := net.ParseIP(network.IPAddress).To4()
			if len(IPObject) > 0 {
				d.Set("ipoctet", strconv.Itoa(int(IPObject[3])))
			}

			return nil
		}
	}

	return diag.Errorf("failed to find instance ip that belongs to the given network")
}
