// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMCISFirewallRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCISFirewallRulesRead,

		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full url-encoded cloud resource name (CRN) of resource instance.",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_firewall_rules",
					"cis_id"),
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
							Description: "The firewall action to perform, \"log\" action is only available for enterprise plans instances.",
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
func DataSourceIBMCISFirewallRulesValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISFirewallRulesValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_firewall_rules",
		Schema:       validateSchema}
	return &iBMCISFirewallRulesValidator
}
func dataSourceIBMCISFirewallRulesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(conns.ClientSession).CisFirewallRulesSession()
	if err != nil {
		return diag.FromErr(err)
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))

	result, resp, err := cisClient.ListAllFirewallRules(cisClient.NewListAllFirewallRulesOptions(xAuthtoken, crn, zoneID))
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error listing the  firewall rules %s:%s", err, resp))
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
