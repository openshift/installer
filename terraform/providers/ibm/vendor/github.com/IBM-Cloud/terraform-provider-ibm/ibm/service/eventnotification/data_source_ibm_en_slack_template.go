// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMEnSlackTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSlackTemplateRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for Template.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Templaten description.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template type slack.notification.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
			"subscription_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of subscriptions.",
			},
			"subscription_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of subscriptions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIBMEnSlackTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.GetTemplateOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("template_id").(string))

	result, response, err := enClient.GetTemplateWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("GetTemplate failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *options.ID))

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
		}
	}

	if err = d.Set("type", result.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}

	if result.SubscriptionNames != nil {
		err = d.Set("subscription_names", result.SubscriptionNames)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_names %s", err))
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	if err = d.Set("subscription_count", flex.IntValue(result.SubscriptionCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_count: %s", err))
	}

	return nil
}
