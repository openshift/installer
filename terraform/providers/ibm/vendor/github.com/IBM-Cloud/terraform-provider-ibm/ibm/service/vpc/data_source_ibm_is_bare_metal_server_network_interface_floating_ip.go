// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerNetworkInterfaceFloatingIPID = "floating_ip"
)

func DataSourceIBMIsBareMetalServerNetworkInterfaceFloatingIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerNetworkInterfaceFloatingIPRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier",
			},
			isBareMetalServerNetworkInterface: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network interface identifier of bare metal server",
			},
			isBareMetalServerNetworkInterfaceFloatingIPID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The floating ip identifier of the network interface associated with the bare metal server",
			},
			floatingIPName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the floating IP",
			},

			floatingIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP address",
			},

			floatingIPStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP status",
			},

			floatingIPZone: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone name",
			},

			floatingIPTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target info",
			},

			floatingIPCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP crn",
			},
		},
	}
}

func dataSourceIBMISBareMetalServerNetworkInterfaceFloatingIPRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	nicID := d.Get(isBareMetalServerNetworkInterface).(string)
	fipID := d.Get(isBareMetalServerNetworkInterfaceFloatingIPID).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.GetBareMetalServerNetworkInterfaceFloatingIPOptions{
		BareMetalServerID:  &bareMetalServerID,
		NetworkInterfaceID: &nicID,
		ID:                 &fipID,
	}

	ip, response, err := sess.GetBareMetalServerNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error fetching floating IP for bare metal server %s\n%s", err, response))
	}
	d.Set(floatingIPName, *ip.Name)
	d.Set(floatingIPAddress, *ip.Address)
	d.Set(floatingIPStatus, *ip.Status)
	d.Set(floatingIPZone, *ip.Zone.Name)

	d.Set(floatingIPCRN, *ip.CRN)

	target, ok := ip.Target.(*vpcv1.FloatingIPTarget)
	if ok {
		d.Set(floatingIPTarget, target.ID)
	}

	d.SetId(*ip.ID)

	return nil
}
