// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"

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
				ExactlyOneOf: []string{"source_service_name", "subject_attributes"},
				Description:  "The source service name",
				ForceNew:     true,
			},

			"target_service_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"target_service_name", "resource_attributes"},
				ForceNew:     true,
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
			},

			"target_resource_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"resource_attributes"},
				Description:   "The target resource group Id",
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
				ConflictsWith: []string{"source_resource_instance_id", "source_resource_group_id", "source_resource_type", "source_service_account"},
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
					},
				},
			},

			"resource_attributes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Description:   "Set resource attributes.",
				ConflictsWith: []string{"target_resource_instance_id", "target_resource_group_id", "target_resource_type"},
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
		},
	}
}

func resourceIBMIAMAuthorizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	var sourceServiceName, targetServiceName string
	policyType := "authorization"
	policySubject := &iampolicymanagementv1.PolicySubject{}
	policyResource := &iampolicymanagementv1.PolicyResource{}

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
			if name == "serviceName" {
				sourceServiceName = value
			}
			at := iampolicymanagementv1.SubjectAttribute{
				Name:  &name,
				Value: &value,
			}
			policySubject.Attributes = append(policySubject.Attributes, at)
		}
	} else {

		sourceServiceName = d.Get("source_service_name").(string)

		serviceNameSubjectAttribute := &iampolicymanagementv1.SubjectAttribute{
			Name:  core.StringPtr("serviceName"),
			Value: &sourceServiceName,
		}
		policySubject.Attributes = append(policySubject.Attributes, *serviceNameSubjectAttribute)

		sourceServiceAccount := userDetails.UserAccount
		if account, ok := d.GetOk("source_service_account"); ok {
			sourceServiceAccount = account.(string)
		}

		accountIdSubjectAttribute := &iampolicymanagementv1.SubjectAttribute{
			Name:  core.StringPtr("accountId"),
			Value: &sourceServiceAccount,
		}

		policySubject.Attributes = append(policySubject.Attributes, *accountIdSubjectAttribute)

		if sID, ok := d.GetOk("source_resource_instance_id"); ok {
			serviceInstanceSubjectAttribute := iampolicymanagementv1.SubjectAttribute{
				Name:  core.StringPtr("serviceInstance"),
				Value: core.StringPtr(sID.(string)),
			}
			policySubject.Attributes = append(policySubject.Attributes, serviceInstanceSubjectAttribute)
		}

		if sType, ok := d.GetOk("source_resource_type"); ok {
			resourceTypeSubjectAttribute := iampolicymanagementv1.SubjectAttribute{
				Name:  core.StringPtr("resourceType"),
				Value: core.StringPtr(sType.(string)),
			}
			policySubject.Attributes = append(policySubject.Attributes, resourceTypeSubjectAttribute)
		}

		if sResGrpID, ok := d.GetOk("source_resource_group_id"); ok {
			resourceGroupSubjectAttribute := iampolicymanagementv1.SubjectAttribute{
				Name:  core.StringPtr("resourceGroupId"),
				Value: core.StringPtr(sResGrpID.(string)),
			}
			policySubject.Attributes = append(policySubject.Attributes, resourceGroupSubjectAttribute)
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
			at := iampolicymanagementv1.ResourceAttribute{
				Name:     &name,
				Value:    &value,
				Operator: &operator,
			}
			policyResource.Attributes = append(policyResource.Attributes, at)
		}
	} else {
		targetServiceName = d.Get("target_service_name").(string)
		serviceNameResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
			Name:     core.StringPtr("serviceName"),
			Value:    core.StringPtr(targetServiceName),
			Operator: core.StringPtr("stringEquals"),
		}
		policyResource.Attributes = append(policyResource.Attributes, *serviceNameResourceAttribute)

		accountIDResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
			Name:     core.StringPtr("accountId"),
			Value:    core.StringPtr(userDetails.UserAccount),
			Operator: core.StringPtr("stringEquals"),
		}

		policyResource.Attributes = append(policyResource.Attributes, *accountIDResourceAttribute)

		if tID, ok := d.GetOk("target_resource_instance_id"); ok {
			serviceInstanceResourceAttribute := iampolicymanagementv1.ResourceAttribute{
				Name:  core.StringPtr("serviceInstance"),
				Value: core.StringPtr(tID.(string)),
			}
			policyResource.Attributes = append(policyResource.Attributes, serviceInstanceResourceAttribute)
		}

		if tType, ok := d.GetOk("target_resource_type"); ok {
			resourceTypeResourceAttribute := iampolicymanagementv1.ResourceAttribute{
				Name:  core.StringPtr("resourceType"),
				Value: core.StringPtr(tType.(string)),
			}
			policyResource.Attributes = append(policyResource.Attributes, resourceTypeResourceAttribute)
		}

		if tResGrpID, ok := d.GetOk("target_resource_group_id"); ok {
			resourceGroupResourceAttribute := iampolicymanagementv1.ResourceAttribute{
				Name:  core.StringPtr("resourceGroupId"),
				Value: core.StringPtr(tResGrpID.(string)),
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

	createPolicyOptions := iampapClient.NewCreatePolicyOptions(
		"authorization",
		[]iampolicymanagementv1.PolicySubject{*policySubject},
		roles,
		[]iampolicymanagementv1.PolicyResource{*policyResource},
	)

	if description, ok := d.GetOk("description"); ok {
		des := description.(string)
		createPolicyOptions.Description = &des
	}

	authPolicy, resp, err := iampapClient.CreatePolicy(createPolicyOptions)
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

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: core.StringPtr(d.Id()),
	}

	authorizationPolicy, resp, err := iampapClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving authorizationPolicy: %s %s", err, resp)
	}
	roles := make([]string, len(authorizationPolicy.Roles))
	for i, role := range authorizationPolicy.Roles {
		roles[i] = *role.DisplayName
	}
	if authorizationPolicy.Description != nil {
		d.Set("description", *authorizationPolicy.Description)
	}
	d.Set("roles", roles)
	source := authorizationPolicy.Subjects[0]
	target := authorizationPolicy.Resources[0]

	d.Set("resource_attributes", setAuthorizationResourceAttributes(target))
	d.Set("target_resource_instance_id", flex.GetResourceAttribute("serviceInstance", target))
	d.Set("target_resource_type", flex.GetResourceAttribute("resourceType", target))
	d.Set("target_resource_group_id", flex.GetResourceAttribute("resourceGroupId", target))
	d.Set("target_service_name", flex.GetResourceAttribute("serviceName", target))

	d.Set("subject_attributes", setAuthorizationSubjectAttributes(source))
	d.Set("source_service_name", flex.GetSubjectAttribute("serviceName", source))
	d.Set("source_resource_instance_id", flex.GetSubjectAttribute("serviceInstance", source))
	d.Set("source_resource_type", flex.GetSubjectAttribute("resourceType", source))
	d.Set("source_service_account", flex.GetSubjectAttribute("accountId", source))
	d.Set("source_resource_group_id", flex.GetSubjectAttribute("resourceGroupId", source))

	return nil
}

// Returns nil, because ibmcloud iam cli authorization policy does not have an update command
func resourceIBMIAMAuthorizationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
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

func setAuthorizationResourceAttributes(list iampolicymanagementv1.PolicyResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, attribute := range list.Attributes {
		l := map[string]interface{}{
			"name":     attribute.Name,
			"value":    attribute.Value,
			"operator": attribute.Operator,
		}
		result = append(result, l)
	}
	return result
}

func setAuthorizationSubjectAttributes(list iampolicymanagementv1.PolicySubject) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, attribute := range list.Attributes {
		l := map[string]interface{}{
			"name":  attribute.Name,
			"value": attribute.Value,
		}
		result = append(result, l)
	}
	return result
}
