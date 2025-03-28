// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"net"
	"strconv"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPIInstanceIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIInstancesIPRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_InstanceName: {
				Description:  "The unique identifier or name of the instance.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkName: {
				Description:  "The subnet that the instance belongs to.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_ExternalIP: {
				Computed:    true,
				Description: "The external IP of the network that is attached to this instance.",
				Type:        schema.TypeString,
			},
			Attr_IP: {
				Computed:    true,
				Description: "The IP address that is attached to this instance from the subnet.",
				Type:        schema.TypeString,
			},
			Attr_IPOctet: {
				Computed:    true,
				Description: "The IP octet of the network that is attached to this instance.",
				Type:        schema.TypeString,
			},
			Attr_MacAddress: {
				Computed:    true,
				Description: "The MAC address of the network that is attached to this instance.",
				Type:        schema.TypeString,
			},
			Attr_Macaddress: {
				Computed:    true,
				Deprecated:  "Deprecated, use mac_address instead",
				Description: "The MAC address of the network that is attached to this instance.",
				Type:        schema.TypeString,
			},
			Attr_NetworkID: {
				Computed:    true,
				Description: "ID of the network.",
				Type:        schema.TypeString,
			},
			Attr_NetworkInterfaceID: {
				Computed:    true,
				Description: "ID of the network interface.",
				Type:        schema.TypeString,
			},
			Attr_NetworkSecurityGroupIDs: {
				Computed:    true,
				Description: "IDs of the network necurity groups that the network interface is a member of.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeSet,
			},
			Attr_NetworkSecurityGroupsHref: {
				Computed:    true,
				Description: "Links to the network security groups that the network interface is a member of.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeList,
			},
			Attr_Type: {
				Computed:    true,
				Description: "The type of the network that is attached to this instance.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIInstancesIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	networkName := d.Get(Arg_NetworkName).(string)
	powerC := instance.NewIBMPIInstanceClient(ctx, sess, cloudInstanceID)

	powervmdata, err := powerC.Get(d.Get(Arg_InstanceName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	for _, network := range powervmdata.Networks {
		if network.NetworkName == networkName {
			d.SetId(network.NetworkID)
			d.Set(Attr_ExternalIP, network.ExternalIP)
			d.Set(Attr_IP, network.IPAddress)
			d.Set(Attr_MacAddress, network.MacAddress)
			d.Set(Attr_Macaddress, network.MacAddress)
			d.Set(Attr_NetworkID, network.NetworkID)
			d.Set(Attr_NetworkInterfaceID, network.NetworkInterfaceID)
			d.Set(Attr_Type, network.Type)

			IPObject := net.ParseIP(network.IPAddress).To4()
			if len(IPObject) > 0 {
				d.Set(Attr_IPOctet, strconv.Itoa(int(IPObject[3])))
			}
			if len(network.NetworkSecurityGroupIDs) > 0 {
				d.Set(Attr_NetworkSecurityGroupIDs, network.NetworkSecurityGroupIDs)
			}
			if len(network.NetworkSecurityGroupsHref) > 0 {
				d.Set(Attr_NetworkSecurityGroupsHref, network.NetworkSecurityGroupsHref)
			}
			return nil
		}
	}

	return diag.Errorf("failed to find instance ip that belongs to the given network")
}
