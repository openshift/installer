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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsStream() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsStreamRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Event stream.",
			},
			"logs_streams_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Event stream.",
			},
			"is_active": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the Event stream is active.",
			},
			"dpxl_expression": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DPXL expression of the Event stream.",
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
			"compression_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The compression type of the stream.",
			},
			"ibm_event_streams": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Configuration for IBM Event Streams.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"brokers": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The brokers of the IBM Event Streams.",
						},
						"topic": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The topic of the IBM Event Streams.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmLogsStreamRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_logs_stream", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	streamsID, _ := strconv.ParseInt(d.Get("logs_streams_id").(string), 10, 64)

	getEventStreamTargetsOptions := &logsv0.GetEventStreamTargetsOptions{}

	streams, _, err := logsClient.GetEventStreamTargetsWithContext(context, getEventStreamTargetsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetEventStreamTargetsWithContext failed: %s", err.Error()), "(Data) ibm_logs_stream", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if streams != nil {
		streamIds := make(map[int64]interface{}, 0)
		for _, stream := range streams.Streams {
			streamIds[*stream.ID] = nil
			if *stream.ID == streamsID {
				d.SetId(fmt.Sprintf("%d", *stream.ID))

				if err = d.Set("name", stream.Name); err != nil {
					return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_logs_stream", "read", "set-name").GetDiag()
				}

				if !core.IsNil(stream.IsActive) {
					if err = d.Set("is_active", stream.IsActive); err != nil {
						return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting is_active: %s", err), "(Data) ibm_logs_stream", "read", "set-is_active").GetDiag()
					}
				}

				if err = d.Set("dpxl_expression", stream.DpxlExpression); err != nil {
					return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting dpxl_expression: %s", err), "(Data) ibm_logs_stream", "read", "set-dpxl_expression").GetDiag()
				}

				if !core.IsNil(stream.CreatedAt) {
					if err = d.Set("created_at", flex.DateTimeToString(stream.CreatedAt)); err != nil {
						return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_logs_stream", "read", "set-created_at").GetDiag()
					}
				}

				if !core.IsNil(stream.UpdatedAt) {
					if err = d.Set("updated_at", flex.DateTimeToString(stream.UpdatedAt)); err != nil {
						return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_logs_stream", "read", "set-updated_at").GetDiag()
					}
				}

				if !core.IsNil(stream.CompressionType) {
					if err = d.Set("compression_type", stream.CompressionType); err != nil {
						return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting compression_type: %s", err), "(Data) ibm_logs_stream", "read", "set-compression_type").GetDiag()
					}
				}

				if !core.IsNil(stream.IbmEventStreams) {
					ibmEventStreams := []map[string]interface{}{}
					ibmEventStreamsMap, err := DataSourceIbmLogsStreamIbmEventStreamsToMap(stream.IbmEventStreams)
					if err != nil {
						return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_logs_stream", "read", "ibm_event_streams-to-map").GetDiag()
					}
					ibmEventStreams = append(ibmEventStreams, ibmEventStreamsMap)
					if err = d.Set("ibm_event_streams", ibmEventStreams); err != nil {
						return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ibm_event_streams: %s", err), "(Data) ibm_logs_stream", "read", "set-ibm_event_streams").GetDiag()
					}
				}
			}
		}
		if _, ok := streamIds[streamsID]; !ok {
			d.SetId("")
			return flex.TerraformErrorf(err, fmt.Sprintf("Stream ID (%d) not found ", streamsID), "(Data) ibm_logs_stream", "read").GetDiag()
		}
	}

	return nil
}

func DataSourceIbmLogsStreamIbmEventStreamsToMap(model *logsv0.IbmEventStreams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["brokers"] = *model.Brokers
	modelMap["topic"] = *model.Topic
	return modelMap, nil
}
