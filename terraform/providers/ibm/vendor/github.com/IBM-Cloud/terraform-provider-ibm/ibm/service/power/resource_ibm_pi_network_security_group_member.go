// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

func ResourceIBMPINetworkSecurityGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPINetworkSecurityGroupMemberCreate,
		ReadContext:   resourceIBMPINetworkSecurityGroupMemberRead,
		DeleteContext: resourceIBMPINetworkSecurityGroupMemberDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_NetworkSecurityGroupID: {
				Description: "network security group ID.",
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
			},
			Arg_NetworkSecurityGroupMemberID: {
				ConflictsWith: []string{Arg_Target, Arg_Type},
				Description:   "network security group member ID.",
				ForceNew:      true,
				Optional:      true,
				Type:          schema.TypeString,
			},
			Arg_Target: {
				ConflictsWith: []string{Arg_NetworkSecurityGroupMemberID},
				Description:   "The target member to add. An IP4 address if ipv4-address type or a network interface ID if network-interface type.",
				ForceNew:      true,
				Optional:      true,
				RequiredWith:  []string{Arg_Type},
				Type:          schema.TypeString,
			},
			Arg_Type: {
				ConflictsWith: []string{Arg_NetworkSecurityGroupMemberID},
				Description:   "The type of member.",
				ForceNew:      true,
				Optional:      true,
				RequiredWith:  []string{Arg_Target},
				Type:          schema.TypeString,
				ValidateFunc:  validate.ValidateAllowedStringValues([]string{IPV4_Address, Network_Interface}),
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
			Attr_NetworkSecurityGroupMemberID: {
				Computed:    true,
				Description: "The ID of the network security group.",
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
										Description: "IIf icmp type, a ICMP packet type affected by ICMP rules and if not present then all types are matched.",
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
				Optional: true,
				Type:     schema.TypeList,
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

func resourceIBMPINetworkSecurityGroupMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	nsgID := d.Get(Arg_NetworkSecurityGroupID).(string)
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)
	if mbrID, ok := d.GetOk(Arg_NetworkSecurityGroupMemberID); ok {
		err = nsgClient.DeleteMember(nsgID, mbrID.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPINetworkSecurityGroupMemberDeleted(ctx, nsgClient, nsgID, mbrID.(string), d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, nsgID))
	} else {
		target := d.Get(Arg_Target).(string)
		mbrType := d.Get(Arg_Type).(string)
		body := &models.NetworkSecurityGroupAddMember{
			Target: &target,
			Type:   &mbrType,
		}
		member, err := nsgClient.AddMember(nsgID, body)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%s/%s/%s", cloudInstanceID, nsgID, *member.ID))
	}

	return resourceIBMPINetworkSecurityGroupMemberRead(ctx, d, meta)
}

func resourceIBMPINetworkSecurityGroupMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, parts[0])
	networkSecurityGroup, err := nsgClient.Get(parts[1])
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), NotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if networkSecurityGroup.Crn != nil {
		d.Set(Attr_CRN, networkSecurityGroup.Crn)
		userTags, err := flex.GetTagsUsingCRN(meta, string(*networkSecurityGroup.Crn))
		if err != nil {
			log.Printf("Error on get of network security group (%s) user_tags: %s", parts[1], err)
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
		d.Set(Attr_Members, nil)
	}
	d.Set(Attr_Name, networkSecurityGroup.Name)
	d.Set(Attr_NetworkSecurityGroupMemberID, networkSecurityGroup.ID)

	if len(networkSecurityGroup.Rules) > 0 {
		rules := []map[string]interface{}{}
		for _, rule := range networkSecurityGroup.Rules {
			ruleMap := networkSecurityGroupRuleToMap(rule)
			rules = append(rules, ruleMap)
		}
		d.Set(Attr_Rules, rules)
	} else {
		d.Set(Attr_Rules, nil)
	}

	return nil
}

func resourceIBMPINetworkSecurityGroupMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if len(parts) > 2 {
		nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, parts[0])
		err = nsgClient.DeleteMember(parts[1], parts[2])
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = isWaitForIBMPINetworkSecurityGroupMemberDeleted(ctx, nsgClient, parts[1], parts[2], d.Timeout(schema.TimeoutDelete))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return nil
}
func isWaitForIBMPINetworkSecurityGroupMemberDeleted(ctx context.Context, client *instance.IBMPINetworkSecurityGroupClient, nsgID, nsgMemberID string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting},
		Target:     []string{State_NotFound},
		Refresh:    isIBMPINetworkSecurityGroupMemberDeleteRefreshFunc(client, nsgID, nsgMemberID),
		Delay:      10 * time.Second,
		MinTimeout: Timeout_Active,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPINetworkSecurityGroupMemberDeleteRefreshFunc(client *instance.IBMPINetworkSecurityGroupClient, nsgID, nsgMemberID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		nsg, err := client.Get(nsgID)
		if err != nil {
			return nsg, "", err
		}
		var mbrIDs []string
		for _, mbr := range nsg.Members {
			mbrIDs = append(mbrIDs, *mbr.ID)
		}
		if !slices.Contains(mbrIDs, nsgMemberID) {
			return nsg, State_NotFound, nil
		}
		return nsg, State_Deleting, nil
	}
}
