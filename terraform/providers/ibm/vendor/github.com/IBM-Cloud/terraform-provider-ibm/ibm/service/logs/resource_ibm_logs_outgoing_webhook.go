// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsOutgoingWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsOutgoingWebhookCreate,
		ReadContext:   resourceIbmLogsOutgoingWebhookRead,
		UpdateContext: resourceIbmLogsOutgoingWebhookUpdate,
		DeleteContext: resourceIbmLogsOutgoingWebhookDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_outgoing_webhook", "type"),
				Description:  "The type of the deployed Outbound Integrations to list.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_outgoing_webhook", "name"),
				Description:  "The name of the Outbound Integration.",
			},
			"url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_outgoing_webhook", "url"),
				Description:  "The URL of the Outbound Integration. Null for IBM Event Notifications integration.",
			},
			"ibm_event_notifications": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The configuration of the IBM Event Notifications Outbound Integration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_notifications_instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the selected IBM Event Notifications instance.",
						},
						"region_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The region ID of the selected IBM Event Notifications instance.",
						},
						"endpoint_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The endpoint type of integration.",
						},
						"source_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the created source in the IBM Event Notifications instance. Corresponds to the Cloud Logs instance crn. Not required when creating an Outbound Integration.",
						},
						"source_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the created source in the IBM Event Notifications instance. Not required when creating an Outbound Integration.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the Outbound Integration.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the Outbound Integration.",
			},
			"external_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The external ID of the Outbound Integration, for connecting with other parts of the system.",
			},
			"webhook_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Outgoing Webhook Id.",
			},
		},
	}
}

func ResourceIbmLogsOutgoingWebhookValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "ibm_event_notifications",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "url",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_outgoing_webhook", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsOutgoingWebhookCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_outgoing_webhook", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	bodyModelMap := map[string]interface{}{}
	createOutgoingWebhookOptions := &logsv0.CreateOutgoingWebhookOptions{}

	bodyModelMap["type"] = d.Get("type")
	bodyModelMap["name"] = d.Get("name")
	if _, ok := d.GetOk("url"); ok {
		bodyModelMap["url"] = d.Get("url")
	}
	if _, ok := d.GetOk("ibm_event_notifications"); ok {
		bodyModelMap["ibm_event_notifications"] = d.Get("ibm_event_notifications")
	}
	convertedModel, err := ResourceIbmLogsOutgoingWebhookMapToOutgoingWebhookPrototype(bodyModelMap)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_outgoing_webhook", "create")
		return tfErr.GetDiag()
	}
	createOutgoingWebhookOptions.OutgoingWebhookPrototype = convertedModel

	outgoingWebhookIntf, _, err := logsClient.CreateOutgoingWebhookWithContext(context, createOutgoingWebhookOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateOutgoingWebhookWithContext failed: %s", err.Error()), "ibm_logs_outgoing_webhook", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	outgoingWebhook := outgoingWebhookIntf.(*logsv0.OutgoingWebhook)

	webhookID := fmt.Sprintf("%s/%s/%s", region, instanceId, *outgoingWebhook.ID)
	d.SetId(webhookID)

	return resourceIbmLogsOutgoingWebhookRead(context, d, meta)
}

func resourceIbmLogsOutgoingWebhookRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_outgoing_webhook", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, region, instanceId, webhookId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	getOutgoingWebhookOptions := &logsv0.GetOutgoingWebhookOptions{}

	getOutgoingWebhookOptions.SetID(core.UUIDPtr(strfmt.UUID(webhookId)))

	outgoingWebhookIntf, response, err := logsClient.GetOutgoingWebhookWithContext(context, getOutgoingWebhookOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetOutgoingWebhookWithContext failed: %s", err.Error()), "ibm_logs_outgoing_webhook", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	outgoingWebhook := outgoingWebhookIntf.(*logsv0.OutgoingWebhook)

	if err = d.Set("webhook_id", webhookId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting webhook_id: %s", err))
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}
	if err = d.Set("type", outgoingWebhook.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("name", outgoingWebhook.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(outgoingWebhook.URL) {
		if err = d.Set("url", outgoingWebhook.URL); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting url: %s", err))
		}
	}
	if !core.IsNil(outgoingWebhook.IbmEventNotifications) {
		ibmEventNotificationsMap, err := ResourceIbmLogsOutgoingWebhookOutgoingWebhooksV1IbmEventNotificationsConfigToMap(outgoingWebhook.IbmEventNotifications)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("ibm_event_notifications", []map[string]interface{}{ibmEventNotificationsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting ibm_event_notifications: %s", err))
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(outgoingWebhook.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(outgoingWebhook.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("external_id", flex.IntValue(outgoingWebhook.ExternalID)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting external_id: %s", err))
	}

	return nil
}

func resourceIbmLogsOutgoingWebhookUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_outgoing_webhook", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, webhookId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateOutgoingWebhookOptions := &logsv0.UpdateOutgoingWebhookOptions{}

	updateOutgoingWebhookOptions.SetID(core.UUIDPtr(strfmt.UUID(webhookId)))

	hasChange := false

	if d.HasChange("type") ||
		d.HasChange("name") ||
		d.HasChange("url") ||
		d.HasChange("ibm_event_notifications") {

		bodyModelMap := map[string]interface{}{}

		bodyModelMap["type"] = d.Get("type")
		bodyModelMap["name"] = d.Get("name")

		if _, ok := d.GetOk("url"); ok {
			bodyModelMap["url"] = d.Get("url")
		}
		if _, ok := d.GetOk("ibm_event_notifications"); ok {
			bodyModelMap["ibm_event_notifications"] = d.Get("ibm_event_notifications")
		}
		convertedModel, err := ResourceIbmLogsOutgoingWebhookMapToOutgoingWebhookPrototype(bodyModelMap)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_outgoing_webhook", "create")
			return tfErr.GetDiag()
		}
		updateOutgoingWebhookOptions.OutgoingWebhookPrototype = convertedModel

		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.UpdateOutgoingWebhookWithContext(context, updateOutgoingWebhookOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateOutgoingWebhookWithContext failed: %s", err.Error()), "ibm_logs_outgoing_webhook", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsOutgoingWebhookRead(context, d, meta)
}

func resourceIbmLogsOutgoingWebhookDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_outgoing_webhook", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	logsClient, _, _, webhookId, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteOutgoingWebhookOptions := &logsv0.DeleteOutgoingWebhookOptions{}

	deleteOutgoingWebhookOptions.SetID(core.UUIDPtr(strfmt.UUID(webhookId)))

	_, err = logsClient.DeleteOutgoingWebhookWithContext(context, deleteOutgoingWebhookOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteOutgoingWebhookWithContext failed: %s", err.Error()), "ibm_logs_outgoing_webhook", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsOutgoingWebhookMapToOutgoingWebhooksV1IbmEventNotificationsConfig(modelMap map[string]interface{}) (*logsv0.OutgoingWebhooksV1IbmEventNotificationsConfig, error) {
	model := &logsv0.OutgoingWebhooksV1IbmEventNotificationsConfig{}
	model.EventNotificationsInstanceID = core.UUIDPtr(strfmt.UUID(modelMap["event_notifications_instance_id"].(string)))
	model.RegionID = core.StringPtr(modelMap["region_id"].(string))
	if modelMap["endpoint_type"] != nil && modelMap["endpoint_type"].(string) != "" {
		model.EndpointType = core.StringPtr(modelMap["endpoint_type"].(string))
	}
	if modelMap["source_id"] != nil && modelMap["source_id"].(string) != "" {
		model.SourceID = core.StringPtr(modelMap["source_id"].(string))
	}
	if modelMap["source_name"] != nil && modelMap["source_name"].(string) != "" {
		model.SourceName = core.StringPtr(modelMap["source_name"].(string))
	}
	return model, nil
}

func ResourceIbmLogsOutgoingWebhookMapToOutgoingWebhookPrototype(modelMap map[string]interface{}) (logsv0.OutgoingWebhookPrototypeIntf, error) {
	model := &logsv0.OutgoingWebhookPrototype{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["ibm_event_notifications"] != nil && len(modelMap["ibm_event_notifications"].([]interface{})) > 0 {
		IbmEventNotificationsModel, err := ResourceIbmLogsOutgoingWebhookMapToOutgoingWebhooksV1IbmEventNotificationsConfig(modelMap["ibm_event_notifications"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.IbmEventNotifications = IbmEventNotificationsModel
	}
	return model, nil
}

func ResourceIbmLogsOutgoingWebhookMapToOutgoingWebhookPrototypeOutgoingWebhooksV1OutgoingWebhookInputDataConfigIbmEventNotifications(modelMap map[string]interface{}) (*logsv0.OutgoingWebhookPrototypeOutgoingWebhooksV1OutgoingWebhookInputDataConfigIbmEventNotifications, error) {
	model := &logsv0.OutgoingWebhookPrototypeOutgoingWebhooksV1OutgoingWebhookInputDataConfigIbmEventNotifications{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["ibm_event_notifications"] != nil && len(modelMap["ibm_event_notifications"].([]interface{})) > 0 {
		IbmEventNotificationsModel, err := ResourceIbmLogsOutgoingWebhookMapToOutgoingWebhooksV1IbmEventNotificationsConfig(modelMap["ibm_event_notifications"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.IbmEventNotifications = IbmEventNotificationsModel
	}
	return model, nil
}

func ResourceIbmLogsOutgoingWebhookOutgoingWebhooksV1IbmEventNotificationsConfigToMap(model *logsv0.OutgoingWebhooksV1IbmEventNotificationsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["event_notifications_instance_id"] = model.EventNotificationsInstanceID.String()
	modelMap["region_id"] = *model.RegionID
	if model.EndpointType != nil {
		modelMap["endpoint_type"] = *model.EndpointType
	}
	if model.SourceID != nil {
		modelMap["source_id"] = *model.SourceID
	}
	if model.SourceName != nil {
		modelMap["source_name"] = *model.SourceName
	}
	return modelMap, nil
}
