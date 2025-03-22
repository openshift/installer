// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIbmSchematicsPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSchematicsPoliciesRead,

		Schema: map[string]*schema.Schema{
			"policy_kind": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy kind or categories for managing and deriving policy decision  * `agent_assignment_policy` Agent assignment policy for job execution.",
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of policy records.",
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of policy records returned.",
			},
			"offset": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The skipped number of policy records.",
			},
			"policies": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of Schematics policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of Schematics customization policy.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The system generated Policy Id.",
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
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of Schematics customization policy.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource-group name for the Policy.  By default, Policy will be created in Default Resource Group.",
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tags for the Schematics customization policy.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"location": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User defined status of the Schematics object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"state": &schema.Schema{
										Type:        schema.TypeString,
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
						"policy_kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy kind or categories for managing and deriving policy decision  * `agent_assignment_policy` Agent assignment policy for job execution.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy creation time.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user who created the Policy.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy updation time.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmSchematicsPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPoliciesRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listPolicyOptions := &schematicsv1.ListPolicyOptions{}

	policyList, response, err := schematicsClient.ListPolicyWithContext(context, listPolicyOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPoliciesRead ListPolicyWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchPolicies []schematicsv1.PolicyLite
	var policyKind string
	var suppliedFilter bool

	if v, ok := d.GetOk("policy_kind"); ok {
		policyKind = v.(string)
		suppliedFilter = true
		for _, data := range policyList.Policies {
			if data.PolicyKind != nil {
				if *data.PolicyKind == policyKind {
					matchPolicies = append(matchPolicies, data)
				}
			}
		}
	} else {
		matchPolicies = policyList.Policies
	}
	policyList.Policies = matchPolicies

	if suppliedFilter {
		if len(policyList.Policies) == 0 {

			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPoliciesRead failed with error: %s", err), "ibm_schematics_policies", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		d.SetId(policyKind)
	} else {
		d.SetId(dataSourceIbmSchematicsPoliciesID(d))
	}

	if err = d.Set("total_count", flex.IntValue(policyList.TotalCount)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPoliciesRead failed with error: %s", err), "ibm_schematics_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("limit", flex.IntValue(policyList.Limit)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPoliciesRead failed with error: %s", err), "ibm_schematics_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("offset", flex.IntValue(policyList.Offset)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPoliciesRead failed with error: %s", err), "ibm_schematics_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	policies := []map[string]interface{}{}
	if policyList.Policies != nil {
		for _, modelItem := range policyList.Policies {
			modelMap, err := dataSourceIbmSchematicsPoliciesPolicyLiteToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPoliciesRead failed: %s", err.Error()), "ibm_schematics_policies", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			policies = append(policies, modelMap)
		}
	}
	if err = d.Set("policies", policies); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPoliciesRead failed with error: %s", err), "ibm_schematics_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmSchematicsPoliciesID returns a reasonable ID for the list.
func dataSourceIbmSchematicsPoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSchematicsPoliciesPolicyLiteToMap(model *schematicsv1.PolicyLite) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.Account != nil {
		modelMap["account"] = *model.Account
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.ResourceGroup != nil {
		modelMap["resource_group"] = *model.ResourceGroup
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Location != nil {
		modelMap["location"] = *model.Location
	}
	if model.State != nil {
		stateMap, err := dataSourceIbmSchematicsPoliciesUserStateToMap(model.State)
		if err != nil {
			return modelMap, err
		}
		modelMap["state"] = []map[string]interface{}{stateMap}
	}
	if model.PolicyKind != nil {
		modelMap["policy_kind"] = *model.PolicyKind
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsPoliciesUserStateToMap(model *schematicsv1.UserState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.SetBy != nil {
		modelMap["set_by"] = *model.SetBy
	}
	if model.SetAt != nil {
		modelMap["set_at"] = model.SetAt.String()
	}
	return modelMap, nil
}
