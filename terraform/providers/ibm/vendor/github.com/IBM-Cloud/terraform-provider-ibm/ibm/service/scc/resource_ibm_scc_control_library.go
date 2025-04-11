// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func ResourceIbmSccControlLibrary() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		CreateContext: resourceIbmSccControlLibraryCreate,
		ReadContext:   resourceIbmSccControlLibraryRead,
		UpdateContext: resourceIbmSccControlLibraryUpdate,
		DeleteContext: resourceIbmSccControlLibraryDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"control_library_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The control library ID.",
			},
			"control_library_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_control_library", "control_library_name"),
				Description:  "The control library name.",
			},
			"control_library_description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_control_library", "control_library_description"),
				Description:  "The control library description.",
			},
			"control_library_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_control_library", "control_library_type"),
				Description:  "The control library type. This should be set to custom",
			},
			"version_group_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the version update of the control library.",
			},
			"control_library_version": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  func() (any, error) { return "0.0.0", nil },
				ValidateFunc: validate.InvokeValidator("ibm_scc_control_library", "control_library_version"),
				Description:  "The control library version.",
			},
			"latest": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The latest version of the control library.",
			},
			"controls_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of controls.",
			},
			"controls": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of controls in a control library.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"control_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The control name.",
						},
						"control_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the control library that contains the profile.",
						},
						"control_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The control description.",
						},
						"control_category": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The control category.",
						},
						"control_parent": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The parent control.",
						},
						"control_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The control tags.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"control_specifications": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The control specifications.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"control_specification_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The control specification ID.",
									},
									"responsibility": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The responsibility for managing the control.",
									},
									"component_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The component ID.",
									},
									"component_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The component name.",
									},
									"environment": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The control specifications environment.",
									},
									"control_specification_description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The control specifications description.",
									},
									"assessments_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of assessments.",
									},
									"assessments": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "The assessments.",
										Set:         assessmentsSchemaSetFunc("assessment_id"),
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"assessment_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The assessment ID.",
												},
												"assessment_method": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The assessment method.",
												},
												"assessment_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The assessment type.",
												},
												"assessment_description": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The assessment description.",
												},
												"parameter_count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The parameter count.",
												},
												"parameters": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"parameter_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The parameter name.",
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
											},
										},
									},
									"assessments_map": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"control_docs": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
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
						"control_requirement": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is this a control that can be automated or manually evaluated.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The control status.",
						},
					},
				},
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID.",
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
			"hierarchy_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The indication of whether hierarchy is enabled for the control library.",
			},
			"control_parents_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of parent controls in the control library.",
			},
		},
	})
}

func ResourceIbmSccControlLibraryValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "control_library_name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9_\s\-]*$`,
			MinValueLength:             2,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "control_library_description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[A-Za-z0-9]+`,
			MinValueLength:             2,
			MaxValueLength:             256,
		},
		validate.ValidateSchema{
			Identifier:                 "control_library_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "custom, predefined",
		},
		validate.ValidateSchema{
			Identifier:                 "version_group_label",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "control_library_version",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9_\-.]*$`,
			MinValueLength:             5,
			MaxValueLength:             64,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_control_library", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSccControlLibraryCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	bodyModelMap := map[string]interface{}{}
	createCustomControlLibraryOptions := &securityandcompliancecenterapiv3.CreateCustomControlLibraryOptions{}

	instance_id := d.Get("instance_id").(string)
	bodyModelMap["instance_id"] = instance_id
	bodyModelMap["control_library_name"] = d.Get("control_library_name")
	bodyModelMap["control_library_description"] = d.Get("control_library_description")
	bodyModelMap["control_library_type"] = d.Get("control_library_type")
	if _, ok := d.GetOk("version_group_label"); ok {
		bodyModelMap["version_group_label"] = d.Get("version_group_label")
	}
	if _, ok := d.GetOk("control_library_version"); ok {
		bodyModelMap["control_library_version"] = d.Get("control_library_version")
	}
	if _, ok := d.GetOk("latest"); ok {
		bodyModelMap["latest"] = d.Get("latest")
	}
	if _, ok := d.GetOk("controls_count"); ok {
		bodyModelMap["controls_count"] = d.Get("controls_count")
	}
	bodyModelMap["controls"] = d.Get("controls")

	convertedModel, err := resourceIbmSccControlLibraryMapToControlLibraryOptions(bodyModelMap)
	if err != nil {
		log.Printf("[DEBUG] CreateCustomControlLibraryWithContext failed %s\n", err)
		return diag.FromErr(flex.FmtErrorf("CreateCustomControlLibraryWithContext failed %s\n", err))
	}
	createCustomControlLibraryOptions = convertedModel
	controlLibrary, response, err := securityandcompliancecenterapiClient.CreateCustomControlLibraryWithContext(context, createCustomControlLibraryOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateCustomControlLibraryWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("CreateCustomControlLibraryWithContext failed %s\n%s", err, response))
	}

	d.SetId(instance_id + "/" + *controlLibrary.ID)

	return resourceIbmSccControlLibraryRead(context, d, meta)
}

func resourceIbmSccControlLibraryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getControlLibraryOptions := &securityandcompliancecenterapiv3.GetControlLibraryOptions{}
	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	getControlLibraryOptions.SetInstanceID(parts[0])
	getControlLibraryOptions.SetControlLibraryID(parts[1])

	controlLibrary, response, err := securityandcompliancecenterapiClient.GetControlLibraryWithContext(context, getControlLibraryOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetControlLibraryWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetControlLibraryWithContext failed %s\n%s", err, response))
	}
	if err = d.Set("instance_id", parts[0]); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("control_library_id", controlLibrary.ID); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting control_library_id: %s", err))
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
	if !core.IsNil(controlLibrary.VersionGroupLabel) {
		if err = d.Set("version_group_label", controlLibrary.VersionGroupLabel); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting version_group_label: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.ControlLibraryVersion) {
		if err = d.Set("control_library_version", controlLibrary.ControlLibraryVersion); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting control_library_version: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.Latest) {
		if err = d.Set("latest", controlLibrary.Latest); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting latest: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.ControlsCount) {
		if err = d.Set("controls_count", flex.IntValue(controlLibrary.ControlsCount)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting controls_count: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.Controls) {
		controls := []map[string]interface{}{}
		for _, controlsItem := range controlLibrary.Controls {
			controlsItemMap, err := resourceIbmSccControlLibraryControlsInControlLibToMap(&controlsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			controls = append(controls, controlsItemMap)
		}
		if err = d.Set("controls", controls); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting controls: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.AccountID) {
		if err = d.Set("account_id", controlLibrary.AccountID); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting account_id: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.CreatedOn) {
		if err = d.Set("created_on", flex.DateTimeToString(controlLibrary.CreatedOn)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting created_on: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.CreatedBy) {
		if err = d.Set("created_by", controlLibrary.CreatedBy); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting created_by: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.UpdatedOn) {
		if err = d.Set("updated_on", flex.DateTimeToString(controlLibrary.UpdatedOn)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting updated_on: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.UpdatedBy) {
		if err = d.Set("updated_by", controlLibrary.UpdatedBy); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting updated_by: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.HierarchyEnabled) {
		if err = d.Set("hierarchy_enabled", controlLibrary.HierarchyEnabled); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting hierarchy_enabled: %s", err))
		}
	}
	if !core.IsNil(controlLibrary.ControlParentsCount) {
		if err = d.Set("control_parents_count", flex.IntValue(controlLibrary.ControlParentsCount)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting control_parents_count: %s", err))
		}
	}

	return nil
}

func resourceIbmSccControlLibraryUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceCustomControlLibraryOptions := &securityandcompliancecenterapiv3.ReplaceCustomControlLibraryOptions{}
	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	replaceCustomControlLibraryOptions.SetInstanceID(parts[0])
	replaceCustomControlLibraryOptions.SetControlLibraryID(parts[1])

	hasChange := false

	if d.HasChange("control_library") {
		controlLibrary, err := resourceIbmSccControlLibraryMapToControlLibrary(d.Get("control_library.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceCustomControlLibraryOptions.SetControlLibraryID(*controlLibrary.ID)
		hasChange = true
	}

	if d.HasChange("controls") {
		for _, controlsItem := range d.Get("controls").([]interface{}) {
			controlsItemModel, err := resourceIbmSccControlLibraryMapToControl(controlsItem.(map[string]interface{}))
			if err != nil {
				return diag.FromErr(flex.FmtErrorf("ReplaceCustomControlLibraryWithContext failed %s\n", err))
			}
			replaceCustomControlLibraryOptions.Controls = append(replaceCustomControlLibraryOptions.Controls, *controlsItemModel)
		}
		hasChange = true
	}
	if d.HasChange("control_library_version") {
		replaceCustomControlLibraryOptions.SetControlLibraryVersion(d.Get("control_library_version").(string))
		hasChange = true
	}
	if d.HasChange("version_group_label") {
		replaceCustomControlLibraryOptions.SetVersionGroupLabel(d.Get("version_group_label").(string))
		hasChange = true
	}

	if hasChange {

		if replaceCustomControlLibraryOptions.ControlLibraryName == nil {
			replaceCustomControlLibraryOptions.SetControlLibraryName(d.Get("control_library_name").(string))
		}
		if replaceCustomControlLibraryOptions.VersionGroupLabel == nil {
			replaceCustomControlLibraryOptions.SetVersionGroupLabel(d.Get("version_group_label").(string))
		}
		if replaceCustomControlLibraryOptions.ControlLibraryDescription == nil {
			replaceCustomControlLibraryOptions.SetControlLibraryDescription(d.Get("control_library_description").(string))
		}
		if replaceCustomControlLibraryOptions.ControlLibraryVersion == nil {
			replaceCustomControlLibraryOptions.SetControlLibraryVersion(d.Get("control_library_version").(string))
		}
		if replaceCustomControlLibraryOptions.ControlLibraryType == nil {
			replaceCustomControlLibraryOptions.SetControlLibraryType("custom")
		}
		if len(replaceCustomControlLibraryOptions.Controls) == 0 {
			for _, controlsItem := range d.Get("controls").([]interface{}) {
				controlsItemModel, err := resourceIbmSccControlLibraryMapToControl(controlsItem.(map[string]interface{}))
				if err != nil {
					return diag.FromErr(flex.FmtErrorf("ReplaceCustomControlLibraryWithContext failed %s\n", err))
				}
				replaceCustomControlLibraryOptions.Controls = append(replaceCustomControlLibraryOptions.Controls, *controlsItemModel)
			}
		}
		_, response, err := securityandcompliancecenterapiClient.ReplaceCustomControlLibraryWithContext(context, replaceCustomControlLibraryOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceCustomControlLibraryWithContext failed %s\n%s", err, response)
			return diag.FromErr(flex.FmtErrorf("ReplaceCustomControlLibraryWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmSccControlLibraryRead(context, d, meta)
}

func resourceIbmSccControlLibraryDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityandcompliancecenterapiClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteCustomControlLibraryOptions := &securityandcompliancecenterapiv3.DeleteCustomControlLibraryOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	deleteCustomControlLibraryOptions.SetInstanceID(parts[0])
	deleteCustomControlLibraryOptions.SetControlLibraryID(parts[1])

	_, response, err := securityandcompliancecenterapiClient.DeleteCustomControlLibraryWithContext(context, deleteCustomControlLibraryOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteCustomControlLibraryWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("DeleteCustomControlLibraryWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmSccControlLibraryMapToControl(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Control, error) {
	model := &securityandcompliancecenterapiv3.Control{}
	if modelMap["control_name"] != nil && modelMap["control_name"].(string) != "" {
		model.ControlName = core.StringPtr(modelMap["control_name"].(string))
	}
	if modelMap["control_description"] != nil && modelMap["control_description"].(string) != "" {
		model.ControlDescription = core.StringPtr(modelMap["control_description"].(string))
	}
	if modelMap["control_id"] != nil && modelMap["control_id"].(string) != "" {
		model.ControlID = core.StringPtr(modelMap["control_id"].(string))
	}
	if modelMap["control_category"] != nil && modelMap["control_category"].(string) != "" {
		model.ControlCategory = core.StringPtr(modelMap["control_category"].(string))
	}
	if modelMap["control_parent"] != nil && modelMap["control_parent"].(string) != "" {
		model.ControlParent = core.StringPtr(modelMap["control_parent"].(string))
	}
	if modelMap["control_tags"] != nil {
		controlTags := []string{}
		for _, controlTagsItem := range modelMap["control_tags"].([]interface{}) {
			controlTags = append(controlTags, controlTagsItem.(string))
		}
		model.ControlTags = controlTags
	}
	if modelMap["control_specifications"] != nil {
		controlSpecifications := []securityandcompliancecenterapiv3.ControlSpecification{}
		for _, controlSpecificationsItem := range modelMap["control_specifications"].([]interface{}) {
			controlSpecificationsItemModel, err := resourceIbmSccControlLibraryMapToControlSpecification(controlSpecificationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			controlSpecifications = append(controlSpecifications, *controlSpecificationsItemModel)
		}
		model.ControlSpecifications = controlSpecifications
	}
	if modelMap["control_docs"].([]interface{})[0] != nil && len(modelMap["control_docs"].([]interface{})[0].(map[string]interface{})) > 0 {
		ControlDocsModel, err := resourceIbmSccControlLibraryMapToControlDocs(modelMap["control_docs"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ControlDocs = ControlDocsModel
	}
	if modelMap["control_requirement"] != nil {
		model.ControlRequirement = core.BoolPtr(modelMap["control_requirement"].(bool))
	}
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToControlPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Control, error) {
	model := &securityandcompliancecenterapiv3.Control{}
	if modelMap["control_name"] != nil && modelMap["control_name"].(string) != "" {
		model.ControlName = core.StringPtr(modelMap["control_name"].(string))
	}
	if modelMap["control_id"] != nil && modelMap["control_id"].(string) != "" {
		model.ControlID = core.StringPtr(modelMap["control_id"].(string))
	}
	if modelMap["control_description"] != nil && modelMap["control_description"].(string) != "" {
		model.ControlDescription = core.StringPtr(modelMap["control_description"].(string))
	}
	if modelMap["control_category"] != nil && modelMap["control_category"].(string) != "" {
		model.ControlCategory = core.StringPtr(modelMap["control_category"].(string))
	}
	if modelMap["control_parent"] != nil && modelMap["control_parent"].(string) != "" {
		model.ControlParent = core.StringPtr(modelMap["control_parent"].(string))
	}
	if modelMap["control_tags"] != nil {
		controlTags := []string{}
		for _, controlTagsItem := range modelMap["control_tags"].([]interface{}) {
			controlTags = append(controlTags, controlTagsItem.(string))
		}
		model.ControlTags = controlTags
	}
	if modelMap["control_specifications"] != nil {
		controlSpecifications := []securityandcompliancecenterapiv3.ControlSpecification{}
		for _, controlSpecificationsItem := range modelMap["control_specifications"].([]interface{}) {
			controlSpecificationsItemModel, err := resourceIbmSccControlLibraryMapToControlSpecificationPrototype(controlSpecificationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			controlSpecifications = append(controlSpecifications, *controlSpecificationsItemModel)
		}
		model.ControlSpecifications = controlSpecifications
	}
	if modelMap["control_docs"].([]interface{})[0] != nil && len(modelMap["control_docs"].([]interface{})[0].(map[string]interface{})) > 0 {
		ControlDocsModel, err := resourceIbmSccControlLibraryMapToControlDocs(modelMap["control_docs"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ControlDocs = ControlDocsModel
	}
	if modelMap["control_requirement"] != nil {
		model.ControlRequirement = core.BoolPtr(modelMap["control_requirement"].(bool))
	}
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToControlSpecificationPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ControlSpecification, error) {
	model := &securityandcompliancecenterapiv3.ControlSpecification{}
	if modelMap["control_specification_id"] != nil && modelMap["control_specification_id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["control_specification_id"].(string))
	}
	if modelMap["control_specification_description"] != nil && modelMap["control_specification_description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["control_specification_description"].(string))
	}
	if modelMap["component_id"] != nil && modelMap["component_id"].(string) != "" {
		model.ComponentID = core.StringPtr(modelMap["component_id"].(string))
	}
	if modelMap["environment"] != nil && modelMap["environment"].(string) != "" {
		model.Environment = core.StringPtr(modelMap["environment"].(string))
	}
	if modelMap["assessments"] != nil {
		assessments := []securityandcompliancecenterapiv3.Assessment{}
		for _, assessmentsItem := range modelMap["assessments"].(*schema.Set).List() {
			assessmentsItemModel, err := resourceIbmSccControlLibraryMapToAssessmentPrototype(assessmentsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			assessments = append(assessments, *assessmentsItemModel)
		}
		model.Assessments = assessments
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToControlSpecification(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ControlSpecification, error) {
	model := &securityandcompliancecenterapiv3.ControlSpecification{}
	if modelMap["control_specification_id"] != nil && modelMap["control_specification_id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["control_specification_id"].(string))
	}
	if modelMap["responsibility"] != nil && modelMap["responsibility"].(string) != "" {
		model.Responsibility = core.StringPtr(modelMap["responsibility"].(string))
	}
	if modelMap["component_id"] != nil && modelMap["component_id"].(string) != "" {
		model.ComponentID = core.StringPtr(modelMap["component_id"].(string))
	}
	if modelMap["component_name"] != nil && modelMap["component_name"].(string) != "" {
		model.ComponentName = core.StringPtr(modelMap["component_name"].(string))
	}
	if modelMap["environment"] != nil && modelMap["environment"].(string) != "" {
		model.Environment = core.StringPtr(modelMap["environment"].(string))
	}
	if modelMap["control_specification_description"] != nil && modelMap["control_specification_description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["control_specification_description"].(string))
	}
	if modelMap["assessments_count"] != nil {
		model.AssessmentsCount = core.Int64Ptr(int64(modelMap["assessments_count"].(int)))
	}
	if modelMap["assessments"] != nil {
		assessments := []securityandcompliancecenterapiv3.Assessment{}
		for _, assessmentsItem := range modelMap["assessments"].(*schema.Set).List() {
			assessmentsItemModel, err := resourceIbmSccControlLibraryMapToAssessment(assessmentsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			assessments = append(assessments, *assessmentsItemModel)
		}
		model.Assessments = assessments
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToAssessment(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Assessment, error) {
	model := &securityandcompliancecenterapiv3.Assessment{}
	if modelMap["assessment_id"] != nil && modelMap["assessment_id"].(string) != "" {
		model.AssessmentID = core.StringPtr(modelMap["assessment_id"].(string))
	}
	if modelMap["assessment_method"] != nil && modelMap["assessment_method"].(string) != "" {
		model.AssessmentMethod = core.StringPtr(modelMap["assessment_method"].(string))
	}
	if modelMap["assessment_type"] != nil && modelMap["assessment_type"].(string) != "" {
		model.AssessmentType = core.StringPtr(modelMap["assessment_type"].(string))
	}
	if modelMap["assessment_description"] != nil && modelMap["assessment_description"].(string) != "" {
		model.AssessmentDescription = core.StringPtr(modelMap["assessment_description"].(string))
	}
	if modelMap["parameter_count"] != nil {
		model.ParameterCount = core.Int64Ptr(int64(modelMap["parameter_count"].(int)))
	}
	if modelMap["parameters"] != nil {
		parameters := []securityandcompliancecenterapiv3.Parameter{}
		for _, parametersItem := range modelMap["parameters"].([]interface{}) {
			if parametersItem != nil {
				parametersItemModel, err := resourceIbmSccControlLibraryMapToParameterInfo(parametersItem.(map[string]interface{}))
				if err != nil {
					return model, err
				}
				parameters = append(parameters, *parametersItemModel)
			}
		}
		model.Parameters = parameters
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToAssessmentPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Assessment, error) {
	model := &securityandcompliancecenterapiv3.Assessment{}
	if modelMap["assessment_id"] != nil && modelMap["assessment_id"].(string) != "" {
		model.AssessmentID = core.StringPtr(modelMap["assessment_id"].(string))
	}
	if modelMap["assessment_description"] != nil && modelMap["assessment_description"].(string) != "" {
		model.AssessmentDescription = core.StringPtr(modelMap["assessment_description"].(string))
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToParameterInfo(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Parameter, error) {
	model := &securityandcompliancecenterapiv3.Parameter{}
	if modelMap["parameter_name"] != nil && modelMap["parameter_name"].(string) != "" {
		model.ParameterName = core.StringPtr(modelMap["parameter_name"].(string))
	}
	if modelMap["parameter_display_name"] != nil && modelMap["parameter_display_name"].(string) != "" {
		model.ParameterDisplayName = core.StringPtr(modelMap["parameter_display_name"].(string))
	}
	if modelMap["parameter_type"] != nil && modelMap["parameter_type"].(string) != "" {
		model.ParameterType = core.StringPtr(modelMap["parameter_type"].(string))
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToControlDocs(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ControlDoc, error) {
	model := &securityandcompliancecenterapiv3.ControlDoc{}
	if modelMap["control_docs_id"] != nil && modelMap["control_docs_id"].(string) != "" {
		model.ControlDocsID = core.StringPtr(modelMap["control_docs_id"].(string))
	}
	if modelMap["control_docs_type"] != nil && modelMap["control_docs_type"].(string) != "" {
		model.ControlDocsType = core.StringPtr(modelMap["control_docs_type"].(string))
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToControlLibrary(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ControlLibrary, error) {
	model := &securityandcompliancecenterapiv3.ControlLibrary{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["account_id"] != nil && modelMap["account_id"].(string) != "" {
		model.AccountID = core.StringPtr(modelMap["account_id"].(string))
	}
	if modelMap["control_library_name"] != nil && modelMap["control_library_name"].(string) != "" {
		model.ControlLibraryName = core.StringPtr(modelMap["control_library_name"].(string))
	}
	if modelMap["control_library_description"] != nil && modelMap["control_library_description"].(string) != "" {
		model.ControlLibraryDescription = core.StringPtr(modelMap["control_library_description"].(string))
	}
	if modelMap["control_library_type"] != nil && modelMap["control_library_type"].(string) != "" {
		model.ControlLibraryType = core.StringPtr(modelMap["control_library_type"].(string))
	}
	if modelMap["version_group_label"] != nil && modelMap["version_group_label"].(string) != "" {
		model.VersionGroupLabel = core.StringPtr(modelMap["version_group_label"].(string))
	}
	if modelMap["control_library_version"] != nil && modelMap["control_library_version"].(string) != "" {
		model.ControlLibraryVersion = core.StringPtr(modelMap["control_library_version"].(string))
	}
	if modelMap["created_on"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["created_on"].(string))
		if err != nil {
			return model, err
		}
		model.CreatedOn = &dateTime
	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["updated_on"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["updated_on"].(string))
		if err != nil {
			return model, err
		}
		model.UpdatedOn = &dateTime
	}
	if modelMap["updated_by"] != nil && modelMap["updated_by"].(string) != "" {
		model.UpdatedBy = core.StringPtr(modelMap["updated_by"].(string))
	}
	if modelMap["latest"] != nil {
		model.Latest = core.BoolPtr(modelMap["latest"].(bool))
	}
	if modelMap["hierarchy_enabled"] != nil {
		model.HierarchyEnabled = core.BoolPtr(modelMap["hierarchy_enabled"].(bool))
	}
	if modelMap["controls_count"] != nil {
		model.ControlsCount = core.Int64Ptr(int64(modelMap["controls_count"].(int)))
	}
	if modelMap["control_parents_count"] != nil {
		model.ControlParentsCount = core.Int64Ptr(int64(modelMap["control_parents_count"].(int)))
	}
	if modelMap["controls"] != nil {
		controls := []securityandcompliancecenterapiv3.Control{}
		for _, controlsItem := range modelMap["controls"].([]interface{}) {
			controlsItemModel, err := resourceIbmSccControlLibraryMapToControl(controlsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			controls = append(controls, *controlsItemModel)
		}
		model.Controls = controls
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToControlLibraryPrototype(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.ControlLibrary, error) {
	model := &securityandcompliancecenterapiv3.ControlLibrary{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["account_id"] != nil && modelMap["account_id"].(string) != "" {
		model.AccountID = core.StringPtr(modelMap["account_id"].(string))
	}
	if modelMap["control_library_name"] != nil && modelMap["control_library_name"].(string) != "" {
		model.ControlLibraryName = core.StringPtr(modelMap["control_library_name"].(string))
	}
	if modelMap["control_library_description"] != nil && modelMap["control_library_description"].(string) != "" {
		model.ControlLibraryDescription = core.StringPtr(modelMap["control_library_description"].(string))
	}
	if modelMap["control_library_type"] != nil && modelMap["control_library_type"].(string) != "" {
		model.ControlLibraryType = core.StringPtr(modelMap["control_library_type"].(string))
	}
	if modelMap["version_group_label"] != nil && modelMap["version_group_label"].(string) != "" {
		model.VersionGroupLabel = core.StringPtr(modelMap["version_group_label"].(string))
	}
	if modelMap["control_library_version"] != nil && modelMap["control_library_version"].(string) != "" {
		model.ControlLibraryVersion = core.StringPtr(modelMap["control_library_version"].(string))
	}
	if modelMap["created_on"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["created_on"].(string))
		if err != nil {
			return model, err
		}
		model.CreatedOn = &dateTime
	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["updated_on"] != nil {
		dateTime, err := core.ParseDateTime(modelMap["updated_on"].(string))
		if err != nil {
			return model, err
		}
		model.UpdatedOn = &dateTime
	}
	if modelMap["updated_by"] != nil && modelMap["updated_by"].(string) != "" {
		model.UpdatedBy = core.StringPtr(modelMap["updated_by"].(string))
	}
	if modelMap["latest"] != nil {
		model.Latest = core.BoolPtr(modelMap["latest"].(bool))
	}
	if modelMap["hierarchy_enabled"] != nil {
		model.HierarchyEnabled = core.BoolPtr(modelMap["hierarchy_enabled"].(bool))
	}
	if modelMap["controls_count"] != nil {
		model.ControlsCount = core.Int64Ptr(int64(modelMap["controls_count"].(int)))
	}
	if modelMap["control_parents_count"] != nil {
		model.ControlParentsCount = core.Int64Ptr(int64(modelMap["control_parents_count"].(int)))
	}
	if modelMap["controls"] != nil {
		controls := []securityandcompliancecenterapiv3.Control{}
		for _, controlsItem := range modelMap["controls"].([]interface{}) {
			controlsItemModel, err := resourceIbmSccControlLibraryMapToControl(controlsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			controls = append(controls, *controlsItemModel)
		}
		model.Controls = controls
	}
	return model, nil
}

func resourceIbmSccControlLibraryMapToControlLibraryOptions(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.CreateCustomControlLibraryOptions, error) {
	model := &securityandcompliancecenterapiv3.CreateCustomControlLibraryOptions{}
	model.InstanceID = core.StringPtr(modelMap["instance_id"].(string))
	model.ControlLibraryName = core.StringPtr(modelMap["control_library_name"].(string))
	model.ControlLibraryDescription = core.StringPtr(modelMap["control_library_description"].(string))
	model.ControlLibraryType = core.StringPtr(modelMap["control_library_type"].(string))
	if modelMap["version_group_label"] != nil && modelMap["version_group_label"].(string) != "" {
		model.VersionGroupLabel = core.StringPtr(modelMap["version_group_label"].(string))
	}
	if modelMap["control_library_version"] != nil && modelMap["control_library_version"].(string) != "" {
		model.ControlLibraryVersion = core.StringPtr(modelMap["control_library_version"].(string))
	}
	if modelMap["latest"] != nil {
		model.Latest = core.BoolPtr(modelMap["latest"].(bool))
	}
	controls := []securityandcompliancecenterapiv3.Control{}
	// iterate through controls
	for _, controlsItem := range modelMap["controls"].([]interface{}) {
		controlsItemModel, err := resourceIbmSccControlLibraryMapToControlPrototype(controlsItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		controls = append(controls, *controlsItemModel)
	}
	model.Controls = controls
	return model, nil
}

// From the GET function, popularize the controls of a Control Library
func resourceIbmSccControlLibraryControlsInControlLibToMap(model *securityandcompliancecenterapiv3.Control) (map[string]interface{}, error) {
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
			controlSpecificationsItemMap, err := resourceIbmSccControlLibraryControlSpecificationsToMap(&controlSpecificationsItem)
			if err != nil {
				return modelMap, err
			}
			controlSpecifications = append(controlSpecifications, controlSpecificationsItemMap)
		}
		modelMap["control_specifications"] = controlSpecifications
	}
	if model.ControlDocs != nil {
		controlDocsMap, err := resourceIbmSccControlLibraryControlDocsToMap(model.ControlDocs)
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

// using the assessment_id for comparison
func compareAssessmentSetFunc(v interface{}) int {
	if v == nil {
		return 0
	}
	m := v.(map[string]interface{})
	id := (m["assessment_id"]).(*string)
	assId := (*id)[5:18]
	var i big.Int
	i.SetString(strings.Replace(assId, "-", "", 4), 16)
	val, err := strconv.Atoi(i.String())
	if err != nil {
		log.Printf("[ERROR] Setting the Assessments for Control Library failed %s\n", err)
	}
	return val
}

// From the GET of a control library, popularize the control specifications in a control
func resourceIbmSccControlLibraryControlSpecificationsToMap(model *securityandcompliancecenterapiv3.ControlSpecification) (map[string]interface{}, error) {
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
		assessmentSet := &schema.Set{
			F: assessmentsSchemaSetFunc("assessment_id"),
		}
		for _, assessmentsItem := range model.Assessments {
			assessmentsItemMap, err := resourceIbmSccControlLibraryImplementationToMap(&assessmentsItem)
			if err != nil {
				return modelMap, err
			}
			assessmentSet.Add(assessmentsItemMap)
		}
		modelMap["assessments"] = assessmentSet
	}
	return modelMap, nil
}

// assessmentsSchemaSetFunc determines how to hash the Assessments schema.Resource
// It uses the assessment_id in order to determine if the assessment was present
func assessmentsSchemaSetFunc(keys ...string) schema.SchemaSetFunc {
	return func(v interface{}) int {
		var str strings.Builder

		if m, ok := v.(map[string]interface{}); ok {
			for _, key := range keys {
				if v, ok := m[key]; ok {
					switch v := v.(type) {
					case string:
						str.WriteRune('-')
						str.WriteString(v)
					case *string:
						str.WriteRune('-')
						str.WriteString(*v)
					case []interface{}:
						str.WriteRune('-')
						s := make([]string, len(v))
						for i, v := range v {
							s[i] = fmt.Sprint(v)
						}
						str.WriteString(fmt.Sprintf("[%s]", strings.Join(s, ",")))
					}
				}
			}
		}
		log.Printf("[DEBUG] assessmentsSchemaSet hashcode string: %s\n", str.String())

		return stringHashcode(str.String())
	}
}

func resourceIbmSccControlLibraryImplementationToMap(model *securityandcompliancecenterapiv3.Assessment) (map[string]interface{}, error) {
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
			parametersItemMap, err := resourceIbmSccControlLibraryParameterInfoToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func resourceIbmSccControlLibraryParameterInfoToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
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

func resourceIbmSccControlLibraryControlDocsToMap(model *securityandcompliancecenterapiv3.ControlDoc) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ControlDocsID != nil {
		modelMap["control_docs_id"] = model.ControlDocsID
	}
	if model.ControlDocsType != nil {
		modelMap["control_docs_type"] = model.ControlDocsType
	}
	return modelMap, nil
}
