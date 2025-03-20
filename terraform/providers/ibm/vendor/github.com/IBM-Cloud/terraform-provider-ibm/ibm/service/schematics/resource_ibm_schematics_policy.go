// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func ResourceIbmSchematicsPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSchematicsPolicyCreate,
		ReadContext:   resourceIbmSchematicsPolicyRead,
		UpdateContext: resourceIbmSchematicsPolicyUpdate,
		DeleteContext: resourceIbmSchematicsPolicyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Schematics customization policy.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of Schematics customization policy.",
			},
			"resource_group": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "The resource group name for the policy.  By default, Policy will be created in `default` Resource Group.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags for the Schematics customization policy.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "User defined status of the Schematics object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "User-defined states  * `draft` Object can be modified; can be used by Jobs run by the author, during execution  * `live` Object can be modified; can be used by Jobs during execution  * `locked` Object cannot be modified; can be used by Jobs during execution  * `disable` Object can be modified. cannot be used by Jobs during execution.",
						},
						"set_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the User who set the state of the Object.",
						},
						"set_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the User who set the state of the Object.",
						},
					},
				},
			},
			"kind": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_schematics_policy", "kind"),
				Description:  "Policy kind or categories for managing and deriving policy decision  * `agent_assignment_policy` Agent assignment policy for job execution.",
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The objects for the Schematics policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"selector_kind": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Types of schematics object selector.",
						},
						"selector_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Static selectors of schematics object ids (agent, workspace, action or blueprint) for the Schematics policy.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"selector_scope": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Selectors to dynamically list of schematics object ids (agent, workspace, action or blueprint) for the Schematics policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Name of the Schematics automation resource.",
									},
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The tag based selector.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"resource_groups": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The resource group based selector.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"locations": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The location based selector.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"parameter": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The parameter to tune the Schematics policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_assignment_policy_parameter": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Parameters for the `agent_assignment_policy`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selector_kind": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Types of schematics object selector.",
									},
									"selector_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The static selectors of schematics object ids (workspace, action or blueprint) for the Schematics policy.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"selector_scope": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "The selectors to dynamically list of schematics object ids (workspace, action or blueprint) for the Schematics policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Name of the Schematics automation resource.",
												},
												"tags": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "The tag based selector.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"resource_groups": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "The resource group based selector.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"locations": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "The location based selector.",
													Elem:        &schema.Schema{Type: schema.TypeString},
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
			"scoped_resources": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of scoped Schematics resources targeted by the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the Schematics automation resource.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Schematics resource Id.",
						},
					},
				},
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy CRN.",
			},
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Account id.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy creation time.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the policy.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy updation time.",
			},
		},
	}
}

func ResourceIbmSchematicsPolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "kind",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "agent_assignment_policy",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_schematics_policy", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSchematicsPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyCreate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_policy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createPolicyOptions := &schematicsv1.CreatePolicyOptions{}

	if _, ok := d.GetOk("name"); ok {
		createPolicyOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		createPolicyOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("resource_group"); ok {
		createPolicyOptions.SetResourceGroup(d.Get("resource_group").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		createPolicyOptions.SetTags(flex.ExpandStringList(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("location"); ok {
		createPolicyOptions.SetLocation(d.Get("location").(string))
	}
	if _, ok := d.GetOk("state"); ok {
		stateModel, err := resourceIbmSchematicsPolicyMapToUserState(d.Get("state.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyCreate failed: %s", err.Error()), "ibm_schematics_policy", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createPolicyOptions.SetState(stateModel)
	}
	if _, ok := d.GetOk("kind"); ok {
		createPolicyOptions.SetKind(d.Get("kind").(string))
	}
	if _, ok := d.GetOk("target"); ok {
		targetModel, err := resourceIbmSchematicsPolicyMapToPolicyObjects(d.Get("target.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyCreate failed: %s", err.Error()), "ibm_schematics_policy", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createPolicyOptions.SetTarget(targetModel)
	}
	if _, ok := d.GetOk("parameter"); ok {
		parameterModel, err := resourceIbmSchematicsPolicyMapToPolicyParameter(d.Get("parameter.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyCreate failed: %s", err.Error()), "ibm_schematics_policy", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		createPolicyOptions.SetParameter(parameterModel)
	}
	if _, ok := d.GetOk("scoped_resources"); ok {
		var scopedResources []schematicsv1.ScopedResource
		for _, e := range d.Get("scoped_resources").([]interface{}) {
			value := e.(map[string]interface{})
			scopedResourcesItem, err := resourceIbmSchematicsPolicyMapToScopedResource(value)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyCreate failed: %s", err.Error()), "ibm_schematics_policy", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			scopedResources = append(scopedResources, *scopedResourcesItem)
		}
		createPolicyOptions.SetScopedResources(scopedResources)
	}

	policy, response, err := schematicsClient.CreatePolicyWithContext(context, createPolicyOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyCreate CreatePolicyWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_policy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*policy.ID)

	return resourceIbmSchematicsPolicyRead(context, d, meta)
}

func resourceIbmSchematicsPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getPolicyOptions := &schematicsv1.GetPolicyOptions{}

	getPolicyOptions.SetPolicyID(d.Id())

	policy, response, err := schematicsClient.GetPolicyWithContext(context, getPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead GetPolicyWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", policy.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("description", policy.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("resource_group", policy.ResourceGroup); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if policy.Tags != nil {
		if err = d.Set("tags", policy.Tags); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("location", policy.Location); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if policy.State != nil {
		stateMap, err := resourceIbmSchematicsPolicyUserStateToMap(policy.State)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed: %s", err.Error()), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("state", []map[string]interface{}{stateMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("kind", policy.Kind); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if policy.Target != nil {
		targetMap, err := resourceIbmSchematicsPolicyPolicyObjectsToMap(policy.Target)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed: %s", err.Error()), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("target", []map[string]interface{}{targetMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if policy.Parameter != nil {
		parameterMap, err := resourceIbmSchematicsPolicyPolicyParameterToMap(policy.Parameter)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed: %s", err.Error()), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("parameter", []map[string]interface{}{parameterMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	scopedResources := []map[string]interface{}{}
	if policy.ScopedResources != nil {
		for _, scopedResourcesItem := range policy.ScopedResources {
			scopedResourcesItemMap, err := resourceIbmSchematicsPolicyScopedResourceToMap(&scopedResourcesItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed: %s", err.Error()), "ibm_schematics_policy", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			scopedResources = append(scopedResources, scopedResourcesItemMap)
		}
	}
	if err = d.Set("scoped_resources", scopedResources); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("crn", policy.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("account", policy.Account); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(policy.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", policy.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(policy.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmSchematicsPolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyUpdate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_policy", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updatePolicyOptions := &schematicsv1.UpdatePolicyOptions{}

	updatePolicyOptions.SetPolicyID(d.Id())

	hasChange := false

	if d.HasChange("name") {
		updatePolicyOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updatePolicyOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("resource_group") {
		updatePolicyOptions.SetResourceGroup(d.Get("resource_group").(string))
		hasChange = true
	}
	if d.HasChange("tags") {
		updatePolicyOptions.SetTags(flex.ExpandStringList(d.Get("tags").([]interface{})))
		hasChange = true
	}
	if d.HasChange("state") {
		state, err := resourceIbmSchematicsPolicyMapToUserState(d.Get("state.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyUpdate failed: %s", err.Error()), "ibm_schematics_policy", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		updatePolicyOptions.SetState(state)
		hasChange = true
	}
	if d.HasChange("kind") {
		updatePolicyOptions.SetKind(d.Get("kind").(string))
		hasChange = true
	}
	if d.HasChange("target") {
		target, err := resourceIbmSchematicsPolicyMapToPolicyObjects(d.Get("target.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyUpdate failed: %s", err.Error()), "ibm_schematics_policy", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		updatePolicyOptions.SetTarget(target)
		hasChange = true
	}
	if d.HasChange("parameter") {
		parameter, err := resourceIbmSchematicsPolicyMapToPolicyParameter(d.Get("parameter.0").(map[string]interface{}))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyUpdate failed: %s", err.Error()), "ibm_schematics_policy", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		updatePolicyOptions.SetParameter(parameter)
		hasChange = true
	}
	if d.HasChange("scoped_resources") {
		// TODO: handle ScopedResources of type TypeList -- not primitive, not model
		hasChange = true
	}

	if hasChange {
		_, response, err := schematicsClient.UpdatePolicyWithContext(context, updatePolicyOptions)
		if err != nil {

			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyUpdate UpdatePolicyWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_policy", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSchematicsPolicyRead(context, d, meta)
}

func resourceIbmSchematicsPolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyDelete schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_policy", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deletePolicyOptions := &schematicsv1.DeletePolicyOptions{}

	deletePolicyOptions.SetPolicyID(d.Id())

	response, err := schematicsClient.DeletePolicyWithContext(context, deletePolicyOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIbmSchematicsPolicyDelete DeletePolicyWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_policy", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSchematicsPolicyMapToUserState(modelMap map[string]interface{}) (*schematicsv1.UserState, error) {
	model := &schematicsv1.UserState{}
	if modelMap["state"] != nil && modelMap["state"].(string) != "" {
		model.State = core.StringPtr(modelMap["state"].(string))
	}
	if modelMap["set_by"] != nil && modelMap["set_by"].(string) != "" {
		model.SetBy = core.StringPtr(modelMap["set_by"].(string))
	}
	if modelMap["set_at"] != nil {

	}
	return model, nil
}

func resourceIbmSchematicsPolicyMapToPolicyObjects(modelMap map[string]interface{}) (*schematicsv1.PolicyObjects, error) {
	model := &schematicsv1.PolicyObjects{}
	if modelMap["selector_kind"] != nil && modelMap["selector_kind"].(string) != "" {
		model.SelectorKind = core.StringPtr(modelMap["selector_kind"].(string))
	}
	if modelMap["selector_ids"] != nil {
		selectorIds := []string{}
		for _, selectorIdsItem := range modelMap["selector_ids"].([]interface{}) {
			selectorIds = append(selectorIds, selectorIdsItem.(string))
		}
		model.SelectorIds = selectorIds
	}
	if modelMap["selector_scope"] != nil {
		selectorScope := []schematicsv1.PolicyObjectSelector{}
		for _, selectorScopeItem := range modelMap["selector_scope"].([]interface{}) {
			selectorScopeItemModel, err := resourceIbmSchematicsPolicyMapToPolicyObjectSelector(selectorScopeItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			selectorScope = append(selectorScope, *selectorScopeItemModel)
		}
		model.SelectorScope = selectorScope
	}
	return model, nil
}

func resourceIbmSchematicsPolicyMapToPolicyObjectSelector(modelMap map[string]interface{}) (*schematicsv1.PolicyObjectSelector, error) {
	model := &schematicsv1.PolicyObjectSelector{}
	if modelMap["kind"] != nil && modelMap["kind"].(string) != "" {
		model.Kind = core.StringPtr(modelMap["kind"].(string))
	}
	if modelMap["tags"] != nil {
		tags := []string{}
		for _, tagsItem := range modelMap["tags"].([]interface{}) {
			tags = append(tags, tagsItem.(string))
		}
		model.Tags = tags
	}
	if modelMap["resource_groups"] != nil {
		resourceGroups := []string{}
		for _, resourceGroupsItem := range modelMap["resource_groups"].([]interface{}) {
			resourceGroups = append(resourceGroups, resourceGroupsItem.(string))
		}
		model.ResourceGroups = resourceGroups
	}
	if modelMap["locations"] != nil {
		locations := []string{}
		for _, locationsItem := range modelMap["locations"].([]interface{}) {
			locations = append(locations, locationsItem.(string))
		}
		model.Locations = locations
	}
	return model, nil
}

func resourceIbmSchematicsPolicyMapToPolicyParameter(modelMap map[string]interface{}) (*schematicsv1.PolicyParameter, error) {
	model := &schematicsv1.PolicyParameter{}
	if modelMap["agent_assignment_policy_parameter"] != nil && len(modelMap["agent_assignment_policy_parameter"].([]interface{})) > 0 {
		AgentAssignmentPolicyParameterModel, err := resourceIbmSchematicsPolicyMapToAgentAssignmentPolicyParameter(modelMap["agent_assignment_policy_parameter"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AgentAssignmentPolicyParameter = AgentAssignmentPolicyParameterModel
	}
	return model, nil
}

func resourceIbmSchematicsPolicyMapToAgentAssignmentPolicyParameter(modelMap map[string]interface{}) (*schematicsv1.AgentAssignmentPolicyParameter, error) {
	model := &schematicsv1.AgentAssignmentPolicyParameter{}
	if modelMap["selector_kind"] != nil && modelMap["selector_kind"].(string) != "" {
		model.SelectorKind = core.StringPtr(modelMap["selector_kind"].(string))
	}
	if modelMap["selector_ids"] != nil {
		selectorIds := []string{}
		for _, selectorIdsItem := range modelMap["selector_ids"].([]interface{}) {
			selectorIds = append(selectorIds, selectorIdsItem.(string))
		}
		model.SelectorIds = selectorIds
	}
	if modelMap["selector_scope"] != nil {
		selectorScope := []schematicsv1.PolicyObjectSelector{}
		for _, selectorScopeItem := range modelMap["selector_scope"].([]interface{}) {
			selectorScopeItemModel, err := resourceIbmSchematicsPolicyMapToPolicyObjectSelector(selectorScopeItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			selectorScope = append(selectorScope, *selectorScopeItemModel)
		}
		model.SelectorScope = selectorScope
	}
	return model, nil
}

func resourceIbmSchematicsPolicyMapToScopedResource(modelMap map[string]interface{}) (*schematicsv1.ScopedResource, error) {
	model := &schematicsv1.ScopedResource{}
	if modelMap["kind"] != nil && modelMap["kind"].(string) != "" {
		model.Kind = core.StringPtr(modelMap["kind"].(string))
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIbmSchematicsPolicyUserStateToMap(model *schematicsv1.UserState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.State != nil {
		modelMap["state"] = model.State
	}
	if model.SetBy != nil {
		modelMap["set_by"] = model.SetBy
	}
	if model.SetAt != nil {
		modelMap["set_at"] = model.SetAt.String()
	}
	return modelMap, nil
}

func resourceIbmSchematicsPolicyPolicyObjectsToMap(model *schematicsv1.PolicyObjects) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SelectorKind != nil {
		modelMap["selector_kind"] = model.SelectorKind
	}
	if model.SelectorIds != nil {
		modelMap["selector_ids"] = model.SelectorIds
	}
	if model.SelectorScope != nil {
		selectorScope := []map[string]interface{}{}
		for _, selectorScopeItem := range model.SelectorScope {
			selectorScopeItemMap, err := resourceIbmSchematicsPolicyPolicyObjectSelectorToMap(&selectorScopeItem)
			if err != nil {
				return modelMap, err
			}
			selectorScope = append(selectorScope, selectorScopeItemMap)
		}
		modelMap["selector_scope"] = selectorScope
	}
	return modelMap, nil
}

func resourceIbmSchematicsPolicyPolicyObjectSelectorToMap(model *schematicsv1.PolicyObjectSelector) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Kind != nil {
		modelMap["kind"] = model.Kind
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.ResourceGroups != nil {
		modelMap["resource_groups"] = model.ResourceGroups
	}
	if model.Locations != nil {
		modelMap["locations"] = model.Locations
	}
	return modelMap, nil
}

func resourceIbmSchematicsPolicyPolicyParameterToMap(model *schematicsv1.PolicyParameter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentAssignmentPolicyParameter != nil {
		agentAssignmentPolicyParameterMap, err := resourceIbmSchematicsPolicyAgentAssignmentPolicyParameterToMap(model.AgentAssignmentPolicyParameter)
		if err != nil {
			return modelMap, err
		}
		modelMap["agent_assignment_policy_parameter"] = []map[string]interface{}{agentAssignmentPolicyParameterMap}
	}
	return modelMap, nil
}

func resourceIbmSchematicsPolicyAgentAssignmentPolicyParameterToMap(model *schematicsv1.AgentAssignmentPolicyParameter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SelectorKind != nil {
		modelMap["selector_kind"] = model.SelectorKind
	}
	if model.SelectorIds != nil {
		modelMap["selector_ids"] = model.SelectorIds
	}
	if model.SelectorScope != nil {
		selectorScope := []map[string]interface{}{}
		for _, selectorScopeItem := range model.SelectorScope {
			selectorScopeItemMap, err := resourceIbmSchematicsPolicyPolicyObjectSelectorToMap(&selectorScopeItem)
			if err != nil {
				return modelMap, err
			}
			selectorScope = append(selectorScope, selectorScopeItemMap)
		}
		modelMap["selector_scope"] = selectorScope
	}
	return modelMap, nil
}

func resourceIbmSchematicsPolicyScopedResourceToMap(model *schematicsv1.ScopedResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Kind != nil {
		modelMap["kind"] = model.Kind
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	return modelMap, nil
}
