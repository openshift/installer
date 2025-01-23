// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmSecretGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmSecretGroupCreate,
		ReadContext:   resourceIbmSmSecretGroupRead,
		UpdateContext: resourceIbmSmSecretGroupUpdate,
		DeleteContext: resourceIbmSmSecretGroupDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A v4 UUID identifier.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator(SecretGroupResourceName, "name"),
				Description:  "The name of your secret group.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator(SecretGroupResourceName, "description"),
				Description:  "An extended description of your secret group.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date that a resource was created. The date format follows RFC 3339.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date that a resource was recently modified. The date format follows RFC 3339.",
			},
		},
	}
}

func ResourceIbmSmSecretGroupValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[A-Za-z0-9][A-Za-z0-9]*(?:_*-*\\.*[A-Za-z0-9]+)*$`,
			MinValueLength:             2,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.*?)`,
			MinValueLength:             0,
			MaxValueLength:             1024,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: SecretGroupResourceName, Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSmSecretGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", SecretGroupResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createSecretGroupOptions := &secretsmanagerv2.CreateSecretGroupOptions{}

	createSecretGroupOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("description"); ok {
		createSecretGroupOptions.SetDescription(d.Get("description").(string))
	}

	secretGroup, response, err := secretsManagerClient.CreateSecretGroupWithContext(context, createSecretGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSecretGroupWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSecretGroupWithContext failed: %s\n%s", err.Error(), response), SecretGroupResourceName, "create")
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *secretGroup.ID))
	d.Set("secret_group_id", *secretGroup.ID)

	return resourceIbmSmSecretGroupRead(context, d, meta)
}

func resourceIbmSmSecretGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a secret group use the format `<region>/<instance_id>/<secret_group_id>`", SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}
	region := id[0]
	instanceId := id[1]
	secretGroupId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getSecretGroupOptions := &secretsmanagerv2.GetSecretGroupOptions{}

	getSecretGroupOptions.SetID(secretGroupId)

	secretGroup, response, err := secretsManagerClient.GetSecretGroupWithContext(context, getSecretGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSecretGroupWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretGroupWithContext failed %s\n%s", err, response), SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", secretGroupId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", secretGroup.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("description", secretGroup.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", DateTimeToRFC3339(secretGroup.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(secretGroup.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), SecretGroupResourceName, "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmSmSecretGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", SecretGroupResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	secretGroupId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	updateSecretGroupOptions := &secretsmanagerv2.UpdateSecretGroupOptions{}

	updateSecretGroupOptions.SetID(secretGroupId)

	hasChange := false

	patchVals := &secretsmanagerv2.SecretGroupPatch{}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("description") {
		newDescription := d.Get("description").(string)
		patchVals.Description = &newDescription
		hasChange = true
	}

	if hasChange {
		updateSecretGroupOptions.SecretGroupPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateSecretGroupWithContext(context, updateSecretGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretGroupWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretGroupWithContext failed %s\n%s", err, response), SecretGroupResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmSecretGroupRead(context, d, meta)
}

func resourceIbmSmSecretGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", SecretGroupResourceName, "delete")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	secretGroupId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	deleteSecretGroupOptions := &secretsmanagerv2.DeleteSecretGroupOptions{}

	deleteSecretGroupOptions.SetID(secretGroupId)

	response, err := secretsManagerClient.DeleteSecretGroupWithContext(context, deleteSecretGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteSecretGroupWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSecretGroupWithContext failed %s\n%s", err, response), SecretGroupResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
