// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMEnSubscription() *schema.Resource {
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
				Optional: true,
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
						"reply_to_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the email address user to reply to.",
						},
						"from_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The email address from.",
						},
						"signing_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Signing webhook attributes.",
						},
						"remove": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Email address to remove the recepient from smtp_ibm.",
							Elem:        &schema.Schema{Type: schema.TypeString},
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
		DeprecationMessage: "This resource will be deprecated. To create subscription for email, sms, webhook destination kindly use ibm_en_subscription_email, ibm_en_subscription_sms, ibm_en_subscription_webhook subscription resources",
	}
}

func resourceIBMEnSubscriptionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
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
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.GetSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
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
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("subscription_id", result.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_guid: %s", err))
	}

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
		}
	}

	if result.From != nil {
		if err = d.Set("from", result.From); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting from: %s", err))
		}
	}

	if err = d.Set("destination_id", result.DestinationID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_id: %s", err))
	}

	if err = d.Set("destination_type", result.DestinationType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_type: %s", err))
	}

	if result.DestinationName != nil {
		if err = d.Set("destination_name", result.DestinationName); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_name: %s", err))
		}
	}

	if err = d.Set("topic_id", result.TopicID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting topic_id: %s", err))
	}

	if result.TopicName != nil {
		if err = d.Set("topic_name", result.TopicName); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting topic_name: %s", err))
		}
	}

	if err = d.Set("updated_at", result.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	return nil
}

func resourceIBMEnSubscriptionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.UpdateSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	destinationtype := d.Get("destination_type").(string)

	if ok := d.HasChanges("name", "description", "attributes"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		if destinationtype == "smtp_ibm" {
			attributes := attributesMapToEmailAttributes(d.Get("attributes.0").(map[string]interface{}))
			options.SetAttributes(&attributes)
		} else {
			_, attributes := attributesMapToAttributes(d.Get("attributes.0").(map[string]interface{}))
			options.SetAttributes(&attributes)
		}

		_, response, err := enClient.UpdateSubscriptionWithContext(context, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("UpdateSubscriptionWithContext failed %s\n%s", err, response))
		}

		return resourceIBMEnSubscriptionRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnSubscriptionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.DeleteSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
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
		attributesCreate.ReplyToMail = core.StringPtr(attributeMap["reply_to"].(string))
	}

	if attributeMap["signing_enabled"] != nil {
		attributesCreate.SigningEnabled = core.BoolPtr(attributeMap["signing_enabled"].(bool))
		attributesUpdate.SigningEnabled = core.BoolPtr(attributeMap["signing_enabled"].(bool))
	}

	if attributeMap["reply_to_name"] != nil {
		attributesCreate.ReplyToName = core.StringPtr(attributeMap["reply_to_name"].(string))
	}

	if attributeMap["from_name"] != nil {
		attributesCreate.FromName = core.StringPtr(attributeMap["from_name"].(string))
	}

	return attributesCreate, attributesUpdate
}

func attributesMapToEmailAttributes(attributeMap map[string]interface{}) en.SubscriptionUpdateAttributesEmailUpdateAttributes {
	updateattributes := en.SubscriptionUpdateAttributesEmailUpdateAttributes{}

	addemail := new(en.EmailUpdateAttributesTo)
	if attributeMap["to"] != nil {
		to := []string{}
		for _, toItem := range attributeMap["to"].([]interface{}) {
			to = append(to, toItem.(string))
		}
		addemail.Add = to
	}
	updateattributes.To = addemail

	if attributeMap["add_notification_payload"] != nil {
		updateattributes.AddNotificationPayload = core.BoolPtr(attributeMap["add_notification_payload"].(bool))
	}

	if attributeMap["reply_to"] != nil {
		updateattributes.ReplyToMail = core.StringPtr(attributeMap["reply_to"].(string))
	}

	if attributeMap["reply_to_name"] != nil {
		updateattributes.ReplyToName = core.StringPtr(attributeMap["reply_to_name"].(string))
	}

	if attributeMap["from_name"] != nil {
		updateattributes.FromName = core.StringPtr(attributeMap["from_name"].(string))
	}

	if attributeMap["remove"] != nil {
		removed := []string{}
		for _, removedItem := range attributeMap["remove"].([]interface{}) {
			removed = append(removed, removedItem.(string))
		}
		addemail.Remove = removed
	}
	unsubscribed := new(en.EmailUpdateAttributesUnsubscribed)
	if attributeMap["unsubscribed"] != nil {
		unsubscribe := []string{}
		for _, unsubscribeItem := range attributeMap["unsubscribed"].([]interface{}) {
			unsubscribe = append(unsubscribe, unsubscribeItem.(string))
		}
		unsubscribed.Remove = unsubscribe
	}
	updateattributes.Unsubscribed = unsubscribed

	return updateattributes
}
