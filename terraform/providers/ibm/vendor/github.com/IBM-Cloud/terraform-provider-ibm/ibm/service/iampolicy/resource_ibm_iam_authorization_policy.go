// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIAMAuthorizationPolicy() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMAuthorizationPolicyCreate,
		Read:     resourceIBMIAMAuthorizationPolicyRead,
		Update:   resourceIBMIAMAuthorizationPolicyUpdate,
		Delete:   resourceIBMIAMAuthorizationPolicyDelete,
		Exists:   resourceIBMIAMAuthorizationPolicyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"source_service_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"source_service_name", "source_resource_group_id", "subject_attributes"},
				Description:  "The source service name",
			},

			"target_service_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"target_service_name", "target_resource_type", "resource_attributes"},
				Description:  "The target service name",
			},

			"roles": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Role names of the policy definition",
			},

			"source_resource_instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"subject_attributes"},
				Description:   "The source resource instance Id",
			},

			"target_resource_instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"resource_attributes"},
				Description:   "The target resource instance Id",
			},

			"source_resource_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"subject_attributes"},
				Description:   "The source resource group Id",
				ValidateFunc: validate.InvokeValidator("ibm_iam_authorization_policy",
					"source_resource_group_id"),
			},

			"target_resource_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"resource_attributes"},
				Description:   "The target resource group Id",
				ValidateFunc: validate.InvokeValidator("ibm_iam_authorization_policy",
					"target_resource_group_id"),
			},

			"source_resource_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"subject_attributes"},
				Description:   "Resource type of source service",
			},

			"target_resource_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"resource_attributes"},
				Description:   "Resource type of target service",
			},

			"source_service_account": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"subject_attributes"},
				Description:   "Account GUID of source service",
			},

			"subject_attributes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Description:   "Set subject attributes.",
				ConflictsWith: []string{"source_service_name", "source_resource_instance_id", "source_resource_group_id", "source_resource_type", "source_service_account"},
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
							Description: "Operator of attribute.",
						},
					},
				},
			},

			"resource_attributes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Description:   "Set resource attributes.",
				ConflictsWith: []string{"target_service_name", "target_resource_instance_id", "target_resource_group_id", "target_resource_type"},
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

			"version": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

func ResourceIBMIAMAuthorizationPolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "source_resource_group_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_group",
			CloudDataRange:             []string{"resolved_to:id"},
			Optional:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "target_resource_group_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_group",
			CloudDataRange:             []string{"resolved_to:id"},
			Optional:                   true})

	iBMIAMAuthorizationPolicyValidator := validate.ResourceValidator{ResourceName: "ibm_iam_authorization_policy", Schema: validateSchema}
	return &iBMIAMAuthorizationPolicyValidator
}
func resourceIBMIAMAuthorizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	var sourceServiceName, targetServiceName string
	policyType := "authorization"
	policySubject := &iampolicymanagementv1.V2PolicySubject{}
	policyResource := &iampolicymanagementv1.V2PolicyResource{}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	iampapClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	// check subject_attributes exists
	if attributes, ok := d.GetOk("subject_attributes"); ok {
		for _, attribute := range attributes.(*schema.Set).List() {
			a := attribute.(map[string]interface{})
			name := a["name"].(string)
			value := a["value"].(string)
			operator := a["operator"].(string)
			if name == "serviceName" {
				sourceServiceName = value
			}

			if operator == "" && value == "*" && name == "resourceGroupId" {
				at := iampolicymanagementv1.V2PolicySubjectAttribute{
					Key:      &name,
					Value:    true,
					Operator: core.StringPtr("stringExists"),
				}
				policySubject.Attributes = append(policySubject.Attributes, at)
			} else if operator == "stringExists" {
				var resourceValue bool
				if value == "true" {
					resourceValue = true
				} else if value == "false" {
					resourceValue = false
				} else {
					return fmt.Errorf("[ERROR] Only values \"true\" and \"false\" are allowed when operator is \"stringExists\". Received %s.", value)
				}
				at := iampolicymanagementv1.V2PolicySubjectAttribute{
					Key:      &name,
					Value:    &resourceValue,
					Operator: &operator,
				}
				policySubject.Attributes = append(policySubject.Attributes, at)
			} else {
				if operator == "" {
					operator = "stringEquals"
				}
				at := iampolicymanagementv1.V2PolicySubjectAttribute{
					Key:      &name,
					Value:    &value,
					Operator: &operator,
				}
				policySubject.Attributes = append(policySubject.Attributes, at)
			}
		}
	} else {

		if name, ok := d.GetOk("source_service_name"); ok {
			sourceServiceName = name.(string)

			serviceNameSubjectAttribute := &iampolicymanagementv1.V2PolicySubjectAttribute{
				Key:      core.StringPtr("serviceName"),
				Value:    &sourceServiceName,
				Operator: core.StringPtr("stringEquals"),
			}
			policySubject.Attributes = append(policySubject.Attributes, *serviceNameSubjectAttribute)
		}

		sourceServiceAccount := userDetails.UserAccount
		if account, ok := d.GetOk("source_service_account"); ok {
			sourceServiceAccount = account.(string)
		}

		accountIdSubjectAttribute := &iampolicymanagementv1.V2PolicySubjectAttribute{
			Key:      core.StringPtr("accountId"),
			Value:    &sourceServiceAccount,
			Operator: core.StringPtr("stringEquals"),
		}

		policySubject.Attributes = append(policySubject.Attributes, *accountIdSubjectAttribute)

		if sID, ok := d.GetOk("source_resource_instance_id"); ok {
			serviceInstanceSubjectAttribute := iampolicymanagementv1.V2PolicySubjectAttribute{
				Key:      core.StringPtr("serviceInstance"),
				Value:    core.StringPtr(sID.(string)),
				Operator: core.StringPtr("stringEquals"),
			}
			policySubject.Attributes = append(policySubject.Attributes, serviceInstanceSubjectAttribute)
		}

		if sType, ok := d.GetOk("source_resource_type"); ok {
			resourceTypeSubjectAttribute := iampolicymanagementv1.V2PolicySubjectAttribute{
				Key:      core.StringPtr("resourceType"),
				Value:    core.StringPtr(sType.(string)),
				Operator: core.StringPtr("stringEquals"),
			}
			policySubject.Attributes = append(policySubject.Attributes, resourceTypeSubjectAttribute)
		}

		if sResGrpID, ok := d.GetOk("source_resource_group_id"); ok {
			if sResGrpID == "*" {
				resourceGroupSubjectAttribute := iampolicymanagementv1.V2PolicySubjectAttribute{
					Key:      core.StringPtr("resourceGroupId"),
					Value:    true,
					Operator: core.StringPtr("stringExists"),
				}
				policySubject.Attributes = append(policySubject.Attributes, resourceGroupSubjectAttribute)
			} else {
				resourceGroupSubjectAttribute := iampolicymanagementv1.V2PolicySubjectAttribute{
					Key:      core.StringPtr("resourceGroupId"),
					Value:    core.StringPtr(sResGrpID.(string)),
					Operator: core.StringPtr("stringEquals"),
				}
				policySubject.Attributes = append(policySubject.Attributes, resourceGroupSubjectAttribute)
			}
		}
	}

	// check resource_attributes exists
	if attributes, ok := d.GetOk("resource_attributes"); ok {
		for _, attribute := range attributes.(*schema.Set).List() {
			a := attribute.(map[string]interface{})
			name := a["name"].(string)
			value := a["value"].(string)
			operator := a["operator"].(string)
			if name == "serviceName" {
				targetServiceName = value
			}
			if name == "resourceType" && targetServiceName == "" {
				targetServiceName = "resource-controller"
			}
			if operator == "stringExists" {
				var resourceValue bool
				if value == "true" {
					resourceValue = true
				} else if value == "false" {
					resourceValue = false
				} else {
					return fmt.Errorf("[ERROR] When operator equals stringExists, value should be either \"true\" or \"false\", instead of %s", value)
				}
				at := iampolicymanagementv1.V2PolicyResourceAttribute{
					Key:      &name,
					Value:    &resourceValue,
					Operator: &operator,
				}
				policyResource.Attributes = append(policyResource.Attributes, at)
			} else {
				at := iampolicymanagementv1.V2PolicyResourceAttribute{
					Key:      &name,
					Value:    &value,
					Operator: &operator,
				}
				policyResource.Attributes = append(policyResource.Attributes, at)
			}
		}
	} else {
		if name, ok := d.GetOk("target_service_name"); ok {
			targetServiceName = name.(string)
			serviceNameResourceAttribute := &iampolicymanagementv1.V2PolicyResourceAttribute{
				Key:      core.StringPtr("serviceName"),
				Value:    core.StringPtr(targetServiceName),
				Operator: core.StringPtr("stringEquals"),
			}
			policyResource.Attributes = append(policyResource.Attributes, *serviceNameResourceAttribute)
		}

		accountIDResourceAttribute := &iampolicymanagementv1.V2PolicyResourceAttribute{
			Key:      core.StringPtr("accountId"),
			Value:    core.StringPtr(userDetails.UserAccount),
			Operator: core.StringPtr("stringEquals"),
		}

		policyResource.Attributes = append(policyResource.Attributes, *accountIDResourceAttribute)

		if tID, ok := d.GetOk("target_resource_instance_id"); ok {
			serviceInstanceResourceAttribute := iampolicymanagementv1.V2PolicyResourceAttribute{
				Key:      core.StringPtr("serviceInstance"),
				Value:    core.StringPtr(tID.(string)),
				Operator: core.StringPtr("stringEquals"),
			}
			policyResource.Attributes = append(policyResource.Attributes, serviceInstanceResourceAttribute)
		}

		if tType, ok := d.GetOk("target_resource_type"); ok {
			resourceTypeResourceAttribute := iampolicymanagementv1.V2PolicyResourceAttribute{
				Key:      core.StringPtr("resourceType"),
				Value:    core.StringPtr(tType.(string)),
				Operator: core.StringPtr("stringEquals"),
			}
			policyResource.Attributes = append(policyResource.Attributes, resourceTypeResourceAttribute)
			if targetServiceName == "" {
				targetServiceName = "resource-controller"
			}
		}

		if tResGrpID, ok := d.GetOk("target_resource_group_id"); ok {
			resourceGroupResourceAttribute := iampolicymanagementv1.V2PolicyResourceAttribute{
				Key:      core.StringPtr("resourceGroupId"),
				Value:    core.StringPtr(tResGrpID.(string)),
				Operator: core.StringPtr("stringEquals"),
			}
			policyResource.Attributes = append(policyResource.Attributes, resourceGroupResourceAttribute)
		}
	}

	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
		ServiceName:       &targetServiceName,
		SourceServiceName: &sourceServiceName,
		PolicyType:        &policyType,
	}
	roleList, resp, err := iampapClient.ListRoles(listRoleOptions)

	if err != nil || roleList == nil {
		return fmt.Errorf("[ERROR] Error in listing roles %s, %s", err, resp)
	}

	policyRoles := flex.MapRoleListToPolicyRoles(*roleList)
	roles, err := flex.GetRolesFromRoleNames(flex.ExpandStringList(d.Get("roles").([]interface{})), policyRoles)

	if err != nil {
		return err
	}

	policyGrant := &iampolicymanagementv1.Grant{
		Roles: flex.MapPolicyRolesToRoles(roles),
	}
	policyControl := &iampolicymanagementv1.Control{
		Grant: policyGrant,
	}

	createPolicyOptions := iampapClient.NewCreateV2PolicyOptions(
		policyControl,
		"authorization",
	)

	createPolicyOptions.SetSubject(policySubject)
	createPolicyOptions.SetResource(policyResource)

	if description, ok := d.GetOk("description"); ok {
		des := description.(string)
		createPolicyOptions.Description = &des
	}

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		createPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	authPolicy, resp, err := iampapClient.CreateV2Policy(createPolicyOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating authorization policy: %s %s", err, resp)
	}

	d.SetId(*authPolicy.ID)

	return resourceIBMIAMAuthorizationPolicyRead(d, meta)
}

func resourceIBMIAMAuthorizationPolicyRead(d *schema.ResourceData, meta interface{}) error {

	iampapClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	getPolicyOptions := &iampolicymanagementv1.GetV2PolicyOptions{
		ID: core.StringPtr(d.Id()),
	}

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		getPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	authorizationPolicy, resp, err := iampapClient.GetV2Policy(getPolicyOptions)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		authorizationPolicy, resp, err = iampapClient.GetV2Policy(getPolicyOptions)
		if err != nil || authorizationPolicy == nil {
			if resp != nil && resp.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		authorizationPolicy, resp, err = iampapClient.GetV2Policy(getPolicyOptions)
	}
	if err != nil || resp == nil {
		return fmt.Errorf("[ERROR] Error retrieving authorizationPolicy: %s %s", err, resp)
	}
	roles, err := flex.GetRoleNamesFromPolicyResponse(*authorizationPolicy, d, meta)
	if authorizationPolicy.Description != nil {
		d.Set("description", *authorizationPolicy.Description)
	}
	if len(resp.Headers["Transaction-Id"]) > 0 && resp.Headers["Transaction-Id"][0] != "" {
		d.Set("transaction_id", resp.Headers["Transaction-Id"][0])
	}
	d.Set("roles", roles)
	source := authorizationPolicy.Subject
	target := authorizationPolicy.Resource

	d.Set("resource_attributes", setAuthorizationResourceAttributes(*target))
	d.Set("target_resource_instance_id", flex.GetV2PolicyResourceAttribute("serviceInstance", *target))
	d.Set("target_resource_type", flex.GetV2PolicyResourceAttribute("resourceType", *target))
	d.Set("target_resource_group_id", flex.GetV2PolicyResourceAttribute("resourceGroupId", *target))
	d.Set("target_service_name", flex.GetV2PolicyResourceAttribute("serviceName", *target))
	if a, ok := d.GetOk("subject_attributes"); ok {
		d.Set("subject_attributes", setAuthorizationSubjectAttributes(*source, a.(*schema.Set)))
	}
	d.Set("source_service_name", flex.GetV2PolicySubjectAttribute("serviceName", *source))
	d.Set("source_resource_instance_id", flex.GetV2PolicySubjectAttribute("serviceInstance", *source))
	d.Set("source_resource_type", flex.GetV2PolicySubjectAttribute("resourceType", *source))
	d.Set("source_service_account", flex.GetV2PolicySubjectAttribute("accountId", *source))
	d.Set("source_resource_group_id", flex.GetV2PolicySubjectAttribute("resourceGroupId", *source))

	return nil
}

// Returns nil, because ibmcloud iam cli authorization policy does not have an update command
func resourceIBMIAMAuthorizationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("source_service_name") || d.HasChange("target_service_name") || d.HasChange("roles") || d.HasChange("source_resource_instance_id") || d.HasChange("target_resource_instance_id") ||
		d.HasChange("source_resource_group_id") || d.HasChange("source_resource_type") || d.HasChange("target_resource_type") || d.HasChange("source_service_account") ||
		d.HasChange("subject_attributes") || d.HasChange("resource_attributes") {
		var sourceServiceName, targetServiceName string
		policyType := "authorization"
		policySubject := &iampolicymanagementv1.V2PolicySubject{}
		policyResource := &iampolicymanagementv1.V2PolicyResource{}

		userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
		if err != nil {
			return err
		}

		iampapClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
		if err != nil {
			return err
		}

		// check subject_attributes exists
		if attributes, ok := d.GetOk("subject_attributes"); ok {
			for _, attribute := range attributes.(*schema.Set).List() {
				a := attribute.(map[string]interface{})
				name := a["name"].(string)
				value := a["value"].(string)
				operator := a["operator"].(string)
				if name == "serviceName" {
					sourceServiceName = value
				}

				if operator == "" && value == "*" && name == "resourceGroupId" {
					at := iampolicymanagementv1.V2PolicySubjectAttribute{
						Key:      &name,
						Value:    true,
						Operator: core.StringPtr("stringExists"),
					}
					policySubject.Attributes = append(policySubject.Attributes, at)
				} else if operator == "stringExists" {
					var resourceValue bool
					if value == "true" {
						resourceValue = true
					} else if value == "false" {
						resourceValue = false
					} else {
						return fmt.Errorf("[ERROR] Only values \"true\" and \"false\" are allowed when operator is \"stringExists\". Received %s.", value)
					}
					at := iampolicymanagementv1.V2PolicySubjectAttribute{
						Key:      &name,
						Value:    &resourceValue,
						Operator: &operator,
					}
					policySubject.Attributes = append(policySubject.Attributes, at)
				} else {
					if operator == "" {
						operator = "stringEquals"
					}
					at := iampolicymanagementv1.V2PolicySubjectAttribute{
						Key:      &name,
						Value:    &value,
						Operator: &operator,
					}
					policySubject.Attributes = append(policySubject.Attributes, at)
				}
			}
		} else {

			if name, ok := d.GetOk("source_service_name"); ok {
				sourceServiceName = name.(string)

				serviceNameSubjectAttribute := &iampolicymanagementv1.V2PolicySubjectAttribute{
					Key:      core.StringPtr("serviceName"),
					Value:    &sourceServiceName,
					Operator: core.StringPtr("stringEquals"),
				}
				policySubject.Attributes = append(policySubject.Attributes, *serviceNameSubjectAttribute)
			}

			sourceServiceAccount := userDetails.UserAccount
			if account, ok := d.GetOk("source_service_account"); ok {
				sourceServiceAccount = account.(string)
			}

			accountIdSubjectAttribute := &iampolicymanagementv1.V2PolicySubjectAttribute{
				Key:      core.StringPtr("accountId"),
				Value:    &sourceServiceAccount,
				Operator: core.StringPtr("stringEquals"),
			}

			policySubject.Attributes = append(policySubject.Attributes, *accountIdSubjectAttribute)

			if sID, ok := d.GetOk("source_resource_instance_id"); ok {
				serviceInstanceSubjectAttribute := iampolicymanagementv1.V2PolicySubjectAttribute{
					Key:      core.StringPtr("serviceInstance"),
					Value:    core.StringPtr(sID.(string)),
					Operator: core.StringPtr("stringEquals"),
				}
				policySubject.Attributes = append(policySubject.Attributes, serviceInstanceSubjectAttribute)
			}

			if sType, ok := d.GetOk("source_resource_type"); ok {
				resourceTypeSubjectAttribute := iampolicymanagementv1.V2PolicySubjectAttribute{
					Key:      core.StringPtr("resourceType"),
					Value:    core.StringPtr(sType.(string)),
					Operator: core.StringPtr("stringEquals"),
				}
				policySubject.Attributes = append(policySubject.Attributes, resourceTypeSubjectAttribute)
			}

			if sResGrpID, ok := d.GetOk("source_resource_group_id"); ok {
				if sResGrpID == "*" {
					resourceGroupSubjectAttribute := iampolicymanagementv1.V2PolicySubjectAttribute{
						Key:      core.StringPtr("resourceGroupId"),
						Value:    true,
						Operator: core.StringPtr("stringExists"),
					}
					policySubject.Attributes = append(policySubject.Attributes, resourceGroupSubjectAttribute)
				} else {
					resourceGroupSubjectAttribute := iampolicymanagementv1.V2PolicySubjectAttribute{
						Key:      core.StringPtr("resourceGroupId"),
						Value:    core.StringPtr(sResGrpID.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					policySubject.Attributes = append(policySubject.Attributes, resourceGroupSubjectAttribute)
				}
			}
		}

		// check resource_attributes exists
		if attributes, ok := d.GetOk("resource_attributes"); ok {
			for _, attribute := range attributes.(*schema.Set).List() {
				a := attribute.(map[string]interface{})
				name := a["name"].(string)
				value := a["value"].(string)
				operator := a["operator"].(string)
				if name == "serviceName" {
					targetServiceName = value
				}
				if name == "resourceType" && targetServiceName == "" {
					targetServiceName = "resource-controller"
				}
				if operator == "stringExists" {
					var resourceValue bool
					if value == "true" {
						resourceValue = true
					} else if value == "false" {
						resourceValue = false
					} else {
						return fmt.Errorf("[ERROR] When operator equals stringExists, value should be either \"true\" or \"false\", instead of %s", value)
					}
					at := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      &name,
						Value:    &resourceValue,
						Operator: &operator,
					}
					policyResource.Attributes = append(policyResource.Attributes, at)
				} else {
					at := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      &name,
						Value:    &value,
						Operator: &operator,
					}
					policyResource.Attributes = append(policyResource.Attributes, at)
				}
			}
		} else {
			if name, ok := d.GetOk("target_service_name"); ok {
				targetServiceName = name.(string)
				serviceNameResourceAttribute := &iampolicymanagementv1.V2PolicyResourceAttribute{
					Key:      core.StringPtr("serviceName"),
					Value:    core.StringPtr(targetServiceName),
					Operator: core.StringPtr("stringEquals"),
				}
				policyResource.Attributes = append(policyResource.Attributes, *serviceNameResourceAttribute)
			}

			accountIDResourceAttribute := &iampolicymanagementv1.V2PolicyResourceAttribute{
				Key:      core.StringPtr("accountId"),
				Value:    core.StringPtr(userDetails.UserAccount),
				Operator: core.StringPtr("stringEquals"),
			}

			policyResource.Attributes = append(policyResource.Attributes, *accountIDResourceAttribute)

			if tID, ok := d.GetOk("target_resource_instance_id"); ok {
				serviceInstanceResourceAttribute := iampolicymanagementv1.V2PolicyResourceAttribute{
					Key:      core.StringPtr("serviceInstance"),
					Value:    core.StringPtr(tID.(string)),
					Operator: core.StringPtr("stringEquals"),
				}
				policyResource.Attributes = append(policyResource.Attributes, serviceInstanceResourceAttribute)
			}

			if tType, ok := d.GetOk("target_resource_type"); ok {
				resourceTypeResourceAttribute := iampolicymanagementv1.V2PolicyResourceAttribute{
					Key:      core.StringPtr("resourceType"),
					Value:    core.StringPtr(tType.(string)),
					Operator: core.StringPtr("stringEquals"),
				}
				policyResource.Attributes = append(policyResource.Attributes, resourceTypeResourceAttribute)
				if targetServiceName == "" {
					targetServiceName = "resource-controller"
				}
			}

			if tResGrpID, ok := d.GetOk("target_resource_group_id"); ok {
				resourceGroupResourceAttribute := iampolicymanagementv1.V2PolicyResourceAttribute{
					Key:      core.StringPtr("resourceGroupId"),
					Value:    core.StringPtr(tResGrpID.(string)),
					Operator: core.StringPtr("stringEquals"),
				}
				policyResource.Attributes = append(policyResource.Attributes, resourceGroupResourceAttribute)
			}
		}

		listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
			ServiceName:       &targetServiceName,
			SourceServiceName: &sourceServiceName,
			PolicyType:        &policyType,
		}
		roleList, resp, err := iampapClient.ListRoles(listRoleOptions)

		if err != nil || roleList == nil {
			return fmt.Errorf("[ERROR] Error in listing roles %s, %s", err, resp)
		}

		policyRoles := flex.MapRoleListToPolicyRoles(*roleList)
		roles, err := flex.GetRolesFromRoleNames(flex.ExpandStringList(d.Get("roles").([]interface{})), policyRoles)

		if err != nil {
			return err
		}

		policyGrant := &iampolicymanagementv1.Grant{
			Roles: flex.MapPolicyRolesToRoles(roles),
		}
		policyControl := &iampolicymanagementv1.Control{
			Grant: policyGrant,
		}

		getPolicyOptions := &iampolicymanagementv1.GetV2PolicyOptions{
			ID: core.StringPtr(d.Id()),
		}

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			getPolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}

		authorizationPolicy, response, err := iampapClient.GetV2Policy(getPolicyOptions)

		if err != nil || authorizationPolicy == nil {
			if response != nil && response.StatusCode == 404 {
				return nil
			}
			return fmt.Errorf("[ERROR] Error retrieving Authorization Policy: %s\n%s", err, response)
		}
		policyETag := response.Headers.Get("ETag")

		updatePolicyOptions := iampapClient.NewReplaceV2PolicyOptions(
			d.Id(),
			policyETag,
			policyControl,
			"authorization",
		)
		updatePolicyOptions.SetSubject(policySubject)
		updatePolicyOptions.SetResource(policyResource)

		if description, ok := d.GetOk("description"); ok {
			des := description.(string)
			updatePolicyOptions.Description = &des
		}

		if transactionID, ok := d.GetOk("transaction_id"); ok {
			updatePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
		}
		_, resp, err = iampapClient.ReplaceV2Policy(updatePolicyOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating authorization policy: %s %s", err, resp)
		}
	}
	return resourceIBMIAMAuthorizationPolicyRead(d, meta)
}

func resourceIBMIAMAuthorizationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iampapClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	authorizationPolicyID := d.Id()

	deletePolicyOptions := &iampolicymanagementv1.DeletePolicyOptions{
		PolicyID: core.StringPtr(authorizationPolicyID),
	}

	if transactionID, ok := d.GetOk("transaction_id"); ok {
		deletePolicyOptions.SetHeaders(map[string]string{"Transaction-Id": transactionID.(string)})
	}

	resp, err := iampapClient.DeletePolicy(deletePolicyOptions)
	if err != nil {
		log.Printf(
			"Error deleting authorization policy: %s, %s", err, resp)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMAuthorizationPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iampapClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: core.StringPtr(d.Id()),
	}
	authorizationPolicy, resp, err := iampapClient.GetPolicy(getPolicyOptions)
	if err != nil || authorizationPolicy == nil {
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting authorisation policy: %s\n%s", err, resp)
	}

	if authorizationPolicy != nil && authorizationPolicy.State != nil && *authorizationPolicy.State == "deleted" {
		return false, nil
	}

	return *authorizationPolicy.ID == d.Id(), nil
}

func setAuthorizationResourceAttributes(list iampolicymanagementv1.V2PolicyResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, attribute := range list.Attributes {
		l := map[string]interface{}{
			"name":     attribute.Key,
			"value":    fmt.Sprintf("%v", attribute.Value),
			"operator": attribute.Operator,
		}
		result = append(result, l)
	}
	return result
}

func setAuthorizationSubjectAttributes(list iampolicymanagementv1.V2PolicySubject, a *schema.Set) []map[string]interface{} {
	previousOperators := make([]string, 0)

	for _, item := range a.List() {
		i := item.(map[string]interface{})

		previousOperators = append(previousOperators, i["operator"].(string))
	}

	result := make([]map[string]interface{}, 0)
	for i, attribute := range list.Attributes {
		var l map[string]interface{}
		if previousOperators[i] == "" && attribute.Value == true && *attribute.Operator == "stringExists" {
			l = map[string]interface{}{
				"name":  attribute.Key,
				"value": "*",
			}
		} else if previousOperators[i] == "" {
			l = map[string]interface{}{
				"name":  attribute.Key,
				"value": fmt.Sprintf("%v", attribute.Value),
			}
		} else {
			l = map[string]interface{}{
				"name":     attribute.Key,
				"value":    fmt.Sprintf("%v", attribute.Value),
				"operator": attribute.Operator,
			}
		}
		result = append(result, l)
	}
	return result
}
