// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isNetworkACLRuleHref = "href"
)

func DataSourceIBMISNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISNetworkACLRuleRead,

		Schema: map[string]*schema.Schema{
			isNwACLID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network ACL id",
			},
			isNwACLRuleId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network ACL rule id",
			},
			isNwACLRuleBefore: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule that this rule is immediately before. If absent, this is the last rule.",
			},
			isNetworkACLRuleName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_network_acl", isNetworkACLRuleName),
				Description:  "The user-defined name for this rule",
			},
			isNetworkACLRuleProtocol: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol to enforce.",
			},
			isNetworkACLRuleHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this network ACL rule",
			},
			isNetworkACLRuleAction: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether to allow or deny matching traffic.",
			},
			isNetworkACLRuleIPVersion: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP version for this rule.",
			},
			isNetworkACLRuleSource: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source IP address or CIDR block.",
			},
			isNetworkACLRuleDestination: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The destination IP address or CIDR block.",
			},
			isNetworkACLRuleDirection: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the traffic to be matched is inbound or outbound.",
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
	}
}

func dataSourceIBMISNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	nwACLID := d.Get(isNwACLID).(string)
	name := d.Get(isNetworkACLRuleName).(string)
	err := nawaclRuleDataGet(d, meta, name, nwACLID)
	if err != nil {
		return err
	}

	return nil
}

func nawaclRuleDataGet(d *schema.ResourceData, meta interface{}, name, nwACLID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	start := ""
	allrecs := []vpcv1.NetworkACLRuleItemIntf{}
	for {
		listNetworkACLRulesOptions := &vpcv1.ListNetworkACLRulesOptions{
			NetworkACLID: &nwACLID,
		}
		if start != "" {
			listNetworkACLRulesOptions.Start = &start
		}

		ruleList, response, err := sess.ListNetworkACLRules(listNetworkACLRulesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching network acl ruless %s\n%s", err, response)
		}
		start = flex.GetNext(ruleList.Next)

		allrecs = append(allrecs, ruleList.Rules...)
		if start == "" {
			break
		}
	}

	for _, rule := range allrecs {
		switch reflect.TypeOf(rule).String() {
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp":
			{
				rulex := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp)
				if *rulex.Name == name {
					d.SetId(makeTerraformACLRuleID(nwACLID, *rulex.ID))
					d.Set(isNwACLRuleId, *rulex.ID)
					if rulex.Before != nil {
						d.Set(isNwACLRuleBefore, *rulex.Before.ID)
					}
					d.Set(isNetworkACLRuleHref, *rulex.Href)
					d.Set(isNetworkACLRuleProtocol, *rulex.Protocol)
					d.Set(isNetworkACLRuleName, *rulex.Name)
					d.Set(isNetworkACLRuleAction, *rulex.Action)
					d.Set(isNetworkACLRuleIPVersion, *rulex.IPVersion)
					d.Set(isNetworkACLRuleSource, *rulex.Source)
					d.Set(isNetworkACLRuleDestination, *rulex.Destination)
					d.Set(isNetworkACLRuleDirection, *rulex.Direction)
					d.Set(isNetworkACLRuleTCP, make([]map[string]int, 0, 0))
					d.Set(isNetworkACLRuleUDP, make([]map[string]int, 0, 0))
					icmp := make([]map[string]int, 1, 1)
					if rulex.Code != nil && rulex.Type != nil {
						icmp[0] = map[string]int{
							isNetworkACLRuleICMPCode: int(*rulex.Code),
							isNetworkACLRuleICMPType: int(*rulex.Type),
						}
					}
					d.Set(isNetworkACLRuleICMP, icmp)
					break
				}
			}
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp":
			{
				rulex := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp)
				if *rulex.Name == name {
					d.SetId(makeTerraformACLRuleID(nwACLID, *rulex.ID))
					d.Set(isNwACLRuleId, *rulex.ID)
					if rulex.Before != nil {
						d.Set(isNwACLRuleBefore, *rulex.Before.ID)
					}
					d.Set(isNetworkACLRuleHref, *rulex.Href)
					d.Set(isNetworkACLRuleProtocol, *rulex.Protocol)
					d.Set(isNetworkACLRuleName, *rulex.Name)
					d.Set(isNetworkACLRuleAction, *rulex.Action)
					d.Set(isNetworkACLRuleIPVersion, *rulex.IPVersion)
					d.Set(isNetworkACLRuleSource, *rulex.Source)
					d.Set(isNetworkACLRuleDestination, *rulex.Destination)
					d.Set(isNetworkACLRuleDirection, *rulex.Direction)
					if *rulex.Protocol == "tcp" {
						d.Set(isNetworkACLRuleICMP, make([]map[string]int, 0, 0))
						d.Set(isNetworkACLRuleUDP, make([]map[string]int, 0, 0))
						tcp := make([]map[string]int, 1, 1)
						tcp[0] = map[string]int{
							isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
							isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
						}
						tcp[0][isNetworkACLRulePortMax] = checkNetworkACLNil(rulex.DestinationPortMax)
						tcp[0][isNetworkACLRulePortMin] = checkNetworkACLNil(rulex.DestinationPortMin)
						d.Set(isNetworkACLRuleTCP, tcp)
					} else if *rulex.Protocol == "udp" {
						d.Set(isNetworkACLRuleICMP, make([]map[string]int, 0, 0))
						d.Set(isNetworkACLRuleTCP, make([]map[string]int, 0, 0))
						udp := make([]map[string]int, 1, 1)
						udp[0] = map[string]int{
							isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
							isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
						}
						udp[0][isNetworkACLRulePortMax] = checkNetworkACLNil(rulex.DestinationPortMax)
						udp[0][isNetworkACLRulePortMin] = checkNetworkACLNil(rulex.DestinationPortMin)
						d.Set(isNetworkACLRuleUDP, udp)
						break
					}
				}
			}
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll":
			{
				rulex := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll)
				if *rulex.Name == name {
					d.SetId(makeTerraformACLRuleID(nwACLID, *rulex.ID))
					d.Set(isNwACLRuleId, *rulex.ID)
					if rulex.Before != nil {
						d.Set(isNwACLRuleBefore, *rulex.Before.ID)
					}
					d.Set(isNetworkACLRuleHref, *rulex.Href)
					d.Set(isNetworkACLRuleProtocol, *rulex.Protocol)
					d.Set(isNetworkACLRuleName, *rulex.Name)
					d.Set(isNetworkACLRuleAction, *rulex.Action)
					d.Set(isNetworkACLRuleIPVersion, *rulex.IPVersion)
					d.Set(isNetworkACLRuleSource, *rulex.Source)
					d.Set(isNetworkACLRuleDestination, *rulex.Destination)
					d.Set(isNetworkACLRuleDirection, *rulex.Direction)
					d.Set(isNetworkACLRuleICMP, make([]map[string]int, 0, 0))
					d.Set(isNetworkACLRuleTCP, make([]map[string]int, 0, 0))
					d.Set(isNetworkACLRuleUDP, make([]map[string]int, 0, 0))
					break
				}
			}
		}
	}
	return nil
}
