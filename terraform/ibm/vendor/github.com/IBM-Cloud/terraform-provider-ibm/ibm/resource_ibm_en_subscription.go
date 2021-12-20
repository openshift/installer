// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func resourceIBMEnSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnSubscriptionCreate,
		ReadContext:   resourceIBMEnSubscriptionRead,
		UpdateContext: resourceIBMEnSubscriptionUpdate,
		DeleteContext: resourceIBMEnSubscriptionDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subscription name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subscription description.",
			},
			"destination_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination ID.",
			},
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Topic ID.",
			},
			"attributes": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"to": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The phone number to send the SMS to in case of sms_ibm. The email id in case of smtp_ibm destination type.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"add_notification_payload": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to add the notification payload to the email.",
						},
						"reply_to": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The email address to reply to.",
						},
						"signing_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Signing webhook attributes.",
						},
					},
				},
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subscription ID.",
			},
			"destination_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Destination Webhook.",
			},
			"destination_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Destintion name.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the topic.",
			},
			"from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "From Email ID (it will be displayed only in case of smtp_ibm destination type).",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
		},
	}
}

func resourceIBMEnSubscriptionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.CreateSubscriptionOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))

	options.SetName(d.Get("name").(string))
	options.SetTopicID(d.Get("topic_id").(string))
	options.SetDestinationID(d.Get("destination_id").(string))

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}

	attributes, _ := attributesMapToAttributes(d.Get("attributes.0").(map[string]interface{}))
	options.SetAttributes(&attributes)

	result, response, err := enClient.CreateSubscriptionWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("CreateSubscriptionWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnSubscriptionRead(context, d, meta)
}

func resourceIBMEnSubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.GetSubscriptionOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	result, response, err := enClient.GetSubscriptionWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("GetSubscriptionWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_guid", options.InstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting instance_guid: %s", err))
	}

	if err = d.Set("subscription_id", result.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting instance_guid: %s", err))
	}

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return diag.FromErr(fmt.Errorf("error setting description: %s", err))
		}
	}

	if result.From != nil {
		if err = d.Set("from", result.From); err != nil {
			return diag.FromErr(fmt.Errorf("error setting from: %s", err))
		}
	}

	if err = d.Set("destination_id", result.DestinationID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting destination_id: %s", err))
	}

	if err = d.Set("destination_type", result.DestinationType); err != nil {
		return diag.FromErr(fmt.Errorf("error setting destination_type: %s", err))
	}

	if result.DestinationName != nil {
		if err = d.Set("destination_name", result.DestinationName); err != nil {
			return diag.FromErr(fmt.Errorf("error setting destination_name: %s", err))
		}
	}

	if err = d.Set("topic_id", result.TopicID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting topic_id: %s", err))
	}

	if result.TopicName != nil {
		if err = d.Set("topic_name", result.TopicName); err != nil {
			return diag.FromErr(fmt.Errorf("error setting topic_name: %s", err))
		}
	}

	if err = d.Set("updated_at", result.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting updated_at: %s", err))
	}

	return nil
}

func resourceIBMEnSubscriptionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.UpdateSubscriptionOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	if ok := d.HasChanges("name", "description", "attributes"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		_, attributes := attributesMapToAttributes(d.Get("attributes.0").(map[string]interface{}))
		options.SetAttributes(&attributes)

		_, response, err := enClient.UpdateSubscriptionWithContext(context, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("UpdateSubscriptionWithContext failed %s\n%s", err, response))
		}

		return resourceIBMEnSubscriptionRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnSubscriptionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.DeleteSubscriptionOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	response, err := enClient.DeleteSubscriptionWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("DeleteSubscriptionWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func attributesMapToAttributes(attributeMap map[string]interface{}) (en.SubscriptionCreateAttributes, en.SubscriptionUpdateAttributes) {
	attributesCreate := en.SubscriptionCreateAttributes{}
	attributesUpdate := en.SubscriptionUpdateAttributes{}

	if attributeMap["to"] != nil {
		to := []string{}
		for _, toItem := range attributeMap["to"].([]interface{}) {
			to = append(to, toItem.(string))
		}
		attributesCreate.To = to
		attributesUpdate.To = to
	}

	if attributeMap["add_notification_payload"] != nil {
		attributesCreate.AddNotificationPayload = core.BoolPtr(attributeMap["add_notification_payload"].(bool))
		attributesUpdate.AddNotificationPayload = core.BoolPtr(attributeMap["add_notification_payload"].(bool))
	}

	if attributeMap["reply_to"] != nil {
		attributesCreate.ReplyTo = core.StringPtr(attributeMap["reply_to"].(string))
		attributesUpdate.ReplyTo = core.StringPtr(attributeMap["reply_to"].(string))
	}

	if attributeMap["signing_enabled"] != nil {
		attributesCreate.SigningEnabled = core.BoolPtr(attributeMap["signing_enabled"].(bool))
		attributesUpdate.SigningEnabled = core.BoolPtr(attributeMap["signing_enabled"].(bool))
	}

	return attributesCreate, attributesUpdate
}
