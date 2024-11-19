// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIsSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsSecurityGroupRuleRead,

		Schema: map[string]*schema.Schema{
			"security_group": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The security group identifier.",
			},
			"security_group_rule": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule identifier.",
			},
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
	}
}

func dataSourceIBMIsSecurityGroupRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getSecurityGroupRuleOptions := &vpcv1.GetSecurityGroupRuleOptions{}

	getSecurityGroupRuleOptions.SetSecurityGroupID(d.Get("security_group").(string))
	getSecurityGroupRuleOptions.SetID(d.Get("security_group_rule").(string))

	securityGroupRuleIntf, response, err := vpcClient.GetSecurityGroupRuleWithContext(context, getSecurityGroupRuleOptions)
	if err != nil || securityGroupRuleIntf == nil {
		log.Printf("[DEBUG] GetSecurityGroupRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSecurityGroupRuleWithContext failed %s\n%s", err, response))
	}

	switch reflect.TypeOf(securityGroupRuleIntf).String() {
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll":
		{
			securityGroupRule := securityGroupRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll)

			d.SetId(*securityGroupRule.ID)
			if err = d.Set("direction", securityGroupRule.Direction); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting direction: %s", err))
			}
			if err = d.Set("href", securityGroupRule.Href); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
			}
			if err = d.Set("ip_version", securityGroupRule.IPVersion); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting ip_version: %s", err))
			}
			if err = d.Set("protocol", securityGroupRule.Protocol); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting protocol: %s", err))
			}
			if securityGroupRule.Remote != nil {
				securityGroupRuleRemote, err := dataSourceSecurityGroupRuleFlattenRemote(securityGroupRule.Remote)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error flattening securityGroupRule.Remote %s", err))
				}
				err = d.Set("remote", securityGroupRuleRemote)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error setting remote %s", err))
				}
			}
			if securityGroupRule.Local != nil {
				securityGroupRuleLocal, err := dataSourceSecurityGroupRuleFlattenLocal(securityGroupRule.Local)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error flattening securityGroupRule.Local %s", err))
				}
				err = d.Set("local", securityGroupRuleLocal)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error setting local %s", err))
				}
			}

		}
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp":
		{
			securityGroupRule := securityGroupRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp)

			d.SetId(*securityGroupRule.ID)
			if err = d.Set("direction", securityGroupRule.Direction); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting direction: %s", err))
			}
			if err = d.Set("href", securityGroupRule.Href); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
			}
			if err = d.Set("ip_version", securityGroupRule.IPVersion); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting ip_version: %s", err))
			}
			if err = d.Set("protocol", securityGroupRule.Protocol); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting protocol: %s", err))
			}
			if securityGroupRule.Remote != nil {
				securityGroupRuleRemote, err := dataSourceSecurityGroupRuleFlattenRemote(securityGroupRule.Remote)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error flattening securityGroupRule.Remote %s", err))
				}
				err = d.Set("remote", securityGroupRuleRemote)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error setting remote %s", err))
				}
			}
			if securityGroupRule.Local != nil {
				securityGroupRuleLocal, err := dataSourceSecurityGroupRuleFlattenLocal(securityGroupRule.Local)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error flattening securityGroupRule.Local %s", err))
				}
				err = d.Set("local", securityGroupRuleLocal)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error setting local %s", err))
				}
			}

			if err = d.Set("code", flex.IntValue(securityGroupRule.Code)); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting code: %s", err))
			}
			if err = d.Set("type", flex.IntValue(securityGroupRule.Type)); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
			}
		}
	case "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp":
		{
			securityGroupRule := securityGroupRuleIntf.(*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp)

			d.SetId(*securityGroupRule.ID)
			if err = d.Set("direction", securityGroupRule.Direction); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting direction: %s", err))
			}
			if err = d.Set("href", securityGroupRule.Href); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
			}
			if err = d.Set("ip_version", securityGroupRule.IPVersion); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting ip_version: %s", err))
			}
			if err = d.Set("protocol", securityGroupRule.Protocol); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting protocol: %s", err))
			}
			if securityGroupRule.Remote != nil {
				securityGroupRuleRemote, err := dataSourceSecurityGroupRuleFlattenRemote(securityGroupRule.Remote)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error flattening securityGroupRule.Remote %s", err))
				}
				err = d.Set("remote", securityGroupRuleRemote)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error setting remote %s", err))
				}
			}
			if securityGroupRule.Local != nil {
				securityGroupRuleLocal, err := dataSourceSecurityGroupRuleFlattenLocal(securityGroupRule.Local)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error flattening securityGroupRule.Local %s", err))
				}
				err = d.Set("local", securityGroupRuleLocal)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error setting local %s", err))
				}
			}
			if err = d.Set("port_max", flex.IntValue(securityGroupRule.PortMax)); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting port_max: %s", err))
			}
			if err = d.Set("port_min", flex.IntValue(securityGroupRule.PortMin)); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting port_min: %s", err))
			}
		}
	}

	return nil
}

func dataSourceSecurityGroupRuleFlattenRemote(m vpcv1.SecurityGroupRuleRemoteIntf) ([]map[string]interface{}, error) {
	var ruleList []map[string]interface{}
	ruleMap := dataSourceSecurityGroupRuleRemoteToMap(m.(*vpcv1.SecurityGroupRuleRemote))
	ruleList = append(ruleList, ruleMap)
	return ruleList, nil
}

func dataSourceSecurityGroupRuleRemoteToMap(remoteItem *vpcv1.SecurityGroupRuleRemote) (remoteMap map[string]interface{}) {
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
		remoteDeletedMap := dataSourceSecurityGroupRuleRemoteDeletedToMap(remoteItem.Deleted)
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

func dataSourceSecurityGroupRuleFlattenLocal(m vpcv1.SecurityGroupRuleLocalIntf) ([]map[string]interface{}, error) {
	var ruleList []map[string]interface{}
	ruleMap := dataSourceSecurityGroupRuleLocalToMap(m.(*vpcv1.SecurityGroupRuleLocal))
	ruleList = append(ruleList, ruleMap)
	return ruleList, nil
}

func dataSourceSecurityGroupRuleLocalToMap(localItem *vpcv1.SecurityGroupRuleLocal) (localMap map[string]interface{}) {
	localMap = map[string]interface{}{}
	if localItem.Address != nil {
		localMap["address"] = *localItem.Address
	}
	if localItem.CIDRBlock != nil {
		localMap["cidr_block"] = *localItem.CIDRBlock
	}
	return localMap
}

func dataSourceSecurityGroupRuleRemoteDeletedToMap(deletedItem *vpcv1.Deleted) (resultMap map[string]interface{}) {
	resultMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		resultMap["more_info"] = deletedItem.MoreInfo
	}

	return resultMap
}
