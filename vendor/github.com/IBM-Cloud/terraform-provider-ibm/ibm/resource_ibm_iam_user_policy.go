// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
)

func resourceIBMIAMUserPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMIAMUserPolicyCreate,
		Read:   resourceIBMIAMUserPolicyRead,
		Update: resourceIBMIAMUserPolicyUpdate,
		Delete: resourceIBMIAMUserPolicyDelete,
		Exists: resourceIBMIAMUserPolicyExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				resources, resourceAttributes, err := importServicePolicy(d, meta)
				if err != nil {
					return nil, fmt.Errorf("Error reading resource ID: %s", err)
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
		},
	}
}

func resourceIBMIAMUserPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	var policyOptions iampolicymanagementv1.CreatePolicyOptions
	policyOptions, err = generatePolicyOptions(d, meta)

	if err != nil {
		return err
	}

	userEmail := d.Get("ibm_id").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	accountID := userDetails.userAccount

	ibmUniqueID, err := getIBMUniqueId(accountID, userEmail, meta)
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
	}

	createPolicyOptions := iamPolicyManagementClient.NewCreatePolicyOptions(
		"access",
		[]iampolicymanagementv1.PolicySubject{*policySubjects},
		policyOptions.Roles,
		[]iampolicymanagementv1.PolicyResource{policyResources},
	)

	userPolicy, _, err := iamPolicyManagementClient.CreatePolicy(createPolicyOptions)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", userEmail, *userPolicy.ID))

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: userPolicy.ID,
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, _, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)

		if err != nil {
			if apiErr, ok := err.(bmxerror.RequestFailure); ok {
				if apiErr.StatusCode() == 404 {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})

	if isResourceTimeoutError(err) {
		_, _, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	}
	if err != nil {
		return fmt.Errorf("error fetching user policy: %w", err)
	}

	return resourceIBMIAMUserPolicyRead(d, meta)
}

func resourceIBMIAMUserPolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
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

	userPolicy, _, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return err
	}
	d.Set("ibm_id", userEmail)
	roles := make([]string, len(userPolicy.Roles))

	for i, role := range userPolicy.Roles {
		roles[i] = *role.DisplayName
	}
	d.Set("roles", roles)
	if _, ok := d.GetOk("resources"); ok {
		d.Set("resources", flattenPolicyResource(userPolicy.Resources))
	}
	if _, ok := d.GetOk("resource_attributes"); ok {
		d.Set("resource_attributes", flattenPolicyResourceAttributes(userPolicy.Resources))
	}
	if len(userPolicy.Resources) > 0 {
		if *getResourceAttribute("serviceType", userPolicy.Resources[0]) == "service" {
			d.Set("account_management", false)
		}
		if *getResourceAttribute("serviceType", userPolicy.Resources[0]) == "platform_service" {
			d.Set("account_management", true)
		}
	}
	return nil
}

func resourceIBMIAMUserPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}
	if d.HasChange("roles") || d.HasChange("resources") || d.HasChange("resource_attributes") || d.HasChange("account_management") {
		parts, err := idParts(d.Id())
		if err != nil {
			return err
		}
		userEmail := parts[0]
		userPolicyID := parts[1]

		userDetails, err := meta.(ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}

		accountID := userDetails.userAccount

		createPolicyOptions, err := generatePolicyOptions(d, meta)
		if err != nil {
			return err
		}

		ibmUniqueID, err := getIBMUniqueId(accountID, userEmail, meta)
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
			return fmt.Errorf("Error retrieving Policy: %s\n%s", err, response)
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

		policy, _, err = iamPolicyManagementClient.UpdatePolicy(updatePolicyOptions)
		if err != nil {
			return fmt.Errorf("Error updating user policy: %s", err)
		}
	}
	return resourceIBMIAMUserPolicyRead(d, meta)
}

func resourceIBMIAMUserPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
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
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	userEmail := parts[0]
	userPolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		userPolicyID,
	)

	userPolicy, _, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	tempID := fmt.Sprintf("%s/%s", userEmail, *userPolicy.ID)

	return tempID == d.Id(), nil

}

func importUserPolicy(d *schema.ResourceData, meta interface{}) (interface{}, interface{}, error) {

	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return nil, nil, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return nil, nil, err
	}
	userPolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		userPolicyID,
	)
	userPolicy, _, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("Error retrieving User Policy: %s", err)
	}
	resources := flattenPolicyResource(userPolicy.Resources)
	resource_attributes := flattenPolicyResourceAttributes(userPolicy.Resources)
	return resources, resource_attributes, nil
}
