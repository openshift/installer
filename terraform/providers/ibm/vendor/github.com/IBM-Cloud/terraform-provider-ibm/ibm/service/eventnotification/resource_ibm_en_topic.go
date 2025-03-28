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

	"github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	en "github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMEnTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMEnTopicCreate,
		ReadContext:   resourceIBMEnTopicRead,
		UpdateContext: resourceIBMEnTopicUpdate,
		DeleteContext: resourceIBMEnTopicDelete,
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
				Description: "Name of the topic.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the topic.",
			},
			"sources": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of sources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the source.",
						},
						"rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Whether the rule is enabled or not.",
									},
									"event_type_filter": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Event type filter.",
									},
									"notification_filter": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "",
										Description: "Notification filter.",
									},
									"event_schedule_filter": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Event Schedule Filter for periodic-timer source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"starts_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "event schedule start time.",
												},
												"ends_at": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "event schedule end time.",
												},
												"expression": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "cron schedule expression.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"topic_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Topic ID.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last time the topic was updated.",
			},
			"source_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of sources.",
			},
			"subscription_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of subscriptions.",
			},
			"subscriptions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of subscriptions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subscription.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the subscription.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the subscription.",
						},
						"destination_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the destination.",
						},
						"destination_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of destination.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the topic.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last updated time of the subscription.",
						},
					},
				},
			},
		},
	}
}

func resourceIBMEnTopicCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_subscription_webhook", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.CreateTopicOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetName(d.Get("name").(string))

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}

	if _, ok := d.GetOk("sources"); ok {
		var sources []en.SourcesItems
		for _, e := range d.Get("sources").([]interface{}) {
			value := e.(map[string]interface{})
			sourcesItem := enTopicUpdateSourcesItem(value)
			sources = append(sources, sourcesItem)
		}
		options.SetSources(sources)
	}

	result, _, err := enClient.CreateTopicWithContext(context, options)

	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateTopicWithContext failed: %s", err.Error()), "ibm_en_topic", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *result.ID))

	return resourceIBMEnTopicRead(context, d, meta)
}

func resourceIBMEnTopicRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_topic", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.GetTopicOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_topic", "read")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])
	options.SetHeaders(map[string]string{"include": "subscriptions"})

	result, response, err := enClient.GetTopicWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTopicWithContext failed: %s", err.Error()), "ibm_en_topic", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.Set("instance_guid", options.InstanceID)
	d.Set("topic_id", result.ID)

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if err = d.Set("description", result.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}

	if result.Sources != nil {
		sources := []map[string]interface{}{}
		for _, sourcesItem := range result.Sources {
			sourcesItemMap := enTopicUpdateSourcesItemToMap(sourcesItem)
			sources = append(sources, sourcesItemMap)
		}
		if err = d.Set("sources", sources); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting sources: %s", err))
		}
	}

	if err = d.Set("updated_at", result.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	if err = d.Set("source_count", flex.IntValue(result.SourceCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting source_count: %s", err))
	}

	if err = d.Set("subscription_count", flex.IntValue(result.SubscriptionCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_count: %s", err))
	}

	subscriptions := []map[string]interface{}{}
	for _, subscriptionsItem := range result.Subscriptions {
		subscriptionsItemMap := enTopicSubscriptionToMap(subscriptionsItem)
		subscriptions = append(subscriptions, subscriptionsItemMap)
	}

	if err = d.Set("subscriptions", subscriptions); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscriptions: %s", err))
	}

	return nil
}

func resourceIBMEnTopicUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_topic", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.ReplaceTopicOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_topic", "update")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])
	options.SetInstanceID(d.Get("instance_guid").(string))

	options.SetName(d.Get("name").(string))

	if _, ok := d.GetOk("description"); ok {
		options.SetDescription(d.Get("description").(string))
	}

	if _, ok := d.GetOk("sources"); ok {
		var sources []en.SourcesItems
		for _, e := range d.Get("sources").([]interface{}) {
			value := e.(map[string]interface{})
			sourcesItem := enTopicUpdateSourcesItem(value)
			sources = append(sources, sourcesItem)
		}
		options.SetSources(sources)
	}

	_, _, err = enClient.ReplaceTopicWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceTopicWithContext failed: %s", err.Error()), "ibm_en_topic", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMEnTopicRead(context, d, meta)
}

func resourceIBMEnTopicDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_topic", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.DeleteTopicOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_en_topic", "delete")
		return tfErr.GetDiag()
	}

	options.SetInstanceID(parts[0])
	options.SetID(parts[1])

	response, err := enClient.DeleteTopicWithContext(context, options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteTopicWithContext: failed: %s", err.Error()), "ibm_en_topic", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func enTopicUpdateSourcesItem(topicUpdateSourcesItemMap map[string]interface{}) en.SourcesItems {
	topicUpdateSourcesItem := en.SourcesItems{}

	if topicUpdateSourcesItemMap["id"] != nil {
		topicUpdateSourcesItem.ID = core.StringPtr(topicUpdateSourcesItemMap["id"].(string))
	}
	if topicUpdateSourcesItemMap["rules"] != nil {
		rules := []en.Rules{}
		for _, rulesItem := range topicUpdateSourcesItemMap["rules"].([]interface{}) {
			rulesItemModel := resourceIBMEnTopicMapToRules(rulesItem.(map[string]interface{}))
			rules = append(rules, rulesItemModel)
		}
		topicUpdateSourcesItem.Rules = rules
	}

	return topicUpdateSourcesItem
}

func resourceIBMEnTopicMapToRules(rulesMap map[string]interface{}) en.Rules {
	rules := en.Rules{}

	if rulesMap["enabled"] != nil {
		rules.Enabled = core.BoolPtr(rulesMap["enabled"].(bool))
	}

	rules.EventTypeFilter = core.StringPtr(rulesMap["event_type_filter"].(string))

	if rulesMap["notification_filter"] != nil {
		rules.NotificationFilter = core.StringPtr(rulesMap["notification_filter"].(string))
	}

	if rulesMap["event_schedule_filter"] != nil && len(rulesMap["event_schedule_filter"].([]interface{})) > 0 {
		EventScheduleFilterModel, err := resourceIBMEnTopicsMapToEventScheduleFilterAttributes(rulesMap["event_schedule_filter"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return rules
		}
		rules.EventScheduleFilter = EventScheduleFilterModel
	}

	return rules
}

func resourceIBMEnTopicsMapToEventScheduleFilterAttributes(modelMap map[string]interface{}) (*eventnotificationsv1.EventScheduleFilterAttributes, error) {
	model := &en.EventScheduleFilterAttributes{}
	if modelMap["starts_at"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["starts_at"].(string))
		if err != nil {
			return model, err
		}
		model.StartsAt = &dateTime
	}
	if modelMap["ends_at"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["ends_at"].(string))
		if err != nil {
			return model, err
		}
		model.EndsAt = &dateTime
	}
	if modelMap["expression"] != nil && modelMap["expression"].(string) != "" {
		model.Expression = core.StringPtr(modelMap["expression"].(string))
	}
	return model, nil
}

func enTopicUpdateSourcesItemToMap(source en.SourcesListItems) map[string]interface{} {
	sourceMap := map[string]interface{}{}

	if source.ID != nil {
		sourceMap["id"] = source.ID
	}

	if source.Rules != nil {
		rules := []map[string]interface{}{}
		for _, rulesItem := range source.Rules {
			rulesItemMap := resourceIBMEnTopicRulesToMap(rulesItem)
			rules = append(rules, rulesItemMap)
		}
		sourceMap["rules"] = rules
	}

	return sourceMap
}

func resourceIBMEnTopicRulesToMap(rules en.RulesGet) map[string]interface{} {
	rulesMap := map[string]interface{}{}

	if rules.Enabled != nil {
		rulesMap["enabled"] = rules.Enabled
	}

	rulesMap["event_type_filter"] = rules.EventTypeFilter

	if rules.NotificationFilter != nil {
		rulesMap["notification_filter"] = rules.NotificationFilter
	}

	if rules.EventScheduleFilter != nil {
		eventScheduleFilterMap, err := resourceIBMEnTopicsEventScheduleFilterAttributesToMap(rules.EventScheduleFilter)
		if err != nil {
			return rulesMap
		}
		rulesMap["event_schedule_filter"] = []map[string]interface{}{eventScheduleFilterMap}
	}

	return rulesMap
}

func resourceIBMEnTopicsEventScheduleFilterAttributesToMap(model *eventnotificationsv1.EventScheduleFilterAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.StartsAt != nil {
		modelMap["starts_at"] = model.StartsAt.String()
	}
	if model.EndsAt != nil {
		modelMap["ends_at"] = model.EndsAt.String()
	}
	if model.Expression != nil {
		modelMap["expression"] = model.Expression
	}
	return modelMap, nil
}

func enTopicSubscriptionToMap(subscriptionListItem en.SubscriptionListItem) map[string]interface{} {
	subscriptionMap := map[string]interface{}{}

	subscriptionMap["id"] = subscriptionListItem.ID

	subscriptionMap["name"] = subscriptionListItem.Name

	subscriptionMap["destination_id"] = subscriptionListItem.DestinationID

	if subscriptionListItem.Description != nil {
		subscriptionMap["description"] = subscriptionListItem.Description
	}

	subscriptionMap["destination_type"] = subscriptionListItem.DestinationType

	subscriptionMap["topic_id"] = subscriptionListItem.TopicID

	subscriptionMap["updated_at"] = subscriptionListItem.UpdatedAt.String()

	return subscriptionMap
}

// func enEventScheduleFilterMap(eventschedulerfilter en.EventScheduleFilterAttributes) map[string]interface{} {
// 	schedulerFilter := map[string]interface{}{}

// 	if eventschedulerfilter.StartsAt != nil {
// 		schedulerFilter["enabled"] = eventschedulerfilter.StartsAt
// 	}

// 	if eventschedulerfilter.EndsAt != nil {
// 		schedulerFilter["ends_at"] = eventschedulerfilter.EndsAt
// 	}

// 	if eventschedulerfilter.Expression != nil {
// 		schedulerFilter["expression"] = eventschedulerfilter.Expression
// 	}

// 	return schedulerFilter
// }
