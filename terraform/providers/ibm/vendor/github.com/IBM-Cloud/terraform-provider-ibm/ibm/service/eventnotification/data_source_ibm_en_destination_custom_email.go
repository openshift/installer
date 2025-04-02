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

func DataSourceIBMEnCustomEmailDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMEnCustomEmailDestinationRead,

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
				Description: "Destination type slack.",
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
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The custom doamin",
									},
									"dkim": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The DKIM attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"public_key": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "dkim public key.",
												},
												"selector": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "dkim selector.",
												},
												"verification": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "dkim verification.",
												},
											},
										},
									},
									"spf": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The SPF attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"txt_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "spf text name.",
												},
												"txt_value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "spf text value.",
												},
												"verification": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "spf verification.",
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

func dataSourceIBMEnCustomEmailDestinationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enClient, err := meta.(conns.ClientSession).EventNotificationsApiV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_en_destination_custom_email", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &en.GetDestinationOptions{}

	options.SetInstanceID(d.Get("instance_guid").(string))
	options.SetID(d.Get("destination_id").(string))

	result, _, err := enClient.GetDestinationWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDestinationWithContext failed: %s", err.Error()), "(Data) ibm_en_destination_custom_email", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *options.InstanceID, *options.ID))

	if err = d.Set("name", result.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_en_destination_custom_email", "read")
		return tfErr.GetDiag()
	}

	if result.Description != nil {
		if err = d.Set("description", result.Description); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_en_destination_custom_email", "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("type", result.Type); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_en_destination_custom_email", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("collect_failed_events", result.CollectFailedEvents); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting CollectFailedEvents: %s", err), "(Data) ibm_en_destination_custom_email", "read")
		return tfErr.GetDiag()
	}

	if result.Config != nil {
		err = d.Set("config", enCustomEmailDestinationFlattenConfig(*result.Config))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config: %s", err), "(Data) ibm_en_destination_custom_email", "read")
			return tfErr.GetDiag()
		}
	}

	if result.SubscriptionNames != nil {
		err = d.Set("subscription_names", result.SubscriptionNames)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting subscription_names: %s", err), "(Data) ibm_en_destination_custom_email", "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(result.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_en_destination_custom_email", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("subscription_count", flex.IntValue(result.SubscriptionCount)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting subscription_count: %s", err), "(Data) ibm_en_destination_custom_email", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func enCustomEmailDestinationFlattenConfig(result en.DestinationConfig) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := enCustomEmailDestinationConfigToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func enCustomEmailDestinationConfigToMap(configItem en.DestinationConfig) (configMap map[string]interface{}) {
	configMap = map[string]interface{}{}

	if configItem.Params != nil {
		paramsList := []map[string]interface{}{}
		paramsMap := enCustomEmailDestinationConfigParamsToMap(configItem.Params)
		paramsList = append(paramsList, paramsMap)
		configMap["params"] = paramsList
	}

	return configMap
}

func enCustomEmailDestinationConfigParamsToMap(paramsItem en.DestinationConfigOneOfIntf) (paramsMap map[string]interface{}) {
	paramsMap = map[string]interface{}{}

	params := paramsItem.(*en.DestinationConfigOneOf)

	if params.Domain != nil {
		paramsMap["domain"] = params.Domain
	}

	if params.Dkim != nil {
		dkimMap, err := dataSourceIBMEnDestinationDkimAttributesToMap(params.Dkim)
		if err != nil {
			return paramsMap
		}
		paramsMap["dkim"] = []map[string]interface{}{dkimMap}
	}
	if params.Spf != nil {
		spfMap, err := dataSourceIBMEnDestinationSpfAttributesToMap(params.Spf)
		if err != nil {
			return paramsMap
		}
		paramsMap["spf"] = []map[string]interface{}{spfMap}
	}

	return paramsMap
}

func dataSourceIBMEnDestinationDkimAttributesToMap(model *en.DkimAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PublicKey != nil {
		modelMap["public_key"] = model.PublicKey
	}
	if model.Selector != nil {
		modelMap["selector"] = model.Selector
	}
	if model.Verification != nil {
		modelMap["verification"] = model.Verification
	}
	return modelMap, nil
}

func dataSourceIBMEnDestinationSpfAttributesToMap(model *en.SpfAttributes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TxtName != nil {
		modelMap["txt_name"] = model.TxtName
	}
	if model.TxtValue != nil {
		modelMap["txt_value"] = model.TxtValue
	}
	if model.Verification != nil {
		modelMap["verification"] = model.Verification
	}
	return modelMap, nil
}
