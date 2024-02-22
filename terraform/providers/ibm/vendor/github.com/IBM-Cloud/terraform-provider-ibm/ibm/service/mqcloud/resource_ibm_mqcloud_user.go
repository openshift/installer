// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

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
				Description:  "The GUID that uniquely identifies the MQ on Cloud service instance.",
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
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Create User failed %s", err))
	}
	createUserOptions := &mqcloudv1.CreateUserOptions{}

	createUserOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createUserOptions.SetEmail(d.Get("email").(string))
	createUserOptions.SetName(d.Get("name").(string))

	userDetails, response, err := mqcloudClient.CreateUserWithContext(context, createUserOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateUserWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateUserWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createUserOptions.ServiceInstanceGuid, *userDetails.ID))

	return resourceIbmMqcloudUserRead(context, d, meta)
}

func resourceIbmMqcloudUserRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getUserOptions := &mqcloudv1.GetUserOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getUserOptions.SetServiceInstanceGuid(parts[0])
	getUserOptions.SetUserID(parts[1])

	userDetails, response, err := mqcloudClient.GetUserWithContext(context, getUserOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetUserWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetUserWithContext failed %s\n%s", err, response))
	}
	if err = d.Set("service_instance_guid", parts[0]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_instance_guid: %s", err))
	}
	if err = d.Set("name", userDetails.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("email", userDetails.Email); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting email: %s", err))
	}
	if err = d.Set("href", userDetails.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("user_id", userDetails.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting user_id: %s", err))
	}

	return nil
}

func resourceIbmMqcloudUserDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		return diag.FromErr(err)
	}
	err = checkSIPlan(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Delete User failed %s", err))
	}
	deleteUserOptions := &mqcloudv1.DeleteUserOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteUserOptions.SetServiceInstanceGuid(parts[0])
	deleteUserOptions.SetUserID(parts[1])

	response, err := mqcloudClient.DeleteUserWithContext(context, deleteUserOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteUserWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteUserWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
