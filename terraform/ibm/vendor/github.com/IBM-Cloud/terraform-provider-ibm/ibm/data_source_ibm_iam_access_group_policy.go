// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to find all the policies for an access group in a particular account
func dataSourceIBMIAMAccessGroupPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIAMAccessGroupPolicyRead,

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Description: "ID of access group",
				Type:        schema.TypeString,
				Required:    true,
			},
			"sort": {
				Description: "Sort query for policies",
				Type:        schema.TypeString,
				Optional:    true,
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
										Optional:    true,
										Description: "Service type of the policy definition",
									},
								},
							},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the Policy",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIAMAccessGroupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	accessGroupId := d.Get("access_group_id").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := userDetails.userAccount

	listPoliciesOptions := &iampolicymanagementv1.ListPoliciesOptions{
		AccountID:     core.StringPtr(accountID),
		AccessGroupID: core.StringPtr(accessGroupId),
		Type:          core.StringPtr("access"),
	}

	if v, ok := d.GetOk("sort"); ok {
		listPoliciesOptions.Sort = core.StringPtr(v.(string))
	}

	policyList, _, err := iamPolicyManagementClient.ListPolicies(listPoliciesOptions)
	policies := policyList.Policies

	if err != nil {
		return err
	}

	accessGroupPolicies := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		roles := make([]string, len(policy.Roles))
		for i, role := range policy.Roles {
			roles[i] = *role.DisplayName
		}
		resources := flattenPolicyResource(policy.Resources)
		p := map[string]interface{}{
			"id":        fmt.Sprintf("%s/%s", accessGroupId, *policy.ID),
			"roles":     roles,
			"resources": resources,
		}
		if policy.Description != nil {
			p["description"] = policy.Description
		}
		accessGroupPolicies = append(accessGroupPolicies, p)
	}
	d.SetId(accessGroupId)
	d.Set("policies", accessGroupPolicies)

	return nil
}
