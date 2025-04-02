// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package logs

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsStreamCreate,
		ReadContext:   resourceIbmLogsStreamRead,
		UpdateContext: resourceIbmLogsStreamUpdate,
		DeleteContext: resourceIbmLogsStreamDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_stream", "name"),
				Description:  "The name of the Event stream.",
			},
			"is_active": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the Event stream is active.",
			},
			"dpxl_expression": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_stream", "dpxl_expression"),
				Description:  "The DPXL expression of the Event stream.",
			},
			"compression_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_logs_stream", "compression_type"),
				Description:  "The compression type of the stream.",
			},
			"ibm_event_streams": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Configuration for IBM Event Streams.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"brokers": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The brokers of the IBM Event Streams.",
						},
						"topic": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The topic of the IBM Event Streams.",
						},
					},
				},
			},
			"streams_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Event stream.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the Event stream.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the Event stream.",
			},
		},
	}
}

func ResourceIbmLogsStreamValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-ZÀ-ÖØ-öø-ÿĀ-ſΑ-ωА-я一-龥ぁ-ゔァ-ヴー々〆〤0-9_\.,\-"{}()\[\]=!:#\/$|' ]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "dpxl_expression",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\p{L}\p{N}\p{P}\p{Z}\p{S}\p{M}]+$`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "compression_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "gzip, unspecified",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_logs_stream", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmLogsStreamCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_rule_group", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	upsertEventStreamTargetOptions := &logsv0.CreateEventStreamTargetOptions{}

	upsertEventStreamTargetOptions.SetName(d.Get("name").(string))
	upsertEventStreamTargetOptions.SetDpxlExpression(d.Get("dpxl_expression").(string))
	if _, ok := d.GetOk("is_active"); ok {
		upsertEventStreamTargetOptions.SetIsActive(d.Get("is_active").(bool))
	}
	if _, ok := d.GetOk("compression_type"); ok {
		upsertEventStreamTargetOptions.SetCompressionType(d.Get("compression_type").(string))
	}
	if _, ok := d.GetOk("ibm_event_streams"); ok {
		ibmEventStreamsModel, err := ResourceIbmLogsStreamMapToIbmEventStreams(d.Get("ibm_event_streams.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "create", "parse-ibm_event_streams").GetDiag()
		}
		upsertEventStreamTargetOptions.SetIbmEventStreams(ibmEventStreamsModel)
	}

	stream, _, err := logsClient.CreateEventStreamTargetWithContext(context, upsertEventStreamTargetOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateEventStreamTargetWithContext failed: %s", err.Error()), "ibm_logs_stream", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	streamsID := fmt.Sprintf("%s/%s/%d", region, instanceId, *stream.ID)
	d.SetId(streamsID)

	return resourceIbmLogsStreamRead(context, d, meta)
}

func resourceIbmLogsStreamRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, region, instanceId, streamsID, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	streamsIDInt, _ := strconv.ParseInt(streamsID, 10, 64)

	getEventStreamTargetsOptions := &logsv0.GetEventStreamTargetsOptions{}

	// getEventStreamTargetsOptions.SetID(streamsID)

	streamCollection, response, err := logsClient.GetEventStreamTargetsWithContext(context, getEventStreamTargetsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetEventStreamTargetsWithContext failed: %s", err.Error()), "ibm_logs_stream", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if streamCollection != nil {
		streamIds := make(map[int64]interface{}, 0)
		for _, stream := range streamCollection.Streams {
			streamIds[*stream.ID] = nil
			if *stream.ID == streamsIDInt {
				if err = d.Set("streams_id", streamsID); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting rule_group_id: %s", err))
				}
				if err = d.Set("instance_id", instanceId); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
				}
				if err = d.Set("region", region); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
				}

				if err = d.Set("name", stream.Name); err != nil {
					err = fmt.Errorf("Error setting name: %s", err)
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "set-name").GetDiag()
				}
				if !core.IsNil(stream.IsActive) {
					if err = d.Set("is_active", stream.IsActive); err != nil {
						err = fmt.Errorf("Error setting is_active: %s", err)
						return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "set-is_active").GetDiag()
					}
				}
				if err = d.Set("dpxl_expression", stream.DpxlExpression); err != nil {
					err = fmt.Errorf("Error setting dpxl_expression: %s", err)
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "set-dpxl_expression").GetDiag()
				}
				if !core.IsNil(stream.CompressionType) {
					if err = d.Set("compression_type", stream.CompressionType); err != nil {
						err = fmt.Errorf("Error setting compression_type: %s", err)
						return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "set-compression_type").GetDiag()
					}
				}
				if !core.IsNil(stream.IbmEventStreams) {
					ibmEventStreamsMap, err := ResourceIbmLogsStreamIbmEventStreamsToMap(stream.IbmEventStreams)
					if err != nil {
						return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "ibm_event_streams-to-map").GetDiag()
					}
					if err = d.Set("ibm_event_streams", []map[string]interface{}{ibmEventStreamsMap}); err != nil {
						err = fmt.Errorf("Error setting ibm_event_streams: %s", err)
						return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "set-ibm_event_streams").GetDiag()
					}
				}
				if !core.IsNil(stream.CreatedAt) {
					if err = d.Set("created_at", flex.DateTimeToString(stream.CreatedAt)); err != nil {
						err = fmt.Errorf("Error setting created_at: %s", err)
						return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "set-created_at").GetDiag()
					}
				}
				if !core.IsNil(stream.UpdatedAt) {
					if err = d.Set("updated_at", flex.DateTimeToString(stream.UpdatedAt)); err != nil {
						err = fmt.Errorf("Error setting updated_at: %s", err)
						return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "read", "set-updated_at").GetDiag()
					}
				}
			}
		}
		if _, ok := streamIds[streamsIDInt]; !ok {
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceIbmLogsStreamUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, _, _, streamsID, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	streamsIDInt, _ := strconv.ParseInt(streamsID, 10, 64)
	updateEventStreamTargetOptions := &logsv0.UpdateEventStreamTargetOptions{}

	updateEventStreamTargetOptions.SetID(streamsIDInt)

	hasChange := false

	if d.HasChange("name") ||
		d.HasChange("dpxl_expression") ||
		d.HasChange("is_active") ||
		d.HasChange("compression_type") ||
		d.HasChange("ibm_event_streams") {

		updateEventStreamTargetOptions.SetName(d.Get("name").(string))

		updateEventStreamTargetOptions.SetDpxlExpression(d.Get("dpxl_expression").(string))

		updateEventStreamTargetOptions.SetIsActive(d.Get("is_active").(bool))

		updateEventStreamTargetOptions.SetCompressionType(d.Get("compression_type").(string))

		ibmEventStreams, err := ResourceIbmLogsStreamMapToIbmEventStreams(d.Get("ibm_event_streams.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "update", "parse-ibm_event_streams").GetDiag()
		}
		updateEventStreamTargetOptions.SetIbmEventStreams(ibmEventStreams)

		hasChange = true
	}

	if hasChange {
		_, _, err = logsClient.UpdateEventStreamTargetWithContext(context, updateEventStreamTargetOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateEventStreamTargetWithContext failed: %s", err.Error()), "ibm_logs_stream", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsStreamRead(context, d, meta)
}

func resourceIbmLogsStreamDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_logs_stream", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, _, _, streamsID, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	streamsIDInt, _ := strconv.ParseInt(streamsID, 10, 64)
	deleteEventStreamTargetOptions := &logsv0.DeleteEventStreamTargetOptions{}

	deleteEventStreamTargetOptions.SetID(streamsIDInt)

	_, err = logsClient.DeleteEventStreamTargetWithContext(context, deleteEventStreamTargetOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteEventStreamTargetWithContext failed: %s", err.Error()), "ibm_logs_stream", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmLogsStreamMapToIbmEventStreams(modelMap map[string]interface{}) (*logsv0.IbmEventStreams, error) {
	model := &logsv0.IbmEventStreams{}
	model.Brokers = core.StringPtr(modelMap["brokers"].(string))
	model.Topic = core.StringPtr(modelMap["topic"].(string))
	return model, nil
}

func ResourceIbmLogsStreamIbmEventStreamsToMap(model *logsv0.IbmEventStreams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["brokers"] = *model.Brokers
	modelMap["topic"] = *model.Topic
	return modelMap, nil
}
