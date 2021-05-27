// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
)

func resourceIBMIAMServicePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMIAMServicePolicyCreate,
		Read:   resourceIBMIAMServicePolicyRead,
		Update: resourceIBMIAMServicePolicyUpdate,
		Delete: resourceIBMIAMServicePolicyDelete,
		Exists: resourceIBMIAMServicePolicyExists,
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
			"iam_service_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"iam_service_id", "iam_id"},
				Description:  "UUID of ServiceID",
				ForceNew:     true,
			},
			"iam_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"iam_service_id", "iam_id"},
				Description:  "IAM ID of ServiceID",
				ForceNew:     true,
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

func resourceIBMIAMServicePolicyCreate(d *schema.ResourceData, meta interface{}) error {

	var iamID string
	if v, ok := d.GetOk("iam_service_id"); ok && v != nil {
		serviceIDUUID := v.(string)

		iamClient, err := meta.(ClientSession).IAMAPI()
		if err != nil {
			return err
		}
		serviceID, err := iamClient.ServiceIds().Get(serviceIDUUID)
		if err != nil {
			return err
		}
		iamID = serviceID.IAMID
	}
	if v, ok := d.GetOk("iam_id"); ok && v != nil {
		iamID = v.(string)
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	policyOptions, err := generatePolicyOptions(d, meta)
	if err != nil {
		return err
	}

	subjectAttribute := &iampolicymanagementv1.SubjectAttribute{
		Name:  core.StringPtr("iam_id"),
		Value: &iamID,
	}

	policySubjects := &iampolicymanagementv1.PolicySubject{
		Attributes: []iampolicymanagementv1.SubjectAttribute{*subjectAttribute},
	}

	accountIDResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
		Name:     core.StringPtr("accountId"),
		Value:    core.StringPtr(userDetails.userAccount),
		Operator: core.StringPtr("stringEquals"),
	}

	policyResources := iampolicymanagementv1.PolicyResource{
		Attributes: append(policyOptions.Resources[0].Attributes, *accountIDResourceAttribute),
	}

	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	createPolicyOptions := iamPolicyManagementClient.NewCreatePolicyOptions(
		"access",
		[]iampolicymanagementv1.PolicySubject{*policySubjects},
		policyOptions.Roles,
		[]iampolicymanagementv1.PolicyResource{policyResources},
	)

	servicePolicy, _, err := iamPolicyManagementClient.CreatePolicy(createPolicyOptions)

	if err != nil {
		return fmt.Errorf("Error creating servicePolicy: %s", err)
	}
	if v, ok := d.GetOk("iam_service_id"); ok && v != nil {
		serviceIDUUID := v.(string)
		d.SetId(fmt.Sprintf("%s/%s", serviceIDUUID, *servicePolicy.ID))
	} else if v, ok := d.GetOk("iam_id"); ok && v != nil {
		iamID := v.(string)
		d.SetId(fmt.Sprintf("%s/%s", iamID, *servicePolicy.ID))
	}

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		*servicePolicy.ID,
	)

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
		return fmt.Errorf("error fetching service  policy: %w", err)
	}

	return resourceIBMIAMServicePolicyRead(d, meta)
}

func resourceIBMIAMServicePolicyRead(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	serviceIDUUID := parts[0]
	servicePolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		servicePolicyID,
	)
	servicePolicy, _, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return fmt.Errorf("Error retrieving servicePolicy: %s", err)
	}
	if strings.HasPrefix(serviceIDUUID, "iam-") {
		d.Set("iam_id", serviceIDUUID)
	} else {
		d.Set("iam_service_id", serviceIDUUID)
	}

	roles := make([]string, len(servicePolicy.Roles))
	for i, role := range servicePolicy.Roles {
		roles[i] = *role.DisplayName
	}
	d.Set("roles", roles)

	if _, ok := d.GetOk("resources"); ok {
		d.Set("resources", flattenPolicyResource(servicePolicy.Resources))
	}
	if _, ok := d.GetOk("resource_attributes"); ok {
		d.Set("resource_attributes", flattenPolicyResourceAttributes(servicePolicy.Resources))
	}
	if len(servicePolicy.Resources) > 0 {
		if *getResourceAttribute("serviceType", servicePolicy.Resources[0]) == "service" {
			d.Set("account_management", false)
		}
		if *getResourceAttribute("serviceType", servicePolicy.Resources[0]) == "platform_service" {
			d.Set("account_management", true)
		}
	}

	return nil
}

func resourceIBMIAMServicePolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("roles") || d.HasChange("resources") || d.HasChange("resource_attributes") || d.HasChange("account_management") {

		parts, err := idParts(d.Id())
		if err != nil {
			return err
		}
		servicePolicyID := parts[1]

		var iamID string
		if v, ok := d.GetOk("iam_service_id"); ok && v != nil {
			serviceIDUUID := v.(string)

			iamClient, err := meta.(ClientSession).IAMAPI()
			if err != nil {
				return err
			}
			serviceID, err := iamClient.ServiceIds().Get(serviceIDUUID)
			if err != nil {
				return err
			}
			iamID = serviceID.IAMID
		}
		if v, ok := d.GetOk("iam_id"); ok && v != nil {
			iamID = v.(string)
		}

		userDetails, err := meta.(ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}

		createPolicyOptions, err := generatePolicyOptions(d, meta)
		if err != nil {
			return err
		}

		accountIDResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
			Name:     core.StringPtr("accountId"),
			Value:    core.StringPtr(userDetails.userAccount),
			Operator: core.StringPtr("stringEquals"),
		}

		policyResources := iampolicymanagementv1.PolicyResource{
			Attributes: append(createPolicyOptions.Resources[0].Attributes, *accountIDResourceAttribute),
		}

		subjectAttribute := &iampolicymanagementv1.SubjectAttribute{
			Name:  core.StringPtr("iam_id"),
			Value: &iamID,
		}
		policySubjects := &iampolicymanagementv1.PolicySubject{
			Attributes: []iampolicymanagementv1.SubjectAttribute{*subjectAttribute},
		}

		iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
		if err != nil {
			return err
		}

		getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
			servicePolicyID,
		)
		policy, response, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
		if err != nil || policy == nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("Error retrieving Policy: %s\n%s", err, response)
		}

		servicePolicyETag := response.Headers.Get("ETag")
		updatePolicyOptions := iamPolicyManagementClient.NewUpdatePolicyOptions(
			servicePolicyID,
			servicePolicyETag,
			"access",
			[]iampolicymanagementv1.PolicySubject{*policySubjects},
			createPolicyOptions.Roles,
			[]iampolicymanagementv1.PolicyResource{policyResources},
		)

		_, _, err = iamPolicyManagementClient.UpdatePolicy(updatePolicyOptions)
		if err != nil {
			return fmt.Errorf("Error updating service policy: %s", err)
		}

	}

	return resourceIBMIAMServicePolicyRead(d, meta)

}

func resourceIBMIAMServicePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	servicePolicyID := parts[1]

	deletePolicyOptions := iamPolicyManagementClient.NewDeletePolicyOptions(
		servicePolicyID,
	)

	_, err = iamPolicyManagementClient.DeletePolicy(deletePolicyOptions)
	if err != nil {
		return fmt.Errorf("Error deleting service policy: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMServicePolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	serviceIDUUID := parts[0]
	servicePolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		servicePolicyID,
	)

	servicePolicy, _, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	tempID := fmt.Sprintf("%s/%s", serviceIDUUID, *servicePolicy.ID)

	return tempID == d.Id(), nil
}

func importServicePolicy(d *schema.ResourceData, meta interface{}) (interface{}, interface{}, error) {

	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return nil, nil, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return nil, nil, err
	}
	servicePolicyID := parts[1]
	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		servicePolicyID,
	)
	servicePolicy, _, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("Error retrieving servicePolicy: %s", err)
	}
	resources := flattenPolicyResource(servicePolicy.Resources)
	resource_attributes := flattenPolicyResourceAttributes(servicePolicy.Resources)
	return resources, resource_attributes, nil
}
