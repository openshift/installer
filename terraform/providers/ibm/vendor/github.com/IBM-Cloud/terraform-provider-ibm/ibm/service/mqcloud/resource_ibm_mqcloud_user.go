// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package mqcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func ResourceIbmMqcloudUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmMqcloudUserCreate,
		ReadContext:   resourceIbmMqcloudUserRead,
		DeleteContext: resourceIbmMqcloudUserDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_user", "service_instance_guid"),
				Description:  "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_user", "name"),
				Description:  "The shortname of the user that will be used as the IBM MQ administrator in interactions with a queue manager for this service instance.",
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_user", "email"),
				Description:  "The email of the user.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for the user details.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the user which was allocated on creation, and can be used for delete calls.",
			},
		},
	}
}

func ResourceIbmMqcloudUserValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_instance_guid",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z][-a-z0-9]*$`,
			MinValueLength:             1,
			MaxValueLength:             12,
		},
		validate.ValidateSchema{
			Identifier:                 "email",
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			Required:                   true,
			MinValueLength:             5,
			MaxValueLength:             253,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_mqcloud_user", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmMqcloudUserCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Create User failed: %s", err.Error()), "ibm_mqcloud_user", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createUserOptions := &mqcloudv1.CreateUserOptions{}

	createUserOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createUserOptions.SetEmail(d.Get("email").(string))
	createUserOptions.SetName(d.Get("name").(string))

	userDetails, _, err := mqcloudClient.CreateUserWithContext(context, createUserOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateUserWithContext failed: %s", err.Error()), "ibm_mqcloud_user", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createUserOptions.ServiceInstanceGuid, *userDetails.ID))

	return resourceIbmMqcloudUserRead(context, d, meta)
}

func resourceIbmMqcloudUserRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getUserOptions := &mqcloudv1.GetUserOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "read", "sep-id-parts").GetDiag()
	}

	getUserOptions.SetServiceInstanceGuid(parts[0])
	getUserOptions.SetUserID(parts[1])

	userDetails, response, err := mqcloudClient.GetUserWithContext(context, getUserOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetUserWithContext failed: %s", err.Error()), "ibm_mqcloud_user", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("service_instance_guid", parts[0]); err != nil {
		err = fmt.Errorf("Error setting service_instance_guid: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "read", "set-service_instance_guid").GetDiag()
	}
	if err = d.Set("name", userDetails.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "read", "set-name").GetDiag()
	}
	if err = d.Set("email", userDetails.Email); err != nil {
		err = fmt.Errorf("Error setting email: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "read", "set-email").GetDiag()
	}
	if err = d.Set("href", userDetails.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "read", "set-href").GetDiag()
	}
	if err = d.Set("user_id", userDetails.ID); err != nil {
		err = fmt.Errorf("Error setting user_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "read", "set-user_id").GetDiag()
	}

	return nil
}

func resourceIbmMqcloudUserDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Delete User failed: %s", err.Error()), "ibm_mqcloud_user", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteUserOptions := &mqcloudv1.DeleteUserOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_user", "delete", "sep-id-parts").GetDiag()
	}

	deleteUserOptions.SetServiceInstanceGuid(parts[0])
	deleteUserOptions.SetUserID(parts[1])

	_, err = mqcloudClient.DeleteUserWithContext(context, deleteUserOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteUserWithContext failed: %s", err.Error()), "ibm_mqcloud_user", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
