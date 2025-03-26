// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func DataSourceIBMIsVirtualNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVirtualNetworkInterfacesRead,

		Schema: map[string]*schema.Schema{
			"virtual_network_interfaces": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of virtual network interfaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						// vni p2 changes
						"allow_ip_spoofing": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether source IP spoofing is allowed on this interface. If `false`, source IP spoofing is prevented on this interface. If `true`, source IP spoofing is allowed on this interface.",
						},
						"tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "UserTags for the vni instance",
						},
						"access_tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "Access management tags for the vni instance",
						},
						"enable_infrastructure_nat": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the virtual network interface,  allowing the workload to perform any needed NAT operations.- `allow_ip_spoofing` must be `false`.- If the virtual network interface is attached:  - The target `resource_type` must be `bare_metal_server_network_attachment`.  - The target `interface_type` must not be `hipersocket`.",
						},
						"ips": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
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
						"mac_address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC address of the virtual network interface. May be absent if `lifecycle_state` is `pending`.",
						},

						"auto_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this virtual network interface will be automatically deleted when`target` is deleted.",
						},

						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the virtual network interface was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this virtual network interface.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual network interface.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this virtual network interface.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the virtual network interface.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
						},
						"primary_ip": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reserved IP for this virtual network interface.May be absent when `lifecycle_state` is `pending`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
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
						"protocol_state_filtering_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol state filtering mode used for this virtual network interface.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this virtual network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this resource group.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"security_groups": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The security groups for this virtual network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The security group's CRN.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The security group's canonical URL.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this security group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this security group. The name is unique across all security groups for the VPC.",
									},
								},
							},
						},
						"subnet": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The associated subnet.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this subnet.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this subnet.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this subnet.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this subnet. The name is unique across all subnets in the VPC.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"target": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The target of this virtual network interface.If absent, this virtual network interface is not attached to a target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this share mount target.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this share mount target.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this share mount target. The name is unique across all targets for the file share.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"vpc": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC this virtual network interface resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this VPC.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this VPC. The name is unique across all VPCs in the region.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"zone": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The zone this virtual network interface resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this zone.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this zone.",
									},
								},
							},
						},
					},
				},
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the resource group these virtual network interfaces belong to",
			},
		},
	}
}

func dataSourceIBMIsVirtualNetworkInterfacesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listVirtualNetworkInterfacesOptions := &vpcv1.ListVirtualNetworkInterfacesOptions{}
	if resgroupintf, ok := d.GetOk("resource_group"); ok {
		resGroup := resgroupintf.(string)
		listVirtualNetworkInterfacesOptions.ResourceGroupID = &resGroup
	}
	var pager *vpcv1.VirtualNetworkInterfacesPager
	pager, err = vpcClient.NewVirtualNetworkInterfacesPager(listVirtualNetworkInterfacesOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] VirtualNetworkInterfacesPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("VirtualNetworkInterfacesPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIBMIsVirtualNetworkInterfacesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIsVirtualNetworkInterfacesVirtualNetworkInterfaceToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}

		tags, err := flex.GetGlobalTagsUsingCRN(meta, *modelItem.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on get of datasources vni (%s) tags: %s", *modelItem.ID, err)
		}

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *modelItem.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on get of datasources vni (%s) access tags: %s", *modelItem.ID, err)
		}

		modelMap["tags"] = tags
		modelMap["access_tags"] = accesstags
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("virtual_network_interfaces", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting virtual_network_interfaces %s", err))
	}

	return nil
}

// dataSourceIBMIsVirtualNetworkInterfacesID returns a reasonable ID for the list.
func dataSourceIBMIsVirtualNetworkInterfacesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsVirtualNetworkInterfacesVirtualNetworkInterfaceToMap(model *vpcv1.VirtualNetworkInterface) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoDelete != nil {
		modelMap["auto_delete"] = *model.AutoDelete
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = flex.DateTimeToString(model.CreatedAt)
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.LifecycleState != nil {
		modelMap["lifecycle_state"] = *model.LifecycleState
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.PrimaryIP != nil {
		primaryIPMap, err := dataSourceIBMIsVirtualNetworkInterfaceReservedIPReferenceToMap(model.PrimaryIP)
		if err != nil {
			return modelMap, err
		}
		modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}
	}
	if model.ProtocolStateFilteringMode != nil {
		modelMap["protocol_state_filtering_mode"] = *model.ProtocolStateFilteringMode
	}
	if model.ResourceGroup != nil {
		resourceGroupMap, err := dataSourceIBMIsVirtualNetworkInterfaceResourceGroupReferenceToMap(model.ResourceGroup)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	if model.SecurityGroups != nil {
		securityGroups := []map[string]interface{}{}
		for _, securityGroupsItem := range model.SecurityGroups {
			securityGroupsItemMap, err := dataSourceIBMIsVirtualNetworkInterfaceSecurityGroupReferenceToMap(&securityGroupsItem)
			if err != nil {
				return modelMap, err
			}
			securityGroups = append(securityGroups, securityGroupsItemMap)
		}
		modelMap["security_groups"] = securityGroups
	}
	if model.Subnet != nil {
		subnetMap, err := dataSourceIBMIsVirtualNetworkInterfaceSubnetReferenceToMap(model.Subnet)
		if err != nil {
			return modelMap, err
		}
		modelMap["subnet"] = []map[string]interface{}{subnetMap}
	}
	if model.Target != nil {
		targetMap, err := dataSourceIBMIsVirtualNetworkInterfaceVirtualNetworkInterfaceTargetToMap(model.Target)
		if err != nil {
			return modelMap, err
		}
		modelMap["target"] = []map[string]interface{}{targetMap}
	}
	if model.VPC != nil {
		vpcMap, err := dataSourceIBMIsVirtualNetworkInterfaceVPCReferenceToMap(model.VPC)
		if err != nil {
			return modelMap, err
		}
		modelMap["vpc"] = []map[string]interface{}{vpcMap}
	}
	if model.Zone != nil {
		zoneMap, err := dataSourceIBMIsVirtualNetworkInterfaceZoneReferenceToMap(model.Zone)
		if err != nil {
			return modelMap, err
		}
		modelMap["zone"] = []map[string]interface{}{zoneMap}
	}
	// vni p2 changes

	modelMap["mac_address"] = model.MacAddress
	modelMap["allow_ip_spoofing"] = model.AllowIPSpoofing
	modelMap["enable_infrastructure_nat"] = model.EnableInfrastructureNat

	ips := []map[string]interface{}{}
	if model.Ips != nil {
		for _, modelItem := range model.Ips {
			if *modelItem.ID != *model.PrimaryIP.ID {
				modelMap, err := dataSourceIBMIsVirtualNetworkInterfaceReservedIPReferenceToMap(&modelItem)
				if err != nil {
					return modelMap, err
				}
				ips = append(ips, modelMap)
			}
		}
	}
	modelMap["ips"] = ips

	return modelMap, nil
}
