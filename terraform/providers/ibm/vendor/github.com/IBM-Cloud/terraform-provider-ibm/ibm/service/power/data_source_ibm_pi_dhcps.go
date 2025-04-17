// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
Datasource to get the list of dhcp servers in a power instance
*/

func DataSourceIBMPIDhcps() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPIDhcpServersRead,
		Schema: map[string]*schema.Schema{

			// Required Arguments
			Arg_CloudInstanceID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_DhcpServers: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of all the DHCP Servers",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_DhcpID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the DHCP Server",
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
				},
			},
		},
	}
}

func dataSourceIBMPIDhcpServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// session and client
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	// client
	client := st.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)

	// get all dhcp
	dhcpServers, err := client.GetAll()
	if err != nil {
		log.Printf("[DEBUG] get all DHCP failed %v", err)
		return diag.FromErr(err)
	}

	// set attributes
	servers := make([]map[string]interface{}, 0, len(dhcpServers))
	for _, dhcpServer := range dhcpServers {
		server := map[string]interface{}{
			Attr_DhcpID:     *dhcpServer.ID,
			Attr_DhcpStatus: *dhcpServer.Status,
		}
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
		servers = append(servers, server)
	}
	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_DhcpServers, servers)

	return nil
}
