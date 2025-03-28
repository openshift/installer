// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func ResourceIbmSccProfile() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		CreateContext: resourceIbmSccProfileCreate,
		ReadContext:   resourceIbmSccProfileRead,
		UpdateContext: resourceIbmSccProfileUpdate,
		DeleteContext: resourceIbmSccProfileDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The profile name.",
			},
			"profile_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_profile", "profile_name"),
				Description:  "The profile name.",
			},
			"profile_description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_profile", "profile_description"),
				Description:  "The profile description.",
			},
			"profile_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_profile", "profile_type"),
				Description:  "The profile type, such as custom or predefined.",
			},
			"controls": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The array of controls that are used to create the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_library_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the control library that contains the profile.",
						},
						"control_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique ID of the control inside the control library.",
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
										Optional:    true,
										Description: "The ID of the control documentation.",
									},
									"control_docs_type": {
										Type:        schema.TypeString,
										Optional:    true,
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
				Optional:    true,
				Description: "The default parameters of the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assessment_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of the implementation.",
						},
						"assessment_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The implementation ID of the parameter.",
						},
						"parameter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The parameter name.",
						},
						"parameter_default_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The default value of the parameter.",
						},
						"parameter_display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The parameter display name.",
						},
						"parameter_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The parameter type.",
						},
					},
				},
			},
			"profile_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0.0.0",
				Description: "The version status of the profile.",
			},
			"version_group_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version group label of the profile.",
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
		},
	})
}

func ResourceIbmSccProfileValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9_\s\-]*$`,
			MinValueLength:             2,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "profile_description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9_,'"\s\-\[\]]+$`,
			MinValueLength:             2,
			MaxValueLength:             256,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_profile", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSccProfileCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Print("[DEBUG] Starting resourceIbmSccProfileCreate")
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	bodyModelMap := map[string]interface{}{}
	createProfileOptions := &securityandcompliancecenterapiv3.CreateProfileOptions{}

	instance_id := d.Get("instance_id").(string)
	bodyModelMap["instance_id"] = instance_id
	bodyModelMap["profile_name"] = d.Get("profile_name")
	bodyModelMap["profile_description"] = d.Get("profile_description")
	bodyModelMap["profile_type"] = "custom"
	// manual change for profile_version
	bodyModelMap["profile_version"] = d.Get("profile_version")
	if _, ok := d.GetOk("controls"); ok {
		bodyModelMap["controls"] = d.Get("controls")
	} else {
		bodyModelMap["controls"] = []interface{}{}
	}
	if _, ok := d.GetOk("default_parameters"); ok {
		bodyModelMap["default_parameters"] = d.Get("default_parameters")
	} else {
		bodyModelMap["default_parameters"] = []interface{}{}
	}
	convertedModel, err := resourceIbmSccProfileMapToProfilePrototype(bodyModelMap)
	if err != nil {
		return diag.FromErr(err)
	}
	createProfileOptions = convertedModel
	createProfileOptions.SetInstanceID(instance_id)

	profile, response, err := securityandcompliancecenterapiClient.CreateProfileWithContext(context, createProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProfileWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("CreateProfileWithContext failed %s\n%s", err, response))
	}

	d.SetId(instance_id + "/" + *profile.ID)

	return resourceIbmSccProfileRead(context, d, meta)
}

func resourceIbmSccProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Print("[DEBUG] Starting resourceIbmSccProfileRead")
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileOptions := &securityandcompliancecenterapiv3.GetProfileOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	getProfileOptions.SetInstanceID(parts[0])
	getProfileOptions.SetProfileID(parts[1])

	profile, response, err := securityandcompliancecenterapiClient.GetProfileWithContext(context, getProfileOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProfileWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetProfileWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_id", parts[0]); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("profile_id", parts[1]); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile_id: %s", err))
	}
	if err = d.Set("profile_name", profile.ProfileName); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile_name: %s", err))
	}
	if err = d.Set("profile_description", profile.ProfileDescription); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile_description: %s", err))
	}
	if err = d.Set("profile_type", profile.ProfileType); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting profile_type: %s", err))
	}
	controls := []map[string]interface{}{}
	for _, controlsItem := range profile.Controls {
		controlsItemMap, err := resourceIbmSccProfileProfileControlsToMap(&controlsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		controls = append(controls, controlsItemMap)
	}
	if err = d.Set("controls", controls); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting controls: %s", err))
	}
	if len(profile.DefaultParameters) > 0 {
		defaultParameters := []map[string]interface{}{}
		for _, defaultParametersItem := range profile.DefaultParameters {
			defaultParametersItemMap, err := resourceIbmSccProfileDefaultParametersPrototypeToMap(&defaultParametersItem)
			if err != nil {
				return diag.FromErr(err)
			}
			defaultParameters = append(defaultParameters, defaultParametersItemMap)
		}
		if err = d.Set("default_parameters", defaultParameters); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting default_parameters: %s", err))
		}
	}
	if !core.IsNil(profile.ProfileVersion) {
		if err = d.Set("profile_version", profile.ProfileVersion); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting profile_version: %s", err))
		}
	}
	if !core.IsNil(profile.VersionGroupLabel) {
		if err = d.Set("version_group_label", profile.VersionGroupLabel); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting version_group_label: %s", err))
		}
	}
	if !core.IsNil(profile.InstanceID) {
		if err = d.Set("instance_id", profile.InstanceID); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting instance_id: %s", err))
		}
	}
	if !core.IsNil(profile.Latest) {
		if err = d.Set("latest", profile.Latest); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting latest: %s", err))
		}
	}
	if !core.IsNil(profile.HierarchyEnabled) {
		if err = d.Set("hierarchy_enabled", profile.HierarchyEnabled); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting hierarchy_enabled: %s", err))
		}
	}
	if !core.IsNil(profile.CreatedBy) {
		if err = d.Set("created_by", profile.CreatedBy); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting created_by: %s", err))
		}
	}
	if !core.IsNil(profile.CreatedOn) {
		if err = d.Set("created_on", flex.DateTimeToString(profile.CreatedOn)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting created_on: %s", err))
		}
	}
	if !core.IsNil(profile.UpdatedBy) {
		if err = d.Set("updated_by", profile.UpdatedBy); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting updated_by: %s", err))
		}
	}
	if !core.IsNil(profile.UpdatedOn) {
		if err = d.Set("updated_on", flex.DateTimeToString(profile.UpdatedOn)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting updated_on: %s", err))
		}
	}
	if !core.IsNil(profile.ControlsCount) {
		if err = d.Set("controls_count", flex.IntValue(profile.ControlsCount)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting controls_count: %s", err))
		}
	}
	if !core.IsNil(profile.ControlParentsCount) {
		if err = d.Set("control_parents_count", flex.IntValue(profile.ControlParentsCount)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting control_parents_count: %s", err))
		}
	}
	if !core.IsNil(profile.AttachmentsCount) {
		if err = d.Set("attachments_count", flex.IntValue(profile.AttachmentsCount)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting attachments_count: %s", err))
		}
	}

	return nil
}

func resourceIbmSccProfileUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceProfileOptions := &securityandcompliancecenterapiv3.ReplaceProfileOptions{}
	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	replaceProfileOptions.SetInstanceID(parts[0])
	replaceProfileOptions.SetProfileID(parts[1])
	hasChange := false
	bodyModelMap := map[string]interface{}{}

	if d.HasChange("controls") {
		hasChange = true
	}
	if d.HasChange("default_parameters") {
		hasChange = true
	}
	if d.HasChange("profile_name") {
		hasChange = true
	}
	if d.HasChange("profile_description") {
		hasChange = true
	}
	if d.HasChange("profile_version") {
		hasChange = true
	}
	if hasChange {
		if _, ok := d.GetOk("controls"); ok {
			bodyModelMap["controls"] = d.Get("controls")
		}
		if _, ok := d.GetOk("default_parameters"); ok {
			bodyModelMap["default_parameters"] = d.Get("default_parameters")
		}
		if _, ok := d.GetOk("profile_name"); ok {
			bodyModelMap["profile_name"] = d.Get("profile_name")
		}
		if _, ok := d.GetOk("profile_description"); ok {
			bodyModelMap["profile_description"] = d.Get("profile_description")
		}
		if _, ok := d.GetOk("profile_version"); ok {
			bodyModelMap["profile_version"] = d.Get("profile_version")
		}
		convertedModel, err := resourceIbmSccProfileMapToReplaceProfileOptions(bodyModelMap)
		if err != nil {
			return diag.FromErr(err)
		}

		replaceProfileOptions = convertedModel
		replaceProfileOptions.SetProfileID(d.Get("profile_id").(string))
		_, response, err := securityandcompliancecenterapiClient.ReplaceProfileWithContext(context, replaceProfileOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceProfileWithContext failed %s\n%s", err, response)
			return diag.FromErr(flex.FmtErrorf("ReplaceProfileWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmSccProfileRead(context, d, meta)
}

func resourceIbmSccProfileDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteCustomProfileOptions := &securityandcompliancecenterapiv3.DeleteCustomProfileOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	deleteCustomProfileOptions.SetInstanceID(parts[0])
	deleteCustomProfileOptions.SetProfileID(parts[1])

	_, response, err := securityandcompliancecenterapiClient.DeleteCustomProfileWithContext(context, deleteCustomProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteCustomProfileWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("DeleteCustomProfileWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmSccProfileMapToProfileControlsPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ProfileControlsPrototype, error) {
	model := &securityandcompliancecenterapiv3.ProfileControlsPrototype{}
	if modelMap["control_library_id"] != nil && modelMap["control_library_id"].(string) != "" {
		model.ControlLibraryID = core.StringPtr(modelMap["control_library_id"].(string))
	}
	if modelMap["control_id"] != nil && modelMap["control_id"].(string) != "" {
		model.ControlID = core.StringPtr(modelMap["control_id"].(string))
	}
	return model, nil
}

func resourceIbmSccProfileMapToDefaultParametersPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.DefaultParametersPrototype, error) {
	model := &securityandcompliancecenterapiv3.DefaultParametersPrototype{}
	if modelMap["assessment_type"] != nil && modelMap["assessment_type"].(string) != "" {
		model.AssessmentType = core.StringPtr(modelMap["assessment_type"].(string))
	}
	if modelMap["assessment_id"] != nil && modelMap["assessment_id"].(string) != "" {
		model.AssessmentID = core.StringPtr(modelMap["assessment_id"].(string))
	}
	if modelMap["parameter_name"] != nil && modelMap["parameter_name"].(string) != "" {
		model.ParameterName = core.StringPtr(modelMap["parameter_name"].(string))
	}
	if modelMap["parameter_default_value"] != nil && modelMap["parameter_default_value"].(string) != "" {
		model.ParameterDefaultValue = core.StringPtr(modelMap["parameter_default_value"].(string))
	}
	if modelMap["parameter_display_name"] != nil && modelMap["parameter_display_name"].(string) != "" {
		model.ParameterDisplayName = core.StringPtr(modelMap["parameter_display_name"].(string))
	}
	if modelMap["parameter_type"] != nil && modelMap["parameter_type"].(string) != "" {
		model.ParameterType = core.StringPtr(modelMap["parameter_type"].(string))
	}
	return model, nil
}

func resourceIbmSccProfileMapToProfilePrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.CreateProfileOptions, error) {
	model := &securityandcompliancecenterapiv3.CreateProfileOptions{}
	model.ProfileName = core.StringPtr(modelMap["profile_name"].(string))
	model.ProfileDescription = core.StringPtr(modelMap["profile_description"].(string))
	model.ProfileVersion = core.StringPtr(modelMap["profile_version"].(string))
	controls := []securityandcompliancecenterapiv3.ProfileControlsPrototype{}
	for _, controlsItem := range modelMap["controls"].([]interface{}) {
		controlsItemModel, err := resourceIbmSccProfileMapToProfileControlsPrototype(controlsItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		controls = append(controls, *controlsItemModel)
	}
	model.Controls = controls
	defaultParameters := []securityandcompliancecenterapiv3.DefaultParametersPrototype{}
	for _, defaultParametersItem := range modelMap["default_parameters"].([]interface{}) {
		if defaultParametersItem != nil {
			defaultParametersItemModel, err := resourceIbmSccProfileMapToDefaultParametersPrototype(defaultParametersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			defaultParameters = append(defaultParameters, *defaultParametersItemModel)
		}
	}
	model.DefaultParameters = defaultParameters
	// TODO: Validate all the Controls have default Parameters for any parameters found
	// Use the instance_id associated
	return model, nil
}

func resourceIbmSccProfileMapToReplaceProfileOptions(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ReplaceProfileOptions, error) {
	model := &securityandcompliancecenterapiv3.ReplaceProfileOptions{}
	model.ProfileName = core.StringPtr(modelMap["profile_name"].(string))
	model.ProfileDescription = core.StringPtr(modelMap["profile_description"].(string))
	model.ProfileType = core.StringPtr(modelMap["profile_type"].(string))
	model.ProfileVersion = core.StringPtr(modelMap["profile_version"].(string))
	controls := []securityandcompliancecenterapiv3.ProfileControlsPrototype{}
	for _, controlsItem := range modelMap["controls"].([]interface{}) {
		controlsItemModel, err := resourceIbmSccProfileMapToProfileControlsPrototype(controlsItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		controls = append(controls, *controlsItemModel)
	}
	model.Controls = controls
	defaultParameters := []securityandcompliancecenterapiv3.DefaultParametersPrototype{}
	for _, defaultParametersItem := range modelMap["default_parameters"].([]interface{}) {
		if defaultParametersItem != nil {
			defaultParametersItemModel, err := resourceIbmSccProfileMapToDefaultParametersPrototype(defaultParametersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			defaultParameters = append(defaultParameters, *defaultParametersItemModel)
		}
	}
	model.DefaultParameters = defaultParameters
	return model, nil
}

func resourceIbmSccProfileProfileControlsToMap(model *securityandcompliancecenterapiv3.ProfileControlsInResponse) (map[string]interface{}, error) {
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
		controlDocsMap, err := resourceIbmSccProfileControlDocsToMap(model.ControlDocs)
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
			controlSpecificationsItemMap, err := resourceIbmSccProfileControlSpecificationsToMap(&controlSpecificationsItem)
			if err != nil {
				return modelMap, err
			}
			controlSpecifications = append(controlSpecifications, controlSpecificationsItemMap)
		}
		modelMap["control_specifications"] = controlSpecifications
	}
	return modelMap, nil
}

func resourceIbmSccProfileControlDocsToMap(model *securityandcompliancecenterapiv3.ControlDoc) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ControlDocsID != nil {
		modelMap["control_docs_id"] = model.ControlDocsID
	}
	if model.ControlDocsType != nil {
		modelMap["control_docs_type"] = model.ControlDocsType
	}
	return modelMap, nil
}

func resourceIbmSccProfileControlSpecificationsToMap(model *securityandcompliancecenterapiv3.ControlSpecification) (map[string]interface{}, error) {
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
			assessmentsItemMap, err := resourceIbmSccProfileImplementationToMap(&assessmentsItem)
			if err != nil {
				return modelMap, err
			}
			assessments = append(assessments, assessmentsItemMap)
		}
		modelMap["assessments"] = assessments
	}
	return modelMap, nil
}

func resourceIbmSccProfileImplementationToMap(model *securityandcompliancecenterapiv3.Assessment) (map[string]interface{}, error) {
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
			parametersItemMap, err := resourceIbmSccProfileParameterInfoToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func resourceIbmSccProfileParameterInfoToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
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

func resourceIbmSccProfileDefaultParametersPrototypeToMap(model *securityandcompliancecenterapiv3.DefaultParameters) (map[string]interface{}, error) {
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
