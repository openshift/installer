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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
)

func DataSourceIBMTektonPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMTektonPipelineRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of current instance.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "String.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pipeline status.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID.",
			},
			"toolchain": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Toolchain object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for the toolchain that containing the tekton pipeline.",
						},
					},
				},
			},
			"definitions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Definition list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scm_source": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Scm source for tekton pipeline defintion.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "General href URL.",
									},
									"branch": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A branch of the repo, branch field doesn't coexist with tag field.",
									},
									"tag": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A tag of the repo.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path to the definitions yaml files.",
									},
								},
							},
						},
						"service_instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID.",
						},
					},
				},
			},
			"properties": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tekton pipeline's environment properties.",
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
					},
				},
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Standard RFC 3339 Date Time String.",
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Standard RFC 3339 Date Time String.",
			},
			"pipeline_definition": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tekton pipeline definition document detail object. If this property is absent, the pipeline has no definitions added.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of pipeline definition status.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID.",
						},
					},
				},
			},
			"triggers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tekton pipeline triggers list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID.",
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
				},
			},
			"worker": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Default pipeline worker used to run the pipeline.",
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
			"html_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dashboard URL of this pipeline.",
			},
			"build_number": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The latest pipeline run build number. If this property is absent, the pipeline hasn't had any pipelineRuns.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag whether this pipeline enabled.",
			},
		},
	}
}

func DataSourceIBMTektonPipelineRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineOptions := &cdtektonpipelinev2.GetTektonPipelineOptions{}

	getTektonPipelineOptions.SetID(d.Get("pipeline_id").(string))

	tektonPipeline, response, err := cdTektonPipelineClient.GetTektonPipelineWithContext(context, getTektonPipelineOptions)
	if err != nil {
		log.Printf("[DEBUG] GetTektonPipelineWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getTektonPipelineOptions.ID))

	if err = d.Set("name", tektonPipeline.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("status", tektonPipeline.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	if err = d.Set("resource_group_id", tektonPipeline.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}

	toolchain := []map[string]interface{}{}
	if tektonPipeline.Toolchain != nil {
		modelMap, err := DataSourceIBMTektonPipelineToolchainToMap(tektonPipeline.Toolchain)
		if err != nil {
			return diag.FromErr(err)
		}
		toolchain = append(toolchain, modelMap)
	}
	if err = d.Set("toolchain", toolchain); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain %s", err))
	}

	definitions := []map[string]interface{}{}
	if tektonPipeline.Definitions != nil {
		for _, modelItem := range tektonPipeline.Definitions {
			modelMap, err := DataSourceIBMTektonPipelineDefinitionToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			definitions = append(definitions, modelMap)
		}
	}
	if err = d.Set("definitions", definitions); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting definitions %s", err))
	}

	properties := []map[string]interface{}{}
	if tektonPipeline.Properties != nil {
		for _, modelItem := range tektonPipeline.Properties {
			modelMap, err := DataSourceIBMTektonPipelinePropertyToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			properties = append(properties, modelMap)
		}
	}
	if err = d.Set("properties", properties); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting properties %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(tektonPipeline.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	if err = d.Set("created", flex.DateTimeToString(tektonPipeline.Created)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created: %s", err))
	}

	pipelineDefinition := []map[string]interface{}{}
	if tektonPipeline.PipelineDefinition != nil {
		modelMap, err := DataSourceIBMTektonPipelineTektonPipelinePipelineDefinitionToMap(tektonPipeline.PipelineDefinition)
		if err != nil {
			return diag.FromErr(err)
		}
		pipelineDefinition = append(pipelineDefinition, modelMap)
	}
	if err = d.Set("pipeline_definition", pipelineDefinition); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pipeline_definition %s", err))
	}

	triggers := []map[string]interface{}{}
	if tektonPipeline.Triggers != nil {
		for _, modelItem := range tektonPipeline.Triggers {
			modelMap, err := DataSourceIBMTektonPipelineTriggerToMap(modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			triggers = append(triggers, modelMap)
		}
	}
	if err = d.Set("triggers", triggers); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting triggers %s", err))
	}

	worker := []map[string]interface{}{}
	if tektonPipeline.Worker != nil {
		modelMap, err := DataSourceIBMTektonPipelineWorkerToMap(tektonPipeline.Worker)
		if err != nil {
			return diag.FromErr(err)
		}
		worker = append(worker, modelMap)
	}
	if err = d.Set("worker", worker); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting worker %s", err))
	}

	if err = d.Set("html_url", tektonPipeline.HTMLURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting html_url: %s", err))
	}

	if err = d.Set("build_number", flex.IntValue(tektonPipeline.BuildNumber)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting build_number: %s", err))
	}

	if err = d.Set("enabled", tektonPipeline.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
	}

	return nil
}

func DataSourceIBMTektonPipelineToolchainToMap(model *cdtektonpipelinev2.Toolchain) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineDefinitionToMap(model *cdtektonpipelinev2.Definition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ScmSource != nil {
		scmSourceMap, err := DataSourceIBMTektonPipelineDefinitionScmSourceToMap(model.ScmSource)
		if err != nil {
			return modelMap, err
		}
		modelMap["scm_source"] = []map[string]interface{}{scmSourceMap}
	}
	if model.ServiceInstanceID != nil {
		modelMap["service_instance_id"] = *model.ServiceInstanceID
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineDefinitionScmSourceToMap(model *cdtektonpipelinev2.DefinitionScmSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.Branch != nil {
		modelMap["branch"] = *model.Branch
	}
	if model.Tag != nil {
		modelMap["tag"] = *model.Tag
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelinePropertyToMap(model *cdtektonpipelinev2.Property) (map[string]interface{}, error) {
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
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTektonPipelinePipelineDefinitionToMap(model *cdtektonpipelinev2.TektonPipelinePipelineDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerToMap(model cdtektonpipelinev2.TriggerIntf) (map[string]interface{}, error) {
	if _, ok := model.(*cdtektonpipelinev2.TriggerManualTrigger); ok {
		return DataSourceIBMTektonPipelineTriggerManualTriggerToMap(model.(*cdtektonpipelinev2.TriggerManualTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerScmTrigger); ok {
		return DataSourceIBMTektonPipelineTriggerScmTriggerToMap(model.(*cdtektonpipelinev2.TriggerScmTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerTimerTrigger); ok {
		return DataSourceIBMTektonPipelineTriggerTimerTriggerToMap(model.(*cdtektonpipelinev2.TriggerTimerTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerGenericTrigger); ok {
		return DataSourceIBMTektonPipelineTriggerGenericTriggerToMap(model.(*cdtektonpipelinev2.TriggerGenericTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.Trigger); ok {
		modelMap := make(map[string]interface{})
		model := model.(*cdtektonpipelinev2.Trigger)

		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.EventListener != nil {
			modelMap["event_listener"] = *model.EventListener
		}
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Properties != nil {
			properties := []map[string]interface{}{}
			for _, propertiesItem := range model.Properties {
				propertiesItemMap, err := DataSourceIBMTektonPipelineTriggerPropertiesItemToMap(&propertiesItem)
				if err != nil {
					return modelMap, err
				}
				properties = append(properties, propertiesItemMap)
			}
			modelMap["properties"] = properties
		}
		if model.Tags != nil {
			modelMap["tags"] = model.Tags
		}
		if model.Worker != nil {
			workerMap, err := DataSourceIBMTektonPipelineWorkerToMap(model.Worker)
			if err != nil {
				return modelMap, err
			}
			modelMap["worker"] = []map[string]interface{}{workerMap}
		}
		if model.Concurrency != nil {
			concurrencyMap, err := DataSourceIBMTektonPipelineConcurrencyToMap(model.Concurrency)
			if err != nil {
				return modelMap, err
			}
			modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
		}
		if model.Disabled != nil {
			modelMap["disabled"] = *model.Disabled
		}
		if model.ScmSource != nil {
			scmSourceMap, err := DataSourceIBMTektonPipelineTriggerScmSourceToMap(model.ScmSource)
			if err != nil {
				return modelMap, err
			}
			modelMap["scm_source"] = []map[string]interface{}{scmSourceMap}
		}
		if model.Events != nil {
			eventsMap, err := DataSourceIBMTektonPipelineEventsToMap(model.Events)
			if err != nil {
				return modelMap, err
			}
			modelMap["events"] = []map[string]interface{}{eventsMap}
		}
		if model.ServiceInstanceID != nil {
			modelMap["service_instance_id"] = *model.ServiceInstanceID
		}
		if model.Cron != nil {
			modelMap["cron"] = *model.Cron
		}
		if model.Timezone != nil {
			modelMap["timezone"] = *model.Timezone
		}
		if model.Secret != nil {
			secretMap, err := DataSourceIBMTektonPipelineGenericSecretToMap(model.Secret)
			if err != nil {
				return modelMap, err
			}
			modelMap["secret"] = []map[string]interface{}{secretMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized cdtektonpipelinev2.TriggerIntf subtype encountered")
	}
}

func DataSourceIBMTektonPipelineTriggerPropertiesItemToMap(model *cdtektonpipelinev2.TriggerPropertiesItem) (map[string]interface{}, error) {
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

func DataSourceIBMTektonPipelineWorkerToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
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

func DataSourceIBMTektonPipelineConcurrencyToMap(model *cdtektonpipelinev2.Concurrency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = *model.MaxConcurrentRuns
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerScmSourceToMap(model *cdtektonpipelinev2.TriggerScmSource) (map[string]interface{}, error) {
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

func DataSourceIBMTektonPipelineEventsToMap(model *cdtektonpipelinev2.Events) (map[string]interface{}, error) {
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

func DataSourceIBMTektonPipelineGenericSecretToMap(model *cdtektonpipelinev2.GenericSecret) (map[string]interface{}, error) {
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

func DataSourceIBMTektonPipelineTriggerManualTriggerToMap(model *cdtektonpipelinev2.TriggerManualTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.EventListener != nil {
		modelMap["event_listener"] = *model.EventListener
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := DataSourceIBMTektonPipelineTriggerManualTriggerPropertiesItemToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := DataSourceIBMTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := DataSourceIBMTektonPipelineConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	if model.Disabled != nil {
		modelMap["disabled"] = *model.Disabled
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerManualTriggerPropertiesItemToMap(model *cdtektonpipelinev2.TriggerManualTriggerPropertiesItem) (map[string]interface{}, error) {
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

func DataSourceIBMTektonPipelineTriggerScmTriggerToMap(model *cdtektonpipelinev2.TriggerScmTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.EventListener != nil {
		modelMap["event_listener"] = *model.EventListener
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := DataSourceIBMTektonPipelineTriggerScmTriggerPropertiesItemToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := DataSourceIBMTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := DataSourceIBMTektonPipelineConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	if model.Disabled != nil {
		modelMap["disabled"] = *model.Disabled
	}
	if model.ScmSource != nil {
		scmSourceMap, err := DataSourceIBMTektonPipelineTriggerScmSourceToMap(model.ScmSource)
		if err != nil {
			return modelMap, err
		}
		modelMap["scm_source"] = []map[string]interface{}{scmSourceMap}
	}
	if model.Events != nil {
		eventsMap, err := DataSourceIBMTektonPipelineEventsToMap(model.Events)
		if err != nil {
			return modelMap, err
		}
		modelMap["events"] = []map[string]interface{}{eventsMap}
	}
	if model.ServiceInstanceID != nil {
		modelMap["service_instance_id"] = *model.ServiceInstanceID
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerScmTriggerPropertiesItemToMap(model *cdtektonpipelinev2.TriggerScmTriggerPropertiesItem) (map[string]interface{}, error) {
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

func DataSourceIBMTektonPipelineTriggerTimerTriggerToMap(model *cdtektonpipelinev2.TriggerTimerTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.EventListener != nil {
		modelMap["event_listener"] = *model.EventListener
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := DataSourceIBMTektonPipelineTriggerTimerTriggerPropertiesItemToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := DataSourceIBMTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := DataSourceIBMTektonPipelineConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	if model.Disabled != nil {
		modelMap["disabled"] = *model.Disabled
	}
	if model.Cron != nil {
		modelMap["cron"] = *model.Cron
	}
	if model.Timezone != nil {
		modelMap["timezone"] = *model.Timezone
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerTimerTriggerPropertiesItemToMap(model *cdtektonpipelinev2.TriggerTimerTriggerPropertiesItem) (map[string]interface{}, error) {
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

func DataSourceIBMTektonPipelineTriggerGenericTriggerToMap(model *cdtektonpipelinev2.TriggerGenericTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.EventListener != nil {
		modelMap["event_listener"] = *model.EventListener
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := DataSourceIBMTektonPipelineTriggerGenericTriggerPropertiesItemToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := DataSourceIBMTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := DataSourceIBMTektonPipelineConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	if model.Disabled != nil {
		modelMap["disabled"] = *model.Disabled
	}
	if model.Secret != nil {
		secretMap, err := DataSourceIBMTektonPipelineGenericSecretToMap(model.Secret)
		if err != nil {
			return modelMap, err
		}
		modelMap["secret"] = []map[string]interface{}{secretMap}
	}
	return modelMap, nil
}

func DataSourceIBMTektonPipelineTriggerGenericTriggerPropertiesItemToMap(model *cdtektonpipelinev2.TriggerGenericTriggerPropertiesItem) (map[string]interface{}, error) {
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
