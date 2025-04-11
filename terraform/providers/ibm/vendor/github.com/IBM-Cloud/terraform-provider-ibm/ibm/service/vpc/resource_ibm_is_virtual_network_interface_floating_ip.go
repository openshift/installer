// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsVirtualNetworkInterfaceFloatingIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsVirtualNetworkInterfaceFloatingIPCreate,
		ReadContext:   resourceIBMIsVirtualNetworkInterfaceFloatingIPRead,
		DeleteContext: resourceIBMIsVirtualNetworkInterfaceFloatingIPDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The virtual network interface identifier",
			},
			"floating_ip": &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
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
	}
}

func resourceIBMIsVirtualNetworkInterfaceFloatingIPCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createVirtualNetworkInterfaceFloatingIPOptions := &vpcv1.AddNetworkInterfaceFloatingIPOptions{}
	vniId := d.Get("virtual_network_interface").(string)
	fipId := d.Get("floating_ip").(string)
	createVirtualNetworkInterfaceFloatingIPOptions.VirtualNetworkInterfaceID = &vniId
	createVirtualNetworkInterfaceFloatingIPOptions.ID = &fipId

	floatingIP, response, err := sess.AddNetworkInterfaceFloatingIPWithContext(context, createVirtualNetworkInterfaceFloatingIPOptions)
	if err != nil {
		log.Printf("[DEBUG] AddNetworkInterfaceFloatingIPOptions failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("AddNetworkInterfaceFloatingIPOptions failed %s\n%s", err, response))
	}
	d.SetId(MakeTerraformVNIFloatingIpID(vniId, fipId))
	resourceIBMIsVirtualNetworkInterfaceFloatingIPGet(d, floatingIP)

	return nil
}

func resourceIBMIsVirtualNetworkInterfaceFloatingIPRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vniId, fipId, err := ParseVNIFloatingIpTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	getNetworkInterfaceFloatingIPOptions := &vpcv1.GetNetworkInterfaceFloatingIPOptions{}
	getNetworkInterfaceFloatingIPOptions.SetVirtualNetworkInterfaceID(vniId)
	getNetworkInterfaceFloatingIPOptions.SetID(fipId)

	floatingIP, response, err := sess.GetNetworkInterfaceFloatingIPWithContext(context, getNetworkInterfaceFloatingIPOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVirtualNetworkInterfaceFloatingIPWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVirtualNetworkInterfaceFloatingIPWithContext failed %s\n%s", err, response))
	}
	resourceIBMIsVirtualNetworkInterfaceFloatingIPGet(d, floatingIP)

	return nil
}
func resourceIBMIsVirtualNetworkInterfaceFloatingIPGet(d *schema.ResourceData, floatingIP *vpcv1.FloatingIPReference) diag.Diagnostics {
	if !core.IsNil(floatingIP.Name) {
		if err := d.Set("name", floatingIP.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if err := d.Set("address", floatingIP.Address); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting address: %s", err))
	}

	if err := d.Set("crn", floatingIP.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err := d.Set("href", floatingIP.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	deleted := make(map[string]interface{})

	if floatingIP.Deleted != nil && floatingIP.Deleted.MoreInfo != nil {
		deleted["more_info"] = floatingIP.Deleted
	}
	d.Set("deleted", []map[string]interface{}{deleted})

	return nil
}

func resourceIBMIsVirtualNetworkInterfaceFloatingIPDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	vniId, fipId, err := ParseVNIFloatingIpTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	removeNetworkInterfaceFloatingIPOptions := &vpcv1.RemoveNetworkInterfaceFloatingIPOptions{}

	removeNetworkInterfaceFloatingIPOptions.SetVirtualNetworkInterfaceID(vniId)
	removeNetworkInterfaceFloatingIPOptions.SetID(fipId)

	response, err := sess.RemoveNetworkInterfaceFloatingIPWithContext(context, removeNetworkInterfaceFloatingIPOptions)
	if err != nil {
		log.Printf("[DEBUG] RemoveNetworkInterfaceFloatingIPWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("RemoveNetworkInterfaceFloatingIPWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func MakeTerraformVNIFloatingIpID(id1, id2 string) string {
	// Include both virtual network interface id and floating ip id to create a unique Terraform id.  As a bonus,
	// we can extract the bare metal sever id as needed for API calls such as READ.
	return fmt.Sprintf("%s/%s", id1, id2)
}

func ParseVNIFloatingIpTerraformID(s string) (string, string, error) {
	segments := strings.Split(s, "/")
	if len(segments) != 2 {
		return "", "", fmt.Errorf("invalid terraform Id %s (incorrect number of segments for vitual network interface floating ip)", s)
	}
	if segments[0] == "" || segments[1] == "" {
		return "", "", fmt.Errorf("invalid terraform Id %s (one or more empty segments) for vitual network interface floating ip", s)
	}
	return segments[0], segments[1], nil
}
