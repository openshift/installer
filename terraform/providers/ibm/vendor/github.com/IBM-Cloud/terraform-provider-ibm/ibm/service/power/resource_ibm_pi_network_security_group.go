// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

func ResourceIBMPINetworkSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPINetworkSecurityGroupCreate,
		ReadContext:   resourceIBMPINetworkSecurityGroupRead,
		UpdateContext: resourceIBMPINetworkSecurityGroupUpdate,
		DeleteContext: resourceIBMPINetworkSecurityGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourcePowerUserTagsCustomizeDiff(diff)
			},
		),
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_Name: {
				Description:  "The name of the network security group.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_UserTags: {
				Computed:    true,
				Description: "The user tags associated with this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
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
							Description: "ist of source port",
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
		},
	}
}

func resourceIBMPINetworkSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	name := d.Get(Arg_Name).(string)
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)

	body := &models.NetworkSecurityGroupCreate{
		Name: &name,
	}
	if v, ok := d.GetOk(Arg_UserTags); ok {
		userTags := flex.FlattenSet(v.(*schema.Set))
		body.UserTags = userTags
	}

	networkSecurityGroup, err := nsgClient.Create(body)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, ok := d.GetOk(Arg_UserTags); ok {
		if networkSecurityGroup.Crn != nil {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(*networkSecurityGroup.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi network security group (%s) pi_user_tags during creation: %s", *networkSecurityGroup.ID, err)
			}
		}
	}
	nsgID := *networkSecurityGroup.ID
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, nsgID))

	return resourceIBMPINetworkSecurityGroupRead(ctx, d, meta)
}

func resourceIBMPINetworkSecurityGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		if strings.Contains(strings.ToLower(err.Error()), NotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set(Arg_Name, networkSecurityGroup.Name)
	crn := networkSecurityGroup.Crn
	if crn != nil {
		d.Set(Attr_CRN, networkSecurityGroup.Crn)
		userTags, err := flex.GetGlobalTagsUsingCRN(meta, string(*networkSecurityGroup.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on get of network security group (%s) pi_user_tags: %s", nsgID, err)
		}
		d.Set(Arg_UserTags, userTags)
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

func resourceIBMPINetworkSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, nsgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi network security group (%s) pi_user_tags: %s", nsgID, err)
			}
		}
	}
	if d.HasChange(Arg_Name) {
		body := &models.NetworkSecurityGroupUpdate{
			Name: d.Get(Arg_Name).(string),
		}
		nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)
		_, err = nsgClient.Update(nsgID, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceIBMPINetworkSecurityGroupRead(ctx, d, meta)
}

func resourceIBMPINetworkSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID, nsgID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	nsgClient := instance.NewIBMIPINetworkSecurityGroupClient(ctx, sess, cloudInstanceID)
	err = nsgClient.Delete(nsgID)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = isWaitForIBMPINetworkSecurityGroupDeleted(ctx, nsgClient, nsgID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}
func isWaitForIBMPINetworkSecurityGroupDeleted(ctx context.Context, client *instance.IBMPINetworkSecurityGroupClient, nsgID string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Deleting},
		Target:     []string{State_NotFound},
		Refresh:    isIBMPINetworkSecurityGroupDeleteRefreshFunc(client, nsgID),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Timeout:    timeout,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPINetworkSecurityGroupDeleteRefreshFunc(client *instance.IBMPINetworkSecurityGroupClient, nsgID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		nsg, err := client.Get(nsgID)
		if err != nil {
			return nsg, State_NotFound, nil
		}
		return nsg, State_Deleting, nil
	}
}
