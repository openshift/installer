// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/logs-go-sdk/logsv0"
)

func ResourceIbmLogsDataUsageMetrics() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmLogsDataUsageMetricsCreate,
		ReadContext:   resourceIbmLogsDataUsageMetricsRead,
		UpdateContext: resourceIbmLogsDataUsageMetricsUpdate,
		DeleteContext: resourceIbmLogsDataUsageMetricsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The \"enabled\" parameter for metrics export.",
			},
		},
	}
}

func resourceIbmLogsDataUsageMetricsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_data_usage_metrics", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	updateDataUsageMetricsExportStatusOptions := &logsv0.UpdateDataUsageMetricsExportStatusOptions{}

	updateDataUsageMetricsExportStatusOptions.SetEnabled(d.Get("enabled").(bool))

	_, _, err = logsClient.UpdateDataUsageMetricsExportStatusWithContext(context, updateDataUsageMetricsExportStatusOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateDataUsageMetricsExportStatusWithContext failed: %s", err.Error()), "ibm_logs_data_usage_metrics", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	usageMetricsId := fmt.Sprintf("%s/%s", region, instanceId)
	d.SetId(usageMetricsId)

	return resourceIbmLogsDataUsageMetricsRead(context, d, meta)
}

func resourceIbmLogsDataUsageMetricsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_data_usage_metrics", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, region, instanceId, _, err := updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getDataUsageMetricsExportStatusOptions := &logsv0.GetDataUsageMetricsExportStatusOptions{}

	dataUsageMetricsExportStatus, response, err := logsClient.GetDataUsageMetricsExportStatusWithContext(context, getDataUsageMetricsExportStatusOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDataUsageMetricsExportStatusWithContext failed: %s", err.Error()), "ibm_logs_data_usage_metrics", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}

	if err = d.Set("enabled", dataUsageMetricsExportStatus.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
	}

	return nil
}

func resourceIbmLogsDataUsageMetricsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_logs_data_usage_metrics", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	logsClient, _, _, _, err = updateClientURLWithInstanceEndpoint(d.Id(), logsClient, d)
	if err != nil {
		return diag.FromErr(err)
	}
	updateDataUsageMetricsExportStatusOptions := &logsv0.UpdateDataUsageMetricsExportStatusOptions{}

	hasChange := false

	if d.HasChange("enabled") {
		updateDataUsageMetricsExportStatusOptions.SetEnabled(d.Get("enabled").(bool))
		hasChange = true
	}
	if hasChange {
		_, _, err = logsClient.UpdateDataUsageMetricsExportStatusWithContext(context, updateDataUsageMetricsExportStatusOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateDataUsageMetricsExportStatusWithContext failed: %s", err.Error()), "ibm_logs_data_usage_metrics", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmLogsDataUsageMetricsRead(context, d, meta)
}

func resourceIbmLogsDataUsageMetricsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "delete" operation.
	d.SetId("")
	return nil
}
