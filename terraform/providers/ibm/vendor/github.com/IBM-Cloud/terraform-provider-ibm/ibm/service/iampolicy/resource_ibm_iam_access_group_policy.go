// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMAccessGroupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMIAMAccessGroupPolicyCreate,
		Read:   resourceIBMIAMAccessGroupPolicyRead,
		Update: resourceIBMIAMAccessGroupPolicyUpdate,
		Delete: resourceIBMIAMAccessGroupPolicyDelete,
		Exists: resourceIBMIAMAccessGroupPolicyExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				resources, resourceAttributes, err := importAccessGroupPolicy(d, meta)
				if err != nil {
					return nil, fmt.Errorf("[ERROR] Error reading resource ID: %s", err)
				}
				d.Set("resources", resources)
				d.Set("resource_attributes", resourceAttributes)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of access group",
				ForceNew:    true,
			},

			"roles": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Role names of the policy definition",
			},

			"resources": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"account_management", "resource_attributes"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service name of the policy definition",
						},

						"resource_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of resource instance of the policy definition",
						},

						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region of the policy definition",
						},

						"resource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource type of the policy definition",
						},

						"resource": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource of the policy definition",
						},

						"resource_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the resource group.",
						},

						"service_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service type of the policy definition",
						},

						"attributes": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Set resource attributes in the form of 'name=value,name=value....",
							Elem:        schema.TypeString,
						},
					},
				},
			},

			"resource_attributes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   "Set resource attributes.",
				ConflictsWith: []string{"resources", "account_management"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of attribute.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of attribute.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "stringEquals",
							Description: "Operator of attribute.",
						},
					},
				},
			},
			"account_management": {
				Type:          schema.TypeBool,
				Default:       false,
				Optional:      true,
				Description:   "Give access to all account management services",
				ConflictsWith: []string{"resources", "resource_attributes"},
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"resource_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set access management tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of attribute.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of attribute.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "stringEquals",
							Description: "Operator of attribute.",
						},
					},
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Policy",
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMIAMAccessGroupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	accessGroupId := d.Get("access_group_id").(string)

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	var policyOptions iampolicymanagementv1.CreatePolicyOptions
	policyOptions, err = flex.GeneratePolicyOptions(d, meta)
	if err != nil {
		return err
	}

	// Keep configuring the policy options by adding subject part
	accessGroupIdSubject := &iampolicymanagementv1.PolicySubject{
		Attributes: []iampolicymanagementv1.SubjectAttribute{
			{
				Name:  core.StringPtr("access_group_id"),
				Value: &accessGroupId,
			},
		},
	}

	accountIdResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
		Name:  core.StringPtr("accountId"),
		Value: &userDetails.UserAccount,
	}

	policyResource := &iampolicymanagementv1.PolicyResource{
		Attributes: append(policyOptions.Resources[0].Attributes, *accountIdResourceAttribute),
		Tags:       flex.SetTags(d),
	}

	createPolicyOptions := iamPolicyManagementClient.NewCreatePolicyOptions(
		"access",
		[]iampolicymanagementv1.PolicySubject{*accessGroupIdSubject},
		policyOptions.Roles,
		[]iampolicymanagementv1.PolicyResource{*policyResource},
	)

	if desc, ok := d.GetOk("description"); ok {
		des := desc.(string)
		createPolicyOptions.Description = &des
	}

	accessGroupPolicy, res, err := iamPolicyManagementClient.CreatePolicy(createPolicyOptions)
	if err != nil || accessGroupPolicy == nil {
		return fmt.Errorf("[ERROR] Error creating access group policy: %s\n%s", err, res)
	}

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: accessGroupPolicy.ID,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		policy, res, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
		if err != nil || policy == nil {
			if res != nil && res.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		_, res, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	}
	if err != nil {
		d.SetId(fmt.Sprintf("%s/%s", accessGroupId, *accessGroupPolicy.ID))
		return fmt.Errorf("[ERROR] Error fetching access group policy: %s\n%s", err, res)
	}
	d.SetId(fmt.Sprintf("%s/%s", accessGroupId, *accessGroupPolicy.ID))

	return resourceIBMIAMAccessGroupPolicyRead(d, meta)
}

func resourceIBMIAMAccessGroupPolicyRead(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	accessGroupId := parts[0]
	accessGroupPolicyId := parts[1]

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: &accessGroupPolicyId,
	}
	accessGroupPolicy := &iampolicymanagementv1.Policy{}
	res := &core.DetailedResponse{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		accessGroupPolicy, res, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)
		if err != nil || accessGroupPolicy == nil {
			if res != nil && res.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		accessGroupPolicy, res, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	}
	if err != nil || accessGroupPolicy == nil {
		return fmt.Errorf("[ERROR] Error retrieving access group policy: %s\n%s", err, res)
	}

	retrievedAttribute := flex.GetSubjectAttribute("access_group_id", accessGroupPolicy.Subjects[0])
	if accessGroupId != *retrievedAttribute {
		return fmt.Errorf("[ERROR] Policy %s does not belong to access group %s, retrievedAttr: %s", accessGroupPolicyId, accessGroupId, *retrievedAttribute)
	}

	d.Set("access_group_id", accessGroupId)
	roles := make([]string, len(accessGroupPolicy.Roles))
	for i, role := range accessGroupPolicy.Roles {
		roles[i] = *role.DisplayName
	}
	d.Set("roles", roles)
	d.Set("version", res.Headers.Get("ETag"))

	if _, ok := d.GetOk("resources"); ok {
		d.Set("resources", flex.FlattenPolicyResource(accessGroupPolicy.Resources))
	}
	if _, ok := d.GetOk("resource_attributes"); ok {
		d.Set("resource_attributes", flex.FlattenPolicyResourceAttributes(accessGroupPolicy.Resources))
	}

	if _, ok := d.GetOk("resource_tags"); ok {
		d.Set("resource_tags", flex.FlattenPolicyResourceTags(accessGroupPolicy.Resources))
	}

	if len(accessGroupPolicy.Resources) > 0 {
		if *flex.GetResourceAttribute("serviceType", accessGroupPolicy.Resources[0]) == "service" {
			d.Set("account_management", false)
		}
		if *flex.GetResourceAttribute("serviceType", accessGroupPolicy.Resources[0]) == "platform_service" {
			d.Set("account_management", true)
		}
	}

	if accessGroupPolicy.Description != nil {
		d.Set("description", *accessGroupPolicy.Description)
	}

	return nil
}

func resourceIBMIAMAccessGroupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}
	if d.HasChange("roles") || d.HasChange("resources") || d.HasChange("resource_attributes") || d.HasChange("account_management") || d.HasChange("description") || d.HasChange("resource_tags") {
		parts, err := flex.IdParts(d.Id())
		if err != nil {
			return err
		}
		accessGroupId := parts[0]
		accessGroupPolicyId := parts[1]

		userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}

		var policyOptions iampolicymanagementv1.CreatePolicyOptions
		policyOptions, err = flex.GeneratePolicyOptions(d, meta)
		if err != nil {
			return err
		}

		accessGroupIdSubject := &iampolicymanagementv1.PolicySubject{
			Attributes: []iampolicymanagementv1.SubjectAttribute{
				{
					Name:  core.StringPtr("access_group_id"),
					Value: &accessGroupId,
				},
			},
		}

		accountIdResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
			Name:  core.StringPtr("accountId"),
			Value: &userDetails.UserAccount,
		}

		policyResource := &iampolicymanagementv1.PolicyResource{
			Attributes: append(policyOptions.Resources[0].Attributes, *accountIdResourceAttribute),
			Tags:       flex.SetTags(d),
		}

		updatePolicyOptions := iamPolicyManagementClient.NewUpdatePolicyOptions(
			accessGroupPolicyId,
			d.Get("version").(string),
			"access",
			[]iampolicymanagementv1.PolicySubject{*accessGroupIdSubject},
			policyOptions.Roles,
			[]iampolicymanagementv1.PolicyResource{*policyResource},
		)

		if desc, ok := d.GetOk("description"); ok {
			des := desc.(string)
			updatePolicyOptions.Description = &des
		}

		_, res, err := iamPolicyManagementClient.UpdatePolicy(updatePolicyOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating access group policy: %s\n%s", err, res)
		}
	}

	return resourceIBMIAMAccessGroupPolicyRead(d, meta)
}

func resourceIBMIAMAccessGroupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	accessGroupPolicyId := parts[1]

	deletePolicyOptions := iamPolicyManagementClient.NewDeletePolicyOptions(
		accessGroupPolicyId,
	)

	res, err := iamPolicyManagementClient.DeletePolicy(deletePolicyOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting access group policy: %s\n%s", err, res)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMAccessGroupPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of accessGroupID/PolicyID", d.Id())
	}

	accessGroupPolicyId := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		accessGroupPolicyId,
	)

	accessGroupPolicy, resp, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil || accessGroupPolicy == nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting access group policy: %s\n%s", err, resp)
	}

	if accessGroupPolicy != nil && accessGroupPolicy.State != nil && *accessGroupPolicy.State == "deleted" {
		return false, nil
	}

	tempID := fmt.Sprintf("%s/%s", *flex.GetSubjectAttribute("access_group_id", accessGroupPolicy.Subjects[0]), *accessGroupPolicy.ID)

	return tempID == d.Id(), nil
}
func importAccessGroupPolicy(d *schema.ResourceData, meta interface{}) (interface{}, interface{}, error) {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return nil, nil, err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil, nil, err
	}
	accgrpPolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		accgrpPolicyID,
	)

	accessGroupPolicy, res, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] Error retrieving access group policy: %s\n%s", err, res)
	}

	resources := flex.FlattenPolicyResource(accessGroupPolicy.Resources)
	resource_attributes := flex.FlattenPolicyResourceAttributes(accessGroupPolicy.Resources)
	d.Set("resource_tags", flex.FlattenPolicyResourceTags(accessGroupPolicy.Resources))

	return resources, resource_attributes, nil
}
