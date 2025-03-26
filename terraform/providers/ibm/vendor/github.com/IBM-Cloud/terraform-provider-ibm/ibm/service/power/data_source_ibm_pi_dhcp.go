// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIDhcp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDhcpRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_DhcpID: {
				Description: "ID of the DHCP Server.",
				Required:    true,
				Type:        schema.TypeString,
			},

			// Attributes
			Attr_DhcpID: {
				Computed:    true,
				Deprecated:  "The field is deprecated, use pi_dhcp_id instead.",
				Description: "ID of the DHCP Server.",
				Type:        schema.TypeString,
			},
			Attr_Leases: {
				Computed:    true,
				Description: "List of DHCP Server PVM Instance leases.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_InstanceIP: {
							Computed:    true,
							Description: "IP of the PVM Instance.",
							Type:        schema.TypeString,
						},
						Attr_InstanceMac: {
							Computed:    true,
							Description: "MAC Address of the PVM Instance.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_NetworkID: {
				Computed:    true,
				Description: "ID of the DHCP Server private network.",
				Type:        schema.TypeString,
			},
			Attr_NetworkName: {
				Computed:    true,
				Description: "Name of the DHCP Server private network.",
				Type:        schema.TypeString,
			},
			Attr_Status: {
				Computed:    true,
				Description: "Status of the DHCP Server.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIDhcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	dhcpID := d.Get(Arg_DhcpID).(string)
	client := instance.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)
	dhcpServer, err := client.Get(dhcpID)
	if err != nil {
		log.Printf("[DEBUG] get DHCP failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *dhcpServer.ID))
	d.Set(Attr_DhcpID, *dhcpServer.ID)
	d.Set(Attr_Status, *dhcpServer.Status)

	if dhcpServer.Network != nil {
		dhcpNetwork := dhcpServer.Network
		if dhcpNetwork.ID != nil {
			d.Set(Attr_NetworkID, *dhcpNetwork.ID)
		}
		if dhcpNetwork.Name != nil {
			d.Set(Attr_NetworkName, *dhcpNetwork.Name)
		}
	}

	if dhcpServer.Leases != nil {
		leaseList := make([]map[string]string, len(dhcpServer.Leases))
		for i, lease := range dhcpServer.Leases {
			leaseList[i] = map[string]string{
				Attr_InstanceIP:  *lease.InstanceIP,
				Attr_InstanceMac: *lease.InstanceMacAddress,
			}
		}
		d.Set(Attr_Leases, leaseList)
	}

	return nil
}
