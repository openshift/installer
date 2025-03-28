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

func DataSourceIBMEnCodeEngineDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnCodeEngineDestinationRead,

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
				Description: "Destination type Webhook.",
			},
			"collect_failed_events": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to collect the failed event in Cloud Object Storage bucket",
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
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL of webhook.",
									},
									"verb": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP method of webhook.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The code engine destination type.",
									},
									"project_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CRN of the code engine project.",
									},
									"job_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the code engine job.",
									},
									"custom_headers": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "Custom headers (Key-Value pair) for webhook call.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"sensitive_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of sensitive headers from custom headers.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
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

func dataSourceIBMEnCodeEngineDestinationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_destination_ce", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.GetDestinationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("destination_id").(string))

	result, _, err := enClient.GetDestinationWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDestinationWithContext failed: %s", err.Error()), "(Data) ibm_en_destination_ce", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *options.ID))

	if err = d.Set("name", result.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_destination_ce", "read")
		return tfErr.GetDiag()
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_en_destination_ce", "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("type", result.Type); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_en_destination_ce", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("collect_failed_events", result.CollectFailedEvents); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting CollectFailedEvents: %s", err), "(Data) ibm_en_destination_ce", "read")
		return tfErr.GetDiag()
	}

	if result.Config != nil {
		err = d.Set("config", enCodeEngineDestinationFlattenConfig(*result.Config))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config: %s", err), "(Data) ibm_en_destination_ce", "read")
			return tfErr.GetDiag()
		}
	}

	if result.SubscriptionNames != nil {
		err = d.Set("subscription_names", result.SubscriptionNames)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting subscription_names: %s", err), "(Data) ibm_en_destination_ce", "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_en_destination_ce", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("subscription_count", flex.IntValue(result.SubscriptionCount)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting subscription_count: %s", err), "(Data) ibm_en_destination_ce", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func enCodeEngineDestinationFlattenConfig(result en.DestinationConfig) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := enCodeEngineDestinationConfigToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func enCodeEngineDestinationConfigToMap(configItem en.DestinationConfig) (configMap map[string]interface{}) {
	configMap = map[string]interface{}{}

	if configItem.Params != nil {
		paramsList := []map[string]interface{}{}
		paramsMap := enCodeEngineDestinationConfigParamsToMap(configItem.Params)
		paramsList = append(paramsList, paramsMap)
		configMap["params"] = paramsList
	}

	return configMap
}

func enCodeEngineDestinationConfigParamsToMap(paramsItem en.DestinationConfigOneOfIntf) (paramsMap map[string]interface{}) {
	paramsMap = map[string]interface{}{}

	params := paramsItem.(*en.DestinationConfigOneOf)

	paramsMap["type"] = params.Type

	if *params.Type == "application" {

		if params.URL != nil {
			paramsMap["url"] = params.URL
		}
		if params.Verb != nil {
			paramsMap["verb"] = params.Verb
		}

		if params.CustomHeaders != nil {
			paramsMap["custom_headers"] = params.CustomHeaders
		}
		if params.SensitiveHeaders != nil {
			paramsMap["sensitive_headers"] = params.SensitiveHeaders
		}

		return paramsMap

	} else {

		if params.ProjectCRN != nil {
			paramsMap["project_crn"] = params.ProjectCRN
		}
		if params.JobName != nil {
			paramsMap["job_name"] = params.JobName
		}

		return paramsMap

	}

}
