// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccControlLibrary() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccControlLibraryRead,

		Schema: map[string]*schema.Schema{
			"control_library_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The control library ID.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID.",
			},
			"control_library_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The control library name.",
			},
			"control_library_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The control library description.",
			},
			"control_library_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The control library type.",
			},
			"version_group_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version group label.",
			},
			"control_library_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The control library version.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the control library was created.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the control library.",
			},
			"updated_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the control library was updated.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who updated the control library.",
			},
			"latest": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The latest version of the control library.",
			},
			"hierarchy_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The indication of whether hierarchy is enabled for the control library.",
			},
			"controls_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of controls.",
			},
			"control_parents_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of parent controls in the control library.",
			},
			"controls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of controls in a control library.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the control library that contains the profile.",
						},
						"control_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control name.",
						},
						"control_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control description.",
						},
						"control_category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control category.",
						},
						"control_parent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parent control.",
						},
						"control_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The control tags.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"control_specifications": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The control specifications.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"control_specification_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The control specification ID.",
									},
									"responsibility": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The responsibility for managing the control.",
									},
									"component_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The component ID.",
									},
									"component_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The component name.",
									},
									"environment": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The control specifications environment.",
									},
									"control_specification_description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The control specifications description.",
									},
									"assessments_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of assessments.",
									},
									"assessments": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The assessments.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"assessment_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The assessment ID.",
												},
												"assessment_method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The assessment method.",
												},
												"assessment_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The assessment type.",
												},
												"assessment_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The assessment description.",
												},
												"parameter_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The parameter count.",
												},
												"parameters": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The parameters.",
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
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"control_docs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The control documentation.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"control_docs_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the control documentation.",
									},
									"control_docs_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of control documentation.",
									},
								},
							},
						},
						"control_requirement": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is this a control that can be automated or manually evaluated.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The control status.",
						},
					},
				},
			},
		},
	})
}

func dataSourceIbmSccControlLibraryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getControlLibraryOptions := &securityandcompliancecenterapiv3.GetControlLibraryOptions{}

	getControlLibraryOptions.SetControlLibraryID(d.Get("control_library_id").(string))
	getControlLibraryOptions.SetInstanceID(d.Get("instance_id").(string))

	controlLibrary, response, err := securityandcompliancecenterapiClient.GetControlLibraryWithContext(context, getControlLibraryOptions)
	if err != nil {
		log.Printf("[DEBUG] GetControlLibraryWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetControlLibraryWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getControlLibraryOptions.ControlLibraryID))

	if err = d.Set("account_id", controlLibrary.AccountID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting account_id: %s", err))
	}

	if err = d.Set("control_library_name", controlLibrary.ControlLibraryName); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting control_library_name: %s", err))
	}

	if err = d.Set("control_library_description", controlLibrary.ControlLibraryDescription); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting control_library_description: %s", err))
	}

	if err = d.Set("control_library_type", controlLibrary.ControlLibraryType); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting control_library_type: %s", err))
	}

	if err = d.Set("version_group_label", controlLibrary.VersionGroupLabel); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting version_group_label: %s", err))
	}

	if err = d.Set("control_library_version", controlLibrary.ControlLibraryVersion); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting control_library_version: %s", err))
	}

	if err = d.Set("created_on", flex.DateTimeToString(controlLibrary.CreatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_on: %s", err))
	}

	if err = d.Set("created_by", controlLibrary.CreatedBy); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_by: %s", err))
	}

	if err = d.Set("updated_on", flex.DateTimeToString(controlLibrary.UpdatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_on: %s", err))
	}

	if err = d.Set("updated_by", controlLibrary.UpdatedBy); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_by: %s", err))
	}

	if err = d.Set("latest", controlLibrary.Latest); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting latest: %s", err))
	}

	if err = d.Set("hierarchy_enabled", controlLibrary.HierarchyEnabled); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting hierarchy_enabled: %s", err))
	}

	if err = d.Set("controls_count", flex.IntValue(controlLibrary.ControlsCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting controls_count: %s", err))
	}

	if err = d.Set("control_parents_count", flex.IntValue(controlLibrary.ControlParentsCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting control_parents_count: %s", err))
	}

	controls := []map[string]interface{}{}
	if controlLibrary.Controls != nil {
		for _, modelItem := range controlLibrary.Controls {
			modelMap, err := dataSourceIbmSccControlLibraryControlsInControlLibToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			controls = append(controls, modelMap)
		}
	}
	if err = d.Set("controls", controls); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting controls %s", err))
	}

	return nil
}

func dataSourceIbmSccControlLibraryControlsInControlLibToMap(model *securityandcompliancecenterapiv3.Control) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ControlName != nil {
		modelMap["control_name"] = model.ControlName
	}
	if model.ControlID != nil {
		modelMap["control_id"] = model.ControlID
	}
	if model.ControlDescription != nil {
		modelMap["control_description"] = model.ControlDescription
	}
	if model.ControlCategory != nil {
		modelMap["control_category"] = model.ControlCategory
	}
	if model.ControlParent != nil {
		modelMap["control_parent"] = model.ControlParent
	}
	if model.ControlTags != nil {
		modelMap["control_tags"] = model.ControlTags
	}
	if model.ControlSpecifications != nil {
		controlSpecifications := []map[string]interface{}{}
		for _, controlSpecificationsItem := range model.ControlSpecifications {
			controlSpecificationsItemMap, err := dataSourceIbmSccControlLibraryControlSpecificationsToMap(&controlSpecificationsItem)
			if err != nil {
				return modelMap, err
			}
			controlSpecifications = append(controlSpecifications, controlSpecificationsItemMap)
		}
		modelMap["control_specifications"] = controlSpecifications
	}
	if model.ControlDocs != nil {
		controlDocsMap, err := dataSourceIbmSccControlLibraryControlDocsToMap(model.ControlDocs)
		if err != nil {
			return modelMap, err
		}
		modelMap["control_docs"] = []map[string]interface{}{controlDocsMap}
	}
	if model.ControlRequirement != nil {
		modelMap["control_requirement"] = model.ControlRequirement
	}
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	return modelMap, nil
}

func dataSourceIbmSccControlLibraryControlSpecificationsToMap(model *securityandcompliancecenterapiv3.ControlSpecification) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["control_specification_id"] = model.ID
	}
	if model.Responsibility != nil {
		modelMap["responsibility"] = model.Responsibility
	}
	if model.ComponentID != nil {
		modelMap["component_id"] = model.ComponentID
	}
	if model.ComponentName != nil {
		modelMap["component_name"] = model.ComponentName
	}
	if model.Environment != nil {
		modelMap["environment"] = model.Environment
	}
	if model.Description != nil {
		modelMap["control_specification_description"] = model.Description
	}
	if model.AssessmentsCount != nil {
		modelMap["assessments_count"] = flex.IntValue(model.AssessmentsCount)
	}
	if model.Assessments != nil {
		assessments := []map[string]interface{}{}
		for _, assessmentsItem := range model.Assessments {
			assessmentsItemMap, err := dataSourceIbmSccControlLibraryImplementationToMap(&assessmentsItem)
			if err != nil {
				return modelMap, err
			}
			assessments = append(assessments, assessmentsItemMap)
		}
		modelMap["assessments"] = assessments
	}
	return modelMap, nil
}

func dataSourceIbmSccControlLibraryImplementationToMap(model *securityandcompliancecenterapiv3.Assessment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AssessmentID != nil {
		modelMap["assessment_id"] = model.AssessmentID
	}
	if model.AssessmentMethod != nil {
		modelMap["assessment_method"] = model.AssessmentMethod
	}
	if model.AssessmentType != nil {
		modelMap["assessment_type"] = model.AssessmentType
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
			parametersItemMap, err := dataSourceIbmSccControlLibraryParameterInfoToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func dataSourceIbmSccControlLibraryParameterInfoToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
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
	return modelMap, nil
}

func dataSourceIbmSccControlLibraryControlDocsToMap(model *securityandcompliancecenterapiv3.ControlDoc) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ControlDocsID != nil {
		modelMap["control_docs_id"] = model.ControlDocsID
	}
	if model.ControlDocsType != nil {
		modelMap["control_docs_type"] = model.ControlDocsType
	}
	return modelMap, nil
}
