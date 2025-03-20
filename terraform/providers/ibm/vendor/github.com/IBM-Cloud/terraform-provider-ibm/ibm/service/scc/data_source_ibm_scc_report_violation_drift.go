// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccReportViolationDrift() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccReportViolationDriftRead,

		Schema: map[string]*schema.Schema{
			"report_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the scan that is associated with a report.",
			},
			"scan_time_duration": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The duration of the `scan_time` timestamp in number of days.",
			},
			"home_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the home account.",
			},
			"data_points": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of report violations data points.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the report.",
						},
						"report_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group ID that is associated with the report. The group ID combines profile, scope, and attachment IDs.",
						},
						"scan_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the scan was run.",
						},
						"controls": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The compliance stats.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The allowed values of an aggregated status for controls, specifications, assessments, and resources.",
									},
									"total_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of checks.",
									},
									"compliant_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of compliant checks.",
									},
									"not_compliant_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of checks that are not compliant.",
									},
									"unable_to_perform_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of checks that are unable to perform.",
									},
									"user_evaluation_required_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of checks that require a user evaluation.",
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccReportViolationDriftRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getReportViolationsDriftOptions := &securityandcompliancecenterapiv3.GetReportViolationsDriftOptions{}

	getReportViolationsDriftOptions.SetReportID(d.Get("report_id").(string))
	getReportViolationsDriftOptions.SetInstanceID(d.Get("instance_id").(string))
	if _, ok := d.GetOk("scan_time_duration"); ok {
		getReportViolationsDriftOptions.SetScanTimeDuration(int64(d.Get("scan_time_duration").(int)))
	}

	reportViolationsDrift, response, err := resultsClient.GetReportViolationsDriftWithContext(context, getReportViolationsDriftOptions)
	if err != nil {
		log.Printf("[DEBUG] GetReportViolationsDriftWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetReportViolationsDriftWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIbmSccReportViolationDriftID(d))

	if err = d.Set("home_account_id", reportViolationsDrift.HomeAccountID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting home_account_id: %s", err))
	}

	dataPoints := []map[string]interface{}{}
	if reportViolationsDrift.DataPoints != nil {
		for _, modelItem := range reportViolationsDrift.DataPoints {
			modelMap, err := dataSourceIbmSccReportViolationDriftReportViolationDataPointToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			dataPoints = append(dataPoints, modelMap)
		}
	}
	if err = d.Set("data_points", dataPoints); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting data_points %s", err))
	}

	return nil
}

// dataSourceIbmSccReportViolationDriftID returns a reasonable ID for the list.
func dataSourceIbmSccReportViolationDriftID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccReportViolationDriftReportViolationDataPointToMap(model *securityandcompliancecenterapiv3.ReportViolationDataPoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReportID != nil {
		modelMap["report_id"] = model.ReportID
	}
	if model.ReportGroupID != nil {
		modelMap["report_group_id"] = model.ReportGroupID
	}
	if model.ScanTime != nil {
		modelMap["scan_time"] = model.ScanTime
	}
	if model.ControlsSummary != nil {
		controlsMap, err := dataSourceIbmSccReportViolationDriftComplianceStatsToMap(model.ControlsSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["controls"] = []map[string]interface{}{controlsMap}
	}
	return modelMap, nil
}

func dataSourceIbmSccReportViolationDriftComplianceStatsToMap(model *securityandcompliancecenterapiv3.ComplianceStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.TotalCount != nil {
		modelMap["total_count"] = flex.IntValue(model.TotalCount)
	}
	if model.CompliantCount != nil {
		modelMap["compliant_count"] = flex.IntValue(model.CompliantCount)
	}
	if model.NotCompliantCount != nil {
		modelMap["not_compliant_count"] = flex.IntValue(model.NotCompliantCount)
	}
	if model.UnableToPerformCount != nil {
		modelMap["unable_to_perform_count"] = flex.IntValue(model.UnableToPerformCount)
	}
	if model.UserEvaluationRequiredCount != nil {
		modelMap["user_evaluation_required_count"] = flex.IntValue(model.UserEvaluationRequiredCount)
	}
	return modelMap, nil
}
