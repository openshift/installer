// Copyright IBM Corp. 2023 All Rights VirtualNetworkInterface.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsVirtualNetworkInterfaceIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsVirtualNetworkInterfaceIPCreate,
		ReadContext:   resourceIBMIsVirtualNetworkInterfaceIPRead,
		DeleteContext: resourceIBMIsVirtualNetworkInterfaceIPDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The virtual network interface identifier.",
			},
			"reserved_ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The reserved ip identifier.",
			},
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
			},

			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this reserved IP.",
			},

			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}
func resourceIBMIsVirtualNetworkInterfaceIPCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	addVirtualNetworkInterfaceIPOptions := &vpcv1.AddVirtualNetworkInterfaceIPOptions{}

	addVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(d.Get("virtual_network_interface").(string))
	addVirtualNetworkInterfaceIPOptions.SetID(d.Get("reserved_ip").(string))

	reservedIP, response, err := vpcClient.AddVirtualNetworkInterfaceIPWithContext(context, addVirtualNetworkInterfaceIPOptions)
	if err != nil {
		log.Printf("[DEBUG] AddVirtualNetworkInterfaceIPWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("AddVirtualNetworkInterfaceIPWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *addVirtualNetworkInterfaceIPOptions.VirtualNetworkInterfaceID, *reservedIP.ID))

	return resourceIBMIsVirtualNetworkInterfaceIPRead(context, d, meta)
}

func resourceIBMIsVirtualNetworkInterfaceIPRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getVirtualNetworkInterfaceIPOptions := &vpcv1.GetVirtualNetworkInterfaceIPOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(parts[0])
	getVirtualNetworkInterfaceIPOptions.SetID(parts[1])

	reservedIP, response, err := vpcClient.GetVirtualNetworkInterfaceIPWithContext(context, getVirtualNetworkInterfaceIPOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVirtualNetworkInterfaceIPWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVirtualNetworkInterfaceIPWithContext failed for resource %s\n%s", err, response))
	}

	if !core.IsNil(reservedIP.Address) {
		if err = d.Set("address", reservedIP.Address); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting address: %s", err))
		}
	}
	if !core.IsNil(reservedIP.Name) {
		if err = d.Set("name", reservedIP.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if err = d.Set("href", reservedIP.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("resource_type", reservedIP.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	return nil
}

func resourceIBMIsVirtualNetworkInterfaceIPDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	removeVirtualNetworkInterfaceIPOptions := &vpcv1.RemoveVirtualNetworkInterfaceIPOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	removeVirtualNetworkInterfaceIPOptions.SetVirtualNetworkInterfaceID(parts[0])
	removeVirtualNetworkInterfaceIPOptions.SetID(parts[1])

	response, err := vpcClient.RemoveVirtualNetworkInterfaceIPWithContext(context, removeVirtualNetworkInterfaceIPOptions)
	if err != nil {
		log.Printf("[DEBUG] RemoveVirtualNetworkInterfaceIPWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("RemoveVirtualNetworkInterfaceIPWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
