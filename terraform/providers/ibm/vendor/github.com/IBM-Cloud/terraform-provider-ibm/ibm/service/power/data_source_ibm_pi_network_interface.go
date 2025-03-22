// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetworkInterface() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkInterfaceRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkID: {
				Description:  "Network ID.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkInterfaceID: {
				Description:  "Network interface ID.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The network interface's crn.",
				Type:        schema.TypeString,
			},
			Attr_Instance: {
				Computed:    true,
				Description: "The attached instance to this network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Href: {
							Computed:    true,
							Description: "Link to instance resource.",
							Type:        schema.TypeString,
						},
						Attr_InstanceID: {
							Computed:    true,
							Description: "The attached instance ID.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_IPAddress: {
				Computed:    true,
				Description: "The ip address of this Network Interface.",
				Type:        schema.TypeString,
			},
			Attr_MacAddress: {
				Computed:    true,
				Description: "The mac address of the Network Interface.",
				Type:        schema.TypeString,
			},
			Attr_Name: {
				Computed:    true,
				Description: "Name of the Network Interface (not unique or indexable).",
				Type:        schema.TypeString,
			},
			Attr_NetworkInterfaceID: {
				Computed:    true,
				Description: "ID of the network interface.",
				Type:        schema.TypeString,
			},
			Attr_NetworkSecurityGroupID: {
				Computed:    true,
				Deprecated:  "Deprecated, use network_security_group_ids instead.",
				Description: "ID of the network security group the network interface will be added to.",
				Type:        schema.TypeString,
			},
			Attr_NetworkSecurityGroupIDs: {
				Computed:    true,
				Description: "List of network security groups that the network interface is a member of.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeSet,
			},
			Attr_Status: {
				Computed:    true,
				Description: "The status of the network interface.",
				Type:        schema.TypeString,
			},
			Attr_UserTags: {
				Computed:    true,
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
		},
	}
}

func dataSourceIBMPINetworkInterfaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	networkID := d.Get(Arg_NetworkID).(string)
	networkInterfaceID := d.Get(Arg_NetworkInterfaceID).(string)
	networkC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	networkInterface, err := networkC.GetNetworkInterface(networkID, networkInterfaceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", networkID, *networkInterface.ID))
	d.Set(Attr_IPAddress, networkInterface.IPAddress)
	d.Set(Attr_MacAddress, networkInterface.MacAddress)
	d.Set(Attr_Name, networkInterface.Name)
	d.Set(Attr_NetworkInterfaceID, *networkInterface.ID)
	d.Set(Attr_NetworkSecurityGroupID, networkInterface.NetworkSecurityGroupID)
	d.Set(Attr_NetworkSecurityGroupIDs, networkInterface.NetworkSecurityGroupIDs)
	if networkInterface.Instance != nil {
		instance := []map[string]interface{}{}
		instanceMap := pvmInstanceToMap(networkInterface.Instance)
		instance = append(instance, instanceMap)
		d.Set(Attr_Instance, instance)
	}
	d.Set(Attr_Status, networkInterface.Status)
	if networkInterface.Crn != nil {
		d.Set(Attr_CRN, networkInterface.Crn)
		userTags, err := flex.GetTagsUsingCRN(meta, string(*networkInterface.Crn))
		if err != nil {
			log.Printf("Error on get of network interface (%s) user_tags: %s", *networkInterface.ID, err)
		}
		d.Set(Attr_UserTags, userTags)
	}

	return nil
}

func pvmInstanceToMap(pvm *models.NetworkInterfaceInstance) map[string]interface{} {
	instanceMap := make(map[string]interface{})
	if pvm.Href != "" {
		instanceMap[Attr_Href] = pvm.Href
	}
	if pvm.InstanceID != "" {
		instanceMap[Attr_InstanceID] = pvm.InstanceID
	}
	return instanceMap
}
