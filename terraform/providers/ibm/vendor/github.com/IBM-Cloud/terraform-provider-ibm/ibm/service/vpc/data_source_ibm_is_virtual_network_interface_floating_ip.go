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

func DataSourceIBMIsVirtualNetworkInterfaceFloatingIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVirtualNetworkInterfaceFloatingIPRead,

		Schema: map[string]*schema.Schema{
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual network interface identifier",
			},
			"floating_ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The floating IP identifier",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this floating IP. The name is unique across all floating IPs in the region.",
			},

			"deleted": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Link to documentation about deleted resources",
						},
					},
				},
			},
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique IP address.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this floating IP.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this floating IP.",
			},
		},
	}
}

func dataSourceIBMIsVirtualNetworkInterfaceFloatingIPRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vniId := d.Get("virtual_network_interface").(string)
	fipId := d.Get("floating_ip").(string)

	getNetworkInterfaceFloatingIPOptions := &vpcv1.GetNetworkInterfaceFloatingIPOptions{}
	getNetworkInterfaceFloatingIPOptions.SetVirtualNetworkInterfaceID(vniId)
	getNetworkInterfaceFloatingIPOptions.SetID(fipId)

	floatingIP, response, err := sess.GetNetworkInterfaceFloatingIPWithContext(context, getNetworkInterfaceFloatingIPOptions)
	if err != nil {
		log.Printf("[DEBUG] GetVirtualNetworkInterfaceFloatingIPWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVirtualNetworkInterfaceFloatingIPWithContext failed %s\n%s", err, response))
	}
	d.SetId(*floatingIP.ID)
	resourceIBMIsVirtualNetworkInterfaceFloatingIPGet(d, floatingIP)
	return nil
}
