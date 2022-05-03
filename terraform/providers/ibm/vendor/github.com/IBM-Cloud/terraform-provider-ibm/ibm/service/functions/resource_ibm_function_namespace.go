// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package functions

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/api/functions"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	funcNamespaceName      = "name"
	funcNamespaceResGrpId  = "resource_group_id"
	funcNamespaceResPlanId = "resource_plan_id"
	funcNamespaceDesc      = "description"
	funcNamespaceLoc       = "location"
)

func ResourceIBMFunctionNamespace() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMFunctionNamespaceCreate,
		Read:     resourceIBMFunctionNamespaceRead,
		Update:   resourceIBMFunctionNamespaceUpdate,
		Delete:   resourceIBMFunctionNamespaceDelete,
		Exists:   resourceIBMFunctionNamespaceExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			funcNamespaceName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of namespace.",
				ValidateFunc: validate.InvokeValidator("ibm_function_namespace", funcNamespaceName),
			},
			funcNamespaceDesc: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace Description.",
			},
			funcNamespaceResGrpId: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Resource Group ID.",
				ValidateFunc: validate.InvokeValidator("ibm_function_namespace", funcNamespaceResGrpId),
			},
			funcNamespaceLoc: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Namespace Location.",
			},
		},
	}
}

func ResourceIBMFuncNamespaceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 funcNamespaceName,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
			Required:                   true})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 funcNamespaceResGrpId,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
			Required:                   true})

	ibmFuncNamespaceResourceValidator := validate.ResourceValidator{ResourceName: "ibm_function_namespace", Schema: validateSchema}
	return &ibmFuncNamespaceResourceValidator
}

func resourceIBMFunctionNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	createNamespaceOptions := functions.CreateNamespaceOptions{}

	name := d.Get(funcNamespaceName).(string)
	createNamespaceOptions.Name = &name
	resourceGroupID := d.Get(funcNamespaceResGrpId).(string)
	createNamespaceOptions.ResourceGroupID = &resourceGroupID
	resourcePlanID := "functions-base-plan"
	createNamespaceOptions.ResourcePlanID = &resourcePlanID

	if _, ok := d.GetOk(funcNamespaceDesc); ok {
		description := d.Get(funcNamespaceDesc).(string)
		createNamespaceOptions.Description = &description
	}

	namespace, err := functionNamespaceAPI.Namespaces().CreateNamespace(createNamespaceOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating Namespace: %s", err)
	}

	d.SetId(*namespace.ID)
	log.Printf("[INFO] Created namespace (IAM) : %s", *namespace.Name)

	return resourceIBMFunctionNamespaceRead(d, meta)
}

func resourceIBMFunctionNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	id := d.Id()

	getOptions := functions.GetNamespaceOptions{
		ID: &id,
	}
	instance, err := functionNamespaceAPI.Namespaces().GetNamespace(getOptions)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting namesapce (IAM): %s", err)
	}

	if instance.Name != nil {
		d.Set(funcNamespaceName, *instance.Name)
	}

	if instance.ResourceGroupID != nil {
		d.Set(funcNamespaceResGrpId, *instance.ResourceGroupID)
	}

	if instance.Location != nil {
		d.Set(funcNamespaceLoc, *instance.Location)
	}
	if instance.Description != nil {
		d.Set(funcNamespaceDesc, *instance.Description)
	}

	return nil
}

func resourceIBMFunctionNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	nsClient, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	ID := d.Id()
	updateNamespaceOptions := functions.UpdateNamespaceOptions{}
	if d.HasChange(funcNamespaceName) {
		name := d.Get(funcNamespaceName).(string)
		updateNamespaceOptions.Name = &name
	}

	if d.HasChange(funcNamespaceDesc) {
		description := d.Get(funcNamespaceDesc).(string)
		updateNamespaceOptions.Description = &description
	}

	updateNamespaceOptions.ID = &ID
	namespace, err := nsClient.Namespaces().UpdateNamespace(updateNamespaceOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Updating Namespace: %s", err)
	}

	log.Printf("[INFO] Updated namespace (IAM) : %s", *namespace.Name)

	return resourceIBMFunctionNamespaceRead(d, meta)
}

func resourceIBMFunctionNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	nsClient, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	ID := d.Id()
	_, err = nsClient.Namespaces().DeleteNamespace(ID)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Namespace: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMFunctionNamespaceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	nsClient, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return false, err
	}

	ID := d.Id()

	getOptions := functions.GetNamespaceOptions{
		ID: &ID,
	}
	_, err = nsClient.Namespaces().GetNamespace(getOptions)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting existing namesapce (IAM): %s", err)
	}

	return true, nil

}
