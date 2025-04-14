// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	iamCRDisplayName = "display_name"
	iamCRName        = "name"
	iamCRDescription = "description"
	iamCRActions     = "actions"
	iamCRServiceName = "service"
)

func ResourceIBMIAMCustomRole() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMCustomRoleCreate,
		Read:     resourceIBMIAMCustomRoleRead,
		Update:   resourceIBMIAMCustomRoleUpdate,
		Delete:   resourceIBMIAMCustomRoleDelete,
		Exists:   resourceIBMIAMCustomRoleExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			iamCRDisplayName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Display Name of the Custom Role",
				ValidateFunc: validate.InvokeValidator("ibm_iam_custom_role", iamCRDisplayName),
			},

			iamCRName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the custom Role",
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_custom_role", iamCRName),
			},
			iamCRDescription: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The description of the role",
				ValidateFunc: validate.InvokeValidator("ibm_iam_custom_role", iamCRDescription),
			},
			iamCRServiceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Service Name",
				ForceNew:    true,
			},
			iamCRActions: {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The actions of the role",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "crn of the Custom Role",
			},
			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about the resource",
			},
		},
	}
}

func ResourceIBMIAMCustomRoleValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 iamCRName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Z]{1}[A-Za-z0-9]{0,29}$`,
			MinValueLength:             1,
			MaxValueLength:             30})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 iamCRDisplayName,
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			Optional:                   true,
			MinValueLength:             1,
			MaxValueLength:             50})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 iamCRDescription,
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			Optional:                   true,
			MinValueLength:             1,
			MaxValueLength:             250})

	ibmIAMCustomRoleResourceValidator := validate.ResourceValidator{ResourceName: "ibm_iam_custom_role", Schema: validateSchema}
	return &ibmIAMCustomRoleResourceValidator
}

func resourceIBMIAMCustomRoleCreate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	displayName := d.Get(iamCRDisplayName).(string)
	name := d.Get(iamCRName).(string)
	description := d.Get(iamCRDescription).(string)
	serviceName := d.Get(iamCRServiceName).(string)
	actionList := flex.ExpandStringList(d.Get(iamCRActions).([]interface{}))

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	roleOptions := &iampolicymanagementv1.CreateRoleOptions{
		DisplayName: &displayName,
		Actions:     actionList,
		Name:        &name,
		AccountID:   &userDetails.UserAccount,
		ServiceName: &serviceName,
		Description: &description,
	}

	role, response, err := iamPolicyManagementClient.CreateRole(roleOptions)
	if err != nil || role == nil {
		return fmt.Errorf("[ERROR] Error creating Custom Roles: %s\n%s", err, response)
	}

	d.SetId(*role.ID)

	return resourceIBMIAMCustomRoleRead(d, meta)
}

func resourceIBMIAMCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	roleID := d.Id()
	roleOptions := &iampolicymanagementv1.GetRoleOptions{
		RoleID: &roleID,
	}

	role, response, err := iamPolicyManagementClient.GetRole(roleOptions)
	if err != nil || role == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error retrieving Custom Roles: %s\n%s", err, response)
	}

	d.Set(iamCRDisplayName, role.DisplayName)
	d.Set(iamCRName, role.Name)
	d.Set(iamCRDescription, role.Description)
	d.Set(iamCRServiceName, role.ServiceName)
	d.Set(iamCRActions, role.Actions)
	d.Set("crn", role.CRN)

	d.Set(flex.ResourceName, role.Name)
	d.Set(flex.ResourceCRN, role.CRN)
	rcontroller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}

	d.Set(flex.ResourceControllerURL, rcontroller+"/iam/roles")

	return nil
}

func resourceIBMIAMCustomRoleUpdate(d *schema.ResourceData, meta interface{}) error {

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}
	roleID := d.Id()

	updatedDescription := d.Get(iamCRDescription).(string)
	updatedActions := flex.ExpandStringList(d.Get(iamCRActions).([]interface{}))
	updatedDisplayName := d.Get(iamCRDisplayName).(string)

	if d.HasChange("display_name") || d.HasChange("description") || d.HasChange("actions") {
		roleGetOptions := &iampolicymanagementv1.GetRoleOptions{
			RoleID: &roleID,
		}

		role, response, err := iamPolicyManagementClient.GetRole(roleGetOptions)
		if err != nil || role == nil {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[ERROR] Error retrieving Custom Roles: %s\n%s", err, response)
		}

		roleETag := response.Headers.Get("ETag")
		roleUpdateOptions := &iampolicymanagementv1.ReplaceRoleOptions{
			RoleID:      &roleID,
			IfMatch:     &roleETag,
			DisplayName: &updatedDisplayName,
			Description: &updatedDescription,
			Actions:     updatedActions,
		}

		_, response, err = iamPolicyManagementClient.ReplaceRole(roleUpdateOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating Custom Roles: %s\n%s", err, response)
		}
	}

	return resourceIBMIAMCustomRoleRead(d, meta)
}

func resourceIBMIAMCustomRoleDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	roleID := d.Id()
	roleDeleteOptions := &iampolicymanagementv1.DeleteRoleOptions{
		RoleID: &roleID,
	}

	response, err := iamPolicyManagementClient.DeleteRole(roleDeleteOptions)
	if err != nil && !strings.Contains(err.Error(), "404") {
		return fmt.Errorf("[ERROR] Error deleting Custom Roles: %s\n%s", err, response)
	}

	d.SetId("")

	return nil
}

func resourceIBMIAMCustomRoleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return false, err
	}
	roleID := d.Id()

	roleGetOptions := &iampolicymanagementv1.GetRoleOptions{
		RoleID: &roleID,
	}

	role, response, err := iamPolicyManagementClient.GetRole(roleGetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error retrieving Custom Roles: %s\n%s", err, response)
	}

	return *role.ID == roleID, nil
}
