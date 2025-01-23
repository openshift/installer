// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package atracker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/atrackerv2"
)

const COS_CRN_PARTS = 8

func ResourceIBMAtrackerTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAtrackerTargetCreate,
		ReadContext:   resourceIBMAtrackerTargetRead,
		UpdateContext: resourceIBMAtrackerTargetUpdate,
		DeleteContext: resourceIBMAtrackerTargetDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_atracker_target", "name"),
				Description:  "The name of the target. The name must be 1000 characters or less, and cannot include any special characters other than `(space) - . _ :`.",
			},
			"target_type": {
				Type:             schema.TypeString,
				DiffSuppressFunc: flex.ApplyOnce,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validate.InvokeValidator("ibm_atracker_target", "target_type"),
				Description:      "The type of the target. It can be cloud_object_storage, logdna, event_streams, or cloud_logs. Based on this type you must include cos_endpoint, logdna_endpoint, eventstreams_endpoint or cloudlogs_endpoint.",
			},
			"cos_endpoint": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Property values for a Cloud Object Storage Endpoint.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The host name of the Cloud Object Storage endpoint.",
						},
						"target_crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN of the Cloud Object Storage instance.",
						},
						"bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The bucket name under the Cloud Object Storage instance.",
						},
						"api_key": {
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							Description:      "The IAM API key that has writer access to the Cloud Object Storage instance. This credential is masked in the response. This is required if service_to_service is not enabled.",
							DiffSuppressFunc: flex.ApplyOnce,
						},
						"service_to_service_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Determines if IBM Cloud Activity Tracker Event Routing has service to service authentication enabled. Set this flag to true if service to service is enabled and do not supply an apikey.",
						},
					},
				},
			},
			"logdna_endpoint": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Property values for a LogDNA Endpoint.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN of the LogDNA instance.",
						},
						"ingestion_key": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							DiffSuppressFunc: flex.ApplyOnce,
							Description:      "The LogDNA ingestion key is used for routing logs to a specific LogDNA instance.",
						},
					},
				},
			},
			"eventstreams_endpoint": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Property values for an Event Streams Endpoint in requests.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN of the Event Streams instance.",
						},
						"brokers": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of broker endpoints.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"topic": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The messsage hub topic defined in the Event Streams instance.",
						},
						"api_key": &schema.Schema{ // pragma: allowlist secret
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							DiffSuppressFunc: flex.ApplyOnce,
							Description:      "The user password (api key) for the message hub topic in the Event Streams instance. This is required if service_to_service is not enabled.",
						},
						"service_to_service_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Determines if IBM Cloud Activity Tracker Event Routing has service to service authentication enabled. Set this flag to true if service to service is enabled and do not supply an apikey.",
						},
					},
				},
			},
			"cloudlogs_endpoint": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Property values for an IBM Cloud Logs endpoint in requests.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN of the IBM Cloud Logs instance.",
						},
					},
				},
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_atracker_target", "region"),
				Description:  "Include this optional field if you want to create a target in a different region other than the one you are connected.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the target resource.",
			},
			"encrypt_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "use encryption_key instead",
				Description: "The encryption key that is used to encrypt events before Activity Tracker services buffer them on storage. This credential is masked in the response.",
			},
			"encryption_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The encryption key that is used to encrypt events before Activity Tracker services buffer them on storage. This credential is masked in the response.",
			},
			"cos_write_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Deprecated:  "use write_status instead",
				Description: "The status of the write attempt with the provided cos_endpoint parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The status such as failed or success.",
						},
						"last_failure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The timestamp of the failure.",
						},
						"reason_for_last_failure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detailed description of the cause of the failure.",
						},
					},
				},
			},
			"write_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The status of the write attempt to the target with the provided endpoint parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The status such as failed or success.",
						},
						"last_failure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The timestamp of the failure.",
						},
						"reason_for_last_failure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detailed description of the cause of the failure.",
						},
					},
				},
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "use created_at instead",
				Description: "The timestamp of the target creation time.",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "use updated_at instead",
				Description: "The timestamp of the target last updated time.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the target creation time.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp of the target last updated time.",
			},
			"api_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The API version of the target.",
			},
		},
	}
}

func ResourceIBMAtrackerTargetValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9 -._:]+$`,
			MinValueLength:             1,
			MaxValueLength:             1000,
		},
		validate.ValidateSchema{
			Identifier:                 "target_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "cloud_object_storage, logdna, event_streams, cloud_logs",
		},
		validate.ValidateSchema{
			Identifier:                 "region",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 -._:]+$`,
			MinValueLength:             3,
			MaxValueLength:             1000,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_atracker_target", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMAtrackerTargetCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createTargetOptions := &atrackerv2.CreateTargetOptions{}

	createTargetOptions.SetName(d.Get("name").(string))
	createTargetOptions.SetTargetType(d.Get("target_type").(string))
	if _, ok := d.GetOk("cos_endpoint"); ok {
		cosEndpointModel, err := resourceIBMAtrackerTargetMapToCosEndpointPrototype(d.Get("cos_endpoint.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createTargetOptions.SetCosEndpoint(cosEndpointModel)
	}
	if _, ok := d.GetOk("logdna_endpoint"); ok {
		logdnaEndpointModel, err := resourceIBMAtrackerTargetMapToLogdnaEndpointPrototype(d.Get("logdna_endpoint.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createTargetOptions.SetLogdnaEndpoint(logdnaEndpointModel)
	}
	if _, ok := d.GetOk("eventstreams_endpoint"); ok {
		eventstreamsEndpointModel, err := resourceIBMAtrackerTargetMapToEventstreamsEndpointPrototype(d.Get("eventstreams_endpoint.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createTargetOptions.SetEventstreamsEndpoint(eventstreamsEndpointModel)
	}
	if _, ok := d.GetOk("cloudlogs_endpoint"); ok {
		cloudLogsEndpointModel, err := resourceIBMAtrackerTargetMapToCloudLogsEndpointPrototype(d.Get("cloudlogs_endpoint.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createTargetOptions.SetCloudlogsEndpoint(cloudLogsEndpointModel)
	}
	if _, ok := d.GetOk("region"); ok {
		createTargetOptions.SetRegion(d.Get("region").(string))
	}

	target, response, err := atrackerClient.CreateTargetWithContext(context, createTargetOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTargetWithContext failed %s\n%s", err, response))
	}

	d.SetId(*target.ID)

	return resourceIBMAtrackerTargetRead(context, d, meta)
}

func resourceIBMAtrackerTargetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getTargetOptions := &atrackerv2.GetTargetOptions{}

	getTargetOptions.SetID(d.Id())

	target, response, err := atrackerClient.GetTargetWithContext(context, getTargetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTargetWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("name", target.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("target_type", target.TargetType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting target_type: %s", err))
	}
	// Don't report difference if the last parts of CRN are different
	if target.CosEndpoint != nil {
		cosEndpointMap, err := resourceIBMAtrackerTargetCosEndpointPrototypeToMap(target.CosEndpoint)
		if cosInterface, ok := d.GetOk("cos_endpoint.0"); ok {
			targetCrnExisting := cosInterface.(map[string]interface{})["target_crn"].(string)
			targetCrnIncoming := cosEndpointMap["target_crn"].(*string)
			if len(targetCrnExisting) > 0 && targetCrnIncoming != nil {
				targetCrnExistingParts := strings.Split(targetCrnExisting, ":")
				targetCrnIncomingParts := strings.Split(*targetCrnIncoming, ":")
				isDifferent := false
				for i := 0; i < COS_CRN_PARTS && len(targetCrnExistingParts) > COS_CRN_PARTS-1 && len(targetCrnIncomingParts) > COS_CRN_PARTS-1; i++ {
					if targetCrnExistingParts[i] != targetCrnIncomingParts[i] {
						isDifferent = true
					}
				}
				if !isDifferent {
					cosEndpointMap["target_crn"] = targetCrnExisting
				}
			}
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("cos_endpoint", []map[string]interface{}{cosEndpointMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cos_endpoint: %s", err))
		}
	}
	if target.LogdnaEndpoint != nil {
		logdnaEndpointMap, err := resourceIBMAtrackerTargetLogdnaEndpointPrototypeToMap(target.LogdnaEndpoint)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("logdna_endpoint", []map[string]interface{}{logdnaEndpointMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting logdna_endpoint: %s", err))
		}
	}
	if target.EventstreamsEndpoint != nil {
		eventstreamsEndpointMap, err := resourceIBMAtrackerTargetEventstreamsEndpointPrototypeToMap(target.EventstreamsEndpoint)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("eventstreams_endpoint", []map[string]interface{}{eventstreamsEndpointMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting eventstreams_endpoint: %s", err))
		}
	}

	if target.CloudlogsEndpoint != nil {
		cloudLogsEndpointMap, err := resourceIBMAtrackerTargetCloudLogsEndpointPrototypeToMap(target.CloudlogsEndpoint)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("cloudlogs_endpoint", []map[string]interface{}{cloudLogsEndpointMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cloudlogs_endpoint: %s", err))
		}
	}

	if target.CRN != nil {
		if err = d.Set("crn", target.CRN); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
		}
	}

	if _, exists := d.GetOk("region"); exists {
		if target.Region != nil && len(*target.Region) > 0 {
			d.Set("region", *target.Region)
			if err = d.Set("region", *target.Region); err != nil {
				return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
			}
		}
	}

	writeStatusMap, err := resourceIBMAtrackerTargetWriteStatusToMap(target.WriteStatus)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("write_status", []map[string]interface{}{writeStatusMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting write_status: %s", err))
	}

	// TODO: will be removed
	if err = d.Set("cos_write_status", []map[string]interface{}{writeStatusMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cos_write_status: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(target.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(target.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("api_version", flex.IntValue(target.APIVersion)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting api_version: %s", err))
	}

	return nil
}

func resourceIBMAtrackerTargetUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTargetOptions := &atrackerv2.ReplaceTargetOptions{}

	replaceTargetOptions.SetID(d.Id())

	hasChange := false

	if d.HasChange("name") || d.HasChange("cos_endpoint") || d.HasChange("region") || d.HasChange("logdna_endpoint") || d.HasChange("eventstreams_endpoint") || d.HasChange("cloudlogs_endpoint") {
		replaceTargetOptions.SetName(d.Get("name").(string))

		_, hasCosEndpoint := d.GetOk("cos_endpoint.0")
		if hasCosEndpoint {
			cosEndpoint, err := resourceIBMAtrackerTargetMapToCosEndpointPrototype(d.Get("cos_endpoint.0").(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			replaceTargetOptions.SetCosEndpoint(cosEndpoint)
		}

		_, hasLogDNAEndpoint := d.GetOk("logdna_endpoint.0")
		if hasLogDNAEndpoint {
			logdnaEndpoint, err := resourceIBMAtrackerTargetMapToLogdnaEndpointPrototype(d.Get("logdna_endpoint.0").(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			replaceTargetOptions.SetLogdnaEndpoint(logdnaEndpoint)
		}
		_, hasEventstreamsEndpoint := d.GetOk("eventstreams_endpoint.0")
		if hasEventstreamsEndpoint {
			eventstreamsEndpoint, err := resourceIBMAtrackerTargetMapToEventstreamsEndpointPrototype(d.Get("eventstreams_endpoint.0").(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			replaceTargetOptions.SetEventstreamsEndpoint(eventstreamsEndpoint)
		}
		_, hasCloudLogsEndpoint := d.GetOk("cloudlogs_endpoint.0")
		if hasCloudLogsEndpoint {
			cloudlogsEndpoint, err := resourceIBMAtrackerTargetMapToCloudLogsEndpointPrototype(d.Get("cloudlogs_endpoint.0").(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			replaceTargetOptions.SetCloudlogsEndpoint(cloudlogsEndpoint)
		}

		hasChange = true
	}

	if hasChange {
		_, response, err := atrackerClient.ReplaceTargetWithContext(context, replaceTargetOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceTargetWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ReplaceTargetWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMAtrackerTargetRead(context, d, meta)
}

func resourceIBMAtrackerTargetDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	atrackerClient, err := getAtrackerClients(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTargetOptions := &atrackerv2.DeleteTargetOptions{}

	deleteTargetOptions.SetID(d.Id())

	_, response, err := atrackerClient.DeleteTargetWithContext(context, deleteTargetOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTargetWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMAtrackerTargetMapToCosEndpointPrototype(modelMap map[string]interface{}) (*atrackerv2.CosEndpointPrototype, error) {
	model := &atrackerv2.CosEndpointPrototype{}
	model.Endpoint = core.StringPtr(modelMap["endpoint"].(string))
	model.TargetCRN = core.StringPtr(modelMap["target_crn"].(string))
	model.Bucket = core.StringPtr(modelMap["bucket"].(string))
	if modelMap["api_key"] != nil && modelMap["api_key"].(string) != "" {
		model.APIKey = core.StringPtr(modelMap["api_key"].(string))
	}
	model.ServiceToServiceEnabled = core.BoolPtr(modelMap["service_to_service_enabled"].(bool))
	return model, nil
}

func resourceIBMAtrackerTargetMapToLogdnaEndpointPrototype(modelMap map[string]interface{}) (*atrackerv2.LogdnaEndpointPrototype, error) {
	model := &atrackerv2.LogdnaEndpointPrototype{}
	model.TargetCRN = core.StringPtr(modelMap["target_crn"].(string))
	model.IngestionKey = core.StringPtr(modelMap["ingestion_key"].(string)) // pragma: whitelist secret
	return model, nil
}

func resourceIBMAtrackerTargetMapToEventstreamsEndpointPrototype(modelMap map[string]interface{}) (*atrackerv2.EventstreamsEndpointPrototype, error) {
	model := &atrackerv2.EventstreamsEndpointPrototype{}
	model.TargetCRN = core.StringPtr(modelMap["target_crn"].(string))
	model.Topic = core.StringPtr(modelMap["topic"].(string))
	brokers := []string{}
	for _, brokersItem := range modelMap["brokers"].([]interface{}) {
		brokers = append(brokers, brokersItem.(string))
	}
	model.Brokers = brokers
	if modelMap["api_key"] != nil && modelMap["api_key"].(string) != "" {
		model.APIKey = core.StringPtr(modelMap["api_key"].(string)) // pragma: whitelist secret
	}
	model.ServiceToServiceEnabled = core.BoolPtr(modelMap["service_to_service_enabled"].(bool))
	return model, nil
}

func resourceIBMAtrackerTargetMapToCloudLogsEndpointPrototype(modelMap map[string]interface{}) (*atrackerv2.CloudLogsEndpointPrototype, error) {
	model := &atrackerv2.CloudLogsEndpointPrototype{}
	model.TargetCRN = core.StringPtr(modelMap["target_crn"].(string))
	return model, nil
}

func resourceIBMAtrackerTargetCosEndpointPrototypeToMap(model *atrackerv2.CosEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["endpoint"] = model.Endpoint
	modelMap["target_crn"] = model.TargetCRN
	modelMap["bucket"] = model.Bucket
	// TODO: remove after deprecation
	modelMap["api_key"] = REDACTED_TEXT // pragma: whitelist secret
	modelMap["service_to_service_enabled"] = model.ServiceToServiceEnabled
	return modelMap, nil
}

func resourceIBMAtrackerTargetLogdnaEndpointPrototypeToMap(model *atrackerv2.LogdnaEndpoint) (map[string]interface{}, error) {

	modelMap := make(map[string]interface{})
	modelMap["target_crn"] = model.TargetCRN
	modelMap["ingestion_key"] = REDACTED_TEXT // pragma: whitelist secret
	return modelMap, nil
}

func resourceIBMAtrackerTargetEventstreamsEndpointPrototypeToMap(model *atrackerv2.EventstreamsEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_crn"] = model.TargetCRN
	modelMap["brokers"] = model.Brokers
	modelMap["topic"] = model.Topic
	// TODO: remove after deprecation
	modelMap["api_key"] = REDACTED_TEXT // pragma: whitelist secret
	modelMap["service_to_service_enabled"] = model.ServiceToServiceEnabled
	return modelMap, nil
}

func resourceIBMAtrackerTargetCloudLogsEndpointPrototypeToMap(model *atrackerv2.CloudLogsEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_crn"] = model.TargetCRN
	return modelMap, nil
}

func resourceIBMAtrackerTargetWriteStatusToMap(model *atrackerv2.WriteStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["status"] = model.Status
	if model.LastFailure != nil {
		modelMap["last_failure"] = model.LastFailure.String()
	}
	if model.ReasonForLastFailure != nil {
		modelMap["reason_for_last_failure"] = model.ReasonForLastFailure
	}
	return modelMap, nil
}
