// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVirtualNetworkInterfaceIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIsVirtualNetworkInterfaceIPRead,

		Schema: map[string]*schema.Schema{
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual network interface identifier",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The reserved IP identifier.",
			},
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this reserved IP.",
			},
			"reserved_ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier for this reserved IP.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIsVirtualNetworkInterfaceIPRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getVirtualNetworkInterfaceIPOptions := &vpcv1.GetVirtualNetworkInterfaceIPOptions{}

	getVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(d.Get("virtual_network_interface").(string))
	getVirtualNetworkInterfaceIPOptions.SetID(d.Get("reserved_ip").(string))

	reservedIP, response, err := vpcClient.GetVirtualNetworkInterfaceIPWithContext(context, getVirtualNetworkInterfaceIPOptions)
	if err != nil {
		log.Printf("[DEBUG] GetVirtualNetworkInterfaceIPWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVirtualNetworkInterfaceIPWithContext failed for datasource %s\n%s", err, response))
	}

	d.SetId(*reservedIP.ID)

	if err = d.Set("address", reservedIP.Address); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting address: %s", err))
	}
	if err = d.Set("href", reservedIP.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("reserved_ip", reservedIP.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting reserved_ip: %s", err))
	}

	if err = d.Set("name", reservedIP.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("resource_type", reservedIP.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}
