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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtektonpipelinev2"
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
			"worker": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Details of the worker used to run the pipeline.",
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
			"next_build_number": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The build number that will be used for the next pipeline run.",
			},
			"enable_notifications": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag to enable notifications for this pipeline. If enabled, the Tekton pipeline run events will be published to all the destinations specified by the Slack and Event Notifications integrations in the parent toolchain. If omitted, this feature is disabled by default.",
			},
			"enable_partial_cloning": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Flag to enable partial cloning for this pipeline. When partial clone is enabled, only the files contained within the paths specified in definition repositories are read and cloned, this means that symbolic links might not work. If omitted, this feature is disabled by default.",
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
							Computed:    true,
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
							Computed:    true,
							Description: "Universally Unique Identifier.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
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
							Computed:    true,
							Description: "Source repository containing the Tekton pipeline definition.",
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
										Description: "Properties of the source, which define the URL of the repository and a branch or tag.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													ForceNew:    true,
													Description: "URL of the definition repository.",
												},
												"branch": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "A branch from the repo, specify one of branch or tag only.",
												},
												"tag": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "A tag from the repo, specify one of branch or tag only.",
												},
												"path": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The path to the definition's YAML files.",
												},
												"tool": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
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
							Optional:    true,
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
							Computed:    true,
							ForceNew:    true,
							Description: "Property name.",
						},
						"value": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: flex.SuppressPipelinePropertyRawSecret,
							Description:      "Property value. Any string value is valid.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "API URL for interacting with the property.",
						},
						"enum": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Options for `single_select` property type. Only needed when using `single_select` property type.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    true,
							Description: "Property type.",
						},
						"locked": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "When true, this property cannot be overridden by a trigger property or at runtime. Attempting to override it will result in run requests being rejected. The default is false.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Trigger type.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Trigger name.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.",
						},
						"event_listener": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Event listener name. The name of the event listener to which the trigger is associated. The event listeners are defined in the definition repositories of the Tekton pipeline.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The Trigger ID.",
						},
						"properties": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Optional trigger properties are used to override or supplement the pipeline properties when triggering a pipeline run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Description: "Property name.",
									},
									"value": &schema.Schema{
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										DiffSuppressFunc: flex.SuppressTriggerPropertyRawSecret,
										Description:      "Property value. Any string value is valid.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "API URL for interacting with the trigger property.",
									},
									"enum": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Options for `single_select` property type. Only needed for `single_select` property type.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Description: "Property type.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left blank the full tool integration data will be used.",
									},
									"locked": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "When true, this property cannot be overridden at runtime. Attempting to override it will result in run requests being rejected. The default is false.",
									},
								},
							},
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Optional trigger tags array.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"worker": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled for this trigger.",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Flag to check if the trigger is enabled.",
						},
						"favorite": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Mark the trigger as a favorite.",
						},
						"enable_events_from_forks": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "When enabled, pull request events from forks of the selected repository will trigger a pipeline run.",
						},
						"source": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
													ForceNew:    true,
													Description: "URL of the repository to which the trigger is listening.",
												},
												"branch": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Name of a branch from the repo. Only one of branch, pattern, or filter should be specified.",
												},
												"pattern": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
							Type:             schema.TypeList,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: flex.SuppressTriggerEvents,
							Description:      "Either 'events' or 'filter' is required specifically for Git triggers. Stores a list of events that a Git trigger listens to. Choose one or more from 'push', 'pull_request', and 'pull_request_closed'. If SCM repositories use the 'merge request' term, they correspond to the generic term i.e. 'pull request'.",
							Elem:             &schema.Schema{Type: schema.TypeString},
						},
						"filter": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Either 'events' or 'filter' can be used. Stores the CEL (Common Expression Language) expression value which is used for event filtering against the Git webhook payloads.",
						},
						"cron": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Only needed for timer triggers. CRON expression that indicates when this trigger will activate. Maximum frequency is every 5 minutes. The string is based on UNIX crontab syntax: minute, hour, day of month, month, day of week. Example: The CRON expression 0 *_/2 * * * - translates to - every 2 hours.",
						},
						"timezone": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Only used for timer triggers. Specify the timezone used for this timer trigger, which will ensure the CRON activates this trigger relative to the specified timezone. If no timezone is specified, the default timezone used is UTC. Valid timezones are those listed in the IANA timezone database, https://www.iana.org/time-zones.",
						},
						"secret": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Only needed for Generic Webhook trigger type. The secret is used to start the Generic Webhook trigger.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Secret type.",
									},
									"value": &schema.Schema{
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										DiffSuppressFunc: flex.SuppressGenericWebhookRawSecret,
										Description:      "Secret value, not needed if secret type is `internal_validation`.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Secret location, not needed if secret type is `internal_validation`.",
									},
									"key_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Secret name, not needed if type is `internal_validation`.",
									},
									"algorithm": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Algorithm used for `digest_matches` secret type. Only needed for `digest_matches` secret type.",
									},
								},
							},
						},
						"webhook_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
				Description: "Flag to check if the trigger is enabled.",
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createTektonPipelineOptions := &cdtektonpipelinev2.CreateTektonPipelineOptions{}

	if _, ok := d.GetOk("next_build_number"); ok {
		createTektonPipelineOptions.SetNextBuildNumber(int64(d.Get("next_build_number").(int)))
	}
	if _, ok := d.GetOk("enable_notifications"); ok {
		createTektonPipelineOptions.SetEnableNotifications(d.Get("enable_notifications").(bool))
	}
	if _, ok := d.GetOk("enable_partial_cloning"); ok {
		createTektonPipelineOptions.SetEnablePartialCloning(d.Get("enable_partial_cloning").(bool))
	}
	if _, ok := d.GetOk("worker"); ok {
		workerModel, err := ResourceIBMCdTektonPipelineMapToWorkerIdentity(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "create", "parse-worker").GetDiag()
		}
		createTektonPipelineOptions.SetWorker(workerModel)
	}

	if _, ok := d.GetOk("pipeline_id"); ok {
		createTektonPipelineOptions.SetID(d.Get("pipeline_id").(string))
	}
	tektonPipeline, _, err := cdTektonPipelineClient.CreateTektonPipelineWithContext(context, createTektonPipelineOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateTektonPipelineWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*tektonPipeline.ID)

	return resourceIBMCdTektonPipelineRead(context, d, meta)
}

func resourceIBMCdTektonPipelineRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getTektonPipelineOptions := &cdtektonpipelinev2.GetTektonPipelineOptions{}

	getTektonPipelineOptions.SetID(d.Id())

	tektonPipeline, response, err := cdTektonPipelineClient.GetTektonPipelineWithContext(context, getTektonPipelineOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTektonPipelineWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(tektonPipeline.Worker) {
		workerMap, err := ResourceIBMCdTektonPipelineWorkerToMap(tektonPipeline.Worker)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "worker-to-map").GetDiag()
		}
		if err = d.Set("worker", []map[string]interface{}{workerMap}); err != nil {
			err = fmt.Errorf("Error setting worker: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-worker").GetDiag()
		}
	}
	if !core.IsNil(tektonPipeline.NextBuildNumber) {
		if err = d.Set("next_build_number", flex.IntValue(tektonPipeline.NextBuildNumber)); err != nil {
			err = fmt.Errorf("Error setting next_build_number: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-next_build_number").GetDiag()
		}
	}
	if !core.IsNil(tektonPipeline.EnableNotifications) {
		if err = d.Set("enable_notifications", tektonPipeline.EnableNotifications); err != nil {
			err = fmt.Errorf("Error setting enable_notifications: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-enable_notifications").GetDiag()
		}
	}
	if !core.IsNil(tektonPipeline.EnablePartialCloning) {
		if err = d.Set("enable_partial_cloning", tektonPipeline.EnablePartialCloning); err != nil {
			err = fmt.Errorf("Error setting enable_partial_cloning: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-enable_partial_cloning").GetDiag()
		}
	}
	if err = d.Set("pipeline_id", tektonPipeline.ID); err != nil {
		err = fmt.Errorf("Error setting pipeline_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-pipeline_id").GetDiag()
	}
	if err = d.Set("name", tektonPipeline.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-name").GetDiag()
	}
	if err = d.Set("status", tektonPipeline.Status); err != nil {
		err = fmt.Errorf("Error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-status").GetDiag()
	}
	resourceGroupMap, err := ResourceIBMCdTektonPipelineResourceGroupReferenceToMap(tektonPipeline.ResourceGroup)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "resource_group-to-map").GetDiag()
	}
	if err = d.Set("resource_group", []map[string]interface{}{resourceGroupMap}); err != nil {
		err = fmt.Errorf("Error setting resource_group: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-resource_group").GetDiag()
	}
	toolchainMap, err := ResourceIBMCdTektonPipelineToolchainReferenceToMap(tektonPipeline.Toolchain)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "toolchain-to-map").GetDiag()
	}
	if err = d.Set("toolchain", []map[string]interface{}{toolchainMap}); err != nil {
		err = fmt.Errorf("Error setting toolchain: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-toolchain").GetDiag()
	}
	definitions := []map[string]interface{}{}
	for _, definitionsItem := range tektonPipeline.Definitions {
		definitionsItemMap, err := ResourceIBMCdTektonPipelineDefinitionToMap(&definitionsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "definitions-to-map").GetDiag()
		}
		definitions = append(definitions, definitionsItemMap)
	}
	if err = d.Set("definitions", definitions); err != nil {
		err = fmt.Errorf("Error setting definitions: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-definitions").GetDiag()
	}
	properties := []map[string]interface{}{}
	for _, propertiesItem := range tektonPipeline.Properties {
		propertiesItemMap, err := ResourceIBMCdTektonPipelinePropertyToMap(&propertiesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "properties-to-map").GetDiag()
		}
		properties = append(properties, propertiesItemMap)
	}
	if err = d.Set("properties", properties); err != nil {
		err = fmt.Errorf("Error setting properties: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-properties").GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(tektonPipeline.UpdatedAt)); err != nil {
		err = fmt.Errorf("Error setting updated_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-updated_at").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(tektonPipeline.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-created_at").GetDiag()
	}
	triggers := []map[string]interface{}{}
	for _, triggersItem := range tektonPipeline.Triggers {
		triggersItemMap, err := ResourceIBMCdTektonPipelineTriggerToMap(triggersItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "triggers-to-map").GetDiag()
		}
		triggers = append(triggers, triggersItemMap)
	}
	if err = d.Set("triggers", triggers); err != nil {
		err = fmt.Errorf("Error setting triggers: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-triggers").GetDiag()
	}
	if err = d.Set("runs_url", tektonPipeline.RunsURL); err != nil {
		err = fmt.Errorf("Error setting runs_url: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-runs_url").GetDiag()
	}
	if !core.IsNil(tektonPipeline.Href) {
		if err = d.Set("href", tektonPipeline.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-href").GetDiag()
		}
	}
	if err = d.Set("build_number", flex.IntValue(tektonPipeline.BuildNumber)); err != nil {
		err = fmt.Errorf("Error setting build_number: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-build_number").GetDiag()
	}
	if err = d.Set("enabled", tektonPipeline.Enabled); err != nil {
		err = fmt.Errorf("Error setting enabled: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "read", "set-enabled").GetDiag()
	}

	return nil
}

func resourceIBMCdTektonPipelineUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
		worker, err := ResourceIBMCdTektonPipelineMapToWorkerIdentity(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "update", "parse-worker").GetDiag()
		}
		patchVals.Worker = worker
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateTektonPipelineOptions.TektonPipelinePatch = ResourceIBMCdTektonPipelineTektonPipelinePatchAsPatch(patchVals, d)

		_, _, err = cdTektonPipelineClient.UpdateTektonPipelineWithContext(context, updateTektonPipelineOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateTektonPipelineWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCdTektonPipelineRead(context, d, meta)
}

func resourceIBMCdTektonPipelineDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteTektonPipelineOptions := &cdtektonpipelinev2.DeleteTektonPipelineOptions{}

	deleteTektonPipelineOptions.SetID(d.Id())

	_, err = cdTektonPipelineClient.DeleteTektonPipelineWithContext(context, deleteTektonPipelineOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteTektonPipelineWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMCdTektonPipelineMapToWorkerIdentity(modelMap map[string]interface{}) (*cdtektonpipelinev2.WorkerIdentity, error) {
	model := &cdtektonpipelinev2.WorkerIdentity{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMCdTektonPipelineWorkerToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
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

func ResourceIBMCdTektonPipelineResourceGroupReferenceToMap(model *cdtektonpipelinev2.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineToolchainReferenceToMap(model *cdtektonpipelinev2.ToolchainReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["crn"] = *model.CRN
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineDefinitionToMap(model *cdtektonpipelinev2.Definition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	sourceMap, err := ResourceIBMCdTektonPipelineDefinitionSourceToMap(model.Source)
	if err != nil {
		return modelMap, err
	}
	modelMap["source"] = []map[string]interface{}{sourceMap}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineDefinitionSourceToMap(model *cdtektonpipelinev2.DefinitionSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	propertiesMap, err := ResourceIBMCdTektonPipelineDefinitionSourcePropertiesToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineDefinitionSourcePropertiesToMap(model *cdtektonpipelinev2.DefinitionSourceProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["url"] = *model.URL
	if model.Branch != nil {
		modelMap["branch"] = *model.Branch
	}
	if model.Tag != nil {
		modelMap["tag"] = *model.Tag
	}
	modelMap["path"] = *model.Path
	if model.Tool != nil {
		toolMap, err := ResourceIBMCdTektonPipelineToolToMap(model.Tool)
		if err != nil {
			return modelMap, err
		}
		modelMap["tool"] = []map[string]interface{}{toolMap}
	}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineToolToMap(model *cdtektonpipelinev2.Tool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func ResourceIBMCdTektonPipelinePropertyToMap(model *cdtektonpipelinev2.Property) (map[string]interface{}, error) {
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
	if model.Locked != nil {
		modelMap["locked"] = *model.Locked
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTriggerToMap(model cdtektonpipelinev2.TriggerIntf) (map[string]interface{}, error) {
	if _, ok := model.(*cdtektonpipelinev2.TriggerManualTrigger); ok {
		return ResourceIBMCdTektonPipelineTriggerManualTriggerToMap(model.(*cdtektonpipelinev2.TriggerManualTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerScmTrigger); ok {
		return ResourceIBMCdTektonPipelineTriggerScmTriggerToMap(model.(*cdtektonpipelinev2.TriggerScmTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerTimerTrigger); ok {
		return ResourceIBMCdTektonPipelineTriggerTimerTriggerToMap(model.(*cdtektonpipelinev2.TriggerTimerTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerGenericTrigger); ok {
		return ResourceIBMCdTektonPipelineTriggerGenericTriggerToMap(model.(*cdtektonpipelinev2.TriggerGenericTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.Trigger); ok {
		modelMap := make(map[string]interface{})
		model := model.(*cdtektonpipelinev2.Trigger)
		if model.Type != nil {
			modelMap["type"] = *model.Type
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
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
				propertiesItemMap, err := ResourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem) // #nosec G601
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
			workerMap, err := ResourceIBMCdTektonPipelineWorkerToMap(model.Worker)
			if err != nil {
				return modelMap, err
			}
			modelMap["worker"] = []map[string]interface{}{workerMap}
		}
		if model.MaxConcurrentRuns != nil {
			modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
		}
		if model.Enabled != nil {
			modelMap["enabled"] = *model.Enabled
		}
		if model.Favorite != nil {
			modelMap["favorite"] = *model.Favorite
		}
		if model.EnableEventsFromForks != nil {
			modelMap["enable_events_from_forks"] = *model.EnableEventsFromForks
		}
		if model.Source != nil {
			sourceMap, err := ResourceIBMCdTektonPipelineTriggerSourceToMap(model.Source)
			if err != nil {
				return modelMap, err
			}
			modelMap["source"] = []map[string]interface{}{sourceMap}
		}
		if model.Events != nil {
			modelMap["events"] = model.Events
		}
		if model.Filter != nil {
			modelMap["filter"] = *model.Filter
		}
		if model.Cron != nil {
			modelMap["cron"] = *model.Cron
		}
		if model.Timezone != nil {
			modelMap["timezone"] = *model.Timezone
		}
		if model.Secret != nil {
			secretMap, err := ResourceIBMCdTektonPipelineGenericSecretToMap(model.Secret)
			if err != nil {
				return modelMap, err
			}
			modelMap["secret"] = []map[string]interface{}{secretMap}
		}
		if model.WebhookURL != nil {
			modelMap["webhook_url"] = *model.WebhookURL
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized cdtektonpipelinev2.TriggerIntf subtype encountered")
	}
}

func ResourceIBMCdTektonPipelineTriggerPropertyToMap(model *cdtektonpipelinev2.TriggerProperty) (map[string]interface{}, error) {
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

func ResourceIBMCdTektonPipelineTriggerSourceToMap(model *cdtektonpipelinev2.TriggerSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	propertiesMap, err := ResourceIBMCdTektonPipelineTriggerSourcePropertiesToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTriggerSourcePropertiesToMap(model *cdtektonpipelinev2.TriggerSourceProperties) (map[string]interface{}, error) {
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
	toolMap, err := ResourceIBMCdTektonPipelineToolToMap(model.Tool)
	if err != nil {
		return modelMap, err
	}
	modelMap["tool"] = []map[string]interface{}{toolMap}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineGenericSecretToMap(model *cdtektonpipelinev2.GenericSecret) (map[string]interface{}, error) {
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

func ResourceIBMCdTektonPipelineTriggerManualTriggerToMap(model *cdtektonpipelinev2.TriggerManualTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["name"] = *model.Name
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	modelMap["event_listener"] = *model.EventListener
	modelMap["id"] = *model.ID
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := ResourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem) // #nosec G601
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
		workerMap, err := ResourceIBMCdTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	modelMap["enabled"] = *model.Enabled
	if model.Favorite != nil {
		modelMap["favorite"] = *model.Favorite
	}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTriggerScmTriggerToMap(model *cdtektonpipelinev2.TriggerScmTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["name"] = *model.Name
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	modelMap["event_listener"] = *model.EventListener
	modelMap["id"] = *model.ID
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := ResourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem) // #nosec G601
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
		workerMap, err := ResourceIBMCdTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	modelMap["enabled"] = *model.Enabled
	if model.Favorite != nil {
		modelMap["favorite"] = *model.Favorite
	}
	if model.EnableEventsFromForks != nil {
		modelMap["enable_events_from_forks"] = *model.EnableEventsFromForks
	}
	if model.Source != nil {
		sourceMap, err := ResourceIBMCdTektonPipelineTriggerSourceToMap(model.Source)
		if err != nil {
			return modelMap, err
		}
		modelMap["source"] = []map[string]interface{}{sourceMap}
	}
	if model.Events != nil {
		modelMap["events"] = model.Events
	}
	if model.Filter != nil {
		modelMap["filter"] = *model.Filter
	}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTriggerTimerTriggerToMap(model *cdtektonpipelinev2.TriggerTimerTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["name"] = *model.Name
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	modelMap["event_listener"] = *model.EventListener
	modelMap["id"] = *model.ID
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := ResourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem) // #nosec G601
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
		workerMap, err := ResourceIBMCdTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	modelMap["enabled"] = *model.Enabled
	if model.Favorite != nil {
		modelMap["favorite"] = *model.Favorite
	}
	if model.Cron != nil {
		modelMap["cron"] = *model.Cron
	}
	if model.Timezone != nil {
		modelMap["timezone"] = *model.Timezone
	}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTriggerGenericTriggerToMap(model *cdtektonpipelinev2.TriggerGenericTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	modelMap["name"] = *model.Name
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	modelMap["event_listener"] = *model.EventListener
	modelMap["id"] = *model.ID
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := ResourceIBMCdTektonPipelineTriggerPropertyToMap(&propertiesItem) // #nosec G601
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
		workerMap, err := ResourceIBMCdTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	modelMap["enabled"] = *model.Enabled
	if model.Favorite != nil {
		modelMap["favorite"] = *model.Favorite
	}
	if model.Secret != nil {
		secretMap, err := ResourceIBMCdTektonPipelineGenericSecretToMap(model.Secret)
		if err != nil {
			return modelMap, err
		}
		modelMap["secret"] = []map[string]interface{}{secretMap}
	}
	if model.WebhookURL != nil {
		modelMap["webhook_url"] = *model.WebhookURL
	}
	if model.Filter != nil {
		modelMap["filter"] = *model.Filter
	}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTektonPipelinePatchAsPatch(patchVals *cdtektonpipelinev2.TektonPipelinePatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "next_build_number"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["next_build_number"] = nil
	}
	path = "enable_notifications"
	if _, exists := d.GetOkExists(path); d.HasChange(path) && !exists {
		patch["enable_notifications"] = nil
	}
	path = "enable_partial_cloning"
	if _, exists := d.GetOkExists(path); d.HasChange(path) && !exists {
		patch["enable_partial_cloning"] = nil
	}
	path = "worker"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["worker"] = nil
	}

	return patch
}
