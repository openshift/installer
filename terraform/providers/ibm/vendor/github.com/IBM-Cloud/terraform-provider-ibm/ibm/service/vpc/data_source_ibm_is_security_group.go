// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isSgName          = "name"
	isSgRules         = "rules"
	isSgRuleID        = "rule_id"
	isSgRuleDirection = "direction"
	isSgRuleIPVersion = "ip_version"
	isSgRuleRemote    = "remote"
	isSgRuleLocal     = "local"
	isSgRuleType      = "type"
	isSgRuleCode      = "code"
	isSgRulePortMax   = "port_max"
	isSgRulePortMin   = "port_min"
	isSgRuleProtocol  = "protocol"
	isSgVPC           = "vpc"
	isSgVPCName       = "vpc_name"
	isSgTags          = "tags"
	isSgCRN           = "crn"
)

func DataSourceIBMISSecurityGroup() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceIBMISSecurityGroupRuleRead,

		Schema: map[string]*schema.Schema{

			isSgName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security group name",
			},

			isSecurityGroupVPC: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Security group's vpc id",
			},
			isSgVPCName: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{isSecurityGroupVPC},
				Description:   "Security group's vpc name",
			},
			isSecurityGroupResourceGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Security group's resource group id",
			},

			isSgRules: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Security Rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						isSgRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule id",
						},

						isSgRuleDirection: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Direction of traffic to enforce, either inbound or outbound",
						},

						isSgRuleIPVersion: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version: ipv4",
						},

						isSgRuleRemote: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group id: an IP address, a CIDR block, or a single security group identifier",
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

						isSgRuleType: {
							Type:     schema.TypeInt,
							Computed: true,
						},

						isSgRuleCode: {
							Type:     schema.TypeInt,
							Computed: true,
						},

						isSgRulePortMin: {
							Type:     schema.TypeInt,
							Computed: true,
						},

						isSgRulePortMax: {
							Type:     schema.TypeInt,
							Computed: true,
						},

						isSgRuleProtocol: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},

			isSgTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags",
			},

			isSecurityGroupAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},

			isSgCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},
		},
	}
}

func dataSourceIBMISSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {

	sgName := d.Get(isSgName).(string)
	vpcId := ""
	vpcName := ""
	rgId := ""
	if vpcIdOk, ok := d.GetOk(isSgVPC); ok {
		vpcId = vpcIdOk.(string)
	}
	if rgIdOk, ok := d.GetOk(isSecurityGroupResourceGroup); ok {
		rgId = rgIdOk.(string)
	}
	if vpcNameOk, ok := d.GetOk(isSgVPCName); ok {
		vpcName = vpcNameOk.(string)
	}
	err := securityGroupGet(d, meta, sgName, vpcId, vpcName, rgId)
	if err != nil {
		return err
	}
	return nil
}

func securityGroupGet(d *schema.ResourceData, meta interface{}, name, vpcId, vpcName, rgId string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	// Support for pagination
	start := ""
	allrecs := []vpcv1.SecurityGroup{}

	listSgOptions := &vpcv1.ListSecurityGroupsOptions{}
	if vpcId != "" {
		listSgOptions.VPCID = &vpcId
	}
	if vpcName != "" {
		listSgOptions.VPCName = &vpcName
	}
	if rgId != "" {
		listSgOptions.ResourceGroupID = &rgId
	}
	for {
		if start != "" {
			listSgOptions.Start = &start
		}
		sgs, response, err := sess.ListSecurityGroups(listSgOptions)
		if err != nil || sgs == nil {
			return fmt.Errorf("[ERROR] Error Getting Security Groups %s\n%s", err, response)
		}
		if *sgs.TotalCount == int64(0) {
			break
		}
		start = flex.GetNext(sgs.Next)
		allrecs = append(allrecs, sgs.SecurityGroups...)

		if start == "" {
			break
		}

	}

	for _, group := range allrecs {
		if *group.Name == name {

			d.Set(isSgName, *group.Name)
			d.Set(isSgVPC, *group.VPC.ID)
			d.Set(isSgVPCName, group.VPC.Name)
			d.Set(isSecurityGroupResourceGroup, group.ResourceGroup.ID)
			d.Set(isSgCRN, *group.CRN)
			tags, err := flex.GetGlobalTagsUsingCRN(meta, *group.CRN, "", isUserTagType)
			if err != nil {
				log.Printf(
					"An error occured during reading of security group (%s) tags : %s", *group.ID, err)
			}
			d.Set(isSgTags, tags)
			accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *group.CRN, "", isAccessTagType)
			if err != nil {
				log.Printf(
					"Error on get of security group (%s) access tags: %s", d.Id(), err)
			}
			d.Set(isSecurityGroupAccessTags, accesstags)
			rules := make([]map[string]interface{}, 0)
			for _, sgrule := range group.Rules {
				switch reflect.TypeOf(sgrule).String() {
				case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
					{
						rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
						r := make(map[string]interface{})
						if rule.Code != nil {
							r[isSgRuleCode] = int(*rule.Code)
						}
						if rule.Type != nil {
							r[isSgRuleType] = int(*rule.Type)
						}
						r[isSgRuleDirection] = *rule.Direction
						r[isSgRuleIPVersion] = *rule.IPVersion
						if rule.Protocol != nil {
							r[isSgRuleProtocol] = *rule.Protocol
						}
						r[isSgRuleID] = *rule.ID
						remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
						if ok {
							if remote != nil && reflect.ValueOf(remote).IsNil() == false {
								if remote.ID != nil {
									r[isSgRuleRemote] = remote.ID
								} else if remote.Address != nil {
									r[isSgRuleRemote] = remote.Address
								} else if remote.CIDRBlock != nil {
									r[isSgRuleRemote] = remote.CIDRBlock
								}
							}
						}
						local, ok := rule.Local.(*vpcv1.SecurityGroupRuleLocal)
						if ok {
							if local != nil && !reflect.ValueOf(local).IsNil() {
								localList := []map[string]interface{}{}
								localMap := dataSourceSecurityGroupRuleLocalToMap(local)
								localList = append(localList, localMap)
								r["local"] = localList
							}
						}
						rules = append(rules, r)
					}

				case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
					{
						rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
						r := make(map[string]interface{})
						r[isSgRuleDirection] = *rule.Direction
						r[isSgRuleIPVersion] = *rule.IPVersion
						if rule.Protocol != nil {
							r[isSgRuleProtocol] = *rule.Protocol
						}
						r[isSgRuleID] = *rule.ID
						remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
						if ok {
							if remote != nil && reflect.ValueOf(remote).IsNil() == false {
								if remote.ID != nil {
									r[isSgRuleRemote] = remote.ID
								} else if remote.Address != nil {
									r[isSgRuleRemote] = remote.Address
								} else if remote.CIDRBlock != nil {
									r[isSgRuleRemote] = remote.CIDRBlock
								}
							}
						}
						local, ok := rule.Local.(*vpcv1.SecurityGroupRuleLocal)
						if ok {
							if local != nil && !reflect.ValueOf(local).IsNil() {
								localList := []map[string]interface{}{}
								localMap := dataSourceSecurityGroupRuleLocalToMap(local)
								localList = append(localList, localMap)
								r["local"] = localList
							}
						}
						rules = append(rules, r)
					}

				case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
					{
						rule := sgrule.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
						r := make(map[string]interface{})
						if rule.PortMin != nil {
							r[isSgRulePortMin] = int(*rule.PortMin)
						}
						if rule.PortMax != nil {
							r[isSgRulePortMax] = int(*rule.PortMax)
						}
						r[isSgRuleDirection] = *rule.Direction
						r[isSgRuleIPVersion] = *rule.IPVersion
						if rule.Protocol != nil {
							r[isSgRuleProtocol] = *rule.Protocol
						}
						remote, ok := rule.Remote.(*vpcv1.SecurityGroupRuleRemote)
						if ok {
							if remote != nil && reflect.ValueOf(remote).IsNil() == false {
								if remote.ID != nil {
									r[isSgRuleRemote] = remote.ID
								} else if remote.Address != nil {
									r[isSgRuleRemote] = remote.Address
								} else if remote.CIDRBlock != nil {
									r[isSgRuleRemote] = remote.CIDRBlock
								}
							}
						}
						local, ok := rule.Local.(*vpcv1.SecurityGroupRuleLocal)
						if ok {
							if local != nil && !reflect.ValueOf(local).IsNil() {
								localList := []map[string]interface{}{}
								localMap := dataSourceSecurityGroupRuleLocalToMap(local)
								localList = append(localList, localMap)
								r["local"] = localList
							}
						}
						rules = append(rules, r)
					}
				}
			}

			d.Set(isSgRules, rules)
			d.SetId(*group.ID)

			if group.ResourceGroup != nil {
				if group.ResourceGroup.Name != nil {
					d.Set(flex.ResourceGroupName, *group.ResourceGroup.Name)
				}
			}

			controller, err := flex.GetBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(flex.ResourceControllerURL, controller+"/vpc/network/securityGroups")
			if group.Name != nil {
				d.Set(flex.ResourceName, *group.Name)
			}

			if group.CRN != nil {
				d.Set(flex.ResourceCRN, *group.CRN)
			}
			return nil
		}
	}
	return fmt.Errorf("[ERROR] No Security Group found with name %s", name)

}
