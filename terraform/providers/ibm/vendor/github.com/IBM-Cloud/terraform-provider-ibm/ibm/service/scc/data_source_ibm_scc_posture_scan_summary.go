// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureScansSummary() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureScansSummaryRead,

		Schema: map[string]*schema.Schema{
			"scan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Your Scan ID.",
			},
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The profile ID. This can be obtained from the Security and Compliance Center UI by clicking on the profile name. The URL contains the ID.",
			},
			"discover_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scan discovery ID.",
			},
			"profile_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scan profile name.",
			},
			"scope_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The scan summary scope ID.",
			},
			"controls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of controls on the scan summary.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The scan summary control ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control status.",
						},
						"external_control_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The external control ID.",
						},
						"desciption": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The scan profile name.",
						},
						"goals": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of goals on the control.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the goal.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The goal ID.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The goal status.",
									},
									"severity": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The severity of the goal.",
									},
									"completed_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The report completed time.",
									},
									"error": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The error on goal validation.",
									},
									"resource_result": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The list of resource results.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource name.",
												},
												"types": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource control result status.",
												},
												"display_expected_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The expected results of a resource.",
												},
												"actual_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The actual results of a resource.",
												},
												"results_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The results information.",
												},
												"not_applicable_reason": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The reason for goal not applicable for a resource.",
												},
											},
										},
									},
								},
							},
						},
						"resource_statistics": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A scans summary controls.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pass_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The resource count of pass controls.",
									},
									"fail_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The resource count of fail controls.",
									},
									"unable_to_perform_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of resources that were unable to be scanned against a control.",
									},
									"not_applicable_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The resource count of not applicable(na) controls.",
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

func dataSourceIBMSccPostureScansSummaryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	scansSummaryOptions := &posturemanagementv2.ScansSummaryOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	scansSummaryOptions.SetAccountID(accountID)

	scansSummaryOptions.SetScanID(d.Get("scan_id").(string))
	scansSummaryOptions.SetProfileID(d.Get("profile_id").(string))

	summary, response, err := postureManagementClient.ScansSummaryWithContext(context, scansSummaryOptions)
	if err != nil {
		log.Printf("[DEBUG] ScansSummaryWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ScansSummaryWithContext failed %s\n%s", err, response))
	}

	d.SetId(*summary.ID)

	if err = d.Set("discover_id", summary.DiscoverID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting discover_id: %s", err))
	}
	if err = d.Set("profile_name", summary.ProfileName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting profile_name: %s", err))
	}
	if err = d.Set("scope_id", summary.ScopeID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting scope_id: %s", err))
	}

	if summary.Controls != nil {
		err = d.Set("controls", dataSourceSummaryFlattenControlsv2(summary.Controls))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting controls %s", err))
		}
	}

	return nil
}

func dataSourceSummaryFlattenControlsv2(result []posturemanagementv2.Control) (controls []map[string]interface{}) {
	for _, controlsItem := range result {
		controls = append(controls, dataSourceSummaryControlsToMapv2(controlsItem))
	}

	return controls
}

func dataSourceSummaryControlsToMapv2(controlsItem posturemanagementv2.Control) (controlsMap map[string]interface{}) {
	controlsMap = map[string]interface{}{}

	if controlsItem.ID != nil {
		controlsMap["id"] = controlsItem.ID
	}
	if controlsItem.Status != nil {
		controlsMap["status"] = controlsItem.Status
	}
	if controlsItem.ExternalControlID != nil {
		controlsMap["external_control_id"] = controlsItem.ExternalControlID
	}
	if controlsItem.Description != nil {
		controlsMap["description"] = controlsItem.Description
	}
	if controlsItem.Goals != nil {
		goalsList := []map[string]interface{}{}
		for _, goalsItem := range controlsItem.Goals {
			goalsList = append(goalsList, dataSourceSummaryControlsGoalsToMapv2(goalsItem))
		}
		controlsMap["goals"] = goalsList
	}
	if controlsItem.ResourceStatistics != nil {
		resourceStatisticsList := []map[string]interface{}{}
		resourceStatisticsMap := dataSourceSummaryControlsResourceStatisticsToMapv2(*controlsItem.ResourceStatistics)
		resourceStatisticsList = append(resourceStatisticsList, resourceStatisticsMap)
		controlsMap["resource_statistics"] = resourceStatisticsList
	}

	return controlsMap
}

func dataSourceSummaryControlsGoalsToMapv2(goalsItem posturemanagementv2.Goal) (goalsMap map[string]interface{}) {
	goalsMap = map[string]interface{}{}

	if goalsItem.Description != nil {
		goalsMap["description"] = goalsItem.Description
	}
	if goalsItem.ID != nil {
		goalsMap["id"] = goalsItem.ID
	}
	if goalsItem.Status != nil {
		goalsMap["status"] = goalsItem.Status
	}
	if goalsItem.Severity != nil {
		goalsMap["severity"] = goalsItem.Severity
	}
	if goalsItem.CompletedTime != nil {
		goalsMap["completed_time"] = goalsItem.CompletedTime.String()
	}
	if goalsItem.Error != nil {
		goalsMap["error"] = goalsItem.Error
	}
	if goalsItem.ResourceResult != nil {
		resourceResultList := []map[string]interface{}{}
		for _, resourceResultItem := range goalsItem.ResourceResult {
			resourceResultList = append(resourceResultList, dataSourceSummaryGoalsResourceResultToMapv2(resourceResultItem))
		}
		goalsMap["resource_result"] = resourceResultList
	}

	return goalsMap
}

func dataSourceSummaryGoalsResourceResultToMapv2(resourceResultItem posturemanagementv2.ResourceResult) (resourceResultMap map[string]interface{}) {
	resourceResultMap = map[string]interface{}{}

	if resourceResultItem.Name != nil {
		resourceResultMap["name"] = resourceResultItem.Name
	}
	if resourceResultItem.Types != nil {
		resourceResultMap["types"] = resourceResultItem.Types
	}
	if resourceResultItem.Status != nil {
		resourceResultMap["status"] = resourceResultItem.Status
	}
	if resourceResultItem.DisplayExpectedValue != nil {
		resourceResultMap["display_expected_value"] = resourceResultItem.DisplayExpectedValue
	}
	if resourceResultItem.ActualValue != nil {
		resourceResultMap["actual_value"] = resourceResultItem.ActualValue
	}
	if resourceResultItem.ResultsInfo != nil {
		resourceResultMap["results_info"] = resourceResultItem.ResultsInfo
	}
	if resourceResultItem.NotApplicableReason != nil {
		resourceResultMap["not_applicable_reason"] = resourceResultItem.NotApplicableReason
	}

	return resourceResultMap
}

func dataSourceSummaryControlsResourceStatisticsToMapv2(resourceStatisticsItem posturemanagementv2.ResourceStatistics) (resourceStatisticsMap map[string]interface{}) {
	resourceStatisticsMap = map[string]interface{}{}

	if resourceStatisticsItem.PassCount != nil {
		resourceStatisticsMap["pass_count"] = resourceStatisticsItem.PassCount
	}
	if resourceStatisticsItem.FailCount != nil {
		resourceStatisticsMap["fail_count"] = resourceStatisticsItem.FailCount
	}
	if resourceStatisticsItem.UnableToPerformCount != nil {
		resourceStatisticsMap["unable_to_perform_count"] = resourceStatisticsItem.UnableToPerformCount
	}
	if resourceStatisticsItem.NotApplicableCount != nil {
		resourceStatisticsMap["not_applicable_count"] = resourceStatisticsItem.NotApplicableCount
	}

	return resourceStatisticsMap
}
