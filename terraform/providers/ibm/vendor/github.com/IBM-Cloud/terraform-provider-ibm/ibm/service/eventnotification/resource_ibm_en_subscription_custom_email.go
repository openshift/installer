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
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMEnCustomEmailSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnCustomEmailSubscriptionCreate,
		ReadContext:   resourceIBMEnCustomEmailSubscriptionRead,
		UpdateContext: resourceIBMEnCustomEmailSubscriptionUpdate,
		DeleteContext: resourceIBMEnCustomEmailSubscriptionDelete,
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
				ForceNew:    true,
				Description: "Destination ID.",
			},
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Topic ID.",
			},
			"attributes": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add_notification_payload": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to add the notification payload to the email.",
						},
						"reply_to_mail": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The email address to reply to.",
						},
						"reply_to_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The  name of the email address user to reply to.",
						},
						"from_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of email address from which email is sourced.",
						},
						"from_email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The email from where it is sourced",
						},
						"template_id_notification": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The templete id for notification",
						},
						"template_id_invitation": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The templete id for invitation",
						},
						"invited": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Email address send the invite to in case of smtp_ibm.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"add": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Email address which should be added to smtp_ibm.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"remove": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The email id to be removed in case of smtp_ibm destination type.",
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
				Description: "The type of Destination.",
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

func resourceIBMEnCustomEmailSubscriptionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription_custom_email", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.CreateSubscriptionOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))

	options.SetName(d.Get("name").(string))
	options.SetTopicID(d.Get("topic_id").(string))
	options.SetDestinationID(d.Get("destination_id").(string))

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}

	attributes := CustomEmailattributesMapToAttributes(d.Get("attributes.0").(map[string]interface{}))
	options.SetAttributes(&attributes)

	result, _, err := enClient.CreateSubscriptionWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSubscriptionWithContext failed: %s", err.Error()), "ibm_en_subscription_custom_email", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnCustomEmailSubscriptionRead(context, d, meta)
}

func resourceIBMEnCustomEmailSubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription_custom_email", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.GetSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription_custom_email", "read")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	result, response, err := enClient.GetSubscriptionWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSubscriptionWithContext failed: %s", err.Error()), "ibm_en_subscription_custom_email", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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

func resourceIBMEnCustomEmailSubscriptionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription_custom_email", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.UpdateSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription_custom_email", "update")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	if ok := d.HasChanges("name", "description", "attributes"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		attributes := CustomEmailattributesupdateMapToAttributes(d.Get("attributes.0").(map[string]interface{}))
		options.SetAttributes(&attributes)

		_, _, err := enClient.UpdateSubscriptionWithContext(context, options)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSubscriptionWithContext failed: %s", err.Error()), "ibm_en_subscription_custom_email", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		return resourceIBMEnCustomEmailSubscriptionRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnCustomEmailSubscriptionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription_custom_email", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.DeleteSubscriptionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription_custom_email", "delete")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	response, err := enClient.DeleteSubscriptionWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSubscriptionWithContext: failed: %s", err.Error()), "ibm_en_subscription_custom_email", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func CustomEmailattributesMapToAttributes(attributeMap map[string]interface{}) en.SubscriptionCreateAttributes {
	attributesCreate := en.SubscriptionCreateAttributes{}
	if attributeMap["invited"] != nil {
		invited := []string{}
		for _, invitedItem := range attributeMap["invited"].([]interface{}) {
			invited = append(invited, invitedItem.(string))
		}
		attributesCreate.Invited = invited
	}

	if attributeMap["add_notification_payload"] != nil {
		attributesCreate.AddNotificationPayload = core.BoolPtr(attributeMap["add_notification_payload"].(bool))
	}

	if attributeMap["reply_to_mail"] != nil {
		attributesCreate.ReplyToMail = core.StringPtr(attributeMap["reply_to_mail"].(string))
	}

	if attributeMap["reply_to_name"] != nil {
		attributesCreate.ReplyToName = core.StringPtr(attributeMap["reply_to_name"].(string))
	}

	if attributeMap["from_name"] != nil {
		attributesCreate.FromName = core.StringPtr(attributeMap["from_name"].(string))
	}

	if attributeMap["from_email"] != nil {
		attributesCreate.FromEmail = core.StringPtr(attributeMap["from_email"].(string))
	}

	if attributeMap["template_id_notification"] != nil {
		attributesCreate.TemplateIDNotification = core.StringPtr(attributeMap["template_id_notification"].(string))
	}

	if attributeMap["template_id_invitation"] != nil {
		attributesCreate.TemplateIDInvitation = core.StringPtr(attributeMap["template_id_invitation"].(string))
	}

	return attributesCreate
}

func CustomEmailattributesupdateMapToAttributes(attributeMap map[string]interface{}) en.SubscriptionUpdateAttributesCustomEmailUpdateAttributes {
	updateattributes := en.SubscriptionUpdateAttributesCustomEmailUpdateAttributes{}

	addemail := new(en.UpdateAttributesInvited)
	if attributeMap["add"] != nil {
		to := []string{}
		for _, toItem := range attributeMap["add"].([]interface{}) {
			to = append(to, toItem.(string))
		}
		addemail.Add = to
	}
	updateattributes.Invited = addemail

	if attributeMap["remove"] != nil {
		rmemail := []string{}
		for _, removeitem := range attributeMap["remove"].([]interface{}) {
			rmemail = append(rmemail, removeitem.(string))
		}

		addemail.Remove = rmemail
	}
	updateattributes.Invited = addemail

	if attributeMap["add_notification_payload"] != nil {
		updateattributes.AddNotificationPayload = core.BoolPtr(attributeMap["add_notification_payload"].(bool))
	}

	if attributeMap["reply_to_mail"] != nil {
		updateattributes.ReplyToMail = core.StringPtr(attributeMap["reply_to_mail"].(string))
	}

	if attributeMap["reply_to_name"] != nil {
		updateattributes.ReplyToName = core.StringPtr(attributeMap["reply_to_name"].(string))
	}

	if attributeMap["from_name"] != nil {
		updateattributes.FromName = core.StringPtr(attributeMap["from_name"].(string))
	}

	if attributeMap["from_email"] != nil {
		updateattributes.FromEmail = core.StringPtr(attributeMap["from_email"].(string))
	}

	if attributeMap["template_id_notification"] != nil {
		updateattributes.TemplateIDNotification = core.StringPtr(attributeMap["template_id_notification"].(string))
	}

	if attributeMap["template_id_invitation"] != nil {
		updateattributes.TemplateIDInvitation = core.StringPtr(attributeMap["template_id_invitation"].(string))
	}

	return updateattributes
}
