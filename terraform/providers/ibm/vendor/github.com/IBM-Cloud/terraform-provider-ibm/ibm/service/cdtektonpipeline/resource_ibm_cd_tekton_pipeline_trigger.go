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
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtektonpipelinev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMCdTektonPipelineTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdTektonPipelineTriggerCreate,
		ReadContext:   resourceIBMCdTektonPipelineTriggerRead,
		UpdateContext: resourceIBMCdTektonPipelineTriggerUpdate,
		DeleteContext: resourceIBMCdTektonPipelineTriggerDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger", "pipeline_id"),
				Description:  "The Tekton pipeline ID.",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger", "type"),
				Description:  "Trigger type.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger", "name"),
				Description:  "Trigger name.",
			},
			"event_listener": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger", "event_listener"),
				Description:  "Event listener name. The name of the event listener to which the trigger is associated. The event listeners are defined in the definition repositories of the Tekton pipeline.",
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
				Description: "Flag to check if the trigger is enabled.",
			},
			"favorite": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Mark the trigger as a favorite.",
			},
			"enable_events_from_forks": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When enabled, pull request events from forks of the selected repository will trigger a pipeline run.",
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
										Description: "Name of a branch from the repo. Only one of branch, pattern, or filter should be specified.",
									},
									"pattern": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
				DiffSuppressFunc: flex.SuppressTriggerEvents,
				Description:      "Either 'events' or 'filter' is required specifically for Git triggers. Stores a list of events that a Git trigger listens to. Choose one or more from 'push', 'pull_request', and 'pull_request_closed'. If SCM repositories use the 'merge request' term, they correspond to the generic term i.e. 'pull request'.",
				Elem:             &schema.Schema{Type: schema.TypeString},
			},
			"filter": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger", "filter"),
				Description:  "Either 'events' or 'filter' can be used. Stores the CEL (Common Expression Language) expression value which is used for event filtering against the Git webhook payloads.",
			},
			"cron": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger", "cron"),
				Description:  "Only needed for timer triggers. CRON expression that indicates when this trigger will activate. Maximum frequency is every 5 minutes. The string is based on UNIX crontab syntax: minute, hour, day of month, month, day of week. Example: The CRON expression 0 *_/2 * * * - translates to - every 2 hours.",
			},
			"timezone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger", "timezone"),
				Description:  "Only used for timer triggers. Specify the timezone used for this timer trigger, which will ensure the CRON activates this trigger relative to the specified timezone. If no timezone is specified, the default timezone used is UTC. Valid timezones are those listed in the IANA timezone database, https://www.iana.org/time-zones.",
			},
			"secret": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Only needed for Generic Webhook trigger type. The secret is used to start the Generic Webhook trigger.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.",
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
			"webhook_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Webhook URL that can be used to trigger pipeline runs.",
			},
			"trigger_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Trigger ID.",
			},
		},
	}
}

func ResourceIBMCdTektonPipelineTriggerValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "pipeline_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z]+$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "generic, manual, scm, timer",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-zA-Z0-9]{1,2}|[a-zA-Z0-9][0-9a-zA-Z-_.: \/\(\)\[\]]{1,251}[a-zA-Z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "event_listener",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-zA-Z_.]{1,253}$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "filter",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^.*$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "cron",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[-0-9a-zA-Z,\*\/ ]{5,253}$`,
			MinValueLength:             5,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "timezone",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[-0-9a-zA-Z+_., \/]{1,253}$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_tekton_pipeline_trigger", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdTektonPipelineTriggerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createTektonPipelineTriggerOptions := &cdtektonpipelinev2.CreateTektonPipelineTriggerOptions{}

	createTektonPipelineTriggerOptions.SetPipelineID(d.Get("pipeline_id").(string))
	createTektonPipelineTriggerOptions.SetType(d.Get("type").(string))
	createTektonPipelineTriggerOptions.SetName(d.Get("name").(string))
	createTektonPipelineTriggerOptions.SetEventListener(d.Get("event_listener").(string))
	if _, ok := d.GetOk("tags"); ok {
		var tags []string
		for _, v := range d.Get("tags").([]interface{}) {
			tagsItem := v.(string)
			tags = append(tags, tagsItem)
		}
		createTektonPipelineTriggerOptions.SetTags(tags)
	}
	if _, ok := d.GetOk("worker"); ok {
		workerModel, err := ResourceIBMCdTektonPipelineTriggerMapToWorkerIdentity(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "create", "parse-worker").GetDiag()
		}
		createTektonPipelineTriggerOptions.SetWorker(workerModel)
	}
	if _, ok := d.GetOk("max_concurrent_runs"); ok {
		createTektonPipelineTriggerOptions.SetMaxConcurrentRuns(int64(d.Get("max_concurrent_runs").(int)))
	}
	if _, ok := d.GetOkExists("enabled"); ok {
		createTektonPipelineTriggerOptions.SetEnabled(d.Get("enabled").(bool))
	}
	if _, ok := d.GetOk("secret"); ok {
		secretModel, err := ResourceIBMCdTektonPipelineTriggerMapToGenericSecret(d.Get("secret.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "create", "parse-secret").GetDiag()
		}
		createTektonPipelineTriggerOptions.SetSecret(secretModel)
	}
	if _, ok := d.GetOk("cron"); ok {
		createTektonPipelineTriggerOptions.SetCron(d.Get("cron").(string))
	}
	if _, ok := d.GetOk("timezone"); ok {
		createTektonPipelineTriggerOptions.SetTimezone(d.Get("timezone").(string))
	}
	if _, ok := d.GetOk("source"); ok {
		sourceModel, err := ResourceIBMCdTektonPipelineTriggerMapToTriggerSourcePrototype(d.Get("source.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "create", "parse-source").GetDiag()
		}
		createTektonPipelineTriggerOptions.SetSource(sourceModel)
	}
	if _, ok := d.GetOk("events"); ok {
		eventsInterface := d.Get("events").([]interface{})
		events := make([]string, len(eventsInterface))
		for i, v := range eventsInterface {
			events[i] = fmt.Sprint(v)
		}
		sort.Strings(events)
		createTektonPipelineTriggerOptions.SetEvents(events)
	}
	if _, ok := d.GetOk("filter"); ok {
		createTektonPipelineTriggerOptions.SetFilter(d.Get("filter").(string))
	}
	if _, ok := d.GetOk("favorite"); ok {
		createTektonPipelineTriggerOptions.SetFavorite(d.Get("favorite").(bool))
	}
	if _, ok := d.GetOk("enable_events_from_forks"); ok {
		createTektonPipelineTriggerOptions.SetEnableEventsFromForks(d.Get("enable_events_from_forks").(bool))
	}

	triggerIntf, _, err := cdTektonPipelineClient.CreateTektonPipelineTriggerWithContext(context, createTektonPipelineTriggerOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateTektonPipelineTriggerWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_trigger", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	trigger := triggerIntf.(*cdtektonpipelinev2.Trigger)
	d.SetId(fmt.Sprintf("%s/%s", *createTektonPipelineTriggerOptions.PipelineID, *trigger.ID))

	return resourceIBMCdTektonPipelineTriggerRead(context, d, meta)
}

func resourceIBMCdTektonPipelineTriggerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getTektonPipelineTriggerOptions := &cdtektonpipelinev2.GetTektonPipelineTriggerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "sep-id-parts").GetDiag()
	}

	getTektonPipelineTriggerOptions.SetPipelineID(parts[0])
	getTektonPipelineTriggerOptions.SetTriggerID(parts[1])

	triggerIntf, response, err := cdTektonPipelineClient.GetTektonPipelineTriggerWithContext(context, getTektonPipelineTriggerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTektonPipelineTriggerWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_trigger", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	trigger := triggerIntf.(*cdtektonpipelinev2.Trigger)
	if err = d.Set("type", trigger.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-type").GetDiag()
	}
	if err = d.Set("name", trigger.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-name").GetDiag()
	}
	if err = d.Set("event_listener", trigger.EventListener); err != nil {
		err = fmt.Errorf("Error setting event_listener: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-event_listener").GetDiag()
	}
	if !core.IsNil(trigger.Tags) {
		if err = d.Set("tags", trigger.Tags); err != nil {
			err = fmt.Errorf("Error setting tags: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-tags").GetDiag()
		}
	}
	if !core.IsNil(trigger.Worker) {
		workerMap, err := ResourceIBMCdTektonPipelineTriggerWorkerToMap(trigger.Worker)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "worker-to-map").GetDiag()
		}
		if err = d.Set("worker", []map[string]interface{}{workerMap}); err != nil {
			err = fmt.Errorf("Error setting worker: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-worker").GetDiag()
		}
	}
	if !core.IsNil(trigger.MaxConcurrentRuns) {
		if err = d.Set("max_concurrent_runs", flex.IntValue(trigger.MaxConcurrentRuns)); err != nil {
			err = fmt.Errorf("Error setting max_concurrent_runs: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-max_concurrent_runs").GetDiag()
		}
	}
	if !core.IsNil(trigger.Enabled) {
		if err = d.Set("enabled", trigger.Enabled); err != nil {
			err = fmt.Errorf("Error setting enabled: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-enabled").GetDiag()
		}
	}
	if !core.IsNil(trigger.Favorite) {
		if err = d.Set("favorite", trigger.Favorite); err != nil {
			err = fmt.Errorf("Error setting favorite: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-favorite").GetDiag()
		}
	}
	if !core.IsNil(trigger.EnableEventsFromForks) {
		if err = d.Set("enable_events_from_forks", trigger.EnableEventsFromForks); err != nil {
			err = fmt.Errorf("Error setting enable_events_from_forks: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-enable_events_from_forks").GetDiag()
		}
	}
	if !core.IsNil(trigger.Source) {
		sourceMap, err := ResourceIBMCdTektonPipelineTriggerTriggerSourceToMap(trigger.Source)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "source-to-map").GetDiag()
		}
		if err = d.Set("source", []map[string]interface{}{sourceMap}); err != nil {
			err = fmt.Errorf("Error setting source: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-source").GetDiag()
		}
	}
	if !core.IsNil(trigger.Events) {
		if err = d.Set("events", trigger.Events); err != nil {
			err = fmt.Errorf("Error setting events: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-events").GetDiag()
		}
	}
	if !core.IsNil(trigger.Filter) {
		if err = d.Set("filter", trigger.Filter); err != nil {
			err = fmt.Errorf("Error setting filter: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-filter").GetDiag()
		}
	}
	if !core.IsNil(trigger.Cron) {
		if err = d.Set("cron", trigger.Cron); err != nil {
			err = fmt.Errorf("Error setting cron: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-cron").GetDiag()
		}
	}
	if !core.IsNil(trigger.Timezone) {
		if err = d.Set("timezone", trigger.Timezone); err != nil {
			err = fmt.Errorf("Error setting timezone: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-timezone").GetDiag()
		}
	}
	if !core.IsNil(trigger.Secret) {
		secretMap, err := ResourceIBMCdTektonPipelineTriggerGenericSecretToMap(trigger.Secret)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "secret-to-map").GetDiag()
		}
		if err = d.Set("secret", []map[string]interface{}{secretMap}); err != nil {
			err = fmt.Errorf("Error setting secret: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-secret").GetDiag()
		}
	}
	if !core.IsNil(trigger.Href) {
		if err = d.Set("href", trigger.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-href").GetDiag()
		}
	}
	if !core.IsNil(trigger.Properties) {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range trigger.Properties {
			propertiesItemMap, err := ResourceIBMCdTektonPipelineTriggerTriggerPropertyToMap(&propertiesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "properties-to-map").GetDiag()
			}
			properties = append(properties, propertiesItemMap)
		}
		if err = d.Set("properties", properties); err != nil {
			err = fmt.Errorf("Error setting properties: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-properties").GetDiag()
		}
	}
	if !core.IsNil(trigger.WebhookURL) {
		if err = d.Set("webhook_url", trigger.WebhookURL); err != nil {
			err = fmt.Errorf("Error setting webhook_url: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-webhook_url").GetDiag()
		}
	}
	if !core.IsNil(trigger.ID) {
		if err = d.Set("trigger_id", trigger.ID); err != nil {
			err = fmt.Errorf("Error setting trigger_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "read", "set-trigger_id").GetDiag()
		}
	}

	return nil
}

func resourceIBMCdTektonPipelineTriggerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateTektonPipelineTriggerOptions := &cdtektonpipelinev2.UpdateTektonPipelineTriggerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "update", "sep-id-parts").GetDiag()
	}

	updateTektonPipelineTriggerOptions.SetPipelineID(parts[0])
	updateTektonPipelineTriggerOptions.SetTriggerID(parts[1])

	hasChange := false

	patchVals := &cdtektonpipelinev2.TriggerPatch{}
	if d.HasChange("pipeline_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "pipeline_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_tekton_pipeline_trigger", "update", "pipeline_id-forces-new").GetDiag()
	}
	if d.HasChange("type") {
		newType := d.Get("type").(string)
		patchVals.Type = &newType
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("event_listener") {
		newEventListener := d.Get("event_listener").(string)
		patchVals.EventListener = &newEventListener
		hasChange = true
	}
	if d.HasChange("tags") {
		var tags []string
		for _, v := range d.Get("tags").([]interface{}) {
			tagsItem := v.(string)
			tags = append(tags, tagsItem)
		}
		patchVals.Tags = tags
		hasChange = true
	}
	if d.HasChange("worker") {
		worker, err := ResourceIBMCdTektonPipelineTriggerMapToWorkerIdentity(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "update", "parse-worker").GetDiag()
		}
		patchVals.Worker = worker
		hasChange = true
	}
	if d.HasChange("max_concurrent_runs") {
		newMaxConcurrentRuns := int64(d.Get("max_concurrent_runs").(int))
		patchVals.MaxConcurrentRuns = &newMaxConcurrentRuns
		hasChange = true
	}
	if d.HasChange("enabled") {
		newEnabled := d.Get("enabled").(bool)
		patchVals.Enabled = &newEnabled
		hasChange = true
	}
	if d.HasChange("secret") {
		secret, err := ResourceIBMCdTektonPipelineTriggerMapToGenericSecret(d.Get("secret.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "update", "parse-secret").GetDiag()
		}
		patchVals.Secret = secret
		hasChange = true
	}
	if d.HasChange("cron") {
		newCron := d.Get("cron").(string)
		patchVals.Cron = &newCron
		hasChange = true
	}
	if d.HasChange("timezone") {
		newTimezone := d.Get("timezone").(string)
		patchVals.Timezone = &newTimezone
		hasChange = true
	}
	if d.HasChange("source") {
		source, err := ResourceIBMCdTektonPipelineTriggerMapToTriggerSourcePrototype(d.Get("source.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "update", "parse-source").GetDiag()
		}
		patchVals.Source = source
		hasChange = true
	}

	if d.HasChange("events") {
		events := []string{}
		for _, eventsItem := range d.Get("events").([]interface{}) {
			events = append(events, eventsItem.(string))
		}
		sort.Strings(events)
		patchVals.Events = events
		hasChange = true
	}

	if d.HasChange("filter") {
		newFilter := d.Get("filter").(string)
		patchVals.Filter = &newFilter
		hasChange = true
	}
	if d.HasChange("favorite") {
		newFavorite := d.Get("favorite").(bool)
		patchVals.Favorite = &newFavorite
		hasChange = true
	}
	if d.HasChange("enable_events_from_forks") {
		newEnableEventsFromForks := d.Get("enable_events_from_forks").(bool)
		patchVals.EnableEventsFromForks = &newEnableEventsFromForks
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateTektonPipelineTriggerOptions.TriggerPatch = ResourceIBMCdTektonPipelineTriggerTriggerPatchAsPatch(patchVals, d)

		_, _, err = cdTektonPipelineClient.UpdateTektonPipelineTriggerWithContext(context, updateTektonPipelineTriggerOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateTektonPipelineTriggerWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_trigger", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCdTektonPipelineTriggerRead(context, d, meta)
}

func resourceIBMCdTektonPipelineTriggerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteTektonPipelineTriggerOptions := &cdtektonpipelinev2.DeleteTektonPipelineTriggerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger", "delete", "sep-id-parts").GetDiag()
	}

	deleteTektonPipelineTriggerOptions.SetPipelineID(parts[0])
	deleteTektonPipelineTriggerOptions.SetTriggerID(parts[1])

	_, err = cdTektonPipelineClient.DeleteTektonPipelineTriggerWithContext(context, deleteTektonPipelineTriggerOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteTektonPipelineTriggerWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_trigger", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMCdTektonPipelineTriggerMapToWorkerIdentity(modelMap map[string]interface{}) (*cdtektonpipelinev2.WorkerIdentity, error) {
	model := &cdtektonpipelinev2.WorkerIdentity{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMCdTektonPipelineTriggerMapToGenericSecret(modelMap map[string]interface{}) (*cdtektonpipelinev2.GenericSecret, error) {
	model := &cdtektonpipelinev2.GenericSecret{}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	if modelMap["source"] != nil && modelMap["source"].(string) != "" {
		model.Source = core.StringPtr(modelMap["source"].(string))
	}
	if modelMap["key_name"] != nil && modelMap["key_name"].(string) != "" {
		model.KeyName = core.StringPtr(modelMap["key_name"].(string))
	}
	if modelMap["algorithm"] != nil && modelMap["algorithm"].(string) != "" {
		model.Algorithm = core.StringPtr(modelMap["algorithm"].(string))
	}
	return model, nil
}

func ResourceIBMCdTektonPipelineTriggerMapToTriggerSourcePrototype(modelMap map[string]interface{}) (*cdtektonpipelinev2.TriggerSourcePrototype, error) {
	model := &cdtektonpipelinev2.TriggerSourcePrototype{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	PropertiesModel, err := ResourceIBMCdTektonPipelineTriggerMapToTriggerSourcePropertiesPrototype(modelMap["properties"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Properties = PropertiesModel
	return model, nil
}

func ResourceIBMCdTektonPipelineTriggerMapToTriggerSourcePropertiesPrototype(modelMap map[string]interface{}) (*cdtektonpipelinev2.TriggerSourcePropertiesPrototype, error) {
	model := &cdtektonpipelinev2.TriggerSourcePropertiesPrototype{}
	model.URL = core.StringPtr(modelMap["url"].(string))
	if modelMap["branch"] != nil && modelMap["branch"].(string) != "" {
		model.Branch = core.StringPtr(modelMap["branch"].(string))
	}
	if modelMap["pattern"] != nil && modelMap["pattern"].(string) != "" {
		model.Pattern = core.StringPtr(modelMap["pattern"].(string))
	}
	return model, nil
}

func ResourceIBMCdTektonPipelineTriggerWorkerToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
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

func ResourceIBMCdTektonPipelineTriggerTriggerSourceToMap(model *cdtektonpipelinev2.TriggerSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = *model.Type
	propertiesMap, err := ResourceIBMCdTektonPipelineTriggerTriggerSourcePropertiesToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTriggerTriggerSourcePropertiesToMap(model *cdtektonpipelinev2.TriggerSourceProperties) (map[string]interface{}, error) {
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
	toolMap, err := ResourceIBMCdTektonPipelineTriggerToolToMap(model.Tool)
	if err != nil {
		return modelMap, err
	}
	modelMap["tool"] = []map[string]interface{}{toolMap}
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTriggerToolToMap(model *cdtektonpipelinev2.Tool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func ResourceIBMCdTektonPipelineTriggerGenericSecretToMap(model *cdtektonpipelinev2.GenericSecret) (map[string]interface{}, error) {
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

func ResourceIBMCdTektonPipelineTriggerTriggerPropertyToMap(model *cdtektonpipelinev2.TriggerProperty) (map[string]interface{}, error) {
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

func ResourceIBMCdTektonPipelineTriggerTriggerPatchAsPatch(patchVals *cdtektonpipelinev2.TriggerPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	}
	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}
	path = "event_listener"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["event_listener"] = nil
	}
	path = "tags"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["tags"] = nil
	}
	path = "worker"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["worker"] = nil
	}
	path = "max_concurrent_runs"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["max_concurrent_runs"] = nil
	}
	path = "enabled"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["enabled"] = nil
	}
	path = "secret"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["secret"] = nil
	} else if exists && patch["secret"] != nil {
		ResourceIBMCdTektonPipelineTriggerGenericSecretAsPatch(patch["secret"].(map[string]interface{}), d)
	}
	path = "cron"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["cron"] = nil
	}
	path = "timezone"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["timezone"] = nil
	}
	path = "source"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["source"] = nil
	}
	path = "events"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["events"] = nil
	}
	path = "filter"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["filter"] = nil
	}
	path = "favorite"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["favorite"] = nil
	}
	path = "enable_events_from_forks"
	if _, exists := d.GetOkExists(path); d.HasChange(path) && !exists {
		patch["enable_events_from_forks"] = nil
	}

	return patch
}

func ResourceIBMCdTektonPipelineTriggerGenericSecretAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "secret.0.type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	}
	path = "secret.0.value"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["value"] = nil
	}
	path = "secret.0.source"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["source"] = nil
	}
	path = "secret.0.key_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["key_name"] = nil
	}
	path = "secret.0.algorithm"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["algorithm"] = nil
	}
}
