// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source to find all the authorization policies in a particular account
func dataSourceIBMIAMAuthorizationPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMIAMAuthorizationPoliciesRead,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "The unique ID of an account",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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

						"source_service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source service name",
						},

						"target_service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target service name",
						},

						"roles": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Role names of the policy definition",
						},

						"source_resource_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source resource instance Id",
						},

						"target_resource_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target resource instance Id",
						},

						"source_resource_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source resource group Id",
						},

						"target_resource_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target resource group Id",
						},

						"source_resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type of source service",
						},

						"target_resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type of target service",
						},

						"source_service_account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account GUID of source service",
						},

						"version": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceIBMIAMAuthorizationPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	var accountID string

	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	if account, ok := d.GetOk("account_id"); ok && account.(string) != "" {
		accountID = account.(string)
	} else {
		userDetails, err := meta.(ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}
		accountID = userDetails.userAccount
	}

	listPoliciesOptions := &iampolicymanagementv1.ListPoliciesOptions{
		AccountID: core.StringPtr(accountID),
		Type:      core.StringPtr("authorization"),
	}

	if v, ok := d.GetOk("sort"); ok {
		listPoliciesOptions.Sort = core.StringPtr(v.(string))
	}

	policyList, resp, err := iamPolicyManagementClient.ListPolicies(listPoliciesOptions)

	if err != nil {
		return fmt.Errorf("Error listing authorization policies: %s, %s", err, resp)
	}

	policies := policyList.Policies

	authorizationPolicies := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		roles := make([]string, len(policy.Roles))
		for i, role := range policy.Roles {
			roles[i] = *role.DisplayName
		}
		source := policy.Subjects[0]
		target := policy.Resources[0]

		p := map[string]interface{}{
			"id":                          fmt.Sprintf("%s/%s", accountID, *policy.ID),
			"roles":                       roles,
			"source_service_name":         *getSubjectAttribute("serviceName", source),
			"target_service_name":         *getResourceAttribute("serviceName", target),
			"source_resource_instance_id": *getSubjectAttribute("serviceInstance", source),
			"target_resource_instance_id": *getResourceAttribute("serviceInstance", target),
			"source_resource_type":        *getSubjectAttribute("resourceType", source),
			"target_resource_type":        *getResourceAttribute("resourceType", target),
			"source_service_account":      *getSubjectAttribute("accountId", source),
			"source_resource_group_id":    *getSubjectAttribute("resourceGroupId", source),
			"target_resource_group_id":    *getResourceAttribute("resourceGroupId", target),
		}
		if policy.Description != nil {
			p["description"] = *policy.Description
		}
		authorizationPolicies = append(authorizationPolicies, p)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("account_id", accountID)
	d.Set("policies", authorizationPolicies)

	return nil
}
