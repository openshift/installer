// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVirtualNetworkInterfaceIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIsVirtualNetworkInterfaceIPsRead,

		Schema: map[string]*schema.Schema{
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual network interface identifier",
			},
			"reserved_ips": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of reserved IPs for this virtual network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP.",
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
				},
			},
		},
	}
}

func dataSourceIsVirtualNetworkInterfaceIPsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listVirtualNetworkInterfaceIpsOptions := &vpcv1.ListVirtualNetworkInterfaceIpsOptions{}

	listVirtualNetworkInterfaceIpsOptions.SetVirtualNetworkInterfaceID(d.Get("virtual_network_interface").(string))
	var pager *vpcv1.VirtualNetworkInterfaceIpsPager
	pager, err = vpcClient.NewVirtualNetworkInterfaceIpsPager(listVirtualNetworkInterfaceIpsOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] VirtualNetworkInterfaceIpsPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("VirtualNetworkInterfaceIpsPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceVirtualNetworkInterfaceIPsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIsReservedIpsReservedIPToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("reserved_ips", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting reserved_ips for virtual network interface datasource %s", err))
	}

	return nil
}

// dataSourceVirtualNetworkInterfaceIPsID returns a reasonable ID for the list.
func dataSourceVirtualNetworkInterfaceIPsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsReservedIpsReservedIPCollectionFirstToMap(model *vpcv1.PageLink) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIBMIsReservedIpsReservedIPCollectionNextToMap(model *vpcv1.PageLink) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIBMIsReservedIpsReservedIPToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["reserved_ip"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}
