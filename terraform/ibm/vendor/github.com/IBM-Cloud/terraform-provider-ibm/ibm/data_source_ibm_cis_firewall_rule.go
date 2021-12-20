// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIBMCISFirewallRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCISFirewallRulesRead,

		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full url-encoded cloud resource name (CRN) of resource instance.",
			},
			cisDomainID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone identifier of the zone for which firewall rules are listed.",
			},
			cisFirewallrulesList: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container for response information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisFirewallrulesID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the firewall rule.",
						},
						cisFirewallrulesPaused: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the firewall rule is active.",
						},
						cisFirewallrulesDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "To briefly describe the firewall rule, omitted from object if empty.",
						},
						cisFirewallrulesAction: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The firewall action to perform, \"log\" action is only available for enterprise plan instances.",
						},
						cisFilter: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "An existing filter.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCISFirewallRulesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(ClientSession).CisFirewallRulesSession()
	if err != nil {
		return diag.FromErr(err)
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := convertTftoCisTwoVar(d.Get(cisDomainID).(string))

	result, resp, err := cisClient.ListAllFirewallRules(cisClient.NewListAllFirewallRulesOptions(xAuthtoken, crn, zoneID))
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("Error listing the  firewall rules %s:%s", err, resp))
	}

	fwrList := make([]map[string]interface{}, 0)

	for _, instance := range result.Result {
		firewallrules := map[string]interface{}{}
		fr_filters := map[string]interface{}{}
		firewallrules[cisFirewallrulesID] = *instance.ID
		firewallrules[cisFirewallrulesPaused] = *instance.Paused
		firewallrules[cisFirewallrulesDescription] = instance.Description
		firewallrules[cisFirewallrulesAction] = *instance.Action
		fr_filters[cisFilterID] = *instance.Filter.ID
		if *instance.Filter.Paused {
			fr_filters[cisFilterPaused] = "true"
		} else {
			fr_filters[cisFilterPaused] = "false"
		}
		fr_filters[cisFilterExpression] = *instance.Filter.Expression
		fr_filters[cisFilterDescription] = instance.Filter.Description
		firewallrules[cisFilter] = fr_filters
		fwrList = append(fwrList, firewallrules)
	}
	d.SetId(dataSourceCISFirewallrulesCheckID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisFirewallrulesList, fwrList)
	return nil
}
func dataSourceCISFirewallrulesCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
