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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func DataSourceIbmLogsOutgoingWebhook() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsOutgoingWebhookRead,

		Schema: map[string]*schema.Schema{
			"logs_outgoing_webhook_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Outbound Integration to delete.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the deployed Outbound Integrations to list.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Outbound Integration.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the Outbound Integration. Null for IBM Event Notifications integration.",
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
			"ibm_event_notifications": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The configuration of the IBM Event Notifications Outbound Integration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_notifications_instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the selected IBM Event Notifications instance.",
						},
						"endpoint_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoint type of integration",
						},
						"region_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region ID of the selected IBM Event Notifications instance.",
						},
						"source_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the created source in the IBM Event Notifications instance. Corresponds to the Cloud Logs instance crn. Not required when creating an Outbound Integration.",
						},
						"source_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the created source in the IBM Event Notifications instance. Not required when creating an Outbound Integration.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmLogsOutgoingWebhookRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_outgoing_webhook", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getOutgoingWebhookOptions := &logsv0.GetOutgoingWebhookOptions{}

	getOutgoingWebhookOptions.SetID(core.UUIDPtr(strfmt.UUID(d.Get("logs_outgoing_webhook_id").(string))))

	outgoingWebhookIntf, _, err := logsClient.GetOutgoingWebhookWithContext(context, getOutgoingWebhookOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetOutgoingWebhookWithContext failed: %s", err.Error()), "(Data) ibm_logs_outgoing_webhook", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	outgoingWebhook := outgoingWebhookIntf.(*logsv0.OutgoingWebhook)

	d.SetId(fmt.Sprintf("%s", *getOutgoingWebhookOptions.ID))

	if err = d.Set("type", outgoingWebhook.Type); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_logs_outgoing_webhook", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", outgoingWebhook.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_logs_outgoing_webhook", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("url", outgoingWebhook.URL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting url: %s", err), "(Data) ibm_logs_outgoing_webhook", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(outgoingWebhook.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_logs_outgoing_webhook", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", flex.DateTimeToString(outgoingWebhook.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_logs_outgoing_webhook", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("external_id", flex.IntValue(outgoingWebhook.ExternalID)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting external_id: %s", err), "(Data) ibm_logs_outgoing_webhook", "read")
		return tfErr.GetDiag()
	}

	ibmEventNotifications := []map[string]interface{}{}
	if outgoingWebhook.IbmEventNotifications != nil {
		modelMap, err := DataSourceIbmLogsOutgoingWebhookOutgoingWebhooksV1IbmEventNotificationsConfigToMap(outgoingWebhook.IbmEventNotifications)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_outgoing_webhook", "read")
			return tfErr.GetDiag()
		}
		ibmEventNotifications = append(ibmEventNotifications, modelMap)
	}
	if err = d.Set("ibm_event_notifications", ibmEventNotifications); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ibm_event_notifications: %s", err), "(Data) ibm_logs_outgoing_webhook", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func DataSourceIbmLogsOutgoingWebhookOutgoingWebhooksV1IbmEventNotificationsConfigToMap(model *logsv0.OutgoingWebhooksV1IbmEventNotificationsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["event_notifications_instance_id"] = model.EventNotificationsInstanceID.String()
	modelMap["region_id"] = *model.RegionID
	if model.EndpointType != nil {
		modelMap["endpoint_type"] = *model.SourceID
	}
	if model.SourceID != nil {
		modelMap["source_id"] = *model.SourceID
	}
	if model.SourceName != nil {
		modelMap["source_name"] = *model.SourceName
	}
	return modelMap, nil
}
