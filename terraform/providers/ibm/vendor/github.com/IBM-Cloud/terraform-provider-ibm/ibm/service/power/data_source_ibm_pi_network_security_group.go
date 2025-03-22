// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetworkSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkSecurityGroupRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkSecurityGroupID: {
				Description:  "network security group ID.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The network security group's crn.",
				Type:        schema.TypeString,
			},
			Attr_Default: {
				Computed:    true,
				Description: "Indicates if the network security group is the default network security group in the workspace.",
				Type:        schema.TypeBool,
			},
			Attr_Members: {
				Computed:    true,
				Description: "The list of IPv4 addresses and, or network interfaces in the network security group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_ID: {
							Computed:    true,
							Description: "The ID of the member in a network security group.",
							Type:        schema.TypeString,
						},
						Attr_MacAddress: {
							Computed:    true,
							Description: "The mac address of a network interface included if the type is network-interface.",
							Type:        schema.TypeString,
						},
						Attr_NetworkInterfaceID: {
							Computed:    true,
							Description: "The network ID of a network interface included if the type is network-interface.",
							Type:        schema.TypeString,
						},
						Attr_Target: {
							Computed:    true,
							Description: "If ipv4-address type, then IPv4 address or if network-interface type, then network interface ID.",
							Type:        schema.TypeString,
						},
						Attr_Type: {
							Computed:    true,
							Description: "The type of member.",
							Type:        schema.TypeString,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_Name: {
				Computed:    true,
				Description: "The name of the network security group.",
				Type:        schema.TypeString,
			},
			Attr_Rules: {
				Computed:    true,
				Description: "The list of rules in the network security group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Action: {
							Computed:    true,
							Description: "The action to take if the rule matches network traffic.",
							Type:        schema.TypeString,
						},
						Attr_DestinationPort: {
							Computed:    true,
							Description: "The list of destination port.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Maximum: {
										Computed:    true,
										Description: "The end of the port range, if applicable. If the value is not present then the default value of 65535 will be the maximum port number.",
										Type:        schema.TypeInt,
									},
									Attr_Minimum: {
										Computed:    true,
										Description: "The start of the port range, if applicable. If the value is not present then the default value of 1 will be the minimum port number.",
										Type:        schema.TypeInt,
									},
								},
							},
							Type: schema.TypeList,
						},
						Attr_ID: {
							Computed:    true,
							Description: "The ID of the rule in a network security group.",
							Type:        schema.TypeString,
						},
						Attr_Protocol: {
							Computed:    true,
							Description: "The list of protocol.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_ICMPType: {
										Computed:    true,
										Description: "If icmp type, a ICMP packet type affected by ICMP rules and if not present then all types are matched.",
										Type:        schema.TypeString,
									},
									Attr_TCPFlags: {
										Computed:    true,
										Description: "If tcp type, the list of TCP flags and if not present then all flags are matched.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												Attr_Flag: {
													Computed:    true,
													Description: "TCP flag.",
													Type:        schema.TypeString,
												},
											},
										},
										Type: schema.TypeList,
									},
									Attr_Type: {
										Computed:    true,
										Description: "The protocol of the network traffic.",
										Type:        schema.TypeString,
									},
								},
							},
							Type: schema.TypeList,
						},
						Attr_Remote: {
							Computed:    true,
							Description: "List of remote.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_ID: {
										Computed:    true,
										Description: "The ID of the remote Network Address Group or network security group the rules apply to. Not required for default-network-address-group.",
										Type:        schema.TypeString,
									},
									Attr_Type: {
										Computed:    true,
										Description: "The type of remote group the rules apply to.",
										Type:        schema.TypeString,
									},
								},
							},
							Type: schema.TypeList,
						},
						Attr_SourcePort: {
							Computed:    true,
							Description: "List of source port",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Maximum: {
										Computed:    true,
										Description: "The end of the port range, if applicable. If the value is not present then the default value of 65535 will be the maximum port number.",
										Type:        schema.TypeInt,
									},
									Attr_Minimum: {
										Computed:    true,
										Description: "The start of the port range, if applicable. If the value is not present then the default value of 1 will be the minimum port number.",
										Type:        schema.TypeInt,
									},
								},
							},
							Type: schema.TypeList,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_UserTags: {
				Computed:    true,
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},
		},
	}
}

func dataSourceIBMPINetworkSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)

	networkSecurityGroup, err := nsgClient.Get(d.Get(Arg_NetworkSecurityGroupID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*networkSecurityGroup.ID)
	if networkSecurityGroup.Crn != nil {
		d.Set(Attr_CRN, networkSecurityGroup.Crn)
		userTags, err := flex.GetTagsUsingCRN(meta, string(*networkSecurityGroup.Crn))
		if err != nil {
			log.Printf("Error on get of pi network security group (%s) user_tags: %s", *networkSecurityGroup.ID, err)
		}
		d.Set(Attr_UserTags, userTags)
	}
	d.Set(Attr_Default, networkSecurityGroup.Default)

	if len(networkSecurityGroup.Members) > 0 {
		members := []map[string]interface{}{}
		for _, mbr := range networkSecurityGroup.Members {
			mbrMap := networkSecurityGroupMemberToMap(mbr)
			members = append(members, mbrMap)
		}
		d.Set(Attr_Members, members)
	}

	d.Set(Attr_Name, networkSecurityGroup.Name)

	if len(networkSecurityGroup.Rules) > 0 {
		rules := []map[string]interface{}{}
		for _, rule := range networkSecurityGroup.Rules {
			ruleMap := networkSecurityGroupRuleToMap(rule)
			rules = append(rules, ruleMap)
		}
		d.Set(Attr_Rules, rules)
	}

	return nil
}

func networkSecurityGroupMemberToMap(mbr *models.NetworkSecurityGroupMember) map[string]interface{} {
	mbrMap := make(map[string]interface{})
	mbrMap[Attr_ID] = mbr.ID
	if mbr.MacAddress != "" {
		mbrMap[Attr_MacAddress] = mbr.MacAddress
	}
	if mbr.NetworkInterfaceNetworkID != "" {
		mbrMap[Attr_NetworkInterfaceID] = mbr.NetworkInterfaceNetworkID
	}
	mbrMap[Attr_Target] = mbr.Target
	mbrMap[Attr_Type] = mbr.Type
	return mbrMap
}

func networkSecurityGroupRuleToMap(rule *models.NetworkSecurityGroupRule) map[string]interface{} {
	ruleMap := make(map[string]interface{})
	ruleMap[Attr_Action] = rule.Action
	if rule.DestinationPort != nil {
		destinationPortMap := networkSecurityGroupRulePortToMap(rule.DestinationPort)
		ruleMap[Attr_DestinationPort] = []map[string]interface{}{destinationPortMap}
	}

	ruleMap[Attr_ID] = rule.ID

	protocolMap := networkSecurityGroupRuleProtocolToMap(rule.Protocol)
	ruleMap[Attr_Protocol] = []map[string]interface{}{protocolMap}

	remoteMap := networkSecurityGroupRuleRemoteToMap(rule.Remote)
	ruleMap[Attr_Remote] = []map[string]interface{}{remoteMap}

	if rule.SourcePort != nil {
		sourcePortMap := networkSecurityGroupRulePortToMap(rule.SourcePort)
		ruleMap[Attr_SourcePort] = []map[string]interface{}{sourcePortMap}
	}

	return ruleMap
}

func networkSecurityGroupRulePortToMap(port *models.NetworkSecurityGroupRulePort) map[string]interface{} {
	portMap := make(map[string]interface{})
	portMap[Attr_Maximum] = port.Maximum
	portMap[Attr_Minimum] = port.Minimum
	return portMap
}

func networkSecurityGroupRuleProtocolToMap(protocol *models.NetworkSecurityGroupRuleProtocol) map[string]interface{} {
	protocolMap := make(map[string]interface{})
	if protocol.IcmpType != nil {
		protocolMap[Attr_ICMPType] = protocol.IcmpType
	}
	if len(protocol.TCPFlags) > 0 {
		tcpFlags := []map[string]interface{}{}
		for _, tcpFlagsItem := range protocol.TCPFlags {
			tcpFlagsItemMap := make(map[string]interface{})
			tcpFlagsItemMap[Attr_Flag] = tcpFlagsItem.Flag
			tcpFlags = append(tcpFlags, tcpFlagsItemMap)
		}
		protocolMap[Attr_TCPFlags] = tcpFlags
	}
	if protocol.Type != "" {
		protocolMap[Attr_Type] = protocol.Type
	}
	return protocolMap
}

func networkSecurityGroupRuleRemoteToMap(remote *models.NetworkSecurityGroupRuleRemote) map[string]interface{} {
	remoteMap := make(map[string]interface{})
	if remote.ID != "" {
		remoteMap[Attr_ID] = remote.ID
	}
	if remote.Type != "" {
		remoteMap[Attr_Type] = remote.Type
	}
	return remoteMap
}
