// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to find all the policies for an access group in a particular account
func DataSourceIBMIAMAccessGroupPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIAMAccessGroupPolicyRead,

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Description: "ID of access group",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_access_group_policy",
					"access_group_id"),
			},
			"sort": {
				Description: "Sort query for policies",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"transaction_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Set transactionID for debug",
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roles": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Role names of the policy definition",
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service name of the policy definition",
									},
									"resource_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "ID of resource instance of the policy definition",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region of the policy definition",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource type of the policy definition",
									},
									"resource": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource of the policy definition",
									},
									"resource_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the resource group.",
									},
									"service_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service type of the policy definition",
									},
									"service_group_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Service group id of the policy definition",
									},
									"attributes": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "Set resource attributes in the form of 'name=value,name=value....",
										Elem:        schema.TypeString,
									},
								},
							},
						},
						"resource_tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Set access management tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of attribute.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value of attribute.",
									},
									"operator": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operator of attribute.",
									},
								},
							},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the Policy",
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
						"pattern": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Pattern rule follows for time-based condition",
						},
						"template": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Template meta data created from policy assignment",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy template id",
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy template version",
									},
									"assignment_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "policy assignment id",
									},
									"root_id": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "orchestrator template id",
									},
									"root_version": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "orchestrator template version",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMIAMAccessGroupPolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "access_group_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:access_group", "resolved_to:id"},
			Required:                   true})

	iBMIAMAccessGroupPolicyValidator := validate.ResourceValidator{ResourceName: "ibm_iam_access_group_policy", Schema: validateSchema}
	return &iBMIAMAccessGroupPolicyValidator
}

func dataSourceIBMIAMAccessGroupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	accessGroupId := d.Get("access_group_id").(string)

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := userDetails.UserAccount

	listPoliciesOptions := &iampolicymanagementv1.ListV2PoliciesOptions{
		AccountID:     core.StringPtr(accountID),
		AccessGroupID: core.StringPtr(accessGroupId),
		Type:          core.StringPtr("access"),
	}

	if v, ok := d.GetOk("sort"); ok {
		listPoliciesOptions.Sort = core.StringPtr(v.(string))
	}

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		listPoliciesOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	policyList, resp, err := iamPolicyManagementClient.ListV2Policies(listPoliciesOptions)

	if err != nil || resp == nil {
		return fmt.Errorf("Error listing access group policies: %s, %s", err, resp)
	}

	policies := policyList.Policies
	accessGroupPolicies := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		roles, err := flex.GetRoleNamesFromPolicyResponse(policy, d, meta)
		if err != nil {
			return err
		}
		resources := flex.FlattenV2PolicyResource(*policy.Resource)
		p := map[string]interface{}{
			"id":            fmt.Sprintf("%s/%s", accessGroupId, *policy.ID),
			"roles":         roles,
			"resources":     resources,
			"resource_tags": flex.FlattenV2PolicyResourceTags(*policy.Resource),
		}
		if policy.Description != nil {
			p["description"] = policy.Description
		}
		if policy.Rule != nil {
			p["rule_conditions"] = flex.FlattenRuleConditions(*policy.Rule.(*iampolicymanagementv1.V2PolicyRule))
			if len(policy.Rule.(*iampolicymanagementv1.V2PolicyRule).Conditions) > 0 {
				p["rule_operator"] = policy.Rule.(*iampolicymanagementv1.V2PolicyRule).Operator
			}
		}
		if policy.Pattern != nil {
			p["pattern"] = policy.Pattern
		}

		if policy.Template != nil {
			templateMap := flattenPolicyTemplateMetaData(policy.Template)
			p["template"] = []map[string]interface{}{templateMap}
		}

		accessGroupPolicies = append(accessGroupPolicies, p)
	}
	d.SetId(accessGroupId)

	if len(resp.Headers["Transaction-Id"]) > 0 && resp.Headers["Transaction-Id"][0] != "" {
		d.Set("transaction_id", resp.Headers["Transaction-Id"][0])
	}
	d.Set("policies", accessGroupPolicies)

	return nil
}

func flattenPolicyTemplateMetaData(model *iampolicymanagementv1.TemplateMetadata) map[string]interface{} {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.AssignmentID != nil {
		modelMap["assignment_id"] = model.AssignmentID
	}
	if model.RootID != nil {
		modelMap["root_id"] = model.RootID
	}
	if model.RootVersion != nil {
		modelMap["root_version"] = model.RootVersion
	}
	return modelMap
}
