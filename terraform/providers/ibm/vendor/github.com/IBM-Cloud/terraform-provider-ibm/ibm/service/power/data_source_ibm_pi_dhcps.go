// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Datasource to list dhcp servers in a power instance
func DataSourceIBMPIDhcps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIDhcpServersRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Servers: {
				Computed:    true,
				Description: "List of all the DHCP Servers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_DhcpID: {
							Computed:    true,
							Description: "ID of the DHCP Server.",
							Type:        schema.TypeString,
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
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPIDhcpServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)
	dhcpServers, err := client.GetAll()
	if err != nil {
		log.Printf("[DEBUG] get all DHCP failed %v", err)
		return diag.FromErr(err)
	}

	servers := make([]map[string]interface{}, 0, len(dhcpServers))
	for _, dhcpServer := range dhcpServers {
		server := map[string]interface{}{
			Attr_DhcpID: *dhcpServer.ID,
			Attr_Status: *dhcpServer.Status,
		}
		if dhcpServer.Network != nil {
			dhcpNetwork := dhcpServer.Network
			if dhcpNetwork.ID != nil {
				server[Attr_NetworkID] = *dhcpNetwork.ID
			}
			if dhcpNetwork.Name != nil {
				server[Attr_NetworkName] = *dhcpNetwork.Name
			}
		}
		servers = append(servers, server)
	}
	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_Servers, servers)

	return nil
}
