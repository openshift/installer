// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cisWebhookList = "cis_webhooks"

func DataSourceIBMCISWebhooks() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCISWebhookRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisWebhookList: {
				Type:        schema.TypeList,
				Description: "Collection of Webhook details",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisWebhookID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Webhook ID",
						},
						cisWebhookName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Webhook Name",
						},
						cisWebhookURL: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Webhook URL",
						},
						cisWebhookType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Webhook Type",
						},
					},
				},
			},
		},
	}
}
func dataIBMCISWebhookRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisWebhookSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the cisWebhookession %s", err)
	}
	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)
	opt := sess.NewListWebhooksOptions()

	result, resp, err := sess.ListWebhooks(opt)
	if err != nil || result == nil {
		return fmt.Errorf("[ERROR] Error Listing all Webhooks %q: %s %s", d.Id(), err, resp)
	}

	webhooks := make([]map[string]interface{}, 0)

	for _, instance := range result.Result {
		webhook := map[string]interface{}{}
		webhook[cisWebhookID] = *instance.ID
		webhook[cisWebhookName] = *instance.Name
		webhook[cisWebhookURL] = *instance.URL
		webhook[cisWebhookType] = *instance.Type
		webhooks = append(webhooks, webhook)
	}
	d.SetId(dataSourcecisWebhookCheckID(d))
	d.Set(cisID, crn)
	d.Set(cisWebhookList, webhooks)
	return nil
}

func dataSourcecisWebhookCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
