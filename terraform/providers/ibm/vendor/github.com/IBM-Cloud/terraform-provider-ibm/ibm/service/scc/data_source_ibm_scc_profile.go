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

func DataSourceIbmSccProfile() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccProfileRead,

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The profile ID.",
			},
			"profile_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The profile name.",
			},
			"profile_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The profile description.",
			},
			"profile_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The profile type, such as custom or predefined.",
			},
			"profile_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version status of the profile.",
			},
			"version_group_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version group label of the profile.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID.",
			},
			"latest": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The latest version of the profile.",
			},
			"hierarchy_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The indication of whether hierarchy is enabled for the profile.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the profile.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the profile was created.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who updated the profile.",
			},
			"updated_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the profile was updated.",
			},
			"controls_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of controls for the profile.",
			},
			"control_parents_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of parent controls for the profile.",
			},
			"attachments_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of attachments related to this profile.",
			},
			"controls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The array of controls that are used to create the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_library_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the control library that contains the profile.",
						},
						"control_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the control library that contains the profile.",
						},
						"control_library_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The most recent version of the control library.",
						},
						"control_name": {
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
						"control_requirement": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is this a control that can be automated or manually evaluated.",
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
						"control_specifications_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of control specifications.",
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
					},
				},
			},
			"default_parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The default parameters of the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assessment_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the implementation.",
						},
						"assessment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The implementation ID of the parameter.",
						},
						"parameter_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter name.",
						},
						"parameter_default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default value of the parameter.",
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
	})
}

func dataSourceIbmSccProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileOptions := &securityandcompliancecenterapiv3.GetProfileOptions{}

	getProfileOptions.SetProfileID(d.Get("profile_id").(string))
	getProfileOptions.SetInstanceID(d.Get("instance_id").(string))

	profile, response, err := securityandcompliancecenterapiClient.GetProfileWithContext(context, getProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProfileWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetProfileWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getProfileOptions.ProfileID))

	if err = d.Set("profile_name", profile.ProfileName); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile_name: %s", err))
	}

	if err = d.Set("profile_description", profile.ProfileDescription); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile_description: %s", err))
	}

	if err = d.Set("profile_type", profile.ProfileType); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile_type: %s", err))
	}

	if err = d.Set("profile_version", profile.ProfileVersion); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile_version: %s", err))
	}

	if err = d.Set("version_group_label", profile.VersionGroupLabel); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting version_group_label: %s", err))
	}

	if err = d.Set("latest", profile.Latest); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting latest: %s", err))
	}

	if err = d.Set("hierarchy_enabled", profile.HierarchyEnabled); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting hierarchy_enabled: %s", err))
	}

	if err = d.Set("created_by", profile.CreatedBy); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_by: %s", err))
	}

	if err = d.Set("created_on", flex.DateTimeToString(profile.CreatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting created_on: %s", err))
	}

	if err = d.Set("updated_by", profile.UpdatedBy); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_by: %s", err))
	}

	if err = d.Set("updated_on", flex.DateTimeToString(profile.UpdatedOn)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting updated_on: %s", err))
	}

	if err = d.Set("controls_count", flex.IntValue(profile.ControlsCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting controls_count: %s", err))
	}

	if err = d.Set("control_parents_count", flex.IntValue(profile.ControlParentsCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting control_parents_count: %s", err))
	}

	if err = d.Set("attachments_count", flex.IntValue(profile.AttachmentsCount)); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting attachments_count: %s", err))
	}

	controls := []map[string]interface{}{}
	if profile.Controls != nil {
		for _, modelItem := range profile.Controls {
			modelMap, err := dataSourceIbmSccProfileProfileControlsToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			controls = append(controls, modelMap)
		}
	}
	if err = d.Set("controls", controls); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting controls %s", err))
	}

	defaultParameters := []map[string]interface{}{}
	if profile.DefaultParameters != nil {
		for _, modelItem := range profile.DefaultParameters {
			modelMap, err := dataSourceIbmSccProfileDefaultParametersPrototypeToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			defaultParameters = append(defaultParameters, modelMap)
		}
	}
	if err = d.Set("default_parameters", defaultParameters); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting default_parameters %s", err))
	}

	return nil
}

func dataSourceIbmSccProfileProfileControlsToMap(model *securityandcompliancecenterapiv3.ProfileControlsInResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ControlLibraryID != nil {
		modelMap["control_library_id"] = model.ControlLibraryID
	}
	if model.ControlID != nil {
		modelMap["control_id"] = model.ControlID
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
	if model.ControlParent != nil {
		modelMap["control_parent"] = model.ControlParent
	}
	if model.ControlRequirement != nil {
		modelMap["control_requirement"] = model.ControlRequirement
	}
	if model.ControlDocs != nil {
		controlDocsMap, err := dataSourceIbmSccProfileControlDocsToMap(model.ControlDocs)
		if err != nil {
			return modelMap, err
		}
		modelMap["control_docs"] = []map[string]interface{}{controlDocsMap}
	}
	if model.ControlSpecificationsCount != nil {
		modelMap["control_specifications_count"] = flex.IntValue(model.ControlSpecificationsCount)
	}
	if model.ControlSpecifications != nil {
		controlSpecifications := []map[string]interface{}{}
		for _, controlSpecificationsItem := range model.ControlSpecifications {
			controlSpecificationsItemMap, err := dataSourceIbmSccProfileControlSpecificationsToMap(&controlSpecificationsItem)
			if err != nil {
				return modelMap, err
			}
			controlSpecifications = append(controlSpecifications, controlSpecificationsItemMap)
		}
		modelMap["control_specifications"] = controlSpecifications
	}
	return modelMap, nil
}

func dataSourceIbmSccProfileControlDocsToMap(model *securityandcompliancecenterapiv3.ControlDoc) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ControlDocsID != nil {
		modelMap["control_docs_id"] = model.ControlDocsID
	}
	if model.ControlDocsType != nil {
		modelMap["control_docs_type"] = model.ControlDocsType
	}
	return modelMap, nil
}

func dataSourceIbmSccProfileControlSpecificationsToMap(model *securityandcompliancecenterapiv3.ControlSpecification) (map[string]interface{}, error) {
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
			assessmentsItemMap, err := dataSourceIbmSccProfileImplementationToMap(&assessmentsItem)
			if err != nil {
				return modelMap, err
			}
			assessments = append(assessments, assessmentsItemMap)
		}
		modelMap["assessments"] = assessments
	}
	return modelMap, nil
}

func dataSourceIbmSccProfileImplementationToMap(model *securityandcompliancecenterapiv3.Assessment) (map[string]interface{}, error) {
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
			parametersItemMap, err := dataSourceIbmSccProfileParameterInfoToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func dataSourceIbmSccProfileParameterInfoToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
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

func dataSourceIbmSccProfileDefaultParametersPrototypeToMap(model *securityandcompliancecenterapiv3.DefaultParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AssessmentType != nil {
		modelMap["assessment_type"] = model.AssessmentType
	}
	if model.AssessmentID != nil {
		modelMap["assessment_id"] = model.AssessmentID
	}
	if model.ParameterName != nil {
		modelMap["parameter_name"] = model.ParameterName
	}
	if model.ParameterDefaultValue != nil {
		modelMap["parameter_default_value"] = model.ParameterDefaultValue
	}
	if model.ParameterDisplayName != nil {
		modelMap["parameter_display_name"] = model.ParameterDisplayName
	}
	if model.ParameterType != nil {
		modelMap["parameter_type"] = model.ParameterType
	}
	return modelMap, nil
}
