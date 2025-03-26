// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

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
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
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
				Description: "Trigger tags array.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"worker": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Specify the worker used to run the trigger. Use `worker: { id: 'public' }` to use the IBM Managed workers. The default is to inherit the worker set in the pipeline settings, which can also be explicitly set using `worker: { id: 'inherit' }`.",
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
			"max_concurrent_runs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled for this trigger.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag to check if the trigger is enabled. If omitted the trigger is enabled by default.",
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
				Description: "Only used for SCM triggers. When enabled, pull request events from forks of the selected repository will trigger a pipeline run.",
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
						"locked": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
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
		validate.ValidateSchema{
			Identifier:                 "filter",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^.*$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_tekton_pipeline_trigger", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdTektonPipelineTriggerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelineTriggerOptions := &cdtektonpipelinev2.CreateTektonPipelineTriggerOptions{}

	createTektonPipelineTriggerOptions.SetPipelineID(d.Get("pipeline_id").(string))
	createTektonPipelineTriggerOptions.SetType(d.Get("type").(string))
	createTektonPipelineTriggerOptions.SetName(d.Get("name").(string))
	createTektonPipelineTriggerOptions.SetEventListener(d.Get("event_listener").(string))
	if _, ok := d.GetOk("tags"); ok {
		tags := []string{}
		for _, tagsItem := range d.Get("tags").([]interface{}) {
			tags = append(tags, tagsItem.(string))
		}
		createTektonPipelineTriggerOptions.SetTags(tags)
	}
	if _, ok := d.GetOk("worker"); ok {
		workerModel, err := resourceIBMCdTektonPipelineTriggerMapToWorkerIdentity(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
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
		secretModel, err := resourceIBMCdTektonPipelineTriggerMapToGenericSecret(d.Get("secret.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
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
		sourceModel, err := resourceIBMCdTektonPipelineTriggerMapToTriggerSourcePrototype(d.Get("source.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
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

	triggerIntf, response, err := cdTektonPipelineClient.CreateTektonPipelineTriggerWithContext(context, createTektonPipelineTriggerOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTektonPipelineTriggerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTektonPipelineTriggerWithContext failed %s\n%s", err, response))
	}

	trigger := triggerIntf.(*cdtektonpipelinev2.Trigger)
	d.SetId(fmt.Sprintf("%s/%s", *createTektonPipelineTriggerOptions.PipelineID, *trigger.ID))

	return resourceIBMCdTektonPipelineTriggerRead(context, d, meta)
}

func resourceIBMCdTektonPipelineTriggerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineTriggerOptions := &cdtektonpipelinev2.GetTektonPipelineTriggerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineTriggerOptions.SetPipelineID(parts[0])
	getTektonPipelineTriggerOptions.SetTriggerID(parts[1])

	triggerIntf, response, err := cdTektonPipelineClient.GetTektonPipelineTriggerWithContext(context, getTektonPipelineTriggerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTektonPipelineTriggerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineTriggerWithContext failed %s\n%s", err, response))
	}

	trigger := triggerIntf.(*cdtektonpipelinev2.Trigger)
	if err = d.Set("pipeline_id", getTektonPipelineTriggerOptions.PipelineID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pipeline_id: %s", err))
	}
	if err = d.Set("type", trigger.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("name", trigger.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("event_listener", trigger.EventListener); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting event_listener: %s", err))
	}
	if !core.IsNil(trigger.Tags) {
		if err = d.Set("tags", trigger.Tags); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting tags: %s", err))
		}
	}
	if !core.IsNil(trigger.Worker) {
		workerMap, err := resourceIBMCdTektonPipelineTriggerWorkerIdentityToMap(trigger.Worker)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("worker", []map[string]interface{}{workerMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting worker: %s", err))
		}
	}
	if !core.IsNil(trigger.MaxConcurrentRuns) {
		if err = d.Set("max_concurrent_runs", flex.IntValue(trigger.MaxConcurrentRuns)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting max_concurrent_runs: %s", err))
		}
	}
	if !core.IsNil(trigger.Enabled) {
		if err = d.Set("enabled", trigger.Enabled); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
		}
	}
	if !core.IsNil(trigger.Secret) {
		secretMap, err := resourceIBMCdTektonPipelineTriggerGenericSecretToMap(trigger.Secret)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("secret", []map[string]interface{}{secretMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting secret: %s", err))
		}
	}
	if !core.IsNil(trigger.Cron) {
		if err = d.Set("cron", trigger.Cron); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cron: %s", err))
		}
	}
	if !core.IsNil(trigger.Timezone) {
		if err = d.Set("timezone", trigger.Timezone); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting timezone: %s", err))
		}
	}
	if !core.IsNil(trigger.Source) {
		sourceMap, err := resourceIBMCdTektonPipelineTriggerTriggerSourcePrototypeToMap(trigger.Source)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("source", []map[string]interface{}{sourceMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source: %s", err))
		}
	}
	if !core.IsNil(trigger.Events) {
		if err = d.Set("events", trigger.Events); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting events: %s", err))
		}
	}
	if !core.IsNil(trigger.Filter) {
		if err = d.Set("filter", trigger.Filter); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting filter: %s", err))
		}
	}
	if !core.IsNil(trigger.Favorite) {
		if err = d.Set("favorite", trigger.Favorite); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting favorite: %s", err))
		}
	}
	if !core.IsNil(trigger.EnableEventsFromForks) {
		if err = d.Set("enable_events_from_forks", trigger.EnableEventsFromForks); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enable_events_from_forks: %s", err))
		}
	}
	if !core.IsNil(trigger.Href) {
		if err = d.Set("href", trigger.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if !core.IsNil(trigger.Properties) {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range trigger.Properties {
			propertiesItemMap, err := resourceIBMCdTektonPipelineTriggerTriggerPropertyToMap(&propertiesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			properties = append(properties, propertiesItemMap)
		}
		if err = d.Set("properties", properties); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting properties: %s", err))
		}
	}
	if !core.IsNil(trigger.WebhookURL) {
		if err = d.Set("webhook_url", trigger.WebhookURL); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting webhook_url: %s", err))
		}
	}
	if !core.IsNil(trigger.ID) {
		if err = d.Set("trigger_id", trigger.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting trigger_id: %s", err))
		}
	}

	return nil
}

func resourceIBMCdTektonPipelineTriggerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateTektonPipelineTriggerOptions := &cdtektonpipelinev2.UpdateTektonPipelineTriggerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateTektonPipelineTriggerOptions.SetPipelineID(parts[0])
	updateTektonPipelineTriggerOptions.SetTriggerID(parts[1])

	hasChange := false

	patchVals := &cdtektonpipelinev2.TriggerPatch{}
	if d.HasChange("pipeline_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "pipeline_id"))
	}
	if d.HasChange("type") || d.HasChange("name") || d.HasChange("event_listener") {
		newType := d.Get("type").(string)
		patchVals.Type = &newType
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		newEventListener := d.Get("event_listener").(string)
		patchVals.EventListener = &newEventListener
		hasChange = true
	}
	if d.HasChange("tags") {
		tags := []string{}
		for _, tagsItem := range d.Get("tags").([]interface{}) {
			tags = append(tags, tagsItem.(string))
		}
		patchVals.Tags = tags
		hasChange = true
	}
	if d.HasChange("worker") {
		worker, err := resourceIBMCdTektonPipelineTriggerMapToWorkerIdentity(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
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
		secret, err := resourceIBMCdTektonPipelineTriggerMapToGenericSecret(d.Get("secret.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
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
		source, err := resourceIBMCdTektonPipelineTriggerMapToTriggerSourcePrototype(d.Get("source.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
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
		updateTektonPipelineTriggerOptions.TriggerPatch, _ = patchVals.AsPatch()
		_, response, err := cdTektonPipelineClient.UpdateTektonPipelineTriggerWithContext(context, updateTektonPipelineTriggerOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateTektonPipelineTriggerWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateTektonPipelineTriggerWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCdTektonPipelineTriggerRead(context, d, meta)
}

func resourceIBMCdTektonPipelineTriggerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelineTriggerOptions := &cdtektonpipelinev2.DeleteTektonPipelineTriggerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelineTriggerOptions.SetPipelineID(parts[0])
	deleteTektonPipelineTriggerOptions.SetTriggerID(parts[1])

	response, err := cdTektonPipelineClient.DeleteTektonPipelineTriggerWithContext(context, deleteTektonPipelineTriggerOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTektonPipelineTriggerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTektonPipelineTriggerWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCdTektonPipelineTriggerMapToWorkerIdentity(modelMap map[string]interface{}) (*cdtektonpipelinev2.WorkerIdentity, error) {
	model := &cdtektonpipelinev2.WorkerIdentity{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMCdTektonPipelineTriggerMapToGenericSecret(modelMap map[string]interface{}) (*cdtektonpipelinev2.GenericSecret, error) {
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

func resourceIBMCdTektonPipelineTriggerMapToTriggerSourcePrototype(modelMap map[string]interface{}) (*cdtektonpipelinev2.TriggerSourcePrototype, error) {
	model := &cdtektonpipelinev2.TriggerSourcePrototype{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	PropertiesModel, err := resourceIBMCdTektonPipelineTriggerMapToTriggerSourcePropertiesPrototype(modelMap["properties"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Properties = PropertiesModel
	return model, nil
}

func resourceIBMCdTektonPipelineTriggerMapToTriggerSourcePropertiesPrototype(modelMap map[string]interface{}) (*cdtektonpipelinev2.TriggerSourcePropertiesPrototype, error) {
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

func resourceIBMCdTektonPipelineTriggerWorkerIdentityToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerGenericSecretToMap(model *cdtektonpipelinev2.GenericSecret) (map[string]interface{}, error) {
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

func resourceIBMCdTektonPipelineTriggerTriggerSourcePrototypeToMap(model *cdtektonpipelinev2.TriggerSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	propertiesMap, err := resourceIBMCdTektonPipelineTriggerTriggerSourcePropertiesPrototypeToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerTriggerSourcePropertiesPrototypeToMap(model *cdtektonpipelinev2.TriggerSourceProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["url"] = model.URL
	if model.Branch != nil {
		modelMap["branch"] = model.Branch
	}
	if model.Pattern != nil {
		modelMap["pattern"] = model.Pattern
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineTriggerTriggerPropertyToMap(model *cdtektonpipelinev2.TriggerProperty) (map[string]interface{}, error) {
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
	if model.Locked != nil {
		modelMap["locked"] = model.Locked
	}
	return modelMap, nil
}
