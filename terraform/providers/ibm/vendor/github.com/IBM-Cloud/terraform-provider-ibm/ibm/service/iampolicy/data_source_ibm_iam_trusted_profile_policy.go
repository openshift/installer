// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
)

// Data source to find all the policies for a trusted profile
func DataSourceIBMIAMTrustedProfilePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIAMTrustedProfilePolicyRead,

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"profile_id", "iam_id"},
				Description:  "UUID of trusted profile",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile_policy",
					"profile_id"),
			},
			"iam_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"profile_id", "iam_id"},
				Description:  "IAM ID of trusted profile",
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
							Optional:    true,
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
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
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

func DataSourceIBMIAMTrustedProfilePolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:id"},
			Required:                   true})

	iBMIAMTrustedProfilePolicyValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_policy", Schema: validateSchema}
	return &iBMIAMTrustedProfilePolicyValidator
}

func dataSourceIBMIAMTrustedProfilePolicyRead(d *schema.ResourceData, meta interface{}) error {

	var iamID string
	if v, ok := d.GetOk("profile_id"); ok && v != nil {

		profileUUID := v.(string)
		iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
		if err != nil {
			return err
		}
		getprofileOptions := iamidentityv1.GetProfileOptions{
			ProfileID: &profileUUID,
		}
		profile, resp, err := iamClient.GetProfile(&getprofileOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting profile ID %s %s", err, resp)
		}
		iamID = *profile.IamID
	}
	if v, ok := d.GetOk("iam_id"); ok && v != nil {
		iamID = v.(string)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	listPoliciesOptions := &iampolicymanagementv1.ListV2PoliciesOptions{
		AccountID: core.StringPtr(userDetails.UserAccount),
		IamID:     core.StringPtr(iamID),
		Type:      core.StringPtr("access"),
	}

	if v, ok := d.GetOk("sort"); ok {
		listPoliciesOptions.Sort = core.StringPtr(v.(string))
	}

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		listPoliciesOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	policyList, resp, err := iamPolicyManagementClient.ListV2Policies(listPoliciesOptions)

	if err != nil || resp == nil {
		return fmt.Errorf("Error listing trusted profile policies: %s, %s", err, resp)
	}

	policies := policyList.Policies
	profilePolicies := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		roles, err := flex.GetRoleNamesFromPolicyResponse(policy, d, meta)
		if err != nil {
			return err
		}
		resources := flex.FlattenV2PolicyResource(*policy.Resource)
		p := map[string]interface{}{
			"roles":         roles,
			"resources":     resources,
			"resource_tags": flex.FlattenV2PolicyResourceTags(*policy.Resource),
		}
		if v, ok := d.GetOk("profile_id"); ok && v != nil {
			profileUUID := v.(string)
			p["id"] = fmt.Sprintf("%s/%s", profileUUID, *policy.ID)
		} else if v, ok := d.GetOk("iam_id"); ok && v != nil {
			iamID := v.(string)
			p["id"] = fmt.Sprintf("%s/%s", iamID, *policy.ID)
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
		profilePolicies = append(profilePolicies, p)
	}

	if v, ok := d.GetOk("profile_id"); ok && v != nil {
		profileUUID := v.(string)
		d.SetId(profileUUID)
	} else if v, ok := d.GetOk("iam_id"); ok && v != nil {
		iamID := v.(string)
		d.SetId(iamID)
	}
	if len(resp.Headers["Transaction-Id"]) > 0 && resp.Headers["Transaction-Id"][0] != "" {
		d.Set("transaction_id", resp.Headers["Transaction-Id"][0])
	}
	d.Set("policies", profilePolicies)
	return nil
}
