// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMTrustedProfilePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMIAMTrustedProfilePolicyCreate,
		Read:   resourceIBMIAMTrustedProfilePolicyRead,
		Update: resourceIBMIAMTrustedProfilePolicyUpdate,
		Delete: resourceIBMIAMTrustedProfilePolicyDelete,
		Exists: resourceIBMIAMTrustedProfilePolicyExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				resources, resourceAttributes, err := importTrustedProfilePolicy(d, meta)
				if err != nil {
					return nil, fmt.Errorf("[ERROR] Error reading resource ID: %s", err)
				}
				d.Set("resources", resources)
				d.Set("resource_attributes", resourceAttributes)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"profile_id", "iam_id"},
				Description:  "UUID of Trusted Profile",
				ForceNew:     true,
			},
			"iam_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"profile_id", "iam_id"},
				Description:  "IAM ID of Trusted Profile",
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

func resourceIBMIAMTrustedProfilePolicyCreate(d *schema.ResourceData, meta interface{}) error {

	var iamID string
	if v, ok := d.GetOk("profile_id"); ok && v != nil {
		profileIDUUID := v.(string)

		iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
		if err != nil {
			return err
		}
		getProfileOptions := &iamidentityv1.GetProfileOptions{
			ProfileID: &profileIDUUID,
		}
		profileID, resp, err := iamClient.GetProfile(getProfileOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error] Error getting trusted profile ID %s %s", err, resp)
		}
		iamID = *profileID.IamID
	}
	if v, ok := d.GetOk("iam_id"); ok && v != nil {
		iamID = v.(string)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	policyOptions, err := flex.GeneratePolicyOptions(d, meta)
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
		Value:    core.StringPtr(userDetails.UserAccount),
		Operator: core.StringPtr("stringEquals"),
	}

	policyResources := iampolicymanagementv1.PolicyResource{
		Attributes: append(policyOptions.Resources[0].Attributes, *accountIDResourceAttribute),
		Tags:       flex.SetTags(d),
	}

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	createPolicyOptions := iamPolicyManagementClient.NewCreatePolicyOptions(
		"access",
		[]iampolicymanagementv1.PolicySubject{*policySubjects},
		policyOptions.Roles,
		[]iampolicymanagementv1.PolicyResource{policyResources},
	)

	if desc, ok := d.GetOk("description"); ok {
		des := desc.(string)
		createPolicyOptions.Description = &des
	}

	trustedProfilePolicy, res, err := iamPolicyManagementClient.CreatePolicy(createPolicyOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating trustedProfilePolicy: %s %s", err, res)
	}

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		*trustedProfilePolicy.ID,
	)

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
		if v, ok := d.GetOk("profile_id"); ok && v != nil {
			profileIDUUID := v.(string)
			d.SetId(fmt.Sprintf("%s/%s", profileIDUUID, *trustedProfilePolicy.ID))
		} else if v, ok := d.GetOk("iam_id"); ok && v != nil {
			iamID := v.(string)
			d.SetId(fmt.Sprintf("%s/%s", iamID, *trustedProfilePolicy.ID))
		}
		return fmt.Errorf("[ERROR] Error fetching trusted profile policy: %s %s", err, res)
	}
	if v, ok := d.GetOk("profile_id"); ok && v != nil {
		profileIDUUID := v.(string)
		d.SetId(fmt.Sprintf("%s/%s", profileIDUUID, *trustedProfilePolicy.ID))
	} else if v, ok := d.GetOk("iam_id"); ok && v != nil {
		iamID := v.(string)
		d.SetId(fmt.Sprintf("%s/%s", iamID, *trustedProfilePolicy.ID))
	}

	return resourceIBMIAMTrustedProfilePolicyRead(d, meta)
}

func resourceIBMIAMTrustedProfilePolicyRead(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	profileIDUUID := parts[0]
	trustedProfilePolicyID := parts[1]
	trustedProfilePolicy := &iampolicymanagementv1.Policy{}
	res := &core.DetailedResponse{}
	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		trustedProfilePolicyID,
	)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		trustedProfilePolicy, res, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)

		if err != nil || trustedProfilePolicy == nil {
			if res != nil && res.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		trustedProfilePolicy, res, err = iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	}
	if err != nil || trustedProfilePolicy == nil {
		return fmt.Errorf("[ERROR] Error retrieving trusted profile policy: %s %s", err, res)
	}
	if strings.HasPrefix(profileIDUUID, "iam-") {
		d.Set("iam_id", profileIDUUID)
	} else {
		d.Set("profile_id", profileIDUUID)
	}

	roles := make([]string, len(trustedProfilePolicy.Roles))
	for i, role := range trustedProfilePolicy.Roles {
		roles[i] = *role.DisplayName
	}
	d.Set("roles", roles)

	if _, ok := d.GetOk("resources"); ok {
		d.Set("resources", flex.FlattenPolicyResource(trustedProfilePolicy.Resources))
	}
	if _, ok := d.GetOk("resource_attributes"); ok {
		d.Set("resource_attributes", flex.FlattenPolicyResourceAttributes(trustedProfilePolicy.Resources))
	}

	if _, ok := d.GetOk("resource_tags"); ok {
		d.Set("resource_tags", flex.FlattenPolicyResourceTags(trustedProfilePolicy.Resources))
	}

	if len(trustedProfilePolicy.Resources) > 0 {
		if *flex.GetResourceAttribute("serviceType", trustedProfilePolicy.Resources[0]) == "service" {
			d.Set("account_management", false)
		}
		if *flex.GetResourceAttribute("serviceType", trustedProfilePolicy.Resources[0]) == "platform_service" {
			d.Set("account_management", true)
		}
	}
	if trustedProfilePolicy.Description != nil {
		d.Set("description", *trustedProfilePolicy.Description)
	}

	return nil
}

func resourceIBMIAMTrustedProfilePolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("roles") || d.HasChange("resources") || d.HasChange("resource_attributes") || d.HasChange("account_management") || d.HasChange("description") || d.HasChange("resource_tags") {

		parts, err := flex.IdParts(d.Id())
		if err != nil {
			return err
		}
		trustedProfilePolicyID := parts[1]

		var iamID string
		if v, ok := d.GetOk("profile_id"); ok && v != nil {
			profileIDUUID := v.(string)

			iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
			if err != nil {
				return err
			}
			getProfileIDOptions := iamidentityv1.GetProfileOptions{
				ProfileID: &profileIDUUID,
			}
			profileID, resp, err := iamClient.GetProfile(&getProfileIDOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error] Error getting trusted profile ID %s %s", err, resp)
			}
			iamID = *profileID.IamID
		}
		if v, ok := d.GetOk("iam_id"); ok && v != nil {
			iamID = v.(string)
		}

		userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}

		createPolicyOptions, err := flex.GeneratePolicyOptions(d, meta)
		if err != nil {
			return err
		}

		accountIDResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
			Name:     core.StringPtr("accountId"),
			Value:    core.StringPtr(userDetails.UserAccount),
			Operator: core.StringPtr("stringEquals"),
		}

		policyResources := iampolicymanagementv1.PolicyResource{
			Attributes: append(createPolicyOptions.Resources[0].Attributes, *accountIDResourceAttribute),
			Tags:       flex.SetTags(d),
		}

		subjectAttribute := &iampolicymanagementv1.SubjectAttribute{
			Name:  core.StringPtr("iam_id"),
			Value: &iamID,
		}
		policySubjects := &iampolicymanagementv1.PolicySubject{
			Attributes: []iampolicymanagementv1.SubjectAttribute{*subjectAttribute},
		}

		iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
		if err != nil {
			return err
		}

		getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
			trustedProfilePolicyID,
		)
		policy, response, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
		if err != nil || policy == nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error retrieving Policy: %s\n%s", err, response)
		}

		trustedProfilePolicyETag := response.Headers.Get("ETag")
		updatePolicyOptions := iamPolicyManagementClient.NewUpdatePolicyOptions(
			trustedProfilePolicyID,
			trustedProfilePolicyETag,
			"access",
			[]iampolicymanagementv1.PolicySubject{*policySubjects},
			createPolicyOptions.Roles,
			[]iampolicymanagementv1.PolicyResource{policyResources},
		)

		if desc, ok := d.GetOk("description"); ok {
			des := desc.(string)
			updatePolicyOptions.Description = &des
		}

		_, resp, err := iamPolicyManagementClient.UpdatePolicy(updatePolicyOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating trusted profile policy: %s: %s", err, resp)
		}

	}

	return resourceIBMIAMTrustedProfilePolicyRead(d, meta)

}

func resourceIBMIAMTrustedProfilePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	trustedProfilePolicyID := parts[1]

	deletePolicyOptions := iamPolicyManagementClient.NewDeletePolicyOptions(
		trustedProfilePolicyID,
	)

	resp, err := iamPolicyManagementClient.DeletePolicy(deletePolicyOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting trusted profile policy: %s %s", err, resp)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMTrustedProfilePolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of profileID(OR)iamID/PolicyID", d.Id())
	}
	profileIDUUID := parts[0]
	trustedProfilePolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		trustedProfilePolicyID,
	)

	trustedProfilePolicy, resp, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil || trustedProfilePolicy == nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting trusted profile policy: %s\n%s", err, resp)
	}

	if trustedProfilePolicy != nil && trustedProfilePolicy.State != nil && *trustedProfilePolicy.State == "deleted" {
		return false, nil
	}

	tempID := fmt.Sprintf("%s/%s", profileIDUUID, *trustedProfilePolicy.ID)

	return tempID == d.Id(), nil
}

func importTrustedProfilePolicy(d *schema.ResourceData, meta interface{}) (interface{}, interface{}, error) {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return nil, nil, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil, nil, err
	}
	trustedProfilePolicyID := parts[1]
	getPolicyOptions := iamPolicyManagementClient.NewGetPolicyOptions(
		trustedProfilePolicyID,
	)
	trustedProfilePolicy, resp, err := iamPolicyManagementClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] Error retrieving trusted profile policy: %s %s", err, resp)
	}
	resources := flex.FlattenPolicyResource(trustedProfilePolicy.Resources)
	resource_attributes := flex.FlattenPolicyResourceAttributes(trustedProfilePolicy.Resources)
	d.Set("resource_tags", flex.FlattenPolicyResourceTags(trustedProfilePolicy.Resources))
	return resources, resource_attributes, nil
}
