// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPINetworkSecurityGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPINetworkSecurityGroupsRead,

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			// Attributes
			Attr_NetworkSecurityGroups: {
				Computed:    true,
				Description: "list of Network Security Groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						Attr_ID: {
							Computed:    true,
							Description: "The ID of the network security group.",
							Type:        schema.TypeString,
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
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPINetworkSecurityGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)
	nsgResp, err := nsgClient.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}

	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)

	networkSecurityGroups := []map[string]interface{}{}
	if len(nsgResp.NetworkSecurityGroups) > 0 {
		for _, nsg := range nsgResp.NetworkSecurityGroups {
			networkSecurityGroup := networkSecurityGroupToMap(nsg, meta)
			networkSecurityGroups = append(networkSecurityGroups, networkSecurityGroup)
		}
	}

	d.Set(Attr_NetworkSecurityGroups, networkSecurityGroups)

	return nil
}

func networkSecurityGroupToMap(nsg *models.NetworkSecurityGroup, meta interface{}) map[string]interface{} {
	networkSecurityGroup := make(map[string]interface{})
	if nsg.Crn != nil {
		networkSecurityGroup[Attr_CRN] = nsg.Crn
		userTags, err := flex.GetTagsUsingCRN(meta, string(*nsg.Crn))
		if err != nil {
			log.Printf("Error on get of pi network security group (%s) user_tags: %s", *nsg.ID, err)
		}
		networkSecurityGroup[Attr_UserTags] = userTags
	}
	networkSecurityGroup[Attr_Default] = nsg.Default

	networkSecurityGroup[Attr_ID] = nsg.ID
	if len(nsg.Members) > 0 {
		members := []map[string]interface{}{}
		for _, mbr := range nsg.Members {
			mbrMap := networkSecurityGroupMemberToMap(mbr)
			members = append(members, mbrMap)
		}
		networkSecurityGroup[Attr_Members] = members
	}
	networkSecurityGroup[Attr_Name] = nsg.Name
	if len(nsg.Rules) > 0 {
		rules := []map[string]interface{}{}
		for _, rule := range nsg.Rules {
			rulesItemMap := networkSecurityGroupRuleToMap(rule)
			rules = append(rules, rulesItemMap)
		}
		networkSecurityGroup[Attr_Rules] = rules
	}
	return networkSecurityGroup
}
