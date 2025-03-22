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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsStreams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsStreamsRead,

		Schema: map[string]*schema.Schema{
			"streams": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of Event Streams.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the Event stream.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Event stream.",
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
				},
			},
		},
	}
}

func dataSourceIbmLogsStreamsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_logs_streams", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getEventStreamTargetsOptions := &logsv0.GetEventStreamTargetsOptions{}

	streamCollection, _, err := logsClient.GetEventStreamTargetsWithContext(context, getEventStreamTargetsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetEventStreamTargetsWithContext failed: %s", err.Error()), "(Data) ibm_logs_streams", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmLogsStreamsID(d))

	streams := []map[string]interface{}{}
	for _, streamsItem := range streamCollection.Streams {
		streamsItemMap, err := DataSourceIbmLogsStreamsStreamToMap(&streamsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_logs_streams", "read", "streams-to-map").GetDiag()
		}
		streams = append(streams, streamsItemMap)
	}
	if err = d.Set("streams", streams); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting streams: %s", err), "(Data) ibm_logs_streams", "read", "set-streams").GetDiag()
	}

	return nil
}

// dataSourceIbmLogsStreamsID returns a reasonable ID for the list.
func dataSourceIbmLogsStreamsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmLogsStreamsStreamToMap(model *logsv0.Stream) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	modelMap["name"] = *model.Name
	if model.IsActive != nil {
		modelMap["is_active"] = *model.IsActive
	}
	modelMap["dpxl_expression"] = *model.DpxlExpression
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.CompressionType != nil {
		modelMap["compression_type"] = *model.CompressionType
	}
	if model.IbmEventStreams != nil {
		ibmEventStreamsMap, err := DataSourceIbmLogsStreamsIbmEventStreamsToMap(model.IbmEventStreams)
		if err != nil {
			return modelMap, err
		}
		modelMap["ibm_event_streams"] = []map[string]interface{}{ibmEventStreamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmLogsStreamsIbmEventStreamsToMap(model *logsv0.IbmEventStreams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["brokers"] = *model.Brokers
	modelMap["topic"] = *model.Topic
	return modelMap, nil
}
