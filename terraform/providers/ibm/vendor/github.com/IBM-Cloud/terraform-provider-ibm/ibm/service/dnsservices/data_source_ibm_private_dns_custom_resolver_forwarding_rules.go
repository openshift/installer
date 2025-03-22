// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMPrivateDNSForwardingRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmDnsCrForwardingRulesRead,

		Schema: map[string]*schema.Schema{
			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a service instance.",
			},
			pdnsCRFRResolverID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a custom resolver.",
			},
			pdnsCRForwardRules: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsCRFRRuleID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the forwarding rule.",
						},
						pdnsCRFRDesctiption: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Descriptive text of the forwarding rule.",
						},
						pdnsCRFRType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the forwarding rule.",
						},
						pdnsCRFRMatch: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The matching zone or hostname.",
						},
						pdnsCRFRForwardTo: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The upstream DNS servers will be forwarded to.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						pdnsCRFRViews: {
							Type:        schema.TypeList,
							Description: "An array of views used by forwarding rules.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									pdnsCRFRVName: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unique name of the view.",
									},
									pdnsCRFRVDescription: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the view.",
									},
									pdnsCRFRVExpression: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expression of the view.",
									},
									pdnsCRFRVForwardTo: {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The upstream DNS servers that the matching DNS queries will be forwarded to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
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

func dataSourceIbmDnsCrForwardingRulesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	resolverID := d.Get(pdnsCRFRResolverID).(string)

	opt := sess.NewListForwardingRulesOptions(instanceID, resolverID)

	result, resp, err := sess.ListForwardingRulesWithContext(context, opt)
	if err != nil || result == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error listing the forwarding rules %s:%s", err, resp))
	}

	forwardRules := make([]interface{}, 0)
	for _, rule := range result.ForwardingRules {
		forwardRule := map[string]interface{}{}
		forwardRule[pdnsCRFRRuleID] = *rule.ID
		forwardRule[pdnsCRFRDesctiption] = *rule.Description
		forwardRule[pdnsCRFRType] = *rule.Type
		forwardRule[pdnsCRFRMatch] = *rule.Match
		forwardRule[pdnsCRFRForwardTo] = rule.ForwardTo
		forwardRule[pdnsCRFRViews] = flattenPDNSFRViews(rule.Views)

		forwardRules = append(forwardRules, forwardRule)
	}
	d.SetId(dataSourceIBMPrivateDNSForwardrulesID())
	d.Set(pdnsInstanceID, instanceID)
	d.Set(pdnsCRFRResolverID, resolverID)
	d.Set(pdnsCRForwardRules, forwardRules)
	return nil
}

func dataSourceIBMPrivateDNSForwardrulesID() string {
	return time.Now().UTC().String()
}
