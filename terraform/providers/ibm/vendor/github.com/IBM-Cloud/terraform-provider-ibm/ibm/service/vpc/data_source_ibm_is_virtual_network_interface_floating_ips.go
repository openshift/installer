// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVirtualNetworkInterfaceFloatingIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVirtualNetworkInterfaceFloatingIPsRead,

		Schema: map[string]*schema.Schema{
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual network interface identifier",
			},
			"floating_ips": {
				Type:        schema.TypeList,
				Description: "List of floating ips associated with the virtual network interface id",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The floating IP identifier",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this floating IP. The name is unique across all floating IPs in the region.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceIBMIsVirtualNetworkInterfaceFloatingIPsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vniId := d.Get("virtual_network_interface").(string)

	start := ""
	allrecs := []vpcv1.FloatingIPReference{}
	for {
		listNetworkInterfaceFloatingIpsOptions := &vpcv1.ListNetworkInterfaceFloatingIpsOptions{}
		listNetworkInterfaceFloatingIpsOptions.SetVirtualNetworkInterfaceID(vniId)
		if start != "" {
			listNetworkInterfaceFloatingIpsOptions.Start = &start
		}
		floatingIPCollection, response, err := sess.ListNetworkInterfaceFloatingIpsWithContext(context, listNetworkInterfaceFloatingIpsOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error ListNetworkInterfaceFloatingIpsWithContext %s\n%s", err, response))
		}
		start = flex.GetNext(floatingIPCollection.Next)
		allrecs = append(allrecs, floatingIPCollection.FloatingIps...)
		if start == "" {
			break
		}
	}
	floatingIpsInfo := make([]map[string]interface{}, 0)
	for _, floatingIP := range allrecs {
		l := map[string]interface{}{}
		l["id"] = *floatingIP.ID
		if !core.IsNil(floatingIP.Name) {
			l["name"] = floatingIP.Name
		}
		l["address"] = floatingIP.Address

		l["crn"] = floatingIP.CRN
		l["href"] = floatingIP.Href
		deleted := make(map[string]interface{})

		if floatingIP.Deleted != nil && floatingIP.Deleted.MoreInfo != nil {
			deleted["more_info"] = floatingIP.Deleted
		}
		l["deleted"] = []map[string]interface{}{deleted}
		floatingIpsInfo = append(floatingIpsInfo, l)
	}
	d.SetId(dataSourceIBMISVirtualNetworkInterfaceFloatingIPsID(d))
	d.Set("floating_ips", floatingIpsInfo)

	return nil
}

// dataSourceIBMISVirtualNetworkInterfaceFloatingIPsID returns a reasonable ID for a Virtual Network Interface FloatingIPs ID list.
func dataSourceIBMISVirtualNetworkInterfaceFloatingIPsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
