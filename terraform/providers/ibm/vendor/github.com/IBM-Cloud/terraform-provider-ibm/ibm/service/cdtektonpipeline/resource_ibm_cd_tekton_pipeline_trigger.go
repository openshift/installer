// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtektonpipeline

import (
	"context"
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/google/go-cmp/cmp"
)

func ResourceIBMTektonPipelineTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMTektonPipelineTriggerCreate,
		ReadContext:   ResourceIBMTektonPipelineTriggerRead,
		UpdateContext: ResourceIBMTektonPipelineTriggerUpdate,
		DeleteContext: ResourceIBMTektonPipelineTriggerDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger", "pipeline_id"),
				Description:  "The tekton pipeline ID.",
			},
			"trigger": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Tekton pipeline trigger object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "name of the duplicated trigger.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger type.",
						},
						"event_listener": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event listener name.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "UUID.",
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
							Description: "Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this trigger uses default pipeline worker.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "worker name.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "worker type.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Id.",
									},
								},
							},
						},
						"concurrency": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Concurrency object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_concurrent_runs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Defines the maximum number of concurrent runs for this trigger.",
									},
								},
							},
						},
						"disabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "flag whether the trigger is disabled.",
						},
						"scm_source": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Scm source for git type tekton pipeline trigger.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Needed only for git trigger type. Repo URL that listening to.",
									},
									"branch": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Needed only for git trigger type. Branch name of the repo. Branch field doesn't coexist with pattern field.",
									},
									"pattern": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Needed only for git trigger type. Git branch or tag pattern to listen to. Please refer to https://github.com/micromatch/micromatch for pattern syntax.",
									},
									"blind_connection": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Needed only for git trigger type. Branch name of the repo.",
									},
									"hook_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Webhook ID.",
									},
								},
							},
						},
						"events": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Needed only for git trigger type. Events object defines the events this git trigger listening to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"push": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the trigger starts when tekton pipeline service receive a repo's 'push' git webhook event.",
									},
									"pull_request_closed": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the trigger starts when tekton pipeline service receive a repo pull reqeust's 'close' git webhook event.",
									},
									"pull_request": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the trigger starts when tekton pipeline service receive a repo pull reqeust's 'open' or 'update' git webhook event.",
									},
								},
							},
						},
						"service_instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "UUID.",
						},
						"cron": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Needed only for timer trigger type. Cron expression for timer trigger. Maximum frequency is every 5 minutes.",
						},
						"timezone": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Needed only for timer trigger type. Timezones for timer trigger.",
						},
						"secret": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Needed only for generic trigger type. Secret used to start generic trigger.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Secret type.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret value, not needed if secret type is \"internalValidation\".",
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											segs := []string{d.Get("pipeline_id").(string), d.Get("trigger.0.id").(string)}
											secret := strings.Join(segs, ".")
											mac := hmac.New(sha3.New512, []byte(secret))
											mac.Write([]byte(new))
											secureHmac := hex.EncodeToString(mac.Sum(nil))
											hasEnvChange := !cmp.Equal(strings.Join([]string{"hash", "SHA3-512", secureHmac}, ":"), old)
											if hasEnvChange {
												return false
											}
											return true
										},
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret location, not needed if secret type is \"internalValidation\".",
									},
									"key_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret name, not needed if type is \"internalValidation\".",
									},
									"algorithm": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Algorithm used for \"digestMatches\" secret type.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ResourceIBMTektonPipelineTriggerValidator() *validate.ResourceValidator {
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
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_tekton_pipeline_trigger", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMTektonPipelineTriggerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelineTriggerOptions := &cdtektonpipelinev2.CreateTektonPipelineTriggerOptions{}

	createTektonPipelineTriggerOptions.SetPipelineID(d.Get("pipeline_id").(string))
	if _, ok := d.GetOk("trigger"); ok {
		triggerModel, err := ResourceIBMTektonPipelineTriggerMapToTrigger(d.Get("trigger.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createTektonPipelineTriggerOptions.SetTrigger(triggerModel)
	}

	triggerIntf, response, err := cdTektonPipelineClient.CreateTektonPipelineTriggerWithContext(context, createTektonPipelineTriggerOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTektonPipelineTriggerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTektonPipelineTriggerWithContext failed %s\n%s", err, response))
	}

	trigger := triggerIntf.(*cdtektonpipelinev2.Trigger)
	d.SetId(fmt.Sprintf("%s/%s", *createTektonPipelineTriggerOptions.PipelineID, *trigger.ID))

	return ResourceIBMTektonPipelineTriggerRead(context, d, meta)
}

func ResourceIBMTektonPipelineTriggerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if err = d.Set("pipeline_id", getTektonPipelineTriggerOptions.PipelineID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pipeline_id: %s", err))
	}

	trigger, err := ResourceIBMTektonPipelineTriggerTriggerToMap(triggerIntf)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("trigger", []map[string]interface{}{trigger}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting trigger: %s", err))
	}

	return nil
}

func ResourceIBMTektonPipelineTriggerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.HasChange("trigger.0.name") {
		updateTektonPipelineTriggerOptions.SetName(d.Get("trigger.0.name").(string))
		hasChange = true
	}
	if d.HasChange("trigger.0.event_listener") {
		updateTektonPipelineTriggerOptions.SetEventListener(d.Get("trigger.0.event_listener").(string))
		hasChange = true
	}
	if d.HasChange("trigger.0.tags") {
		tags := []string{}
		for _, tagsItem := range d.Get("trigger.0.tags").([]interface{}) {
			tags = append(tags, tagsItem.(string))
		}
		updateTektonPipelineTriggerOptions.SetTags(tags)
		hasChange = true
	}
	if d.HasChange("trigger.0.worker") {
		worker, err := ResourceIBMTektonPipelineTriggerMapToWorker(d.Get("trigger.0.worker").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateTektonPipelineTriggerOptions.SetWorker(worker)
		hasChange = true
	}

	if d.HasChange("trigger.0.concurrency") {
		concurrency, err := ResourceIBMTektonPipelineTriggerMapToConcurrency(d.Get("trigger.0.concurrency").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateTektonPipelineTriggerOptions.SetConcurrency(concurrency)
		hasChange = true
	}

	if d.HasChange("trigger.0.secret") {
		secret, err := ResourceIBMTektonPipelineTriggerMapToGenericSecret(d.Get("trigger.0.secret").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateTektonPipelineTriggerOptions.SetSecret(secret)
		hasChange = true
	}

	if d.HasChange("trigger.0.scm_source") {
		secret, err := ResourceIBMTektonPipelineTriggerMapToTriggerScmSource(d.Get("trigger.0.scm_source").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateTektonPipelineTriggerOptions.SetScmSource(secret)
		hasChange = true
	}

	if d.HasChange("trigger.0.events") {
		events, err := ResourceIBMTektonPipelineTriggerMapToEvents(d.Get("trigger.0.events").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateTektonPipelineTriggerOptions.SetEvents(events)
		hasChange = true
	}

	if d.HasChange("trigger.0.cron") {
		updateTektonPipelineTriggerOptions.SetCron(d.Get("trigger.0.cron").(string))
		hasChange = true
	}
	if d.HasChange("trigger.0.timezone") {
		updateTektonPipelineTriggerOptions.SetTimezone(d.Get("trigger.0.timezone").(string))
		hasChange = true
	}
	if d.HasChange("trigger.0.disabled") {
		updateTektonPipelineTriggerOptions.SetDisabled(d.Get("trigger.0.disabled").(bool))
		hasChange = true
	}

	if hasChange {
		_, response, err := cdTektonPipelineClient.UpdateTektonPipelineTriggerWithContext(context, updateTektonPipelineTriggerOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateTektonPipelineTriggerWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateTektonPipelineTriggerWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMTektonPipelineTriggerRead(context, d, meta)
}

func ResourceIBMTektonPipelineTriggerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func ResourceIBMTektonPipelineTriggerMapToTrigger(modelMap map[string]interface{}) (cdtektonpipelinev2.TriggerIntf, error) {
	model := &cdtektonpipelinev2.Trigger{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	if modelMap["event_listener"] != nil && modelMap["event_listener"].(string) != "" {
		model.EventListener = core.StringPtr(modelMap["event_listener"].(string))
	}
	if modelMap["tags"] != nil {
		tags := []string{}
		for _, tagsItem := range modelMap["tags"].([]interface{}) {
			tags = append(tags, tagsItem.(string))
		}
		model.Tags = tags
	}
	if modelMap["worker"] != nil && len(modelMap["worker"].([]interface{})) > 0 {
		WorkerModel, err := ResourceIBMTektonPipelineTriggerMapToWorker(modelMap["worker"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Worker = WorkerModel
	}
	if modelMap["concurrency"] != nil && len(modelMap["concurrency"].([]interface{})) > 0 {
		ConcurrencyModel, err := ResourceIBMTektonPipelineTriggerMapToConcurrency(modelMap["concurrency"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Concurrency = ConcurrencyModel
	}
	if modelMap["disabled"] != nil {
		model.Disabled = core.BoolPtr(modelMap["disabled"].(bool))
	}
	if modelMap["scm_source"] != nil && len(modelMap["scm_source"].([]interface{})) > 0 {
		ScmSourceModel, err := ResourceIBMTektonPipelineTriggerMapToTriggerScmSource(modelMap["scm_source"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ScmSource = ScmSourceModel
	}
	if modelMap["events"] != nil && len(modelMap["events"].([]interface{})) > 0 {
		EventsModel, err := ResourceIBMTektonPipelineTriggerMapToEvents(modelMap["events"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Events = EventsModel
	}
	if modelMap["service_instance_id"] != nil && modelMap["service_instance_id"].(string) != "" {
		model.ServiceInstanceID = core.StringPtr(modelMap["service_instance_id"].(string))
	}
	if modelMap["cron"] != nil && modelMap["cron"].(string) != "" {
		model.Cron = core.StringPtr(modelMap["cron"].(string))
	}
	if modelMap["timezone"] != nil && modelMap["timezone"].(string) != "" {
		model.Timezone = core.StringPtr(modelMap["timezone"].(string))
	}
	if modelMap["secret"] != nil && len(modelMap["secret"].([]interface{})) > 0 {
		SecretModel, err := ResourceIBMTektonPipelineTriggerMapToGenericSecret(modelMap["secret"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Secret = SecretModel
	}
	return model, nil
}

func ResourceIBMTektonPipelineTriggerMapToWorker(modelMap map[string]interface{}) (*cdtektonpipelinev2.Worker, error) {
	model := &cdtektonpipelinev2.Worker{}
	//TODO double check if working
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMTektonPipelineTriggerMapToConcurrency(modelMap map[string]interface{}) (*cdtektonpipelinev2.Concurrency, error) {
	model := &cdtektonpipelinev2.Concurrency{}
	if modelMap["max_concurrent_runs"] != nil {
		model.MaxConcurrentRuns = core.Int64Ptr(int64(modelMap["max_concurrent_runs"].(int)))
	}
	return model, nil
}

func ResourceIBMTektonPipelineTriggerMapToTriggerScmSource(modelMap map[string]interface{}) (*cdtektonpipelinev2.TriggerScmSource, error) {
	model := &cdtektonpipelinev2.TriggerScmSource{}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["branch"] != nil && modelMap["branch"].(string) != "" {
		model.Branch = core.StringPtr(modelMap["branch"].(string))
	}
	if modelMap["pattern"] != nil && modelMap["pattern"].(string) != "" {
		model.Pattern = core.StringPtr(modelMap["pattern"].(string))
	}
	return model, nil
}

func ResourceIBMTektonPipelineTriggerMapToEvents(modelMap map[string]interface{}) (*cdtektonpipelinev2.Events, error) {
	model := &cdtektonpipelinev2.Events{}
	if modelMap["push"] != nil {
		model.Push = core.BoolPtr(modelMap["push"].(bool))
	}
	if modelMap["pull_request_closed"] != nil {
		model.PullRequestClosed = core.BoolPtr(modelMap["pull_request_closed"].(bool))
	}
	if modelMap["pull_request"] != nil {
		model.PullRequest = core.BoolPtr(modelMap["pull_request"].(bool))
	}
	return model, nil
}

func ResourceIBMTektonPipelineTriggerMapToGenericSecret(modelMap map[string]interface{}) (*cdtektonpipelinev2.GenericSecret, error) {
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

func ResourceIBMTektonPipelineTriggerTriggerToMap(model cdtektonpipelinev2.TriggerIntf) (map[string]interface{}, error) {
	if _, ok := model.(*cdtektonpipelinev2.TriggerManualTrigger); ok {
		return ResourceIBMTektonPipelineTriggerTriggerManualTriggerToMap(model.(*cdtektonpipelinev2.TriggerManualTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerScmTrigger); ok {
		return ResourceIBMTektonPipelineTriggerTriggerScmTriggerToMap(model.(*cdtektonpipelinev2.TriggerScmTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerTimerTrigger); ok {
		return ResourceIBMTektonPipelineTriggerTriggerTimerTriggerToMap(model.(*cdtektonpipelinev2.TriggerTimerTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.TriggerGenericTrigger); ok {
		return ResourceIBMTektonPipelineTriggerTriggerGenericTriggerToMap(model.(*cdtektonpipelinev2.TriggerGenericTrigger))
	} else if _, ok := model.(*cdtektonpipelinev2.Trigger); ok {
		modelMap := make(map[string]interface{})
		model := model.(*cdtektonpipelinev2.Trigger)
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.Type != nil {
			modelMap["type"] = model.Type
		}
		if model.EventListener != nil {
			modelMap["event_listener"] = model.EventListener
		}
		if model.ID != nil {
			modelMap["id"] = model.ID
		}
		if model.Tags != nil {
			modelMap["tags"] = model.Tags
		}
		if model.Worker != nil {
			workerMap, err := ResourceIBMTektonPipelineTriggerWorkerToMap(model.Worker)
			if err != nil {
				return modelMap, err
			}
			modelMap["worker"] = []map[string]interface{}{workerMap}
		}
		if model.Concurrency != nil {
			concurrencyMap, err := ResourceIBMTektonPipelineTriggerConcurrencyToMap(model.Concurrency)
			if err != nil {
				return modelMap, err
			}
			modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
		}
		if model.Disabled != nil {
			modelMap["disabled"] = model.Disabled
		}
		if model.ScmSource != nil {
			scmSourceMap, err := ResourceIBMTektonPipelineTriggerTriggerScmSourceToMap(model.ScmSource)
			if err != nil {
				return modelMap, err
			}
			modelMap["scm_source"] = []map[string]interface{}{scmSourceMap}
		}
		if model.Events != nil {
			eventsMap, err := ResourceIBMTektonPipelineTriggerEventsToMap(model.Events)
			if err != nil {
				return modelMap, err
			}
			modelMap["events"] = []map[string]interface{}{eventsMap}
		}
		if model.ServiceInstanceID != nil {
			modelMap["service_instance_id"] = model.ServiceInstanceID
		}
		if model.Cron != nil {
			modelMap["cron"] = model.Cron
		}
		if model.Timezone != nil {
			modelMap["timezone"] = model.Timezone
		}
		if model.Secret != nil {
			secretMap, err := ResourceIBMTektonPipelineTriggerGenericSecretToMap(model.Secret)
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

func ResourceIBMTektonPipelineTriggerWorkerToMap(model *cdtektonpipelinev2.Worker) (map[string]interface{}, error) {
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

func ResourceIBMTektonPipelineTriggerConcurrencyToMap(model *cdtektonpipelinev2.Concurrency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerTriggerScmSourceToMap(model *cdtektonpipelinev2.TriggerScmSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.Branch != nil {
		modelMap["branch"] = model.Branch
	}
	if model.Pattern != nil {
		modelMap["pattern"] = model.Pattern
	}
	if model.BlindConnection != nil {
		modelMap["blind_connection"] = model.BlindConnection
	}
	if model.HookID != nil {
		modelMap["hook_id"] = model.HookID
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerEventsToMap(model *cdtektonpipelinev2.Events) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Push != nil {
		modelMap["push"] = model.Push
	}
	if model.PullRequestClosed != nil {
		modelMap["pull_request_closed"] = model.PullRequestClosed
	}
	if model.PullRequest != nil {
		modelMap["pull_request"] = model.PullRequest
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerGenericSecretToMap(model *cdtektonpipelinev2.GenericSecret) (map[string]interface{}, error) {
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

func ResourceIBMTektonPipelineTriggerTriggerManualTriggerToMap(model *cdtektonpipelinev2.TriggerManualTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["event_listener"] = model.EventListener
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineTriggerWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := ResourceIBMTektonPipelineTriggerConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	modelMap["disabled"] = model.Disabled
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerTriggerScmTriggerToMap(model *cdtektonpipelinev2.TriggerScmTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["event_listener"] = model.EventListener
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineTriggerWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := ResourceIBMTektonPipelineTriggerConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	modelMap["disabled"] = model.Disabled
	if model.ScmSource != nil {
		scmSourceMap, err := ResourceIBMTektonPipelineTriggerTriggerScmSourceToMap(model.ScmSource)
		if err != nil {
			return modelMap, err
		}
		modelMap["scm_source"] = []map[string]interface{}{scmSourceMap}
	}
	if model.Events != nil {
		eventsMap, err := ResourceIBMTektonPipelineTriggerEventsToMap(model.Events)
		if err != nil {
			return modelMap, err
		}
		modelMap["events"] = []map[string]interface{}{eventsMap}
	}
	if model.ServiceInstanceID != nil {
		modelMap["service_instance_id"] = model.ServiceInstanceID
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerTriggerTimerTriggerToMap(model *cdtektonpipelinev2.TriggerTimerTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["event_listener"] = model.EventListener
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineTriggerWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := ResourceIBMTektonPipelineTriggerConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	modelMap["disabled"] = model.Disabled
	if model.Cron != nil {
		modelMap["cron"] = model.Cron
	}
	if model.Timezone != nil {
		modelMap["timezone"] = model.Timezone
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerTriggerGenericTriggerToMap(model *cdtektonpipelinev2.TriggerGenericTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["event_listener"] = model.EventListener
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineTriggerWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := ResourceIBMTektonPipelineTriggerConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	modelMap["disabled"] = model.Disabled
	if model.Secret != nil {
		secretMap, err := ResourceIBMTektonPipelineTriggerGenericSecretToMap(model.Secret)
		if err != nil {
			return modelMap, err
		}
		modelMap["secret"] = []map[string]interface{}{secretMap}
	}
	return modelMap, nil
}
