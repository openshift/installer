// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func dataSourceIBMEnSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnSubscriptionRead,

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
			"destination_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of destination.",
			},
			"destination_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The destination ID.",
			},
			"destination_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The destination name.",
			},
			"topic_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Topic ID.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Topic name.",
			},
			"from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "From Email ID (it will be displayed only in case of smtp_ibm destination type).",
			},
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
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
						"recipient_selection": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The recipient selection method.",
						},
						"add_notification_payload": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to add the notification payload to the email.",
						},
						"reply_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The email address to reply to.",
						},
						"signing_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Signing webhook attributes.",
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

func dataSourceIBMEnSubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(ClientSession).EventNotificationsApiV1()
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
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return diag.FromErr(fmt.Errorf("error setting description: %s", err))
		}
	}
	if err = d.Set("updated_at", result.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting updated_at: %s", err))
	}

	if result.DestinationType != nil {
		if err = d.Set("destination_type", result.DestinationType); err != nil {
			return diag.FromErr(fmt.Errorf("error setting destination_type: %s", err))
		}
	}

	if err = d.Set("destination_id", result.DestinationID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting destination_id: %s", err))
	}

	if err = d.Set("destination_name", result.DestinationName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting destination_name: %s", err))
	}

	if err = d.Set("topic_id", result.TopicID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting topic_id: %s", err))
	}

	if err = d.Set("topic_name", result.TopicName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting topic_name: %s", err))
	}

	if result.Attributes != nil {
		if err = d.Set("attributes", enSubscriptionFlattenAttributes(result.Attributes)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting attributes %s", err))
		}
	}

	if result.From != nil {
		if err = d.Set("from", result.From); err != nil {
			return diag.FromErr(fmt.Errorf("error setting from %s", err))
		}
	}

	return nil
}

func enSubscriptionFlattenAttributes(result en.SubscriptionAttributesIntf) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}

	attributes := result.(*en.SubscriptionAttributes)

	finalMap := enSubscriptionToMap(attributes)
	finalList = append(finalList, finalMap)

	return finalList
}

func enSubscriptionToMap(attributeItem *en.SubscriptionAttributes) (attributeMap map[string]interface{}) {
	attributeMap = map[string]interface{}{}

	if attributeItem.AddNotificationPayload != nil {
		attributeMap["add_notification_payload"] = attributeItem.AddNotificationPayload
	}

	if attributeItem.RecipientSelection != nil {
		attributeMap["recipient_selection"] = attributeItem.RecipientSelection
	}

	if attributeItem.ReplyTo != nil {
		attributeMap["reply_to"] = attributeItem.ReplyTo
	}

	if attributeItem.To != nil {
		attributeMap["to"] = attributeItem.To
	}

	if attributeItem.SigningEnabled != nil {
		attributeMap["signing_enabled"] = attributeItem.SigningEnabled
	}

	return attributeMap
}
