// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMIAMAuthorizationPolicy() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMAuthorizationPolicyCreate,
		Read:     resourceIBMIAMAuthorizationPolicyRead,
		Update:   resourceIBMIAMAuthorizationPolicyUpdate,
		Delete:   resourceIBMIAMAuthorizationPolicyDelete,
		Exists:   resourceIBMIAMAuthorizationPolicyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"source_service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The source service name",
				ForceNew:    true,
			},

			"target_service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target service name",
			},

			"roles": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Role names of the policy definition",
			},

			"source_resource_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The source resource instance Id",
			},

			"target_resource_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The target resource instance Id",
			},

			"source_resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The source resource group Id",
			},

			"target_resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The target resource group Id",
			},

			"source_resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Resource type of source service",
			},

			"target_resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Resource type of target service",
			},

			"source_service_account": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Account GUID of source service",
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMIAMAuthorizationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	sourceServiceName := d.Get("source_service_name").(string)
	targetServiceName := d.Get("target_service_name").(string)

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	iampapClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	sourceServiceAccount := userDetails.userAccount

	if account, ok := d.GetOk("source_service_account"); ok {
		sourceServiceAccount = account.(string)
	}

	accountIdSubjectAttribute := &iampolicymanagementv1.SubjectAttribute{
		Name:  core.StringPtr("accountId"),
		Value: &sourceServiceAccount,
	}
	serviceNameSubjectAttribute := &iampolicymanagementv1.SubjectAttribute{
		Name:  core.StringPtr("serviceName"),
		Value: &sourceServiceName,
	}

	policySubject := &iampolicymanagementv1.PolicySubject{
		Attributes: []iampolicymanagementv1.SubjectAttribute{*accountIdSubjectAttribute, *serviceNameSubjectAttribute},
	}

	accountIDResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
		Name:     core.StringPtr("accountId"),
		Value:    core.StringPtr(userDetails.userAccount),
		Operator: core.StringPtr("stringEquals"),
	}

	serviceNameResourceAttribute := &iampolicymanagementv1.ResourceAttribute{
		Name:     core.StringPtr("serviceName"),
		Value:    core.StringPtr(targetServiceName),
		Operator: core.StringPtr("stringEquals"),
	}

	policyResource := &iampolicymanagementv1.PolicyResource{
		Attributes: []iampolicymanagementv1.ResourceAttribute{*accountIDResourceAttribute, *serviceNameResourceAttribute},
	}

	if sID, ok := d.GetOk("source_resource_instance_id"); ok {
		serviceInstanceSubjectAttribute := iampolicymanagementv1.SubjectAttribute{
			Name:  core.StringPtr("serviceInstance"),
			Value: core.StringPtr(sID.(string)),
		}
		policySubject.Attributes = append(policySubject.Attributes, serviceInstanceSubjectAttribute)
	}

	if tID, ok := d.GetOk("target_resource_instance_id"); ok {
		serviceInstanceResourceAttribute := iampolicymanagementv1.ResourceAttribute{
			Name:  core.StringPtr("serviceInstance"),
			Value: core.StringPtr(tID.(string)),
		}
		policyResource.Attributes = append(policyResource.Attributes, serviceInstanceResourceAttribute)
	}

	if sType, ok := d.GetOk("source_resource_type"); ok {
		resourceTypeSubjectAttribute := iampolicymanagementv1.SubjectAttribute{
			Name:  core.StringPtr("resourceType"),
			Value: core.StringPtr(sType.(string)),
		}
		policySubject.Attributes = append(policySubject.Attributes, resourceTypeSubjectAttribute)
	}

	if tType, ok := d.GetOk("target_resource_type"); ok {
		resourceTypeResourceAttribute := iampolicymanagementv1.ResourceAttribute{
			Name:  core.StringPtr("resourceType"),
			Value: core.StringPtr(tType.(string)),
		}
		policyResource.Attributes = append(policyResource.Attributes, resourceTypeResourceAttribute)
	}

	if sResGrpID, ok := d.GetOk("source_resource_group_id"); ok {
		resourceGroupSubjectAttribute := iampolicymanagementv1.SubjectAttribute{
			Name:  core.StringPtr("resourceGroupId"),
			Value: core.StringPtr(sResGrpID.(string)),
		}
		policySubject.Attributes = append(policySubject.Attributes, resourceGroupSubjectAttribute)
	}

	if tResGrpID, ok := d.GetOk("target_resource_group_id"); ok {
		resourceGroupResourceAttribute := iampolicymanagementv1.ResourceAttribute{
			Name:  core.StringPtr("resourceGroupId"),
			Value: core.StringPtr(tResGrpID.(string)),
		}
		policyResource.Attributes = append(policyResource.Attributes, resourceGroupResourceAttribute)
	}

	roles, err := getAuthorizationRolesByName(expandStringList(d.Get("roles").([]interface{})), sourceServiceName, targetServiceName, meta)
	if err != nil {
		return err
	}

	createPolicyOptions := iampapClient.NewCreatePolicyOptions(
		"authorization",
		[]iampolicymanagementv1.PolicySubject{*policySubject},
		roles,
		[]iampolicymanagementv1.PolicyResource{*policyResource},
	)
	authPolicy, _, err := iampapClient.CreatePolicy(createPolicyOptions)

	if err != nil {
		return fmt.Errorf("Error creating authorization policy: %s", err)
	}

	d.SetId(*authPolicy.ID)

	return resourceIBMIAMAuthorizationPolicyRead(d, meta)
}

func resourceIBMIAMAuthorizationPolicyRead(d *schema.ResourceData, meta interface{}) error {

	iampapClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: core.StringPtr(d.Id()),
	}

	authorizationPolicy, _, err := iampapClient.GetPolicy(getPolicyOptions)
	if err != nil {
		return fmt.Errorf("Error retrieving authorizationPolicy: %s", err)
	}
	roles := make([]string, len(authorizationPolicy.Roles))
	for i, role := range authorizationPolicy.Roles {
		roles[i] = *role.DisplayName
	}
	d.Set("roles", roles)
	source := authorizationPolicy.Subjects[0]
	target := authorizationPolicy.Resources[0]
	d.Set("source_service_name", getSubjectAttribute("serviceName", source))
	d.Set("target_service_name", getResourceAttribute("serviceName", target))
	d.Set("source_resource_instance_id", getSubjectAttribute("serviceInstance", source))
	d.Set("target_resource_instance_id", getResourceAttribute("serviceInstance", target))
	d.Set("source_resource_type", getSubjectAttribute("resourceType", source))
	d.Set("target_resource_type", getResourceAttribute("resourceType", target))
	d.Set("source_service_account", getSubjectAttribute("accountId", source))
	d.Set("source_resource_group_id", getSubjectAttribute("resourceGroupId", source))
	d.Set("target_resource_group_id", getResourceAttribute("resourceGroupId", target))
	return nil
}

// Returns nil, because ibmcloud iam cli authoirization policy does not have an update command
func resourceIBMIAMAuthorizationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMIAMAuthorizationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iampapClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	authorizationPolicyID := d.Id()

	deletePolicyOptions := &iampolicymanagementv1.DeletePolicyOptions{
		PolicyID: core.StringPtr(authorizationPolicyID),
	}
	_, err = iampapClient.DeletePolicy(deletePolicyOptions)
	if err != nil {
		log.Printf(
			"Error deleting authorization policy: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMAuthorizationPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iampapClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}

	getPolicyOptions := &iampolicymanagementv1.GetPolicyOptions{
		PolicyID: core.StringPtr(d.Id()),
	}
	authorizationPolicy, _, err := iampapClient.GetPolicy(getPolicyOptions)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return *authorizationPolicy.ID == d.Id(), nil
}

// TODO: Refactor to remove ListAuthorizationRoles (which lives in bluemix-go-sdk) ?
func getAuthorizationRolesByName(roleNames []string, sourceServiceName string, targetServiceName string, meta interface{}) ([]iampolicymanagementv1.PolicyRole, error) {

	iamClient, err := meta.(ClientSession).IAMAPI()
	if err != nil {
		return []iampolicymanagementv1.PolicyRole{}, err
	}

	iamRepo := iamClient.ServiceRoles()
	roles, err := iamRepo.ListAuthorizationRoles(sourceServiceName, targetServiceName)
	convertedRoles := convertRoleModels(roles)
	if err != nil {
		return []iampolicymanagementv1.PolicyRole{}, err
	}

	filteredRoles := []iampolicymanagementv1.PolicyRole{}
	filteredRoles, err = getRolesFromRoleNames(roleNames, convertedRoles)
	if err != nil {
		return []iampolicymanagementv1.PolicyRole{}, err
	}
	return filteredRoles, nil
}

// ConvertRoleModels will transform role models returned from "/v1/roles" to the model used by policy
func convertRoleModels(roles []models.PolicyRole) []iampolicymanagementv1.PolicyRole {
	results := make([]iampolicymanagementv1.PolicyRole, len(roles))
	for i, r := range roles {
		results[i] = iampolicymanagementv1.PolicyRole{
			RoleID:      core.StringPtr(r.ID.String()),
			DisplayName: core.StringPtr(r.DisplayName),
		}
	}
	return results
}
