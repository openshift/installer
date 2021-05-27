// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"reflect"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isSgName          = "name"
	isSgRules         = "rules"
	isSgRuleID        = "rule_id"
	isSgRuleDirection = "direction"
	isSgRuleIPVersion = "ip_version"
	isSgRuleRemote    = "remote"
	isSgRuleType      = "type"
	isSgRuleCode      = "code"
	isSgRulePortMax   = "port_max"
	isSgRulePortMin   = "port_min"
	isSgRuleProtocol  = "protocol"
	isSgVPC           = "vpc"
	isSgTags          = "tags"
	isSgCRN           = "crn"
)

func dataSourceIBMISSecurityGroup() *schema.Resource {
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
				Computed:    true,
				Description: "Security group's resource group id",
				ForceNew:    true,
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
							Description: "IP version: ipv4 or ipv6",
						},

						isSgRuleRemote: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group id: an IP address, a CIDR block, or a single security group identifier",
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

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},

			isSgTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         resourceIBMVPCHash,
				Description: "List of tags",
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
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	sgName := d.Get(isSgName).(string)
	if userDetails.generation == 1 {
		err := classicSecurityGroupGet(d, meta, sgName)
		if err != nil {
			return err
		}
	} else {
		err := securityGroupGet(d, meta, sgName)
		if err != nil {
			return err
		}
	}
	return nil
}

func classicSecurityGroupGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}

	listSgOptions := &vpcclassicv1.ListSecurityGroupsOptions{}
	sgs, _, err := sess.ListSecurityGroups(listSgOptions)
	if err != nil {
		return err
	}

	for _, group := range sgs.SecurityGroups {
		if *group.Name == name {

			d.Set(isSgName, *group.Name)
			d.Set(isSgVPC, *group.VPC.ID)
			d.Set(isSgCRN, *group.CRN)
			tags, err := GetTagsUsingCRN(meta, *group.CRN)
			if err != nil {
				log.Printf(
					"An error occured during reading of security group (%s) tags : %s", *group.ID, err)
			}
			d.Set(isSgTags, tags)
			rules := make([]map[string]interface{}, 0)
			for _, sgrule := range group.Rules {
				switch reflect.TypeOf(sgrule).String() {
				case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
					{

						rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)
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
						remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
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
						rules = append(rules, r)
					}

				case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
					{

						rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)
						r := make(map[string]interface{})
						r[isSgRuleDirection] = *rule.Direction
						r[isSgRuleIPVersion] = *rule.IPVersion
						if rule.Protocol != nil {
							r[isSgRuleProtocol] = *rule.Protocol
						}
						r[isSgRuleID] = *rule.ID
						remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
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
						rules = append(rules, r)
					}

				case "*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
					{
						rule := sgrule.(*vpcclassicv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)
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

						r[isSgRuleID] = *rule.ID
						remote, ok := rule.Remote.(*vpcclassicv1.SecurityGroupRuleRemote)
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
						rules = append(rules, r)
					}
				}
			}

			d.Set(isSgRules, rules)
			d.SetId(*group.ID)

			if group.ResourceGroup != nil {
				rsMangClient, err := meta.(ClientSession).ResourceManagementAPIv2()
				if err != nil {
					return err
				}
				grp, err := rsMangClient.ResourceGroup().Get(*group.ResourceGroup.ID)
				if err != nil {
					return err
				}
				d.Set(ResourceGroupName, grp.Name)
			}

			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc/network/securityGroups")
			if group.Name != nil {
				d.Set(ResourceName, *group.Name)
			}

			if group.CRN != nil {
				d.Set(ResourceCRN, *group.CRN)
			}
			return nil
		}
	}

	return nil
}

func securityGroupGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	listSgOptions := &vpcv1.ListSecurityGroupsOptions{}
	sgs, _, err := sess.ListSecurityGroups(listSgOptions)
	if err != nil {
		return err
	}

	for _, group := range sgs.SecurityGroups {
		if *group.Name == name {

			d.Set(isSgName, *group.Name)
			d.Set(isSgVPC, *group.VPC.ID)
			d.Set(isSgCRN, *group.CRN)
			tags, err := GetTagsUsingCRN(meta, *group.CRN)
			if err != nil {
				log.Printf(
					"An error occured during reading of security group (%s) tags : %s", *group.ID, err)
			}
			d.Set(isSgTags, tags)
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
						rules = append(rules, r)
					}
				}
			}

			d.Set(isSgRules, rules)
			d.SetId(*group.ID)

			if group.ResourceGroup != nil {
				rsMangClient, err := meta.(ClientSession).ResourceManagementAPIv2()
				if err != nil {
					return err
				}
				grp, err := rsMangClient.ResourceGroup().Get(*group.ResourceGroup.ID)
				if err != nil {
					return err
				}
				d.Set(ResourceGroupName, grp.Name)
			}

			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc/network/securityGroups")
			if group.Name != nil {
				d.Set(ResourceName, *group.Name)
			}

			if group.CRN != nil {
				d.Set(ResourceCRN, *group.CRN)
			}
			return nil
		}
	}

	return nil

}
