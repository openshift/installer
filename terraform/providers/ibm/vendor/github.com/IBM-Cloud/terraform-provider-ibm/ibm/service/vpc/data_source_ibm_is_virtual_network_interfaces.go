// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-beta-go-sdk/vpcbetav1"
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
		},
	}
}

func dataSourceIBMIsVirtualNetworkInterfacesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcbetaClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	listVirtualNetworkInterfacesOptions := &vpcbetav1.ListVirtualNetworkInterfacesOptions{}

	vniCollection, response, err := vpcbetaClient.ListVirtualNetworkInterfacesWithContext(context, listVirtualNetworkInterfacesOptions)
	//log.Printf(len(vniCollection.VirtualNetworkInterfaces))
	if err != nil {
		log.Printf("[DEBUG] VirtualNetworkInterfacesPager.GetAll() failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("VirtualNetworkInterfacesPager.GetAll() failed %s\n%s", err, response))
	}
	// var pager *vpcbetav1.VirtualNetworkInterfacesPager
	// pager, err = vpcbetaClient.NewVirtualNetworkInterfacesPager(listVirtualNetworkInterfacesOptions)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// allItems, err := pager.GetAll()
	// if err != nil {
	// 	log.Printf("[DEBUG] VirtualNetworkInterfacesPager.GetAll() failed %s", err)
	// 	return diag.FromErr(fmt.Errorf("VirtualNetworkInterfacesPager.GetAll() failed %s", err))
	// }

	d.SetId(dataSourceIBMIsVirtualNetworkInterfacesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range vniCollection.VirtualNetworkInterfaces {
		modelMap, err := dataSourceIBMIsVirtualNetworkInterfacesVirtualNetworkInterfaceToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
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

func dataSourceIBMIsVirtualNetworkInterfacesVirtualNetworkInterfaceToMap(model *vpcbetav1.VirtualNetworkInterface) (map[string]interface{}, error) {
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
	return modelMap, nil
}
