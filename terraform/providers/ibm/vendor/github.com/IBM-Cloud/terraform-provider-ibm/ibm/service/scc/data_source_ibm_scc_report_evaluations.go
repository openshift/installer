// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccReportEvaluations() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccReportEvaluationsRead,

		Schema: map[string]*schema.Schema{
			"report_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the scan that is associated with a report.",
			},
			"assessment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the assessment.",
			},
			"component_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of component.",
			},
			"target_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the evaluation target.",
			},
			"target_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the evaluation target.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The evaluation status value.",
			},
			"first": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The page reference.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for the first and next page.",
						},
					},
				},
			},
			"home_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the home account.",
			},
			"evaluations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of evaluations that are on the page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"home_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the home account.",
						},
						"report_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the report that is associated to the evaluation.",
						},
						"control_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control ID. (Deprecated field)",
						},
						"component_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The component ID.",
						},
						"component_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of component.",
						},
						"assessment": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The control specification assessment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"assessment_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The assessment ID.",
									},
									"assessment_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The assessment type.",
									},
									"assessment_method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The assessment method.",
									},
									"assessment_description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The assessment description.",
									},
									"parameter_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of parameters of this assessment.",
									},
									"parameters": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The list of parameters of this assessment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameter_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The parameter name.",
												},
												"parameter_display_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The parameter display name.",
												},
												"parameter_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The parameter type.",
												},
												"parameter_value": {
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
						"evaluate_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the evaluation was made.",
						},
						"target": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The evaluation target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The target ID.",
									},
									"account_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The target account ID.",
									},
									"resource_crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The target resource CRN.",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The target resource name.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The target service name.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The allowed values of an evaluation status.",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason for the evaluation failure.",
						},
						"details": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The evaluation details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The evaluation properties.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"property": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property name.",
												},
												"property_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property description.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property operator.",
												},
												"expected_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property value.",
												},
												"found_value": {
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
	})
}

func dataSourceIbmSccReportEvaluationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resultsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	listReportEvaluationsOptions := &securityandcompliancecenterapiv3.ListReportEvaluationsOptions{}

	listReportEvaluationsOptions.SetReportID(d.Get("report_id").(string))
	listReportEvaluationsOptions.SetInstanceID(d.Get("instance_id").(string))
	if _, ok := d.GetOk("assessment_id"); ok {
		listReportEvaluationsOptions.SetAssessmentID(d.Get("assessment_id").(string))
	}
	if _, ok := d.GetOk("component_id"); ok {
		listReportEvaluationsOptions.SetComponentID(d.Get("component_id").(string))
	}
	if _, ok := d.GetOk("target_id"); ok {
		listReportEvaluationsOptions.SetTargetID(d.Get("target_id").(string))
	}
	if _, ok := d.GetOk("target_name"); ok {
		listReportEvaluationsOptions.SetTargetName(d.Get("target_name").(string))
	}
	if _, ok := d.GetOk("status"); ok {
		listReportEvaluationsOptions.SetStatus(d.Get("status").(string))
	}

	var pager *securityandcompliancecenterapiv3.ReportEvaluationsPager
	pager, err = resultsClient.NewReportEvaluationsPager(listReportEvaluationsOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] ReportEvaluationsPager %v:\n%s", pager, err)
	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] ReportEvaluationsPager.GetAll() failed %s", err)
		return diag.FromErr(flex.FmtErrorf("ReportEvaluationsPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIbmSccReportEvaluationsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIbmSccReportEvaluationsEvaluationToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("evaluations", mapSlice); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting evaluations %s", err))
	}

	return nil
}

// dataSourceIbmSccReportEvaluationsID returns a reasonable ID for the list.
func dataSourceIbmSccReportEvaluationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSccReportEvaluationsEvaluationToMap(model *securityandcompliancecenterapiv3.Evaluation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.HomeAccountID != nil {
		modelMap["home_account_id"] = model.HomeAccountID
	}
	if model.ReportID != nil {
		modelMap["report_id"] = model.ReportID
	}
	if model.ComponentName != nil {
		modelMap["component_name"] = model.ComponentName
	}
	if model.ComponentID != nil {
		modelMap["component_id"] = model.ComponentID
	}
	if model.Assessment != nil {
		assessmentMap, err := dataSourceIbmSccReportEvaluationsAssessmentToMap(model.Assessment)
		if err != nil {
			return modelMap, err
		}
		modelMap["assessment"] = []map[string]interface{}{assessmentMap}
	}
	if model.EvaluateTime != nil {
		modelMap["evaluate_time"] = model.EvaluateTime
	}
	if model.Target != nil {
		targetMap, err := dataSourceIbmSccReportEvaluationsTargetInfoToMap(model.Target)
		if err != nil {
			return modelMap, err
		}
		modelMap["target"] = []map[string]interface{}{targetMap}
	}
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.Reason != nil {
		modelMap["reason"] = model.Reason
	}
	if model.Details != nil {
		detailsMap, err := dataSourceIbmSccReportEvaluationsEvalDetailsToMap(model.Details)
		if err != nil {
			return modelMap, err
		}
		modelMap["details"] = []map[string]interface{}{detailsMap}
	}
	return modelMap, nil
}

func dataSourceIbmSccReportEvaluationsAssessmentToMap(model *securityandcompliancecenterapiv3.Assessment) (map[string]interface{}, error) {
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
			parametersItemMap, err := dataSourceIbmSccReportEvaluationsParameterInfoToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func dataSourceIbmSccReportEvaluationsParameterInfoToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
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

func dataSourceIbmSccReportEvaluationsTargetInfoToMap(model *securityandcompliancecenterapiv3.TargetInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.AccountID != nil {
		modelMap["account_id"] = model.AccountID
	}
	if model.ResourceCRN != nil {
		modelMap["resource_crn"] = model.ResourceCRN
	}
	if model.ResourceName != nil {
		modelMap["resource_name"] = model.ResourceName
	}
	if model.ServiceName != nil {
		modelMap["service_name"] = model.ServiceName
	}
	return modelMap, nil
}

func dataSourceIbmSccReportEvaluationsEvalDetailsToMap(model *securityandcompliancecenterapiv3.EvaluationDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := dataSourceIbmSccReportEvaluationsPropertyToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	return modelMap, nil
}

func dataSourceIbmSccReportEvaluationsPropertyToMap(model *securityandcompliancecenterapiv3.EvaluationProperty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Property != nil {
		modelMap["property"] = model.Property
	}
	if model.PropertyDescription != nil {
		modelMap["property_description"] = model.PropertyDescription
	}
	if model.Operator != nil {
		modelMap["operator"] = model.Operator
	}
	if model.ExpectedValue != nil {
		modelMap["expected_value"] = fmt.Sprintf("%v", model.ExpectedValue)
	}
	if model.FoundValue != nil {
		// modelMap["found_value"] = model.FoundValue
		fValIntf := model.FoundValue
		log.Printf("The Found value is = %v", fValIntf)
		modelMap["found_value"] = fmt.Sprintf("%v", fValIntf)
	}
	return modelMap, nil
}
