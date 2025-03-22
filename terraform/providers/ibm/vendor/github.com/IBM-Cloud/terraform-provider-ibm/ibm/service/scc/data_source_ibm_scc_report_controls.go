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

func DataSourceIbmSccReportControls() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccReportControlsRead,

		Schema: map[string]*schema.Schema{
			"report_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the scan that is associated with a report.",
			},
			"control_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the control.",
			},
			"control_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the control.",
			},
			"control_description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the control.",
			},
			"control_category": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A control category value.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The compliance status value.",
			},
			"sort": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This field sorts controls by using a valid sort field. To learn more, see [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).",
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
			"home_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the home account.",
			},
			"controls": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of controls that are in the report.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control ID.",
						},
						"control_library_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control library ID.",
						},
						"control_library_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control library version.",
						},
						"control_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control name.",
						},
						"control_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control description.",
						},
						"control_category": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control category.",
						},
						"control_path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control path.",
						},
						"control_specifications": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of specifications that are on the page.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"control_specification_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The control specification ID.",
									},
									"component_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The component ID.",
									},
									"control_specification_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The component description.",
									},
									"environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The environment.",
									},
									"responsibility": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The responsibility for managing control specifications.",
									},
									"assessments": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The list of assessments.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"assessment_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The assessment ID.",
												},
												"assessment_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The assessment type.",
												},
												"assessment_method": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The assessment method.",
												},
												"assessment_description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The assessment description.",
												},
												"parameter_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of parameters of this assessment.",
												},
												"parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The list of parameters of this assessment.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"parameter_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The parameter name.",
															},
															"parameter_display_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The parameter display name.",
															},
															"parameter_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The parameter type.",
															},
															"parameter_value": &schema.Schema{
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
	})
}

func dataSourceIbmSccReportControlsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getReportControlsOptions := &securityandcompliancecenterapiv3.GetReportControlsOptions{}

	getReportControlsOptions.SetReportID(d.Get("report_id").(string))
	getReportControlsOptions.SetInstanceID(d.Get("instance_id").(string))
	if _, ok := d.GetOk("control_id"); ok {
		getReportControlsOptions.SetControlID(d.Get("control_id").(string))
	}
	if _, ok := d.GetOk("control_name"); ok {
		getReportControlsOptions.SetControlName(d.Get("control_name").(string))
	}
	if _, ok := d.GetOk("control_description"); ok {
		getReportControlsOptions.SetControlDescription(d.Get("control_description").(string))
	}
	if _, ok := d.GetOk("control_category"); ok {
		getReportControlsOptions.SetControlCategory(d.Get("control_category").(string))
	}
	if _, ok := d.GetOk("status"); ok {
		getReportControlsOptions.SetStatus(d.Get("status").(string))
	}
	if _, ok := d.GetOk("sort"); ok {
		getReportControlsOptions.SetSort(d.Get("sort").(string))
	}

	reportControls, response, err := resultsClient.GetReportControlsWithContext(context, getReportControlsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetReportControlsWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetReportControlsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIbmSccReportControlsID(d))

	if err = d.Set("home_account_id", reportControls.HomeAccountID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting home_account_id: %s", err))
	}

	controls := []map[string]interface{}{}
	if reportControls.Controls != nil {
		for _, modelItem := range reportControls.Controls {
			modelMap, err := dataSourceIbmSccReportControlsControlWithStatsToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			controls = append(controls, modelMap)
		}
	}
	if err = d.Set("controls", controls); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting controls %s", err))
	}
	getReportSummaryOptions := &securityandcompliancecenterapiv3.GetReportSummaryOptions{}

	getReportSummaryOptions.SetReportID(d.Get("report_id").(string))
	getReportSummaryOptions.SetInstanceID(d.Get("instance_id").(string))
	reportSummary, response, err := resultsClient.GetReportSummaryWithContext(context, getReportSummaryOptions)

	if err = d.Set("total_count", flex.IntValue(reportSummary.Controls.TotalCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting total_count: %s", err))
	}

	if err = d.Set("compliant_count", flex.IntValue(reportSummary.Controls.CompliantCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting compliant_count: %s", err))
	}

	if err = d.Set("not_compliant_count", flex.IntValue(reportSummary.Controls.NotCompliantCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting not_compliant_count: %s", err))
	}

	if err = d.Set("unable_to_perform_count", flex.IntValue(reportSummary.Controls.UnableToPerformCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting unable_to_perform_count: %s", err))
	}

	if err = d.Set("user_evaluation_required_count", flex.IntValue(reportSummary.Controls.UserEvaluationRequiredCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting user_evaluation_required_count: %s", err))
	}
	return nil
}

// dataSourceIbmSccReportControlsID returns a reasonable ID for the list.
func dataSourceIbmSccReportControlsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccReportControlsControlWithStatsToMap(model *securityandcompliancecenterapiv3.ControlWithStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.ControlLibraryID != nil {
		modelMap["control_library_id"] = model.ControlLibraryID
	}
	if model.ControlLibraryVersion != nil {
		modelMap["control_library_version"] = model.ControlLibraryVersion
	}
	if model.ControlName != nil {
		modelMap["control_name"] = model.ControlName
	}
	if model.ControlDescription != nil {
		modelMap["control_description"] = model.ControlDescription
	}
	if model.ControlCategory != nil {
		modelMap["control_category"] = model.ControlCategory
	}
	if model.ControlSpecifications != nil {
		controlSpecifications := []map[string]interface{}{}
		for _, controlSpecificationsItem := range model.ControlSpecifications {
			controlSpecificationsItemMap, err := dataSourceIbmSccReportControlsControlSpecificationWithStatsToMap(&controlSpecificationsItem)
			if err != nil {
				return modelMap, err
			}
			controlSpecifications = append(controlSpecifications, controlSpecificationsItemMap)
		}
		modelMap["control_specifications"] = controlSpecifications
	}
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

func dataSourceIbmSccReportControlsControlSpecificationWithStatsToMap(model *securityandcompliancecenterapiv3.ControlSpecificationWithStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ControlSpecificationID != nil {
		modelMap["control_specification_id"] = model.ControlSpecificationID
	}
	if model.ComponentID != nil {
		modelMap["component_id"] = model.ComponentID
	}
	if model.ControlSpecificationDescription != nil {
		modelMap["control_specification_description"] = model.ControlSpecificationDescription
	}
	if model.Environment != nil {
		modelMap["environment"] = model.Environment
	}
	if model.Responsibility != nil {
		modelMap["responsibility"] = model.Responsibility
	}
	if model.Assessments != nil {
		assessments := []map[string]interface{}{}
		for _, assessmentsItem := range model.Assessments {
			assessmentsItemMap, err := dataSourceIbmSccReportControlsAssessmentToMap(&assessmentsItem)
			if err != nil {
				return modelMap, err
			}
			assessments = append(assessments, assessmentsItemMap)
		}
		modelMap["assessments"] = assessments
	}
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

func dataSourceIbmSccReportControlsAssessmentToMap(model *securityandcompliancecenterapiv3.AssessmentWithStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AssessmentID != nil {
		modelMap["assessment_id"] = model.AssessmentID
	}
	if model.AssessmentType != nil {
		modelMap["assessment_type"] = model.AssessmentType
	}
	if model.AssessmentMethod != nil {
		modelMap["assessment_method"] = model.AssessmentMethod
	}
	if model.AssessmentDescription != nil {
		modelMap["assessment_description"] = model.AssessmentDescription
	}
	if model.ParameterCount != nil {
		modelMap["parameter_count"] = flex.IntValue(model.ParameterCount)
	}
	if model.Parameters != nil {
		parameters := []map[string]interface{}{}
		for _, parametersItem := range model.Parameters {
			parametersItemMap, err := dataSourceIbmSccReportControlsParameterInfoToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func dataSourceIbmSccReportControlsParameterInfoToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ParameterName != nil {
		modelMap["parameter_name"] = model.ParameterName
	}
	if model.ParameterDisplayName != nil {
		modelMap["parameter_display_name"] = model.ParameterDisplayName
	}
	if model.ParameterType != nil {
		modelMap["parameter_type"] = model.ParameterType
	}
	if model.ParameterValue != nil {
		modelMap["parameter_value"] = model.ParameterValue
	}
	return modelMap, nil
}
