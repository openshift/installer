// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkInterfacesRead,

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
			// Attributes
			Attr_Interfaces: {
				Computed:    true,
				Description: "Network interfaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CRN: {
							Computed:    true,
							Description: "The network interface's crn.",
							Type:        schema.TypeString,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The unique network interface ID.",
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
							Description: "The ip address of this network interface.",
							Type:        schema.TypeString,
						},
						Attr_MacAddress: {
							Computed:    true,
							Description: "The mac address of the network interface.",
							Type:        schema.TypeString,
						},
						Attr_Name: {
							Computed:    true,
							Description: "Name of the network interface (not unique or indexable).",
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
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPINetworkInterfacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()

	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	networkID := d.Get(Arg_NetworkID).(string)
	networkC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	networkInterfaces, err := networkC.GetAllNetworkInterfaces(networkID)
	if err != nil {
		return diag.FromErr(err)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	interfaces := []map[string]interface{}{}
	if len(networkInterfaces.Interfaces) > 0 {
		for _, netInterface := range networkInterfaces.Interfaces {
			interfaceMap := networkInterfaceToMap(netInterface, meta)
			interfaces = append(interfaces, interfaceMap)
		}
	}
	d.Set(Attr_Interfaces, interfaces)

	return nil
}

func networkInterfaceToMap(netInterface *models.NetworkInterface, meta interface{}) map[string]interface{} {
	interfaceMap := make(map[string]interface{})
	interfaceMap[Attr_ID] = netInterface.ID
	interfaceMap[Attr_IPAddress] = netInterface.IPAddress
	interfaceMap[Attr_MacAddress] = netInterface.MacAddress
	interfaceMap[Attr_Name] = netInterface.Name
	interfaceMap[Attr_NetworkSecurityGroupID] = netInterface.NetworkSecurityGroupID
	interfaceMap[Attr_NetworkSecurityGroupIDs] = netInterface.NetworkSecurityGroupIDs
	if netInterface.Instance != nil {
		pvmInstanceMap := pvmInstanceToMap(netInterface.Instance)
		interfaceMap[Attr_Instance] = []map[string]interface{}{pvmInstanceMap}
	}
	interfaceMap[Attr_Status] = netInterface.Status
	if netInterface.Crn != nil {
		interfaceMap[Attr_CRN] = netInterface.Crn
		userTags, err := flex.GetTagsUsingCRN(meta, string(*netInterface.Crn))
		if err != nil {
			log.Printf("Error on get of network interface (%s) user_tags: %s", *netInterface.ID, err)
		}
		interfaceMap[Attr_UserTags] = userTags
	}

	return interfaceMap
}
