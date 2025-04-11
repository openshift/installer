// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventnotification

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
)

func DataSourceIBMEnCustomEmailSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnCustomEmailSubscriptionRead,

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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add_notification_payload": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to add the notification payload to the email.",
						},
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
							Description: "The email address user name to reply to.",
						},
						"from_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The email address username of source email address.",
						},
						"from_email": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The email from where it is sourced",
						},
						"template_id_notification": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The templete id for notification",
						},
						"template_id_invitation": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The templete id for invitation",
						},
						"invited": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The invited item schema",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The email address to reply to.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The updated date of invitation",
									},
									"expires_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The expiry date of invitation mail",
									},
								},
							},
						},
						"subscribed": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The email subscribed items schema",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The email address to reply to.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The updated date of susbcription",
									},
								},
							},
						},
						"unsubscribed": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The unsusbscribed email items schema",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The email address to reply to.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The updated date of unsusbcription",
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

func dataSourceIBMEnCustomEmailSubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_subscription_custom_email", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSubscriptionOptions := &en.GetSubscriptionOptions{}

	getSubscriptionOptions.SetInstanceID(d.Get("instance_guid").(string))
	getSubscriptionOptions.SetID(d.Get("subscription_id").(string))

	result, _, err := enClient.GetSubscriptionWithContext(context, getSubscriptionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSubscriptionWithContext failed: %s", err.Error()), "(Data) ibm_en_subscription_custom_email", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getSubscriptionOptions.InstanceID, *getSubscriptionOptions.ID))

	if err = d.Set("name", result.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_subscription_custom_email", "read")
		return tfErr.GetDiag()
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_en_subscription_custom_email", "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("updated_at", result.UpdatedAt); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_en_subscription_custom_email", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("destination_id", result.DestinationID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting destination_id: %s", err), "(Data) ibm_en_subscription_custom_email", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("topic_id", result.TopicID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting topic_id: %s", err), "(Data) ibm_en_subscription_custom_email", "read")
		return tfErr.GetDiag()
	}

	if result.Attributes != nil {
		if err = d.Set("attributes", enCustomEmailSubscriptionFlattenAttributes(result.Attributes)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf(" Error setting attributes %s", err), "(Data) ibm_en_subscription_custom_email", "read")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func enCustomEmailSubscriptionFlattenAttributes(result en.SubscriptionAttributesIntf) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}

	attributes := result.(*en.SubscriptionAttributes)

	finalMap := enCustomEmailSubscriptionToMap(attributes)
	finalList = append(finalList, finalMap)

	return finalList
}

func enCustomEmailSubscriptionToMap(attributeItem *en.SubscriptionAttributes) (attributeMap map[string]interface{}) {
	attributeMap = map[string]interface{}{}

	if attributeItem.AddNotificationPayload != nil {
		attributeMap["add_notification_payload"] = attributeItem.AddNotificationPayload
	}
	if attributeItem.ReplyToMail != nil {
		attributeMap["reply_to_mail"] = attributeItem.ReplyToMail
	}
	if attributeItem.ReplyToName != nil {
		attributeMap["reply_to_name"] = attributeItem.ReplyToName
	}
	if attributeItem.FromName != nil {
		attributeMap["from_name"] = attributeItem.FromName
	}
	if attributeItem.FromEmail != nil {
		attributeMap["from_email"] = attributeItem.FromEmail
	}
	if attributeItem.TemplateIDNotification != nil {
		attributeMap["template_id_notification"] = attributeItem.TemplateIDNotification
	}
	if attributeItem.TemplateIDInvitation != nil {
		attributeMap["template_id_invitation"] = attributeItem.TemplateIDInvitation
	}
	return attributeMap
}
