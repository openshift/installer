// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIsSecurityGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsSecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "resource group identifier.",
			},
			"vpc_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vpc identifier.",
			},
			"vpc_crn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vpc crn",
			},
			"vpc_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vpc name.",
			},
			"security_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of security groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this security group was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group's CRN.",
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
							Description: "The user-defined name for this security group. Names must be unique within the VPC the security group resides in.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this security group.",
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
										Description: "The user-defined name for this resource group.",
									},
								},
							},
						},
						"rules": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The rules for this security group. If no rules exist, all traffic will be denied.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"direction": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The direction of traffic to enforce, either `inbound` or `outbound`.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this security group rule.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this security group rule.",
									},
									"ip_version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP version to enforce. The format of `remote.address` or `remote.cidr_block` must match this property, if they are used. Alternatively, if `remote` references a security group, then this rule only applies to IP addresses (network interfaces) in that group matching this IP version.",
									},
									"protocol": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The protocol to enforce.",
									},
									"local": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The local IP address or range of local IP addresses to which this rule will allow inbound traffic (or from which, for outbound traffic). A CIDR block of 0.0.0.0/0 allows traffic to all local IP addresses (or from all local IP addresses, for outbound rules).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												"cidr_block": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CIDR block. This property may add support for IPv6 CIDR blocks in the future. When processing a value in this property, verify that the CIDR block is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected CIDR block format was encountered.",
												},
											},
										},
									},
									"remote": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The IP addresses or security groups from which this rule allows traffic (or to which,for outbound rules). Can be specified as an IP address, a CIDR block, or a securitygroup. A CIDR block of `0.0.0.0/0` allows traffic from any source (or to any source,for outbound rules).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												"cidr_block": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CIDR block. This property may add support for IPv6 CIDR blocks in the future. When processing a value in this property, verify that the CIDR block is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected CIDR block format was encountered.",
												},
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The security group's CRN.",
												},
												"deleted": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
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
													Description: "The user-defined name for this security group. Names must be unique within the VPC the security group resides in.",
												},
											},
										},
									},
									"code": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ICMP traffic code to allow.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ICMP traffic type to allow.",
									},
									"port_max": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive upper bound of TCP/UDP port range.",
									},
									"port_min": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The inclusive lower bound of TCP/UDP port range.",
									},
								},
							},
						},
						"targets": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The targets for this security group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
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
										Description: "The URL for this network interface.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this network interface.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this network interface.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The load balancer's CRN.",
									},
								},
							},
						},
						"vpc": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC this security group is a part of.",
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
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
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
										Description: "The unique user-defined name for this VPC.",
									},
								},
							},
						},

						isSecurityGroupAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access management tags",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsSecurityGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	resourceGrp := d.Get("resource_group").(string)
	vpcId := d.Get("vpc_id").(string)
	vpcCrn := d.Get("vpc_crn").(string)
	vpcName := d.Get("vpc_name").(string)

	start := ""
	allrecs := []vpcv1.SecurityGroup{}
	listSecurityGroupsOptions := &vpcv1.ListSecurityGroupsOptions{}
	if resourceGrp != "" {
		listSecurityGroupsOptions.ResourceGroupID = &resourceGrp
	}
	if vpcId != "" {
		listSecurityGroupsOptions.VPCID = &vpcId
	}
	if vpcCrn != "" {
		listSecurityGroupsOptions.VPCCRN = &vpcCrn
	}
	if vpcName != "" {
		listSecurityGroupsOptions.VPCName = &vpcName
	}
	for {

		if start != "" {
			listSecurityGroupsOptions.Start = &start
		}
		securityGroupCollection, response, err := vpcClient.ListSecurityGroupsWithContext(context, listSecurityGroupsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListSecurityGroupsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListSecurityGroupsWithContext failed %s\n%s", err, response))
		}

		start = flex.GetNext(securityGroupCollection.Next)
		allrecs = append(allrecs, securityGroupCollection.SecurityGroups...)

		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsSecurityGroupsID(d))
	err = d.Set("security_groups", dataSourceSecurityGroupCollectionFlattenSecurityGroups(allrecs, d, meta))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error setting security_groups %s", err))
	}
	return nil
}

// dataSourceIBMIsSecurityGroupsID returns a reasonable ID for the list.
func dataSourceIBMIsSecurityGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceSecurityGroupCollectionFlattenSecurityGroups(modelSlice []vpcv1.SecurityGroup, d *schema.ResourceData, meta interface{}) (result []map[string]interface{}) {
	result = []map[string]interface{}{}
	for _, securityGroupsItem := range modelSlice {
		mapItem := dataSourceSecurityGroupCollectionSecurityGroupsToMap(&securityGroupsItem, d, meta)
		result = append(result, mapItem)
	}
	return result
}

func dataSourceSecurityGroupCollectionSecurityGroupsToMap(securityGroupsItem *vpcv1.SecurityGroup, d *schema.ResourceData, meta interface{}) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if securityGroupsItem.CreatedAt != nil {
		resultMap["created_at"] = securityGroupsItem.CreatedAt.String()
	}
	if securityGroupsItem.CRN != nil {
		resultMap["crn"] = securityGroupsItem.CRN
	}
	if securityGroupsItem.Href != nil {
		resultMap["href"] = securityGroupsItem.Href
	}
	if securityGroupsItem.ID != nil {
		resultMap["id"] = securityGroupsItem.ID
	}
	if securityGroupsItem.Name != nil {
		resultMap["name"] = securityGroupsItem.Name
	}
	if securityGroupsItem.ResourceGroup != nil {
		var mapSlice []map[string]interface{}
		modelMap := dataSourceSecurityGroupCollectionSecurityGroupsResourceGroupToMap(securityGroupsItem.ResourceGroup)
		mapSlice = append(mapSlice, modelMap)
		resultMap["resource_group"] = mapSlice
	}
	if securityGroupsItem.Rules != nil {
		var mapSlice []map[string]interface{}
		for _, listElem := range securityGroupsItem.Rules {
			mapElem := dataSourceSecurityGroupCollectionSecurityGroupsRulesToMap(listElem)
			mapSlice = append(mapSlice, mapElem)
		}
		resultMap["rules"] = mapSlice
	}
	if securityGroupsItem.Targets != nil {
		var mapSlice []map[string]interface{}
		for _, listElem := range securityGroupsItem.Targets {
			mapElem := dataSourceSecurityGroupCollectionSecurityGroupsTargetsToMap(listElem)
			mapSlice = append(mapSlice, mapElem)
		}
		resultMap["targets"] = mapSlice
	}
	if securityGroupsItem.VPC != nil {
		var mapSlice []map[string]interface{}
		modelMap := dataSourceSecurityGroupCollectionSecurityGroupsVPCToMap(securityGroupsItem.VPC)
		mapSlice = append(mapSlice, modelMap)
		resultMap["vpc"] = mapSlice
	}
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *securityGroupsItem.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of security group (%s) access tags: %s", d.Id(), err)
	}
	resultMap[isSecurityGroupAccessTags] = accesstags

	return resultMap
}

func dataSourceSecurityGroupCollectionSecurityGroupsResourceGroupToMap(resourceGroupItem *vpcv1.ResourceGroupReference) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resultMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resultMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resultMap["name"] = resourceGroupItem.Name
	}

	return resultMap
}

func dataSourceSecurityGroupCollectionSecurityGroupsRulesToMap(rulesItem vpcv1.SecurityGroupRuleIntf) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}
	switch reflect.TypeOf(rulesItem).String() {
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
		{
			securityGroupRule := rulesItem.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
			if securityGroupRule.ID != nil {
				resultMap["id"] = securityGroupRule.ID
			}
			if securityGroupRule.Direction != nil {
				resultMap["direction"] = securityGroupRule.Direction
			}
			if securityGroupRule.Href != nil {
				resultMap["href"] = securityGroupRule.Href
			}
			if securityGroupRule.IPVersion != nil {
				resultMap["ip_version"] = securityGroupRule.IPVersion
			}
			if securityGroupRule.Protocol != nil {
				resultMap["protocol"] = securityGroupRule.Protocol
			}
			if securityGroupRule.Remote != nil {
				remoteList := []map[string]interface{}{}
				remoteMap := dataSourceSecurityGroupsRemoteToMap(*securityGroupRule.Remote.(*vpcv1.SecurityGroupRuleRemote))
				remoteList = append(remoteList, remoteMap)
				resultMap["remote"] = remoteList
			}
			if securityGroupRule.Local != nil {
				localList := []map[string]interface{}{}
				localMap := dataSourceSecurityGroupsLocalToMap(*securityGroupRule.Local.(*vpcv1.SecurityGroupRuleLocal))
				localList = append(localList, localMap)
				resultMap["local"] = localList
			}
		}
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
		{
			securityGroupRule := rulesItem.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
			if securityGroupRule.ID != nil {
				resultMap["id"] = securityGroupRule.ID
			}
			if securityGroupRule.Direction != nil {
				resultMap["direction"] = securityGroupRule.Direction
			}
			if securityGroupRule.Href != nil {
				resultMap["href"] = securityGroupRule.Href
			}
			if securityGroupRule.IPVersion != nil {
				resultMap["ip_version"] = securityGroupRule.IPVersion
			}
			if securityGroupRule.Protocol != nil {
				resultMap["protocol"] = securityGroupRule.Protocol
			}
			if securityGroupRule.Href != nil {
				resultMap["href"] = securityGroupRule.Href
			}
			if securityGroupRule.Code != nil {
				resultMap["code"] = securityGroupRule.Code
			}
			if securityGroupRule.Type != nil {
				resultMap["type"] = securityGroupRule.Type
			}

			if securityGroupRule.Remote != nil {
				remoteList := []map[string]interface{}{}
				remoteMap := dataSourceSecurityGroupsRemoteToMap(*securityGroupRule.Remote.(*vpcv1.SecurityGroupRuleRemote))
				remoteList = append(remoteList, remoteMap)
				resultMap["remote"] = remoteList
			}
			if securityGroupRule.Local != nil {
				localList := []map[string]interface{}{}
				localMap := dataSourceSecurityGroupsLocalToMap(*securityGroupRule.Local.(*vpcv1.SecurityGroupRuleLocal))
				localList = append(localList, localMap)
				resultMap["local"] = localList
			}
		}
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
		{
			securityGroupRule := rulesItem.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
			if securityGroupRule.ID != nil {
				resultMap["id"] = securityGroupRule.ID
			}
			if securityGroupRule.Direction != nil {
				resultMap["direction"] = securityGroupRule.Direction
			}
			if securityGroupRule.Href != nil {
				resultMap["href"] = securityGroupRule.Href
			}
			if securityGroupRule.IPVersion != nil {
				resultMap["ip_version"] = securityGroupRule.IPVersion
			}
			if securityGroupRule.Protocol != nil {
				resultMap["protocol"] = securityGroupRule.Protocol
			}
			if securityGroupRule.Href != nil {
				resultMap["href"] = securityGroupRule.Href
			}
			if securityGroupRule.PortMax != nil {
				resultMap["port_max"] = securityGroupRule.PortMax
			}
			if securityGroupRule.PortMin != nil {
				resultMap["port_min"] = securityGroupRule.PortMin
			}

			if securityGroupRule.Remote != nil {
				remoteList := []map[string]interface{}{}
				remoteMap := dataSourceSecurityGroupsRemoteToMap(*securityGroupRule.Remote.(*vpcv1.SecurityGroupRuleRemote))
				remoteList = append(remoteList, remoteMap)
				resultMap["remote"] = remoteList
			}
			if securityGroupRule.Local != nil {
				localList := []map[string]interface{}{}
				localMap := dataSourceSecurityGroupsLocalToMap(*securityGroupRule.Local.(*vpcv1.SecurityGroupRuleLocal))
				localList = append(localList, localMap)
				resultMap["local"] = localList
			}
		}
	}

	return resultMap
}

func dataSourceSecurityGroupCollectionRemoteDeletedToMap(deletedItem *vpcv1.Deleted) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		resultMap["more_info"] = deletedItem.MoreInfo
	}

	return resultMap
}

func dataSourceSecurityGroupCollectionSecurityGroupsTargetsToMap(targetsItem vpcv1.SecurityGroupTargetReferenceIntf) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	// SecurityGroupTargetReference
	switch reflect.TypeOf(targetsItem).String() {
	case "*vpcv1.SecurityGroupTargetReference":
		{
			targetx := targetsItem.(*vpcv1.SecurityGroupTargetReference)
			if targetx.Deleted != nil {
				targetxDeletedList := []map[string]interface{}{}
				targetxDeletedMap := dataSourceSecurityGroupCollectionTargetsDeletedToMap(targetx.Deleted)
				targetxDeletedList = append(targetxDeletedList, targetxDeletedMap)
				resultMap["deleted"] = targetxDeletedList
			}
			if targetx.Href != nil {
				resultMap["href"] = *targetx.Href
			}
			if targetx.ID != nil {
				resultMap["id"] = *targetx.ID
			}
			if targetx.Name != nil {
				resultMap["name"] = *targetx.Name
			}
			if targetx.ResourceType != nil {
				resultMap["resource_type"] = *targetx.ResourceType
			}

		}
	case "*vpcv1.SecurityGroupTargetReferenceLoadBalancerReference":
		{
			targety := targetsItem.(*vpcv1.SecurityGroupTargetReferenceLoadBalancerReference)
			if targety.CRN != nil {
				resultMap["crn"] = *targety.CRN
			}
			if targety.Deleted != nil {
				targetyDeletedList := []map[string]interface{}{}
				targetyDeletedMap := dataSourceSecurityGroupCollectionTargetsDeleted2ToMap(targety.Deleted)
				targetyDeletedList = append(targetyDeletedList, targetyDeletedMap)
				resultMap["deleted"] = targetyDeletedList
			}
			if targety.Href != nil {
				resultMap["href"] = *targety.Href
			}
			if targety.ID != nil {
				resultMap["id"] = *targety.ID
			}
			if targety.Name != nil {
				resultMap["name"] = *targety.Name
			}
		}
	}

	return resultMap
}

func dataSourceSecurityGroupCollectionTargetsDeletedToMap(deletedItem *vpcv1.Deleted) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		resultMap["more_info"] = deletedItem.MoreInfo
	}

	return resultMap
}

func dataSourceSecurityGroupCollectionTargetsDeleted2ToMap(deletedItem *vpcv1.Deleted) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		resultMap["more_info"] = deletedItem.MoreInfo
	}

	return resultMap
}

func dataSourceSecurityGroupCollectionSecurityGroupsVPCToMap(vpcItem *vpcv1.VPCReference) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if vpcItem.CRN != nil {
		resultMap["crn"] = vpcItem.CRN
	}
	if vpcItem.Deleted != nil {
		var mapSlice []map[string]interface{}
		modelMap := dataSourceSecurityGroupCollectionVPCDeletedToMap(vpcItem.Deleted)
		mapSlice = append(mapSlice, modelMap)
		resultMap["deleted"] = mapSlice
	}
	if vpcItem.Href != nil {
		resultMap["href"] = vpcItem.Href
	}
	if vpcItem.ID != nil {
		resultMap["id"] = vpcItem.ID
	}
	if vpcItem.Name != nil {
		resultMap["name"] = vpcItem.Name
	}

	return resultMap
}

func dataSourceSecurityGroupCollectionVPCDeletedToMap(deletedItem *vpcv1.Deleted) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		resultMap["more_info"] = deletedItem.MoreInfo
	}

	return resultMap
}

func dataSourceSecurityGroupsRemoteToMap(remoteItem vpcv1.SecurityGroupRuleRemote) (remoteMap map[string]interface{}) {
	remoteMap = map[string]interface{}{}

	if remoteItem.Address != nil {
		remoteMap["address"] = *remoteItem.Address
	}

	if remoteItem.CIDRBlock != nil {
		remoteMap["cidr_block"] = *remoteItem.CIDRBlock
	}
	if remoteItem.CRN != nil {
		remoteMap["crn"] = *remoteItem.CRN
	}
	if remoteItem.Deleted != nil {
		remoteDeletedList := []map[string]interface{}{}
		remoteDeletedMap := dataSourceSecurityGroupRuleCollectionRemoteDeletedToMap(remoteItem.Deleted)
		remoteDeletedList = append(remoteDeletedList, remoteDeletedMap)
		remoteMap["deleted"] = remoteDeletedList
	}

	if remoteItem.Href != nil {
		remoteMap["href"] = *remoteItem.Href
	}
	if remoteItem.ID != nil {
		remoteMap["id"] = *remoteItem.ID
	}
	if remoteItem.Name != nil {
		remoteMap["name"] = *remoteItem.Name
	}
	return remoteMap
}

func dataSourceSecurityGroupsLocalToMap(localItem vpcv1.SecurityGroupRuleLocal) (localMap map[string]interface{}) {
	localMap = map[string]interface{}{}

	if localItem.Address != nil {
		localMap["address"] = *localItem.Address
	}

	if localItem.CIDRBlock != nil {
		localMap["cidr_block"] = *localItem.CIDRBlock
	}
	return localMap
}
