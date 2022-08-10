// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsInstanceNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsInstanceNetworkInterfacesRead,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance name.",
			},
			"network_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of network interfaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_ip_spoofing": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether source IP spoofing is allowed on this interface. If false, source IP spoofing is prevented on this interface. If true, source IP spoofing is allowed on this interface.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the network interface was created.",
						},
						"floating_ips": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The floating IPs associated with this network interface.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique IP address.",
									},
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this floating IP.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this floating IP.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this floating IP.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this floating IP.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network interface.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this network interface.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this network interface.",
						},
						"port_speed": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The network interface port speed in Mbps.",
						},
						"primary_ipv4_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The primary IPv4 address.",
						},
						isInstanceNicPrimaryIP: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The primary IP address to bind to the network interface. This can be specified using an existing reserved IP, or a prototype object for a new reserved IP.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isInstanceNicReservedIpAddress: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address to reserve, which must not already be reserved on the subnet.",
									},
									isInstanceNicReservedIpHref: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this reserved IP",
									},
									isInstanceNicReservedIpName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
									},
									isInstanceNicReservedIpId: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies a reserved IP by a unique property.",
									},
									isInstanceNicReservedIpResourceType: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type",
									},
								},
							},
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of security groups.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The security group's CRN.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The security group's canonical URL.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this security group.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this security group. Names must be unique within the VPC the security group resides in.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the network interface.",
						},
						"subnet": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The associated subnet.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this subnet.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this subnet.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this subnet.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this subnet.",
									},
								},
							},
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of this network interface as it relates to an instance.",
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIBMIsInstanceNetworkInterfacesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	instance_name := d.Get("instance_name").(string)
	listInstancesOptions := &vpcv1.ListInstancesOptions{}

	start := ""
	allrecs := []vpcv1.Instance{}
	for {

		if start != "" {
			listInstancesOptions.Start = &start
		}

		instances, response, err := vpcClient.ListInstances(listInstancesOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error Fetching Instances %s\n%s", err, response))
		}
		start = flex.GetNext(instances.Next)
		allrecs = append(allrecs, instances.Instances...)
		if start == "" {
			break
		}
	}

	ins_id := ""
	for _, instance := range allrecs {
		if *instance.Name == instance_name {
			ins_id = *instance.ID
			listInstanceNetworkInterfacesOptions := &vpcv1.ListInstanceNetworkInterfacesOptions{
				InstanceID: &ins_id,
			}
			networkInterfaceCollection, response, err := vpcClient.ListInstanceNetworkInterfacesWithContext(context, listInstanceNetworkInterfacesOptions)

			if err != nil {
				log.Printf("[DEBUG] ListSecurityGroupNetworkInterfacesWithContext failed %s\n%s", err, response)
				return diag.FromErr(fmt.Errorf("ListSecurityGroupNetworkInterfacesWithContext failed %s\n%s", err, response))
			}

			d.SetId(ins_id)

			if networkInterfaceCollection.NetworkInterfaces != nil {
				err = d.Set("network_interfaces", dataSourceNetworkInterfaceCollectionFlattenNetworkInterfaces(networkInterfaceCollection.NetworkInterfaces))
				if err != nil {
					return diag.FromErr(fmt.Errorf("[ERROR] Error setting network_interfaces %s", err))
				}
			}
			return nil
		}
	}

	return diag.FromErr(fmt.Errorf("Instance %s not found. %s", instance_name, err))
}

// dataSourceIBMIsInstanceNetworkInterfacesID returns a reasonable ID for the list.
func dataSourceIBMIsInstanceNetworkInterfacesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceNetworkInterfaceCollectionFlattenNetworkInterfaces(result []vpcv1.NetworkInterface) (networkInterfaces []map[string]interface{}) {
	for _, networkInterfacesItem := range result {
		networkInterfaces = append(networkInterfaces, dataSourceNetworkInterfaceCollectionNetworkInterfacesToMap(networkInterfacesItem))
	}

	return networkInterfaces
}

func dataSourceNetworkInterfaceCollectionNetworkInterfacesToMap(networkInterfacesItem vpcv1.NetworkInterface) (networkInterfacesMap map[string]interface{}) {
	networkInterfacesMap = map[string]interface{}{}

	if networkInterfacesItem.AllowIPSpoofing != nil {
		networkInterfacesMap["allow_ip_spoofing"] = networkInterfacesItem.AllowIPSpoofing
	}
	if networkInterfacesItem.CreatedAt != nil {
		networkInterfacesMap["created_at"] = networkInterfacesItem.CreatedAt.String()
	}
	if networkInterfacesItem.FloatingIps != nil {
		floatingIpsList := []map[string]interface{}{}
		for _, floatingIpsItem := range networkInterfacesItem.FloatingIps {
			floatingIpsList = append(floatingIpsList, dataSourceNetworkInterfaceFloatingIpsToMap(floatingIpsItem))
		}
		networkInterfacesMap["floating_ips"] = floatingIpsList
	}
	if networkInterfacesItem.Href != nil {
		networkInterfacesMap["href"] = networkInterfacesItem.Href
	}
	if networkInterfacesItem.ID != nil {
		networkInterfacesMap["id"] = networkInterfacesItem.ID
	}
	if networkInterfacesItem.Name != nil {
		networkInterfacesMap["name"] = networkInterfacesItem.Name
	}
	if networkInterfacesItem.PortSpeed != nil {
		networkInterfacesMap["port_speed"] = networkInterfacesItem.PortSpeed
	}
	if networkInterfacesItem.PrimaryIP != nil {
		networkInterfacesMap["primary_ipv4_address"] = networkInterfacesItem.PrimaryIP.Address
	}
	if networkInterfacesItem.PrimaryIP != nil {
		// reserved ip changes
		primaryIpList := make([]map[string]interface{}, 0)
		currentPrimIp := map[string]interface{}{}

		if networkInterfacesItem.PrimaryIP.Address != nil {
			currentPrimIp[isInstanceNicReservedIpAddress] = *networkInterfacesItem.PrimaryIP.Address
		}
		if networkInterfacesItem.PrimaryIP.Href != nil {
			currentPrimIp[isInstanceNicReservedIpHref] = *networkInterfacesItem.PrimaryIP.Href
		}
		if networkInterfacesItem.PrimaryIP.Name != nil {
			currentPrimIp[isInstanceNicReservedIpName] = *networkInterfacesItem.PrimaryIP.Name
		}
		if networkInterfacesItem.PrimaryIP.ID != nil {
			currentPrimIp[isInstanceNicReservedIpId] = *networkInterfacesItem.PrimaryIP.ID
		}
		if networkInterfacesItem.PrimaryIP.ResourceType != nil {
			currentPrimIp[isInstanceNicReservedIpResourceType] = *networkInterfacesItem.PrimaryIP.ResourceType
		}

		primaryIpList = append(primaryIpList, currentPrimIp)
		networkInterfacesMap[isInstanceNicPrimaryIP] = primaryIpList
	}
	if networkInterfacesItem.ResourceType != nil {
		networkInterfacesMap["resource_type"] = networkInterfacesItem.ResourceType
	}
	if networkInterfacesItem.SecurityGroups != nil {
		securityGroupsList := []map[string]interface{}{}
		for _, securityGroupsItem := range networkInterfacesItem.SecurityGroups {
			securityGroupsList = append(securityGroupsList, dataSourceNetworkInterfaceSecurityGroupsToMap(securityGroupsItem))
		}
		networkInterfacesMap["security_groups"] = securityGroupsList
	}
	if networkInterfacesItem.Status != nil {
		networkInterfacesMap["status"] = networkInterfacesItem.Status
	}
	if networkInterfacesItem.Subnet != nil {
		subnetList := []map[string]interface{}{}
		subnetMap := dataSourceNetworkInterfaceSubnetToMap(*networkInterfacesItem.Subnet)
		subnetList = append(subnetList, subnetMap)
		networkInterfacesMap["subnet"] = subnetList
	}
	if networkInterfacesItem.Type != nil {
		networkInterfacesMap["type"] = networkInterfacesItem.Type
	}

	return networkInterfacesMap
}
