// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPINetworkSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPINetworkSecurityGroupRuleCreate,
		ReadContext:   resourceIBMPINetworkSecurityGroupRuleRead,
		DeleteContext: resourceIBMPINetworkSecurityGroupRuleDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_Action: {
				ConflictsWith: []string{Arg_NetworkSecurityGroupRuleID},
				Description:   "The action to take if the rule matches network traffic.",
				ForceNew:      true,
				Optional:      true,
				Type:          schema.TypeString,
				ValidateFunc:  validate.ValidateAllowedStringValues([]string{Allow, Deny}),
			},
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_DestinationPort: {
				ConflictsWith: []string{Arg_DestinationPorts, Arg_NetworkSecurityGroupRuleID},
				Description:   "Destination port ranges.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Maximum: {
							Default:     65535,
							Description: "The end of the port range, if applicable. If the value is not present then the default value of 65535 will be the maximum port number.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						Attr_Minimum: {
							Default:     1,
							Description: "The start of the port range, if applicable. If the value is not present then the default value of 1 will be the minimum port number.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
					},
				},
				ForceNew: true,
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			Arg_DestinationPorts: {
				ConflictsWith: []string{Arg_DestinationPort, Arg_NetworkSecurityGroupRuleID},
				Deprecated:    "This field is deprecated. Please use 'pi_destination_port' instead.",
				Description:   "Destination port ranges.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Maximum: {
							Default:     65535,
							Description: "The end of the port range, if applicable. If the value is not present then the default value of 65535 will be the maximum port number.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						Attr_Minimum: {
							Default:     1,
							Description: "The start of the port range, if applicable. If the value is not present then the default value of 1 will be the minimum port number.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
					},
				},
				ForceNew: true,
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			Arg_NetworkSecurityGroupID: {
				Description: "The unique identifier of the network security group.",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_NetworkSecurityGroupRuleID: {
				ConflictsWith: []string{Arg_Action, Arg_DestinationPort, Arg_DestinationPorts, Arg_Protocol, Arg_Remote, Arg_SourcePort, Arg_SourcePorts},
				Description:   "The network security group rule id to remove.",
				ForceNew:      true,
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_Protocol: {
				ConflictsWith: []string{Arg_NetworkSecurityGroupRuleID},
				Description:   "The protocol of the network traffic.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_ICMPType: {
							Description:  "If icmp type, a ICMP packet type affected by ICMP rules and if not present then all types are matched.",
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{All, DestinationUnreach, Echo, EchoReply, SourceQuench, TimeExceeded}),
						},
						Attr_TCPFlags: {
							Description: "If tcp type, the list of TCP flags and if not present then all flags are matched.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Flag: {
										Description: "TCP flag.",
										Required:    true,
										Type:        schema.TypeString,
									},
								},
							},
							Optional: true,
							Type:     schema.TypeSet,
						},
						Attr_Type: {
							Description:  "The protocol of the network traffic.",
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{All, ICMP, TCP, UDP}),
						},
					},
				},
				ForceNew: true,
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			Arg_Remote: {
				ConflictsWith: []string{Arg_NetworkSecurityGroupRuleID},
				Description:   "The protocol of the network traffic.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_ID: {
							Description: "The ID of the remote network address group or network security group the rules apply to. Not required for default-network-address-group.",
							Optional:    true,
							Type:        schema.TypeString,
						},
						Attr_Type: {
							Description:  "The type of remote group (MAC addresses, IP addresses, CIDRs, external CIDRs) that are the originators of rule's network traffic to match.",
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{DefaultNAG, NAG, NSG}),
						},
					},
				},
				ForceNew: true,
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			Arg_SourcePort: {
				ConflictsWith: []string{Arg_NetworkSecurityGroupRuleID, Arg_SourcePorts},
				Description:   "Source port ranges.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Maximum: {
							Default:     65535,
							Description: "The end of the port range, if applicable. If the value is not present then the default value of 65535 will be the maximum port number.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						Attr_Minimum: {
							Default:     1,
							Description: "The start of the port range, if applicable. If the value is not present then the default value of 1 will be the minimum port number.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
					},
				},
				ForceNew: true,
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
			},
			Arg_SourcePorts: {
				ConflictsWith: []string{Arg_NetworkSecurityGroupRuleID, Arg_SourcePort},
				Deprecated:    "This field is deprecated. 'Please use pi_source_port' instead.",
				Description:   "Source port ranges.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_Maximum: {
							Default:     65535,
							Description: "The end of the port range, if applicable. If the value is not present then the default value of 65535 will be the maximum port number.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
						Attr_Minimum: {
							Default:     1,
							Description: "The start of the port range, if applicable. If the value is not present then the default value of 1 will be the minimum port number.",
							Optional:    true,
							Type:        schema.TypeInt,
						},
					},
				},
				ForceNew: true,
				MaxItems: 1,
				Optional: true,
				Type:     schema.TypeList,
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
			Attr_NetworkSecurityGroupID: {
				Computed:    true,
				Description: "The unique identifier of the network security group.",
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
							Description: "Destination port ranges.",
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
										Description: "The ID of the remote network address group or network security group the rules apply to. Not required for default-network-address-group.",
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
							Description: "Source port ranges.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_Maximum: {
										Computed:    true,
										Description: "The end of the port range, if applicable. If the value is not present then the default value of 65535 will be the maximum port number.",
										Type:        schema.TypeFloat,
									},
									Attr_Minimum: {
										Computed:    true,
										Description: "The start of the port range, if applicable. If the value is not present then the default value of 1 will be the minimum port number.",
										Type:        schema.TypeFloat,
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

func resourceIBMPINetworkSecurityGroupRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)
	nsgID := d.Get(Arg_NetworkSecurityGroupID).(string)

	if v, ok := d.GetOk(Arg_NetworkSecurityGroupRuleID); ok {
		ruleID := v.(string)
		err := nsgClient.DeleteRule(nsgID, ruleID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), NotFound) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPINetworkSecurityGroupRuleRemove(ctx, nsgClient, nsgID, ruleID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, nsgID))
	} else {
		action := d.Get(Arg_Action).(string)

		networkSecurityGroupAddRule := models.NetworkSecurityGroupAddRule{
			Action: &action,
		}

		// Add protocol
		protocol := d.Get(Arg_Protocol + ".0").(map[string]interface{})
		networkSecurityGroupAddRule.Protocol = networkSecurityGroupRuleMapToProtocol(protocol)

		// Add remote
		remote := d.Get(Arg_Remote + ".0").(map[string]interface{})
		networkSecurityGroupAddRule.Remote = networkSecurityGroupRuleMapToRemote(remote)

		// Check against incompatible options
		if networkSecurityGroupAddRule.Protocol.Type == All || networkSecurityGroupAddRule.Protocol.Type == ICMP {
			_, okDestPorts := d.GetOk(Arg_DestinationPorts)
			_, okDestPort := d.GetOk(Arg_DestinationPort)
			_, okSourcePorts := d.GetOk(Arg_SourcePorts)
			_, okSourcePort := d.GetOk(Arg_SourcePort)

			if okDestPorts || okDestPort || okSourcePorts || okSourcePort {
				return diag.Errorf("pi_destination_ports, pi_destination_port, pi_source_ports, and pi_source_port are not allowed with protocol value of %s or %s", All, ICMP)
			}
		}

		// Optional fields
		if _, ok := d.GetOk(Arg_DestinationPorts); ok {
			destinationPort := d.Get(Arg_DestinationPorts + ".0").(map[string]interface{})
			networkSecurityGroupAddRule.DestinationPorts = networkSecurityGroupRuleMapToPort(destinationPort)
		} else if _, ok := d.GetOk(Arg_DestinationPort); ok {
			destinationPort := d.Get(Arg_DestinationPort + ".0").(map[string]interface{})
			networkSecurityGroupAddRule.DestinationPort = networkSecurityGroupRuleMapToPort(destinationPort)
		}

		if _, ok := d.GetOk(Arg_SourcePorts); ok {
			sourcePort := d.Get(Arg_SourcePorts + ".0").(map[string]interface{})
			networkSecurityGroupAddRule.SourcePorts = networkSecurityGroupRuleMapToPort(sourcePort)
		} else if _, ok := d.GetOk(Arg_SourcePort); ok {
			sourcePort := d.Get(Arg_SourcePort + ".0").(map[string]interface{})
			networkSecurityGroupAddRule.SourcePort = networkSecurityGroupRuleMapToPort(sourcePort)
		}

		networkSecurityGroup, err := nsgClient.AddRule(nsgID, &networkSecurityGroupAddRule)
		if err != nil {
			return diag.FromErr(err)
		}
		ruleID := *networkSecurityGroup.ID

		_, err = isWaitForIBMPINetworkSecurityGroupRuleAdd(ctx, nsgClient, nsgID, ruleID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, nsgID, ruleID))
	}

	return resourceIBMPINetworkSecurityGroupRuleRead(ctx, d, meta)
}

func resourceIBMPINetworkSecurityGroupRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, nsgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)
	networkSecurityGroup, err := nsgClient.Get(nsgID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(Attr_Name, networkSecurityGroup.Name)

	if networkSecurityGroup.Crn != nil {
		d.Set(Attr_CRN, networkSecurityGroup.Crn)
		userTags, err := flex.GetTagsUsingCRN(meta, string(*networkSecurityGroup.Crn))
		if err != nil {
			log.Printf("Error on get of network security group (%s) user_tags: %s", nsgID, err)
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
	} else {
		d.Set(Attr_Members, []string{})
	}

	d.Set(Attr_NetworkSecurityGroupID, networkSecurityGroup.ID)
	if len(networkSecurityGroup.Rules) > 0 {
		rules := []map[string]interface{}{}
		for _, rule := range networkSecurityGroup.Rules {
			ruleMap := networkSecurityGroupRuleToMap(rule)
			rules = append(rules, ruleMap)
		}
		d.Set(Attr_Rules, rules)
	} else {
		d.Set(Attr_Rules, []string{})
	}
	return nil
}

func resourceIBMPINetworkSecurityGroupRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ids, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(ids) == 3 {
		cloudInstanceID := ids[0]
		nsgID := ids[1]
		ruleID := ids[2]

		sess, err := meta.(conns.ClientSession).IBMPISession()
		if err != nil {
			return diag.FromErr(err)
		}
		nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)

		err = nsgClient.DeleteRule(nsgID, ruleID)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = isWaitForIBMPINetworkSecurityGroupRuleRemove(ctx, nsgClient, nsgID, ruleID, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return nil
}

func isWaitForIBMPINetworkSecurityGroupRuleAdd(ctx context.Context, client *instance.IBMPINetworkSecurityGroupClient, id, ruleID string, timeout time.Duration) (interface{}, error) {

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Pending},
		Target:     []string{State_Available},
		Refresh:    isIBMPINetworkSecurityGroupRuleAddRefreshFunc(client, id, ruleID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPINetworkSecurityGroupRuleAddRefreshFunc(client *instance.IBMPINetworkSecurityGroupClient, id, ruleID string) retry.StateRefreshFunc {

	return func() (interface{}, string, error) {
		networkSecurityGroup, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if networkSecurityGroup.Rules != nil {
			for _, rule := range networkSecurityGroup.Rules {
				if *rule.ID == ruleID {
					return networkSecurityGroup, State_Available, nil
				}

			}
		}
		return networkSecurityGroup, State_Pending, nil
	}
}

func isWaitForIBMPINetworkSecurityGroupRuleRemove(ctx context.Context, client *instance.IBMPINetworkSecurityGroupClient, id, ruleID string, timeout time.Duration) (interface{}, error) {

	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Pending},
		Target:     []string{State_Removed},
		Refresh:    isIBMPINetworkSecurityGroupRuleRemoveRefreshFunc(client, id, ruleID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: time.Minute,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPINetworkSecurityGroupRuleRemoveRefreshFunc(client *instance.IBMPINetworkSecurityGroupClient, id, ruleID string) retry.StateRefreshFunc {

	return func() (interface{}, string, error) {
		networkSecurityGroup, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if networkSecurityGroup.Rules != nil {
			foundRule := false
			for _, rule := range networkSecurityGroup.Rules {
				if *rule.ID == ruleID {
					foundRule = true
					return networkSecurityGroup, State_Pending, nil
				}
			}
			if !foundRule {
				return networkSecurityGroup, State_Removed, nil
			}
		}
		return networkSecurityGroup, State_Pending, nil
	}
}

func networkSecurityGroupRuleMapToPort(portMap map[string]interface{}) *models.NetworkSecurityGroupRulePort {
	networkSecurityGroupRulePort := models.NetworkSecurityGroupRulePort{}
	if portMap[Attr_Maximum] != nil {
		networkSecurityGroupRulePort.Maximum = int64(portMap[Attr_Maximum].(int))
	}
	if portMap[Attr_Minimum] != nil {
		networkSecurityGroupRulePort.Minimum = int64(portMap[Attr_Minimum].(int))
	}
	return &networkSecurityGroupRulePort
}

func networkSecurityGroupRuleMapToRemote(remoteMap map[string]interface{}) *models.NetworkSecurityGroupRuleRemote {
	networkSecurityGroupRuleRemote := models.NetworkSecurityGroupRuleRemote{}
	if remoteMap[Attr_ID].(string) != "" {
		networkSecurityGroupRuleRemote.ID = remoteMap[Attr_ID].(string)
	}
	networkSecurityGroupRuleRemote.Type = remoteMap[Attr_Type].(string)
	return &networkSecurityGroupRuleRemote
}

func networkSecurityGroupRuleMapToProtocol(protocolMap map[string]interface{}) *models.NetworkSecurityGroupRuleProtocol {
	networkSecurityGroupRuleProtocol := models.NetworkSecurityGroupRuleProtocol{}
	networkSecurityGroupRuleProtocol.Type = protocolMap[Attr_Type].(string)

	if networkSecurityGroupRuleProtocol.Type == ICMP {
		icmpType := protocolMap[Attr_ICMPType].(string)
		networkSecurityGroupRuleProtocol.IcmpType = &icmpType
	} else if networkSecurityGroupRuleProtocol.Type == TCP {
		tcpMaps := protocolMap[Attr_TCPFlags].(*schema.Set)
		networkSecurityGroupRuleProtocolTCPFlagArray := []*models.NetworkSecurityGroupRuleProtocolTCPFlag{}
		for _, tcpMap := range tcpMaps.List() {
			flag := tcpMap.(map[string]interface{})
			networkSecurityGroupRuleProtocolTCPFlag := models.NetworkSecurityGroupRuleProtocolTCPFlag{}
			networkSecurityGroupRuleProtocolTCPFlag.Flag = flag[Attr_Flag].(string)
			networkSecurityGroupRuleProtocolTCPFlagArray = append(networkSecurityGroupRuleProtocolTCPFlagArray, &networkSecurityGroupRuleProtocolTCPFlag)
		}
		networkSecurityGroupRuleProtocol.TCPFlags = networkSecurityGroupRuleProtocolTCPFlagArray
	}

	return &networkSecurityGroupRuleProtocol
}
