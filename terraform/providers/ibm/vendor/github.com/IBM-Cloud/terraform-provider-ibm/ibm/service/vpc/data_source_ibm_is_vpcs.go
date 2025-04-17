// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	//"encoding/json"

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

const (
	isVPCs  = "vpcs"
	isVPCID = "id"
)

func DataSourceIBMISVPCs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISVPCListRead,
		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the resource group this vpc belongs to",
			},
			"classic_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filters the collection to VPCs with the specified classic_access value",
			},
			isVPCs: {
				Type:        schema.TypeList,
				Description: "Collection of VPCs",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isVPCDefaultNetworkACL: {
							Type:     schema.TypeString,
							Computed: true,
						},

						isVPCClassicAccess: {
							Type:     schema.TypeBool,
							Computed: true,
						},

						isVPCDefaultRoutingTable: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default routing table associated with VPC",
						},

						isVPCName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC name",
						},
						isVPCID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC id",
						},

						isVPCDefaultNetworkACLName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default Network ACL name",
						},

						isVPCDefaultSecurityGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default security group name",
						},

						isVPCDefaultSecurityGroupCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default security group CRN",
						},

						isVPCDefaultNetworkACLCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default Network ACL CRN",
						},

						isVPCDefaultRoutingTableName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default routing table name",
						},

						isVPCResourceGroup: {
							Type:     schema.TypeString,
							Computed: true,
						},

						isVPCStatus: {
							Type:     schema.TypeString,
							Computed: true,
						},

						isVPCDefaultSecurityGroup: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group associated with VPC",
						},

						isVPCTags: {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      flex.ResourceIBMVPCHash,
						},

						isVPCAccessTags: {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access tags",
						},

						isVPCCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn of the resource",
						},

						flex.ResourceControllerURL: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
						},
						flex.ResourceName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the resource",
						},

						flex.ResourceCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crn of the resource",
						},

						flex.ResourceStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the resource",
						},

						flex.ResourceGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group name in which resource is provisioned",
						},

						cseSourceAddresses: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud service endpoint IP Address",
									},

									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Location info of CSE Address",
									},
								},
							},
						},

						isVPCSecurityGroupList: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isVPCSecurityGroupName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Security group name",
									},

									isVPCSecurityGroupID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Security group id",
										ForceNew:    true,
									},

									isSecurityGroupRules: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Security Rules",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{

												isVPCSecurityGroupRuleID: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule ID",
												},

												isVPCSecurityGroupRuleDirection: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Direction of traffic to enforce, either inbound or outbound",
												},

												isVPCSecurityGroupRuleIPVersion: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IP version: ipv4 or ipv6",
												},

												isVPCSecurityGroupRuleRemote: {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Security group id: an IP address, a CIDR block, or a single security group identifier",
												},

												isVPCSecurityGroupRuleType: {
													Type:     schema.TypeInt,
													Computed: true,
												},

												isVPCSecurityGroupRuleCode: {
													Type:     schema.TypeInt,
													Computed: true,
												},

												isVPCSecurityGroupRulePortMin: {
													Type:     schema.TypeInt,
													Computed: true,
												},

												isVPCSecurityGroupRulePortMax: {
													Type:     schema.TypeInt,
													Computed: true,
												},

												isVPCSecurityGroupRuleProtocol: {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},

						subnetsList: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subent name",
									},

									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet ID",
									},

									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet status",
									},

									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet location",
									},

									totalIPV4AddressCount: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total IPv4 address count in the subnet",
									},

									availableIPV4AddressCount: {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Available IPv4 address count in the subnet",
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

func dataSourceIBMISVPCListRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	start := ""
	allrecs := []vpcv1.VPC{}
	listOptions := &vpcv1.ListVpcsOptions{}
	if resgroupintf, ok := d.GetOk("resource_group"); ok {
		resGroup := resgroupintf.(string)
		listOptions.ResourceGroupID = &resGroup
	}
	if classicAccessIntf, ok := d.GetOk("classic_access"); ok {
		classicAccess := classicAccessIntf.(bool)
		listOptions.ClassicAccess = &classicAccess
	}
	for {

		if start != "" {
			listOptions.Start = &start
		}
		result, detail, err := sess.ListVpcsWithContext(context, listOptions)
		if err != nil {
			log.Printf("Error reading list of VPCs:%s\n%s", err, detail)
			return diag.FromErr(err)
		}
		start = flex.GetNext(result.Next)
		allrecs = append(allrecs, result.Vpcs...)
		if start == "" {
			break
		}
	}

	vpcs := make([]map[string]interface{}, 0)
	for _, vpc := range allrecs {

		l := map[string]interface{}{
			isVPCID:             *vpc.ID,
			isVPCName:           *vpc.Name,
			isVPCClassicAccess:  *vpc.ClassicAccess,
			isVPCStatus:         *vpc.Status,
			isVPCCRN:            *vpc.CRN,
			flex.ResourceName:   *vpc.Name,
			flex.ResourceCRN:    *vpc.CRN,
			flex.ResourceStatus: *vpc.Status,
			isVPCResourceGroup:  *vpc.ResourceGroup.ID,
		}

		if vpc.DefaultNetworkACL != nil {
			l[isVPCDefaultNetworkACL] = *vpc.DefaultNetworkACL.ID
			l[isVPCDefaultNetworkACLName] = *vpc.DefaultNetworkACL.Name
			l[isVPCDefaultNetworkACLCRN] = vpc.DefaultNetworkACL.CRN
		}
		if vpc.DefaultRoutingTable != nil {
			l[isVPCDefaultRoutingTable] = *vpc.DefaultRoutingTable.ID
			l[isVPCDefaultRoutingTableName] = *vpc.DefaultRoutingTable.Name
		}
		if vpc.DefaultSecurityGroup != nil {
			l[isVPCDefaultSecurityGroup] = *vpc.DefaultSecurityGroup.ID
			l[isVPCDefaultSecurityGroupName] = *vpc.DefaultSecurityGroup.Name
			l[isVPCDefaultSecurityGroupCRN] = vpc.DefaultSecurityGroup.CRN
		}
		tags, err := flex.GetGlobalTagsUsingCRN(meta, *vpc.CRN, "", isVPCUserTagType)
		if err != nil {
			log.Printf(
				"An error occured during reading of vpc (%s) tags : %s", d.Id(), err)
		}
		l[isVPCTags] = tags

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *vpc.CRN, "", isVPCAccessTagType)
		if err != nil {
			log.Printf(
				"An error occured during reading of vpc (%s) access tags: %s", d.Id(), err)
		}
		l[isVPCAccessTags] = accesstags

		controller, err := flex.GetBaseController(meta)
		if err != nil {
			return diag.FromErr(err)
		}
		l[flex.ResourceControllerURL] = controller + "/vpc-ext/network/vpcs"

		if vpc.ResourceGroup != nil {
			l[flex.ResourceGroupName] = *vpc.ResourceGroup.Name
		}
		//set the cse ip addresses info
		if vpc.CseSourceIps != nil {
			cseSourceIpsList := make([]map[string]interface{}, 0)
			for _, sourceIP := range vpc.CseSourceIps {
				currentCseSourceIp := map[string]interface{}{}
				if sourceIP.IP != nil {
					currentCseSourceIp["address"] = *sourceIP.IP.Address
					currentCseSourceIp["zone_name"] = *sourceIP.Zone.Name
					cseSourceIpsList = append(cseSourceIpsList, currentCseSourceIp)
				}
			}
			l[cseSourceAddresses] = cseSourceIpsList
		}

		// adding pagination support for subnets inside vpc

		startSub := ""
		allrecsSub := []vpcv1.Subnet{}
		options := &vpcv1.ListSubnetsOptions{}
		for {
			if startSub != "" {
				options.Start = &startSub
			}
			s, response, err := sess.ListSubnetsWithContext(context, options)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error fetching subnets %s\n%s", err, response))
			}
			startSub = flex.GetNext(s.Next)
			allrecsSub = append(allrecsSub, s.Subnets...)
			if startSub == "" {
				break
			}
		}
		if err == nil {
			subnetsInfo := make([]map[string]interface{}, 0)
			for _, subnet := range allrecsSub {
				if *subnet.VPC.ID == *vpc.ID {
					l := map[string]interface{}{
						"name":                    *subnet.Name,
						"id":                      *subnet.ID,
						"status":                  *subnet.Status,
						"zone":                    *subnet.Zone.Name,
						totalIPV4AddressCount:     *subnet.TotalIpv4AddressCount,
						availableIPV4AddressCount: *subnet.AvailableIpv4AddressCount,
					}
					subnetsInfo = append(subnetsInfo, l)
				}
			}
			l[subnetsList] = subnetsInfo
		}

		// adding pagination support for sg inside vpc

		startSg := ""
		allrecsSg := []vpcv1.SecurityGroup{}

		for {
			vpcId := vpc.ID
			listSgOptions := &vpcv1.ListSecurityGroupsOptions{
				VPCID: vpcId,
			}
			if startSg != "" {
				listSgOptions.Start = &start
			}
			sgs, response, err := sess.ListSecurityGroupsWithContext(context, listSgOptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error fetching Security Groups %s\n%s", err, response))
			}
			if *sgs.TotalCount == int64(0) {
				break
			}
			start = flex.GetNext(sgs.Next)
			allrecsSg = append(allrecsSg, sgs.SecurityGroups...)

			if startSg == "" {
				break
			}

		}

		securityGroupList := make([]map[string]interface{}, 0)

		for _, group := range allrecsSg {
			if *group.VPC.ID == *vpc.ID {
				g := make(map[string]interface{})

				g[isVPCSecurityGroupName] = *group.Name
				g[isVPCSecurityGroupID] = *group.ID

				rules := make([]map[string]interface{}, 0)
				for _, sgrule := range group.Rules {
					switch reflect.TypeOf(sgrule).String() {
					case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
						{
							rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
							r := make(map[string]interface{})
							if rule.Code != nil {
								r[isVPCSecurityGroupRuleCode] = int(*rule.Code)
							}
							if rule.Type != nil {
								r[isVPCSecurityGroupRuleType] = int(*rule.Type)
							}
							r[isVPCSecurityGroupRuleDirection] = *rule.Direction
							r[isVPCSecurityGroupRuleIPVersion] = *rule.IPVersion
							if rule.Protocol != nil {
								r[isVPCSecurityGroupRuleProtocol] = *rule.Protocol
							}
							r[isVPCSecurityGroupRuleID] = *rule.ID
							remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
							if ok {
								if remote != nil && reflect.ValueOf(remote).IsNil() == false {
									if remote.ID != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.ID
									} else if remote.Address != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.Address
									} else if remote.CIDRBlock != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.CIDRBlock
									}
								}
							}
							rules = append(rules, r)
						}

					case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
						{
							rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
							r := make(map[string]interface{})
							r[isVPCSecurityGroupRuleDirection] = *rule.Direction
							r[isVPCSecurityGroupRuleIPVersion] = *rule.IPVersion
							if rule.Protocol != nil {
								r[isVPCSecurityGroupRuleProtocol] = *rule.Protocol
							}
							r[isVPCSecurityGroupRuleID] = *rule.ID
							remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
							if ok {
								if remote != nil && reflect.ValueOf(remote).IsNil() == false {
									if remote.ID != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.ID
									} else if remote.Address != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.Address
									} else if remote.CIDRBlock != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.CIDRBlock
									}
								}
							}
							rules = append(rules, r)
						}

					case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
						{
							rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
							r := make(map[string]interface{})
							r[isVPCSecurityGroupRuleDirection] = *rule.Direction
							r[isVPCSecurityGroupRuleIPVersion] = *rule.IPVersion
							if rule.PortMin != nil {
								r[isVPCSecurityGroupRulePortMin] = int(*rule.PortMin)
							}
							if rule.PortMax != nil {
								r[isVPCSecurityGroupRulePortMax] = int(*rule.PortMax)
							}

							if rule.Protocol != nil {
								r[isVPCSecurityGroupRuleProtocol] = *rule.Protocol
							}

							remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
							if ok {
								if remote != nil && reflect.ValueOf(remote).IsNil() == false {
									if remote.ID != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.ID
									} else if remote.Address != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.Address
									} else if remote.CIDRBlock != nil {
										r[isVPCSecurityGroupRuleRemote] = remote.CIDRBlock
									}
								}
							}
							rules = append(rules, r)
						}
					}
				}
				g[isVPCSgRules] = rules
				securityGroupList = append(securityGroupList, g)
			}
		}
		l[isVPCSecurityGroupList] = securityGroupList
		vpcs = append(vpcs, l)
	}
	d.SetId(dataSourceIBMISVPCsID(d))
	d.Set(isVPCs, vpcs)
	return nil
}

// dataSourceIBMISVPCsID returns a reasonable ID for vpc list.
func dataSourceIBMISVPCsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
