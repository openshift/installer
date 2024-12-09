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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructure() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureCreate,
		ReadContext:   resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead,
		UpdateContext: resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureUpdate,
		DeleteContext: resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureDelete,
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
			"classic_infrastructure_username": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator(PublicCertConfigDnsClassicInfrastructureResourceName, "classic_infrastructure_username"),
				Description:  "The username that is associated with your classic infrastructure account.In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).",
			},
			"classic_infrastructure_password": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator(PublicCertConfigDnsClassicInfrastructureResourceName, "classic_infrastructure_password"),
				Description:  "Your classic infrastructure API key.For information about viewing and accessing your classic infrastructure API key, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).",
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

func ResourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "classic_infrastructure_username",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.*?)`,
			MinValueLength:             2,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "classic_infrastructure_password",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.*?)`,
			MinValueLength:             2,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: PublicCertConfigDnsClassicInfrastructureResourceName, Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsClassicInfrastructureResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))
	bodyModelMap := map[string]interface{}{}
	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	if _, ok := d.GetOk("config_type"); ok {
		bodyModelMap["config_type"] = d.Get("config_type")
	}
	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	if _, ok := d.GetOk("classic_infrastructure_username"); ok {
		bodyModelMap["classic_infrastructure_username"] = d.Get("classic_infrastructure_username")
	}
	if _, ok := d.GetOk("classic_infrastructure_password"); ok {
		bodyModelMap["classic_infrastructure_password"] = d.Get("classic_infrastructure_password")
	}
	convertedModel, err := resourceIbmSmConfigurationPublicCertificateClassicInfrastructureMapToPublicCertificateConfigurationDNSClassicInfrastructurePrototype(bodyModelMap)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsClassicInfrastructureResourceName, "create")
		return tfErr.GetDiag()
	}
	createConfigurationOptions.SetConfigurationPrototype(convertedModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationWithContext failed: %s\n%s", err.Error(), response), PublicCertConfigDnsClassicInfrastructureResourceName, "create")
		return tfErr.GetDiag()
	}

	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructure)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	return resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead(context, d, meta)
}

func resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a DNS configuration use the format `<region>/<instance_id>/<name>`", PublicCertConfigDnsClassicInfrastructureResourceName, "read")
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}

	configuration := configurationIntf.(*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructure)
	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}
	if !core.IsNil(configuration.ConfigType) {
		if err = d.Set("config_type", configuration.ConfigType); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if !core.IsNil(configuration.ClassicInfrastructureUsername) {
		if err = d.Set("classic_infrastructure_username", configuration.ClassicInfrastructureUsername); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting classic_infrastructure_username"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if !core.IsNil(configuration.ClassicInfrastructurePassword) {
		if err = d.Set("classic_infrastructure_password", configuration.ClassicInfrastructurePassword); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting classic_infrastructure_password"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("name", configuration.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", configuration.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", DateTimeToRFC3339(configuration.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(configuration.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), PublicCertConfigDnsClassicInfrastructureResourceName, "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsClassicInfrastructureResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))
	updateConfigurationOptions := &secretsmanagerv2.UpdateConfigurationOptions{}

	updateConfigurationOptions.SetName(configName)

	hasChange := false

	patchVals := &secretsmanagerv2.ConfigurationPatch{}
	if d.HasChange("classic_infrastructure_username") {
		newClassicInfrastructureUsername := d.Get("classic_infrastructure_username").(string)
		patchVals.ClassicInfrastructureUsername = &newClassicInfrastructureUsername
		hasChange = true
	}
	if d.HasChange("classic_infrastructure_password") {
		newClassicInfrastructurePassword := d.Get("classic_infrastructure_password").(string)
		patchVals.ClassicInfrastructurePassword = &newClassicInfrastructurePassword
		hasChange = true
	}

	if hasChange {
		updateConfigurationOptions.ConfigurationPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateConfigurationWithContext(context, updateConfigurationOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateConfigurationWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigDnsClassicInfrastructureResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureRead(context, d, meta)
}

func resourceIbmSmPublicCertificateConfigurationDNSClassicInfrastructureDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertConfigDnsClassicInfrastructureResourceName, "delete")
		return tfErr.GetDiag()
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteConfigurationWithContext failed %s\n%s", err, response), PublicCertConfigDnsClassicInfrastructureResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmConfigurationPublicCertificateClassicInfrastructureMapToPublicCertificateConfigurationDNSClassicInfrastructurePrototype(modelMap map[string]interface{}) (*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructurePrototype, error) {
	model := &secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructurePrototype{}
	model.ConfigType = core.StringPtr("public_cert_configuration_dns_classic_infrastructure")
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.ClassicInfrastructureUsername = core.StringPtr(modelMap["classic_infrastructure_username"].(string))
	model.ClassicInfrastructurePassword = core.StringPtr(modelMap["classic_infrastructure_password"].(string))
	return model, nil
}
