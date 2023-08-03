// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"

	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIDhcp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDhcpRead,
		Schema: map[string]*schema.Schema{

			// Required Arguments
			Arg_CloudInstanceID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_DhcpID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the DHCP Server",
			},

			// Attributes
			Attr_DhcpID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DHCP Server",
			},
			Attr_DhcpLeases: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of DHCP Server PVM Instance leases",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_DhcpLeaseInstanceIP: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP of the PVM Instance",
						},
						Attr_DhcpLeaseInstanceMac: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC Address of the PVM Instance",
						},
					},
				},
			},
			Attr_DhcpNetworkDeprecated: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DHCP Server private network (deprecated - replaced by network_id)",
			},
			Attr_DhcpNetworkID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DHCP Server private network",
			},
			Attr_DhcpNetworkName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the DHCP Server private network",
			},
			Attr_DhcpStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the DHCP Server",
			},
		},
	}
}

func dataSourceIBMPIDhcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	dhcpID := d.Get(Arg_DhcpID).(string)

	// client
	client := st.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)

	// get dhcp
	dhcpServer, err := client.Get(dhcpID)
	if err != nil {
		log.Printf("[DEBUG] get DHCP failed %v", err)
		return diag.FromErr(err)
	}

	// set attributes
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *dhcpServer.ID))
	d.Set(Attr_DhcpID, *dhcpServer.ID)
	d.Set(Attr_DhcpStatus, *dhcpServer.Status)

	if dhcpServer.Network != nil {
		dhcpNetwork := dhcpServer.Network
		if dhcpNetwork.ID != nil {
			d.Set(Attr_DhcpNetworkDeprecated, *dhcpNetwork.ID)
			d.Set(Attr_DhcpNetworkID, *dhcpNetwork.ID)
		}
		if dhcpNetwork.Name != nil {
			d.Set(Attr_DhcpNetworkName, *dhcpNetwork.Name)
		}
	}

	if dhcpServer.Leases != nil {
		leaseList := make([]map[string]string, len(dhcpServer.Leases))
		for i, lease := range dhcpServer.Leases {
			leaseList[i] = map[string]string{
				Attr_DhcpLeaseInstanceIP:  *lease.InstanceIP,
				Attr_DhcpLeaseInstanceMac: *lease.InstanceMacAddress,
			}
		}
		d.Set(Attr_DhcpLeases, leaseList)
	}

	return nil
}
