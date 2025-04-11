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

func DataSourceIBMEnPagerDutySubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnPagerDutySubscriptionRead,

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
						"template_id_notification": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The templete id for notification",
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

func dataSourceIBMEnPagerDutySubscriptionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_subscription_pagerduty", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSubscriptionOptions := &en.GetSubscriptionOptions{}

	getSubscriptionOptions.SetInstanceID(d.Get("instance_guid").(string))
	getSubscriptionOptions.SetID(d.Get("subscription_id").(string))

	result, _, err := enClient.GetSubscriptionWithContext(context, getSubscriptionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetIntegrationWithContext failed: %s", err.Error()), "(Data) ibm_en_subscription_pagerduty", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getSubscriptionOptions.InstanceID, *getSubscriptionOptions.ID))

	if err = d.Set("name", result.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_subscription_pagerduty", "read")
		return tfErr.GetDiag()
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_en_subscription_pagerduty", "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("updated_at", result.UpdatedAt); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_en_subscription_pagerduty", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("destination_id", result.DestinationID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting destination_id: %s", err), "(Data) ibm_en_subscription_pagerduty", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("topic_id", result.TopicID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting topic_id: %s", err), "(Data) ibm_en_subscription_pagerduty", "read")
		return tfErr.GetDiag()
	}

	if result.Attributes != nil {
		if err = d.Set("attributes", enPagerDutySubscriptionFlattenAttributes(result.Attributes)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting attributes: %s", err), "(Data) ibm_en_subscription_pagerduty", "read")
			return tfErr.GetDiag()
		}
	}

	return nil
}

func enPagerDutySubscriptionFlattenAttributes(result en.SubscriptionAttributesIntf) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}

	attributes := result.(*en.SubscriptionAttributes)

	finalMap := enPagerDutySubscriptionToMap(attributes)
	finalList = append(finalList, finalMap)

	return finalList
}

func enPagerDutySubscriptionToMap(attributeItem *en.SubscriptionAttributes) (attributeMap map[string]interface{}) {
	attributeMap = map[string]interface{}{}

	if attributeItem.TemplateIDNotification != nil {
		attributeMap["template_id_notification"] = attributeItem.TemplateIDNotification
	}

	return attributeMap
}
