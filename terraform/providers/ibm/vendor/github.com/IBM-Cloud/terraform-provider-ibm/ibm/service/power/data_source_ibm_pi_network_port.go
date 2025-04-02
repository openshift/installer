// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetworkPort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkPortsRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkName: {
				Description:  "The unique identifier or name of a network.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_NetworkPorts: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Description: {
							Computed:    true,
							Description: "The description for the network port.",
							Type:        schema.TypeString,
						},
						Attr_Href: {
							Computed:    true,
							Description: "Network port href.",
							Type:        schema.TypeString,
						},
						Attr_IPaddress: {
							Computed:    true,
							Description: "The IP address of the port.",
							Type:        schema.TypeString,
						},
						Attr_Macaddress: {
							Computed:    true,
							Deprecated:  "Deprecated, use mac_address instead",
							Description: "The MAC address of the port.",
							Type:        schema.TypeString,
						},
						Attr_MacAddress: {
							Computed:    true,
							Description: "The MAC address of the port.",
							Type:        schema.TypeString,
						},
						Attr_PortID: {
							Computed:    true,
							Description: "The ID of the port.",
							Type:        schema.TypeString,
						},
						Attr_PublicIP: {
							Computed:    true,
							Description: "The public IP associated with the port.",
							Type:        schema.TypeString,
						},
						Attr_Status: {
							Computed:    true,
							Description: "The status of the port.",
							Type:        schema.TypeString,
						},
					},
				},
			},
		},
		DeprecationMessage: "Data source ibm_pi_network_port_attach is deprecated. Use `ibm_pi_network_interface` data source instead.",
	}
}

func dataSourceIBMPINetworkPortsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)

	networkportC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	networkportdata, err := networkportC.GetAllPorts(d.Get(helpers.PINetworkName).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)

	d.Set(Attr_NetworkPorts, flattenNetworkPorts(networkportdata.Ports))

	return nil
}

func flattenNetworkPorts(networkPorts []*models.NetworkPort) interface{} {
	result := make([]map[string]interface{}, 0, len(networkPorts))
	log.Printf("the number of ports is %d", len(networkPorts))
	for _, i := range networkPorts {
		l := map[string]interface{}{
			Attr_Description: i.Description,
			Attr_Href:        i.Href,
			Attr_IPaddress:   *i.IPAddress,
			Attr_Macaddress:  *i.MacAddress,
			Attr_MacAddress:  *i.MacAddress,
			Attr_PortID:      *i.PortID,
			Attr_PublicIP:    i.ExternalIP,
			Attr_Status:      *i.Status,
		}
		result = append(result, l)
	}
	return result
}
