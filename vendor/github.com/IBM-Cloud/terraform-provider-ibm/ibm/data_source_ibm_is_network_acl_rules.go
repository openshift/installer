// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"reflect"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isNwACLRules = "rules"
)

func dataSourceIBMISNetworkACLRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISNetworkACLRulesRead,

		Schema: map[string]*schema.Schema{
			isNwACLID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network ACL id",
			},
			isNwACLRules: {
				Type:        schema.TypeList,
				Description: "List of network acl rules",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isNwACLRuleId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network acl rule id.",
						},
						isNetworkACLRuleName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this rule",
						},
						isNwACLRuleBefore: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule that this rule is immediately before. If absent, this is the last rule.",
						},
						isNetworkACLRuleHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network ACL rule.",
						},
						isNetworkACLRuleProtocol: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol to enforce.",
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
				},
			},
		},
	}
}

func dataSourceIBMISNetworkACLRulesRead(d *schema.ResourceData, meta interface{}) error {
	nwACLID := d.Get(isNwACLID).(string)
	err := networkACLRulesList(d, meta, nwACLID)
	if err != nil {
		return err
	}
	return nil
}

func networkACLRulesList(d *schema.ResourceData, meta interface{}, nwACLID string) error {
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
			return fmt.Errorf("Error Fetching network acl ruless %s\n%s", err, response)
		}
		start = GetNext(ruleList.Next)

		allrecs = append(allrecs, ruleList.Rules...)
		if start == "" {
			break
		}
	}
	rulesInfo := make([]map[string]interface{}, 0)
	for _, rule := range allrecs {
		l := map[string]interface{}{}
		switch reflect.TypeOf(rule).String() {
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp":
			{
				rulex := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolIcmp)
				l[isNwACLRuleId] = *rulex.ID
				l[isNetworkACLRuleHref] = *rulex.Href
				l[isNetworkACLRuleProtocol] = *rulex.Protocol
				if rulex.Before != nil {
					l[isNwACLRuleBefore] = *rulex.Before.ID
				}
				l[isNetworkACLRuleName] = *rulex.Name
				l[isNetworkACLRuleAction] = *rulex.Action
				l[isNetworkACLRuleIPVersion] = *rulex.IPVersion
				l[isNetworkACLRuleSource] = *rulex.Source
				l[isNetworkACLRuleDestination] = *rulex.Destination
				l[isNetworkACLRuleDirection] = *rulex.Direction
				l[isNetworkACLRuleTCP] = make([]map[string]int, 0, 0)
				l[isNetworkACLRuleUDP] = make([]map[string]int, 0, 0)
				icmp := make([]map[string]int, 1, 1)
				if rulex.Code != nil && rulex.Type != nil {
					icmp[0] = map[string]int{
						isNetworkACLRuleICMPCode: int(*rulex.Code),
						isNetworkACLRuleICMPType: int(*rulex.Type),
					}
				}
				l[isNetworkACLRuleICMP] = icmp
			}
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp":
			{
				rulex := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolTcpudp)
				l[isNwACLRuleId] = *rulex.ID
				l[isNetworkACLRuleHref] = *rulex.Href
				l[isNetworkACLRuleProtocol] = *rulex.Protocol
				if rulex.Before != nil {
					l[isNwACLRuleBefore] = *rulex.Before.ID
				}
				l[isNetworkACLRuleName] = *rulex.Name
				l[isNetworkACLRuleAction] = *rulex.Action
				l[isNetworkACLRuleIPVersion] = *rulex.IPVersion
				l[isNetworkACLRuleSource] = *rulex.Source
				l[isNetworkACLRuleDestination] = *rulex.Destination
				l[isNetworkACLRuleDirection] = *rulex.Direction
				if *rulex.Protocol == "tcp" {
					l[isNetworkACLRuleICMP] = make([]map[string]int, 0, 0)
					l[isNetworkACLRuleUDP] = make([]map[string]int, 0, 0)
					tcp := make([]map[string]int, 1, 1)
					tcp[0] = map[string]int{
						isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
						isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
					}
					tcp[0][isNetworkACLRulePortMax] = checkNetworkACLNil(rulex.DestinationPortMax)
					tcp[0][isNetworkACLRulePortMin] = checkNetworkACLNil(rulex.DestinationPortMin)
					l[isNetworkACLRuleTCP] = tcp
				} else if *rulex.Protocol == "udp" {
					l[isNetworkACLRuleICMP] = make([]map[string]int, 0, 0)
					l[isNetworkACLRuleTCP] = make([]map[string]int, 0, 0)
					udp := make([]map[string]int, 1, 1)
					udp[0] = map[string]int{
						isNetworkACLRuleSourcePortMax: checkNetworkACLNil(rulex.SourcePortMax),
						isNetworkACLRuleSourcePortMin: checkNetworkACLNil(rulex.SourcePortMin),
					}
					udp[0][isNetworkACLRulePortMax] = checkNetworkACLNil(rulex.DestinationPortMax)
					udp[0][isNetworkACLRulePortMin] = checkNetworkACLNil(rulex.DestinationPortMin)
					l[isNetworkACLRuleUDP] = udp
				}
			}
		case "*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll":
			{
				rulex := rule.(*vpcv1.NetworkACLRuleItemNetworkACLRuleProtocolAll)
				l[isNwACLRuleId] = *rulex.ID
				l[isNetworkACLRuleHref] = *rulex.Href
				l[isNetworkACLRuleProtocol] = *rulex.Protocol
				if rulex.Before != nil {
					l[isNwACLRuleBefore] = *rulex.Before.ID
				}
				l[isNetworkACLRuleName] = *rulex.Name
				l[isNetworkACLRuleAction] = *rulex.Action
				l[isNetworkACLRuleIPVersion] = *rulex.IPVersion
				l[isNetworkACLRuleSource] = *rulex.Source
				l[isNetworkACLRuleDestination] = *rulex.Destination
				l[isNetworkACLRuleDirection] = *rulex.Direction
				l[isNetworkACLRuleICMP] = make([]map[string]int, 0, 0)
				l[isNetworkACLRuleTCP] = make([]map[string]int, 0, 0)
				l[isNetworkACLRuleUDP] = make([]map[string]int, 0, 0)
			}
		}
		rulesInfo = append(rulesInfo, l)
	}
	d.SetId(dataSourceIBMISNetworkACLRulesId(d))
	d.Set(isNetworkACLRules, rulesInfo)
	return nil
}

// dataSourceIBMISNetworkACLRulesId returns a reasonable ID for a rule list.
func dataSourceIBMISNetworkACLRulesId(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
