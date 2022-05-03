// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
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
					},
				},
			},
		},
	}
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

	listPoliciesOptions := &iampolicymanagementv1.ListPoliciesOptions{
		AccountID:     core.StringPtr(accountID),
		AccessGroupID: core.StringPtr(accessGroupId),
		Type:          core.StringPtr("access"),
	}

	if v, ok := d.GetOk("sort"); ok {
		listPoliciesOptions.Sort = core.StringPtr(v.(string))
	}

	policyList, resp, err := iamPolicyManagementClient.ListPolicies(listPoliciesOptions)
	if err != nil {
		return fmt.Errorf("Error listing access group policies: %s, %s", err, resp)
	}

	policies := policyList.Policies
	accessGroupPolicies := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		roles := make([]string, len(policy.Roles))
		for i, role := range policy.Roles {
			roles[i] = *role.DisplayName
		}
		resources := flex.FlattenPolicyResource(policy.Resources)
		p := map[string]interface{}{
			"id":            fmt.Sprintf("%s/%s", accessGroupId, *policy.ID),
			"roles":         roles,
			"resources":     resources,
			"resource_tags": flex.FlattenPolicyResourceTags(policy.Resources),
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
