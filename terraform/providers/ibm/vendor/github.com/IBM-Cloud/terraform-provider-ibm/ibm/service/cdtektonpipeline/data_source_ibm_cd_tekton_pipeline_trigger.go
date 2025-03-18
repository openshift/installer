// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package cdtektonpipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtektonpipelinev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func DataSourceIBMCdTektonPipelineTrigger() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdTektonPipelineTriggerRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Tekton pipeline ID.",
			},
			"trigger_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The trigger ID.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trigger type.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trigger name.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.",
			},
			"event_listener": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Event listener name. The name of the event listener to which the trigger is associated. The event listeners are defined in the definition repositories of the Tekton pipeline.",
			},
			"properties": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Optional trigger properties are used to override or supplement the pipeline properties when triggering a pipeline run.",
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
							Description: "Property value. Any string value is valid.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API URL for interacting with the trigger property.",
						},
						"enum": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Options for `single_select` property type. Only needed for `single_select` property type.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Property type.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left blank the full tool integration data will be used.",
						},
						"locked": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "When true, this property cannot be overridden at runtime. Attempting to override it will result in run requests being rejected. The default is false.",
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Optional trigger tags array.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"worker": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details of the worker used to run the trigger.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the worker. Computed based on the worker ID.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the worker. Computed based on the worker ID.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the worker.",
						},
					},
				},
			},
			"max_concurrent_runs": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled for this trigger.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to check if the trigger is enabled.",
			},
			"favorite": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Mark the trigger as a favorite.",
			},
			"enable_events_from_forks": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When enabled, pull request events from forks of the selected repository will trigger a pipeline run.",
			},
			"source": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source repository for a Git trigger. Only required for Git triggers. The referenced repository URL must match the URL of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API https://cloud.ibm.com/apidocs/toolchain#list-tools.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The only supported source type is \"git\", indicating that the source is a git repository.",
						},
						"properties": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Properties of the source, which define the URL of the repository and a branch or pattern.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL of the repository to which the trigger is listening.",
									},
									"branch": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of a branch from the repo. Only one of branch, pattern, or filter should be specified.",
									},
									"pattern": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The pattern of Git branch or tag. You can specify a glob pattern such as '!test' or '*master' to match against multiple tags or branches in the repository.The glob pattern used must conform to Bash 4.3 specifications, see bash documentation for more info: https://www.gnu.org/software/bash/manual/bash.html#Pattern-Matching. Only one of branch, pattern, or filter should be specified.",
									},
									"blind_connection": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "True if the repository server is not addressable on the public internet. IBM Cloud will not be able to validate the connection details you provide.",
									},
									"hook_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Repository webhook ID. It is generated upon trigger creation.",
									},
									"tool": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Reference to the repository tool in the parent toolchain.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of the repository tool instance in the parent toolchain.",
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
			"events": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Either 'events' or 'filter' is required specifically for Git triggers. Stores a list of events that a Git trigger listens to. Choose one or more from 'push', 'pull_request', and 'pull_request_closed'. If SCM repositories use the 'merge request' term, they correspond to the generic term i.e. 'pull request'.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"filter": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Either 'events' or 'filter' can be used. Stores the CEL (Common Expression Language) expression value which is used for event filtering against the Git webhook payloads.",
			},
			"cron": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only needed for timer triggers. CRON expression that indicates when this trigger will activate. Maximum frequency is every 5 minutes. The string is based on UNIX crontab syntax: minute, hour, day of month, month, day of week. Example: The CRON expression 0 *_/2 * * * - translates to - every 2 hours.",
			},
			"timezone": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only used for timer triggers. Specify the timezone used for this timer trigger, which will ensure the CRON activates this trigger relative to the specified timezone. If no timezone is specified, the default timezone used is UTC. Valid timezones are those listed in the IANA timezone database, https://www.iana.org/time-zones.",
			},
			"secret": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Only needed for Generic Webhook trigger type. The secret is used to start the Generic Webhook trigger.",
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
							Description: "Secret value, not needed if secret type is `internal_validation`.",
						},
						"source": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secret location, not needed if secret type is `internal_validation`.",
						},
						"key_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secret name, not needed if type is `internal_validation`.",
						},
						"algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm used for `digest_matches` secret type. Only needed for `digest_matches` secret type.",
						},
					},
				},
			},
			"webhook_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Webhook URL that can be used to trigger pipeline runs.",
			},
		},
	}
}

func dataSourceIBMCdTektonPipelineTriggerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getTektonPipelineTriggerOptions := &cdtektonpipelinev2.GetTektonPipelineTriggerOptions{}

	getTektonPipelineTriggerOptions.SetPipelineID(d.Get("pipeline_id").(string))
	getTektonPipelineTriggerOptions.SetTriggerID(d.Get("trigger_id").(string))

	triggerIntf, _, err := cdTektonPipelineClient.GetTektonPipelineTriggerWithContext(context, getTektonPipelineTriggerOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTektonPipelineTriggerWithContext failed: %s", err.Error()), "(Data) ibm_cd_tekton_pipeline_trigger", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	trigger := triggerIntf.(*cdtektonpipelinev2.Trigger)

	d.SetId(fmt.Sprintf("%s/%s", *getTektonPipelineTriggerOptions.PipelineID, *getTektonPipelineTriggerOptions.TriggerID))

	if !core.IsNil(trigger.Type) {
		if err = d.Set("type", trigger.Type); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-type").GetDiag()
		}
	}

	if !core.IsNil(trigger.Name) {
		if err = d.Set("name", trigger.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-name").GetDiag()
		}
	}

	if !core.IsNil(trigger.Href) {
		if err = d.Set("href", trigger.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-href").GetDiag()
		}
	}

	if !core.IsNil(trigger.EventListener) {
		if err = d.Set("event_listener", trigger.EventListener); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting event_listener: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-event_listener").GetDiag()
		}
	}

	if !core.IsNil(trigger.Properties) {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range trigger.Properties {
			propertiesItemMap, err := DataSourceIBMCdTektonPipelineTriggerTriggerPropertyToMap(&propertiesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "properties-to-map").GetDiag()
			}
			properties = append(properties, propertiesItemMap)
		}
		if err = d.Set("properties", properties); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting properties: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-properties").GetDiag()
		}
	}

	if !core.IsNil(trigger.Tags) {
		tags := []interface{}{}
		for _, tagsItem := range trigger.Tags {
			tags = append(tags, tagsItem)
		}
		if err = d.Set("tags", tags); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tags: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-tags").GetDiag()
		}
	}

	if !core.IsNil(trigger.Worker) {
		worker := []map[string]interface{}{}
		workerMap, err := DataSourceIBMCdTektonPipelineTriggerWorkerToMap(trigger.Worker)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "worker-to-map").GetDiag()
		}
		worker = append(worker, workerMap)
		if err = d.Set("worker", worker); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting worker: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-worker").GetDiag()
		}
	}

	if !core.IsNil(trigger.MaxConcurrentRuns) {
		if err = d.Set("max_concurrent_runs", flex.IntValue(trigger.MaxConcurrentRuns)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting max_concurrent_runs: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-max_concurrent_runs").GetDiag()
		}
	}

	if !core.IsNil(trigger.Enabled) {
		if err = d.Set("enabled", trigger.Enabled); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting enabled: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-enabled").GetDiag()
		}
	}

	if !core.IsNil(trigger.Favorite) {
		if err = d.Set("favorite", trigger.Favorite); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting favorite: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-favorite").GetDiag()
		}
	}

	if !core.IsNil(trigger.EnableEventsFromForks) {
		if err = d.Set("enable_events_from_forks", trigger.EnableEventsFromForks); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting enable_events_from_forks: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-enable_events_from_forks").GetDiag()
		}
	}

	if !core.IsNil(trigger.Source) {
		source := []map[string]interface{}{}
		sourceMap, err := DataSourceIBMCdTektonPipelineTriggerTriggerSourceToMap(trigger.Source)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "source-to-map").GetDiag()
		}
		source = append(source, sourceMap)
		if err = d.Set("source", source); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting source: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-source").GetDiag()
		}
	}

	if !core.IsNil(trigger.Events) {
		events := []interface{}{}
		for _, eventsItem := range trigger.Events {
			events = append(events, eventsItem)
		}
		if err = d.Set("events", events); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting events: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-events").GetDiag()
		}
	}

	if !core.IsNil(trigger.Filter) {
		if err = d.Set("filter", trigger.Filter); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting filter: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-filter").GetDiag()
		}
	}

	if !core.IsNil(trigger.Cron) {
		if err = d.Set("cron", trigger.Cron); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cron: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-cron").GetDiag()
		}
	}

	if !core.IsNil(trigger.Timezone) {
		if err = d.Set("timezone", trigger.Timezone); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting timezone: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-timezone").GetDiag()
		}
	}

	if !core.IsNil(trigger.Secret) {
		secret := []map[string]interface{}{}
		secretMap, err := DataSourceIBMCdTektonPipelineTriggerGenericSecretToMap(trigger.Secret)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "secret-to-map").GetDiag()
		}
		secret = append(secret, secretMap)
		if err = d.Set("secret", secret); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting secret: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-secret").GetDiag()
		}
	}

	if !core.IsNil(trigger.WebhookURL) {
		if err = d.Set("webhook_url", trigger.WebhookURL); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting webhook_url: %s", err), "(Data) ibm_cd_tekton_pipeline_trigger", "read", "set-webhook_url").GetDiag()
		}
	}

	return nil
}

func DataSourceIBMCdTektonPipelineTriggerTriggerPropertyToMap(model *cdtektonpipelinev2.TriggerProperty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.Enum != nil {
		modelMap["enum"] = model.Enum
	}
	modelMap["type"] = *model.Type
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.Locked != nil {
		modelMap["locked"] = *model.Locked
	}
	return modelMap, nil
}

func DataSourceIBMCdTektonPipelineTriggerWorkerToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func DataSourceIBMCdTektonPipelineTriggerTriggerSourceToMap(model *cdtektonpipelinev2.TriggerSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	propertiesMap, err := DataSourceIBMCdTektonPipelineTriggerTriggerSourcePropertiesToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func DataSourceIBMCdTektonPipelineTriggerTriggerSourcePropertiesToMap(model *cdtektonpipelinev2.TriggerSourceProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["url"] = *model.URL
	if model.Branch != nil {
		modelMap["branch"] = *model.Branch
	}
	if model.Pattern != nil {
		modelMap["pattern"] = *model.Pattern
	}
	modelMap["blind_connection"] = *model.BlindConnection
	if model.HookID != nil {
		modelMap["hook_id"] = *model.HookID
	}
	toolMap, err := DataSourceIBMCdTektonPipelineTriggerToolToMap(model.Tool)
	if err != nil {
		return modelMap, err
	}
	modelMap["tool"] = []map[string]interface{}{toolMap}
	return modelMap, nil
}

func DataSourceIBMCdTektonPipelineTriggerToolToMap(model *cdtektonpipelinev2.Tool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func DataSourceIBMCdTektonPipelineTriggerGenericSecretToMap(model *cdtektonpipelinev2.GenericSecret) (map[string]interface{}, error) {
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
