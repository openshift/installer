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

func ResourceIBMEnSlackSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnSlackSubscriptionCreate,
		ReadContext:   resourceIBMEnSlackSubscriptionRead,
		UpdateContext: resourceIBMEnSlackSubscriptionUpdate,
		DeleteContext: resourceIBMEnSlackSubscriptionDelete,
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
						"attachment_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "attachment color code",
						},
						"template_id_notification": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The templete id for notification",
						},
						"channels": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of channels.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "channel id.",
									},
									"operation": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The channel operation type. The values are add/remove",
									},
								},
							},
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
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time.",
			},
		},
	}
}

func resourceIBMEnSlackSubscriptionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	attributes, _ := slackattributesMapToAttributes(d.Get("attributes.0").(map[string]interface{}))
	options.SetAttributes(&attributes)

	result, response, err := enClient.CreateSubscriptionWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("CreateSubscriptionWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnSlackSubscriptionRead(context, d, meta)
}

func resourceIBMEnSlackSubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceIBMEnSlackSubscriptionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if ok := d.HasChanges("name", "description", "attributes"); ok {
		options.SetName(d.Get("name").(string))

		if _, ok := d.GetOk("description"); ok {
			options.SetDescription(d.Get("description").(string))
		}

		attributes, err := resourceIBMEnSubscriptionMapToSubscriptionUpdateAttributes(d.Get("attributes.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		options.SetAttributes(&attributes)

		_, response, err := enClient.UpdateSubscriptionWithContext(context, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("UpdateSubscriptionWithContext failed %s\n%s", err, response))
		}

		return resourceIBMEnSlackSubscriptionRead(context, d, meta)
	}

	return nil
}

func resourceIBMEnSlackSubscriptionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func slackattributesMapToAttributes(modelMap map[string]interface{}) (en.SubscriptionCreateAttributes, error) {
	model := en.SubscriptionCreateAttributes{}

	if modelMap["template_id_notification"] != nil && modelMap["template_id_notification"].(string) != "" {
		model.TemplateIDNotification = core.StringPtr(modelMap["template_id_notification"].(string))
	}
	if modelMap["template_id_invitation"] != nil && modelMap["template_id_invitation"].(string) != "" {
		model.TemplateIDInvitation = core.StringPtr(modelMap["template_id_invitation"].(string))
	}

	if modelMap["attachment_color"] != nil && modelMap["attachment_color"].(string) != "" {
		model.AttachmentColor = core.StringPtr(modelMap["attachment_color"].(string))
	}
	if modelMap["channels"] != nil {
		channels := []en.ChannelCreateAttributes{}
		for _, channelsItem := range modelMap["channels"].([]interface{}) {
			channelsItemModel, err := resourceIBMEnSubscriptionMapToChannelCreateAttributes(channelsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			channels = append(channels, *channelsItemModel)
		}
		model.Channels = channels
	}
	return model, nil
}

func resourceIBMEnSubscriptionMapToChannelCreateAttributes(modelMap map[string]interface{}) (*en.ChannelCreateAttributes, error) {
	model := &en.ChannelCreateAttributes{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMEnSubscriptionMapToSubscriptionUpdateAttributes(modelMap map[string]interface{}) (en.SubscriptionUpdateAttributes, error) {
	model := en.SubscriptionUpdateAttributes{}

	if modelMap["template_id_notification"] != nil && modelMap["template_id_notification"].(string) != "" {
		model.TemplateIDNotification = core.StringPtr(modelMap["template_id_notification"].(string))
	}

	if modelMap["attachment_color"] != nil && modelMap["attachment_color"].(string) != "" {
		model.AttachmentColor = core.StringPtr(modelMap["attachment_color"].(string))
	}
	if modelMap["channels"] != nil {
		channels := []en.ChannelUpdateAttributes{}
		for _, channelsItem := range modelMap["channels"].([]interface{}) {
			channelsItemModel, err := resourceIBMEnSubscriptionMapToChannelUpdateAttributes(channelsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			channels = append(channels, *channelsItemModel)
		}
		model.Channels = channels
	}
	return model, nil
}

func resourceIBMEnSubscriptionChannelCreateAttributesToMap(model *en.ChannelCreateAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceIBMEnSubscriptionMapToChannelUpdateAttributes(modelMap map[string]interface{}) (*en.ChannelUpdateAttributes, error) {
	model := &en.ChannelUpdateAttributes{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	model.Operation = core.StringPtr(modelMap["operation"].(string))
	return model, nil
}
