// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtektonpipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
)

func DataSourceIBMTektonPipelineTrigger() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMTektonPipelineTriggerRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tekton pipeline ID.",
			},
			"trigger_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The trigger ID.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the duplicated trigger.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trigger type.",
			},
			"event_listener": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Event listener name.",
			},
			"properties": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Trigger properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Property name.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "String format property value.",
						},
						"enum": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Options for SINGLE_SELECT property type.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"default": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default option for SINGLE_SELECT property type.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Property type.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "property path for INTEGRATION type properties.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "General href URL.",
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Trigger tags array.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"worker": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this trigger uses default pipeline worker.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "worker name.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "worker type.",
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"concurrency": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Concurrency object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_concurrent_runs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Defines the maximum number of concurrent runs for this trigger.",
						},
					},
				},
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "flag whether the trigger is disabled.",
			},
			"scm_source": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Scm source for git type tekton pipeline trigger.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Needed only for git trigger type. Repo URL that listening to.",
						},
						"branch": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Needed only for git trigger type. Branch name of the repo. Branch field doesn't coexist with pattern field.",
						},
						"pattern": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Needed only for git trigger type. Git branch or tag pattern to listen to. Please refer to https://github.com/micromatch/micromatch for pattern syntax.",
						},
						"blind_connection": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Needed only for git trigger type. Branch name of the repo.",
						},
						"hook_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Webhook ID.",
						},
					},
				},
			},
			"events": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Needed only for git trigger type. Events object defines the events this git trigger listening to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"push": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If true, the trigger starts when tekton pipeline service receive a repo's 'push' git webhook event.",
						},
						"pull_request_closed": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If true, the trigger starts when tekton pipeline service receive a repo pull reqeust's 'close' git webhook event.",
						},
						"pull_request": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If true, the trigger starts when tekton pipeline service receive a repo pull reqeust's 'open' or 'update' git webhook event.",
						},
					},
				},
			},
			"service_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID.",
			},
			"cron": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Needed only for timer trigger type. Cron expression for timer trigger. Maximum frequency is every 5 minutes.",
			},
			"timezone": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Needed only for timer trigger type. Timezones for timer trigger.",
			},
			"secret": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Needed only for generic trigger type. Secret used to start generic trigger.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secret type.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secret value, not needed if secret type is \"internalValidation\".",
						},
						"source": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secret location, not needed if secret type is \"internalValidation\".",
						},
						"key_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secret name, not needed if type is \"internalValidation\".",
						},
						"algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm used for \"digestMatches\" secret type.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMTektonPipelineTriggerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineTriggerOptions := &cdtektonpipelinev2.GetTektonPipelineTriggerOptions{}

	getTektonPipelineTriggerOptions.SetPipelineID(d.Get("pipeline_id").(string))
	getTektonPipelineTriggerOptions.SetTriggerID(d.Get("trigger_id").(string))

	TriggerIntf, response, err := cdTektonPipelineClient.GetTektonPipelineTriggerWithContext(context, getTektonPipelineTriggerOptions)
	if err != nil {
		log.Printf("[DEBUG] GetTektonPipelineTriggerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineTriggerWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getTektonPipelineTriggerOptions.PipelineID, *getTektonPipelineTriggerOptions.TriggerID))

	trigger := TriggerIntf.(*cdtektonpipelinev2.Trigger)

	if err = d.Set("type", trigger.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	if err = d.Set("name", trigger.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("event_listener", trigger.EventListener); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting event_listener: %s", err))
	}

	properties := []map[string]interface{}{}
	if trigger.Properties != nil {
		for _, modelItem := range trigger.Properties {
			modelMap, err := DataSourceIBMTektonPipelineTriggerTriggerPropertiesItemToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			properties = append(properties, modelMap)
		}
	}
	if err = d.Set("properties", properties); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting properties %s", err))
	}

	worker := []map[string]interface{}{}
	if trigger.Worker != nil {
		modelMap, err := DataSourceIBMTektonPipelineTriggerWorkerToMap(trigger.Worker)
		if err != nil {
			return diag.FromErr(err)
		}
		worker = append(worker, modelMap)
	}
	if err = d.Set("worker", worker); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting worker %s", err))
	}

	concurrency := []map[string]interface{}{}
	if trigger.Concurrency != nil {
		modelMap, err := DataSourceIBMTektonPipelineTriggerConcurrencyToMap(trigger.Concurrency)
		if err != nil {
			return diag.FromErr(err)
		}
		concurrency = append(concurrency, modelMap)
	}
	if err = d.Set("concurrency", concurrency); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting concurrency %s", err))
	}

	if err = d.Set("disabled", trigger.Disabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting disabled: %s", err))
	}

	scmSource := []map[string]interface{}{}
	if trigger.ScmSource != nil {
		modelMap, err := DataSourceIBMTektonPipelineTriggerTriggerScmSourceToMap(trigger.ScmSource)
		if err != nil {
			return diag.FromErr(err)
		}
		scmSource = append(scmSource, modelMap)
	}
	if err = d.Set("scm_source", scmSource); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting scm_source %s", err))
	}

	events := []map[string]interface{}{}
	if trigger.Events != nil {
		modelMap, err := DataSourceIBMTektonPipelineTriggerEventsToMap(trigger.Events)
		if err != nil {
			return diag.FromErr(err)
		}
		events = append(events, modelMap)
	}
	if err = d.Set("events", events); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting events %s", err))
	}

	if err = d.Set("service_instance_id", trigger.ServiceInstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_instance_id: %s", err))
	}

	if err = d.Set("cron", trigger.Cron); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cron: %s", err))
	}

	if err = d.Set("timezone", trigger.Timezone); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting timezone: %s", err))
	}

	secret := []map[string]interface{}{}
	if trigger.Secret != nil {
		modelMap, err := DataSourceIBMTektonPipelineTriggerGenericSecretToMap(trigger.Secret)
		if err != nil {
			return diag.FromErr(err)
		}
		secret = append(secret, modelMap)
	}
	if err = d.Set("secret", secret); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret %s", err))
	}

	return nil
}

func DataSourceIBMTektonPipelineTriggerTriggerPropertiesItemToMap(model *cdtektonpipelinev2.TriggerPropertiesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.Enum != nil {
		modelMap["enum"] = model.Enum
	}
	if model.Default != nil {
		modelMap["default"] = *model.Default
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerWorkerToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerConcurrencyToMap(model *cdtektonpipelinev2.Concurrency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = *model.MaxConcurrentRuns
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerTriggerScmSourceToMap(model *cdtektonpipelinev2.TriggerScmSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.Branch != nil {
		modelMap["branch"] = *model.Branch
	}
	if model.Pattern != nil {
		modelMap["pattern"] = *model.Pattern
	}
	if model.BlindConnection != nil {
		modelMap["blind_connection"] = *model.BlindConnection
	}
	if model.HookID != nil {
		modelMap["hook_id"] = *model.HookID
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerEventsToMap(model *cdtektonpipelinev2.Events) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Push != nil {
		modelMap["push"] = *model.Push
	}
	if model.PullRequestClosed != nil {
		modelMap["pull_request_closed"] = *model.PullRequestClosed
	}
	if model.PullRequest != nil {
		modelMap["pull_request"] = *model.PullRequest
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerGenericSecretToMap(model *cdtektonpipelinev2.GenericSecret) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.Source != nil {
		modelMap["source"] = *model.Source
	}
	if model.KeyName != nil {
		modelMap["key_name"] = *model.KeyName
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = *model.Algorithm
	}
	return modelMap, nil
}
