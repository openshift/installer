// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMServicePolicy() *schema.Resource {
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
					return nil, fmt.Errorf("[ERROR] Error reading resource ID: %s", err)
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
				ValidateFunc: validate.InvokeValidator("ibm_iam_service_policy",
					"iam_service_id"),
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

						"service_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service type of the policy definition",
						},

						"service_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Service group id of the policy definition",
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

			"transaction_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Set transactionID for debug",
			},

			"rule_conditions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Rule conditions enforced by the policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key of the condition",
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Operator of the condition",
						},
						"value": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Value of the condition",
						},
						"conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional Rule conditions enforced by the policy",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Key of the condition",
									},
									"operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Operator of the condition",
									},
									"value": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Value of the condition",
									},
								},
							},
						},
					},
				},
			},

			"rule_operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operator that multiple rule conditions are evaluated over",
			},

			"pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pattern rule follows for time-based condition",
			},
		},
	}
}

func ResourceIBMIAMServicePolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "iam_service_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:service_id", "resolved_to:id"},
			Optional:                   true})

	iBMIAMServicePolicyValidator := validate.ResourceValidator{ResourceName: "ibm_iam_service_policy", Schema: validateSchema}
	return &iBMIAMServicePolicyValidator
}

func resourceIBMIAMServicePolicyCreate(d *schema.ResourceData, meta interface{}) error {

	var iamID string
	if v, ok := d.GetOk("iam_service_id"); ok && v != nil {
		serviceIDUUID := v.(string)

		iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
		if err != nil {
			return err
		}
		getServiceIDOptions := iamidentityv1.GetServiceIDOptions{
			ID: &serviceIDUUID,
		}
		serviceID, resp, err := iamClient.GetServiceID(&getServiceIDOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error] Error Getting Service Id %s %s", err, resp)
		}
		iamID = *serviceID.IamID
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
	var policyID string
	ruleConditions, ruleConditionsDefined := d.GetOk("rule_conditions")
	pattern, patternDefined := d.GetOk("pattern")
	if ruleConditionsDefined || patternDefined {
		policyOptions, err := flex.GenerateV2PolicyOptions(d, meta)

		subjectAttribute := &iampolicymanagementv1.V2PolicySubjectAttribute{
			Key:      core.StringPtr("iam_id"),
			Value:    &iamID,
			Operator: core.StringPtr("stringEquals"),
		}

		policySubject := &iampolicymanagementv1.V2PolicySubject{
			Attributes: []iampolicymanagementv1.V2PolicySubjectAttribute{*subjectAttribute},
		}

		accountIDResourceAttribute := &iampolicymanagementv1.V2PolicyResourceAttribute{
			Key:      core.StringPtr("accountId"),
			Value:    core.StringPtr(userDetails.UserAccount),
			Operator: core.StringPtr("stringEquals"),
		}

		policyResource := &iampolicymanagementv1.V2PolicyResource{
			Attributes: append(policyOptions.Resource.Attributes, *accountIDResourceAttribute),
			Tags:       flex.SetV2PolicyTags(d),
		}

		createPolicyOptions := iamPolicyManagementClient.NewCreateV2PolicyOptions(
			policyOptions.Control,
			"access",
		)

		createPolicyOptions.SetSubject(policySubject)
		createPolicyOptions.SetResource(policyResource)

		if patternDefined {
			createPolicyOptions.SetPattern(pattern.(string))
		}

		if ruleConditionsDefined {
			createPolicyOptions.SetRule(flex.GeneratePolicyRule(d, ruleConditions))
		}

		if description, ok := d.GetOk("description"); ok {
			des := description.(string)
			createPolicyOptions.Description = &des
		}

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			createPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		servicePolicy, resp, err := iamPolicyManagementClient.CreateV2Policy(createPolicyOptions)

		if err != nil {
			return fmt.Errorf("Error creating user policies: %s, %s", err, resp)
		}

		policyID = *servicePolicy.ID
	} else {
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

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			createPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		servicePolicy, res, err := iamPolicyManagementClient.CreatePolicy(createPolicyOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating servicePolicy: %s %s", err, res)
		}

		policyID = *servicePolicy.ID
	}

	getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
		policyID,
	)

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		getPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		policy, res, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)

		if err != nil || policy == nil {
			if res != nil && res.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		_, _, err = iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
	}
	if err != nil {
		if v, ok := d.GetOk("iam_service_id"); ok && v != nil {
			serviceIDUUID := v.(string)
			d.SetId(fmt.Sprintf("%s/%s", serviceIDUUID, policyID))
		} else if v, ok := d.GetOk("iam_id"); ok && v != nil {
			iamID := v.(string)
			d.SetId(fmt.Sprintf("%s/%s", iamID, policyID))
		}
		return fmt.Errorf("[ERROR] Error fetching service  policy: %s", err)
	}
	if v, ok := d.GetOk("iam_service_id"); ok && v != nil {
		serviceIDUUID := v.(string)
		d.SetId(fmt.Sprintf("%s/%s", serviceIDUUID, policyID))
	} else if v, ok := d.GetOk("iam_id"); ok && v != nil {
		iamID := v.(string)
		d.SetId(fmt.Sprintf("%s/%s", iamID, policyID))
	}

	return resourceIBMIAMServicePolicyRead(d, meta)
}

func resourceIBMIAMServicePolicyRead(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	serviceIDUUID := parts[0]
	servicePolicyID := parts[1]
	servicePolicy := &iampolicymanagementv1.V2PolicyTemplateMetaData{}
	res := &core.DetailedResponse{}
	getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
		servicePolicyID,
	)

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		getPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		servicePolicy, res, err = iamPolicyManagementClient.GetV2Policy(getPolicyOptions)

		if err != nil || servicePolicy == nil {
			if res != nil && res.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		servicePolicy, res, err = iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
	}
	if err != nil || servicePolicy == nil || res == nil {
		return fmt.Errorf("[ERROR] Error retrieving servicePolicy: %s %s", err, res)
	}
	if strings.HasPrefix(serviceIDUUID, "iam-") {
		d.Set("iam_id", serviceIDUUID)
	} else {
		d.Set("iam_service_id", serviceIDUUID)
	}

	roles, err := flex.GetRoleNamesFromPolicyResponse(*servicePolicy, d, meta)
	d.Set("roles", roles)

	if _, ok := d.GetOk("resources"); ok {
		d.Set("resources", flex.FlattenV2PolicyResource(*servicePolicy.Resource))
	}
	if _, ok := d.GetOk("resource_attributes"); ok {
		d.Set("resource_attributes", flex.FlattenV2PolicyResourceAttributes(servicePolicy.Resource.Attributes))
	}

	if _, ok := d.GetOk("resource_tags"); ok {
		d.Set("resource_tags", flex.FlattenV2PolicyResourceTags(*servicePolicy.Resource))
	}

	if _, ok := d.GetOk("rule_conditions"); ok {
		d.Set("rule_conditions", flex.FlattenRuleConditions(*servicePolicy.Rule.(*iampolicymanagementv1.V2PolicyRule)))
	}

	if _, ok := d.GetOk("rule_operator"); ok {
		d.Set("rule_operator", *servicePolicy.Rule.(*iampolicymanagementv1.V2PolicyRule).Operator)
	}

	if (&iampolicymanagementv1.V2PolicyResource{}) != servicePolicy.Resource {
		if flex.GetV2PolicyResourceAttribute("serviceType", *servicePolicy.Resource) == "service" {
			d.Set("account_management", false)
		}
		if flex.GetV2PolicyResourceAttribute("serviceType", *servicePolicy.Resource) == "platform_service" {
			d.Set("account_management", true)
		}
	}
	if servicePolicy.Description != nil {
		d.Set("description", *servicePolicy.Description)
	}

	if len(res.Headers["Transaction-Id"]) > 0 && res.Headers["Transaction-Id"][0] != "" {
		d.Set("transaction_id", res.Headers["Transaction-Id"][0])
	}

	return nil
}

func resourceIBMIAMServicePolicyUpdate(d *schema.ResourceData, meta interface{}) error {

	if d.HasChange("roles") || d.HasChange("resources") || d.HasChange("resource_attributes") || d.HasChange("account_management") || d.HasChange("description") || d.HasChange("resource_tags") || d.HasChange("rule_conditions") || d.HasChange("rule_operator") || d.HasChange("pattern") {

		parts, err := flex.IdParts(d.Id())
		if err != nil {
			return err
		}
		servicePolicyID := parts[1]

		var iamID string
		if v, ok := d.GetOk("iam_service_id"); ok && v != nil {
			serviceIDUUID := v.(string)

			iamClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
			if err != nil {
				return err
			}
			getServiceIDOptions := iamidentityv1.GetServiceIDOptions{
				ID: &serviceIDUUID,
			}

			if transactionID, ok := d.GetOk("transaction_id"); ok {
				getServiceIDOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
			}

			serviceID, resp, err := iamClient.GetServiceID(&getServiceIDOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error] Error Getting Service Id %s %s", err, resp)
			}
			iamID = *serviceID.IamID
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

		getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
			servicePolicyID,
		)

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			getPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		policy, response, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
		if err != nil || policy == nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error retrieving Policy: %s\n%s", err, response)
		}

		servicePolicyETag := response.Headers.Get("ETag")

		if strings.Contains(*policy.Href, "/v2/policies") {
			createPolicyOptions, err := flex.GenerateV2PolicyOptions(d, meta)
			if err != nil {
				return err
			}
			subjectAttribute := &iampolicymanagementv1.V2PolicySubjectAttribute{
				Key:      core.StringPtr("iam_id"),
				Value:    &iamID,
				Operator: core.StringPtr("stringEquals"),
			}

			policySubject := &iampolicymanagementv1.V2PolicySubject{
				Attributes: []iampolicymanagementv1.V2PolicySubjectAttribute{*subjectAttribute},
			}

			accountIDResourceAttribute := &iampolicymanagementv1.V2PolicyResourceAttribute{
				Key:      core.StringPtr("accountId"),
				Value:    core.StringPtr(userDetails.UserAccount),
				Operator: core.StringPtr("stringEquals"),
			}

			policyResource := &iampolicymanagementv1.V2PolicyResource{
				Attributes: append(createPolicyOptions.Resource.Attributes, *accountIDResourceAttribute),
				Tags:       flex.SetV2PolicyTags(d),
			}

			updatePolicyOptions := iamPolicyManagementClient.NewReplaceV2PolicyOptions(
				servicePolicyID,
				servicePolicyETag,
				createPolicyOptions.Control,
				"access",
			)
			updatePolicyOptions.SetSubject(policySubject)
			updatePolicyOptions.SetResource(policyResource)

			if pattern, ok := d.GetOk("pattern"); ok {
				updatePolicyOptions.SetPattern(pattern.(string))
			}

			if ruleConditions, ok := d.GetOk("rule_conditions"); ok {
				updatePolicyOptions.SetRule(flex.GeneratePolicyRule(d, ruleConditions))
			}

			if description, ok := d.GetOk("description"); ok {
				des := description.(string)
				updatePolicyOptions.Description = &des
			}

			if transactionID, ok := d.GetOk("transaction_id"); ok {
				updatePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
			}

			_, resp, err := iamPolicyManagementClient.ReplaceV2Policy(updatePolicyOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating user policy: %s, %s", err, resp)
			}
		} else {
			_, ruleConditionsDefined := d.GetOk("rule_conditions")
			_, patternDefined := d.GetOk("pattern")
			if ruleConditionsDefined || patternDefined {
				return fmt.Errorf("Cannot use rule_conditions, rule_operator, or pattern when updating v1/policy. Delete existing v1/policy and create using rule_conditions and pattern.")
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

			updatePolicyOptions := iamPolicyManagementClient.NewReplacePolicyOptions(
				servicePolicyID,
				servicePolicyETag,
				"access",
				[]iampolicymanagementv1.PolicySubject{*policySubjects},
				createPolicyOptions.Roles,
				[]iampolicymanagementv1.PolicyResource{policyResources},
			)

			if desc, ok := d.GetOk("description"); ok {
				des := desc.(string)
				updatePolicyOptions.Description = &des
			}

			if transactionID, ok := d.GetOk("transaction_id"); ok {
				updatePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
			}
			_, _, err = iamPolicyManagementClient.ReplacePolicy(updatePolicyOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating service policy: %s", err)
			}
		}
	}

	return resourceIBMIAMServicePolicyRead(d, meta)

}

func resourceIBMIAMServicePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	servicePolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
		servicePolicyID,
	)

	// Get policy to find version in href
	policy, _, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)

	if err != nil {
		return err
	}

	if strings.Contains(*policy.Href, "/v2/policies") {
		deletePolicyOptions := iamPolicyManagementClient.NewDeleteV2PolicyOptions(
			servicePolicyID,
		)
		if transactionID, ok := d.GetOk("transaction_id"); ok {
			deletePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		_, err = iamPolicyManagementClient.DeleteV2Policy(deletePolicyOptions)
	} else {
		deletePolicyOptions := iamPolicyManagementClient.NewDeletePolicyOptions(
			servicePolicyID,
		)

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			deletePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		_, err = iamPolicyManagementClient.DeletePolicy(deletePolicyOptions)
	}

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting service policy: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMServicePolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of serviceID(OR)iamID/PolicyID", d.Id())
	}
	serviceIDUUID := parts[0]
	servicePolicyID := parts[1]

	getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
		servicePolicyID,
	)

	servicePolicy, resp, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
	if err != nil || servicePolicy == nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting service policy: %s\n%s", err, resp)
	}

	if servicePolicy != nil && servicePolicy.State != nil && *servicePolicy.State == "deleted" {
		return false, nil
	}

	tempID := fmt.Sprintf("%s/%s", serviceIDUUID, *servicePolicy.ID)

	return tempID == d.Id(), nil
}

func importServicePolicy(d *schema.ResourceData, meta interface{}) (interface{}, interface{}, error) {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return nil, nil, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return nil, nil, err
	}
	servicePolicyID := parts[1]
	getPolicyOptions := iamPolicyManagementClient.NewGetV2PolicyOptions(
		servicePolicyID,
	)
	servicePolicy, _, err := iamPolicyManagementClient.GetV2Policy(getPolicyOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("[ERROR] Error retrieving servicePolicy: %s", err)
	}
	resources := flex.FlattenV2PolicyResource(*servicePolicy.Resource)
	resource_attributes := flex.FlattenV2PolicyResourceAttributes(servicePolicy.Resource.Attributes)
	d.Set("resource_tags", flex.FlattenV2PolicyResourceTags(*servicePolicy.Resource))
	return resources, resource_attributes, nil
}
