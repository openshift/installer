// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
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

	listPoliciesOptions := &iampolicymanagementv1.ListPoliciesOptions{
		AccountID: core.StringPtr(userDetails.UserAccount),
		IamID:     core.StringPtr(iamID),
		Type:      core.StringPtr("access"),
	}

	if v, ok := d.GetOk("sort"); ok {
		listPoliciesOptions.Sort = core.StringPtr(v.(string))
	}

	policyList, resp, err := iamPolicyManagementClient.ListPolicies(listPoliciesOptions)
	if err != nil {
		return fmt.Errorf("Error listing trusted profile policies: %s, %s", err, resp)
	}

	policies := policyList.Policies
	profilePolicies := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		roles := make([]string, len(policy.Roles))
		for i, role := range policy.Roles {
			roles[i] = *role.DisplayName
		}
		resources := flex.FlattenPolicyResource(policy.Resources)
		p := map[string]interface{}{
			"roles":         roles,
			"resources":     resources,
			"resource_tags": flex.FlattenPolicyResourceTags(policy.Resources),
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
		profilePolicies = append(profilePolicies, p)
	}

	if v, ok := d.GetOk("profile_id"); ok && v != nil {
		profileUUID := v.(string)
		d.SetId(profileUUID)
	} else if v, ok := d.GetOk("iam_id"); ok && v != nil {
		iamID := v.(string)
		d.SetId(iamID)
	}
	d.Set("policies", profilePolicies)
	return nil
}
