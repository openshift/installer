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

func DataSourceIbmSccLatestReports() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccLatestReportsRead,

		Schema: map[string]*schema.Schema{
			"sort": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This field sorts results by using a valid sort field. To learn more, see [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).",
			},
			"home_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the home account.",
			},
			"controls_summary": &schema.Schema{
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
			"evaluations_summary": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The evaluation stats.",
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
							Description: "The total number of evaluations.",
						},
						"pass_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of passed evaluations.",
						},
						"failure_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of failed evaluations.",
						},
						"error_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of evaluations that started, but did not finish, and ended with errors.",
						},
						"completed_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of completed evaluations.",
						},
					},
				},
			},
			"score": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The compliance score.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"passed": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of successful evaluations.",
						},
						"total_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of evaluations.",
						},
						"percent": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The percentage of successful evaluations.",
						},
					},
				},
			},
			"reports": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of reports.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the report.",
						},
						"group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group ID that is associated with the report. The group ID combines profile, scope, and attachment IDs.",
						},
						"created_on": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the report was created.",
						},
						"scan_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the scan was run.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the scan.",
						},
						"cos_object": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Object Storage object that is associated with the report.",
						},
						"instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"account": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The account that is associated with a report.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The account ID.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The account name.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The account type.",
									},
								},
							},
						},
						"profile": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The profile information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The profile ID.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The profile name.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The profile version.",
									},
								},
							},
						},
						"attachment": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The attachment that is associated with a report.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The attachment ID.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the attachment.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the attachment.",
									},
									"schedule": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The attachment schedule.",
									},
									"scope": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The scope of the attachment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this scope.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The environment that relates to this scope.",
												},
												"properties": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The properties that are supported for scoping by this environment.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The property name.",
															},
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The property value.",
															},
														},
													},
												},
											},
										},
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

func dataSourceIbmSccLatestReportsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getLatestReportsOptions := &securityandcompliancecenterapiv3.GetLatestReportsOptions{}
	getLatestReportsOptions.SetInstanceID(d.Get("instance_id").(string))

	if _, ok := d.GetOk("sort"); ok {
		getLatestReportsOptions.SetSort(d.Get("sort").(string))
	}

	reportLatest, response, err := resultsClient.GetLatestReportsWithContext(context, getLatestReportsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetLatestReportsWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetLatestReportsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIbmSccLatestReportsID(d))

	if err = d.Set("home_account_id", reportLatest.HomeAccountID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting home_account_id: %s", err))
	}

	controlsSummary := []map[string]interface{}{}
	if reportLatest.ControlsSummary != nil {
		modelMap, err := dataSourceIbmSccLatestReportsComplianceStatsToMap(reportLatest.ControlsSummary)
		if err != nil {
			return diag.FromErr(err)
		}
		controlsSummary = append(controlsSummary, modelMap)
	}
	if err = d.Set("controls_summary", controlsSummary); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting controls_summary %s", err))
	}

	evaluationsSummary := []map[string]interface{}{}
	if reportLatest.EvaluationsSummary != nil {
		modelMap, err := dataSourceIbmSccLatestReportsEvalStatsToMap(reportLatest.EvaluationsSummary)
		if err != nil {
			return diag.FromErr(err)
		}
		evaluationsSummary = append(evaluationsSummary, modelMap)
	}
	if err = d.Set("evaluations_summary", evaluationsSummary); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting evaluations_summary %s", err))
	}

	score := []map[string]interface{}{}
	if reportLatest.Score != nil {
		modelMap, err := dataSourceIbmSccLatestReportsComplianceScoreToMap(reportLatest.Score)
		if err != nil {
			return diag.FromErr(err)
		}
		score = append(score, modelMap)
	}
	if err = d.Set("score", score); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting score %s", err))
	}

	reports := []map[string]interface{}{}
	if reportLatest.Reports != nil {
		for _, modelItem := range reportLatest.Reports {
			modelMap, err := dataSourceIbmSccLatestReportsReportToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			reports = append(reports, modelMap)
		}
	}
	if err = d.Set("reports", reports); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting reports %s", err))
	}

	return nil
}

// dataSourceIbmSccLatestReportsID returns a reasonable ID for the list.
func dataSourceIbmSccLatestReportsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccLatestReportsComplianceStatsToMap(model *securityandcompliancecenterapiv3.ComplianceStats) (map[string]interface{}, error) {
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

func dataSourceIbmSccLatestReportsEvalStatsToMap(model *securityandcompliancecenterapiv3.EvalStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.TotalCount != nil {
		modelMap["total_count"] = flex.IntValue(model.TotalCount)
	}
	if model.PassCount != nil {
		modelMap["pass_count"] = flex.IntValue(model.PassCount)
	}
	if model.FailureCount != nil {
		modelMap["failure_count"] = flex.IntValue(model.FailureCount)
	}
	if model.ErrorCount != nil {
		modelMap["error_count"] = flex.IntValue(model.ErrorCount)
	}
	if model.CompletedCount != nil {
		modelMap["completed_count"] = flex.IntValue(model.CompletedCount)
	}
	return modelMap, nil
}

func dataSourceIbmSccLatestReportsComplianceScoreToMap(model *securityandcompliancecenterapiv3.ComplianceScore) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Passed != nil {
		modelMap["passed"] = flex.IntValue(model.Passed)
	}
	if model.TotalCount != nil {
		modelMap["total_count"] = flex.IntValue(model.TotalCount)
	}
	if model.Percent != nil {
		modelMap["percent"] = flex.IntValue(model.Percent)
	}
	return modelMap, nil
}

func dataSourceIbmSccLatestReportsReportToMap(model *securityandcompliancecenterapiv3.Report) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.GroupID != nil {
		modelMap["group_id"] = model.GroupID
	}
	if model.CreatedOn != nil {
		modelMap["created_on"] = model.CreatedOn.String()
	}
	if model.ScanTime != nil {
		modelMap["scan_time"] = model.ScanTime
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.CosObject != nil {
		modelMap["cos_object"] = model.CosObject
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = model.InstanceID
	}
	if model.Account != nil {
		accountMap, err := dataSourceIbmSccLatestReportsAccountToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Profile != nil {
		profileMap, err := dataSourceIbmSccLatestReportsProfileInfoToMap(model.Profile)
		if err != nil {
			return modelMap, err
		}
		modelMap["profile"] = []map[string]interface{}{profileMap}
	}
	if model.Attachment != nil {
		attachmentMap, err := dataSourceIbmSccLatestReportsAttachmentToMap(model.Attachment)
		if err != nil {
			return modelMap, err
		}
		modelMap["attachment"] = []map[string]interface{}{attachmentMap}
	}
	return modelMap, nil
}

func dataSourceIbmSccLatestReportsAccountToMap(model *securityandcompliancecenterapiv3.Account) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func dataSourceIbmSccLatestReportsProfileInfoToMap(model *securityandcompliancecenterapiv3.ProfileInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	return modelMap, nil
}

func dataSourceIbmSccLatestReportsAttachmentToMap(model *securityandcompliancecenterapiv3.Attachment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Schedule != nil {
		modelMap["schedule"] = model.Schedule
	}
	if model.Scopes != nil {
		scope := []map[string]interface{}{}
		for _, scopeItem := range model.Scopes {
			scopeItemMap, err := dataSourceIbmSccLatestReportsAttachmentScopeToMap(&scopeItem)
			if err != nil {
				return modelMap, err
			}
			scope = append(scope, scopeItemMap)
		}
		modelMap["scope"] = scope
	}
	return modelMap, nil
}

func dataSourceIbmSccLatestReportsAttachmentScopeToMap(model *securityandcompliancecenterapiv3.Scope) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Environment != nil {
		modelMap["environment"] = model.Environment
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := dataSourceIbmSccLatestReportsScopePropertyToMap(propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	return modelMap, nil
}

func dataSourceIbmSccLatestReportsScopePropertyToMap(model securityandcompliancecenterapiv3.ScopePropertyIntf) (map[string]interface{}, error) {
	return scopePropertiesToMap(model)
}
