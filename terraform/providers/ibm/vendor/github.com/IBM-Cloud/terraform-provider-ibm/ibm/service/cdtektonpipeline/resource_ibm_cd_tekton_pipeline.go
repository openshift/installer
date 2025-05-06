// Copyright IBM Corp. 2023 All Rights Reserved.
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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMCdTektonPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdTektonPipelineCreate,
		ReadContext:   resourceIBMCdTektonPipelineRead,
		UpdateContext: resourceIBMCdTektonPipelineUpdate,
		DeleteContext: resourceIBMCdTektonPipelineDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"next_build_number": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specify the build number that will be used for the next pipeline run. Build numbers can be any positive whole number between 0 and 100000000000000.",
			},
			"enable_notifications": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag whether to enable notifications for this pipeline. When enabled, pipeline run events are published on all slack integration specified channels in the parent toolchain.",
			},
			"enable_partial_cloning": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag whether to enable partial cloning for this pipeline. When partial clone is enabled, only the files contained within the paths specified in definition repositories are read and cloned, this means that symbolic links might not work.",
			},
			"worker": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Worker object containing worker ID only. If omitted the IBM Managed shared workers are used by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the worker.",
						},
					},
				},
			},
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "String.",
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
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group in which the pipeline was created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID.",
						},
					},
				},
			},
			"toolchain": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Toolchain object containing references to the parent toolchain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "UUID.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN for the toolchain that contains the Tekton pipeline.",
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
						"source": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Source repository containing the Tekton pipeline definition.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The only supported source type is \"git\", indicating that the source is a git repository.",
									},
									"properties": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Properties of the source, which define the URL of the repository and a branch or tag.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "URL of the definition repository.",
												},
												"branch": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "A branch from the repo, specify one of branch or tag only.",
												},
												"tag": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "A tag from the repo, specify one of branch or tag only.",
												},
												"path": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The path to the definition's YAML files.",
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
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API URL for interacting with the definition.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The aggregated definition ID.",
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
							Required:    true,
							ForceNew:    true,
							Description: "Property name.",
						},
						"value": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: flex.SuppressPipelinePropertyRawSecret,
							Description:      "Property value. Any string value is valid.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API URL for interacting with the property.",
						},
						"enum": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Options for `single_select` property type. Only needed when using `single_select` property type.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Property type.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left blank the full tool integration data will be used.",
						},
					},
				},
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Standard RFC 3339 Date Time String.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Standard RFC 3339 Date Time String.",
			},
			"triggers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tekton pipeline triggers list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger type.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger name.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.",
						},
						"event_listener": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event listener name. The name of the event listener to which the trigger is associated. The event listeners are defined in the definition repositories of the Tekton pipeline.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Trigger ID.",
						},
						"properties": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Optional trigger properties used to override or supplement the pipeline properties when triggering a pipeline run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Property name.",
									},
									"value": &schema.Schema{
										Type:             schema.TypeString,
										Optional:         true,
										DiffSuppressFunc: flex.SuppressTriggerPropertyRawSecret,
										Description:      "Property value. Any string value is valid.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API URL for interacting with the trigger property.",
									},
									"enum": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Options for `single_select` property type. Only needed for `single_select` property type.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Property type.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left blank the full tool integration data will be used.",
									},
								},
							},
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Optional trigger tags array.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"worker": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Worker used to run the trigger. If not specified the trigger will use the default pipeline worker.",
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
										Required:    true,
										Description: "ID of the worker.",
									},
								},
							},
						},
						"max_concurrent_runs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled for this trigger.",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Flag whether the trigger is enabled.",
						},
						"source": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Source repository for a Git trigger. Only required for Git triggers. The referenced repository URL must match the URL of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API https://cloud.ibm.com/apidocs/toolchain#list-tools.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The only supported source type is \"git\", indicating that the source is a git repository.",
									},
									"properties": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Properties of the source, which define the URL of the repository and a branch or pattern.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "URL of the repository to which the trigger is listening.",
												},
												"branch": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of a branch from the repo. One of branch or pattern must be specified, but only one or the other.",
												},
												"pattern": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The pattern of Git branch or tag to which to listen. You can specify a glob pattern such as '!test' or '*master' to match against multiple tags/branches in the repository. The glob pattern used must conform to Bash 4.3 specifications, see bash documentation for more info: https://www.gnu.org/software/bash/manual/bash.html#Pattern-Matching. One of branch or pattern must be specified, but only one or the other.",
												},
												"blind_connection": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "True if the repository server is not addressable on the public internet. IBM Cloud will not be able to validate the connection details you provide.",
												},
												"hook_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of the webhook from the repo. Computed upon creation of the trigger.",
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
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: flex.SuppressTriggerEvents,
							Description:      "Only needed for Git triggers. List of events to which a Git trigger listens. Choose one or more from: 'push', 'pull_request' and 'pull_request_closed'. For SCM repositories that use 'merge request' events, such events map to the equivalent 'pull request' events.",
							Elem:             &schema.Schema{Type: schema.TypeString},
						},
						"cron": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Only needed for timer triggers. Cron expression that indicates when this trigger will activate. Maximum frequency is every 5 minutes. The string is based on UNIX crontab syntax: minute, hour, day of month, month, day of week. Example: 0 *_/2 * * * - every 2 hours.",
						},
						"timezone": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Only used for timer triggers. Specify the timezone used for this timer trigger, which will ensure the cron activates this trigger relative to the specified timezone. If no timezone is specified, the default timezone used is UTC. Valid timezones are those listed in the IANA timezone database, https://www.iana.org/time-zones.",
						},
						"secret": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Only needed for generic webhook trigger type. Secret used to start generic webhook trigger.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret type.",
									},
									"value": &schema.Schema{
										Type:             schema.TypeString,
										Optional:         true,
										DiffSuppressFunc: flex.SuppressGenericWebhookRawSecret,
										Description:      "Secret value, not needed if secret type is `internal_validation`.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret location, not needed if secret type is `internal_validation`.",
									},
									"key_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret name, not needed if type is `internal_validation`.",
									},
									"algorithm": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
				},
			},
			"runs_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for this pipeline showing the list of pipeline runs.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API URL for interacting with the pipeline.",
			},
			"build_number": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The latest pipeline run build number. If this property is absent, the pipeline hasn't had any pipeline runs.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag whether this pipeline is enabled.",
			},
		},
	}
}

func ResourceIBMCdTektonPipelineValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "next_build_number",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "99999999999999",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_tekton_pipeline", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdTektonPipelineCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelineOptions := &cdtektonpipelinev2.CreateTektonPipelineOptions{}

	if _, ok := d.GetOk("next_build_number"); ok {
		createTektonPipelineOptions.SetNextBuildNumber(int64(d.Get("next_build_number").(int)))
	}
	if _, ok := d.GetOkExists("enable_notifications"); ok {
		createTektonPipelineOptions.SetEnableNotifications(d.Get("enable_notifications").(bool))
	}
	if _, ok := d.GetOkExists("enable_partial_cloning"); ok {
		createTektonPipelineOptions.SetEnablePartialCloning(d.Get("enable_partial_cloning").(bool))
	}
	if _, ok := d.GetOk("worker"); ok {
		workerModel, err := resourceIBMCdTektonPipelineMapToWorkerIdentity(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createTektonPipelineOptions.SetWorker(workerModel)
	}

	if _, ok := d.GetOk("pipeline_id"); ok {
		createTektonPipelineOptions.SetID(d.Get("pipeline_id").(string))
	}
	tektonPipeline, response, err := cdTektonPipelineClient.CreateTektonPipelineWithContext(context, createTektonPipelineOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTektonPipelineWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTektonPipelineWithContext failed %s\n%s", err, response))
	}

	d.SetId(*tektonPipeline.ID)

	return resourceIBMCdTektonPipelineRead(context, d, meta)
}

func resourceIBMCdTektonPipelineRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineOptions := &cdtektonpipelinev2.GetTektonPipelineOptions{}

	getTektonPipelineOptions.SetID(d.Id())

	tektonPipeline, response, err := cdTektonPipelineClient.GetTektonPipelineWithContext(context, getTektonPipelineOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTektonPipelineWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineWithContext failed %s\n%s", err, response))
	}

	if !core.IsNil(tektonPipeline.NextBuildNumber) {
		if err = d.Set("next_build_number", flex.IntValue(tektonPipeline.NextBuildNumber)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting next_build_number: %s", err))
		}
	}
	if !core.IsNil(tektonPipeline.EnableNotifications) {
		if err = d.Set("enable_notifications", tektonPipeline.EnableNotifications); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enable_notifications: %s", err))
		}
	}
	if !core.IsNil(tektonPipeline.EnablePartialCloning) {
		if err = d.Set("enable_partial_cloning", tektonPipeline.EnablePartialCloning); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enable_partial_cloning: %s", err))
		}
	}
	if !core.IsNil(tektonPipeline.Worker) {
		workerMap, err := resourceIBMCdTektonPipelineWorkerIdentityToMap(tektonPipeline.Worker)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("worker", []map[string]interface{}{workerMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting worker: %s", err))
		}
	}
	if err = d.Set("pipeline_id", tektonPipeline.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pipeline_id: %s", err))
	}
	if err = d.Set("name", tektonPipeline.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("status", tektonPipeline.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}
	resourceGroupMap, err := resourceIBMCdTektonPipelineResourceGroupReferenceToMap(tektonPipeline.ResourceGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("resource_group", []map[string]interface{}{resourceGroupMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
	}
	toolchainMap, err := resourceIBMCdTektonPipelineToolchainReferenceToMap(tektonPipeline.Toolchain)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("toolchain", []map[string]interface{}{toolchainMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain: %s", err))
	}
	definitions := []map[string]interface{}{}
	for _, definitionsItem := range tektonPipeline.Definitions {
		definitionsItemMap, err := resourceIBMCdTektonPipelineDefinitionToMap(&definitionsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		definitions = append(definitions, definitionsItemMap)
	}
	if err = d.Set("definitions", definitions); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting definitions: %s", err))
	}
	properties := []map[string]interface{}{}
	for _, propertiesItem := range tektonPipeline.Properties {
		propertiesItemMap, err := resourceIBMCdTektonPipelinePropertyToMap(&propertiesItem)
		if err != nil {
			return diag.FromErr(err)
		}
		properties = append(properties, propertiesItemMap)
	}
	if err = d.Set("properties", properties); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting properties: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(tektonPipeline.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(tektonPipeline.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	triggers := []map[string]interface{}{}
	for _, triggersItem := range tektonPipeline.Triggers {
		triggersItemMap, err := resourceIBMCdTektonPipelineTriggerToMap(triggersItem)
		if err != nil {
			return diag.FromErr(err)
		}
		triggers = append(triggers, triggersItemMap)
	}
	if err = d.Set("triggers", triggers); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting triggers: %s", err))
	}
	if err = d.Set("runs_url", tektonPipeline.RunsURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting runs_url: %s", err))
	}
	if !core.IsNil(tektonPipeline.Href) {
		if err = d.Set("href", tektonPipeline.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if err = d.Set("build_number", flex.IntValue(tektonPipeline.BuildNumber)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting build_number: %s", err))
	}
	if err = d.Set("enabled", tektonPipeline.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
	}

	return nil
}

func resourceIBMCdTektonPipelineUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateTektonPipelineOptions := &cdtektonpipelinev2.UpdateTektonPipelineOptions{}

	updateTektonPipelineOptions.SetID(d.Id())

	hasChange := false

	patchVals := &cdtektonpipelinev2.TektonPipelinePatch{}
	if d.HasChange("next_build_number") {
		newNextBuildNumber := int64(d.Get("next_build_number").(int))
		patchVals.NextBuildNumber = &newNextBuildNumber
		hasChange = true
	}
	if d.HasChange("enable_notifications") {
		newEnableNotifications := d.Get("enable_notifications").(bool)
		patchVals.EnableNotifications = &newEnableNotifications
		hasChange = true
	}
	if d.HasChange("enable_partial_cloning") {
		newEnablePartialCloning := d.Get("enable_partial_cloning").(bool)
		patchVals.EnablePartialCloning = &newEnablePartialCloning
		hasChange = true
	}
	if d.HasChange("worker") {
		worker, err := resourceIBMCdTektonPipelineMapToWorkerIdentity(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		patchVals.Worker = worker
		hasChange = true
	}

	if hasChange {
		updateTektonPipelineOptions.TektonPipelinePatch, _ = patchVals.AsPatch()
		_, response, err := cdTektonPipelineClient.UpdateTektonPipelineWithContext(context, updateTektonPipelineOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateTektonPipelineWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateTektonPipelineWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCdTektonPipelineRead(context, d, meta)
}

func resourceIBMCdTektonPipelineDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelineOptions := &cdtektonpipelinev2.DeleteTektonPipelineOptions{}

	deleteTektonPipelineOptions.SetID(d.Id())

	response, err := cdTektonPipelineClient.DeleteTektonPipelineWithContext(context, deleteTektonPipelineOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTektonPipelineWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTektonPipelineWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCdTektonPipelineMapToWorkerIdentity(modelMap map[string]interface{}) (*cdtektonpipelinev2.WorkerIdentity, error) {
	model := &cdtektonpipelinev2.WorkerIdentity{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMCdTektonPipelineWorkerIdentityToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceIBMCdTektonPipelineResourceGroupReferenceToMap(model *cdtektonpipelinev2.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineToolchainReferenceToMap(model *cdtektonpipelinev2.ToolchainReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["crn"] = model.CRN
	return modelMap, nil
}

func resourceIBMCdTektonPipelineDefinitionToMap(model *cdtektonpipelinev2.Definition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	sourceMap, err := resourceIBMCdTektonPipelineDefinitionSourceToMap(model.Source)
	if err != nil {
		return modelMap, err
	}
	modelMap["source"] = []map[string]interface{}{sourceMap}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceIBMCdTektonPipelineDefinitionSourceToMap(model *cdtektonpipelinev2.DefinitionSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	propertiesMap, err := resourceIBMCdTektonPipelineDefinitionSourcePropertiesToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineDefinitionSourcePropertiesToMap(model *cdtektonpipelinev2.DefinitionSourceProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["url"] = model.URL
	if model.Branch != nil {
		modelMap["branch"] = model.Branch
	}
	if model.Tag != nil {
		modelMap["tag"] = model.Tag
	}
	modelMap["path"] = model.Path
	if model.Tool != nil {
		toolMap, err := resourceIBMCdTektonPipelineToolToMap(model.Tool)
		if err != nil {
			return modelMap, err
		}
		modelMap["tool"] = []map[string]interface{}{toolMap}
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineToolToMap(model *cdtektonpipelinev2.Tool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceIBMCdTektonPipelinePropertyToMap(model *cdtektonpipelinev2.Property) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	if model.Enum != nil {
		modelMap["enum"] = model.Enum
	}
	modelMap["type"] = model.Type
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerToMap(model cdtektonpipelinev2.TriggerIntf) (map[string]interface{}, error) {
	if _, ok := model.(*cdtektonpipelinev2.TriggerManualTrigger); ok {
		return resourceIBMCdTektonPipelineTriggerManualTriggerToMap(model.(*cdtektonpipelinev2.TriggerManualTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerScmTrigger); ok {
		return resourceIBMCdTektonPipelineTriggerScmTriggerToMap(model.(*cdtektonpipelinev2.TriggerScmTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerTimerTrigger); ok {
		return resourceIBMCdTektonPipelineTriggerTimerTriggerToMap(model.(*cdtektonpipelinev2.TriggerTimerTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerGenericTrigger); ok {
		return resourceIBMCdTektonPipelineTriggerGenericTriggerToMap(model.(*cdtektonpipelinev2.TriggerGenericTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.Trigger); ok {
		modelMap := make(map[string]interface{})
		model := model.(*cdtektonpipelinev2.Trigger)
		if model.Type != nil {
			modelMap["type"] = model.Type
		}
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.Href != nil {
			modelMap["href"] = model.Href
		}
		if model.EventListener != nil {
			modelMap["event_listener"] = model.EventListener
		}
		if model.ID != nil {
			modelMap["id"] = model.ID
		}
		if model.Properties != nil {
			properties := []map[string]interface{}{}
			for _, propertiesItem := range model.Properties {
				propertiesItemMap, err := resourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem)
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
			workerMap, err := resourceIBMCdTektonPipelineWorkerToMap(model.Worker)
			if err != nil {
				return modelMap, err
			}
			modelMap["worker"] = []map[string]interface{}{workerMap}
		}
		if model.MaxConcurrentRuns != nil {
			modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
		}
		if model.Enabled != nil {
			modelMap["enabled"] = model.Enabled
		}
		if model.Source != nil {
			sourceMap, err := resourceIBMCdTektonPipelineTriggerSourceToMap(model.Source)
			if err != nil {
				return modelMap, err
			}
			modelMap["source"] = []map[string]interface{}{sourceMap}
		}
		if model.Events != nil {
			modelMap["events"] = model.Events
		}
		if model.Cron != nil {
			modelMap["cron"] = model.Cron
		}
		if model.Timezone != nil {
			modelMap["timezone"] = model.Timezone
		}
		if model.Secret != nil {
			secretMap, err := resourceIBMCdTektonPipelineGenericSecretToMap(model.Secret)
			if err != nil {
				return modelMap, err
			}
			modelMap["secret"] = []map[string]interface{}{secretMap}
		}
		if model.WebhookURL != nil {
			modelMap["webhook_url"] = model.WebhookURL
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized cdtektonpipelinev2.TriggerIntf subtype encountered")
	}
}

func resourceIBMCdTektonPipelineTriggerPropertyToMap(model *cdtektonpipelinev2.TriggerProperty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	if model.Enum != nil {
		modelMap["enum"] = model.Enum
	}
	modelMap["type"] = model.Type
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineWorkerToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerSourceToMap(model *cdtektonpipelinev2.TriggerSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	propertiesMap, err := resourceIBMCdTektonPipelineTriggerSourcePropertiesToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerSourcePropertiesToMap(model *cdtektonpipelinev2.TriggerSourceProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["url"] = model.URL
	if model.Branch != nil {
		modelMap["branch"] = model.Branch
	}
	if model.Pattern != nil {
		modelMap["pattern"] = model.Pattern
	}
	modelMap["blind_connection"] = model.BlindConnection
	if model.HookID != nil {
		modelMap["hook_id"] = model.HookID
	}
	toolMap, err := resourceIBMCdTektonPipelineToolToMap(model.Tool)
	if err != nil {
		return modelMap, err
	}
	modelMap["tool"] = []map[string]interface{}{toolMap}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineGenericSecretToMap(model *cdtektonpipelinev2.GenericSecret) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Source != nil {
		modelMap["source"] = model.Source
	}
	if model.KeyName != nil {
		modelMap["key_name"] = model.KeyName
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = model.Algorithm
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerManualTriggerToMap(model *cdtektonpipelinev2.TriggerManualTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	modelMap["event_listener"] = model.EventListener
	modelMap["id"] = model.ID
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := resourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem)
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
		workerMap, err := resourceIBMCdTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	modelMap["enabled"] = model.Enabled
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerScmTriggerToMap(model *cdtektonpipelinev2.TriggerScmTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	modelMap["event_listener"] = model.EventListener
	modelMap["id"] = model.ID
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := resourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem)
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
		workerMap, err := resourceIBMCdTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	modelMap["enabled"] = model.Enabled
	if model.Source != nil {
		sourceMap, err := resourceIBMCdTektonPipelineTriggerSourceToMap(model.Source)
		if err != nil {
			return modelMap, err
		}
		modelMap["source"] = []map[string]interface{}{sourceMap}
	}
	if model.Events != nil {
		modelMap["events"] = model.Events
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerTimerTriggerToMap(model *cdtektonpipelinev2.TriggerTimerTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	modelMap["event_listener"] = model.EventListener
	modelMap["id"] = model.ID
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := resourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem)
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
		workerMap, err := resourceIBMCdTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	modelMap["enabled"] = model.Enabled
	if model.Cron != nil {
		modelMap["cron"] = model.Cron
	}
	if model.Timezone != nil {
		modelMap["timezone"] = model.Timezone
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerGenericTriggerToMap(model *cdtektonpipelinev2.TriggerGenericTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	modelMap["event_listener"] = model.EventListener
	modelMap["id"] = model.ID
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := resourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem)
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
		workerMap, err := resourceIBMCdTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	modelMap["enabled"] = model.Enabled
	if model.Secret != nil {
		secretMap, err := resourceIBMCdTektonPipelineGenericSecretToMap(model.Secret)
		if err != nil {
			return modelMap, err
		}
		modelMap["secret"] = []map[string]interface{}{secretMap}
	}
	if model.WebhookURL != nil {
		modelMap["webhook_url"] = model.WebhookURL
	}
	return modelMap, nil
}
