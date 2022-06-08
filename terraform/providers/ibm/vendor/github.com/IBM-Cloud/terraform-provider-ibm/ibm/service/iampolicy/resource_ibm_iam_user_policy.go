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

func ResourceIBMIAMUserPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMIAMUserPolicyCreate,
		Read:   resourceIBMIAMUserPolicyRead,
		Update: resourceIBMIAMUserPolicyUpdate,
		Delete: resourceIBMIAMUserPolicyDelete,
		Exists: resourceIBMIAMUserPolicyExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				resources, resourceAttributes, err := importUserPolicy(d, meta)
				if err != nil {
					return nil, fmt.Errorf("[ERROR] Error reading resource ID: %s", err)
				}
				d.Set("resources", resources)
				d.Set("resource_attributes", resourceAttributes)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{

			"ibm_id": {
				Description: "The ibm id or email of user",
				Type:        schema.TypeString,
				Required:    true,
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
		},
	}
}

func resourceIBMIAMUserPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	var policyOptions iampolicymanagementv1.CreatePolicyOptions
	policyOptions, err = flex.GeneratePolicyOptions(d, meta)

	if err != nil {
		return err
	}

	userEmail := d.Get("ibm_id").(string)

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := userDetails.UserAccount

	ibmUniqueID, err := flex.GetIBMUniqueId(accountID, userEmail, meta)
	if err != nil {
		return err
	}

	subjectAttribute := &iampolicymanagementv1.SubjectAttribute{
		Name:  core.StringPtr("iam_id"),
		Value: &ibmUniqueID,
	}

	policySubjects := &iampolicymanagementv1.PolicySubject{
		Attributes: []iampolicymanagementv1.SubjectAttribute{*subjectAttribute},
	}

	accountIDResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
		Name:     core.StringPtr("accountId"),
		Value:    core.StringPtr(accountID),
		Operator: core.StringPtr("stringEquals"),
	}

	policyResources := iampolicymanagementv1.PolicyResource{
		Attributes: append(policyOptions.Resources[0].Attributes, *accountIDResourceAttribute),
		Tags:       flex.SetTags(d),
	}

	createPolicyOptions := iamPolicyManagementClient.NewCreatePolicyOptions(
		"access",
		[]iampolicymanagementv1.PolicySubject{*policySubjects},
		policyOptions.Roles,
		[]iampolicymanagementv1.PolicyResource{policyResources},
	)

	if description, ok := d.GetOk("description"); ok {
		des := description.(string)
		createPolicyOptions.Description = &des
	}

	userPolicy, resp, err := iamPolicyManagementClient.CreatePolicy(createPolicyOptions)

	if err != nil {
		return fmt.Errorf("Error creating user policies: %s, %s", err, resp)
	}

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: userPolicy.ID,
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
		_, _, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	}
	if err != nil {
		d.SetId(fmt.Sprintf("%s/%s", userEmail, *userPolicy.ID))
		return fmt.Errorf("[ERROR] Error fetching user policy: %w", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", userEmail, *userPolicy.ID))

	return resourceIBMIAMUserPolicyRead(d, meta)
}

func resourceIBMIAMUserPolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	userEmail := parts[0]
	userPolicyID := parts[1]

	if err != nil {
		return err
	}

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: core.StringPtr(userPolicyID),
	}
	userPolicy := &iampolicymanagementv1.Policy{}
	res := &core.DetailedResponse{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		userPolicy, res, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)

		if err != nil || userPolicy == nil {
			if res != nil && res.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		userPolicy, res, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	}
	if err != nil || userPolicy == nil {
		return fmt.Errorf("[ERROR] Error retrieving userPolicy: %s %s", err, res)
	}
	d.Set("ibm_id", userEmail)
	roles := make([]string, len(userPolicy.Roles))

	for i, role := range userPolicy.Roles {
		roles[i] = *role.DisplayName
	}
	d.Set("roles", roles)
	if _, ok := d.GetOk("resources"); ok {
		d.Set("resources", flex.FlattenPolicyResource(userPolicy.Resources))
	}
	if _, ok := d.GetOk("resource_attributes"); ok {
		d.Set("resource_attributes", flex.FlattenPolicyResourceAttributes(userPolicy.Resources))
	}

	if _, ok := d.GetOk("resource_tags"); ok {
		d.Set("resource_tags", flex.FlattenPolicyResourceTags(userPolicy.Resources))
	}

	if len(userPolicy.Resources) > 0 {
		if *flex.GetResourceAttribute("serviceType", userPolicy.Resources[0]) == "service" {
			d.Set("account_management", false)
		}
		if *flex.GetResourceAttribute("serviceType", userPolicy.Resources[0]) == "platform_service" {
			d.Set("account_management", true)
		}
	}
	if userPolicy.Description != nil {
		d.Set("description", *userPolicy.Description)
	}
	return nil
}

func resourceIBMIAMUserPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}
	if d.HasChange("roles") || d.HasChange("resources") || d.HasChange("resource_attributes") || d.HasChange("account_management") || d.HasChange("description") || d.HasChange("resource_tags") {
		parts, err := flex.IdParts(d.Id())
		if err != nil {
			return err
		}
		userEmail := parts[0]
		userPolicyID := parts[1]

		userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}

		accountID := userDetails.UserAccount

		createPolicyOptions, err := flex.GeneratePolicyOptions(d, meta)
		if err != nil {
			return err
		}

		ibmUniqueID, err := flex.GetIBMUniqueId(accountID, userEmail, meta)
		if err != nil {
			return err
		}

		accountIDResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
			Name:     core.StringPtr("accountId"),
			Value:    core.StringPtr(accountID),
			Operator: core.StringPtr("stringEquals"),
		}

		policyResources := iampolicymanagementv1.PolicyResource{
			Attributes: append(createPolicyOptions.Resources[0].Attributes, *accountIDResourceAttribute),
			Tags:       flex.SetTags(d),
		}

		subjectAttribute := &iampolicymanagementv1.SubjectAttribute{
			Name:  core.StringPtr("iam_id"),
			Value: &ibmUniqueID,
		}
		policySubjects := &iampolicymanagementv1.PolicySubject{
			Attributes: []iampolicymanagementv1.SubjectAttribute{*subjectAttribute},
		}

		getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
			PolicyID: &userPolicyID,
		}
		policy, response, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
		if err != nil || policy == nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error retrieving Policy: %s\n%s", err, response)
		}

		userPolicyETag := response.Headers.Get("ETag")
		updatePolicyOptions := iamPolicyManagementClient.NewUpdatePolicyOptions(
			userPolicyID,
			userPolicyETag,
			"access",
			[]iampolicymanagementv1.PolicySubject{*policySubjects},
			createPolicyOptions.Roles,
			[]iampolicymanagementv1.PolicyResource{policyResources},
		)

		if description, ok := d.GetOk("description"); ok {
			des := description.(string)
			updatePolicyOptions.Description = &des
		}

		_, resp, err := iamPolicyManagementClient.UpdatePolicy(updatePolicyOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating user policy: %s, %s", err, resp)
		}
	}
	return resourceIBMIAMUserPolicyRead(d, meta)
}

func resourceIBMIAMUserPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	userPolicyID := parts[1]

	deletePolicyOptions := iamPolicyManagementClient.NewDeletePolicyOptions(
		userPolicyID,
	)
	_, err = iamPolicyManagementClient.DeletePolicy(deletePolicyOptions)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceIBMIAMUserPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of userEmail/PolicyID", d.Id())
	}
	userEmail := parts[0]
	userPolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		userPolicyID,
	)

	userPolicy, resp, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil || userPolicy == nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting user policy: %s\n%s", err, resp)
	}

	if userPolicy != nil && userPolicy.State != nil && *userPolicy.State == "deleted" {
		return false, nil
	}

	tempID := fmt.Sprintf("%s/%s", userEmail, *userPolicy.ID)

	return tempID == d.Id(), nil
}

func importUserPolicy(d *schema.ResourceData, meta interface{}) (interface{}, interface{}, error) {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return nil, nil, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil, nil, err
	}
	userPolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		userPolicyID,
	)
	userPolicy, _, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] Error retrieving User Policy: %s", err)
	}
	resources := flex.FlattenPolicyResource(userPolicy.Resources)
	resource_attributes := flex.FlattenPolicyResourceAttributes(userPolicy.Resources)
	d.Set("resource_tags", flex.FlattenPolicyResourceTags(userPolicy.Resources))
	return resources, resource_attributes, nil
}
