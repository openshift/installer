// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMPrivateDNSForwardingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIbmDnsCrForwardingRulesRead,

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
					},
				},
			},
		},
	}
}

func dataSourceIbmDnsCrForwardingRulesRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	resolverID := d.Get(pdnsCRFRResolverID).(string)

	opt := sess.NewListForwardingRulesOptions(instanceID, resolverID)

	result, resp, err := sess.ListForwardingRulesWithContext(context.TODO(), opt)
	if err != nil || result == nil {
		return fmt.Errorf("Error listing the forwarding rules %s:%s", err, resp)
	}

	forwardRules := make([]interface{}, 0)
	for _, instance := range result.ForwardingRules {
		forwardRule := map[string]interface{}{}
		forwardRule[pdnsCRFRRuleID] = *instance.ID
		forwardRule[pdnsCRFRDesctiption] = *instance.Description
		forwardRule[pdnsCRFRType] = *instance.Type
		forwardRule[pdnsCRFRMatch] = *instance.Match
		forwardRule[pdnsCRFRForwardTo] = instance.ForwardTo

		forwardRules = append(forwardRules, forwardRule)
	}
	d.SetId(dataSourceIBMPrivateDNSForwardrulesID(d))
	d.Set(pdnsInstanceID, instanceID)
	d.Set(pdnsCRFRResolverID, resolverID)
	d.Set(pdnsCRForwardRules, forwardRules)
	return nil
}

func dataSourceIBMPrivateDNSForwardrulesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
