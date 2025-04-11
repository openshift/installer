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

func DataSourceIbmSccReportSummary() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccReportSummaryRead,

		Schema: map[string]*schema.Schema{
			"report_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the scan that is associated with a report.",
			},
			"isntance_id": &schema.Schema{
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
			"evaluations": &schema.Schema{
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
			"resources": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource summary.",
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
						"top_failed": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The top 10 resources that have the most failures.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource name.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource ID.",
									},
									"service": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The service that is managing the resource.",
									},
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The collection of different types of tags.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The collection of user tags.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"access": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The collection of access tags.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"service": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The collection of service tags.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"account": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The account that owns the resource.",
									},
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
					},
				},
			},
		},
	})
}

func dataSourceIbmSccReportSummaryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getReportSummaryOptions := &securityandcompliancecenterapiv3.GetReportSummaryOptions{}

	getReportSummaryOptions.SetReportID(d.Get("report_id").(string))
	getReportSummaryOptions.SetInstanceID(d.Get("instance_id").(string))

	reportSummary, response, err := resultsClient.GetReportSummaryWithContext(context, getReportSummaryOptions)
	if err != nil {
		log.Printf("[DEBUG] GetReportSummaryWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetReportSummaryWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIbmSccReportSummaryID(d))

	if err = d.Set("isntance_id", reportSummary.IsntanceID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting isntance_id: %s", err))
	}

	account := []map[string]interface{}{}
	if reportSummary.Account != nil {
		modelMap, err := dataSourceIbmSccReportSummaryAccountToMap(reportSummary.Account)
		if err != nil {
			return diag.FromErr(err)
		}
		account = append(account, modelMap)
	}
	if err = d.Set("account", account); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting account %s", err))
	}

	score := []map[string]interface{}{}
	if reportSummary.Score != nil {
		modelMap, err := dataSourceIbmSccReportSummaryComplianceScoreToMap(reportSummary.Score)
		if err != nil {
			return diag.FromErr(err)
		}
		score = append(score, modelMap)
	}
	if err = d.Set("score", score); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting score %s", err))
	}

	controls := []map[string]interface{}{}
	if reportSummary.Controls != nil {
		modelMap, err := dataSourceIbmSccReportSummaryComplianceStatsToMap(reportSummary.Controls)
		if err != nil {
			return diag.FromErr(err)
		}
		controls = append(controls, modelMap)
	}
	if err = d.Set("controls", controls); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting controls %s", err))
	}

	evaluations := []map[string]interface{}{}
	if reportSummary.Evaluations != nil {
		modelMap, err := dataSourceIbmSccReportSummaryEvalStatsToMap(reportSummary.Evaluations)
		if err != nil {
			return diag.FromErr(err)
		}
		evaluations = append(evaluations, modelMap)
	}
	if err = d.Set("evaluations", evaluations); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting evaluations %s", err))
	}

	resources := []map[string]interface{}{}
	if reportSummary.Resources != nil {
		modelMap, err := dataSourceIbmSccReportSummaryResourceSummaryToMap(reportSummary.Resources)
		if err != nil {
			return diag.FromErr(err)
		}
		resources = append(resources, modelMap)
	}
	if err = d.Set("resources", resources); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting resources %s", err))
	}

	return nil
}

// dataSourceIbmSccReportSummaryID returns a reasonable ID for the list.
func dataSourceIbmSccReportSummaryID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccReportSummaryAccountToMap(model *securityandcompliancecenterapiv3.Account) (map[string]interface{}, error) {
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

func dataSourceIbmSccReportSummaryComplianceScoreToMap(model *securityandcompliancecenterapiv3.ComplianceScore) (map[string]interface{}, error) {
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

func dataSourceIbmSccReportSummaryComplianceStatsToMap(model *securityandcompliancecenterapiv3.ComplianceStats) (map[string]interface{}, error) {
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

func dataSourceIbmSccReportSummaryEvalStatsToMap(model *securityandcompliancecenterapiv3.EvalStats) (map[string]interface{}, error) {
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

func dataSourceIbmSccReportSummaryResourceSummaryToMap(model *securityandcompliancecenterapiv3.ResourceSummary) (map[string]interface{}, error) {
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
	if model.TopFailed != nil {
		topFailed := []map[string]interface{}{}
		for _, topFailedItem := range model.TopFailed {
			topFailedItemMap, err := dataSourceIbmSccReportSummaryResourceSummaryItemToMap(&topFailedItem)
			if err != nil {
				return modelMap, err
			}
			topFailed = append(topFailed, topFailedItemMap)
		}
		modelMap["top_failed"] = topFailed
	}
	return modelMap, nil
}

func dataSourceIbmSccReportSummaryResourceSummaryItemToMap(model *securityandcompliancecenterapiv3.ResourceSummaryItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Service != nil {
		modelMap["service"] = model.Service
	}
	if model.Tags != nil {
		tagsMap, err := dataSourceIbmSccReportSummaryTagsToMap(model.Tags)
		if err != nil {
			return modelMap, err
		}
		modelMap["tags"] = []map[string]interface{}{tagsMap}
	}
	if model.Account != nil {
		modelMap["account"] = model.Account
	}
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

func dataSourceIbmSccReportSummaryTagsToMap(model *securityandcompliancecenterapiv3.Tags) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.User != nil {
		modelMap["user"] = model.User
	}
	if model.Access != nil {
		modelMap["access"] = model.Access
	}
	if model.Service != nil {
		modelMap["service"] = model.Service
	}
	return modelMap, nil
}
