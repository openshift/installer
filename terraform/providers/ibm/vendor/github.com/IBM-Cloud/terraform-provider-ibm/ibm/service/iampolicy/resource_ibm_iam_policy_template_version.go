// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

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
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
)

func ResourceIBMIAMPolicyTemplateVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIAMPolicyTemplateVersionCreate,
		ReadContext:   resourceIBMIAMPolicyTemplateVersionRead,
		UpdateContext: resourceIBMIAMPolicyTemplateVersionUpdate,
		DeleteContext: resourceIBMIAMPolicyTemplateVersionDelete,
		Exists:        resourceIBMIAMPolicyTemplateVersionExists,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy template ID and Version.",
			},
			"policy": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The core set of properties associated with the template's policy objet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The policy type; either 'access' or 'authorization'.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Allows the customer to use their own words to record the purpose/context related to a policy.",
						},
						"resource": {
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Optional:    true,
							Description: "The resource attributes to which the policy grants access.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attributes": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of resource attributes to which the policy grants access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of a resource attribute.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator of an attribute.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.",
												},
											},
										},
									},
									"tags": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Optional list of resource tags to which the policy grants access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of an access management tag.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The value of an access management tag.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator of an access management tag.",
												},
											},
										},
									},
								},
							},
						},
						"subject": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "The subject attributes for authorization type templates",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attributes": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of resource attributes to which the policy grants access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of a resource attribute.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator of an attribute.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.",
												},
											},
										},
									},
								},
							},
						},
						"pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates pattern of rule, either 'time-based-conditions:once', 'time-based-conditions:weekly:all-day', or 'time-based-conditions:weekly:custom-hours'.",
						},
						"rule_conditions": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Rule conditions enforced by the policy",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key of the condition",
									},
									"operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Operator of the condition",
									},
									"value": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Value of the condition",
									},
									"conditions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Additional Rule conditions enforced by the policy",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Key of the condition",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Operator of the condition",
												},
												"value": {
													Type:        schema.TypeList,
													Required:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Value of the condition",
												},
											},
										},
									},
								},
							},
						},
						"rule_operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator that multiple rule conditions are evaluated over",
						},
						"roles": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Role names of the policy definition",
						},
					},
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_policy_template_version", "description"),
				Description:  "description of template purpose.",
			},
			"committed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Template version committed status.",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template Version.",
			},
		},
	}
}

func ResourceIBMIAMPolicyTemplateVersionValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "template_id",
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			Required:                   true,
			MinValueLength:             1,
			MaxValueLength:             51,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: " ibm_iam_policy_template_version", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIAMPolicyTemplateVersionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createPolicyTemplateVersionOptions := &iampolicymanagementv1.CreatePolicyTemplateVersionOptions{}

	createPolicyTemplateVersionOptions.SetPolicyTemplateID(d.Get("template_id").(string))

	policyModel, err := generateTemplatePolicy(d, iamPolicyManagementClient)
	if err != nil {
		return diag.FromErr(err)
	}
	createPolicyTemplateVersionOptions.SetPolicy(policyModel)
	if _, ok := d.GetOk("description"); ok {
		createPolicyTemplateVersionOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("committed"); ok {
		createPolicyTemplateVersionOptions.SetCommitted(d.Get("committed").(bool))
	}
	if _, ok := d.GetOk("name"); ok {
		createPolicyTemplateVersionOptions.SetName(d.Get("name").(string))
	}

	policyTemplate, response, err := iamPolicyManagementClient.CreatePolicyTemplateVersion(createPolicyTemplateVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] CreatePolicyTemplateVersion failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreatePolicyTemplateVersion failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createPolicyTemplateVersionOptions.PolicyTemplateID, *policyTemplate.Version))

	return resourceIBMIAMPolicyTemplateVersionRead(context, d, meta)
}

func resourceIBMIAMPolicyTemplateVersionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getPolicyTemplateVersionOptions := &iampolicymanagementv1.GetPolicyTemplateVersionOptions{
		PolicyTemplateID: &parts[0],
		Version:          &parts[1],
	}

	policyTemplate, response, err := iamPolicyManagementClient.GetPolicyTemplateVersion(getPolicyTemplateVersionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetPolicyTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetPolicyTemplateWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("name", policyTemplate.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("account_id", policyTemplate.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}

	policyMap, err := flattenTemplatePolicy(policyTemplate.Policy, iamPolicyManagementClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("policy", []map[string]interface{}{policyMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting policy: %s", err))
	}
	if !core.IsNil(policyTemplate.Description) {
		if err = d.Set("description", policyTemplate.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if !core.IsNil(policyTemplate.Committed) {
		if err = d.Set("committed", policyTemplate.Committed); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting committed: %s", err))
		}
	}

	if !core.IsNil(policyTemplate.ID) {
		if err = d.Set("template_id", policyTemplate.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting committed: %s", err))
		}
	}

	if !core.IsNil(policyTemplate.Version) {
		if err = d.Set("version", policyTemplate.Version); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting committed: %s", err))
		}
	}

	return nil
}

func resourceIBMIAMPolicyTemplateVersionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("policy") || d.HasChange("description") || d.HasChange("committed") || d.HasChange("name") {
		iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
		if err != nil {
			return diag.FromErr(err)
		}

		replacePolicyTemplateOptions := &iampolicymanagementv1.ReplacePolicyTemplateOptions{}

		parts, err := flex.SepIdParts(d.Id(), "/")
		if err != nil {
			return diag.FromErr(err)
		}

		getPolicyTemplateVersionOptions := &iampolicymanagementv1.GetPolicyTemplateVersionOptions{
			PolicyTemplateID: &parts[0],
			Version:          &parts[1],
		}

		policyTemplate, response, err := iamPolicyManagementClient.GetPolicyTemplateVersionWithContext(context, getPolicyTemplateVersionOptions)

		if err != nil || policyTemplate == nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return diag.FromErr(fmt.Errorf("[ERROR] Error retrieving Policy Template: %s\n%s", err, response))
		}

		replacePolicyTemplateOptions.SetPolicyTemplateID(parts[0])
		replacePolicyTemplateOptions.SetVersion(parts[1])
		replacePolicyTemplateOptions.SetIfMatch(response.Headers.Get("ETag"))

		if description, ok := d.GetOk("description"); ok {
			replacePolicyTemplateOptions.SetDescription(description.(string))
		}

		if committed, ok := d.GetOk("committed"); ok {
			replacePolicyTemplateOptions.SetCommitted(committed.(bool))
		}

		if name, ok := d.GetOk("name"); ok {
			replacePolicyTemplateOptions.SetName(name.(string))
		}

		policy, err := generateTemplatePolicy(d, iamPolicyManagementClient)
		if err != nil {
			return diag.FromErr(err)
		}
		replacePolicyTemplateOptions.SetPolicy(policy)

		_, response, err = iamPolicyManagementClient.ReplacePolicyTemplateWithContext(context, replacePolicyTemplateOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplacePolicyTemplateWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ReplacePolicyTemplateWithContext failed %s\n%s", err, response))
		}
	}
	return resourceIBMIAMPolicyTemplateVersionRead(context, d, meta)
}

func resourceIBMIAMPolicyTemplateVersionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deletePolicyTemplateVersionOptions := &iampolicymanagementv1.DeletePolicyTemplateVersionOptions{
		PolicyTemplateID: &parts[0],
		Version:          &parts[1],
	}

	response, err := iamPolicyManagementClient.DeletePolicyTemplateVersionWithContext(context, deletePolicyTemplateVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] DeletePolicyTemplateVersion failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeletePolicyTemplateVersion failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMPolicyTemplateVersionExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of template id/version", d.Id())
	}

	getPolicyTemplateVersionOptions := &iampolicymanagementv1.GetPolicyTemplateVersionOptions{
		PolicyTemplateID: &parts[0],
		Version:          &parts[1],
	}

	policyTemplate, resp, err := iamPolicyManagementClient.GetPolicyTemplateVersion(getPolicyTemplateVersionOptions)

	if err != nil || policyTemplate == nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting policy template: %s\n%s", err, resp)
	}

	tempID := fmt.Sprintf("%s/%s", *policyTemplate.ID, *policyTemplate.Version)

	return tempID == d.Id(), nil
}
