// Copyright IBM Corp. 2023 All Rights Reserved.
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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmIamCredentialsConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmIamCredentialsConfigurationCreate,
		ReadContext:   resourceIbmSmIamCredentialsConfigurationRead,
		UpdateContext: resourceIbmSmIamCredentialsConfigurationUpdate,
		DeleteContext: resourceIbmSmIamCredentialsConfigurationDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-readable unique name to assign to your configuration.To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.",
			},
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration type.",
			},
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "An IBM Cloud API key that can create and manage service IDs. The API key must be assigned the Editor platform role on the Access Groups Service and the Operator platform role on the IAM Identity Service. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-configure-iam-engine).",
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "This attribute indicates whether the API key configuration is disabled. If it is set to `true`, the IAM credentials engine doesn't use the configured API key for credentials management.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier that is associated with the entity that created the secret.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was created. The date format follows RFC 3339.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
			},
		},
	}
}

func resourceIbmSmIamCredentialsConfigurationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", IAMCredentialsConfigResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	configurationPrototypeModel, err := resourceIbmSmIamCredentialsConfigurationMapToConfigurationPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", IAMCredentialsConfigResourceName, "create")
		return tfErr.GetDiag()
	}
	createConfigurationOptions.SetConfigurationPrototype(configurationPrototypeModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationWithContext failed %s\n%s", err, response), IAMCredentialsConfigResourceName, "create")
		return tfErr.GetDiag()
	}
	configuration := configurationIntf.(*secretsmanagerv2.IAMCredentialsConfiguration)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	return resourceIbmSmIamCredentialsConfigurationRead(context, d, meta)
}

func resourceIbmSmIamCredentialsConfigurationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import IAM credentials configuration use the format `<region>/<instance_id>/<name>`", IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(configName)

	configurationIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	configuration := configurationIntf.(*secretsmanagerv2.IAMCredentialsConfiguration)

	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", configuration.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", configuration.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", DateTimeToRFC3339(configuration.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(configuration.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("api_key", configuration.ApiKey); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting api_key"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("disabled", configuration.Disabled); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting disabled"), IAMCredentialsConfigResourceName, "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmSmIamCredentialsConfigurationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", IAMCredentialsConfigResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	updateConfigurationOptions := &secretsmanagerv2.UpdateConfigurationOptions{}

	updateConfigurationOptions.SetName(configName)
	updateConfigurationOptions.SetXSmAcceptConfigurationType("iam_credentials_configuration")

	hasChange := false

	patchVals := &secretsmanagerv2.ConfigurationPatch{}

	if d.HasChange("api_key") {
		patchVals.ApiKey = core.StringPtr(d.Get("api_key").(string))
		hasChange = true
	}

	if d.HasChange("disabled") {
		patchVals.Disabled = core.BoolPtr(d.Get("disabled").(bool))
		hasChange = true
	}

	if hasChange {
		updateConfigurationOptions.ConfigurationPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateConfigurationWithContext(context, updateConfigurationOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateConfigurationWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateConfigurationWithContext failed %s\n%s", err, response), IAMCredentialsConfigResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmIamCredentialsConfigurationRead(context, d, meta)
}

func resourceIbmSmIamCredentialsConfigurationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	deleteConfigurationOptions := &secretsmanagerv2.DeleteConfigurationOptions{}

	deleteConfigurationOptions.SetName(configName)

	response, err := secretsManagerClient.DeleteConfigurationWithContext(context, deleteConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteConfigurationWithContext failed %s\n%s", err, response), IAMCredentialsConfigResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmIamCredentialsConfigurationMapToConfigurationPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationPrototypeIntf, error) {
	model := &secretsmanagerv2.IAMCredentialsConfigurationPrototype{}

	model.ConfigType = core.StringPtr("iam_credentials_configuration")

	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("api_key"); ok {
		model.ApiKey = core.StringPtr(d.Get("api_key").(string))
	}
	if _, ok := d.GetOk("disabled"); ok {
		model.Disabled = core.BoolPtr(d.Get("disabled").(bool))
	}
	return model, nil
}
