// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnSMSSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSMSSubscriptionRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for result.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscription name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscription description.",
			},
			"destination_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The destination ID.",
			},
			"topic_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Topic ID.",
			},
			"attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The additional attributes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Additional attributes for sms and webhook subscription.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"to": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The phone number to send the SMS to.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
		},
	}
}

func dataSourceIBMEnSMSSubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getSubscriptionOptions := &en.GetSubscriptionOptions{}

	getSubscriptionOptions.SetInstanceID(d.Get("instance_guid").(string))
	getSubscriptionOptions.SetID(d.Get("subscription_id").(string))

	result, response, err := enClient.GetSubscriptionWithContext(context, getSubscriptionOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("GetSubscriptionWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getSubscriptionOptions.InstanceID, *getSubscriptionOptions.ID))

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
		}
	}
	if err = d.Set("updated_at", result.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	if err = d.Set("destination_id", result.DestinationID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_id: %s", err))
	}

	if err = d.Set("topic_id", result.TopicID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting topic_id: %s", err))
	}

	if result.Attributes != nil {
		if err = d.Set("attributes", enSMSSubscriptionFlattenAttributes(result.Attributes)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting attributes %s", err))
		}
	}

	return nil
}

func enSMSSubscriptionFlattenAttributes(result en.SubscriptionAttributesIntf) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}

	attributes := result.(*en.SubscriptionAttributes)

	finalMap := enSMSSubscriptionToMap(attributes)
	finalList = append(finalList, finalMap)

	return finalList
}

func enSMSSubscriptionToMap(attributeItem *en.SubscriptionAttributes) (attributeMap map[string]interface{}) {
	attributeMap = map[string]interface{}{}

	prop := []map[string]interface{}{}

	b := attributeItem.GetProperties()
	m := make(map[string]interface{})
	if len(b) > 0 {
		for k, v := range b {
			m[k] = v
		}
	}
	prop = append(prop, m)
	attributeMap["additional_properties"] = prop

	return attributeMap
}
