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

func DataSourceIbmLogsDataUsageMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmLogsDataUsageMetricsRead,

		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The \"enabled\" parameter for metrics export.",
			},
		},
	}
}

func dataSourceIbmLogsDataUsageMetricsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logsClient, err := meta.(conns.ClientSession).LogsV0()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_logs_data_usage_metrics", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	region := getLogsInstanceRegion(logsClient, d)
	instanceId := d.Get("instance_id").(string)
	logsClient = getClientWithLogsInstanceEndpoint(logsClient, instanceId, region, getLogsInstanceEndpointType(logsClient, d))

	getDataUsageMetricsExportStatusOptions := &logsv0.GetDataUsageMetricsExportStatusOptions{}

	dataUsageMetricsExportStatus, _, err := logsClient.GetDataUsageMetricsExportStatusWithContext(context, getDataUsageMetricsExportStatusOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDataUsageMetricsExportStatusWithContext failed: %s", err.Error()), "(Data) ibm_logs_data_usage_metrics", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	usageMetricsId := fmt.Sprintf("%s/%s", region, instanceId)
	d.SetId(usageMetricsId)

	if err = d.Set("enabled", dataUsageMetricsExportStatus.Enabled); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting enabled: %s", err), "(Data) ibm_logs_data_usage_metrics", "read")
		return tfErr.GetDiag()
	}

	return nil
}
