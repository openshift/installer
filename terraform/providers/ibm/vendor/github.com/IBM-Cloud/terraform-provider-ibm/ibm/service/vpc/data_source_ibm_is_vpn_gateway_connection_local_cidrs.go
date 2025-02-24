// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVPNGatewayConnectionLocalCidrs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNGatewayConnectionLocalCidrsRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN gateway identifier.",
			},
			"vpn_gateway_connection": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN gateway connection identifier.",
			},
			"cidrs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The CIDRs for this resource.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIBMIsVPNGatewayConnectionLocalCidrsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listVPNGatewayConnectionsLocalCidrsOptions := &vpcv1.ListVPNGatewayConnectionsLocalCIDRsOptions{}

	listVPNGatewayConnectionsLocalCidrsOptions.SetVPNGatewayID(d.Get("vpn_gateway").(string))
	listVPNGatewayConnectionsLocalCidrsOptions.SetID(d.Get("vpn_gateway_connection").(string))

	vpnGatewayConnectionCidRs, response, err := vpcClient.ListVPNGatewayConnectionsLocalCIDRsWithContext(context, listVPNGatewayConnectionsLocalCidrsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListVPNGatewayConnectionsLocalCidrsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListVPNGatewayConnectionsLocalCidrsWithContext failed %s\n%s", err, response))
	}
	d.SetId(dataSourceIBMIsVPNGatewayConnectionLocalCidrsID(d))
	d.Set("cidrs", vpnGatewayConnectionCidRs.CIDRs)

	return nil
}

// dataSourceIBMIsVPNGatewayConnectionLocalCidrsID returns a reasonable ID for the list.
func dataSourceIBMIsVPNGatewayConnectionLocalCidrsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
