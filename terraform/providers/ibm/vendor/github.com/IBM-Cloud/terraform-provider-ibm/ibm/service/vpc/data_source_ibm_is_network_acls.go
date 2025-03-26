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

func DataSourceIBMIsNetworkAcls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsNetworkAclsRead,

		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource group identifier.",
			},
			"network_acls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of network ACLs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the network ACL was created.",
						},
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this network ACL.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network ACL.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this network ACL.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this network ACL.",
						},
						"resource_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this network ACL.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this resource group.",
									},
								},
							},
						},
						isNetworkACLRules: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The ordered rules for this network ACL. If no rules exist, all traffic will be denied.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isNetworkACLRuleAction: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether to allow or deny matching traffic.",
									},
									"before": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The rule that this rule is immediately before. In a rule collection, this alwaysrefers to the next item in the collection. If absent, this is the last rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
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
													Description: "The URL for this network ACL rule.",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this network ACL rule.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The user-defined name for this network ACL rule.",
												},
											},
										},
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the rule was created.",
									},
									isNetworkACLRuleDestination: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The destination CIDR block. The CIDR block `0.0.0.0/0` applies to all addresses.",
									},
									isNetworkACLRuleDirection: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether the traffic to be matched is `inbound` or `outbound`.",
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this network ACL rule.",
									},
									isNetworkACLRuleID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this network ACL rule.",
									},
									isNetworkACLRuleIPVersion: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP version for this rule.",
									},
									isNetworkACLRuleName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this rule. Names must be unique within the network ACL the rule resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									isNetworkACLRuleProtocol: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The protocol to enforce.",
									},
									isNetworkACLRuleSource: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The source CIDR block. The CIDR block `0.0.0.0/0` applies to all addresses.",
									},
									isNetworkACLRuleICMP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The protocol ICMP",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isNetworkACLRuleICMPCode: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ICMP traffic code to allow. Valid values from 0 to 255.",
												},
												isNetworkACLRuleICMPType: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ICMP traffic type to allow. Valid values from 0 to 254.",
												},
											},
										},
									},

									isNetworkACLRuleTCP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "TCP protocol",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isNetworkACLRulePortMax: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The highest port in the range of ports to be matched",
												},
												isNetworkACLRulePortMin: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The lowest port in the range of ports to be matched",
												},
												isNetworkACLRuleSourcePortMax: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The highest port in the range of ports to be matched",
												},
												isNetworkACLRuleSourcePortMin: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The lowest port in the range of ports to be matched",
												},
											},
										},
									},

									isNetworkACLRuleUDP: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "UDP protocol",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												isNetworkACLRulePortMax: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The highest port in the range of ports to be matched",
												},
												isNetworkACLRulePortMin: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The lowest port in the range of ports to be matched",
												},
												isNetworkACLRuleSourcePortMax: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The highest port in the range of ports to be matched",
												},
												isNetworkACLRuleSourcePortMin: {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The lowest port in the range of ports to be matched",
												},
											},
										},
									},
								},
							},
						},
						"subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subnets to which this network ACL is attached.",
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
						"vpc": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC this network ACL is a part of.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
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
										Description: "The URL for this VPC.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this VPC.",
									},
								},
							},
						},
						isNetworkACLAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access tags",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsNetworkAclsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	resource_group_id := d.Get("resource_group").(string)
	start := ""
	allrecs := []vpcv1.NetworkACL{}
	listNetworkAclsOptions := &vpcv1.ListNetworkAclsOptions{}
	if resource_group_id != "" {
		listNetworkAclsOptions.ResourceGroupID = &resource_group_id
	}
	for {
		if start != "" {
			listNetworkAclsOptions.Start = &start
		}
		networkACLCollection, response, err := vpcClient.ListNetworkAclsWithContext(context, listNetworkAclsOptions)
		if err != nil || networkACLCollection == nil {
			log.Printf("[DEBUG] ListNetworkAclsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListNetworkAclsWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(networkACLCollection.Next)
		allrecs = append(allrecs, networkACLCollection.NetworkAcls...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsNetworkAclsID(d))

	err = d.Set("network_acls", dataSourceNetworkACLCollectionFlattenNetworkAcls(allrecs, d, meta))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting network_acls %s", err))
	}

	return nil
}

// dataSourceIBMIsNetworkAclsID returns a reasonable ID for the list.
func dataSourceIBMIsNetworkAclsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceNetworkACLCollectionFlattenNetworkAcls(result []vpcv1.NetworkACL, d *schema.ResourceData, meta interface{}) (networkAcls []map[string]interface{}) {
	for _, networkAclsItem := range result {
		networkAcls = append(networkAcls, dataSourceNetworkACLCollectionNetworkAclsToMap(networkAclsItem, d, meta))
	}

	return networkAcls
}

func dataSourceNetworkACLCollectionNetworkAclsToMap(networkAclsItem vpcv1.NetworkACL, d *schema.ResourceData, meta interface{}) (networkAclsMap map[string]interface{}) {
	networkAclsMap = map[string]interface{}{}

	if networkAclsItem.CreatedAt != nil {
		networkAclsMap["created_at"] = networkAclsItem.CreatedAt.String()
	}
	if networkAclsItem.CRN != nil {
		networkAclsMap["crn"] = networkAclsItem.CRN
	}
	if networkAclsItem.Href != nil {
		networkAclsMap["href"] = networkAclsItem.Href
	}
	if networkAclsItem.ID != nil {
		networkAclsMap["id"] = networkAclsItem.ID
	}
	if networkAclsItem.Name != nil {
		networkAclsMap["name"] = networkAclsItem.Name
	}
	if networkAclsItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceNetworkACLCollectionNetworkAclsResourceGroupToMap(*networkAclsItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		networkAclsMap["resource_group"] = resourceGroupList
	}
	if networkAclsItem.Rules != nil {
		rulesList := []map[string]interface{}{}
		for _, rulesItem := range networkAclsItem.Rules {
			rulesList = append(rulesList, dataSourceNetworkACLRulesToMap(rulesItem))
		}
		networkAclsMap[isNetworkACLRules] = rulesList
	}
	if networkAclsItem.Subnets != nil {
		subnetsList := []map[string]interface{}{}
		for _, subnetsItem := range networkAclsItem.Subnets {
			subnetsList = append(subnetsList, dataSourceNetworkACLCollectionNetworkAclsSubnetsToMap(subnetsItem))
		}
		networkAclsMap["subnets"] = subnetsList
	}
	if networkAclsItem.VPC != nil {
		vpcList := []map[string]interface{}{}
		vpcMap := dataSourceNetworkACLCollectionNetworkAclsVPCToMap(*networkAclsItem.VPC)
		vpcList = append(vpcList, vpcMap)
		networkAclsMap["vpc"] = vpcList
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *networkAclsItem.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource Network ACL (%s) access tags: %s", d.Id(), err)
	}
	networkAclsMap[isNetworkACLAccessTags] = accesstags

	return networkAclsMap
}

func dataSourceNetworkACLCollectionNetworkAclsResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}

func dataSourceNetworkACLCollectionRulesBeforeToMap(beforeItem vpcv1.NetworkACLRuleReference) (beforeMap map[string]interface{}) {
	beforeMap = map[string]interface{}{}

	if beforeItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkACLCollectionBeforeDeletedToMap(*beforeItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		beforeMap["deleted"] = deletedList
	}
	if beforeItem.Href != nil {
		beforeMap["href"] = beforeItem.Href
	}
	if beforeItem.ID != nil {
		beforeMap["id"] = beforeItem.ID
	}
	if beforeItem.Name != nil {
		beforeMap["name"] = beforeItem.Name
	}

	return beforeMap
}

func dataSourceNetworkACLCollectionBeforeDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceNetworkACLCollectionNetworkAclsSubnetsToMap(subnetsItem vpcv1.SubnetReference) (subnetsMap map[string]interface{}) {
	subnetsMap = map[string]interface{}{}

	if subnetsItem.CRN != nil {
		subnetsMap["crn"] = subnetsItem.CRN
	}
	if subnetsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkACLCollectionSubnetsDeletedToMap(*subnetsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		subnetsMap["deleted"] = deletedList
	}
	if subnetsItem.Href != nil {
		subnetsMap["href"] = subnetsItem.Href
	}
	if subnetsItem.ID != nil {
		subnetsMap["id"] = subnetsItem.ID
	}
	if subnetsItem.Name != nil {
		subnetsMap["name"] = subnetsItem.Name
	}

	return subnetsMap
}

func dataSourceNetworkACLCollectionSubnetsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceNetworkACLCollectionNetworkAclsVPCToMap(vpcItem vpcv1.VPCReference) (vpcMap map[string]interface{}) {
	vpcMap = map[string]interface{}{}

	if vpcItem.CRN != nil {
		vpcMap["crn"] = vpcItem.CRN
	}
	if vpcItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceNetworkACLCollectionVPCDeletedToMap(*vpcItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		vpcMap["deleted"] = deletedList
	}
	if vpcItem.Href != nil {
		vpcMap["href"] = vpcItem.Href
	}
	if vpcItem.ID != nil {
		vpcMap["id"] = vpcItem.ID
	}
	if vpcItem.Name != nil {
		vpcMap["name"] = vpcItem.Name
	}

	return vpcMap
}

func dataSourceNetworkACLCollectionVPCDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
