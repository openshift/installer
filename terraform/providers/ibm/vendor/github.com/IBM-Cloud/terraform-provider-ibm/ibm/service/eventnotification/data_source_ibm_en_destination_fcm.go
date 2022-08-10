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
)

func DataSourceIBMEnFCMDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnFCMDestinationRead,

		Schema: map[string]*schema.Schema{
			"instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for IBM Cloud Event Notifications instance.",
			},
			"destination_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier for Destination.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination description.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Destination type push_android.",
			},
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Payload describing a destination configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sender_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Sender_id value for FCM project.",
									},
									"server_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Server_key value for FCM project.",
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
			"subscription_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of subscriptions.",
			},
			"subscription_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of subscriptions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIBMEnFCMDestinationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		return diag.FromErr(err)
	}

	options := &en.GetDestinationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("destination_id").(string))

	result, response, err := enClient.GetDestinationWithContext(context, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("GetDestination failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *options.ID))

	if err = d.Set("name", result.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
		}
	}

	if err = d.Set("type", result.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}

	if result.Config != nil {
		err = d.Set("config", enFCMDestinationFlattenConfig(*result.Config))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting config %s", err))
		}
	}

	if result.SubscriptionNames != nil {
		err = d.Set("subscription_names", result.SubscriptionNames)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_names %s", err))
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	if err = d.Set("subscription_count", flex.IntValue(result.SubscriptionCount)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subscription_count: %s", err))
	}

	return nil
}

func enFCMDestinationFlattenConfig(result en.DestinationConfig) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := enFCMDestinationConfigToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func enFCMDestinationConfigToMap(configItem en.DestinationConfig) (configMap map[string]interface{}) {
	configMap = map[string]interface{}{}

	if configItem.Params != nil {
		paramsList := []map[string]interface{}{}
		paramsMap := enFCMDestinationConfigParamsToMap(configItem.Params)
		paramsList = append(paramsList, paramsMap)
		configMap["params"] = paramsList
	}

	return configMap
}

func enFCMDestinationConfigParamsToMap(paramsItem en.DestinationConfigParamsIntf) (paramsMap map[string]interface{}) {
	paramsMap = map[string]interface{}{}

	params := paramsItem.(*en.DestinationConfigParams)

	if params.SenderID != nil {
		paramsMap["sender_id"] = params.SenderID
	}
	if params.ServerKey != nil {
		paramsMap["server_key"] = params.ServerKey
	}

	return paramsMap
}
