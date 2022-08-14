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

func DataSourceIBMEnSubscription() *schema.Resource {
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
						"add_notification_payload": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to add the notification payload to the email.",
						},
						"signing_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Signing webhook attributes.",
						},
						"additionalproperties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Additional attributes for sms and webhook subscription.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"to": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The phone number to send the SMS to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"additional_properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: "Additional attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reply_to_mail": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The email address to reply to.",
									},
									"reply_to_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The email address to reply to.",
									},
									"from_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The email address user name to reply to.",
									},
									"to": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The email id in case of smtp_ibm destination type.",
										Elem:        &schema.Schema{Type: schema.TypeMap},
									},
									"invited": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The email id in case of smtp_ibm destination type.",
										Elem:        &schema.Schema{Type: schema.TypeMap},
									},
									"unsubscribed": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The Email address which should be unsubscribed from smtp_ibm.",
										Elem:        &schema.Schema{Type: schema.TypeMap},
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
		DeprecationMessage: "This data source would be deprecated. For subscription data sources for email, sms, webhook destination kindly use ibm_en_subscription_email, ibm_en_subscription_sms, ibm_en_subscription_webhook subscription data sources",
	}
}

func dataSourceIBMEnSubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if result.DestinationType != nil {
		if err = d.Set("destination_type", result.DestinationType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_type: %s", err))
		}
	}
	destinationtype := d.Get("destination_type").(string)

	if err = d.Set("destination_id", result.DestinationID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_id: %s", err))
	}

	if err = d.Set("destination_name", result.DestinationName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting destination_name: %s", err))
	}

	if err = d.Set("topic_id", result.TopicID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting topic_id: %s", err))
	}

	if err = d.Set("topic_name", result.TopicName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting topic_name: %s", err))
	}

	if result.Attributes != nil {
		if err = d.Set("attributes", enSubscriptionFlattenAttributes(result.Attributes, destinationtype)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting attributes %s", err))
		}
	}

	if result.From != nil {
		if err = d.Set("from", result.From); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting from %s", err))
		}
	}

	return nil
}

func enSubscriptionFlattenAttributes(result en.SubscriptionAttributesIntf, destinationtype string) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}

	attributes := result.(*en.SubscriptionAttributes)

	destination_type := destinationtype

	finalMap := enSubscriptionToMap(attributes, destination_type)
	finalList = append(finalList, finalMap)

	return finalList
}

func enSubscriptionToMap(attributeItem *en.SubscriptionAttributes, destinationtype string) (attributeMap map[string]interface{}) {
	attributeMap = map[string]interface{}{}

	if attributeItem.AddNotificationPayload != nil {
		attributeMap["add_notification_payload"] = attributeItem.AddNotificationPayload
	}

	if attributeItem.SigningEnabled != nil {
		attributeMap["signing_enabled"] = attributeItem.SigningEnabled
	}

	if destinationtype == "smtp_ibm" {

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
	} else if destinationtype == "sms_ibm" {
		prop := []map[string]interface{}{}

		b := attributeItem.GetProperties()
		m := make(map[string]interface{})
		if len(b) > 0 {
			for k, v := range b {
				m[k] = v
			}
		}
		prop = append(prop, m)
		attributeMap["additionalproperties"] = prop
	}

	return attributeMap
}
