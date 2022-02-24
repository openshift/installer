// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"

	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"

	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	PIDhcpID = "pi_dhcp_id"
)

func dataSourceIBMPIDhcp() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPIDhcpRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			PIDhcpID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the DHCP Server",
			},
			// Computed Attributes
			PIDhcpStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the DHCP Server",
			},
			PIDhcpNetwork: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DHCP Server private network",
			},
			PIDhcpLeases: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of DHCP Server PVM Instance leases",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PIDhcpInstanceIp: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP of the PVM Instance",
						},
						PIDhcpInstanceMac: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC Address of the PVM Instance",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPIDhcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	dhcpID := d.Get(PIDhcpID).(string)

	client := st.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)
	dhcpServer, err := client.Get(dhcpID)
	if err != nil {
		log.Printf("[DEBUG] get DHCP failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(*dhcpServer.ID)
	d.Set(PIDhcpStatus, *dhcpServer.Status)
	dhcpNetwork := dhcpServer.Network
	if dhcpNetwork != nil {
		d.Set(PIDhcpNetwork, *dhcpNetwork.ID)
	}
	dhcpLeases := dhcpServer.Leases
	if dhcpLeases != nil {
		leaseList := make([]map[string]string, len(dhcpLeases))
		for i, lease := range dhcpLeases {
			leaseList[i] = map[string]string{
				PIDhcpInstanceIp:  *lease.InstanceIP,
				PIDhcpInstanceMac: *lease.InstanceMacAddress,
			}
		}
		d.Set(PIDhcpLeases, leaseList)
	}

	return nil
}
