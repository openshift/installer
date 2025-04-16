// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define all the constants that matches with the given terrafrom attribute
const (
	// Request Param Constants
	isBareMetalServerNICReservedIPLimit  = "limit"
	isBareMetalServerNICReservedIPSort   = "sort"
	isBareMetalServerNICReservedIPs      = "reserved_ips"
	isBareMetalServerNICReservedIPsCount = "total_count"
)

func DataSourceIBMISBareMetalServerNICReservedIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerNICReservedIPsRead,
		Schema: map[string]*schema.Schema{
			/*
				Request Parameters
				==================
				These are mandatory req parameters
			*/
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The BareMetalServer identifier.",
			},
			isBareMetalServerNicID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The BareMetalServer network interface identifier.",
			},
			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.
			*/

			isBareMetalServerNICReservedIPs: {
				Type:        schema.TypeList,
				Description: "Collection of all reserved IPs bound to a network interface.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerNicIpAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address",
						},
						isBareMetalServerNicIpAutoDelete: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If reserved ip shall be deleted automatically",
						},
						isBareMetalServerNICReservedIPCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the reserved IP was created.",
						},
						isBareMetalServerNICReservedIPhref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP.",
						},
						isBareMetalServerNicIpID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP",
						},
						isBareMetalServerNicIpName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined or system-provided name for this reserved IP.",
						},
						isBareMetalServerNICReservedIPOwner: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The owner of a reserved IP, defining whether it is managed by the user or the provider.",
						},
						isBareMetalServerNICReservedIPType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						isBareMetalServerNICReservedIPTarget: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reserved IP target id",
						},
					},
				},
			},
			isBareMetalServerNICReservedIPsCount: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages",
			},
		},
	}
}

func dataSourceIBMISBareMetalServerNICReservedIPsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	nicID := d.Get(isBareMetalServerNicID).(string)

	// Flatten all the reserved IPs
	allrecs := []vpcv1.ReservedIP{}
	options := &vpcv1.ListBareMetalServerNetworkInterfaceIpsOptions{
		BareMetalServerID:  &bareMetalServerID,
		NetworkInterfaceID: &nicID,
	}

	result, response, err := sess.ListBareMetalServerNetworkInterfaceIpsWithContext(context, options)
	if err != nil || response == nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error fetching reserved ips %s\n%s", err, response))
	}
	allrecs = append(allrecs, result.Ips...)

	// Now store all the reserved IP info with their response tags
	reservedIPs := []map[string]interface{}{}
	for _, data := range allrecs {
		ipsOutput := map[string]interface{}{}
		ipsOutput[isBareMetalServerNicIpAddress] = *data.Address
		ipsOutput[isBareMetalServerNicIpAutoDelete] = *data.AutoDelete
		ipsOutput[isBareMetalServerNICReservedIPCreatedAt] = (*data.CreatedAt).String()
		ipsOutput[isBareMetalServerNICReservedIPhref] = *data.Href
		ipsOutput[isBareMetalServerNicIpID] = *data.ID
		ipsOutput[isBareMetalServerNicIpName] = *data.Name
		ipsOutput[isBareMetalServerNICReservedIPOwner] = *data.Owner
		ipsOutput[isBareMetalServerNICReservedIPType] = *data.ResourceType
		target, ok := data.Target.(*vpcv1.ReservedIPTarget)
		if ok {
			ipsOutput[isReservedIPTarget] = target.ID
		}
		reservedIPs = append(reservedIPs, ipsOutput)
	}

	d.SetId(time.Now().UTC().String()) // This is not any reserved ip or BareMetalServer id but state id
	d.Set(isBareMetalServerNICReservedIPs, reservedIPs)
	d.Set(isBareMetalServerNICReservedIPsCount, len(reservedIPs))
	d.Set(isBareMetalServerID, bareMetalServerID)
	d.Set(isBareMetalServerNicID, nicID)
	return nil
}
